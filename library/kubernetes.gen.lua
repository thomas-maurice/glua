---@meta kubernetes

---@alias v1.Time string RFC3339 timestamp (metav1.Time)
---@alias v1.MicroTime string RFC3339 timestamp with microsecond precision (metav1.MicroTime)
---@alias resource.Quantity string Kubernetes resource quantity, e.g. "100Mi", "500m"
---@alias intstr.IntOrString string|number value that can be either an int or a string
---@alias v1.FieldsV1 table opaque managed-fields data

---@class kubernetes.GVKMatcher
---@field group string
---@field kind string
---@field version string

---@class admissionregistrationv1.MatchCondition
---@field expression string expression represents the expression which will be evaluated by CEL. Must evaluate to bool. CEL expressions have access to the contents of the AdmissionRequest and Authorizer, organized into CEL variables: 'object' - The object from the incoming request. The value is null for DELETE requests. 'oldObject' - The existing object. The value is null for CREATE requests. 'request' - Attributes of the admission request(/pkg/apis/admission/types.go#AdmissionRequest). 'authorizer' - A CEL Authorizer. May be used to perform authorization checks for the principal (user or service account) of the request. See https://pkg.go.dev/k8s.io/apiserver/pkg/cel/library#Authz 'authorizer.requestResource' - A CEL ResourceCheck constructed from the 'authorizer' and configured with the request resource. Documentation on CEL: https://kubernetes.io/docs/reference/using-api/cel/ Required.
---@field name string name is an identifier for this match condition, used for strategic merging of MatchConditions, as well as providing an identifier for logging purposes. A good name should be descriptive of the associated expression. Name must be a qualified name consisting of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyName', or 'my.name', or '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]') with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName') Required.

---@class admissionregistrationv1.MutatingWebhook
---@field admissionReviewVersions string[] admissionReviewVersions is an ordered list of preferred `AdmissionReview` versions the Webhook expects. API server will try to use first version in the list which it supports. If none of the versions specified in this list supported by API server, validation will fail for this object. If a persisted webhook configuration specifies allowed versions and does not include any versions known to the API Server, calls to the webhook will fail and be subject to the failure policy. +listType=atomic
---@field clientConfig admissionregistrationv1.WebhookClientConfig clientConfig defines how to communicate with the hook. Required
---@field failurePolicy string failurePolicy defines how unrecognized errors from the admission endpoint are handled - allowed values are Ignore or Fail. Defaults to Fail. +optional
---@field matchConditions admissionregistrationv1.MatchCondition[] matchConditions is a list of conditions that must be met for a request to be sent to this webhook. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. The exact matching logic is (in order): 1. If ANY matchCondition evaluates to FALSE, the webhook is skipped. 2. If ALL matchConditions evaluate to TRUE, the webhook is called. 3. If any matchCondition evaluates to an error (but none are FALSE): - If failurePolicy=Fail, reject the request - If failurePolicy=Ignore, the error is ignored and the webhook is skipped +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name +optional
---@field matchPolicy string matchPolicy defines how the "rules" list is used to match incoming requests. Allowed values are "Exact" or "Equivalent". - Exact: match a request only if it exactly matches a specified rule. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, but "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`, a request to apps/v1beta1 or extensions/v1beta1 would not be sent to the webhook. - Equivalent: match a request if modifies a resource listed in rules, even via another API group or version. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, and "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`, a request to apps/v1beta1 or extensions/v1beta1 would be converted to apps/v1 and sent to the webhook. Defaults to "Equivalent" +optional
---@field name string name is the name of the admission webhook. Name should be fully qualified, e.g., imagepolicy.kubernetes.io, where "imagepolicy" is the name of the webhook, and kubernetes.io is the name of the organization. Required.
---@field namespaceSelector v1.LabelSelector namespaceSelector decides whether to run the webhook on an object based on whether the namespace for that object matches the selector. If the object itself is a namespace, the matching is performed on object.metadata.labels. If the object is another cluster scoped resource, it never skips the webhook. For example, to run the webhook on any objects whose namespace is not associated with "runlevel" of "0" or "1"; you will set the selector as follows: "namespaceSelector": { "matchExpressions": [ { "key": "runlevel", "operator": "NotIn", "values": [ "0", "1" ] } ] } If instead you want to only run the webhook on any objects whose namespace is associated with the "environment" of "prod" or "staging"; you will set the selector as follows: "namespaceSelector": { "matchExpressions": [ { "key": "environment", "operator": "In", "values": [ "prod", "staging" ] } ] } See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ for more examples of label selectors. Default to the empty LabelSelector, which matches everything. +optional
---@field objectSelector v1.LabelSelector objectSelector decides whether to run the webhook based on if the object has matching labels. objectSelector is evaluated against both the oldObject and newObject that would be sent to the webhook, and is considered to match if either object matches the selector. A null object (oldObject in the case of create, or newObject in the case of delete) or an object that cannot have labels (like a DeploymentRollback or a PodProxyOptions object) is not considered to match. Use the object selector only if the webhook is opt-in, because end users may skip the admission webhook by setting the labels. Default to the empty LabelSelector, which matches everything. +optional
---@field reinvocationPolicy string reinvocationPolicy indicates whether this webhook should be called multiple times as part of a single admission evaluation. Allowed values are "Never" and "IfNeeded". Never: the webhook will not be called more than once in a single admission evaluation. IfNeeded: the webhook will be called at least one additional time as part of the admission evaluation if the object being admitted is modified by other admission plugins after the initial webhook call. Webhooks that specify this option *must* be idempotent, able to process objects they previously admitted. Note: * the number of additional invocations is not guaranteed to be exactly one. * if additional invocations result in further modifications to the object, webhooks are not guaranteed to be invoked again. * webhooks that use this option may be reordered to minimize the number of additional invocations. * to validate an object after all mutations are guaranteed complete, use a validating admission webhook instead. Defaults to "Never". +optional
---@field rules admissionregistrationv1.RuleWithOperations[] rules describes what operations on what resources/subresources the webhook cares about. The webhook cares about an operation if it matches _any_ Rule. However, in order to prevent ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks from putting the cluster in a state which cannot be recovered from without completely disabling the plugin, ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks are never called on admission requests for ValidatingWebhookConfiguration and MutatingWebhookConfiguration objects. +listType=atomic
---@field sideEffects string sideEffects states whether this webhook has side effects. Acceptable values are: None, NoneOnDryRun (webhooks created via v1beta1 may also specify Some or Unknown). Webhooks with side effects MUST implement a reconciliation system, since a request may be rejected by a future step in the admission chain and the side effects therefore need to be undone. Requests with the dryRun attribute will be auto-rejected if they match a webhook with sideEffects == Unknown or Some.
---@field timeoutSeconds number timeoutSeconds specifies the timeout for this webhook. After the timeout passes, the webhook call will be ignored or the API call will fail based on the failure policy. The timeout value must be between 1 and 30 seconds. Default to 10 seconds. +optional

---@class admissionregistrationv1.MutatingWebhookConfiguration
---@field metadata v1.ObjectMeta metadata is the standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata. +optional
---@field webhooks admissionregistrationv1.MutatingWebhook[] webhooks is a list of webhooks and the affected resources and operations. +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name

---@class admissionregistrationv1.MutatingWebhookConfigurationList
---@field items admissionregistrationv1.MutatingWebhookConfiguration[] List of MutatingWebhookConfiguration.
---@field metadata v1.ListMeta metadata is the standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class admissionregistrationv1.RuleWithOperations
---@field operations string[] operations is the operations the admission hook cares about - CREATE, UPDATE, DELETE, CONNECT or * for all of those operations and any future admission operations that are added. If '*' is present, the length of the slice must be one. Required. +listType=atomic

---@class admissionregistrationv1.ServiceReference
---@field name string name is the name of the service. Required
---@field namespace string namespace is the namespace of the service. Required
---@field path string path is an optional URL path which will be sent in any request to this service. +optional
---@field port number port is the port on the service that hosts the webhook. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive). +optional

---@class admissionregistrationv1.ValidatingWebhook
---@field admissionReviewVersions string[] admissionReviewVersions is an ordered list of preferred `AdmissionReview` versions the Webhook expects. API server will try to use first version in the list which it supports. If none of the versions specified in this list supported by API server, validation will fail for this object. If a persisted webhook configuration specifies allowed versions and does not include any versions known to the API Server, calls to the webhook will fail and be subject to the failure policy. +listType=atomic
---@field clientConfig admissionregistrationv1.WebhookClientConfig clientConfig defines how to communicate with the hook. Required
---@field failurePolicy string failurePolicy defines how unrecognized errors from the admission endpoint are handled - allowed values are Ignore or Fail. Defaults to Fail. +optional
---@field matchConditions admissionregistrationv1.MatchCondition[] matchConditions is a list of conditions that must be met for a request to be sent to this webhook. Match conditions filter requests that have already been matched by the rules, namespaceSelector, and objectSelector. An empty list of matchConditions matches all requests. There are a maximum of 64 match conditions allowed. The exact matching logic is (in order): 1. If ANY matchCondition evaluates to FALSE, the webhook is skipped. 2. If ALL matchConditions evaluate to TRUE, the webhook is called. 3. If any matchCondition evaluates to an error (but none are FALSE): - If failurePolicy=Fail, reject the request - If failurePolicy=Ignore, the error is ignored and the webhook is skipped +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name +optional
---@field matchPolicy string matchPolicy defines how the "rules" list is used to match incoming requests. Allowed values are "Exact" or "Equivalent". - Exact: match a request only if it exactly matches a specified rule. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, but "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`, a request to apps/v1beta1 or extensions/v1beta1 would not be sent to the webhook. - Equivalent: match a request if modifies a resource listed in rules, even via another API group or version. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, and "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`, a request to apps/v1beta1 or extensions/v1beta1 would be converted to apps/v1 and sent to the webhook. Defaults to "Equivalent" +optional
---@field name string name is the name of the admission webhook. Name should be fully qualified, e.g., imagepolicy.kubernetes.io, where "imagepolicy" is the name of the webhook, and kubernetes.io is the name of the organization. Required.
---@field namespaceSelector v1.LabelSelector namespaceSelector decides whether to run the webhook on an object based on whether the namespace for that object matches the selector. If the object itself is a namespace, the matching is performed on object.metadata.labels. If the object is another cluster scoped resource, it never skips the webhook. For example, to run the webhook on any objects whose namespace is not associated with "runlevel" of "0" or "1"; you will set the selector as follows: "namespaceSelector": { "matchExpressions": [ { "key": "runlevel", "operator": "NotIn", "values": [ "0", "1" ] } ] } If instead you want to only run the webhook on any objects whose namespace is associated with the "environment" of "prod" or "staging"; you will set the selector as follows: "namespaceSelector": { "matchExpressions": [ { "key": "environment", "operator": "In", "values": [ "prod", "staging" ] } ] } See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels for more examples of label selectors. Default to the empty LabelSelector, which matches everything. +optional
---@field objectSelector v1.LabelSelector objectSelector decides whether to run the webhook based on if the object has matching labels. objectSelector is evaluated against both the oldObject and newObject that would be sent to the webhook, and is considered to match if either object matches the selector. A null object (oldObject in the case of create, or newObject in the case of delete) or an object that cannot have labels (like a DeploymentRollback or a PodProxyOptions object) is not considered to match. Use the object selector only if the webhook is opt-in, because end users may skip the admission webhook by setting the labels. Default to the empty LabelSelector, which matches everything. +optional
---@field rules admissionregistrationv1.RuleWithOperations[] rules describes what operations on what resources/subresources the webhook cares about. The webhook cares about an operation if it matches _any_ Rule. However, in order to prevent ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks from putting the cluster in a state which cannot be recovered from without completely disabling the plugin, ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks are never called on admission requests for ValidatingWebhookConfiguration and MutatingWebhookConfiguration objects. +listType=atomic
---@field sideEffects string sideEffects states whether this webhook has side effects. Acceptable values are: None, NoneOnDryRun (webhooks created via v1beta1 may also specify Some or Unknown). Webhooks with side effects MUST implement a reconciliation system, since a request may be rejected by a future step in the admission chain and the side effects therefore need to be undone. Requests with the dryRun attribute will be auto-rejected if they match a webhook with sideEffects == Unknown or Some.
---@field timeoutSeconds number timeoutSeconds specifies the timeout for this webhook. After the timeout passes, the webhook call will be ignored or the API call will fail based on the failure policy. The timeout value must be between 1 and 30 seconds. Default to 10 seconds. +optional

---@class admissionregistrationv1.ValidatingWebhookConfiguration
---@field metadata v1.ObjectMeta metadata is the standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata. +optional
---@field webhooks admissionregistrationv1.ValidatingWebhook[] webhooks is a list of webhooks and the affected resources and operations. +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name

---@class admissionregistrationv1.ValidatingWebhookConfigurationList
---@field items admissionregistrationv1.ValidatingWebhookConfiguration[] List of ValidatingWebhookConfiguration.
---@field metadata v1.ListMeta metadata is the standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class admissionregistrationv1.WebhookClientConfig
---@field caBundle number[] caBundle is a PEM encoded CA bundle which will be used to validate the webhook's server certificate. If unspecified, system trust roots on the apiserver are used. +optional
---@field service admissionregistrationv1.ServiceReference service is a reference to the service for this webhook. Either `service` or `url` must be specified. If the webhook is running within the cluster, then you should use `service`. +optional
---@field url string url gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified. The `host` should not refer to a service running in the cluster; use the `service` field instead. The host might be resolved via external DNS in some apiservers (e.g., `kube-apiserver` cannot resolve in-cluster DNS as that would be a layering violation). `host` may also be an IP address. Please note that using `localhost` or `127.0.0.1` as a `host` is risky unless you take great care to run this webhook on all hosts which run an apiserver which might need to make calls to this webhook. Such installs are likely to be non-portable, i.e., not easy to turn up in a new cluster. The scheme must be "https"; the URL must begin with "https://". A path is optional, and if present may be any string permissible in a URL. You may use the path to pass an arbitrary string to the webhook, for example, a cluster identifier. Attempting to use a user or basic auth e.g. "user:password@" is not allowed. Fragments ("#...") and query parameters ("?...") are not allowed, either. +optional

---@class appsv1.DaemonSet
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec appsv1.DaemonSetSpec The desired behavior of this daemon set. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +required
---@field status appsv1.DaemonSetStatus The current status of this daemon set. This data may be out of date by some window of time. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class appsv1.DaemonSetCondition
---@field lastTransitionTime v1.Time Last time the condition transitioned from one status to another. +optional
---@field message string A human readable message indicating details about the transition. +optional
---@field reason string The reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown. +optional
---@field type string Type of DaemonSet condition. +optional

---@class appsv1.DaemonSetList
---@field items appsv1.DaemonSet[] A list of daemon sets.
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class appsv1.DaemonSetSpec
---@field minReadySeconds number The minimum number of seconds for which a newly created DaemonSet pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready). +optional
---@field revisionHistoryLimit number The number of old history to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified. Defaults to 10. +optional
---@field selector v1.LabelSelector A label query over pods that are managed by the daemon set. Must match in order to be controlled. It must match the pod template's labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors +required
---@field template corev1.PodTemplateSpec An object that describes the pod that will be created. The DaemonSet will create exactly one copy of this pod on every node that matches the template's node selector (or on every node if no node selector is specified). The only allowed template.spec.restartPolicy value is "Always". More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#pod-template +required
---@field updateStrategy appsv1.DaemonSetUpdateStrategy An update strategy to replace existing DaemonSet pods with new pods. +optional

---@class appsv1.DaemonSetStatus
---@field collisionCount number Count of hash collisions for the DaemonSet. The DaemonSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision. +optional
---@field conditions appsv1.DaemonSetCondition[] Represents the latest available observations of a DaemonSet's current state. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field currentNumberScheduled number The number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod. More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
---@field desiredNumberScheduled number The total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod). More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
---@field numberAvailable number The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available (ready for at least spec.minReadySeconds) +optional
---@field numberMisscheduled number The number of nodes that are running the daemon pod, but are not supposed to run the daemon pod. More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
---@field numberReady number numberReady is the number of nodes that should be running the daemon pod and have one or more of the daemon pod running with a Ready Condition.
---@field numberUnavailable number The number of nodes that should be running the daemon pod and have none of the daemon pod running and available (ready for at least spec.minReadySeconds) +optional
---@field observedGeneration number The most recent generation observed by the daemon set controller. +optional
---@field updatedNumberScheduled number The total number of nodes that are running updated daemon pod +optional

---@class appsv1.DaemonSetUpdateStrategy
---@field rollingUpdate appsv1.RollingUpdateDaemonSet Rolling update config params. Present only if type = "RollingUpdate". --- TODO: Update this to follow our convention for oneOf, whatever we decide it to be. Same as Deployment `strategy.rollingUpdate`. See https://github.com/kubernetes/kubernetes/issues/35345 +optional
---@field type string Type of daemon set update. Can be "RollingUpdate" or "OnDelete". Default is RollingUpdate. +optional

---@class appsv1.Deployment
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec appsv1.DeploymentSpec Specification of the desired behavior of the Deployment. +required
---@field status appsv1.DeploymentStatus Most recently observed status of the Deployment. +optional

---@class appsv1.DeploymentCondition
---@field lastTransitionTime v1.Time Last time the condition transitioned from one status to another. +optional
---@field lastUpdateTime v1.Time The last time this condition was updated. +optional
---@field message string A human readable message indicating details about the transition. +optional
---@field reason string The reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown. +optional
---@field type string Type of deployment condition. +optional

---@class appsv1.DeploymentList
---@field items appsv1.Deployment[] Items is the list of Deployments.
---@field metadata v1.ListMeta Standard list metadata. +optional

---@class appsv1.DeploymentSpec
---@field minReadySeconds number Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready) +optional
---@field paused boolean Indicates that the deployment is paused. +optional
---@field progressDeadlineSeconds number The maximum time in seconds for a deployment to make progress before it is considered to be failed. The deployment controller will continue to process failed deployments and a condition with a ProgressDeadlineExceeded reason will be surfaced in the deployment status. Note that progress will not be estimated during the time a deployment is paused. Defaults to 600s. +optional
---@field replicas number Number of desired pods. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1. +optional
---@field revisionHistoryLimit number The number of old ReplicaSets to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified. Defaults to 10. +optional
---@field selector v1.LabelSelector Label selector for pods. Existing ReplicaSets whose pods are selected by this will be the ones affected by this deployment. It must match the pod template's labels. +required
---@field strategy appsv1.DeploymentStrategy The deployment strategy to use to replace existing pods with new ones. +optional +patchStrategy=retainKeys
---@field template corev1.PodTemplateSpec Template describes the pods that will be created. The only allowed template.spec.restartPolicy value is "Always". +required

---@class appsv1.DeploymentStatus
---@field availableReplicas number Total number of available non-terminating pods (ready for at least minReadySeconds) targeted by this deployment. +optional
---@field collisionCount number Count of hash collisions for the Deployment. The Deployment controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ReplicaSet. +optional
---@field conditions appsv1.DeploymentCondition[] Represents the latest available observations of a deployment's current state. +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field observedGeneration number The generation observed by the deployment controller. +optional
---@field readyReplicas number Total number of non-terminating pods targeted by this Deployment with a Ready Condition. +optional
---@field replicas number Total number of non-terminating pods targeted by this deployment (their labels match the selector). +optional
---@field terminatingReplicas number Total number of terminating pods targeted by this deployment. Terminating pods have a non-null .metadata.deletionTimestamp and have not yet reached the Failed or Succeeded .status.phase. This is a beta field and requires enabling DeploymentReplicaSetTerminatingReplicas feature (enabled by default). +optional
---@field unavailableReplicas number Total number of unavailable pods targeted by this deployment. This is the total number of pods that are still required for the deployment to have 100% available capacity. They may either be pods that are running but not yet available or pods that still have not been created. +optional
---@field updatedReplicas number Total number of non-terminating pods targeted by this deployment that have the desired template spec. +optional

---@class appsv1.DeploymentStrategy
---@field rollingUpdate appsv1.RollingUpdateDeployment Rolling update config params. Present only if DeploymentStrategyType = RollingUpdate. --- TODO: Update this to follow our convention for oneOf, whatever we decide it to be. +optional
---@field type string Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate. +optional

---@class appsv1.ReplicaSet
---@field metadata v1.ObjectMeta If the Labels of a ReplicaSet are empty, they are defaulted to be the same as the Pod(s) that the ReplicaSet manages. Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec appsv1.ReplicaSetSpec Spec defines the specification of the desired behavior of the ReplicaSet. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +required
---@field status appsv1.ReplicaSetStatus Status is the most recently observed status of the ReplicaSet. This data may be out of date by some window of time. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class appsv1.ReplicaSetCondition
---@field lastTransitionTime v1.Time The last time the condition transitioned from one status to another. +optional
---@field message string A human readable message indicating details about the transition. +optional
---@field reason string The reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown. +optional
---@field type string Type of replica set condition. +optional

---@class appsv1.ReplicaSetList
---@field items appsv1.ReplicaSet[] List of ReplicaSets. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class appsv1.ReplicaSetSpec
---@field minReadySeconds number Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready) +optional
---@field replicas number Replicas is the number of desired pods. This is a pointer to distinguish between explicit zero and unspecified. Defaults to 1. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset +optional
---@field selector v1.LabelSelector Selector is a label query over pods that should match the replica count. Label keys and values that must match in order to be controlled by this replica set. It must match the pod template's labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors +required
---@field template corev1.PodTemplateSpec Template is the object that describes the pod that will be created if insufficient replicas are detected. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/#pod-template +optional

---@class appsv1.ReplicaSetStatus
---@field availableReplicas number The number of available non-terminating pods (ready for at least minReadySeconds) for this replica set. +optional
---@field conditions appsv1.ReplicaSetCondition[] Represents the latest available observations of a replica set's current state. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field fullyLabeledReplicas number The number of non-terminating pods that have labels matching the labels of the pod template of the replicaset. +optional
---@field observedGeneration number ObservedGeneration reflects the generation of the most recently observed ReplicaSet. +optional
---@field readyReplicas number The number of non-terminating pods targeted by this ReplicaSet with a Ready Condition. +optional
---@field replicas number Replicas is the most recently observed number of non-terminating pods. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset
---@field terminatingReplicas number The number of terminating pods for this replica set. Terminating pods have a non-null .metadata.deletionTimestamp and have not yet reached the Failed or Succeeded .status.phase. This is a beta field and requires enabling DeploymentReplicaSetTerminatingReplicas feature (enabled by default). +optional

---@class appsv1.RollingUpdateDaemonSet
---@field maxSurge intstr.IntOrString The maximum number of nodes with an existing available DaemonSet pod that can have an updated DaemonSet pod during during an update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up to a minimum of 1. Default value is 0. Example: when this is set to 30%, at most 30% of the total number of nodes that should be running the daemon pod (i.e. status.desiredNumberScheduled) can have their a new pod created before the old pod is marked as deleted. The update starts by launching new pods on 30% of nodes. Once an updated pod is available (Ready for at least minReadySeconds) the old DaemonSet pod on that node is marked deleted. If the old pod becomes unavailable for any reason (Ready transitions to false, is evicted, or is drained) an updated pod is immediately created on that node without considering surge limits. Allowing surge implies the possibility that the resources consumed by the daemonset on any given node can double if the readiness check fails, and so resource intensive daemonsets should take into account that they may cause evictions during disruption. +optional
---@field maxUnavailable intstr.IntOrString The maximum number of DaemonSet pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of total number of DaemonSet pods at the start of the update (ex: 10%). Absolute number is calculated from percentage by rounding up. This cannot be 0 if MaxSurge is 0 Default value is 1. Example: when this is set to 30%, at most 30% of the total number of nodes that should be running the daemon pod (i.e. status.desiredNumberScheduled) can have their pods stopped for an update at any given time. The update starts by stopping at most 30% of those DaemonSet pods and then brings up new DaemonSet pods in their place. Once the new pods are available, it then proceeds onto other DaemonSet pods, thus ensuring that at least 70% of original number of DaemonSet pods are available at all times during the update. +optional

---@class appsv1.RollingUpdateDeployment
---@field maxSurge intstr.IntOrString The maximum number of pods that can be scheduled above the desired number of pods. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up. Defaults to 25%. Example: when this is set to 30%, the new ReplicaSet can be scaled up immediately when the rolling update starts, such that the total number of old and new pods do not exceed 130% of desired pods. Once old pods have been killed, new ReplicaSet can be scaled up further, ensuring that total number of pods running at any time during the update is at most 130% of desired pods. +optional
---@field maxUnavailable intstr.IntOrString The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding down. This can not be 0 if MaxSurge is 0. Defaults to 25%. Example: when this is set to 30%, the old ReplicaSet can be scaled down to 70% of desired pods immediately when the rolling update starts. Once new pods are ready, old ReplicaSet can be scaled down further, followed by scaling up the new ReplicaSet, ensuring that the total number of pods available at all times during the update is at least 70% of desired pods. +optional

---@class appsv1.RollingUpdateStatefulSetStrategy
---@field maxUnavailable intstr.IntOrString The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding up. This can not be 0. Defaults to 1. This field is beta-level and is enabled by default. The field applies to all pods in the range 0 to Replicas-1. That means if there is any unavailable pod in the range 0 to Replicas-1, it will be counted towards MaxUnavailable. This setting might not be effective for the OrderedReady podManagementPolicy. That policy ensures pods are created and become ready one at a time. +featureGate=MaxUnavailableStatefulSet +optional
---@field partition number Partition indicates the ordinal at which the StatefulSet should be partitioned for updates. During a rolling update, all pods from ordinal Replicas-1 to Partition are updated. All pods from ordinal Partition-1 to 0 remain untouched. This is helpful in being able to do a canary based deployment. The default value is 0. +optional

---@class appsv1.StatefulSet
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec appsv1.StatefulSetSpec Spec defines the desired identities of pods in this set. +required
---@field status appsv1.StatefulSetStatus Status is the current status of Pods in this StatefulSet. This data may be out of date by some window of time. +optional

---@class appsv1.StatefulSetCondition
---@field lastTransitionTime v1.Time Last time the condition transitioned from one status to another. +optional
---@field message string A human readable message indicating details about the transition. +optional
---@field reason string The reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown. +optional
---@field type string Type of statefulset condition. +optional

---@class appsv1.StatefulSetList
---@field items appsv1.StatefulSet[] Items is the list of stateful sets.
---@field metadata v1.ListMeta Standard list's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class appsv1.StatefulSetOrdinals
---@field start number start is the number representing the first replica's index. It may be used to number replicas from an alternate index (eg: 1-indexed) over the default 0-indexed names, or to orchestrate progressive movement of replicas from one StatefulSet to another. If set, replica indices will be in the range: [.spec.ordinals.start, .spec.ordinals.start + .spec.replicas). If unset, defaults to 0. Replica indices will be in the range: [0, .spec.replicas). +optional

