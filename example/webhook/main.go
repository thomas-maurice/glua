package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

type Config struct {
	Address    string
	CertFile   string
	KeyFile    string
	ScriptPath string
}

type WebhookServer struct {
	config *Config
	logger *slog.Logger
	engine *gin.Engine
}

// NewWebhookServer: creates a new webhook server instance
func NewWebhookServer(cfg *Config) *WebhookServer {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	ws := &WebhookServer{
		config: cfg,
		logger: logger,
		engine: engine,
	}

	// Register routes
	engine.POST("/mutate", ws.handleMutate)
	engine.GET("/healthz", ws.handleHealth)

	return ws
}

// handleHealth: health check endpoint
func (ws *WebhookServer) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleMutate: processes admission review requests and applies Lua-based mutations
func (ws *WebhookServer) handleMutate(c *gin.Context) {
	var admissionReview admissionv1.AdmissionReview

	if err := c.ShouldBindJSON(&admissionReview); err != nil {
		ws.logger.Error("failed to decode admission review", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if admissionReview.Request == nil {
		ws.logger.Error("admission review request is nil")
		c.JSON(http.StatusBadRequest, gin.H{"error": "admission review request is nil"})
		return
	}

	ws.logger.Info("received admission request",
		"uid", admissionReview.Request.UID,
		"kind", admissionReview.Request.Kind.Kind,
		"namespace", admissionReview.Request.Namespace,
		"name", admissionReview.Request.Name,
	)

	response := ws.mutate(admissionReview.Request)
	admissionReview.Response = response

	c.JSON(http.StatusOK, admissionReview)
}

// mutate: processes the admission request and returns a response with patches
func (ws *WebhookServer) mutate(req *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	response := &admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
	}

	// Only handle Pod resources
	if req.Kind.Kind != "Pod" {
		return response
	}

	// Decode the pod
	pod := &corev1.Pod{}
	if err := json.Unmarshal(req.Object.Raw, pod); err != nil {
		ws.logger.Error("failed to unmarshal pod", "error", err)
		return response
	}

	// Check if pod should be mutated
	if pod.Labels["thomas.maurice/mutate"] != "true" {
		ws.logger.Info("pod doesn't have mutation label, skipping",
			"namespace", pod.Namespace,
			"name", pod.Name,
		)
		return response
	}

	ws.logger.Info("mutating pod",
		"namespace", pod.Namespace,
		"name", pod.Name,
	)

	// Run Lua mutation script
	patches, err := ws.runLuaMutation(pod)
	if err != nil {
		ws.logger.Error("lua mutation failed", "error", err)
		return response
	}

	if len(patches) == 0 {
		ws.logger.Info("no patches generated")
		return response
	}

	// Marshal patches to JSON
	patchBytes, err := json.Marshal(patches)
	if err != nil {
		ws.logger.Error("failed to marshal patches", "error", err)
		return response
	}

	patchType := admissionv1.PatchTypeJSONPatch
	response.Patch = patchBytes
	response.PatchType = &patchType

	ws.logger.Info("mutation successful",
		"patches", len(patches),
	)

	return response
}

// runLuaMutation: executes the Lua script to generate JSON patches
func (ws *WebhookServer) runLuaMutation(pod *corev1.Pod) ([]map[string]interface{}, error) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	// Convert pod to Lua table
	podTable, err := translator.ToLua(L, pod)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pod to lua: %w", err)
	}

	// Set pod as global variable
	L.SetGlobal("pod", podTable)

	// Create empty patches table
	L.SetGlobal("patches", L.NewTable())

	// Execute the Lua script
	if err := L.DoFile(ws.config.ScriptPath); err != nil {
		return nil, fmt.Errorf("failed to execute lua script: %w", err)
	}

	// Get patches from Lua
	patchesValue := L.GetGlobal("patches")
	if patchesValue.Type() == lua.LTNil {
		return []map[string]interface{}{}, nil
	}

	// Convert Lua table to Go slice
	var patches []map[string]interface{}
	if err := translator.FromLua(L, patchesValue, &patches); err != nil {
		return nil, fmt.Errorf("failed to convert patches from lua: %w", err)
	}

	return patches, nil
}

// Serve: starts the webhook server with TLS
func (ws *WebhookServer) Serve() error {
	ws.logger.Info("starting webhook server",
		"address", ws.config.Address,
		"script", ws.config.ScriptPath,
	)

	return ws.engine.RunTLS(ws.config.Address, ws.config.CertFile, ws.config.KeyFile)
}

func main() {
	var (
		address    = flag.String("address", ":8443", "Address to listen on")
		certFile   = flag.String("cert", "/etc/webhook/certs/tls.crt", "TLS certificate file")
		keyFile    = flag.String("key", "/etc/webhook/certs/tls.key", "TLS key file")
		scriptPath = flag.String("script", "/etc/webhook/scripts/mutate.lua", "Lua mutation script path")
	)
	flag.Parse()

	// Check if script exists
	if _, err := os.Stat(*scriptPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Script file not found: %s\n", *scriptPath)
		os.Exit(1)
	}

	cfg := &Config{
		Address:    *address,
		CertFile:   *certFile,
		KeyFile:    *keyFile,
		ScriptPath: *scriptPath,
	}

	server := NewWebhookServer(cfg)
	if err := server.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
