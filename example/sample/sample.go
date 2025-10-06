package sample

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
)

func GetPod() *corev1.Pod {
	var pod corev1.Pod

	err := json.Unmarshal([]byte(podJson), &pod)
	if err != nil {
		panic(err)
	}

	return &pod
}

var podJson = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2025-10-03T16:39:00Z",
        "generateName": "sample-application-6c998859b7-",
        "generation": 1,
        "labels": {
            "app": "sample-application",
            "team": "platform",
            "environment": "production",
            "mutate": "true",
            "pod-template-hash": "6c998859b7"
        },
        "name": "sample-application-6c998859b7-6bxxr",
        "namespace": "default",
        "ownerReferences": [
            {
                "apiVersion": "apps/v1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "sample-application-6c998859b7",
                "uid": "c37260e1-7929-45f9-8fca-5e4bf6a75861"
            }
        ],
        "resourceVersion": "499",
        "uid": "e2c3f346-dd4c-4a5d-a40a-04b5395776ec"
    },
    "spec": {
        "containers": [
            {
                "args": [
                    "-c",
                    "sleep 3600"
                ],
                "command": [
                    "/bin/sh"
                ],
                "env": [
                    {
                        "name": "GOMAXPROCS",
                        "value": "2"
                    },
                    {
                        "name": "some_env",
                        "value": "foo"
                    }
                ],
                "image": "alpine",
                "imagePullPolicy": "Always",
                "name": "sample-application",
                "resources": {
                    "limits": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    },
                    "requests": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "kube-api-access-87h5h",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "nodeName": "kind-control-plane",
        "preemptionPolicy": "PreemptLowerPriority",
        "priority": 0,
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/not-ready",
                "operator": "Exists",
                "tolerationSeconds": 300
            },
            {
                "effect": "NoExecute",
                "key": "node.kubernetes.io/unreachable",
                "operator": "Exists",
                "tolerationSeconds": 300
            }
        ],
        "volumes": [
            {
                "name": "kube-api-access-87h5h",
                "projected": {
                    "defaultMode": 420,
                    "sources": [
                        {
                            "serviceAccountToken": {
                                "expirationSeconds": 3607,
                                "path": "token"
                            }
                        },
                        {
                            "configMap": {
                                "items": [
                                    {
                                        "key": "ca.crt",
                                        "path": "ca.crt"
                                    }
                                ],
                                "name": "kube-root-ca.crt"
                            }
                        },
                        {
                            "downwardAPI": {
                                "items": [
                                    {
                                        "fieldRef": {
                                            "apiVersion": "v1",
                                            "fieldPath": "metadata.namespace"
                                        },
                                        "path": "namespace"
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2025-10-03T16:39:11Z",
                "observedGeneration": 1,
                "status": "True",
                "type": "PodReadyToStartContainers"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2025-10-03T16:39:06Z",
                "observedGeneration": 1,
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2025-10-03T16:39:11Z",
                "observedGeneration": 1,
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2025-10-03T16:39:11Z",
                "observedGeneration": 1,
                "status": "True",
                "type": "ContainersReady"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2025-10-03T16:39:06Z",
                "observedGeneration": 1,
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "allocatedResources": {
                    "cpu": "100m",
                    "memory": "100Mi"
                },
                "containerID": "containerd://127a3ffa5acba60f9d12b02192e800815ae5e6a59161a6c5c2877b4f61b57004",
                "image": "docker.io/library/alpine:latest",
                "imageID": "docker.io/library/alpine@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1",
                "lastState": {},
                "name": "sample-application",
                "ready": true,
                "resources": {
                    "limits": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    },
                    "requests": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    }
                },
                "restartCount": 0,
                "started": true,
                "state": {
                    "running": {
                        "startedAt": "2025-10-03T16:39:10Z"
                    }
                },
                "user": {
                    "linux": {
                        "gid": 0,
                        "supplementalGroups": [
                            0,
                            1,
                            2,
                            3,
                            4,
                            6,
                            10,
                            11,
                            20,
                            26,
                            27
                        ],
                        "uid": 0
                    }
                },
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "kube-api-access-87h5h",
                        "readOnly": true,
                        "recursiveReadOnly": "Disabled"
                    }
                ]
            }
        ],
        "hostIP": "172.20.0.2",
        "hostIPs": [
            {
                "ip": "172.20.0.2"
            }
        ],
        "observedGeneration": 1,
        "phase": "Running",
        "podIP": "10.244.0.7",
        "podIPs": [
            {
                "ip": "10.244.0.7"
            }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2025-10-03T16:39:06Z"
    }
}
`