---@class appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy
---@field whenDeleted string WhenDeleted specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is deleted. The default policy of `Retain` causes PVCs to not be affected by StatefulSet deletion. The `Delete` policy causes those PVCs to be deleted. +optional
---@field whenScaled string WhenScaled specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is scaled down. The default policy of `Retain` causes PVCs to not be affected by a scaledown. The `Delete` policy causes the associated PVCs for any excess pods above the replica count to be deleted. +optional

---@class appsv1.StatefulSetSpec
---@field minReadySeconds number Minimum number of seconds for which a newly created pod should be ready without any of its container crashing for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready) +optional
---@field ordinals appsv1.StatefulSetOrdinals ordinals controls the numbering of replica indices in a StatefulSet. The default ordinals behavior assigns a "0" index to the first replica and increments the index by one for each additional replica requested. +optional
---@field persistentVolumeClaimRetentionPolicy appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent volume claims created from volumeClaimTemplates. By default, all persistent volume claims are created as needed and retained until manually deleted. This policy allows the lifecycle to be altered, for example by deleting persistent volume claims when their stateful set is deleted, or when their pod is scaled down. +optional
---@field podManagementPolicy string podManagementPolicy controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down. The default policy is `OrderedReady`, where pods are created in increasing order (pod-0, then pod-1, etc) and the controller will wait until each pod is ready before continuing. When scaling down, the pods are removed in the opposite order. The alternative policy is `Parallel` which will create pods in parallel to match the desired scale without waiting, and on scale down will delete all pods at once. +optional +k8s:alpha(since: "1.37")=+k8s:immutable +k8s:alpha(since: "1.37")=+k8s:optional
---@field replicas number replicas is the desired number of replicas of the given Template. These are replicas in the sense that they are instantiations of the same Template, but individual replicas also have a consistent identity. If unspecified, defaults to 1. TODO: Consider a rename of this field. +optional
---@field revisionHistoryLimit number revisionHistoryLimit is the maximum number of revisions that will be maintained in the StatefulSet's revision history. The revision history consists of all revisions not represented by a currently applied StatefulSetSpec version. The default value is 10. +optional
---@field selector v1.LabelSelector selector is a label query over pods that should match the replica count. It must match the pod template's labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors +required +k8s:alpha(since: "1.37")=+k8s:required +k8s:alpha(since: "1.37")=+k8s:immutable
---@field serviceName string serviceName is the name of the service that governs this StatefulSet. This service must exist before the StatefulSet, and is responsible for the network identity of the set. Pods get DNS/hostnames that follow the pattern: pod-specific-string.serviceName.default.svc.cluster.local where "pod-specific-string" is managed by the StatefulSet controller. +optional +k8s:alpha(since: "1.37")=+k8s:immutable +k8s:alpha(since: "1.37")=+k8s:optional
---@field template corev1.PodTemplateSpec template is the object that describes the pod that will be created if insufficient replicas are detected. Each pod stamped out by the StatefulSet will fulfill this Template, but have a unique identity from the rest of the StatefulSet. Each pod will be named with the format <statefulsetname>-<podindex>. For example, a pod in a StatefulSet named "web" with index number "3" would be named "web-3". The only allowed template.spec.restartPolicy value is "Always". +required
---@field updateStrategy appsv1.StatefulSetUpdateStrategy updateStrategy indicates the StatefulSetUpdateStrategy that will be employed to update Pods in the StatefulSet when a revision is made to Template. +optional
---@field volumeClaimTemplates corev1.PersistentVolumeClaim[] volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. A claim in this list takes precedence over any volumes in the template, with the same name. TODO: Define the behavior if a claim already exists with the same name. +optional +k8s:alpha(since: "1.37")=+k8s:immutable +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:eachVal=+k8s:opaqueType +listType=atomic

---@class appsv1.StatefulSetStatus
---@field availableReplicas number Total number of available pods (ready for at least minReadySeconds) targeted by this statefulset. +optional
---@field collisionCount number collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision. +optional
---@field conditions appsv1.StatefulSetCondition[] Represents the latest available observations of a statefulset's current state. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field currentReplicas number currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.
---@field currentRevision string currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).
---@field observedGeneration number observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the StatefulSet's generation, which is updated on mutation by the API Server. +optional
---@field readyReplicas number readyReplicas is the number of pods created for this StatefulSet with a Ready Condition.
---@field replicas number replicas is the number of Pods created by the StatefulSet controller.
---@field updateRevision string updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas)
---@field updatedReplicas number updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.

---@class appsv1.StatefulSetUpdateStrategy
---@field rollingUpdate appsv1.RollingUpdateStatefulSetStrategy RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType. +optional
---@field type string Type indicates the type of the StatefulSetUpdateStrategy. Default is RollingUpdate. +optional

---@class autoscalingv2.ContainerResourceMetricSource
---@field container string container is the name of the container in the pods of the scaling target
---@field name string name is the name of the resource in question.
---@field target autoscalingv2.MetricTarget target specifies the target value for the given metric

---@class autoscalingv2.ContainerResourceMetricStatus
---@field container string container is the name of the container in the pods of the scaling target
---@field current autoscalingv2.MetricValueStatus current contains the current value for the given metric
---@field name string name is the name of the resource in question.

---@class autoscalingv2.CrossVersionObjectReference
---@field apiVersion string apiVersion is the API version of the referent +optional
---@field kind string kind is the kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +k8s:alpha(since: "1.37")=+k8s:required
---@field name string name is the name of the referent; More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names +k8s:alpha(since: "1.37")=+k8s:required

---@class autoscalingv2.ExternalMetricSource
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector
---@field target autoscalingv2.MetricTarget target specifies the target value for the given metric

---@class autoscalingv2.ExternalMetricStatus
---@field current autoscalingv2.MetricValueStatus current contains the current value for the given metric
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector

---@class autoscalingv2.HPAScalingPolicy
---@field periodSeconds number periodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
---@field type string type is used to specify the scaling policy.
---@field value number value contains the amount of change which is permitted by the policy. It must be greater than zero

---@class autoscalingv2.HPAScalingRules
---@field policies autoscalingv2.HPAScalingPolicy[] policies is a list of potential scaling polices which can be used during scaling. If not set, use the default values: - For scale up: allow doubling the number of pods, or an absolute change of 4 pods in a 15s window. - For scale down: allow all pods to be removed in a 15s window. +listType=atomic +optional
---@field selectPolicy string selectPolicy is used to specify which policy should be used. If not set, the default value Max is used. +optional
---@field stabilizationWindowSeconds number stabilizationWindowSeconds is the number of seconds for which past recommendations should be considered while scaling up or scaling down. StabilizationWindowSeconds must be greater than or equal to zero and less than or equal to 3600 (one hour). If not set, use the default values: - For scale up: 0 (i.e. no stabilization is done). - For scale down: 300 (i.e. the stabilization window is 300 seconds long). +optional
---@field tolerance resource.Quantity tolerance is the tolerance on the ratio between the current and desired metric value under which no updates are made to the desired number of replicas (e.g. 0.01 for 1%). Must be greater than or equal to zero. If not set, the default cluster-wide tolerance is applied (by default 10%). For example, if autoscaling is configured with a memory consumption target of 100Mi, and scale-down and scale-up tolerances of 5% and 1% respectively, scaling will be triggered when the actual consumption falls below 95Mi or exceeds 101Mi. +featureGate=HPAConfigurableTolerance +optional

---@class autoscalingv2.HorizontalPodAutoscaler
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec autoscalingv2.HorizontalPodAutoscalerSpec spec is the specification for the behaviour of the autoscaler. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status. +required
---@field status autoscalingv2.HorizontalPodAutoscalerStatus status is the current information about the autoscaler. +optional

---@class autoscalingv2.HorizontalPodAutoscalerBehavior
---@field scaleDown autoscalingv2.HPAScalingRules scaleDown is scaling policy for scaling Down. If not set, the default value is to allow to scale down to minReplicas pods, with a 300 second stabilization window (i.e., the highest recommendation for the last 300sec is used). +optional
---@field scaleUp autoscalingv2.HPAScalingRules scaleUp is scaling policy for scaling Up. If not set, the default value is the higher of: * increase no more than 4 pods per 60 seconds * double the number of pods per 60 seconds No stabilization is used. +optional

---@class autoscalingv2.HorizontalPodAutoscalerCondition
---@field lastTransitionTime v1.Time lastTransitionTime is the last time the condition transitioned from one status to another +optional
---@field message string message is a human-readable explanation containing details about the transition +optional
---@field observedGeneration number observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. +optional
---@field reason string reason is the reason for the condition's last transition. +optional
---@field status string status is the status of the condition (True, False, Unknown)
---@field type string type describes the current condition

---@class autoscalingv2.HorizontalPodAutoscalerList
---@field items autoscalingv2.HorizontalPodAutoscaler[] items is the list of horizontal pod autoscaler objects.
---@field metadata v1.ListMeta metadata is the standard list metadata. +optional

---@class autoscalingv2.HorizontalPodAutoscalerSpec
---@field behavior autoscalingv2.HorizontalPodAutoscalerBehavior behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default HPAScalingRules for scale up and scale down are used. +optional
---@field maxReplicas number maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas. +required +k8s:beta(since: "1.37")=+k8s:required +k8s:beta(since: "1.37")=+k8s:minimum=1
---@field metrics autoscalingv2.MetricSpec[] metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used). The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods. Ergo, metrics used must decrease as the pod count is increased, and vice-versa. See the individual metric source types for more information about how each type of metric must respond. If not set, the default metric will be set to 80% average CPU utilization. +listType=atomic +optional +k8s:alpha(since: "1.37")=+k8s:optional
---@field minReplicas number minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down. It defaults to 1 pod. minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured. Scaling is active as long as at least one metric value is available. +optional +k8s:beta(since: "1.37")=+k8s:optional +k8s:beta(since: "1.37")=+k8s:ifEnabled(HPAScaleToZero)=+k8s:minimum=0 +k8s:beta(since: "1.37")=+k8s:ifDisabled(HPAScaleToZero)=+k8s:minimum=1
---@field scaleTargetRef autoscalingv2.CrossVersionObjectReference scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.

---@class autoscalingv2.HorizontalPodAutoscalerStatus
---@field conditions autoscalingv2.HorizontalPodAutoscalerCondition[] conditions is the set of conditions required for this autoscaler to scale its target, and indicates whether or not those conditions are met. +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type +optional
---@field currentMetrics autoscalingv2.MetricStatus[] currentMetrics is the last read state of the metrics used by this autoscaler. +listType=atomic +optional +k8s:alpha(since: "1.37")=+k8s:optional
---@field currentReplicas number currentReplicas is current number of replicas of pods managed by this autoscaler, as last seen by the autoscaler. +optional
---@field desiredReplicas number desiredReplicas is the desired number of replicas of pods managed by this autoscaler, as last calculated by the autoscaler.
---@field lastScaleTime v1.Time lastScaleTime is the last time the HorizontalPodAutoscaler scaled the number of pods, used by the autoscaler to control how often the number of pods is changed. +optional
---@field observedGeneration number observedGeneration is the most recent generation observed by this autoscaler. +optional

---@class autoscalingv2.MetricIdentifier
---@field name string name is the name of the given metric
---@field selector v1.LabelSelector selector is the string-encoded form of a standard kubernetes label selector for the given metric When set, it is passed as an additional parameter to the metrics server for more specific metrics scoping. When unset, just the metricName will be used to gather metrics. +optional

---@class autoscalingv2.MetricSpec
---@field containerResource autoscalingv2.ContainerResourceMetricSource containerResource refers to a resource metric (such as those specified in requests and limits) known to Kubernetes describing a single container in each pod of the current scale target (e.g. CPU or memory). Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the "pods" source. +optional
---@field external autoscalingv2.ExternalMetricSource external refers to a global metric that is not associated with any Kubernetes object. It allows autoscaling based on information coming from components running outside of cluster (for example length of queue in cloud messaging service, or QPS from loadbalancer running outside of cluster). +optional
---@field object autoscalingv2.ObjectMetricSource object refers to a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object). +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:opaqueType
---@field pods autoscalingv2.PodsMetricSource pods refers to a metric describing each pod in the current scale target (for example, transactions-processed-per-second). The values will be averaged together before being compared to the target value. +optional
---@field resource autoscalingv2.ResourceMetricSource resource refers to a resource metric (such as those specified in requests and limits) known to Kubernetes describing each pod in the current scale target (e.g. CPU or memory). Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the "pods" source. +optional
---@field type string type is the type of metric source. It should be one of "ContainerResource", "External", "Object", "Pods" or "Resource", each mapping to a matching field in the object.

---@class autoscalingv2.MetricStatus
---@field containerResource autoscalingv2.ContainerResourceMetricStatus containerResource refers to a resource metric (such as those specified in requests and limits) known to Kubernetes describing a single container in each pod in the current scale target (e.g. CPU or memory). Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the "pods" source. +optional
---@field external autoscalingv2.ExternalMetricStatus external refers to a global metric that is not associated with any Kubernetes object. It allows autoscaling based on information coming from components running outside of cluster (for example length of queue in cloud messaging service, or QPS from loadbalancer running outside of cluster). +optional
---@field object autoscalingv2.ObjectMetricStatus object refers to a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object). +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:opaqueType
---@field pods autoscalingv2.PodsMetricStatus pods refers to a metric describing each pod in the current scale target (for example, transactions-processed-per-second). The values will be averaged together before being compared to the target value. +optional
---@field resource autoscalingv2.ResourceMetricStatus resource refers to a resource metric (such as those specified in requests and limits) known to Kubernetes describing each pod in the current scale target (e.g. CPU or memory). Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the "pods" source. +optional
---@field type string type is the type of metric source. It will be one of "ContainerResource", "External", "Object", "Pods" or "Resource", each corresponds to a matching field in the object.

---@class autoscalingv2.MetricTarget
---@field averageUtilization number averageUtilization is the target value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods. Currently only valid for Resource metric source type +optional
---@field averageValue resource.Quantity averageValue is the target value of the average of the metric across all relevant pods (as a quantity) +optional
---@field type string type represents whether the metric type is Utilization, Value, or AverageValue
---@field value resource.Quantity value is the target value of the metric (as a quantity). +optional

---@class autoscalingv2.MetricValueStatus
---@field averageUtilization number averageUtilization is the current value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods. +optional
---@field averageValue resource.Quantity averageValue is the current value of the average of the metric across all relevant pods (as a quantity) +optional
---@field value resource.Quantity value is the current value of the metric (as a quantity). +optional

---@class autoscalingv2.ObjectMetricSource
---@field describedObject autoscalingv2.CrossVersionObjectReference describedObject specifies the descriptions of a object,such as kind,name apiVersion
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector
---@field target autoscalingv2.MetricTarget target specifies the target value for the given metric

---@class autoscalingv2.ObjectMetricStatus
---@field current autoscalingv2.MetricValueStatus current contains the current value for the given metric
---@field describedObject autoscalingv2.CrossVersionObjectReference describedObject specifies the descriptions of a object,such as kind,name apiVersion
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector

---@class autoscalingv2.PodsMetricSource
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector
---@field target autoscalingv2.MetricTarget target specifies the target value for the given metric

---@class autoscalingv2.PodsMetricStatus
---@field current autoscalingv2.MetricValueStatus current contains the current value for the given metric
---@field metric autoscalingv2.MetricIdentifier metric identifies the target metric by name and selector

---@class autoscalingv2.ResourceMetricSource
---@field name string name is the name of the resource in question.
---@field target autoscalingv2.MetricTarget target specifies the target value for the given metric

---@class autoscalingv2.ResourceMetricStatus
---@field current autoscalingv2.MetricValueStatus current contains the current value for the given metric
---@field name string name is the name of the resource in question.

---@class batchv1.CronJob
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec batchv1.CronJobSpec Specification of the desired behavior of a cron job, including the schedule. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +required
---@field status batchv1.CronJobStatus Current status of a cron job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class batchv1.CronJobList
---@field items batchv1.CronJob[] items is the list of CronJobs.
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class batchv1.CronJobSpec
---@field concurrencyPolicy string Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one +optional
---@field failedJobsHistoryLimit number The number of failed finished jobs to retain. Value must be non-negative integer. Defaults to 1. +optional
---@field jobTemplate batchv1.JobTemplateSpec Specifies the job that will be created when executing a CronJob.
---@field schedule string The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron. +required +k8s:beta(since: "1.37")=+k8s:required
---@field startingDeadlineSeconds number Optional deadline in seconds for starting the job if it misses scheduled time for any reason. Missed jobs executions will be counted as failed ones. +optional
---@field successfulJobsHistoryLimit number The number of successful finished jobs to retain. Value must be non-negative integer. Defaults to 3. +optional
---@field suspend boolean This flag tells the controller to suspend subsequent executions, it does not apply to already started executions. Defaults to false. +optional
---@field timeZone string The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones. If not specified, this will default to the time zone of the kube-controller-manager process. The set of valid time zone names and the time zone offset is loaded from the system-wide time zone database by the API server during CronJob validation and the controller manager during execution. If no system-wide time zone database can be found a bundled version of the database is used instead. If the time zone name becomes invalid during the lifetime of a CronJob or due to a change in host configuration, the controller will stop creating new new Jobs and will create a system event with the reason UnknownTimeZone. More information can be found in https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#time-zones +optional

---@class batchv1.CronJobStatus
---@field active corev1.ObjectReference[] A list of pointers to currently running jobs. +optional +listType=atomic
---@field lastScheduleTime v1.Time Information when was the last time the job was successfully scheduled. +optional
---@field lastSuccessfulTime v1.Time Information when was the last time the job successfully completed. +optional

---@class batchv1.Job
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec batchv1.JobSpec Specification of the desired behavior of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status batchv1.JobStatus Current status of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class batchv1.JobCondition
---@field lastProbeTime v1.Time Last time the condition was checked. +optional
---@field lastTransitionTime v1.Time Last time the condition transit from one status to another. +optional
---@field message string Human readable message indicating details about last transition. +optional
---@field reason string (brief) reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown.
---@field type string Type of job condition, Complete or Failed.

---@class batchv1.JobList
---@field items batchv1.Job[] items is the list of Jobs.
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class batchv1.JobSchedulingConfiguration
---@field disruptionMode schedulingv1alpha3.WorkloadPodGroupDisruptionMode DisruptionMode defines the mode in which the Job's pods can be disrupted. One of Single, All. This field is immutable after creation: it may not be added or removed, and the selected mode may not be changed. +optional +k8s:optional +k8s:immutable
---@field resourceClaims schedulingv1alpha3.WorkloadPodGroupResourceClaim[] ResourceClaims defines which ResourceClaims may be shared among Pods in the Job. Pods consume the devices allocated to a PodGroup's claim by defining a claim in its own Spec.ResourceClaims that matches the PodGroup's claim exactly. The claim must have the same name and refer to the same ResourceClaim or ResourceClaimTemplate. At most 4 claims may be set, matching the limit on the resulting PodGroup. This list is immutable after creation: entries may neither be added, removed, nor modified. +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name +k8s:optional +k8s:listType=map +k8s:listMapKey=name +k8s:maxItems=4 +k8s:immutable
---@field schedulingConstraints schedulingv1alpha3.WorkloadPodGroupSchedulingConstraints SchedulingConstraints defines scheduling constraints (e.g. topology) for the Job's pods. This field is immutable after creation. +optional +k8s:optional +k8s:immutable
---@field schedulingPolicy schedulingv1alpha3.WorkloadPodGroupSchedulingPolicy SchedulingPolicy defines the scheduling policy for this Job. Exactly one of Basic or Gang must be set. This field is immutable after creation: the policy may not be added or removed. The policy variant (basic/gang) is frozen by hand-written validation; only schedulingPolicy.gang.minCount may be changed. +optional +k8s:optional +k8s:update=NoSet +k8s:update=NoUnset

---@class batchv1.JobSpec
---@field activeDeadlineSeconds number Specifies the duration in seconds relative to the startTime that the job may be continuously active before the system tries to terminate it; value must be positive integer. If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again. +optional
---@field backoffLimit number Specifies the number of retries before marking this job failed. Defaults to 6, unless backoffLimitPerIndex (only Indexed Job) is specified. When backoffLimitPerIndex is specified, backoffLimit defaults to 2147483647. +optional
---@field backoffLimitPerIndex number Specifies the limit for the number of retries within an index before marking this index as failed. When enabled the number of failures per index is kept in the pod's batch.kubernetes.io/job-index-failure-count annotation. It can only be set when Job's completionMode=Indexed, and the Pod's restart policy is Never. The field is immutable. +optional
---@field completionMode string completionMode specifies how Pod completions are tracked. It can be `NonIndexed` (default) or `Indexed`. `NonIndexed` means that the Job is considered complete when there have been .spec.completions successfully completed Pods. Each Pod completion is homologous to each other. `Indexed` means that the Pods of a Job get an associated completion index from 0 to (.spec.completions - 1), available in the annotation batch.kubernetes.io/job-completion-index. The Job is considered complete when there is one successfully completed Pod for each index. When value is `Indexed`, .spec.completions must be specified and `.spec.parallelism` must be less than or equal to 10^5. In addition, The Pod name takes the form `$(job-name)-$(index)-$(random-string)`, the Pod hostname takes the form `$(job-name)-$(index)`. More completion modes can be added in the future. If the Job controller observes a mode that it doesn't recognize, which is possible during upgrades due to version skew, the controller skips updates for the Job. +optional
---@field completions number Specifies the desired number of successfully finished pods the job should be run with. Setting to null means that the success of any pod signals the success of all pods, and allows parallelism to have any positive value. Setting to 1 means that parallelism is limited to 1 and the success of that pod signals the success of the job. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/ +optional
---@field managedBy string ManagedBy field indicates the controller that manages a Job. The k8s Job controller reconciles jobs which don't have this field at all or the field value is the reserved string `kubernetes.io/job-controller`, but skips reconciling Jobs with a custom value for this field. The value must be a valid domain-prefixed path (e.g. acme.io/foo) - all characters before the first "/" must be a valid subdomain as defined by RFC 1123. All characters trailing the first "/" must be valid HTTP Path characters as defined by RFC 3986. The value cannot exceed 63 characters. This field is immutable. +optional
---@field manualSelector boolean manualSelector controls generation of pod labels and pod selectors. Leave `manualSelector` unset unless you are certain what you are doing. When false or unset, the system pick labels unique to this job and appends those labels to the pod template. When true, the user is responsible for picking unique labels and specifying the selector. Failure to pick a unique label may cause this and other jobs to not function correctly. However, You may see `manualSelector=true` in jobs that were created with the old `extensions/v1beta1` API. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#specifying-your-own-pod-selector +optional
---@field maxFailedIndexes number Specifies the maximal number of failed indexes before marking the Job as failed, when backoffLimitPerIndex is set. Once the number of failed indexes exceeds this number the entire Job is marked as Failed and its execution is terminated. When left as null the job continues execution of all of its indexes and is marked with the `Complete` Job condition. It can only be specified when backoffLimitPerIndex is set. It can be null or up to completions. It is required and must be less than or equal to 10^4 when is completions greater than 10^5. +optional +k8s:optional +k8s:alpha(since: "1.37")=+k8s:dependentRequired("backoffLimitPerIndex")
---@field parallelism number Specifies the maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/ +optional
---@field podFailurePolicy batchv1.PodFailurePolicy Specifies the policy of handling failed pods. In particular, it allows to specify the set of actions and conditions which need to be satisfied to take the associated action. If empty, the default behaviour applies - the counter of failed pods, represented by the jobs's .status.failed field, is incremented and it is checked against the backoffLimit. This field cannot be used in combination with restartPolicy=OnFailure. +optional
---@field podReplacementPolicy string podReplacementPolicy specifies when to create replacement Pods. Possible values are: - TerminatingOrFailed means that we recreate pods when they are terminating (has a metadata.deletionTimestamp) or failed. - Failed means to wait until a previously created Pod is fully terminated (has phase Failed or Succeeded) before creating a replacement Pod. When using podFailurePolicy, Failed is the the only allowed value. TerminatingOrFailed and Failed are allowed values when podFailurePolicy is not in use. +optional
---@field scheduling batchv1.JobSchedulingConfiguration scheduling defines the Workload-aware Scheduling configuration for this Job. When set, it specifies the scheduling policy (basic or gang), topology constraints, disruption mode, and shared resource claims. When omitted, the Job defaults to the basic scheduling policy, which behaves as standard pod-by-pod scheduling. This field is alpha-level and requires the WorkloadWithJob feature gate. This field is immutable, including whether it is set at all, only policy.gang.minCount may be changed after creation. +featureGate=WorkloadWithJob +optional +k8s:ifDisabled(WorkloadWithJob)=+k8s:forbidden +k8s:optional +k8s:update=NoSet +k8s:update=NoUnset
---@field selector v1.LabelSelector A label query over pods that should match the pod count. Normally, the system sets this field for you. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors +optional
---@field successPolicy batchv1.SuccessPolicy successPolicy specifies the policy when the Job can be declared as succeeded. If empty, the default behavior applies - the Job is declared as succeeded only when the number of succeeded pods equals to the completions. When the field is specified, it must be immutable and works only for the Indexed Jobs. Once the Job meets the SuccessPolicy, the lingering pods are terminated. +optional
---@field suspend boolean suspend specifies whether the Job controller should create Pods or not. If a Job is created with suspend set to true, no Pods are created by the Job controller. If a Job is suspended after creation (i.e. the flag goes from false to true), the Job controller will delete all active Pods associated with this Job. Users must design their workload to gracefully handle this. Suspending a Job will reset the StartTime field of the Job, effectively resetting the ActiveDeadlineSeconds timer too. Defaults to false. +optional
---@field template corev1.PodTemplateSpec Describes the pod that will be created when executing a job. The only allowed template.spec.restartPolicy values are "Never" or "OnFailure". More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
---@field ttlSecondsAfterFinished number ttlSecondsAfterFinished limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted. When the Job is being deleted, its lifecycle guarantees (e.g. finalizers) will be honored. If this field is unset, the Job won't be automatically deleted. If this field is set to zero, the Job becomes eligible to be deleted immediately after it finishes. +optional

