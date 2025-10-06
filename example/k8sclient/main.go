package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thomas-maurice/glua/pkg/modules/k8sclient"
	lua "github.com/yuin/gopher-lua"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var (
		kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file (defaults to $HOME/.kube/config)")
		script     = flag.String("script", "test.lua", "Lua script to execute")
	)
	flag.Parse()

	// Determine kubeconfig path
	kubeconfigPath := *kubeconfig
	if kubeconfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		kubeconfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create Lua state
	L := lua.NewState()
	defer L.Close()

	// Preload k8sclient module
	L.PreloadModule("k8sclient", k8sclient.Loader(config))

	// Execute script
	if err := L.DoFile(*script); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing script: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Script executed successfully")
}
