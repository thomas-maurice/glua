// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	lua "github.com/yuin/gopher-lua"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/runtime"
)

// Config: holds the webhook server configuration
type Config struct {
	Address     string
	CertFile    string
	KeyFile     string
	ScriptsDir  string
	EnableNodes bool
	EnablePods  bool
}

// WebhookServer: represents the generic mutating webhook server
type WebhookServer struct {
	config *Config
	logger *slog.Logger
	engine *gin.Engine
}

// NewWebhookServer: creates a new generic webhook server instance
func NewWebhookServer(cfg *Config) (*WebhookServer, error) {
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

	logger.Info("webhook server initialized",
		"enablePods", cfg.EnablePods,
		"enableNodes", cfg.EnableNodes,
		"scriptsDir", cfg.ScriptsDir,
	)

	return ws, nil
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

	// Route to appropriate handler based on resource kind
	switch req.Kind.Kind {
	case "Pod":
		if !ws.config.EnablePods {
			ws.logger.Info("pod mutations disabled")
			return response
		}
		return ws.mutatePod(req, response)
	case "Node":
		if !ws.config.EnableNodes {
			ws.logger.Info("node mutations disabled")
			return response
		}
		return ws.mutateNode(req, response)
	default:
		ws.logger.Info("unsupported resource kind", "kind", req.Kind.Kind)
		return response
	}
}

// mutatePod: handles Pod resource mutations
func (ws *WebhookServer) mutatePod(req *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse) *admissionv1.AdmissionResponse {
	pod := &corev1.Pod{}
	if err := json.Unmarshal(req.Object.Raw, pod); err != nil {
		ws.logger.Error("failed to unmarshal pod", "error", err)
		return response
	}

	ws.logger.Info("mutating pod",
		"namespace", pod.Namespace,
		"name", pod.Name,
	)

	scriptPath := filepath.Join(ws.config.ScriptsDir, "mutate_pod.lua")
	patches, err := ws.runLuaMutation(scriptPath, "pod", pod)
	if err != nil {
		ws.logger.Error("pod mutation failed", "error", err)
		return response
	}

	return ws.buildPatchResponse(response, patches)
}

// mutateNode: handles Node resource mutations
func (ws *WebhookServer) mutateNode(req *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse) *admissionv1.AdmissionResponse {
	node := &corev1.Node{}
	if err := json.Unmarshal(req.Object.Raw, node); err != nil {
		ws.logger.Error("failed to unmarshal node", "error", err)
		return response
	}

	ws.logger.Info("mutating node",
		"name", node.Name,
	)

	scriptPath := filepath.Join(ws.config.ScriptsDir, "mutate_node.lua")
	patches, err := ws.runLuaMutation(scriptPath, "node", node)
	if err != nil {
		ws.logger.Error("node mutation failed", "error", err)
		return response
	}

	return ws.buildPatchResponse(response, patches)
}

// runLuaMutation: executes the Lua script to generate JSON patches (generic implementation)
func (ws *WebhookServer) runLuaMutation(scriptPath, globalName string, object interface{}) ([]map[string]interface{}, error) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	// Preload kubernetes module for Lua scripts
	L.PreloadModule("kubernetes", kubernetes.Loader)

	// Convert object to Lua table
	objectTable, err := translator.ToLua(L, object)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to lua: %w", globalName, err)
	}

	// Set object as global variable
	L.SetGlobal(globalName, objectTable)

	// Create empty patches table
	L.SetGlobal("patches", L.NewTable())

	// Execute the Lua script
	if err := L.DoFile(scriptPath); err != nil {
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

// buildPatchResponse: constructs the admission response with JSON patches
func (ws *WebhookServer) buildPatchResponse(response *admissionv1.AdmissionResponse, patches []map[string]interface{}) *admissionv1.AdmissionResponse {
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

// Run: starts the webhook server
func (ws *WebhookServer) Run() error {
	ws.logger.Info("starting webhook server",
		"address", ws.config.Address,
	)

	if ws.config.CertFile != "" && ws.config.KeyFile != "" {
		return ws.engine.RunTLS(ws.config.Address, ws.config.CertFile, ws.config.KeyFile)
	}

	ws.logger.Warn("running without TLS - not suitable for production")
	return ws.engine.Run(ws.config.Address)
}

func main() {
	var (
		address     = flag.String("address", ":8443", "Address to listen on")
		certFile    = flag.String("cert", "/etc/webhook/certs/tls.crt", "Path to TLS certificate")
		keyFile     = flag.String("key", "/etc/webhook/certs/tls.key", "Path to TLS private key")
		scriptsDir  = flag.String("scripts", "/etc/webhook/scripts", "Path to Lua scripts directory")
		enableNodes = flag.Bool("enable-nodes", true, "Enable node mutations")
		enablePods  = flag.Bool("enable-pods", true, "Enable pod mutations")
	)
	flag.Parse()

	config := &Config{
		Address:     *address,
		CertFile:    *certFile,
		KeyFile:     *keyFile,
		ScriptsDir:  *scriptsDir,
		EnableNodes: *enableNodes,
		EnablePods:  *enablePods,
	}

	server, err := NewWebhookServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create webhook server: %v\n", err)
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// Ensure we satisfy the runtime.Object interface (unused but good practice)
var (
	_ metav1.Object = &corev1.Pod{}
	_ metav1.Object = &corev1.Node{}
)