---@class batchv1.JobStatus
---@field active number The number of pending and running pods which are not terminating (without a deletionTimestamp). The value is zero for finished jobs. +optional
---@field completedIndexes string completedIndexes holds the completed indexes when .spec.completionMode = "Indexed" in a text format. The indexes are represented as decimal integers separated by commas. The numbers are listed in increasing order. Three or more consecutive numbers are compressed and represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7". +optional
---@field completionTime v1.Time Represents time when the job was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. The completion time is set when the job finishes successfully, and only then. The value cannot be updated or removed. The value indicates the same or later point in time as the startTime field. +optional
---@field conditions batchv1.JobCondition[] The latest available observations of an object's current state. When a Job fails, one of the conditions will have type "Failed" and status true. When a Job is suspended, one of the conditions will have type "Suspended" and status true; when the Job is resumed, the status of this condition will become false. When a Job is completed, one of the conditions will have type "Complete" and status true. A job is considered finished when it is in a terminal condition, either "Complete" or "Failed". A Job cannot have both the "Complete" and "Failed" conditions. Additionally, it cannot be in the "Complete" and "FailureTarget" conditions. The "Complete", "Failed" and "FailureTarget" conditions cannot be disabled. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/ +optional +patchMergeKey=type +patchStrategy=merge +listType=atomic
---@field failed number The number of pods which reached phase Failed. The value increases monotonically. +optional
---@field failedIndexes string FailedIndexes holds the failed indexes when spec.backoffLimitPerIndex is set. The indexes are represented in the text format analogous as for the `completedIndexes` field, ie. they are kept as decimal integers separated by commas. The numbers are listed in increasing order. Three or more consecutive numbers are compressed and represented by the first and last element of the series, separated by a hyphen. For example, if the failed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7". The set of failed indexes cannot overlap with the set of completed indexes. +optional
---@field ready number The number of active pods which have a Ready condition and are not terminating (without a deletionTimestamp).
---@field startTime v1.Time Represents time when the job controller started processing a job. When a Job is created in the suspended state, this field is not set until the first time it is resumed. This field is reset every time a Job is resumed from suspension. It is represented in RFC3339 form and is in UTC. Once set, the field can only be removed when the job is suspended. The field cannot be modified while the job is unsuspended or finished. +optional
---@field succeeded number The number of pods which reached phase Succeeded. The value increases monotonically for a given spec. However, it may decrease in reaction to scale down of elastic indexed jobs. +optional
---@field terminating number The number of pods which are terminating (in phase Pending or Running and have a deletionTimestamp). +optional
---@field uncountedTerminatedPods batchv1.UncountedTerminatedPods uncountedTerminatedPods holds the UIDs of Pods that have terminated but the job controller hasn't yet accounted for in the status counters. The job controller creates pods with a finalizer. When a pod terminates (succeeded or failed), the controller does three steps to account for it in the job status: 1. Add the pod UID to the arrays in this field. 2. Remove the pod finalizer. 3. Remove the pod UID from the arrays while increasing the corresponding counter. Old jobs might not be tracked using this field, in which case the field remains null. The structure is empty for finished jobs. +optional

---@class batchv1.JobTemplateSpec
---@field metadata v1.ObjectMeta Standard object's metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional +k8s:opaqueType
---@field spec batchv1.JobSpec Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class batchv1.PodFailurePolicy
---@field rules batchv1.PodFailurePolicyRule[] A list of pod failure policy rules. The rules are evaluated in order. Once a rule matches a Pod failure, the remaining of the rules are ignored. When no rule matches the Pod failure, the default handling applies - the counter of pod failures is incremented and it is checked against the backoffLimit. At most 20 elements are allowed. +listType=atomic

---@class batchv1.PodFailurePolicyOnExitCodesRequirement
---@field containerName string Restricts the check for exit codes to the container with the specified name. When null, the rule applies to all containers. When specified, it should match one the container or initContainer names in the pod template. +optional
---@field operator string Represents the relationship between the container exit code(s) and the specified values. Containers completed with success (exit code 0) are excluded from the requirement check. Possible values are: - In: the requirement is satisfied if at least one container exit code (might be multiple if there are multiple containers not restricted by the 'containerName' field) is in the set of specified values. - NotIn: the requirement is satisfied if at least one container exit code (might be multiple if there are multiple containers not restricted by the 'containerName' field) is not in the set of specified values. Additional values are considered to be added in the future. Clients should react to an unknown operator by assuming the requirement is not satisfied.
---@field values number[] Specifies the set of values. Each returned container exit code (might be multiple in case of multiple containers) is checked against this set of values with respect to the operator. The list of values must be ordered and must not contain duplicates. Value '0' cannot be used for the In operator. At least one element is required. At most 255 elements are allowed. +listType=set

---@class batchv1.PodFailurePolicyOnPodConditionsPattern
---@field status string Specifies the required Pod condition status. To match a pod condition it is required that the specified status equals the pod condition status. Defaults to True. +optional
---@field type string Specifies the required Pod condition type. To match a pod condition it is required that specified type equals the pod condition type.

---@class batchv1.PodFailurePolicyRule
---@field action string Specifies the action taken on a pod failure when the requirements are satisfied. Possible values are: - FailJob: indicates that the pod's job is marked as Failed and all running pods are terminated. - FailIndex: indicates that the pod's index is marked as Failed and will not be restarted. - Ignore: indicates that the counter towards the .backoffLimit is not incremented and a replacement pod is created. - Count: indicates that the pod is handled in the default way - the counter towards the .backoffLimit is incremented. Additional values are considered to be added in the future. Clients should react to an unknown action by skipping the rule.
---@field onExitCodes batchv1.PodFailurePolicyOnExitCodesRequirement Represents the requirement on the container exit codes. +optional
---@field onPodConditions batchv1.PodFailurePolicyOnPodConditionsPattern[] Represents the requirement on the pod conditions. The requirement is represented as a list of pod condition patterns. The requirement is satisfied if at least one pattern matches an actual pod condition. At most 20 elements are allowed. +listType=atomic +optional

---@class batchv1.SuccessPolicy
---@field rules batchv1.SuccessPolicyRule[] rules represents the list of alternative rules for the declaring the Jobs as successful before `.status.succeeded >= .spec.completions`. Once any of the rules are met, the "SuccessCriteriaMet" condition is added, and the lingering pods are removed. The terminal state for such a Job has the "Complete" condition. Additionally, these rules are evaluated in order; Once the Job meets one of the rules, other rules are ignored. At most 20 elements are allowed. +listType=atomic

---@class batchv1.SuccessPolicyRule
---@field succeededCount number succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is "1-4", succeededCount is "3", and completed indexes are "1", "3", and "5", the Job isn't declared as succeeded because only "1" and "3" indexes are considered in that rules. When this field is null, this doesn't default to any value and is never evaluated at any time. When specified it needs to be a positive integer. +optional
---@field succeededIndexes string succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to ".spec.completions-1" and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7". When this field is null, this field doesn't default to any value and is never evaluated at any time. +optional

---@class batchv1.UncountedTerminatedPods
---@field failed string[] failed holds UIDs of failed Pods. +listType=set +optional
---@field succeeded string[] succeeded holds UIDs of succeeded Pods. +listType=set +optional

---@class coordinationv1.Lease
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec coordinationv1.LeaseSpec spec contains the specification of the Lease. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class coordinationv1.LeaseList
---@field items coordinationv1.Lease[] items is a list of schema objects.
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class coordinationv1.LeaseSpec
---@field acquireTime v1.MicroTime acquireTime is a time when the current lease was acquired. +optional
---@field holderIdentity string holderIdentity contains the identity of the holder of a current lease. If Coordinated Leader Election is used, the holder identity must be equal to the elected LeaseCandidate.metadata.name field. +optional
---@field leaseDurationSeconds number leaseDurationSeconds is a duration that candidates for a lease need to wait to force acquire it. This is measured against the time of last observed renewTime. +optional
---@field leaseTransitions number leaseTransitions is the number of transitions of a lease between holders. +optional
---@field preferredHolder string preferredHolder signals to a lease holder that the lease has a more optimal holder and should be given up. This field can only be set if Strategy is also set. +featureGate=CoordinatedLeaderElection +optional
---@field renewTime v1.MicroTime renewTime is a time when the current holder of a lease has last updated the lease. +optional
---@field strategy string strategy indicates the strategy for picking the leader for coordinated leader election. If the field is not specified, there is no active coordination for this lease. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled. +featureGate=CoordinatedLeaderElection +optional

---@class corev1.Affinity
---@field nodeAffinity corev1.NodeAffinity Describes node affinity scheduling rules for the pod. +optional
---@field podAffinity corev1.PodAffinity Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)). +optional
---@field podAntiAffinity corev1.PodAntiAffinity Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)). +optional

---@class corev1.AppArmorProfile
---@field localhostProfile string localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost". +optional
---@field type string type indicates which kind of AppArmor profile will be applied. Valid options are: Localhost - a profile pre-loaded on the node. RuntimeDefault - the container runtime's default profile. Unconfined - no AppArmor enforcement. +unionDiscriminator

---@class corev1.AttachedVolume
---@field devicePath string DevicePath represents the device path where the volume should be available
---@field name string Name of the attached volume

---@class corev1.Capabilities
---@field add string[] Added capabilities +optional +listType=atomic
---@field drop string[] Removed capabilities +optional +listType=atomic

---@class corev1.ClientIPConfig
---@field timeoutSeconds number timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be >0 && <=86400(for 1 day) if ServiceAffinity == "ClientIP". Default value is 10800(for 3 hours). +optional

---@class corev1.ConfigMap
---@field binaryData table<string, number[]> BinaryData contains the binary data. Each key must consist of alphanumeric characters, '-', '_' or '.'. BinaryData can contain byte sequences that are not in the UTF-8 range. The keys stored in BinaryData must not overlap with the ones in the Data field, this is enforced during validation process. Using this field will require 1.10+ apiserver and kubelet. Note: BinaryData keys are not currently propagated to container env vars via ConfigMapKeyRef or ConfigMapRef env sources; only Data keys are used. +optional
---@field data table<string, string> Data contains the configuration data. Each key must consist of alphanumeric characters, '-', '_' or '.'. Values with non-UTF-8 byte sequences must use the BinaryData field. The keys stored in Data must not overlap with the keys in the BinaryData field, this is enforced during validation process. +optional
---@field immutable boolean Immutable, if set to true, ensures that data stored in the ConfigMap cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil. +optional
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class corev1.ConfigMapEnvSource
---@field optional boolean Specify whether the ConfigMap must be defined +optional

---@class corev1.ConfigMapKeySelector
---@field key string The key to select from the ConfigMap's Data field. Keys in the BinaryData field are not currently propagated to container env vars.
---@field optional boolean Specify whether the ConfigMap or its key must be defined +optional

---@class corev1.ConfigMapList
---@field items corev1.ConfigMap[] Items is the list of ConfigMaps.
---@field metadata v1.ListMeta More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class corev1.ConfigMapNodeConfigSource
---@field kubeletConfigKey string KubeletConfigKey declares which key of the referenced ConfigMap corresponds to the KubeletConfiguration structure This field is required in all cases.
---@field name string Name is the metadata.name of the referenced ConfigMap. This field is required in all cases.
---@field namespace string Namespace is the metadata.namespace of the referenced ConfigMap. This field is required in all cases.
---@field resourceVersion string ResourceVersion is the metadata.ResourceVersion of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status. +optional
---@field uid string UID is the metadata.UID of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status. +optional

---@class corev1.Container
---@field args string[] Arguments to the entrypoint. The container image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell +optional +listType=atomic
---@field command string[] Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell +optional +listType=atomic
---@field env corev1.EnvVar[] List of environment variables to set in the container. Cannot be updated. +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field envFrom corev1.EnvFromSource[] List of sources to populate environment variables in the container. The keys defined within a source may consist of any printable ASCII characters except '='. When a key exists in multiple sources, the value associated with the last source will take precedence. Values defined by an Env with a duplicate key will take precedence. Cannot be updated. +optional +listType=atomic
---@field image string Container image name. More info: https://kubernetes.io/docs/concepts/containers/images This field is optional to allow higher level config management to default or override container images in workload controllers like Deployments and StatefulSets. +optional
---@field imagePullPolicy string Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images +optional
---@field lifecycle corev1.Lifecycle Actions that the management system should take in response to container lifecycle events. Cannot be updated. +optional
---@field livenessProbe corev1.Probe Periodic probe of container liveness. Container will be restarted if the probe fails. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes +optional
---@field name string Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). Cannot be updated.
---@field ports corev1.ContainerPort[] List of ports to expose from the container. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Modifying this array with strategic merge patch may corrupt the data. For more information See https://github.com/kubernetes/kubernetes/issues/108255. Cannot be updated. +optional +patchMergeKey=containerPort +patchStrategy=merge +listType=map +listMapKey=containerPort +listMapKey=protocol
---@field readinessProbe corev1.Probe Periodic probe of container service readiness. Container will be removed from service endpoints if the probe fails. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes +optional
---@field resizePolicy corev1.ContainerResizePolicy[] Resources resize policy for the container. This field cannot be set on ephemeral containers. +featureGate=InPlacePodVerticalScaling +optional +listType=atomic
---@field resources corev1.ResourceRequirements Compute Resources required by this container. Cannot be updated. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ +optional
---@field restartPolicy string RestartPolicy defines the restart behavior of individual containers in a pod. This overrides the pod-level restart policy. When this field is not specified, the restart behavior is defined by the Pod's restart policy and the container type. Additionally, setting the RestartPolicy as "Always" for the init container will have the following effect: this init container will be continually restarted on exit until all regular containers have terminated. Once all regular containers have completed, all init containers with restartPolicy "Always" will be shut down. This lifecycle differs from normal init containers and is often referred to as a "sidecar" container. Although this init container still starts in the init container sequence, it does not wait for the container to complete before proceeding to the next init container. Instead, the next init container starts immediately after this init container is started, or after any startupProbe has successfully completed. +optional
---@field restartPolicyRules corev1.ContainerRestartRule[] Represents a list of rules to be checked to determine if the container should be restarted on exit. The rules are evaluated in order. Once a rule matches a container exit condition, the remaining rules are ignored. If no rule matches the container exit condition, the Container-level restart policy determines the whether the container is restarted or not. Constraints on the rules: - At most 20 rules are allowed. - Rules can have the same action. - Identical rules are not forbidden in validations. When rules are specified, container MUST set RestartPolicy explicitly even it if matches the Pod's RestartPolicy. +featureGate=ContainerRestartRules +optional +listType=atomic
---@field securityContext corev1.SecurityContext SecurityContext defines the security options the container should be run with. If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext. More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ +optional
---@field startupProbe corev1.Probe StartupProbe indicates that the Pod has successfully initialized. If specified, no other probes are executed until this completes successfully. If this probe fails, the Pod will be restarted, just as if the livenessProbe failed. This can be used to provide different probe parameters at the beginning of a Pod's lifecycle, when it might take a long time to load data or warm a cache, than during steady-state operation. This cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes +optional
---@field stdin boolean Whether this container should allocate a buffer for stdin in the container runtime. If this is not set, reads from stdin in the container will always result in EOF. Default is false. +optional
---@field stdinOnce boolean Whether the container runtime should close the stdin channel after it has been opened by a single attach. When stdin is true the stdin stream will remain open across multiple attach sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the first client attaches to stdin, and then remains open and accepts data until the client disconnects, at which time stdin is closed and remains closed until the container is restarted. If this flag is false, a container processes that reads from stdin will never receive an EOF. Default is false +optional
---@field terminationMessagePath string Optional: Path at which the file to which the container's termination message will be written is mounted into the container's filesystem. Message written is intended to be brief final status, such as an assertion failure message. Will be truncated by the node if greater than 4096 bytes. The total message length across all containers will be limited to 12kb. Defaults to /dev/termination-log. Cannot be updated. +optional
---@field terminationMessagePolicy string Indicate how the termination message should be populated. File will use the contents of terminationMessagePath to populate the container status message on both success and failure. FallbackToLogsOnError will use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be updated. +optional
---@field tty boolean Whether this container should allocate a TTY for itself, also requires 'stdin' to be true. Default is false. +optional
---@field volumeDevices corev1.VolumeDevice[] volumeDevices is the list of block devices to be used by the container. +patchMergeKey=devicePath +patchStrategy=merge +listType=map +listMapKey=devicePath +optional
---@field volumeMounts corev1.VolumeMount[] Pod volumes to mount into the container's filesystem. Cannot be updated. +optional +patchMergeKey=mountPath +patchStrategy=merge +listType=map +listMapKey=mountPath
---@field workingDir string Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated. +optional

---@class corev1.ContainerExtendedResourceRequest
---@field containerName string The name of the container requesting resources.
---@field requestName string The name of the request in the special ResourceClaim which corresponds to the extended resource.
---@field resourceName string The name of the extended resource in that container which gets backed by DRA.

---@class corev1.ContainerImage
---@field names string[] Names by which this image is known. e.g. ["kubernetes.example/hyperkube:v1.0.7", "cloud-vendor.registry.example/cloud-vendor/hyperkube:v1.0.7"] +optional +listType=atomic
---@field sizeBytes number The size of the image in bytes. +optional

---@class corev1.ContainerPort
---@field containerPort number Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
---@field hostIP string What host IP to bind the external port to. +optional
---@field hostPort number Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this. +optional
---@field name string If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services. +optional
---@field protocol string Protocol for port. Must be UDP, TCP, or SCTP. Defaults to "TCP". +optional +default="TCP"

---@class corev1.ContainerResizePolicy
---@field resourceName string Name of the resource to which this resource resize policy applies. Supported values: cpu, memory.
---@field restartPolicy string Restart policy to apply when specified resource is resized. If not specified, it defaults to NotRequired.

---@class corev1.ContainerRestartRule
---@field action string Specifies the action taken on a container exit if the requirements are satisfied. The only possible value is "Restart" to restart the container. +required
---@field exitCodes corev1.ContainerRestartRuleOnExitCodes Represents the exit codes to check on container exits. +optional +oneOf=when

---@class corev1.ContainerRestartRuleOnExitCodes
---@field operator string Represents the relationship between the container exit code(s) and the specified values. Possible values are: - In: the requirement is satisfied if the container exit code is in the set of specified values. - NotIn: the requirement is satisfied if the container exit code is not in the set of specified values. +required
---@field values number[] Specifies the set of values to check for container exit codes. At most 255 elements are allowed. +optional +listType=set

---@class corev1.ContainerState
---@field running corev1.ContainerStateRunning Details about a running container +optional
---@field terminated corev1.ContainerStateTerminated Details about a terminated container +optional
---@field waiting corev1.ContainerStateWaiting Details about a waiting container +optional

---@class corev1.ContainerStateRunning
---@field startedAt v1.Time Time at which the container was last (re-)started +optional

---@class corev1.ContainerStateTerminated
---@field containerID string Container's ID in the format '<type>://<container_id>' +optional
---@field exitCode number Exit status from the last termination of the container
---@field finishedAt v1.Time Time at which the container last terminated +optional
---@field message string Message regarding the last termination of the container +optional
---@field reason string (brief) reason from the last termination of the container +optional
---@field signal number Signal from the last termination of the container +optional
---@field startedAt v1.Time Time at which previous execution of the container started +optional

---@class corev1.ContainerStateWaiting
---@field message string Message regarding why the container is not yet running. +optional
---@field reason string (brief) reason the container is not yet running. +optional

---@class corev1.ContainerStatus
---@field allocatedResources table<string, resource.Quantity> AllocatedResources represents the compute resources allocated for this container by the node. Kubelet sets this value to Container.Resources.Requests upon successful pod admission and after successfully admitting desired pod resize. +optional
---@field allocatedResourcesStatus corev1.ResourceStatus[] AllocatedResourcesStatus represents the status of various resources allocated for this Pod. +featureGate=ResourceHealthStatus +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field containerID string ContainerID is the ID of the container in the format '<type>://<container_id>'. Where type is a container runtime identifier, returned from Version call of CRI API (for example "containerd"). +optional
---@field image string Image is the name of container image that the container is running. The container image may not match the image used in the PodSpec, as it may have been resolved by the runtime. More info: https://kubernetes.io/docs/concepts/containers/images.
---@field imageID string ImageID is the image ID of the container's image. The image ID may not match the image ID of the image used in the PodSpec, as it may have been resolved by the runtime.
---@field lastState corev1.ContainerState LastTerminationState holds the last termination state of the container to help debug container crashes and restarts. This field is not populated if the container is still running and RestartCount is 0. +optional
---@field name string Name is a DNS_LABEL representing the unique name of the container. Each container in a pod must have a unique name across all container types. Cannot be updated.
---@field ready boolean Ready specifies whether the container is currently passing its readiness check. The value will change as readiness probes keep executing. If no readiness probes are specified, this field defaults to true once the container is fully started (see Started field). The value is typically used to determine whether a container is ready to accept traffic.
---@field resources corev1.ResourceRequirements Resources represents the compute resource requests and limits that have been successfully enacted on the running container after it has been started or has been successfully resized. +featureGate=InPlacePodVerticalScaling +optional
---@field restartCount number RestartCount holds the number of times the container has been restarted. Kubelet makes an effort to always increment the value, but there are cases when the state may be lost due to node restarts and then the value may be reset to 0. The value is never negative.
---@field started boolean Started indicates whether the container has finished its postStart lifecycle hook and passed its startup probe. Initialized as false, becomes true after startupProbe is considered successful. Resets to false when the container is restarted, or if kubelet loses state temporarily. In both cases, startup probes will run again. Is always true when no startupProbe is defined and container is running and has passed the postStart lifecycle hook. The null value must be treated the same as false. +optional
---@field state corev1.ContainerState State holds details about the container's current condition. +optional
---@field stopSignal string StopSignal reports the effective stop signal for this container +featureGate=ContainerStopSignals +optional
---@field user corev1.ContainerUser User represents user identity information initially attached to the first process of the container +featureGate=SupplementalGroupsPolicy +optional
---@field volumeMounts corev1.VolumeMountStatus[] Status of volume mounts. +optional +patchMergeKey=mountPath +patchStrategy=merge +listType=map +listMapKey=mountPath

---@class corev1.ContainerUser
---@field linux corev1.LinuxContainerUser Linux holds user identity information initially attached to the first process of the containers in Linux. Note that the actual running identity can be changed if the process has enough privilege to do so. +optional

---@class corev1.DaemonEndpoint
---@field Port number Port number of the given endpoint.

---@class corev1.EnvFromSource
---@field configMapRef corev1.ConfigMapEnvSource The ConfigMap to select from +optional
---@field prefix string Optional text to prepend to the name of each environment variable. May consist of any printable ASCII characters except '='. +optional
---@field secretRef corev1.SecretEnvSource The Secret to select from +optional

---@class corev1.EnvVar
---@field name string Name of the environment variable. May consist of any printable ASCII characters except '='.
---@field value string Variable references $(VAR_NAME) are expanded using the previously defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "". +optional
---@field valueFrom corev1.EnvVarSource Source for the environment variable's value. Cannot be used if value is not empty. +optional

