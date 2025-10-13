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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	"github.com/thomas-maurice/glua/pkg/stubgen"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	outputDir := flag.String("output", "library", "Output directory for generated stubs")
	flag.Parse()

	// Get the directory where this source file lives
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error determining source directory\n")
		os.Exit(1)
	}
	moduleDir := filepath.Dir(filepath.Dir(filename))

	// Create generator and generate stubs
	gen := stubgen.NewGenerator()
	outputFile, err := gen.Generate(stubgen.GenerateConfig{
		ScanDir:    moduleDir,
		OutputDir:  *outputDir,
		ModuleName: "kubernetes",
		OutputFile: "kubernetes.gen.lua",
		Types: []interface{}{
			kubernetes.GVKMatcher{},
			// Core resources
			corev1.Pod{},
			corev1.PodList{},
			corev1.Namespace{},
			corev1.NamespaceList{},
			corev1.Node{},
			corev1.NodeList{},
			corev1.ConfigMap{},
			corev1.ConfigMapList{},
			corev1.Secret{},
			corev1.SecretList{},
			corev1.Service{},
			corev1.ServiceList{},
			corev1.ServiceAccount{},
			corev1.ServiceAccountList{},
			corev1.PersistentVolume{},
			corev1.PersistentVolumeList{},
			corev1.PersistentVolumeClaim{},
			corev1.PersistentVolumeClaimList{},
			// Apps resources
			appsv1.Deployment{},
			appsv1.DeploymentList{},
			appsv1.StatefulSet{},
			appsv1.StatefulSetList{},
			appsv1.DaemonSet{},
			appsv1.DaemonSetList{},
			appsv1.ReplicaSet{},
			appsv1.ReplicaSetList{},
			// Batch resources
			batchv1.Job{},
			batchv1.JobList{},
			batchv1.CronJob{},
			batchv1.CronJobList{},
			// Networking resources
			networkingv1.Ingress{},
			networkingv1.IngressList{},
			networkingv1.NetworkPolicy{},
			networkingv1.NetworkPolicyList{},
			// RBAC resources
			rbacv1.Role{},
			rbacv1.RoleList{},
			rbacv1.ClusterRole{},
			rbacv1.ClusterRoleList{},
			rbacv1.RoleBinding{},
			rbacv1.RoleBindingList{},
			rbacv1.ClusterRoleBinding{},
			rbacv1.ClusterRoleBindingList{},
			// Metav1 types
			metav1.ObjectMeta{},
			metav1.TypeMeta{},
			metav1.Time{},
			metav1.MicroTime{},
			metav1.Duration{},
			metav1.Status{},
			metav1.StatusDetails{},
			metav1.StatusCause{},
			metav1.ListMeta{},
			metav1.OwnerReference{},
			metav1.LabelSelector{},
			metav1.LabelSelectorRequirement{},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outputFile)
}
