package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thomas-maurice/glua/pkg/modules/k8sclient"
	lua "github.com/yuin/gopher-lua"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: kubernetes <script.lua>")
	}

	scriptPath := os.Args[1]

	// Load kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create Lua state
	L := lua.NewState()
	defer L.Close()

	// Register k8sclient module
	L.PreloadModule("k8sclient", k8sclient.Loader(config))

	// Run script
	if err := L.DoFile(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