---@class corev1.EnvVarSource
---@field configMapKeyRef corev1.ConfigMapKeySelector Selects a key of a ConfigMap. +optional
---@field fieldRef corev1.ObjectFieldSelector Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`, spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs. +optional
---@field fileKeyRef corev1.FileKeySelector FileKeyRef selects a key of the env file. Requires the EnvFiles feature gate to be enabled. +featureGate=EnvFiles +optional
---@field resourceFieldRef corev1.ResourceFieldSelector Selects a resource of the container: only resources limits and requests (limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported. +optional
---@field secretKeyRef corev1.SecretKeySelector Selects a key of a secret in the pod's namespace +optional

---@class corev1.EphemeralContainer
---@field targetContainerName string If set, the name of the container from PodSpec that this ephemeral container targets. The ephemeral container will be run in the namespaces (IPC, PID, etc) of this container. If not set then the ephemeral container uses the namespaces configured in the Pod spec. The container runtime must implement support for this feature. If the runtime does not support namespace targeting then the result of setting this field is undefined. +optional

---@class corev1.EventSource
---@field component string Component from which the event is generated. +optional
---@field host string Node name on which the event is generated. +optional

---@class corev1.EvictionResponder
---@field name string name allows you to identify the responder responding to the Eviction. It must be a valid domain-prefixed key (such as "acme.io/foo"). Domain names *.k8s.io and *.kubernetes.io are reserved. This field must be unique for each responder. This field is required. +required +k8s:required +k8s:format=k8s-prefixed-label-key +k8s:customValidation
---@field priority number priority for this responder. Higher priorities are selected first by the evictionrequest-controller. If there are responders with the same priority, the responder whose domain name comes first in the alphabetical higher domain order, will be picked. This means that the top domain labels are compared alphabetically first, followed by the lower domain labels. The key is compared last. The responder that is the managing controller of the pod should set the value of this field to 10000 to allow both for preemption or fallback registration by other responders. The minimum value is 0 and the maximum value is 100000. The interval 0-999 is reserved for responders with *.k8s.io suffix. This field is required. +required +k8s:required +k8s:minimum=0 +k8s:maximum=100000 +k8s:customValidation

---@class corev1.ExecAction
---@field command string[] Command is the command line to execute inside the container, the working directory for the command is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy. +optional +listType=atomic

---@class corev1.FileKeySelector
---@field key string The key within the env file. An invalid key will prevent the pod from starting. The keys defined within a source may consist of any printable ASCII characters except '='. During Alpha stage of the EnvFiles feature gate, the key size is limited to 128 characters. +required
---@field optional boolean Specify whether the file or its key must be defined. If the file or key does not exist, then the env var is not published. If optional is set to true and the specified key does not exist, the environment variable will not be set in the Pod's containers. If optional is set to false and the specified key does not exist, an error will be returned during Pod creation. +optional +default=false
---@field path string The path within the volume from which to select the file. Must be relative and may not contain the '..' path or start with '..'. +required
---@field volumeName string The name of the volume mount containing the env file. +required

---@class corev1.HTTPGetAction
---@field host string Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead. +optional
---@field httpHeaders corev1.HTTPHeader[] Custom headers to set in the request. HTTP allows repeated headers. +optional +listType=atomic
---@field path string Path to access on the HTTP server. +optional
---@field port intstr.IntOrString Name or number of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
---@field protocol string Protocol selects the wire protocol for the probe connection. Nil defaults to HTTP/1.1. +optional +featureGate=H2CContainerProbe
---@field scheme string Scheme to use for connecting to the host. Defaults to HTTP. +optional

---@class corev1.HTTPHeader
---@field name string The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
---@field value string The header field value

---@class corev1.HostAlias
---@field hostnames string[] Hostnames for the above IP address. +listType=atomic
---@field ip string IP address of the host file entry. +required

---@class corev1.HostIP
---@field ip string IP is the IP address assigned to the host +required

---@class corev1.ImageVolumeStatus
---@field imageRef string ImageRef is the digest of the image used for this volume. It should have a value that's similar to the pod's status.containerStatuses[i].imageID. The ImageRef length should not exceed 256 characters. +kubebuilder:validation:MaxLength=256 +required

---@class corev1.Lifecycle
---@field postStart corev1.LifecycleHandler PostStart is called immediately after a container is created. If the handler fails, the container is terminated and restarted according to its restart policy. Other management of the container blocks until the hook completes. More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks +optional
---@field preStop corev1.LifecycleHandler PreStop is called immediately before a container is terminated due to an API request or management event such as liveness/startup probe failure, preemption, resource contention, etc. The handler is not called if the container crashes or exits. The Pod's termination grace period countdown begins before the PreStop hook is executed. Regardless of the outcome of the handler, the container will eventually terminate within the Pod's termination grace period (unless delayed by finalizers). Other management of the container blocks until the hook completes or until the termination grace period is reached. More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks +optional
---@field stopSignal string StopSignal defines which signal will be sent to a container when it is being stopped. If not specified, the default is defined by the container runtime in use. StopSignal can only be set for Pods with a non-empty .spec.os.name +optional

---@class corev1.LifecycleHandler
---@field exec corev1.ExecAction Exec specifies a command to execute in the container. +optional
---@field httpGet corev1.HTTPGetAction HTTPGet specifies an HTTP GET request to perform. +optional
---@field sleep corev1.SleepAction Sleep represents a duration that the container should sleep. +optional
---@field tcpSocket corev1.TCPSocketAction Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for backward compatibility. There is no validation of this field and lifecycle hooks will fail at runtime when it is specified. +optional

---@class corev1.LinuxContainerUser
---@field gid number GID is the primary gid initially attached to the first process in the container
---@field supplementalGroups number[] SupplementalGroups are the supplemental groups initially attached to the first process in the container +optional +listType=atomic
---@field uid number UID is the primary uid initially attached to the first process in the container

---@class corev1.LoadBalancerIngress
---@field hostname string Hostname is set for load-balancer ingress points that are DNS based (typically AWS load-balancers) +optional
---@field ip string IP is set for load-balancer ingress points that are IP based (typically GCE or OpenStack load-balancers) +optional
---@field ipMode string IPMode specifies how the load-balancer IP behaves, and may only be specified when the ip field is specified. Setting this to "VIP" indicates that traffic is delivered to the node with the destination set to the load-balancer's IP and port. Setting this to "Proxy" indicates that traffic is delivered to the node or pod with the destination set to the node's IP and node port or the pod's IP and port. Service implementations may use this information to adjust traffic routing. +optional
---@field ports corev1.PortStatus[] Ports is a list of records of service ports If used, every port defined in the service should have an entry in it +listType=atomic +optional

---@class corev1.LoadBalancerStatus
---@field ingress corev1.LoadBalancerIngress[] Ingress is a list containing ingress points for the load-balancer. Traffic intended for the service should be sent to these ingress points. +optional +listType=atomic

---@class corev1.LocalObjectReference
---@field name string Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names +optional +default="" +kubebuilder:default="" TODO: Drop `kubebuilder:default` when controller-gen doesn't need it https://github.com/kubernetes-sigs/kubebuilder/issues/3896.

---@class corev1.ModifyVolumeStatus
---@field status string status is the status of the ControllerModifyVolume operation. It can be in any of following states: - Pending Pending indicates that the PersistentVolumeClaim cannot be modified due to unmet requirements, such as the specified VolumeAttributesClass not existing. - InProgress InProgress indicates that the volume is being modified. - Infeasible Infeasible indicates that the request has been rejected as invalid by the CSI driver. To resolve the error, a valid VolumeAttributesClass needs to be specified. Note: New statuses can be added in the future. Consumers should check for unknown statuses and fail appropriately.
---@field targetVolumeAttributesClassName string targetVolumeAttributesClassName is the name of the VolumeAttributesClass the PVC currently being reconciled

---@class corev1.Namespace
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.NamespaceSpec Spec defines the behavior of the Namespace. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status corev1.NamespaceStatus Status describes the current status of a Namespace. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class corev1.NamespaceCondition
---@field lastTransitionTime v1.Time Last time the condition transitioned from one status to another. +optional
---@field message string Human-readable message indicating details about last transition. +optional
---@field reason string Unique, one-word, CamelCase reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown.
---@field type string Type of namespace controller condition.

---@class corev1.NamespaceList
---@field items corev1.Namespace[] Items is the list of Namespace objects in the list. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.NamespaceSpec
---@field finalizers string[] Finalizers is an opaque list of values that must be empty to permanently remove object from storage. More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/ +optional +listType=atomic

---@class corev1.NamespaceStatus
---@field conditions corev1.NamespaceCondition[] Represents the latest available observations of a namespace's current state. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field phase string Phase is the current lifecycle phase of the namespace. More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/ +optional

---@class corev1.Node
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.NodeSpec Spec defines the behavior of a node. https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status corev1.NodeStatus Most recently observed status of the node. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class corev1.NodeAddress
---@field address string The node address.
---@field type string Node address type, one of Hostname, ExternalIP or InternalIP.

---@class corev1.NodeAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.PreferredSchedulingTerm[] The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node matches the corresponding matchExpressions; the node(s) with the highest sum are the most preferred. +optional +listType=atomic
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.NodeSelector If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to an update), the system may or may not try to eventually evict the pod from its node. +optional

---@class corev1.NodeAllocatableMappedResources
---@field name string Name is the name of the resource (e.g., cpu, memory). +required +k8s:required
---@field quantity resource.Quantity Quantity is the total node allocatable resource capacity allocated for the claim. This claim's allocated devices is shared by all the containers referencing the claim. Kubelet adds this value to both requests and limits at the pod-level cgroup, and to limits at the container-level cgroup for each container referencing the claim. +required +k8s:required

---@class corev1.NodeAllocatableOverheadResources
---@field name string Name is the name of the resource (e.g., cpu, memory). +required +k8s:required
---@field perContainer resource.Quantity PerContainer is the variable overhead quantity applied for each container referencing the claim. The container references are recorded in `nodeAllocatableResourceClaimStatuses.containers`. The total overhead quantity allocated for the claim is computed as: Quantity = PerPod + (PerContainer * NumReferences) Kubelet accounts for this overhead in cgroups: - Pod-level cgroup (requests and limits): Kubelet adds PerPod + (PerContainer * NumReferences). - Container-level cgroup (limits only): Kubelet adds PerPod + PerContainer for each referencing container. This allows any single container to access the pod-level overhead, while the parent cgroup caps the total usage to account for PerPod exactly once. At least one of PerPod or PerContainer must be specified. Specifying neither is an invalid configuration. +optional +k8s:optional
---@field perPod resource.Quantity PerPod is the flat overhead quantity allocated per pod. Adding to each container limit allows individual containers to utilize the overhead, while the parent pod-level cgroup limit caps the total usage at the pod boundary where the overhead is accounted for exactly once. At least one of PerPod or PerContainer must be specified. Specifying neither is an invalid configuration. +optional +k8s:optional

---@class corev1.NodeAllocatableResourceClaimStatus
---@field containers string[] Containers lists the names of all containers in this pod that reference the claim. +optional +listType=set +k8s:optional +k8s:listType=set
---@field mapping corev1.NodeAllocatableMappedResources[] Mapping contains allocations through devices mapped in the device spec's `nodeAllocatableResources[...].mapping` field. This is used by kubelet for pod level and container-level cgroup enforcement. +optional +patchStrategy=merge +patchMergeKey=name +listType=map +listMapKey=name +k8s:optional +k8s:listType=map +k8s:listMapKey=name
---@field overhead corev1.NodeAllocatableOverheadResources[] Overhead contains allocations through devices mapped in the device spec's `nodeAllocatableResources[...].overhead` field. This is used by kubelet for pod level and container-level cgroup enforcement. +optional +patchStrategy=merge +patchMergeKey=name +listType=map +listMapKey=name +k8s:optional +k8s:listType=map +k8s:listMapKey=name
---@field resourceClaimName string ResourceClaimName is the resource claim referenced by the pod that resulted in this node allocatable resource allocation. +required +k8s:required

---@class corev1.NodeCondition
---@field lastHeartbeatTime v1.Time Last time we got an update on a given condition. +optional
---@field lastTransitionTime v1.Time Last time the condition transit from one status to another. +optional
---@field message string Human readable message indicating details about last transition. +optional
---@field reason string (brief) reason for the condition's last transition. +optional
---@field status string Status of the condition, one of True, False, Unknown.
---@field type string Type of node condition.

---@class corev1.NodeConfigSource
---@field configMap corev1.ConfigMapNodeConfigSource ConfigMap is a reference to a Node's ConfigMap

---@class corev1.NodeConfigStatus
---@field active corev1.NodeConfigSource Active reports the checkpointed config the node is actively using. Active will represent either the current version of the Assigned config, or the current LastKnownGood config, depending on whether attempting to use the Assigned config results in an error. +optional
---@field assigned corev1.NodeConfigSource Assigned reports the checkpointed config the node will try to use. When Node.Spec.ConfigSource is updated, the node checkpoints the associated config payload to local disk, along with a record indicating intended config. The node refers to this record to choose its config checkpoint, and reports this record in Assigned. Assigned only updates in the status after the record has been checkpointed to disk. When the Kubelet is restarted, it tries to make the Assigned config the Active config by loading and validating the checkpointed payload identified by Assigned. +optional
---@field error string Error describes any problems reconciling the Spec.ConfigSource to the Active config. Errors may occur, for example, attempting to checkpoint Spec.ConfigSource to the local Assigned record, attempting to checkpoint the payload associated with Spec.ConfigSource, attempting to load or validate the Assigned config, etc. Errors may occur at different points while syncing config. Earlier errors (e.g. download or checkpointing errors) will not result in a rollback to LastKnownGood, and may resolve across Kubelet retries. Later errors (e.g. loading or validating a checkpointed config) will result in a rollback to LastKnownGood. In the latter case, it is usually possible to resolve the error by fixing the config assigned in Spec.ConfigSource. You can find additional information for debugging by searching the error message in the Kubelet log. Error is a human-readable description of the error state; machines can check whether or not Error is empty, but should not rely on the stability of the Error text across Kubelet versions. +optional
---@field lastKnownGood corev1.NodeConfigSource LastKnownGood reports the checkpointed config the node will fall back to when it encounters an error attempting to use the Assigned config. The Assigned config becomes the LastKnownGood config when the node determines that the Assigned config is stable and correct. This is currently implemented as a 10-minute soak period starting when the local record of Assigned config is updated. If the Assigned config is Active at the end of this period, it becomes the LastKnownGood. Note that if Spec.ConfigSource is reset to nil (use local defaults), the LastKnownGood is also immediately reset to nil, because the local default config is always assumed good. You should not make assumptions about the node's method of determining config stability and correctness, as this may change or become configurable in the future. +optional

---@class corev1.NodeDaemonEndpoints
---@field kubeletEndpoint corev1.DaemonEndpoint Endpoint on which Kubelet is listening. +optional

---@class corev1.NodeFeatures
---@field supplementalGroupsPolicy boolean SupplementalGroupsPolicy is set to true if the runtime supports SupplementalGroupsPolicy and ContainerUser. +optional

---@class corev1.NodeList
---@field items corev1.Node[] List of nodes
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.NodePodPreemptionPolicy
---@field disableResizePreemption string[] DisableResizePreemption lists the owners (e.g., autoscalers, operators, administrators) that have requested to disable scheduler and Kubelet preemption for in-place pod resize on this node. If this list is non-empty, resize-induced preemption is disabled on this node. This is an alpha field and requires enabling the InPlacePodVerticalScalingSchedulerPreemption feature gate. +listType=set +k8s:listType=set +optional +k8s:maxItems=20 +k8s:optional +k8s:eachVal=+k8s:format=k8s-label-key

---@class corev1.NodeRuntimeHandler
---@field features corev1.NodeRuntimeHandlerFeatures Supported features. +optional
---@field name string Runtime handler name. Empty for the default runtime handler. +optional

---@class corev1.NodeRuntimeHandlerFeatures
---@field recursiveReadOnlyMounts boolean RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts. +optional
---@field userNamespaces boolean UserNamespaces is set to true if the runtime handler supports UserNamespaces, including for volumes. +optional

---@class corev1.NodeSelector
---@field nodeSelectorTerms corev1.NodeSelectorTerm[] Required. A list of node selector terms. The terms are ORed. +listType=atomic

---@class corev1.NodeSelectorRequirement
---@field key string The label key that the selector applies to.
---@field operator string Represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
---@field values string[] An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch. +optional +listType=atomic

---@class corev1.NodeSelectorTerm
---@field matchExpressions corev1.NodeSelectorRequirement[] A list of node selector requirements by node's labels. +optional +listType=atomic
---@field matchFields corev1.NodeSelectorRequirement[] A list of node selector requirements by node's fields. +optional +listType=atomic

---@class corev1.NodeSpec
---@field configSource corev1.NodeConfigSource Deprecated: Previously used to specify the source of the node's configuration for the DynamicKubeletConfig feature. This feature is removed. +optional
---@field externalID string Deprecated. Not all kubelets will set this field. Remove field after 1.13. see: https://issues.k8s.io/61966 +optional
---@field podCIDR string PodCIDR represents the pod IP range assigned to the node. +optional
---@field podCIDRs string[] podCIDRs represents the IP ranges assigned to the node for usage by Pods on that node. If this field is specified, the 0th entry must match the podCIDR field. It may contain at most 1 value for each of IPv4 and IPv6. +optional +patchStrategy=merge +listType=set
---@field podPreemptionPolicy corev1.NodePodPreemptionPolicy PodPreemptionPolicy controls the node-level preemption behaviors for pods on this node. This is an alpha field and requires enabling the InPlacePodVerticalScalingSchedulerPreemption feature gate. +featureGate=InPlacePodVerticalScalingSchedulerPreemption +optional +k8s:optional +k8s:ifDisabled(InPlacePodVerticalScalingSchedulerPreemption)=+k8s:forbidden
---@field providerID string ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID> +optional +k8s:alpha(since: "1.36")=+k8s:optional +k8s:alpha(since: "1.36")=+k8s:update=NoModify +k8s:alpha(since: "1.36")=+k8s:update=NoUnset
---@field taints corev1.Taint[] If specified, the node's taints. +optional +listType=atomic
---@field unschedulable boolean Unschedulable controls node schedulability of new pods. By default, node is schedulable. More info: https://kubernetes.io/docs/concepts/nodes/node/#manual-node-administration +optional

---@class corev1.NodeStatus
---@field addresses corev1.NodeAddress[] List of addresses reachable to the node. Queried from cloud provider, if available. More info: https://kubernetes.io/docs/reference/node/node-status/#addresses Note: This field is declared as mergeable, but the merge key is not sufficiently unique, which can cause data corruption when it is merged. Callers should instead use a full-replacement patch. See https://pr.k8s.io/79391 for an example. Consumers should assume that addresses can change during the lifetime of a Node. However, there are some exceptions where this may not be possible, such as Pods that inherit a Node's address in its own status or consumers of the downward API (status.hostIP). +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field allocatable table<string, resource.Quantity> Allocatable represents the resources of a node that are available for scheduling. Defaults to Capacity. +optional
---@field capacity table<string, resource.Quantity> Capacity represents the total resources of a node. More info: https://kubernetes.io/docs/reference/node/node-status/#capacity +optional
---@field conditions corev1.NodeCondition[] Conditions is an array of current observed node conditions. More info: https://kubernetes.io/docs/reference/node/node-status/#condition +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field config corev1.NodeConfigStatus Status of the config assigned to the node via the dynamic Kubelet config feature. +optional
---@field daemonEndpoints corev1.NodeDaemonEndpoints Endpoints of daemons running on the Node. +optional
---@field declaredFeatures string[] DeclaredFeatures represents the features related to feature gates that are declared by the node. +featureGate=NodeDeclaredFeatures +optional +listType=atomic
---@field features corev1.NodeFeatures Features describes the set of features implemented by the CRI implementation. +featureGate=SupplementalGroupsPolicy +optional
---@field images corev1.ContainerImage[] List of container images on this node +optional +listType=atomic
---@field nodeInfo corev1.NodeSystemInfo Set of ids/uuids to uniquely identify the node. More info: https://kubernetes.io/docs/reference/node/node-status/#info +optional
---@field phase string NodePhase is the recently observed lifecycle phase of the node. More info: https://kubernetes.io/docs/concepts/nodes/node/#phase The field is never populated, and now is deprecated. +optional
---@field runtimeHandlers corev1.NodeRuntimeHandler[] The available runtime handlers. +optional +listType=atomic
---@field volumesAttached corev1.AttachedVolume[] List of volumes that are attached to the node. +optional +listType=atomic
---@field volumesInUse string[] List of attachable volumes in use (mounted) by the node. +optional +listType=atomic

---@class corev1.NodeSwapStatus
---@field capacity number Total amount of swap memory in bytes. +optional

---@class corev1.NodeSystemInfo
---@field architecture string The Architecture reported by the node
---@field bootID string Boot ID reported by the node.
---@field containerRuntimeVersion string ContainerRuntime Version reported by the node through runtime remote API (e.g. containerd://1.4.2).
---@field kernelVersion string Kernel Version reported by the node from 'uname -r' (e.g. 3.16.0-0.bpo.4-amd64).
---@field kubeProxyVersion string Deprecated: KubeProxy Version reported by the node.
---@field kubeletVersion string Kubelet Version reported by the node.
---@field machineID string MachineID reported by the node. For unique machine identification in the cluster this field is preferred. Learn more from man(5) machine-id: http://man7.org/linux/man-pages/man5/machine-id.5.html
---@field operatingSystem string The Operating System reported by the node
---@field osImage string OS Image reported by the node from /etc/os-release (e.g. Debian GNU/Linux 7 (wheezy)).
---@field runningInUserNamespace boolean Whether the node is running in a user namespace. +featureGate=KubeletInUserNamespace +optional
---@field swap corev1.NodeSwapStatus Swap Info reported by the node.
---@field systemUUID string SystemUUID reported by the node. For unique machine identification MachineID is preferred. This field is specific to Red Hat hosts https://access.redhat.com/documentation/en-us/red_hat_subscription_management/1/html/rhsm/uuid

---@class corev1.ObjectFieldSelector
---@field apiVersion string Version of the schema the FieldPath is written in terms of, defaults to "v1". +optional
---@field fieldPath string Path of the field to select in the specified API version.

---@class corev1.ObjectReference
---@field apiVersion string API version of the referent. +optional
---@field fieldPath string If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: "spec.containers{name}" (where "name" refers to the name of the container that triggered the event) or if no container name is specified "spec.containers[2]" (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object. TODO: this design is not final and this field is subject to change in the future. +optional
---@field kind string Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional
---@field name string Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names +optional
---@field namespace string Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/ +optional
---@field resourceVersion string Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency +optional
---@field uid string UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids +optional

---@class corev1.PersistentVolume
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.PersistentVolumeSpec spec defines a specification of a persistent volume owned by the cluster. Provisioned by an administrator. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistent-volumes +optional
---@field status corev1.PersistentVolumeStatus status represents the current information/status for the persistent volume. Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistent-volumes +optional

---@class corev1.PersistentVolumeClaim
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.PersistentVolumeClaimSpec spec defines the desired characteristics of a volume requested by a pod author. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims +optional
---@field status corev1.PersistentVolumeClaimStatus status represents the current information/status of a persistent volume claim. Read-only. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims +optional

---@class corev1.PersistentVolumeClaimCondition
---@field lastProbeTime v1.Time lastProbeTime is the time we probed the condition. +optional
---@field lastTransitionTime v1.Time lastTransitionTime is the time the condition transitioned from one status to another. +optional
---@field message string message is the human-readable message indicating details about last transition. +optional
---@field reason string reason is a unique, this should be a short, machine understandable string that gives the reason for condition's last transition. If it reports "Resizing" that means the underlying persistent volume is being resized. +optional
---@field status string Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
---@field type string Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about

---@class corev1.PersistentVolumeClaimList
---@field items corev1.PersistentVolumeClaim[] items is a list of persistent volume claims. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.PersistentVolumeClaimSpec
---@field accessModes string[] accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1 +optional +listType=atomic
---@field dataSource corev1.TypedLocalObjectReference dataSource field can be used to specify either: * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot) * An existing PVC (PersistentVolumeClaim) If the provisioner or an external controller can support the specified data source, it will create a new volume based on the contents of the specified data source. dataSource contents will be copied to dataSourceRef, and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified. If the namespace is specified, then dataSourceRef will not be copied to dataSource. +optional
---@field dataSourceRef corev1.TypedObjectReference dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired. This may be any object from a non-empty API group (non core object) or a PersistentVolumeClaim object. When this field is specified, volume binding will only succeed if the type of the specified object matches some installed volume populator or dynamic provisioner. This field will replace the functionality of the dataSource field and as such if both fields are non-empty, they must have the same value. For backwards compatibility, when namespace isn't specified in dataSourceRef, both fields (dataSource and dataSourceRef) will be set to the same value automatically if one of them is empty and the other is non-empty. When namespace is specified in dataSourceRef, dataSource isn't set to the same value and must be empty. There are three important differences between dataSource and dataSourceRef: * While dataSource only allows two specific types of objects, dataSourceRef allows any non-core object, as well as PersistentVolumeClaim objects. * While dataSource ignores disallowed values (dropping them), dataSourceRef preserves all values, and generates an error if a disallowed value is specified. * While dataSource only allows local objects, dataSourceRef allows objects in any namespaces. (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled. +optional
---@field resources corev1.VolumeResourceRequirements resources represents the minimum resources the volume should have. Users are allowed to specify resource requirements that are lower than previous value but must still be higher than capacity recorded in the status field of the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources +optional
---@field selector v1.LabelSelector selector is a label query over volumes to consider for binding. +optional
---@field storageClassName string storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1 +optional
---@field volumeAttributesClassName string volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string or nil value indicates that no VolumeAttributesClass will be applied to the claim. If the claim enters an Infeasible error state, this field can be reset to its previous value (including nil) to cancel the modification. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ +featureGate=VolumeAttributesClass +optional
---@field volumeMode string volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec. +optional
---@field volumeName string volumeName is the binding reference to the PersistentVolume backing this claim. +optional

---@class corev1.PersistentVolumeClaimStatus
---@field accessModes string[] accessModes contains the actual access modes the volume backing the PVC has. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1 +optional +listType=atomic
---@field allocatedResourceStatuses table<string, string> allocatedResourceStatuses stores status of resource being resized for the given PVC. Key names follow standard Kubernetes label syntax. Valid values are either: * Un-prefixed keys: - storage - the capacity of the volume. * Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource" Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used. ClaimResourceStatus can be in any of following states: - ControllerResizeInProgress: State set when resize controller starts resizing the volume in control-plane. - ControllerResizeFailed: State set when resize has failed in resize controller with a terminal error. - NodeResizePending: State set when resize controller has finished resizing the volume but further resizing of volume is needed on the node. - NodeResizeInProgress: State set when kubelet starts resizing the volume. - NodeResizeFailed: State set when resizing has failed in kubelet with a terminal error. Transient errors don't set NodeResizeFailed. For example: if expanding a PVC for more capacity - this field can be one of the following states: - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeInProgress" - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeFailed" - pvc.status.allocatedResourceStatus['storage'] = "NodeResizePending" - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeInProgress" - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeFailed" When this field is not set, it means that no resize operation is in progress for the given PVC. A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC. +mapType=granular +optional
---@field allocatedResources table<string, resource.Quantity> allocatedResources tracks the resources allocated to a PVC including its capacity. Key names follow standard Kubernetes label syntax. Valid values are either: * Un-prefixed keys: - storage - the capacity of the volume. * Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource" Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used. Capacity reported here may be larger than the actual capacity when a volume expansion operation is requested. For storage quota, the larger value from allocatedResources and PVC.spec.resources is used. If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation. If a volume expansion capacity request is lowered, allocatedResources is only lowered if there are no expansion operations in progress and if the actual volume capacity is equal or lower than the requested capacity. A controller that receives PVC update with previously unknown resourceName should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC. +optional
---@field capacity table<string, resource.Quantity> capacity represents the actual resources of the underlying volume. +optional
---@field conditions corev1.PersistentVolumeClaimCondition[] conditions is the current Condition of persistent volume claim. If underlying persistent volume is being resized then the Condition will be set to 'Resizing'. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field currentVolumeAttributesClassName string currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using. When unset, there is no VolumeAttributeClass applied to this PersistentVolumeClaim +featureGate=VolumeAttributesClass +optional
---@field healthStatus corev1.VolumeHealthStatus healthStatus contains the latest controller-reported health information for the volume bound to this claim. +featureGate=CSIVolumeHealth +optional +k8s:optional
---@field modifyVolumeStatus corev1.ModifyVolumeStatus ModifyVolumeStatus represents the status object of ControllerModifyVolume operation. When this is unset, there is no ModifyVolume operation being attempted. +featureGate=VolumeAttributesClass +optional
---@field phase string phase represents the current phase of PersistentVolumeClaim. +optional

---@class corev1.PersistentVolumeList
---@field items corev1.PersistentVolume[] items is a list of persistent volumes. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.PersistentVolumeSpec
---@field accessModes string[] accessModes contains all ways the volume can be mounted. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes +optional +listType=atomic
---@field capacity table<string, resource.Quantity> capacity is the description of the persistent volume's resources and capacity. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity +optional
---@field claimRef corev1.ObjectReference claimRef is part of a bi-directional binding between PersistentVolume and PersistentVolumeClaim. Expected to be non-nil when bound. claim.VolumeName is the authoritative bind between PV and PVC. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#binding +optional +structType=granular
---@field mountOptions string[] mountOptions is the list of mount options, e.g. ["ro", "soft"]. Not validated - mount will simply fail if one is invalid. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options +optional +listType=atomic
---@field nodeAffinity corev1.VolumeNodeAffinity nodeAffinity defines constraints that limit what nodes this volume can be accessed from. This field influences the scheduling of pods that use this volume. This field is mutable if MutablePVNodeAffinity feature gate is enabled. +optional
---@field persistentVolumeReclaimPolicy string persistentVolumeReclaimPolicy defines what happens to a persistent volume when released from its claim. Valid options are Retain (default for manually created PersistentVolumes), Delete (default for dynamically provisioned PersistentVolumes), and Recycle (deprecated). Recycle must be supported by the volume plugin underlying this PersistentVolume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#reclaiming +optional
---@field storageClassName string storageClassName is the name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass. +optional
---@field volumeAttributesClassName string Name of VolumeAttributesClass to which this persistent volume belongs. Empty value is not allowed. When this field is not set, it indicates that this volume does not belong to any VolumeAttributesClass. This field is mutable and can be changed by the CSI driver after a volume has been updated successfully to a new class. For an unbound PersistentVolume, the volumeAttributesClassName will be matched with unbound PersistentVolumeClaims during the binding process. +featureGate=VolumeAttributesClass +optional
---@field volumeMode string volumeMode defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state. Value of Filesystem is implied when not included in spec. +optional

---@class corev1.PersistentVolumeStatus
---@field lastPhaseTransitionTime v1.Time lastPhaseTransitionTime is the time the phase transitioned from one to another and automatically resets to current time everytime a volume phase transitions. +optional
---@field message string message is a human-readable message indicating details about why the volume is in this state. +optional
---@field phase string phase indicates if a volume is available, bound to a claim, or released by a claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#phase +optional
---@field reason string reason is a brief CamelCase string that describes any failure and is meant for machine parsing and tidy display in the CLI. +optional

---@class corev1.Pod
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.PodSpec Specification of the desired behavior of the pod. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status corev1.PodStatus Most recently observed status of the pod. This data may not be up to date. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class corev1.PodAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.WeightedPodAffinityTerm[] The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred. +optional +listType=atomic
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.PodAffinityTerm[] If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied. +optional +listType=atomic

---@class corev1.PodAffinityTerm
---@field labelSelector v1.LabelSelector A label query over a set of resources, in this case pods. If it's null, this PodAffinityTerm matches with no Pods. +optional
---@field matchLabelKeys string[] MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set. +listType=atomic +optional
---@field mismatchLabelKeys string[] MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set. +listType=atomic +optional
---@field namespaceSelector v1.LabelSelector A label query over the set of namespaces that the term applies to. The term is applied to the union of the namespaces selected by this field and the ones listed in the namespaces field. null selector and null or empty namespaces list means "this pod's namespace". An empty selector ({}) matches all namespaces. +optional
---@field namespaces string[] namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace". +optional +listType=atomic
---@field topologyKey string This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.

---@class corev1.PodAntiAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.WeightedPodAffinityTerm[] The scheduler will prefer to schedule pods to nodes that satisfy the anti-affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling anti-affinity expressions, etc.), compute a sum by iterating through the elements of this field and subtracting "weight" from the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred. +optional +listType=atomic
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.PodAffinityTerm[] If the anti-affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the anti-affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied. +optional +listType=atomic

---@class corev1.PodCondition
---@field lastProbeTime v1.Time Last time we probed the condition. +optional
---@field lastTransitionTime v1.Time Last time the condition transitioned from one status to another. +optional
---@field message string Human-readable message indicating details about last transition. +optional
---@field observedGeneration number If set, this represents the .metadata.generation that the pod condition was set based upon. +optional
---@field reason string Unique, one-word, CamelCase reason for the condition's last transition. +optional
---@field status string Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions
---@field type string Type is the type of the condition. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions

---@class corev1.PodDNSConfig
---@field nameservers string[] A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed. +optional +listType=atomic
---@field options corev1.PodDNSConfigOption[] A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy. +optional +listType=atomic
---@field searches string[] A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed. +optional +listType=atomic

---@class corev1.PodDNSConfigOption
---@field name string Name is this DNS resolver option's name. Required.
---@field value string Value is this DNS resolver option's value. +optional

---@class corev1.PodExtendedResourceClaimStatus
---@field requestMappings corev1.ContainerExtendedResourceRequest[] RequestMappings identifies the mapping of <container, extended resource backed by DRA> to device request in the generated ResourceClaim. +listType=atomic
---@field resourceClaimName string ResourceClaimName is the name of the ResourceClaim that was generated for the Pod in the namespace of the Pod.

---@class corev1.PodIP
---@field ip string IP is the IP address assigned to the pod +required

---@class corev1.PodList
---@field items corev1.Pod[] List of pods. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.PodOS
---@field name string Name is the name of the operating system. The currently supported values are linux and windows. Additional value may be defined in future and can be one of: https://github.com/opencontainers/runtime-spec/blob/master/config.md#platform-specific-configuration Clients should expect to handle additional values and treat unrecognized values in this field as os: null

---@class corev1.PodReadinessGate
---@field conditionType string ConditionType refers to a condition in the pod's condition list with matching type.

---@class corev1.PodResourceClaim
---@field name string Name uniquely identifies this resource claim inside the pod. This must be a DNS_LABEL.
---@field resourceClaimName string ResourceClaimName is the name of a ResourceClaim object in the same namespace as this pod. Exactly one of ResourceClaimName and ResourceClaimTemplateName must be set.
---@field resourceClaimTemplateName string ResourceClaimTemplateName is the name of a ResourceClaimTemplate object in the same namespace as this pod. The template will be used to create a new ResourceClaim, which will be bound to this pod. When this pod is deleted, the ResourceClaim will also be deleted. The pod name and resource name, along with a generated component, will be used to form a unique name for the ResourceClaim, which will be recorded in pod.status.resourceClaimStatuses. When the DRAWorkloadResourceClaims feature gate is enabled and the pod belongs to a PodGroup that defines a PodGroupResourceClaim with the same Name and ResourceClaimTemplateName, this PodResourceClaim resolves to the ResourceClaim generated for the PodGroup. All pods in the group that define an equivalent PodResourceClaim matching the PodGroupResourceClaim's Name and ResourceClaimTemplateName share the same generated ResourceClaim. ResourceClaims generated for a PodGroup are owned by the PodGroup and their lifecycles are tied to the PodGroup instead of any individual pod. This field is immutable and no changes will be made to the corresponding ResourceClaim by the control plane after creating the ResourceClaim. Exactly one of ResourceClaimName and ResourceClaimTemplateName must be set.

---@class corev1.PodResourceClaimStatus
---@field name string Name uniquely identifies this resource claim inside the pod. This must match the name of an entry in pod.spec.resourceClaims, which implies that the string must be a DNS_LABEL.
---@field resourceClaimName string ResourceClaimName is the name of the ResourceClaim that was generated for the Pod in the namespace of the Pod. When the DRAWorkloadResourceClaims feature is enabled and the corresponding PodResourceClaim matches a PodGroupResourceClaim made by the Pod's PodGroup, then this is the name of the ResourceClaim generated and reserved for the PodGroup. If this is unset, then generating a ResourceClaim was not necessary. The pod.spec.resourceClaims entry can be ignored in this case. +optional

---@class corev1.PodSchedulingGate
---@field name string Name of the scheduling gate. Each scheduling gate must have a unique name field.

---@class corev1.PodSchedulingGroup
---@field podGroupName string PodGroupName specifies the name of the standalone PodGroup object that represents the runtime instance of this group. Must be a DNS subdomain. +optional +oneOf=GroupSelection

---@class corev1.PodSecurityContext
---@field appArmorProfile corev1.AppArmorProfile appArmorProfile is the AppArmor options to use by the containers in this pod. Note that this field cannot be set when spec.os.name is windows. +optional
---@field fsGroup number A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod: 1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw---- If unset, the Kubelet will not modify the ownership and permissions of any volume. Note that this field cannot be set when spec.os.name is windows. +optional
---@field fsGroupChangePolicy string fsGroupChangePolicy defines behavior of changing ownership and permission of the volume before being exposed inside Pod. This field will only apply to volume types which support fsGroup based ownership(and permissions). It will have no effect on ephemeral volume types such as: secret, configmaps and emptydir. Valid values are "OnRootMismatch" and "Always". If not specified, "Always" is used. Note that this field cannot be set when spec.os.name is windows. +optional
---@field runAsGroup number The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows. +optional
---@field runAsNonRoot boolean Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. +optional
---@field runAsUser number The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows. +optional
---@field seLinuxChangePolicy string seLinuxChangePolicy defines how the container's SELinux label is applied to all volumes used by the Pod. It has no effect on nodes that do not support SELinux or to volumes does not support SELinux. Valid values are "MountOption" and "Recursive". "Recursive" means relabeling of all files on all Pod volumes by the container runtime. This may be slow for large volumes, but allows mixing privileged and unprivileged Pods sharing the same volume on the same node. "MountOption" mounts all eligible Pod volumes with `-o context` mount option. This requires all Pods that share the same volume to use the same SELinux label. It is not possible to share the same volume among privileged and unprivileged Pods. Eligible volumes are in-tree FibreChannel and iSCSI volumes, and all CSI volumes whose CSI driver announces SELinux support by setting spec.seLinuxMount: true in their CSIDriver instance. Other volumes are always re-labelled recursively. If not specified, "MountOption" is used. This field affects only Pods that have SELinux label set, either in PodSecurityContext or in SecurityContext of all containers. All Pods that use the same volume should use the same seLinuxChangePolicy, otherwise some pods can get stuck in ContainerCreating state. Note that this field cannot be set when spec.os.name is windows. +featureGate=SELinuxChangePolicy +optional
---@field seLinuxOptions corev1.SELinuxOptions The SELinux context to be applied to all containers. If unspecified, the container runtime will allocate a random SELinux context for each container. May also be set in SecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows. +optional
---@field seccompProfile corev1.SeccompProfile The seccomp options to use by the containers in this pod. Note that this field cannot be set when spec.os.name is windows. +optional
---@field supplementalGroups number[] A list of groups applied to the first process run in each container, in addition to the container's primary GID and fsGroup (if specified). If the SupplementalGroupsPolicy feature is enabled, the supplementalGroupsPolicy field determines whether these are in addition to or instead of any group memberships defined in the container image. If unspecified, no additional groups are added, though group memberships defined in the container image may still be used, depending on the supplementalGroupsPolicy field. Note that this field cannot be set when spec.os.name is windows. +optional +listType=atomic
---@field supplementalGroupsPolicy string Defines how supplemental groups of the first container processes are calculated. Valid values are "Merge" and "Strict". If not specified, "Merge" is used. (Alpha) Using the field requires the SupplementalGroupsPolicy feature gate to be enabled and the container runtime must implement support for this feature. Note that this field cannot be set when spec.os.name is windows. TODO: update the default value to "Merge" when spec.os.name is not windows in v1.34 +featureGate=SupplementalGroupsPolicy +optional
---@field sysctls corev1.Sysctl[] Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported sysctls (by the container runtime) might fail to launch. Note that this field cannot be set when spec.os.name is windows. +optional +listType=atomic
---@field windowsOptions corev1.WindowsSecurityContextOptions The Windows specific settings applied to all containers. If unspecified, the options within a container's SecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux. +optional

---@class corev1.PodSpec
---@field activeDeadlineSeconds number Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer. +optional
---@field affinity corev1.Affinity If specified, the pod's scheduling constraints +optional
---@field automountServiceAccountToken boolean AutomountServiceAccountToken indicates whether a service account token should be automatically mounted. +optional
---@field containers corev1.Container[] List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated. +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field dnsConfig corev1.PodDNSConfig Specifies the DNS parameters of a pod. Parameters specified here will be merged to the generated DNS configuration based on DNSPolicy. +optional
---@field dnsPolicy string Set DNS policy for the pod. Defaults to "ClusterFirst". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy explicitly to 'ClusterFirstWithHostNet'. +optional
---@field enableServiceLinks boolean EnableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links. Optional: Defaults to true. +optional
---@field ephemeralContainers corev1.EphemeralContainer[] List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing pod to perform user-initiated actions such as debugging. This list cannot be specified when creating a pod, and it cannot be modified by updating the pod spec. In order to add an ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field evictionResponders corev1.EvictionResponder[] evictionResponders reference responders that react to Evictions based on EvictionRequests. Responders should observe and communicate through the Eviction Resource API to help with the graceful termination of a pod. The responders are selected sequentially, according to their specified priority. Responders should periodically report on an eviction progress by updating the .status.responders[].heartbeatTime field of the Eviction object. If this field is not updated within the heartbeat deadline defined by the Eviction API (currently 20 minutes), the eviction is passed over to the next responder with a lower priority. If there is no other responder, the last default imperative-eviction.k8s.io/evictor responder with a priority of 100 will evict the pod using the imperative Eviction API (pods/<name>/eviction subresource). The maximum length of the responders list is 10. Responders are not supported when the pod is part of a PodGroup (.spec.schedulingGroup is set). This field can only be set on creation and is immutable afterwards. +featureGate=EvictionRequestAPI +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name +k8s:optional +k8s:listType=map +k8s:listMapKey=name +k8s:maxItems=10 +k8s:alpha(since: "1.37")=+k8s:dependentForbidden("schedulingGroup")
---@field hostAliases corev1.HostAlias[] HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified. +optional +patchMergeKey=ip +patchStrategy=merge +listType=map +listMapKey=ip
---@field hostIPC boolean Use the host's ipc namespace. Optional: Default to false. +optional
---@field hostNetwork boolean Host networking requested for this pod. Use the host's network namespace. When using HostNetwork you should specify ports so the scheduler is aware. When `hostNetwork` is true, specified `hostPort` fields in port definitions must match `containerPort`, and unspecified `hostPort` fields in port definitions are defaulted to match `containerPort`. Default to false. +optional
---@field hostPID boolean Use the host's pid namespace. Optional: Default to false. +optional
---@field hostUsers boolean Use the host's user namespace. Optional: Default to true. If set to true or not present, the pod will be run in the host user namespace, useful for when the pod needs a feature only available to the host user namespace, such as loading a kernel module with CAP_SYS_MODULE. When set to false, a new userns is created for the pod. Setting false is useful for mitigating container breakout vulnerabilities even allowing users to run their containers as root without actually having root privileges on the host. +optional
---@field hostname string Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value. +optional
---@field hostnameOverride string HostnameOverride specifies an explicit override for the pod's hostname as perceived by the pod. This field only specifies the pod's hostname and does not affect its DNS records. When this field is set to a non-empty string: - It takes precedence over the values set in `hostname` and `subdomain`. - The Pod's hostname will be set to this value. - `setHostnameAsFQDN` must be nil or set to false. - `hostNetwork` must be set to false. This field must be a valid DNS subdomain as defined in RFC 1123 and contain at most 64 characters. +featureGate=HostnameOverride +optional
---@field imagePullSecrets corev1.LocalObjectReference[] ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field initContainers corev1.Container[] List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers. Init containers may not have Lifecycle actions, Readiness probes, Liveness probes, or Startup probes. The resourceRequirements of an init container are taken into account during scheduling by finding the highest request/limit for each resource type, and then using the max of that value or the sum of the normal containers. Limits are applied to init containers in a similar fashion. Init containers cannot currently be added or removed. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/ +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name
---@field nodeName string NodeName indicates in which node this pod is scheduled. If empty, this pod is a candidate for scheduling by the scheduler defined in schedulerName. Once this field is set, the kubelet for this node becomes responsible for the lifecycle of this pod. This field should not be used to express a desire for the pod to be scheduled on a specific node. https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodename +optional
---@field nodeSelector table<string, string> NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ +optional +mapType=atomic
---@field os corev1.PodOS Specifies the OS of the containers in the pod. Some pod and container fields are restricted if this is set. If the OS field is set to linux, the following fields must be unset: -securityContext.windowsOptions If the OS field is set to windows, following fields must be unset: - spec.hostPID - spec.hostIPC - spec.hostUsers - spec.resources - spec.securityContext.appArmorProfile - spec.securityContext.seLinuxOptions - spec.securityContext.seccompProfile - spec.securityContext.fsGroup - spec.securityContext.fsGroupChangePolicy - spec.securityContext.sysctls - spec.shareProcessNamespace - spec.securityContext.runAsUser - spec.securityContext.runAsGroup - spec.securityContext.supplementalGroups - spec.securityContext.supplementalGroupsPolicy - spec.containers[*].securityContext.appArmorProfile - spec.containers[*].securityContext.seLinuxOptions - spec.containers[*].securityContext.seccompProfile - spec.containers[*].securityContext.capabilities - spec.containers[*].securityContext.readOnlyRootFilesystem - spec.containers[*].securityContext.privileged - spec.containers[*].securityContext.allowPrivilegeEscalation - spec.containers[*].securityContext.procMount - spec.containers[*].securityContext.runAsUser - spec.containers[*].securityContext.runAsGroup +optional
---@field overhead table<string, resource.Quantity> Overhead represents the resource overhead associated with running a pod for a given RuntimeClass. This field will be autopopulated at admission time by the RuntimeClass admission controller. If the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests. The RuntimeClass admission controller will reject Pod create requests which have the overhead already set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero. More info: https://git.k8s.io/enhancements/keps/sig-node/688-pod-overhead/README.md +optional
---@field preemptionPolicy string PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. Defaults to PreemptLowerPriority if unset. +optional
---@field priority number The priority value. Various system components use this field to find the priority of the pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority. +optional
---@field priorityClassName string If specified, indicates the pod's priority. "system-node-critical" and "system-cluster-critical" are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default. +optional
---@field readinessGates corev1.PodReadinessGate[] If specified, all readiness gates will be evaluated for pod readiness. A pod is ready when all its containers are ready AND all conditions specified in the readiness gates have status equal to "True" More info: https://git.k8s.io/enhancements/keps/sig-network/580-pod-readiness-gates +optional +listType=atomic
---@field resourceClaims corev1.PodResourceClaim[] ResourceClaims defines which ResourceClaims must be allocated and reserved before the Pod is allowed to start. The resources will be made available to those containers which consume them by name. This is a stable field but requires that the DynamicResourceAllocation feature gate is enabled. This field is immutable. +patchMergeKey=name +patchStrategy=merge,retainKeys +listType=map +listMapKey=name +featureGate=DynamicResourceAllocation +optional
---@field resources corev1.ResourceRequirements Resources is the total amount of CPU and Memory resources required by all containers in the pod. It supports specifying Requests and Limits for "cpu", "memory" and "hugepages-" resource names only. ResourceClaims are not supported. This field enables fine-grained control over resource allocation for the entire pod, allowing resource sharing among containers in a pod. TODO: For beta graduation, expand this comment with a detailed explanation. This is an alpha field and requires enabling the PodLevelResources feature gate. +featureGate=PodLevelResources +optional
---@field restartPolicy string Restart policy for all containers within the pod. One of Always, OnFailure, Never. In some contexts, only a subset of those values may be permitted. Default to Always. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy +optional
---@field runtimeClassName string RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used to run this pod. If no RuntimeClass resource matches the named class, the pod will not be run. If unset or empty, the "legacy" RuntimeClass will be used, which is an implicit class with an empty definition that uses the default runtime handler. More info: https://git.k8s.io/enhancements/keps/sig-node/585-runtime-class +optional
---@field schedulerName string If specified, the pod will be dispatched by specified scheduler. If not specified, the pod will be dispatched by default scheduler. +optional
---@field schedulingGates corev1.PodSchedulingGate[] SchedulingGates is an opaque list of values that if specified will block scheduling the pod. If schedulingGates is not empty, the pod will stay in the SchedulingGated state and the scheduler will not attempt to schedule the pod. SchedulingGates can only be set at pod creation time, and be removed only afterwards. +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name +optional
---@field schedulingGroup corev1.PodSchedulingGroup SchedulingGroup provides a reference to the immediate scheduling runtime grouping object that this Pod belongs to. This field is used by the scheduler to identify the group and apply the correct group scheduling policies. The association with a group also impacts other lifecycle aspects of a Pod that are relevant in a wider context of scheduling like preemption, resource attachment, etc. If not specified, the Pod is treated as a single unit in all of these aspects. The group object referenced by this field may not exist at the time the Pod is created. This field is immutable, but a group object with the same name may be recreated with different policies. Doing this during pod scheduling may result in the placement not conforming to the expected policies. +featureGate=GenericWorkload +optional
---@field securityContext corev1.PodSecurityContext SecurityContext holds pod-level security attributes and common container settings. Optional: Defaults to empty. See type description for default values of each field. +optional
---@field serviceAccount string DeprecatedServiceAccount is a deprecated alias for ServiceAccountName. Deprecated: Use serviceAccountName instead. +optional
---@field serviceAccountName string ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/ +optional
---@field setHostnameAsFQDN boolean If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. If a pod does not have FQDN, this has no effect. Default to false. +optional
---@field shareProcessNamespace boolean Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set. Optional: Default to false. +optional
---@field subdomain string If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>". If not specified, the pod will not have a domainname at all. +optional
---@field terminationGracePeriodSeconds number Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. Defaults to 30 seconds. +optional
---@field tolerations corev1.Toleration[] If specified, the pod's tolerations. +optional +listType=atomic +k8s:alpha(since: "1.37")=+k8s:optional
---@field topologySpreadConstraints corev1.TopologySpreadConstraint[] TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. All topologySpreadConstraints are ANDed. +optional +patchMergeKey=topologyKey +patchStrategy=merge +listType=map +listMapKey=topologyKey +listMapKey=whenUnsatisfiable
---@field volumes corev1.Volume[] List of volumes that can be mounted by containers belonging to the pod. More info: https://kubernetes.io/docs/concepts/storage/volumes +optional +patchMergeKey=name +patchStrategy=merge,retainKeys +listType=map +listMapKey=name

---@class corev1.PodStatus
---@field allocatedResources table<string, resource.Quantity> AllocatedResources is the total requests allocated for this pod by the node. If pod-level requests are not set, this will be the total requests aggregated across containers in the pod. +featureGate=InPlacePodLevelResourcesVerticalScaling +optional
---@field conditions corev1.PodCondition[] Current service state of pod. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type
---@field containerStatuses corev1.ContainerStatus[] Statuses of containers in this pod. Each container in the pod should have at most one status in this list, and all statuses should be for containers in the pod. However this is not enforced. If a status for a non-existent container is present in the list, or the list has duplicate names, the behavior of various Kubernetes components is not defined and those statuses might be ignored. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-and-container-status +optional +listType=atomic
---@field ephemeralContainerStatuses corev1.ContainerStatus[] Statuses for any ephemeral containers that have run in this pod. Each ephemeral container in the pod should have at most one status in this list, and all statuses should be for containers in the pod. However this is not enforced. If a status for a non-existent container is present in the list, or the list has duplicate names, the behavior of various Kubernetes components is not defined and those statuses might be ignored. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-and-container-status +optional +listType=atomic
---@field extendedResourceClaimStatus corev1.PodExtendedResourceClaimStatus Status of extended resource claim backed by DRA. +featureGate=DRAExtendedResource +optional
---@field hostIP string hostIP holds the IP address of the host to which the pod is assigned. Empty if the pod has not started yet. A pod can be assigned to a node that has a problem in kubelet which in turns mean that HostIP will not be updated even if there is a node is assigned to pod +optional
---@field hostIPs corev1.HostIP[] hostIPs holds the IP addresses allocated to the host. If this field is specified, the first entry must match the hostIP field. This list is empty if the pod has not started yet. A pod can be assigned to a node that has a problem in kubelet which in turns means that HostIPs will not be updated even if there is a node is assigned to this pod. +optional +patchStrategy=merge +patchMergeKey=ip +listType=atomic
---@field initContainerStatuses corev1.ContainerStatus[] Statuses of init containers in this pod. The most recent successful non-restartable init container will have ready = true, the most recently started container will have startTime set. Each init container in the pod should have at most one status in this list, and all statuses should be for containers in the pod. However this is not enforced. If a status for a non-existent container is present in the list, or the list has duplicate names, the behavior of various Kubernetes components is not defined and those statuses might be ignored. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-and-container-status +listType=atomic
---@field message string A human readable message indicating details about why the pod is in this condition. +optional
---@field nodeAllocatableResourceClaimStatuses corev1.NodeAllocatableResourceClaimStatus[] NodeAllocatableResourceClaimStatuses contains the status of node-allocatable resources that were allocated for this pod through DRA claims. This includes resources currently reported in v1.Node `status.allocatable` that are not extended resources (see https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#extended-resources). Examples include "cpu", "memory", "ephemeral-storage", and hugepages. +featureGate=DRANodeAllocatableResources +optional +patchStrategy=merge +patchMergeKey=resourceClaimName +listType=map +listMapKey=resourceClaimName +k8s:optional +k8s:listType=map +k8s:listMapKey=resourceClaimName
---@field nominatedNodeName string nominatedNodeName is set only when this pod preempts other pods on the node, but it cannot be scheduled right away as preemption victims receive their graceful termination periods. This field does not guarantee that the pod will be scheduled on this node. Scheduler may decide to place the pod elsewhere if other nodes become available sooner. Scheduler may also decide to give the resources on this node to a higher priority pod that is created after preemption. As a result, this field may be different than PodSpec.nodeName when the pod is scheduled. +optional
---@field observedGeneration number If set, this represents the .metadata.generation that the pod status was set based upon. The PodObservedGenerationTracking feature gate must be enabled to use this field. +optional
---@field phase string The phase of a Pod is a simple, high-level summary of where the Pod is in its lifecycle. The conditions array, the reason and message fields, and the individual container status arrays contain more detail about the pod's status. There are five possible phase values: Pending: The pod has been accepted by the Kubernetes system, but one or more of the container images has not been created. This includes time before being scheduled as well as time spent downloading images over the network, which could take a while. Running: The pod has been bound to a node, and all of the containers have been created. At least one container is still running, or is in the process of starting or restarting. Succeeded: All containers in the pod have terminated in success, and will not be restarted. Failed: All containers in the pod have terminated, and at least one container has terminated in failure. The container either exited with non-zero status or was terminated by the system. Unknown: For some reason the state of the pod could not be obtained, typically due to an error in communicating with the host of the pod. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-phase +optional
---@field podIP string podIP address allocated to the pod. Routable at least within the cluster. Empty if not yet allocated. +optional
---@field podIPs corev1.PodIP[] podIPs holds the IP addresses allocated to the pod. If this field is specified, the 0th entry must match the podIP field. Pods may be allocated at most 1 value for each of IPv4 and IPv6. This list is empty if no IPs have been allocated yet. +optional +patchStrategy=merge +patchMergeKey=ip +listType=map +listMapKey=ip
---@field qosClass string The Quality of Service (QOS) classification assigned to the pod based on resource requirements See PodQOSClass type for available QOS classes More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/#quality-of-service-classes +optional
---@field reason string A brief CamelCase message indicating details about why the pod is in this state. e.g. 'Evicted' +optional
---@field resize string Status of resources resize desired for pod's containers. It is empty if no resources resize is pending. Any changes to container resources will automatically set this to "Proposed" Deprecated: Resize status is moved to two pod conditions PodResizePending and PodResizeInProgress. PodResizePending will track states where the spec has been resized, but the Kubelet has not yet allocated the resources. PodResizeInProgress will track in-progress resizes, and should be present whenever allocated resources != acknowledged resources. +featureGate=InPlacePodVerticalScaling +optional
---@field resourceClaimStatuses corev1.PodResourceClaimStatus[] Status of resource claims. +patchMergeKey=name +patchStrategy=merge,retainKeys +listType=map +listMapKey=name +featureGate=DynamicResourceAllocation +optional
---@field resources corev1.ResourceRequirements Resources represents the compute resource requests and limits that have been applied at the pod level if pod-level requests or limits are set in PodSpec.Resources +featureGate=InPlacePodLevelResourcesVerticalScaling +optional
---@field startTime v1.Time RFC 3339 date and time at which the object was acknowledged by the Kubelet. This is before the Kubelet pulled the container image(s) for the pod. +optional
---@field volumeHealth corev1.PodVolumeHealth[] volumeHealth contains node-reported health for each volume the pod is using. Populated by the kubelet on the pod's node. +featureGate=CSIVolumeHealth +optional +listType=map +listMapKey=name +k8s:optional +k8s:listType=map +k8s:listMapKey=name

---@class corev1.PodTemplateSpec
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional +k8s:opaqueType
---@field spec corev1.PodSpec Specification of the desired behavior of the pod. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class corev1.PodVolumeHealth
---@field healthConditions corev1.VolumeHealthCondition[] conditions is the set of adverse conditions reported by the CSI node plugin for this volume on this node. At most 16 conditions may be reported. +optional +listType=map +listMapKey=status +patchMergeKey=status +patchStrategy=merge +listMapKey=reason +k8s:optional +k8s:listType=map +k8s:listMapKey=status +k8s:listMapKey=reason +k8s:maxItems=16
---@field lastTransitionTime v1.Time lastTransitionTime is when the current set of conditions first appeared. +optional
---@field name string name matches an entry in pod.spec.volumes. +required +k8s:required

---@class corev1.PortStatus
---@field error string Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use CamelCase names - cloud provider specific error values must have names that comply with the format foo.example.com/CamelCase. --- The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) +optional +kubebuilder:validation:Required +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$` +kubebuilder:validation:MaxLength=316
---@field port number Port is the port number of the service port of which status is recorded here
---@field protocol string Protocol is the protocol of the service port of which status is recorded here The supported values are: "TCP", "UDP", "SCTP"

---@class corev1.PreferredSchedulingTerm
---@field preference corev1.NodeSelectorTerm A node selector term, associated with the corresponding weight.
---@field weight number Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.

---@class corev1.Probe
---@field failureThreshold number Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1. +optional
---@field initialDelaySeconds number Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes +optional
---@field periodSeconds number How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1. +optional
---@field successThreshold number Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1. +optional
---@field terminationGracePeriodSeconds number Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset. +optional
---@field timeoutSeconds number Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes +optional

---@class corev1.ResourceClaim
---@field name string Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
---@field request string Request is the name chosen for a request in the referenced claim. If empty, everything from the claim is made available, otherwise only the result of this request. +optional

---@class corev1.ResourceFieldSelector
---@field containerName string Container name: required for volumes, optional for env vars +optional
---@field divisor resource.Quantity Specifies the output format of the exposed resources, defaults to "1" +optional
---@field resource string Required: resource to select

---@class corev1.ResourceHealth
---@field health string Health of the resource. can be one of: - Healthy: operates as normal - Unhealthy: reported unhealthy. We consider this a temporary health issue since we do not have a mechanism today to distinguish temporary and permanent issues. - Unknown: The status cannot be determined. For example, Device Plugin got unregistered and hasn't been re-registered since. In future we may want to introduce the PermanentlyUnhealthy Status.
---@field message string Message provides human-readable context for Health (e.g. "ECC error count exceeded threshold"). This field is populated by the kubelet when ResourceHealthStatusMessage is enabled if the DRA plugin returns a message, and is null otherwise. +featureGate=ResourceHealthStatusMessage +optional
---@field resourceID string ResourceID is the unique identifier of the resource. See the ResourceID type for more information.

---@class corev1.ResourceRequirements
---@field claims corev1.ResourceClaim[] Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container. This field depends on the DynamicResourceAllocation feature gate. This field is immutable. It can only be set for containers. +listType=map +listMapKey=name +featureGate=DynamicResourceAllocation +optional
---@field limits table<string, resource.Quantity> Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ +optional
---@field requests table<string, resource.Quantity> Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ +optional

---@class corev1.ResourceStatus
---@field name string Name of the resource. Must be unique within the pod and in case of non-DRA resource, match one of the resources from the pod spec. For DRA resources, the value must be "claim:<claim_name>/<request>" when container.resources.claims[*].request is set or "claim:<claim_name>" when container.resources.claims[*].request is empty. For DRA-backed extended resources, "claim:<claim_name>/<request>" is used when the claim name and request name are recorded in pod.status.extendedResourceClaimStatus. When this status is reported about a container, the "claim_name" and "request" must match one of the claims of this container. +required
---@field resources corev1.ResourceHealth[] List of unique resources health. Each element in the list contains an unique resource ID and its health. At a minimum, for the lifetime of a Pod, resource ID must uniquely identify the resource allocated to the Pod on the Node. If other Pod on the same Node reports the status with the same resource ID, it must be the same resource they share. See ResourceID type definition for a specific format it has in various use cases. +listType=map +listMapKey=resourceID

---@class corev1.SELinuxOptions
---@field level string Level is SELinux level label that applies to the container. +optional
---@field role string Role is a SELinux role label that applies to the container. +optional
---@field type string Type is a SELinux type label that applies to the container. +optional
---@field user string User is a SELinux user label that applies to the container. +optional

---@class corev1.SeccompProfile
---@field localhostProfile string localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type. +optional
---@field type string type indicates which kind of seccomp profile will be applied. Valid options are: Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied. +unionDiscriminator

---@class corev1.Secret
---@field data table<string, number[]> Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4 +optional
---@field immutable boolean Immutable, if set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil. +optional
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field stringData table<string, string> stringData allows specifying non-binary secret data in string form. It is provided as a write-only input field for convenience. All keys and values are merged into the data field on write, overwriting any existing values. The stringData field is never output when reading from the API. +k8s:conversion-gen=false +optional
---@field type string Used to facilitate programmatic handling of secret data. More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types +optional +k8s:optional +k8s:alpha(since: "1.37")=+k8s:immutable

---@class corev1.SecretEnvSource
---@field optional boolean Specify whether the Secret must be defined +optional

---@class corev1.SecretKeySelector
---@field key string The key of the secret to select from. Must be a valid secret key.
---@field optional boolean Specify whether the Secret or its key must be defined +optional

---@class corev1.SecretList
---@field items corev1.Secret[] Items is a list of secret objects. More info: https://kubernetes.io/docs/concepts/configuration/secret
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.SecurityContext
---@field allowPrivilegeEscalation boolean AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows. +optional
---@field appArmorProfile corev1.AppArmorProfile appArmorProfile is the AppArmor options to use by this container. If set, this profile overrides the pod's appArmorProfile. Note that this field cannot be set when spec.os.name is windows. +optional
---@field capabilities corev1.Capabilities The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime. Note that this field cannot be set when spec.os.name is windows. +optional
---@field privileged boolean Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows. +optional
---@field procMount string procMount denotes the type of proc mount to use for the containers. The default value is Default which uses the container runtime defaults for readonly paths and masked paths. Note that this field cannot be set when spec.os.name is windows. +optional
---@field readOnlyRootFilesystem boolean Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows. +optional
---@field runAsGroup number The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows. +optional
---@field runAsNonRoot boolean Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. +optional
---@field runAsUser number The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows. +optional
---@field seLinuxOptions corev1.SELinuxOptions The SELinux context to be applied to the container. If unspecified, the container runtime will allocate a random SELinux context for each container. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows. +optional
---@field seccompProfile corev1.SeccompProfile The seccomp options to use by this container. If seccomp options are provided at both the pod & container level, the container options override the pod options. Note that this field cannot be set when spec.os.name is windows. +optional
---@field windowsOptions corev1.WindowsSecurityContextOptions The Windows specific settings applied to all containers. If unspecified, the options from the PodSecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux. +optional

---@class corev1.Service
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec corev1.ServiceSpec Spec defines the behavior of a service. https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status corev1.ServiceStatus Most recently observed status of the service. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class corev1.ServiceAccount
---@field automountServiceAccountToken boolean AutomountServiceAccountToken indicates whether pods running as this service account should have an API token automatically mounted. Can be overridden at the pod level. +optional
---@field imagePullSecrets corev1.LocalObjectReference[] ImagePullSecrets is a list of references to secrets in the same namespace to use for pulling any images in pods that reference this ServiceAccount. ImagePullSecrets are distinct from Secrets because Secrets can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet. More info: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod +optional +listType=atomic
---@field metadata v1.ObjectMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field secrets corev1.ObjectReference[] Secrets is a list of the secrets in the same namespace that pods running using this ServiceAccount are allowed to use. Pods are only limited to this list if this service account has a "kubernetes.io/enforce-mountable-secrets" annotation set to "true". The "kubernetes.io/enforce-mountable-secrets" annotation is deprecated since v1.32. Prefer separate namespaces to isolate access to mounted secrets. This field should not be used to find auto-generated service account token secrets for use outside of pods. Instead, tokens can be requested directly using the TokenRequest API, or service account token secrets can be manually created. More info: https://kubernetes.io/docs/concepts/configuration/secret +optional +patchMergeKey=name +patchStrategy=merge +listType=map +listMapKey=name

---@class corev1.ServiceAccountList
---@field items corev1.ServiceAccount[] List of ServiceAccounts. More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.ServiceList
---@field items corev1.Service[] List of services
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class corev1.ServicePort
---@field appProtocol string The application protocol for this port. This is used as a hint for implementations to offer richer behavior for protocols that they understand. This field follows standard Kubernetes label syntax. Valid values are either: * Un-prefixed protocol names - reserved for IANA standard service names (as per RFC-6335 and https://www.iana.org/assignments/service-names). * Kubernetes-defined prefixed names: * 'kubernetes.io/h2c' - HTTP/2 prior knowledge over cleartext as described in https://www.rfc-editor.org/rfc/rfc9113.html#name-starting-http-2-with-prior- * 'kubernetes.io/ws' - WebSocket over cleartext as described in https://www.rfc-editor.org/rfc/rfc6455 * 'kubernetes.io/wss' - WebSocket over TLS as described in https://www.rfc-editor.org/rfc/rfc6455 * Other protocols should use implementation-defined prefixed names such as mycompany.com/my-custom-protocol. +optional
---@field name string The name of this port within the service. This must be a DNS_LABEL. All ports within a ServiceSpec must have unique names. When considering the endpoints for a Service, this must match the 'name' field in the EndpointPort. Optional if only one ServicePort is defined on this service. +optional
---@field nodePort number The port on each node on which this service is exposed when type is NodePort or LoadBalancer. Usually assigned by the system. If a value is specified, in-range, and not in use it will be used, otherwise the operation will fail. If not specified, a port will be allocated if this Service requires one. If this field is specified when creating a Service which does not need it, creation will fail. This field will be wiped when updating a Service to no longer need it (e.g. changing type from NodePort to ClusterIP). More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport +optional
---@field port number The port that will be exposed by this service.
---@field protocol string The IP protocol for this port. Supports "TCP", "UDP", and "SCTP". Default is TCP. +default="TCP" +optional
---@field targetPort intstr.IntOrString Number or name of the port to access on the pods targeted by the service. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME. If this is a string, it will be looked up as a named port in the target Pod's container ports. If this is not specified, the value of the 'port' field is used (an identity map). This field is ignored for services with clusterIP=None, and should be omitted or set equal to the 'port' field. More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service +optional

---@class corev1.ServiceSpec
---@field allocateLoadBalancerNodePorts boolean allocateLoadBalancerNodePorts defines if NodePorts will be automatically allocated for services with type LoadBalancer. Default is "true". It may be set to "false" if the cluster load-balancer does not rely on NodePorts. If the caller requests specific NodePorts (by specifying a value), those requests will be respected, regardless of this field. This field may only be set for services with type LoadBalancer and will be cleared if the type is changed to any other type. +optional
---@field clusterIP string clusterIP is the IP address of the service and is usually assigned randomly. If an address is specified manually, is in-range (as per system configuration), and is not in use, it will be allocated to the service; otherwise creation of the service will fail. This field may not be changed through updates unless the type field is also being changed to ExternalName (which requires this field to be blank) or the type field is being changed from ExternalName (in which case this field may optionally be specified, as describe above). Valid values are "None", empty string (""), or a valid IP address. Setting this to "None" makes a "headless service" (no virtual IP), which is useful when direct endpoint connections are preferred and proxying is not required. Only applies to types ClusterIP, NodePort, and LoadBalancer. If this field is specified when creating a Service of type ExternalName, creation will fail. This field will be wiped when updating a Service to type ExternalName. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies +optional
---@field clusterIPs string[] ClusterIPs is a list of IP addresses assigned to this service, and are usually assigned randomly. If an address is specified manually, is in-range (as per system configuration), and is not in use, it will be allocated to the service; otherwise creation of the service will fail. This field may not be changed through updates unless the type field is also being changed to ExternalName (which requires this field to be empty) or the type field is being changed from ExternalName (in which case this field may optionally be specified, as describe above). Valid values are "None", empty string (""), or a valid IP address. Setting this to "None" makes a "headless service" (no virtual IP), which is useful when direct endpoint connections are preferred and proxying is not required. Only applies to types ClusterIP, NodePort, and LoadBalancer. If this field is specified when creating a Service of type ExternalName, creation will fail. This field will be wiped when updating a Service to type ExternalName. If this field is not specified, it will be initialized from the clusterIP field. If this field is specified, clients must ensure that clusterIPs[0] and clusterIP have the same value. This field may hold a maximum of two entries (dual-stack IPs, in either order). These IPs must correspond to the values of the ipFamilies field. Both clusterIPs and ipFamilies are governed by the ipFamilyPolicy field. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies +listType=atomic +optional
---@field externalIPs string[] externalIPs is a list of IP addresses for which nodes in the cluster will also accept traffic for this service. These IPs are not managed by Kubernetes. The user is responsible for ensuring that traffic arrives at a node with this IP. A common example is external load-balancers that are not part of the Kubernetes system. +optional +listType=atomic
---@field externalName string externalName is the external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record). No proxying will be involved. Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) and requires `type` to be "ExternalName". +optional
---@field externalTrafficPolicy string externalTrafficPolicy describes how nodes distribute service traffic they receive on one of the Service's "externally-facing" addresses (NodePorts, ExternalIPs, and LoadBalancer IPs). If set to "Local", the proxy will configure the service in a way that assumes that external load balancers will take care of balancing the service traffic between nodes, and so each node will deliver traffic only to the node-local endpoints of the service, without masquerading the client source IP. (Traffic mistakenly sent to a node with no endpoints will be dropped.) The default value, "Cluster", uses the standard behavior of routing to all endpoints evenly (possibly modified by topology and other features). Note that traffic sent to an External IP or LoadBalancer IP from within the cluster will always get "Cluster" semantics, but clients sending to a NodePort from within the cluster may need to take traffic policy into account when picking a node. +optional
---@field healthCheckNodePort number healthCheckNodePort specifies the healthcheck nodePort for the service. This only applies when type is set to LoadBalancer and externalTrafficPolicy is set to Local. If a value is specified, is in-range, and is not in use, it will be used. If not specified, a value will be automatically allocated. External systems (e.g. load-balancers) can use this port to determine if a given node holds endpoints for this service or not. If this field is specified when creating a Service which does not need it, creation will fail. This field will be wiped when updating a Service to no longer need it (e.g. changing type). This field cannot be updated once set. +optional
---@field internalTrafficPolicy string InternalTrafficPolicy describes how nodes distribute service traffic they receive on the ClusterIP. If set to "Local", the proxy will assume that pods only want to talk to endpoints of the service on the same node as the pod, dropping the traffic if there are no local endpoints. The default value, "Cluster", uses the standard behavior of routing to all endpoints evenly (possibly modified by topology and other features). +optional
---@field ipFamilies string[] IPFamilies is a list of IP families (e.g. IPv4, IPv6) assigned to this service. This field is usually assigned automatically based on cluster configuration and the ipFamilyPolicy field. If this field is specified manually, the requested family is available in the cluster, and ipFamilyPolicy allows it, it will be used; otherwise creation of the service will fail. This field is conditionally mutable: it allows for adding or removing a secondary IP family, but it does not allow changing the primary IP family of the Service. Valid values are "IPv4" and "IPv6". This field only applies to Services of types ClusterIP, NodePort, and LoadBalancer, and does apply to "headless" services. This field will be wiped when updating a Service to type ExternalName. This field may hold a maximum of two entries (dual-stack families, in either order). These families must correspond to the values of the clusterIPs field, if specified. Both clusterIPs and ipFamilies are governed by the ipFamilyPolicy field. +listType=atomic +optional
---@field ipFamilyPolicy string IPFamilyPolicy represents the dual-stack-ness requested or required by this Service. If there is no value provided, then this field will be set to SingleStack. Services can be "SingleStack" (a single IP family), "PreferDualStack" (two IP families on dual-stack configured clusters or a single IP family on single-stack clusters), or "RequireDualStack" (two IP families on dual-stack configured clusters, otherwise fail). The ipFamilies and clusterIPs fields depend on the value of this field. This field will be wiped when updating a service to type ExternalName. +optional
---@field loadBalancerClass string loadBalancerClass is the class of the load balancer implementation this Service belongs to. If specified, the value of this field must be a label-style identifier, with an optional prefix, e.g. "internal-vip" or "example.com/internal-vip". Unprefixed names are reserved for end-users. This field can only be set when the Service type is 'LoadBalancer'. If not set, the default load balancer implementation is used, today this is typically done through the cloud provider integration, but should apply for any default implementation. If set, it is assumed that a load balancer implementation is watching for Services with a matching class. Any default load balancer implementation (e.g. cloud providers) should ignore Services that set this field. This field can only be set when creating or updating a Service to type 'LoadBalancer'. Once set, it can not be changed. This field will be wiped when a service is updated to a non 'LoadBalancer' type. +optional
---@field loadBalancerIP string Only applies to Service Type: LoadBalancer. This feature depends on whether the underlying cloud-provider supports specifying the loadBalancerIP when a load balancer is created. This field will be ignored if the cloud-provider does not support the feature. Deprecated: This field was under-specified and its meaning varies across implementations. Using it is non-portable and it may not support dual-stack. Users are encouraged to use implementation-specific annotations when available. +optional
---@field loadBalancerSourceRanges string[] If specified and supported by the platform, this will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs. This field will be ignored if the cloud-provider does not support the feature." More info: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/ +optional +listType=atomic
---@field ports corev1.ServicePort[] The list of ports that are exposed by this service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies +patchMergeKey=port +patchStrategy=merge +listType=map +listMapKey=port +listMapKey=protocol
---@field publishNotReadyAddresses boolean publishNotReadyAddresses indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready. The primary use case for setting this field is for a StatefulSet's Headless Service to propagate SRV DNS records for its Pods for the purpose of peer discovery. The Kubernetes controllers that generate Endpoints and EndpointSlice resources for Services interpret this to mean that all endpoints are considered "ready" even if the Pods themselves are not. Agents which consume only Kubernetes generated endpoints through the Endpoints or EndpointSlice resources can safely assume this behavior. +optional
---@field selector table<string, string> Route service traffic to pods with label keys and values matching this selector. If empty or not present, the service is assumed to have an external process managing its endpoints, which Kubernetes will not modify. Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if type is ExternalName. More info: https://kubernetes.io/docs/concepts/services-networking/service/ +optional +mapType=atomic
---@field sessionAffinity string Supports "ClientIP" and "None". Used to maintain session affinity. Enable client IP based session affinity. Must be ClientIP or None. Defaults to None. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies +optional
---@field sessionAffinityConfig corev1.SessionAffinityConfig sessionAffinityConfig contains the configurations of session affinity. +optional
---@field trafficDistribution string TrafficDistribution offers a way to express preferences for how traffic is distributed to Service endpoints. Implementations can use this field as a hint, but are not required to guarantee strict adherence. If the field is not set, the implementation will apply its default routing strategy. If set to "PreferClose", implementations should prioritize endpoints that are in the same zone. +optional
---@field type string type determines how the Service is exposed. Defaults to ClusterIP. Valid options are ExternalName, ClusterIP, NodePort, and LoadBalancer. "ClusterIP" allocates a cluster-internal IP address for load-balancing to endpoints. Endpoints are determined by the selector or if that is not specified, by manual construction of an Endpoints object or EndpointSlice objects. If clusterIP is "None", no virtual IP is allocated and the endpoints are published as a set of endpoints rather than a virtual IP. "NodePort" builds on ClusterIP and allocates a port on every node which routes to the same endpoints as the clusterIP. "LoadBalancer" builds on NodePort and creates an external load-balancer (if supported in the current cloud) which routes to the same endpoints as the clusterIP. "ExternalName" aliases this service to the specified externalName. Several other fields do not apply to ExternalName services. More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types +optional

---@class corev1.ServiceStatus
---@field conditions v1.Condition[] Current service state +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type +k8s:alpha(since: "1.37")=+k8s:eachVal=+k8s:opaqueType
---@field loadBalancer corev1.LoadBalancerStatus LoadBalancer contains the current status of the load-balancer, if one is present. +optional

---@class corev1.SessionAffinityConfig
---@field clientIP corev1.ClientIPConfig clientIP contains the configurations of Client IP based session affinity. +optional

---@class corev1.SleepAction
---@field seconds number Seconds is the number of seconds to sleep.

---@class corev1.Sysctl
---@field name string Name of a property to set
---@field value string Value of a property to set

---@class corev1.TCPSocketAction
---@field host string Optional: Host name to connect to, defaults to the pod IP. +optional
---@field port intstr.IntOrString Number or name of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.

---@class corev1.Taint
---@field effect string Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
---@field key string Required. The taint key to be applied to a node.
---@field timeAdded v1.Time TimeAdded represents the time at which the taint was added. +optional
---@field value string The taint value corresponding to the taint key. +optional

---@class corev1.Toleration
---@field effect string Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute. +optional
---@field key string Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys. +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:format=k8s-label-key
---@field operator string Operator represents a key's relationship to the value. Valid operators are Exists, Equal, Lt, and Gt. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category. Lt and Gt perform numeric comparisons (requires feature gate TaintTolerationComparisonOperators). +optional
---@field tolerationSeconds number TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system. +optional
---@field value string Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string. +optional

---@class corev1.TopologySelectorLabelRequirement
---@field key string The label key that the selector applies to.
---@field values string[] An array of string values. One value must match the label to be selected. Each entry in Values is ORed. +listType=atomic

---@class corev1.TopologySelectorTerm
---@field matchLabelExpressions corev1.TopologySelectorLabelRequirement[] A list of topology selector requirements by labels. +optional +listType=atomic

---@class corev1.TopologySpreadConstraint
---@field labelSelector v1.LabelSelector LabelSelector is used to find matching pods. Pods that match this label selector are counted to determine the number of pods in their corresponding topology domain. +optional
---@field matchLabelKeys string[] MatchLabelKeys is a set of pod label keys to select the pods over which spreading will be calculated. The keys are used to lookup values from the incoming pod labels, those key-value labels are ANDed with labelSelector to select the group of existing pods over which spreading will be calculated for the incoming pod. The same key is forbidden to exist in both MatchLabelKeys and LabelSelector. MatchLabelKeys cannot be set when LabelSelector isn't set. Keys that don't exist in the incoming pod labels will be ignored. A null or empty list means only match against labelSelector. This is a beta field and requires the MatchLabelKeysInPodTopologySpread feature gate to be enabled (enabled by default). +listType=atomic +optional
---@field maxSkew number MaxSkew describes the degree to which pods may be unevenly distributed. When `whenUnsatisfiable=DoNotSchedule`, it is the maximum permitted difference between the number of matching pods in the target topology and the global minimum. The global minimum is the minimum number of matching pods in an eligible domain or zero if the number of eligible domains is less than MinDomains. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 2/2/1: In this case, the global minimum is 1. +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P | P P | P | +-------+-------+-------+ - if MaxSkew is 1, incoming pod can only be scheduled to zone3 to become 2/2/2; scheduling it onto zone1(zone2) would make the ActualSkew(3-1) on zone1(zone2) violate MaxSkew(1). - if MaxSkew is 2, incoming pod can be scheduled onto any zone. When `whenUnsatisfiable=ScheduleAnyway`, it is used to give higher precedence to topologies that satisfy it. It's a required field. Default value is 1 and 0 is not allowed.
---@field minDomains number MinDomains indicates a minimum number of eligible domains. When the number of eligible domains with matching topology keys is less than minDomains, Pod Topology Spread treats "global minimum" as 0, and then the calculation of Skew is performed. And when the number of eligible domains with matching topology keys equals or greater than minDomains, this value has no effect on scheduling. As a result, when the number of eligible domains is less than minDomains, scheduler won't schedule more than maxSkew Pods to those domains. If value is nil, the constraint behaves as if MinDomains is equal to 1. Valid values are integers greater than 0. When value is not nil, WhenUnsatisfiable must be DoNotSchedule. For example, in a 3-zone cluster, MaxSkew is set to 2, MinDomains is set to 5 and pods with the same labelSelector spread as 2/2/2: +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P | P P | P P | +-------+-------+-------+ The number of domains is less than 5(MinDomains), so "global minimum" is treated as 0. In this situation, new pod with the same labelSelector cannot be scheduled, because computed skew will be 3(3 - 0) if new Pod is scheduled to any of the three zones, it will violate MaxSkew. +optional
---@field nodeAffinityPolicy string NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew. Options are: - Honor: only nodes matching nodeAffinity/nodeSelector are included in the calculations. - Ignore: nodeAffinity/nodeSelector are ignored. All nodes are included in the calculations. If this value is nil, the behavior is equivalent to the Honor policy. +optional
---@field nodeTaintsPolicy string NodeTaintsPolicy indicates how we will treat node taints when calculating pod topology spread skew. Options are: - Honor: nodes without taints, along with tainted nodes for which the incoming pod has a toleration, are included. - Ignore: node taints are ignored. All nodes are included. If this value is nil, the behavior is equivalent to the Ignore policy. +optional
---@field topologyKey string TopologyKey is the key of node labels. Nodes that have a label with this key and identical values are considered to be in the same topology. We consider each <key, value> as a "bucket", and try to put balanced number of pods into each bucket. We define a domain as a particular instance of a topology. Also, we define an eligible domain as a domain whose nodes meet the requirements of nodeAffinityPolicy and nodeTaintsPolicy. e.g. If TopologyKey is "kubernetes.io/hostname", each Node is a domain of that topology. And, if TopologyKey is "topology.kubernetes.io/zone", each zone is a domain of that topology. It's a required field.
---@field whenUnsatisfiable string WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint. - DoNotSchedule (default) tells the scheduler not to schedule it. - ScheduleAnyway tells the scheduler to schedule the pod in any location, but giving higher precedence to topologies that would help reduce the skew. A constraint is considered "Unsatisfiable" for an incoming pod if and only if every possible node assignment for that pod would violate "MaxSkew" on some topology. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 3/1/1: +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P P | P | P | +-------+-------+-------+ If WhenUnsatisfiable is set to DoNotSchedule, incoming pod can only be scheduled to zone2(zone3) to become 3/2/1(3/1/2) as ActualSkew(2-1) on zone2(zone3) satisfies MaxSkew(1). In other words, the cluster can still be imbalanced, but scheduler won't make it *more* imbalanced. It's a required field.

---@class corev1.TypedLocalObjectReference
---@field apiGroup string APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required. +optional
---@field kind string Kind is the type of resource being referenced
---@field name string Name is the name of resource being referenced

---@class corev1.TypedObjectReference
---@field apiGroup string APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required. +optional
---@field kind string Kind is the type of resource being referenced
---@field name string Name is the name of resource being referenced
---@field namespace string Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details. (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled. +featureGate=CrossNamespaceVolumeDataSource +optional

---@class corev1.Volume
---@field name string name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names

---@class corev1.VolumeDevice
---@field devicePath string devicePath is the path inside of the container that the device will be mapped to.
---@field name string name must match the name of a persistentVolumeClaim in the pod

---@class corev1.VolumeHealthCondition
---@field message string message is a human-readable description. Maximum permitted length of a message is 1024 bytes. +optional +k8s:optional +k8s:maxBytes=1024
---@field reason string reason is a brief CamelCase machine-parseable reason. Together with status it forms the unique identity of a condition entry. Maximum permitted length of a reason is 256 bytes. +required +k8s:required +k8s:maxBytes=256
---@field status string status is the machine-parseable health category. Possible values: - "Inaccessible": the volume cannot be accessed. - "DataLoss": data loss has been detected on the volume. - "Degraded": the volume is functioning with reduced capability. +required +k8s:required

---@class corev1.VolumeHealthStatus
---@field healthConditions corev1.VolumeHealthCondition[] conditions is the set of adverse conditions reported by the CSI controller plugin. An empty list means no adverse condition. At most 16 conditions may be reported. +optional +listType=map +listMapKey=status +patchMergeKey=status +patchStrategy=merge +listMapKey=reason +k8s:optional +k8s:listType=map +k8s:listMapKey=status +k8s:listMapKey=reason +k8s:maxItems=16
---@field lastTransitionTime v1.Time lastTransitionTime is when the current set of conditions first appeared. +optional

---@class corev1.VolumeMount
---@field bindMountOptions string[] bindMountOptions is the list of additional bind mount options to apply when mounting this volume into the container. Allowed values are noexec, nodev, and nosuid. These are Linux mount options and have no effect on Windows nodes. This field is not supported with image volumes. This is an alpha field and requires enabling the VolumeBindMountOptions feature gate. +featureGate=VolumeBindMountOptions +optional +listType=set
---@field mountPath string Path within the container at which the volume should be mounted.
---@field mountPropagation string mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used. This field is beta in 1.10. When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified (which defaults to None). +optional
---@field name string This must match the Name of a Volume.
---@field readOnly boolean Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false. +optional
---@field recursiveReadOnly string RecursiveReadOnly specifies whether read-only mounts should be handled recursively. If ReadOnly is false, this field has no meaning and must be unspecified. If ReadOnly is true, and this field is set to Disabled, the mount is not made recursively read-only. If this field is set to IfPossible, the mount is made recursively read-only, if it is supported by the container runtime. If this field is set to Enabled, the mount is made recursively read-only if it is supported by the container runtime, otherwise the pod will not be started and an error will be generated to indicate the reason. If this field is set to IfPossible or Enabled, MountPropagation must be set to None (or be unspecified, which defaults to None). If this field is not specified, it is treated as an equivalent of Disabled. +optional
---@field subPath string Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root). +optional
---@field subPathExpr string Expanded path within the volume from which the container's volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment. Defaults to "" (volume's root). SubPathExpr and SubPath are mutually exclusive. +optional

---@class corev1.VolumeMountStatus
---@field mountPath string MountPath corresponds to the original VolumeMount.
---@field name string Name corresponds to the name of the original VolumeMount.
---@field readOnly boolean ReadOnly corresponds to the original VolumeMount. +optional
---@field recursiveReadOnly string RecursiveReadOnly must be set to Disabled, Enabled, or unspecified (for non-readonly mounts). An IfPossible value in the original VolumeMount must be translated to Disabled or Enabled, depending on the mount result. +optional
---@field volumeStatus corev1.VolumeStatus volumeStatus represents volume-type-specific status about the mounted volume. +optional

---@class corev1.VolumeNodeAffinity
---@field required corev1.NodeSelector required specifies hard node constraints that must be met.

---@class corev1.VolumeResourceRequirements
---@field limits table<string, resource.Quantity> Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ +optional
---@field requests table<string, resource.Quantity> Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ +optional

---@class corev1.VolumeStatus
---@field image corev1.ImageVolumeStatus image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet's host machine. +featureGate=ImageVolumeWithDigest +optional

---@class corev1.WeightedPodAffinityTerm
---@field podAffinityTerm corev1.PodAffinityTerm Required. A pod affinity term, associated with the corresponding weight.
---@field weight number weight associated with matching the corresponding podAffinityTerm, in the range 1-100.

---@class corev1.WindowsSecurityContextOptions
---@field gmsaCredentialSpec string GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field. +optional
---@field gmsaCredentialSpecName string GMSACredentialSpecName is the name of the GMSA credential spec to use. +optional
---@field hostProcess boolean HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true. +optional
---@field runAsUserName string The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. +optional

---@class discoveryv1.Endpoint
---@field addresses string[] addresses of this endpoint. For EndpointSlices of addressType "IPv4" or "IPv6", the values are IP addresses in canonical form. The syntax and semantics of other addressType values are not defined. This must contain at least one address but no more than 100. EndpointSlices generated by the EndpointSlice controller will always have exactly 1 address. No semantics are defined for additional addresses beyond the first, and kube-proxy does not look at them. +listType=set +required +k8s:beta(since: "1.37")=+k8s:required +k8s:beta(since: "1.37")=+k8s:maxItems=100
---@field conditions discoveryv1.EndpointConditions conditions contains information about the current status of the endpoint.
---@field deprecatedTopology table<string, string> deprecatedTopology contains topology information part of the v1beta1 API. This field is deprecated, and will be removed when the v1beta1 API is removed (no sooner than kubernetes v1.24). While this field can hold values, it is not writable through the v1 API, and any attempts to write to it will be silently ignored. Topology information can be found in the zone and nodeName fields instead. +optional
---@field hints discoveryv1.EndpointHints hints contains information associated with how an endpoint should be consumed. +optional
---@field hostname string hostname of this endpoint. This field may be used by consumers of endpoints to distinguish endpoints from each other (e.g. in DNS names). Multiple endpoints which use the same hostname should be considered fungible (e.g. multiple A values in DNS). Must be lowercase and pass DNS Label (RFC 1123) validation. +optional
---@field nodeName string nodeName represents the name of the Node hosting this endpoint. This can be used to determine endpoints local to a Node. +optional
---@field targetRef corev1.ObjectReference targetRef is a reference to a Kubernetes object that represents this endpoint. +optional
---@field zone string zone is the name of the Zone this endpoint exists in. +optional

---@class discoveryv1.EndpointConditions
---@field ready boolean ready indicates that this endpoint is ready to receive traffic, according to whatever system is managing the endpoint. A nil value should be interpreted as "true". In general, an endpoint should be marked ready if it is serving and not terminating, though this can be overridden in some cases, such as when the associated Service has set the publishNotReadyAddresses flag. +optional
---@field serving boolean serving indicates that this endpoint is able to receive traffic, according to whatever system is managing the endpoint. For endpoints backed by pods, the EndpointSlice controller will mark the endpoint as serving if the pod's Ready condition is True. A nil value should be interpreted as "true". +optional
---@field terminating boolean terminating indicates that this endpoint is terminating. A nil value should be interpreted as "false". +optional

---@class discoveryv1.EndpointHints
---@field forNodes discoveryv1.ForNode[] forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. +listType=atomic
---@field forZones discoveryv1.ForZone[] forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. +listType=atomic

---@class discoveryv1.EndpointPort
---@field appProtocol string The application protocol for this port. This is used as a hint for implementations to offer richer behavior for protocols that they understand. This field follows standard Kubernetes label syntax. Valid values are either: * Un-prefixed protocol names - reserved for IANA standard service names (as per RFC-6335 and https://www.iana.org/assignments/service-names). * Kubernetes-defined prefixed names: * 'kubernetes.io/h2c' - HTTP/2 prior knowledge over cleartext as described in https://www.rfc-editor.org/rfc/rfc9113.html#name-starting-http-2-with-prior- * 'kubernetes.io/ws' - WebSocket over cleartext as described in https://www.rfc-editor.org/rfc/rfc6455 * 'kubernetes.io/wss' - WebSocket over TLS as described in https://www.rfc-editor.org/rfc/rfc6455 * Other protocols should use implementation-defined prefixed names such as mycompany.com/my-custom-protocol. +optional
---@field name string name represents the name of this port. All ports in an EndpointSlice must have a unique name. If the EndpointSlice is derived from a Kubernetes service, this corresponds to the Service.ports[].name. Name must either be an empty string or pass DNS_LABEL validation: * must be no more than 63 characters long. * must consist of lower case alphanumeric characters or '-'. * must start and end with an alphanumeric character. Default is empty string.
---@field port number port represents the port number of the endpoint. If the EndpointSlice is derived from a Kubernetes service, this must be set to the service's target port. EndpointSlices used for other purposes may have a nil port.
---@field protocol string protocol represents the IP protocol for this port. Must be UDP, TCP, or SCTP. Default is TCP.

---@class discoveryv1.EndpointSlice
---@field addressType string addressType specifies the type of address carried by this EndpointSlice. All addresses in this slice must be the same type. This field is immutable after creation. The following address types are currently supported: * IPv4: Represents an IPv4 Address. * IPv6: Represents an IPv6 Address. * FQDN: Represents a Fully Qualified Domain Name. (Deprecated) The EndpointSlice controller only generates, and kube-proxy only processes, slices of addressType "IPv4" and "IPv6". No semantics are defined for the "FQDN" type. +required +k8s:beta(since: "1.37")=+k8s:required +k8s:beta(since: "1.37")=+k8s:immutable
---@field endpoints discoveryv1.Endpoint[] endpoints is a list of unique endpoints in this slice. Each slice may include a maximum of 1000 endpoints. +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional
---@field metadata v1.ObjectMeta Standard object's metadata. +optional
---@field ports discoveryv1.EndpointPort[] ports specifies the list of network ports exposed by each endpoint in this slice. Each port must have a unique name. Each slice may include a maximum of 100 ports. Services always have at least 1 port, so EndpointSlices generated by the EndpointSlice controller will likewise always have at least 1 port. EndpointSlices used for other purposes may have an empty ports list. +optional +listType=atomic

---@class discoveryv1.EndpointSliceList
---@field items discoveryv1.EndpointSlice[] items is the list of endpoint slices
---@field metadata v1.ListMeta Standard list metadata. +optional

---@class discoveryv1.ForNode
---@field name string name represents the name of the node.

---@class discoveryv1.ForZone
---@field name string name represents the name of the zone.

---@class eventsv1.Event
---@field action string action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
---@field deprecatedCount number deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type. +optional
---@field deprecatedFirstTimestamp v1.Time deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type. +optional
---@field deprecatedLastTimestamp v1.Time deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type. +optional
---@field deprecatedSource corev1.EventSource deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type. +optional
---@field eventTime v1.MicroTime eventTime is the time when this Event was first observed. It is required.
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field note string note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB. +optional
---@field reason string reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
---@field regarding corev1.ObjectReference regarding contains the object this Event is about. In most cases it's an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object. +optional
---@field related corev1.ObjectReference related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object. +optional
---@field reportingController string reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
---@field reportingInstance string reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
---@field series eventsv1.EventSeries series is data about the Event series this event represents or nil if it's a singleton Event. +optional
---@field type string type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.

---@class eventsv1.EventList
---@field items eventsv1.Event[] items is a list of schema objects.
---@field metadata v1.ListMeta metadata is the standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class eventsv1.EventSeries
---@field count number count is the number of occurrences in this series up to the last heartbeat time.
---@field lastObservedTime v1.MicroTime lastObservedTime is the time when last Event from the series was seen before last heartbeat.

---@class networkingv1.IPBlock
---@field cidr string cidr is a string representing the IPBlock Valid examples are "192.168.1.0/24" or "2001:db8::/64" +required +k8s:beta(since: "1.37")=+k8s:required
---@field except string[] except is a slice of CIDRs that should not be included within an IPBlock Valid examples are "192.168.1.0/24" or "2001:db8::/64" Except values will be rejected if they are outside the cidr range +optional +listType=atomic

---@class networkingv1.Ingress
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec networkingv1.IngressSpec spec is the desired state of the Ingress. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional
---@field status networkingv1.IngressStatus status is the current state of the Ingress. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class networkingv1.IngressBackend
---@field resource corev1.TypedLocalObjectReference resource is an ObjectRef to another Kubernetes resource in the namespace of the Ingress object. If resource is specified, a service.Name and service.Port must not be specified. This is a mutually exclusive setting with "Service". +optional
---@field service networkingv1.IngressServiceBackend service references a service as a backend. This is a mutually exclusive setting with "Resource". +optional

---@class networkingv1.IngressList
---@field items networkingv1.Ingress[] items is the list of Ingress.
---@field metadata v1.ListMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class networkingv1.IngressLoadBalancerIngress
---@field hostname string hostname is set for load-balancer ingress points that are DNS based. +optional
---@field ip string ip is set for load-balancer ingress points that are IP based. +optional
---@field ports networkingv1.IngressPortStatus[] ports provides information about the ports exposed by this LoadBalancer. +listType=atomic +optional

---@class networkingv1.IngressLoadBalancerStatus
---@field ingress networkingv1.IngressLoadBalancerIngress[] ingress is a list containing ingress points for the load-balancer. +optional +listType=atomic

---@class networkingv1.IngressPortStatus
---@field error string error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use CamelCase names - cloud provider specific error values must have names that comply with the format foo.example.com/CamelCase. --- The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) +optional +kubebuilder:validation:Required +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$` +kubebuilder:validation:MaxLength=316
---@field port number port is the port number of the ingress port.
---@field protocol string protocol is the protocol of the ingress port. The supported values are: "TCP", "UDP", "SCTP"

---@class networkingv1.IngressRule
---@field host string host is the fully qualified domain name of a network host, as defined by RFC 3986. Note the following deviations from the "host" part of the URI as defined in RFC 3986: 1. IPs are not allowed. Currently an IngressRuleValue can only apply to the IP in the Spec of the parent Ingress. 2. The `:` delimiter is not respected because ports are not allowed. Currently the port of an Ingress is implicitly :80 for http and :443 for https. Both these may change in the future. Incoming requests are matched against the host before the IngressRuleValue. If the host is unspecified, the Ingress routes all traffic based on the specified IngressRuleValue. host can be "precise" which is a domain name without the terminating dot of a network host (e.g. "foo.bar.com") or "wildcard", which is a domain name prefixed with a single wildcard label (e.g. "*.foo.com"). The wildcard character '*' must appear by itself as the first DNS label and matches only a single label. You cannot have a wildcard label by itself (e.g. Host == "*"). Requests will be matched against the Host field in the following way: 1. If host is precise, the request matches this rule if the http host header is equal to Host. 2. If host is a wildcard, then the request matches this rule if the http host header is to equal to the suffix (removing the first label) of the wildcard rule. +optional

---@class networkingv1.IngressServiceBackend
---@field name string name is the referenced service. The service must exist in the same namespace as the Ingress object.
---@field port networkingv1.ServiceBackendPort port of the referenced service. A port name or port number is required for a IngressServiceBackend.

---@class networkingv1.IngressSpec
---@field defaultBackend networkingv1.IngressBackend defaultBackend is the backend that should handle requests that don't match any rule. If Rules are not specified, DefaultBackend must be specified. If DefaultBackend is not set, the handling of requests that do not match any of the rules will be up to the Ingress controller. +optional
---@field ingressClassName string ingressClassName is the name of an IngressClass cluster resource. Ingress controller implementations use this field to know whether they should be serving this Ingress resource, by a transitive connection (controller -> IngressClass -> Ingress resource). Although the `kubernetes.io/ingress.class` annotation (simple constant name) was never formally defined, it was widely supported by Ingress controllers to create a direct binding between Ingress controller and Ingress resources. Newly created Ingress resources should prefer using the field. However, even though the annotation is officially deprecated, for backwards compatibility reasons, ingress controllers should still honor that annotation if present. +optional
---@field rules networkingv1.IngressRule[] rules is a list of host rules used to configure the Ingress. If unspecified, or no rule matches, all traffic is sent to the default backend. +listType=atomic +optional
---@field tls networkingv1.IngressTLS[] tls represents the TLS configuration. Currently the Ingress only supports a single TLS port, 443. If multiple members of this list specify different hosts, they will be multiplexed on the same port according to the hostname specified through the SNI TLS extension, if the ingress controller fulfilling the ingress supports SNI. +listType=atomic +optional

---@class networkingv1.IngressStatus
---@field loadBalancer networkingv1.IngressLoadBalancerStatus loadBalancer contains the current status of the load-balancer. +optional

---@class networkingv1.IngressTLS
---@field hosts string[] hosts is a list of hosts included in the TLS certificate. The values in this list must match the name/s used in the tlsSecret. Defaults to the wildcard host setting for the loadbalancer controller fulfilling this Ingress, if left unspecified. +listType=atomic +optional
---@field secretName string secretName is the name of the secret used to terminate TLS traffic on port 443. Field is left optional to allow TLS routing based on SNI hostname alone. If the SNI host in a listener conflicts with the "Host" header field used by an IngressRule, the SNI host is used for termination and value of the "Host" header is used for routing. +optional

---@class networkingv1.NetworkPolicy
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec networkingv1.NetworkPolicySpec spec represents the specification of the desired behavior for this NetworkPolicy. +optional

---@class networkingv1.NetworkPolicyEgressRule
---@field ports networkingv1.NetworkPolicyPort[] ports is a list of destination ports for outgoing traffic. Each item in this list is combined using a logical OR. If this field is empty or missing, this rule matches all ports (traffic not restricted by port). If this field is present and contains at least one item, then this rule allows traffic only if the traffic matches at least one port in the list. +optional +listType=atomic
---@field to networkingv1.NetworkPolicyPeer[] to is a list of destinations for outgoing traffic of pods selected for this rule. Items in this list are combined using a logical OR operation. If this field is empty or missing, this rule matches all destinations (traffic not restricted by destination). If this field is present and contains at least one item, this rule allows traffic only if the traffic matches at least one item in the to list. +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional

---@class networkingv1.NetworkPolicyIngressRule
---@field from networkingv1.NetworkPolicyPeer[] from is a list of sources which should be able to access the pods selected for this rule. Items in this list are combined using a logical OR operation. If this field is empty or missing, this rule matches all sources (traffic not restricted by source). If this field is present and contains at least one item, this rule allows traffic only if the traffic matches at least one item in the from list. +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional
---@field ports networkingv1.NetworkPolicyPort[] ports is a list of ports which should be made accessible on the pods selected for this rule. Each item in this list is combined using a logical OR. If this field is empty or missing, this rule matches all ports (traffic not restricted by port). If this field is present and contains at least one item, then this rule allows traffic only if the traffic matches at least one port in the list. +optional +listType=atomic

---@class networkingv1.NetworkPolicyList
---@field items networkingv1.NetworkPolicy[] items is a list of schema objects.
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class networkingv1.NetworkPolicyPeer
---@field ipBlock networkingv1.IPBlock ipBlock defines policy on a particular IPBlock. If this field is set then neither of the other fields can be. +optional +k8s:beta(since: "1.37")=+k8s:optional
---@field namespaceSelector v1.LabelSelector namespaceSelector selects namespaces using cluster-scoped labels. This field follows standard label selector semantics; if present but empty, it selects all namespaces. If podSelector is also set, then the NetworkPolicyPeer as a whole selects the pods matching podSelector in the namespaces selected by namespaceSelector. Otherwise it selects all pods in the namespaces selected by namespaceSelector. +optional
---@field podSelector v1.LabelSelector podSelector is a label selector which selects pods. This field follows standard label selector semantics; if present but empty, it selects all pods. If namespaceSelector is also set, then the NetworkPolicyPeer as a whole selects the pods matching podSelector in the Namespaces selected by NamespaceSelector. Otherwise it selects the pods matching podSelector in the policy's own namespace. +optional

---@class networkingv1.NetworkPolicyPort
---@field endPort number endPort indicates that the range of ports from port to endPort if set, inclusive, should be allowed by the policy. This field cannot be defined if the port field is not defined or if the port field is defined as a named (string) port. The endPort must be equal or greater than port. +optional
---@field port intstr.IntOrString port represents the port on the given protocol. This can either be a numerical or named port on a pod. If this field is not provided, this matches all port names and numbers. If present, only traffic on the specified protocol AND port will be matched. +optional
---@field protocol string protocol represents the protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this field defaults to TCP. +optional

---@class networkingv1.NetworkPolicySpec
---@field egress networkingv1.NetworkPolicyEgressRule[] egress is a list of egress rules to be applied to the selected pods. Outgoing traffic is allowed if there are no NetworkPolicies selecting the pod (and cluster policy otherwise allows the traffic), OR if the traffic matches at least one egress rule across all of the NetworkPolicy objects whose podSelector matches the pod. If this field is empty then this NetworkPolicy limits all outgoing traffic (and serves solely to ensure that the pods it selects are isolated by default). This field is beta-level in 1.8 +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional
---@field ingress networkingv1.NetworkPolicyIngressRule[] ingress is a list of ingress rules to be applied to the selected pods. Traffic is allowed to a pod if there are no NetworkPolicies selecting the pod (and cluster policy otherwise allows the traffic), OR if the traffic source is the pod's local node, OR if the traffic matches at least one ingress rule across all of the NetworkPolicy objects whose podSelector matches the pod. If this field is empty then this NetworkPolicy does not allow any traffic (and serves solely to ensure that the pods it selects are isolated by default) +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional
---@field podSelector v1.LabelSelector podSelector selects the pods to which this NetworkPolicy object applies. The array of rules is applied to any pods selected by this field. An empty selector matches all pods in the policy's namespace. Multiple network policies can select the same set of pods. In this case, the ingress rules for each are combined additively. This field is optional. If it is not specified, it defaults to an empty selector. +optional
---@field policyTypes string[] policyTypes is a list of rule types that the NetworkPolicy relates to. Valid options are ["Ingress"], ["Egress"], or ["Ingress", "Egress"]. If this field is not specified, it will default based on the existence of ingress or egress rules; policies that contain an egress section are assumed to affect egress, and all policies (whether or not they contain an ingress section) are assumed to affect ingress. If you want to write an egress-only policy, you must explicitly specify policyTypes [ "Egress" ]. Likewise, if you want to write a policy that specifies that no egress is allowed, you must specify a policyTypes value that include "Egress" (since such a policy would not include an egress section and would otherwise default to just [ "Ingress" ]). This field is beta-level in 1.8 +optional +listType=atomic

---@class networkingv1.ServiceBackendPort
---@field name string name is the name of the port on the Service. This is a mutually exclusive setting with "Number". +optional
---@field number number number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with "Name". +optional

---@class policyv1.PodDisruptionBudget
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec policyv1.PodDisruptionBudgetSpec spec is the specification of the desired behavior of the PodDisruptionBudget. +optional
---@field status policyv1.PodDisruptionBudgetStatus status is the most recently observed status of the PodDisruptionBudget. +optional

---@class policyv1.PodDisruptionBudgetList
---@field items policyv1.PodDisruptionBudget[] Items is a list of PodDisruptionBudgets
---@field metadata v1.ListMeta Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class policyv1.PodDisruptionBudgetSpec
---@field maxUnavailable intstr.IntOrString maxUnavailable indicates that an eviction is allowed if at most "maxUnavailable" pods selected by "selector" are unavailable after the eviction, i.e. even in absence of the evicted pod. For example, one can prevent all voluntary evictions by specifying 0. This is a mutually exclusive setting with "minAvailable". +optional
---@field minAvailable intstr.IntOrString minAvailable indicates that an eviction is allowed if at least "minAvailable" pods selected by "selector" will still be available after the eviction, i.e. even in the absence of the evicted pod. So for example you can prevent all voluntary evictions by specifying "100%". +optional
---@field selector v1.LabelSelector selector is a label query over pods whose evictions are managed by the disruption budget. A null selector will match no pods, while an empty ({}) selector will select all pods within the namespace. +patchStrategy=replace +optional
---@field unhealthyPodEvictionPolicy string unhealthyPodEvictionPolicy defines the criteria for when unhealthy pods should be considered for eviction. Current implementation considers healthy pods, as pods that have status.conditions item with type="Ready",status="True". Valid policies are IfHealthyBudget and AlwaysAllow. If no policy is specified, the default behavior will be used, which corresponds to the IfHealthyBudget policy. IfHealthyBudget policy means that running pods (status.phase="Running"), but not yet healthy can be evicted only if the guarded application is not disrupted (status.currentHealthy is at least equal to status.desiredHealthy). Healthy pods will be subject to the PDB for eviction. AlwaysAllow policy means that all running pods (status.phase="Running"), but not yet healthy are considered disrupted and can be evicted regardless of whether the criteria in a PDB is met. This means perspective running pods of a disrupted application might not get a chance to become healthy. Healthy pods will be subject to the PDB for eviction. Additional policies may be added in the future. Clients making eviction decisions should disallow eviction of unhealthy pods if they encounter an unrecognized policy in this field. +optional

---@class policyv1.PodDisruptionBudgetStatus
---@field conditions v1.Condition[] Conditions contain conditions for PDB. The disruption controller sets the DisruptionAllowed condition. The following are known values for the reason field (additional reasons could be added in the future): - SyncFailed: The controller encountered an error and wasn't able to compute the number of allowed disruptions. Therefore no disruptions are allowed and the status of the condition will be False. - InsufficientPods: The number of pods are either at or below the number required by the PodDisruptionBudget. No disruptions are allowed and the status of the condition will be False. - SufficientPods: There are more pods than required by the PodDisruptionBudget. The condition will be True, and the number of allowed disruptions are provided by the disruptionsAllowed property. +optional +patchMergeKey=type +patchStrategy=merge +listType=map +listMapKey=type +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:listType=map +k8s:alpha(since: "1.37")=+k8s:listMapKey=type
---@field currentHealthy number current number of healthy pods +optional
---@field desiredHealthy number minimum desired number of healthy pods +optional
---@field disruptedPods table<string, v1.Time> DisruptedPods contains information about pods whose eviction was processed by the API server eviction subresource handler but has not yet been observed by the PodDisruptionBudget controller. A pod will be in this map from the time when the API server processed the eviction request to the time when the pod is seen by PDB controller as having been marked for deletion (or after a timeout). The key in the map is the name of the pod and the value is the time when the API server processed the eviction request. If the deletion didn't occur and a pod is still there it will be removed from the list automatically by PodDisruptionBudget controller after some time. If everything goes smooth this map should be empty for the most of the time. Large number of entries in the map may indicate problems with pod deletions. +optional
---@field disruptionsAllowed number Number of pod disruptions that are currently allowed. +optional
---@field expectedPods number total number of pods counted by this disruption budget +optional
---@field observedGeneration number Most recent generation observed when updating this PDB status. DisruptionsAllowed and other status information is valid only if observedGeneration equals to PDB's object generation. +optional

---@class rbacv1.AggregationRule
---@field clusterRoleSelectors v1.LabelSelector[] clusterRoleSelectors holds a list of selectors which will be used to find ClusterRoles and create the rules. If any of the selectors match, then the ClusterRole's permissions will be added +optional +listType=atomic

---@class rbacv1.ClusterRole
---@field aggregationRule rbacv1.AggregationRule aggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller. +optional
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. +optional
---@field rules rbacv1.PolicyRule[] rules holds all the PolicyRules for this ClusterRole +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional

---@class rbacv1.ClusterRoleBinding
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. +optional
---@field roleRef rbacv1.RoleRef roleRef can only reference a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error. This field is immutable. +required +k8s:alpha(since:"1.37")=+k8s:immutable
---@field subjects rbacv1.Subject[] subjects holds references to the objects the role applies to. +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional

---@class rbacv1.ClusterRoleBindingList
---@field items rbacv1.ClusterRoleBinding[] Items is a list of ClusterRoleBindings
---@field metadata v1.ListMeta Standard object's metadata. +optional

---@class rbacv1.ClusterRoleList
---@field items rbacv1.ClusterRole[] Items is a list of ClusterRoles
---@field metadata v1.ListMeta Standard object's metadata. +optional

---@class rbacv1.PolicyRule
---@field apiGroups string[] apiGroups is the name of the APIGroup that contains the resources. If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed. "" represents the core API group and "*" represents all API groups. +optional +listType=atomic
---@field nonResourceURLs string[] nonResourceURLs is a set of partial urls that a user should have access to. *s are allowed, but only as the full, final step in the path Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding. Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"), but not both. +optional +listType=atomic
---@field resourceNames string[] resourceNames is an optional white list of names that the rule applies to. An empty set means that everything is allowed. +optional +listType=atomic
---@field resources string[] resources is a list of resources this rule applies to. '*' represents all resources. +optional +listType=atomic
---@field verbs string[] verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule. '*' represents all verbs. +listType=atomic +required +k8s:beta(since: "1.37")=+k8s:required

---@class rbacv1.Role
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. +optional
---@field rules rbacv1.PolicyRule[] rules holds all the PolicyRules for this Role +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional

---@class rbacv1.RoleBinding
---@field metadata v1.ObjectMeta metadata is the standard object's metadata. +optional
---@field roleRef rbacv1.RoleRef roleRef can reference a Role in the current namespace or a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error. This field is immutable. +required +k8s:alpha(since:"1.37")=+k8s:immutable
---@field subjects rbacv1.Subject[] subjects holds references to the objects the role applies to. +optional +listType=atomic +k8s:beta(since: "1.37")=+k8s:optional

---@class rbacv1.RoleBindingList
---@field items rbacv1.RoleBinding[] Items is a list of RoleBindings
---@field metadata v1.ListMeta Standard object's metadata. +optional

---@class rbacv1.RoleList
---@field items rbacv1.Role[] Items is a list of Roles
---@field metadata v1.ListMeta Standard object's metadata. +optional

---@class rbacv1.RoleRef
---@field apiGroup string apiGroup is the group for the resource being referenced +optional
---@field kind string kind is the type of resource being referenced +required
---@field name string name is the name of resource being referenced +required +k8s:beta(since: "1.37")=+k8s:required

---@class rbacv1.Subject
---@field apiGroup string apiGroup holds the API group of the referenced subject. Defaults to "" for ServiceAccount subjects. Defaults to "rbac.authorization.k8s.io" for User and Group subjects. +optional
---@field kind string kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount". If the Authorizer does not recognized the kind value, the Authorizer should report an error. +required
---@field name string name of the object being referenced. +required +k8s:beta(since: "1.37")=+k8s:required
---@field namespace string namespace of the referenced object. If the object kind is non-namespace, such as "User" or "Group", and this value is not empty the Authorizer should report an error. +optional

---@class schedulingv1alpha3.TopologyConstraint
---@field key string key specifies the key of the node label representing the topology domain. All pods within the PodGroup must be colocated within the same domain instance. Different PodGroups can land on different domain instances even if they derive from the same PodGroupTemplate. Examples: "topology.kubernetes.io/rack" +required +k8s:required +k8s:format=k8s-label-key

---@class schedulingv1alpha3.WorkloadPodGroupDisruptionMode
---@field all schedulingv1alpha3.WorkloadPodGroupAllDisruptionMode all specifies that all pods in the group must be disrupted together. +optional +k8s:optional +k8s:unionMember
---@field single schedulingv1alpha3.WorkloadPodGroupSingleDisruptionMode single specifies that pods can be disrupted independently from each other. +optional +k8s:optional +k8s:unionMember

---@class schedulingv1alpha3.WorkloadPodGroupGangSchedulingPolicy
---@field minCount number minCount is the minimum number of pods that must be scheduled at the same time for the scheduler to admit the entire group. This field is optional. If it is not specified, the controller should inject a context-specific sane default (e.g., parallelism for a Job). If set, it must be a positive integer. +optional +k8s:optional +k8s:minimum=1

---@class schedulingv1alpha3.WorkloadPodGroupResourceClaim
---@field name string name uniquely identifies this resource claim inside the group. This field is required. It must be a DNS_LABEL. +required +k8s:required +k8s:format=k8s-short-name
---@field resourceClaimName string resourceClaimName is the name of a ResourceClaim object in the same namespace. This field is optional. If it is not specified, no resource claim is used. If set, it must be a DNS subdomain. +optional +k8s:optional +k8s:unionMember +k8s:format=k8s-long-name
---@field resourceClaimTemplateName string resourceClaimTemplateName is the name of a ResourceClaimTemplate object in the same namespace. This field is optional. If it is not specified, no resource claim template is used. If set, it must be a DNS subdomain. +optional +k8s:optional +k8s:unionMember +k8s:format=k8s-long-name

---@class schedulingv1alpha3.WorkloadPodGroupSchedulingConstraints
---@field topology schedulingv1alpha3.TopologyConstraint[] topology specifies desired topological placements for all pods within the pod group. If unset, no topology placement is requested. +optional +k8s:optional +k8s:maxItems=1 +listType=atomic +k8s:listType=atomic

---@class schedulingv1alpha3.WorkloadPodGroupSchedulingPolicy
---@field basic schedulingv1alpha3.WorkloadPodGroupBasicSchedulingPolicy basic specifies that standard, pod-by-pod Kubernetes scheduling behavior should be used. +optional +k8s:optional +k8s:unionMember
---@field gang schedulingv1alpha3.WorkloadPodGroupGangSchedulingPolicy gang specifies all-or-nothing scheduling semantics. +optional +k8s:optional +k8s:unionMember

---@class storagev1.StorageClass
---@field allowVolumeExpansion boolean allowVolumeExpansion shows whether the storage class allow volume expand. +optional
---@field allowedTopologies corev1.TopologySelectorTerm[] allowedTopologies restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction. This field is only honored by servers that enable the VolumeScheduling feature. +optional +listType=atomic
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field mountOptions string[] mountOptions controls the mountOptions for dynamically provisioned PersistentVolumes of this storage class. e.g. ["ro", "soft"]. Not validated - mount of the PVs will simply fail if one is invalid. +optional +listType=atomic
---@field parameters table<string, string> parameters holds the parameters for the provisioner that should create volumes of this storage class. +optional +k8s:beta(since: "1.37")=+k8s:immutable +k8s:beta(since: "1.37")=+k8s:optional
---@field provisioner string provisioner indicates the type of the provisioner. +required +k8s:beta(since: "1.37")=+k8s:required +k8s:beta(since: "1.37")=+k8s:immutable
---@field reclaimPolicy string reclaimPolicy controls the reclaimPolicy for dynamically provisioned PersistentVolumes of this storage class. Defaults to Delete. +optional +k8s:beta(since: "1.37")=+k8s:immutable +k8s:beta(since: "1.37")=+k8s:optional
---@field volumeBindingMode string volumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound. When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature. +optional +k8s:beta(since: "1.37")=+k8s:immutable +k8s:beta(since: "1.37")=+k8s:optional

---@class storagev1.StorageClassList
---@field items storagev1.StorageClass[] items is the list of StorageClasses
---@field metadata v1.ListMeta Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class storagev1.VolumeAttachment
---@field metadata v1.ObjectMeta metadata is the standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional
---@field spec storagev1.VolumeAttachmentSpec spec represents specification of the desired attach/detach volume behavior. Populated by the Kubernetes system. +k8s:beta(since: "1.37")=+k8s:immutable +required
---@field status storagev1.VolumeAttachmentStatus status represents status of the VolumeAttachment request. Populated by the entity completing the attach or detach operation, i.e. the external-attacher. +optional

---@class storagev1.VolumeAttachmentList
---@field items storagev1.VolumeAttachment[] items is the list of VolumeAttachments
---@field metadata v1.ListMeta Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional

---@class storagev1.VolumeAttachmentSource
---@field inlineVolumeSpec corev1.PersistentVolumeSpec inlineVolumeSpec contains all the information necessary to attach a persistent volume defined by a pod's inline VolumeSource. This field is populated only for the CSIMigration feature. It contains translated fields from a pod's inline VolumeSource to a PersistentVolumeSpec. This field is beta-level and is only honored by servers that enabled the CSIMigration feature. +optional
---@field persistentVolumeName string persistentVolumeName represents the name of the persistent volume to attach. +optional

---@class storagev1.VolumeAttachmentSpec
---@field attacher string attacher indicates the name of the volume driver that MUST handle this request. This is the name returned by GetPluginName(). +required +k8s:beta(since: "1.37")=+k8s:required +k8s:beta(since: "1.37")=+k8s:format="k8s-long-name-caseless" +k8s:beta(since: "1.37")=+k8s:maxLength=63
---@field nodeName string nodeName represents the node that the volume should be attached to.
---@field source storagev1.VolumeAttachmentSource source represents the volume that should be attached.

---@class storagev1.VolumeAttachmentStatus
---@field attachError storagev1.VolumeError attachError represents the last error encountered during attach operation, if any. This field must only be set by the entity completing the attach operation, i.e. the external-attacher. +optional
---@field attached boolean attached indicates the volume is successfully attached. This field must only be set by the entity completing the attach operation, i.e. the external-attacher.
---@field attachmentMetadata table<string, string> attachmentMetadata is populated with any information returned by the attach operation, upon successful attach, that must be passed into subsequent WaitForAttach or Mount calls. This field must only be set by the entity completing the attach operation, i.e. the external-attacher. +optional
---@field detachError storagev1.VolumeError detachError represents the last error encountered during detach operation, if any. This field must only be set by the entity completing the detach operation, i.e. the external-attacher. +optional

---@class storagev1.VolumeError
---@field errorCode number errorCode is a numeric gRPC code representing the error encountered during Attach or Detach operations. This field requires the MutableCSINodeAllocatableCount feature gate being enabled to be set. +featureGate=MutableCSINodeAllocatableCount +optional
---@field message string message represents the error encountered during Attach or Detach operation. This string may be logged, so it should not contain sensitive information. +optional
---@field time v1.Time time represents the time the error was encountered. +optional

---@class v1.Condition
---@field lastTransitionTime v1.Time lastTransitionTime is the last time the condition transitioned from one status to another. This should be when the underlying condition changed. If that is not known, then using the time when the API field changed is acceptable. +required +kubebuilder:validation:Required +kubebuilder:validation:Type=string +kubebuilder:validation:Format=date-time +k8s:alpha(since: "1.37")=+k8s:customValidation
---@field message string message is a human readable message indicating details about the transition. This may be an empty string. +required +kubebuilder:validation:Required +kubebuilder:validation:MaxLength=32768
---@field observedGeneration number observedGeneration represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date with respect to the current state of the instance. +optional +kubebuilder:validation:Minimum=0 +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:minimum=0
---@field reason string reason contains a programmatic identifier indicating the reason for the condition's last transition. Producers of specific condition types may define expected values and meanings for this field, and whether the values are considered a guaranteed API. The value should be a CamelCase string. This field may not be empty. +required +kubebuilder:validation:Required +kubebuilder:validation:MaxLength=1024 +kubebuilder:validation:MinLength=1 +kubebuilder:validation:Pattern=`^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$` +k8s:alpha(since: "1.37")=+k8s:required +k8s:alpha(since: "1.37")=+k8s:maxBytes=1024
---@field status string status of the condition, one of True, False, Unknown. +required +kubebuilder:validation:Required +kubebuilder:validation:Enum=True;False;Unknown +k8s:alpha(since: "1.37")=+k8s:required
---@field type string type of condition in CamelCase or in foo.example.com/CamelCase. --- Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be useful (see .node.status.conditions), the ability to deconflict is important. The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt) +required +kubebuilder:validation:Required +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$` +kubebuilder:validation:MaxLength=316 +k8s:alpha(since: "1.37")=+k8s:required

---@class v1.LabelSelector
---@field matchExpressions v1.LabelSelectorRequirement[] matchExpressions is a list of label selector requirements. The requirements are ANDed. +optional +listType=atomic
---@field matchLabels table<string, string> matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed. +optional

---@class v1.LabelSelectorRequirement
---@field key string key is the label key that the selector applies to.
---@field operator string operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
---@field values string[] values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch. +optional +listType=atomic

---@class v1.ListMeta
---@field continue string continue may be set if the user set a limit on the number of items returned, and indicates that the server has more data available. The value is opaque and may be used to issue another request to the endpoint that served this list to retrieve the next set of available objects. Continuing a consistent list may not be possible if the server configuration has changed or more than a few minutes have passed. The resourceVersion field returned when using this continue value will be identical to the value in the first response, unless you have received this token from an error message.
---@field remainingItemCount number remainingItemCount is the number of subsequent items in the list which are not included in this list response. If the list request contained label or field selectors, then the number of remaining items is unknown and the field will be left unset and omitted during serialization. If the list is complete (either because it is not chunking or because this is the last chunk), then there are no more remaining items and this field will be left unset and omitted during serialization. Servers older than v1.15 do not set this field. The intended use of the remainingItemCount is *estimating* the size of a collection. Clients should not rely on the remainingItemCount to be set or to be exact. +optional
---@field resourceVersion string String that identifies the server's internal version of this object that can be used by clients to determine when objects have changed. Value must be treated as opaque by clients and passed unmodified back to the server. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency +optional
---@field selfLink string Deprecated: selfLink is a legacy read-only field that is no longer populated by the system. +optional
---@field shardInfo v1.ShardInfo shardInfo is set when the list is a filtered subset of the full collection, as selected by a shard selector on the request. It echoes back the selector so clients can verify which shard they received and merge sharded responses. Clients should not cache sharded list responses as a full representation of the collection. This is an alpha field and requires enabling the ShardedListAndWatch feature gate. +featureGate=ShardedListAndWatch +optional

---@class v1.ManagedFieldsEntry
---@field apiVersion string APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
---@field fieldsType string FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"
---@field fieldsV1 v1.FieldsV1 FieldsV1 holds the first JSON version format as described in the "FieldsV1" type. +optional
---@field manager string Manager is an identifier of the workflow managing these fields.
---@field operation string Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are 'Apply' and 'Update'. +k8s:alpha(since: "1.37")=+k8s:required
---@field subresource string Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
---@field time v1.Time Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over. +optional

---@class v1.ObjectMeta
---@field annotations table<string, string> Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations +optional
---@field creationTimestamp v1.Time CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC. Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional +k8s:alpha(since: "1.37")=+k8s:immutable
---@field deletionGracePeriodSeconds number Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only. +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:immutable
---@field deletionTimestamp v1.Time DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted. This field is set by the server when a graceful deletion is requested by the user, and is not directly settable by a client. The resource is expected to be deleted (no longer visible from resource lists, and not reachable by name) after the time in this field, once the finalizers list is empty. As long as the finalizers list contains items, deletion is blocked. Once the deletionTimestamp is set, this value may not be unset or be set further into the future, although it may be shortened or the resource may be deleted prior to this time. For example, a user may request that a pod is deleted in 30 seconds. The Kubelet will react by sending a graceful termination signal to the containers in the pod. After that 30 seconds, the Kubelet will send a hard termination signal (SIGKILL) to the container and after cleanup, remove the pod from the API. In the presence of network partitions, this object may still exist after this timestamp, until an administrator or automated process can determine the resource is fully terminated. If not set, graceful deletion of the object has not been requested. Populated by the system when a graceful deletion is requested. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:immutable
---@field finalizers string[] Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order. Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list. +optional +patchStrategy=merge +listType=set
---@field generateName string GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server. If this field is specified and the generated name exists, the server will return a 409. Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency +optional
---@field generation number A sequence number representing a specific generation of the desired state. Populated by the system. Read-only. +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:minimum=0
---@field labels table<string, string> Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels +optional
---@field managedFields v1.ManagedFieldsEntry[] ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object. +optional +listType=atomic +k8s:alpha(since: "1.37")=+k8s:optional
---@field name string Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names +optional
---@field namespace string Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty. Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces +optional
---@field ownerReferences v1.OwnerReference[] List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller. +optional +patchMergeKey=uid +patchStrategy=merge +listType=map +listMapKey=uid +k8s:alpha(since:"1.37")=+k8s:optional
---@field resourceVersion string An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources. Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency +optional
---@field selfLink string Deprecated: selfLink is a legacy read-only field that is no longer populated by the system. +optional
---@field uid string UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations. Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids +optional +k8s:alpha(since: "1.37")=+k8s:optional +k8s:alpha(since: "1.37")=+k8s:immutable

---@class v1.OwnerReference
---@field apiVersion string API version of the referent. +k8s:alpha(since:"1.37")=+k8s:required
---@field blockOwnerDeletion boolean If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned. +optional
---@field controller boolean If true, this reference points to the managing controller. +optional
---@field kind string Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +k8s:alpha(since:"1.37")=+k8s:required
---@field name string Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names +k8s:alpha(since:"1.37")=+k8s:required
---@field uid string UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids +k8s:alpha(since:"1.37")=+k8s:required

---@class v1.ShardInfo
---@field selector string selector is the shard selector string from the request, echoed back so clients can verify which shard they received and merge responses from multiple shards. +required

---@class v1.Status
---@field code number Suggested HTTP return code for this status, 0 if not set. +optional
---@field details v1.StatusDetails Extended data associated with the reason. Each reason may define its own extended details. This field is optional and the data returned is not guaranteed to conform to any schema except that defined by the reason type. +optional
---@field message string A human-readable description of the status of this operation. +optional
---@field metadata v1.ListMeta Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional
---@field reason string A machine-readable description of why this operation is in the "Failure" status. If this value is empty there is no information available. A Reason clarifies an HTTP status code but does not override it. +optional
---@field status string Status of the operation. One of: "Success" or "Failure". More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status +optional

---@class v1.StatusCause
---@field field string The field of the resource that has caused this error, as named by its JSON serialization. May include dot and postfix notation for nested attributes. Arrays are zero-indexed. Fields may appear more than once in an array of causes due to fields having multiple errors. Optional. Examples: "name" - the field "name" on the current resource "items[0].name" - the field "name" on the first array entry in "items" +optional
---@field message string A human-readable description of the cause of the error. This field may be presented as-is to a reader. +optional
---@field reason string A machine-readable description of the cause of the error. If this value is empty there is no information available. +optional

---@class v1.StatusDetails
---@field causes v1.StatusCause[] The Causes array includes more details associated with the StatusReason failure. Not all StatusReasons may provide detailed causes. +optional +listType=atomic
---@field group string The group attribute of the resource associated with the status StatusReason. +optional
---@field kind string The kind attribute of the resource associated with the status StatusReason. On some operations may differ from the requested resource Kind. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional
---@field name string The name attribute of the resource associated with the status StatusReason (when there is a single name which can be described). +optional
---@field retryAfterSeconds number If specified, the time in seconds before the operation should be retried. Some errors may indicate the client must take an alternate action - for those errors this field may indicate how long to wait before taking the alternate action. +optional
---@field uid string UID of the resource. (when there is a single resource which can be described). More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids +optional

---@class v1.TypeMeta
---@field apiVersion string APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources +optional
---@field kind string Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds +optional

---@class kubernetes
local kubernetes = {}

--- parse a Kubernetes memory quantity, returns bytes
---@param quantity string a K8s resource.Quantity string, e.g. "1024Mi" or "1Gi"
---@return number bytes the quantity's value in bytes
function kubernetes.parse_memory(quantity) end

--- parse a Kubernetes CPU quantity, returns millicores
---@param quantity string a K8s resource.Quantity string, e.g. "100m" or "1"
---@return number millicores the quantity's value in millicores (1000m = 1 core)
function kubernetes.parse_cpu(quantity) end

--- format a byte count as a canonical K8s memory string (BinarySI: Ki/Mi/Gi/Ti)
---@param bytes number the memory amount in bytes
---@return string quantity the binary-SI quantity string, e.g. "2Gi"
function kubernetes.format_memory(bytes) end

--- format a byte count as a canonical K8s memory string (DecimalSI: k/M/G/T)
---@param bytes number the memory amount in bytes
---@return string quantity the decimal-SI quantity string, e.g. "2G"
function kubernetes.format_memory_si(bytes) end

--- format a millicore count as a canonical K8s CPU string (DecimalSI)
---@param millicores number the CPU amount in millicores (1000m = 1 core)
---@return string quantity the decimal-SI CPU quantity string, e.g. "500m" or "1"
function kubernetes.format_cpu(millicores) end

--- parse an RFC3339 time string, returns Unix timestamp
---@param timestr string an RFC3339-formatted timestamp, e.g. "2025-10-03T16:39:00Z"
---@return number timestamp seconds since the Unix epoch (UTC)
function kubernetes.parse_time(timestr) end

--- convert a Unix timestamp to RFC3339 string
---@param timestamp number seconds since the Unix epoch (UTC)
---@return string timestr the RFC3339-formatted timestamp
function kubernetes.format_time(timestamp) end

--- parse a duration string, returns seconds
---@param duration string a Go-style duration string, e.g. "5m" or "1h30m"
---@return number seconds the duration's length in seconds
function kubernetes.parse_duration(duration) end

--- convert seconds to a duration string
---@param seconds number a duration length in seconds
---@return string duration the Go-style duration string, e.g. "5m0s"
function kubernetes.format_duration(seconds) end

--- check if a Kubernetes object matches a GVK matcher
---@param obj table<string, any> the object to check; must have apiVersion and kind keys
---@param matcher kubernetes.GVKMatcher GVK matcher table with group, version and kind fields
---@return boolean matches true if obj's apiVersion and kind match matcher
function kubernetes.match_gvk(obj, matcher) end

--- ensure metadata.labels and annotations exist, returns updated obj
---@param obj table<string, any> the object whose metadata.labels and metadata.annotations should exist
---@return table<string, any> obj the object with metadata.labels and metadata.annotations guaranteed present
function kubernetes.ensure_metadata(obj) end

--- ensure metadata.labels and annotations exist, returns updated obj
---@param obj table<string, any> the object whose metadata.labels and metadata.annotations should exist
---@return table<string, any> obj the object with metadata.labels and metadata.annotations guaranteed present
function kubernetes.init_defaults(obj) end

--- add a label and return the updated obj
---@param obj table<string, any> the object to modify
---@param key string the label key to set
---@param value string the label value to set
---@return table<string, any> obj the object with the label set
function kubernetes.add_label(obj, key, value) end

--- add multiple labels and return the updated obj
---@param obj table<string, any> the object to modify
---@param labels table<string, any> table of label key to value to merge into obj.metadata.labels
---@return table<string, any> obj the object with the labels set
function kubernetes.add_labels(obj, labels) end

--- remove a label and return the updated obj
---@param obj table<string, any> the object to modify
---@param key string the label key to remove; a no-op if absent
---@return table<string, any> obj the object with the label removed
function kubernetes.remove_label(obj, key) end

--- return true if the label exists
---@param obj table<string, any> the object to check
---@param key string the label key to look for
---@return boolean ok true if obj.metadata.labels contains key
function kubernetes.has_label(obj, key) end

--- return the value of a label, or empty string if absent
---@param obj table<string, any> the object to read from
---@param key string the label key to look up
---@return string value the label value, or empty string if absent
function kubernetes.get_label(obj, key) end

--- add an annotation and return the updated obj
---@param obj table<string, any> the object to modify
---@param key string the annotation key to set
---@param value string the annotation value to set
---@return table<string, any> obj the object with the annotation set
function kubernetes.add_annotation(obj, key, value) end

--- add multiple annotations and return the updated obj
---@param obj table<string, any> the object to modify
---@param annotations table<string, any> table of annotation key to value to merge into obj.metadata.annotations
---@return table<string, any> obj the object with the annotations set
function kubernetes.add_annotations(obj, annotations) end

--- remove an annotation and return the updated obj
---@param obj table<string, any> the object to modify
---@param key string the annotation key to remove; a no-op if absent
---@return table<string, any> obj the object with the annotation removed
function kubernetes.remove_annotation(obj, key) end

--- return true if the annotation exists
---@param obj table<string, any> the object to check
---@param key string the annotation key to look for
---@return boolean ok true if obj.metadata.annotations contains key
function kubernetes.has_annotation(obj, key) end

--- return the value of an annotation, or empty string if absent
---@param obj table<string, any> the object to read from
---@param key string the annotation key to look up
---@return string value the annotation value, or empty string if absent
function kubernetes.get_annotation(obj, key) end

return kubernetes
