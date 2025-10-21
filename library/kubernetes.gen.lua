---@class kubernetes.GVKMatcher
---@field group string
---@field kind string
---@field version string

---@class appsv1.DaemonSet
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec appsv1.DaemonSetSpec
---@field status appsv1.DaemonSetStatus

---@class appsv1.DaemonSetCondition
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class appsv1.DaemonSetList
---@field TypeMeta v1.TypeMeta
---@field items appsv1.DaemonSet[]
---@field metadata v1.ListMeta

---@class appsv1.DaemonSetSpec
---@field minReadySeconds number
---@field revisionHistoryLimit number
---@field selector v1.LabelSelector
---@field template corev1.PodTemplateSpec
---@field updateStrategy appsv1.DaemonSetUpdateStrategy

---@class appsv1.DaemonSetStatus
---@field collisionCount number
---@field conditions appsv1.DaemonSetCondition[]
---@field currentNumberScheduled number
---@field desiredNumberScheduled number
---@field numberAvailable number
---@field numberMisscheduled number
---@field numberReady number
---@field numberUnavailable number
---@field observedGeneration number
---@field updatedNumberScheduled number

---@class appsv1.DaemonSetUpdateStrategy
---@field rollingUpdate appsv1.RollingUpdateDaemonSet
---@field type string

---@class appsv1.Deployment
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec appsv1.DeploymentSpec
---@field status appsv1.DeploymentStatus

---@class appsv1.DeploymentCondition
---@field lastTransitionTime v1.Time
---@field lastUpdateTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class appsv1.DeploymentList
---@field TypeMeta v1.TypeMeta
---@field items appsv1.Deployment[]
---@field metadata v1.ListMeta

---@class appsv1.DeploymentSpec
---@field minReadySeconds number
---@field paused boolean
---@field progressDeadlineSeconds number
---@field replicas number
---@field revisionHistoryLimit number
---@field selector v1.LabelSelector
---@field strategy appsv1.DeploymentStrategy
---@field template corev1.PodTemplateSpec

---@class appsv1.DeploymentStatus
---@field availableReplicas number
---@field collisionCount number
---@field conditions appsv1.DeploymentCondition[]
---@field observedGeneration number
---@field readyReplicas number
---@field replicas number
---@field terminatingReplicas number
---@field unavailableReplicas number
---@field updatedReplicas number

---@class appsv1.DeploymentStrategy
---@field rollingUpdate appsv1.RollingUpdateDeployment
---@field type string

---@class appsv1.ReplicaSet
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec appsv1.ReplicaSetSpec
---@field status appsv1.ReplicaSetStatus

---@class appsv1.ReplicaSetCondition
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class appsv1.ReplicaSetList
---@field TypeMeta v1.TypeMeta
---@field items appsv1.ReplicaSet[]
---@field metadata v1.ListMeta

---@class appsv1.ReplicaSetSpec
---@field minReadySeconds number
---@field replicas number
---@field selector v1.LabelSelector
---@field template corev1.PodTemplateSpec

---@class appsv1.ReplicaSetStatus
---@field availableReplicas number
---@field conditions appsv1.ReplicaSetCondition[]
---@field fullyLabeledReplicas number
---@field observedGeneration number
---@field readyReplicas number
---@field replicas number
---@field terminatingReplicas number

---@class appsv1.RollingUpdateDaemonSet
---@field maxSurge intstr.IntOrString
---@field maxUnavailable intstr.IntOrString

---@class appsv1.RollingUpdateDeployment
---@field maxSurge intstr.IntOrString
---@field maxUnavailable intstr.IntOrString

---@class appsv1.RollingUpdateStatefulSetStrategy
---@field maxUnavailable intstr.IntOrString
---@field partition number

---@class appsv1.StatefulSet
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec appsv1.StatefulSetSpec
---@field status appsv1.StatefulSetStatus

---@class appsv1.StatefulSetCondition
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class appsv1.StatefulSetList
---@field TypeMeta v1.TypeMeta
---@field items appsv1.StatefulSet[]
---@field metadata v1.ListMeta

---@class appsv1.StatefulSetOrdinals
---@field start number

---@class appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy
---@field whenDeleted string
---@field whenScaled string

---@class appsv1.StatefulSetSpec
---@field minReadySeconds number
---@field ordinals appsv1.StatefulSetOrdinals
---@field persistentVolumeClaimRetentionPolicy appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy
---@field podManagementPolicy string
---@field replicas number
---@field revisionHistoryLimit number
---@field selector v1.LabelSelector
---@field serviceName string
---@field template corev1.PodTemplateSpec
---@field updateStrategy appsv1.StatefulSetUpdateStrategy
---@field volumeClaimTemplates corev1.PersistentVolumeClaim[]

---@class appsv1.StatefulSetStatus
---@field availableReplicas number
---@field collisionCount number
---@field conditions appsv1.StatefulSetCondition[]
---@field currentReplicas number
---@field currentRevision string
---@field observedGeneration number
---@field readyReplicas number
---@field replicas number
---@field updateRevision string
---@field updatedReplicas number

---@class appsv1.StatefulSetUpdateStrategy
---@field rollingUpdate appsv1.RollingUpdateStatefulSetStrategy
---@field type string

---@class batchv1.CronJob
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec batchv1.CronJobSpec
---@field status batchv1.CronJobStatus

---@class batchv1.CronJobList
---@field TypeMeta v1.TypeMeta
---@field items batchv1.CronJob[]
---@field metadata v1.ListMeta

---@class batchv1.CronJobSpec
---@field concurrencyPolicy string
---@field failedJobsHistoryLimit number
---@field jobTemplate batchv1.JobTemplateSpec
---@field schedule string
---@field startingDeadlineSeconds number
---@field successfulJobsHistoryLimit number
---@field suspend boolean
---@field timeZone string

---@class batchv1.CronJobStatus
---@field active corev1.ObjectReference[]
---@field lastScheduleTime v1.Time
---@field lastSuccessfulTime v1.Time

---@class batchv1.Job
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec batchv1.JobSpec
---@field status batchv1.JobStatus

---@class batchv1.JobCondition
---@field lastProbeTime v1.Time
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class batchv1.JobList
---@field TypeMeta v1.TypeMeta
---@field items batchv1.Job[]
---@field metadata v1.ListMeta

---@class batchv1.JobSpec
---@field activeDeadlineSeconds number
---@field backoffLimit number
---@field backoffLimitPerIndex number
---@field completionMode string
---@field completions number
---@field managedBy string
---@field manualSelector boolean
---@field maxFailedIndexes number
---@field parallelism number
---@field podFailurePolicy batchv1.PodFailurePolicy
---@field podReplacementPolicy string
---@field selector v1.LabelSelector
---@field successPolicy batchv1.SuccessPolicy
---@field suspend boolean
---@field template corev1.PodTemplateSpec
---@field ttlSecondsAfterFinished number

---@class batchv1.JobStatus
---@field active number
---@field completedIndexes string
---@field completionTime v1.Time
---@field conditions batchv1.JobCondition[]
---@field failed number
---@field failedIndexes string
---@field ready number
---@field startTime v1.Time
---@field succeeded number
---@field terminating number
---@field uncountedTerminatedPods batchv1.UncountedTerminatedPods

---@class batchv1.JobTemplateSpec
---@field metadata v1.ObjectMeta
---@field spec batchv1.JobSpec

---@class batchv1.PodFailurePolicy
---@field rules batchv1.PodFailurePolicyRule[]

---@class batchv1.PodFailurePolicyOnExitCodesRequirement
---@field containerName string
---@field operator string
---@field values number[]

---@class batchv1.PodFailurePolicyOnPodConditionsPattern
---@field status string
---@field type string

---@class batchv1.PodFailurePolicyRule
---@field action string
---@field onExitCodes batchv1.PodFailurePolicyOnExitCodesRequirement
---@field onPodConditions batchv1.PodFailurePolicyOnPodConditionsPattern[]

---@class batchv1.SuccessPolicy
---@field rules batchv1.SuccessPolicyRule[]

---@class batchv1.SuccessPolicyRule
---@field succeededCount number
---@field succeededIndexes string

---@class batchv1.UncountedTerminatedPods
---@field failed string[]
---@field succeeded string[]

---@class corev1.AWSElasticBlockStoreVolumeSource
---@field fsType string
---@field partition number
---@field readOnly boolean
---@field volumeID string

---@class corev1.Affinity
---@field nodeAffinity corev1.NodeAffinity
---@field podAffinity corev1.PodAffinity
---@field podAntiAffinity corev1.PodAntiAffinity

---@class corev1.AppArmorProfile
---@field localhostProfile string
---@field type string

---@class corev1.AttachedVolume
---@field devicePath string
---@field name string

---@class corev1.AzureDiskVolumeSource
---@field cachingMode string
---@field diskName string
---@field diskURI string
---@field fsType string
---@field kind string
---@field readOnly boolean

---@class corev1.AzureFilePersistentVolumeSource
---@field readOnly boolean
---@field secretName string
---@field secretNamespace string
---@field shareName string

---@class corev1.AzureFileVolumeSource
---@field readOnly boolean
---@field secretName string
---@field shareName string

---@class corev1.CSIPersistentVolumeSource
---@field controllerExpandSecretRef corev1.SecretReference
---@field controllerPublishSecretRef corev1.SecretReference
---@field driver string
---@field fsType string
---@field nodeExpandSecretRef corev1.SecretReference
---@field nodePublishSecretRef corev1.SecretReference
---@field nodeStageSecretRef corev1.SecretReference
---@field readOnly boolean
---@field volumeAttributes table<string, string>
---@field volumeHandle string

---@class corev1.CSIVolumeSource
---@field driver string
---@field fsType string
---@field nodePublishSecretRef corev1.LocalObjectReference
---@field readOnly boolean
---@field volumeAttributes table<string, string>

---@class corev1.Capabilities
---@field add string[]
---@field drop string[]

---@class corev1.CephFSPersistentVolumeSource
---@field monitors string[]
---@field path string
---@field readOnly boolean
---@field secretFile string
---@field secretRef corev1.SecretReference
---@field user string

---@class corev1.CephFSVolumeSource
---@field monitors string[]
---@field path string
---@field readOnly boolean
---@field secretFile string
---@field secretRef corev1.LocalObjectReference
---@field user string

---@class corev1.CinderPersistentVolumeSource
---@field fsType string
---@field readOnly boolean
---@field secretRef corev1.SecretReference
---@field volumeID string

---@class corev1.CinderVolumeSource
---@field fsType string
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field volumeID string

---@class corev1.ClientIPConfig
---@field timeoutSeconds number

---@class corev1.ClusterTrustBundleProjection
---@field labelSelector v1.LabelSelector
---@field name string
---@field optional boolean
---@field path string
---@field signerName string

---@class corev1.ConfigMap
---@field TypeMeta v1.TypeMeta
---@field binaryData table<string, number[]>
---@field data table<string, string>
---@field immutable boolean
---@field metadata v1.ObjectMeta

---@class corev1.ConfigMapEnvSource
---@field LocalObjectReference corev1.LocalObjectReference
---@field optional boolean

---@class corev1.ConfigMapKeySelector
---@field LocalObjectReference corev1.LocalObjectReference
---@field key string
---@field optional boolean

---@class corev1.ConfigMapList
---@field TypeMeta v1.TypeMeta
---@field items corev1.ConfigMap[]
---@field metadata v1.ListMeta

---@class corev1.ConfigMapNodeConfigSource
---@field kubeletConfigKey string
---@field name string
---@field namespace string
---@field resourceVersion string
---@field uid string

---@class corev1.ConfigMapProjection
---@field LocalObjectReference corev1.LocalObjectReference
---@field items corev1.KeyToPath[]
---@field optional boolean

---@class corev1.ConfigMapVolumeSource
---@field LocalObjectReference corev1.LocalObjectReference
---@field defaultMode number
---@field items corev1.KeyToPath[]
---@field optional boolean

---@class corev1.Container
---@field args string[]
---@field command string[]
---@field env corev1.EnvVar[]
---@field envFrom corev1.EnvFromSource[]
---@field image string
---@field imagePullPolicy string
---@field lifecycle corev1.Lifecycle
---@field livenessProbe corev1.Probe
---@field name string
---@field ports corev1.ContainerPort[]
---@field readinessProbe corev1.Probe
---@field resizePolicy corev1.ContainerResizePolicy[]
---@field resources corev1.ResourceRequirements
---@field restartPolicy string
---@field restartPolicyRules corev1.ContainerRestartRule[]
---@field securityContext corev1.SecurityContext
---@field startupProbe corev1.Probe
---@field stdin boolean
---@field stdinOnce boolean
---@field terminationMessagePath string
---@field terminationMessagePolicy string
---@field tty boolean
---@field volumeDevices corev1.VolumeDevice[]
---@field volumeMounts corev1.VolumeMount[]
---@field workingDir string

---@class corev1.ContainerExtendedResourceRequest
---@field containerName string
---@field requestName string
---@field resourceName string

---@class corev1.ContainerImage
---@field names string[]
---@field sizeBytes number

---@class corev1.ContainerPort
---@field containerPort number
---@field hostIP string
---@field hostPort number
---@field name string
---@field protocol string

---@class corev1.ContainerResizePolicy
---@field resourceName string
---@field restartPolicy string

---@class corev1.ContainerRestartRule
---@field action string
---@field exitCodes corev1.ContainerRestartRuleOnExitCodes

---@class corev1.ContainerRestartRuleOnExitCodes
---@field operator string
---@field values number[]

---@class corev1.ContainerState
---@field running corev1.ContainerStateRunning
---@field terminated corev1.ContainerStateTerminated
---@field waiting corev1.ContainerStateWaiting

---@class corev1.ContainerStateRunning
---@field startedAt v1.Time

---@class corev1.ContainerStateTerminated
---@field containerID string
---@field exitCode number
---@field finishedAt v1.Time
---@field message string
---@field reason string
---@field signal number
---@field startedAt v1.Time

---@class corev1.ContainerStateWaiting
---@field message string
---@field reason string

---@class corev1.ContainerStatus
---@field allocatedResources table<string, resource.Quantity>
---@field allocatedResourcesStatus corev1.ResourceStatus[]
---@field containerID string
---@field image string
---@field imageID string
---@field lastState corev1.ContainerState
---@field name string
---@field ready boolean
---@field resources corev1.ResourceRequirements
---@field restartCount number
---@field started boolean
---@field state corev1.ContainerState
---@field stopSignal string
---@field user corev1.ContainerUser
---@field volumeMounts corev1.VolumeMountStatus[]

---@class corev1.ContainerUser
---@field linux corev1.LinuxContainerUser

---@class corev1.DaemonEndpoint
---@field Port number

---@class corev1.DownwardAPIProjection
---@field items corev1.DownwardAPIVolumeFile[]

---@class corev1.DownwardAPIVolumeFile
---@field fieldRef corev1.ObjectFieldSelector
---@field mode number
---@field path string
---@field resourceFieldRef corev1.ResourceFieldSelector

---@class corev1.DownwardAPIVolumeSource
---@field defaultMode number
---@field items corev1.DownwardAPIVolumeFile[]

---@class corev1.EmptyDirVolumeSource
---@field medium string
---@field sizeLimit resource.Quantity

---@class corev1.EnvFromSource
---@field configMapRef corev1.ConfigMapEnvSource
---@field prefix string
---@field secretRef corev1.SecretEnvSource

---@class corev1.EnvVar
---@field name string
---@field value string
---@field valueFrom corev1.EnvVarSource

---@class corev1.EnvVarSource
---@field configMapKeyRef corev1.ConfigMapKeySelector
---@field fieldRef corev1.ObjectFieldSelector
---@field fileKeyRef corev1.FileKeySelector
---@field resourceFieldRef corev1.ResourceFieldSelector
---@field secretKeyRef corev1.SecretKeySelector

---@class corev1.EphemeralContainer
---@field EphemeralContainerCommon corev1.EphemeralContainerCommon
---@field targetContainerName string

---@class corev1.EphemeralContainerCommon
---@field args string[]
---@field command string[]
---@field env corev1.EnvVar[]
---@field envFrom corev1.EnvFromSource[]
---@field image string
---@field imagePullPolicy string
---@field lifecycle corev1.Lifecycle
---@field livenessProbe corev1.Probe
---@field name string
---@field ports corev1.ContainerPort[]
---@field readinessProbe corev1.Probe
---@field resizePolicy corev1.ContainerResizePolicy[]
---@field resources corev1.ResourceRequirements
---@field restartPolicy string
---@field restartPolicyRules corev1.ContainerRestartRule[]
---@field securityContext corev1.SecurityContext
---@field startupProbe corev1.Probe
---@field stdin boolean
---@field stdinOnce boolean
---@field terminationMessagePath string
---@field terminationMessagePolicy string
---@field tty boolean
---@field volumeDevices corev1.VolumeDevice[]
---@field volumeMounts corev1.VolumeMount[]
---@field workingDir string

---@class corev1.EphemeralVolumeSource
---@field volumeClaimTemplate corev1.PersistentVolumeClaimTemplate

---@class corev1.ExecAction
---@field command string[]

---@class corev1.FCVolumeSource
---@field fsType string
---@field lun number
---@field readOnly boolean
---@field targetWWNs string[]
---@field wwids string[]

---@class corev1.FileKeySelector
---@field key string
---@field optional boolean
---@field path string
---@field volumeName string

---@class corev1.FlexPersistentVolumeSource
---@field driver string
---@field fsType string
---@field options table<string, string>
---@field readOnly boolean
---@field secretRef corev1.SecretReference

---@class corev1.FlexVolumeSource
---@field driver string
---@field fsType string
---@field options table<string, string>
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference

---@class corev1.FlockerVolumeSource
---@field datasetName string
---@field datasetUUID string

---@class corev1.GCEPersistentDiskVolumeSource
---@field fsType string
---@field partition number
---@field pdName string
---@field readOnly boolean

---@class corev1.GRPCAction
---@field port number
---@field service string

---@class corev1.GitRepoVolumeSource
---@field directory string
---@field repository string
---@field revision string

---@class corev1.GlusterfsPersistentVolumeSource
---@field endpoints string
---@field endpointsNamespace string
---@field path string
---@field readOnly boolean

---@class corev1.GlusterfsVolumeSource
---@field endpoints string
---@field path string
---@field readOnly boolean

---@class corev1.HTTPGetAction
---@field host string
---@field httpHeaders corev1.HTTPHeader[]
---@field path string
---@field port intstr.IntOrString
---@field scheme string

---@class corev1.HTTPHeader
---@field name string
---@field value string

---@class corev1.HostAlias
---@field hostnames string[]
---@field ip string

---@class corev1.HostIP
---@field ip string

---@class corev1.HostPathVolumeSource
---@field path string
---@field type string

---@class corev1.ISCSIPersistentVolumeSource
---@field chapAuthDiscovery boolean
---@field chapAuthSession boolean
---@field fsType string
---@field initiatorName string
---@field iqn string
---@field iscsiInterface string
---@field lun number
---@field portals string[]
---@field readOnly boolean
---@field secretRef corev1.SecretReference
---@field targetPortal string

---@class corev1.ISCSIVolumeSource
---@field chapAuthDiscovery boolean
---@field chapAuthSession boolean
---@field fsType string
---@field initiatorName string
---@field iqn string
---@field iscsiInterface string
---@field lun number
---@field portals string[]
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field targetPortal string

---@class corev1.ImageVolumeSource
---@field pullPolicy string
---@field reference string

---@class corev1.KeyToPath
---@field key string
---@field mode number
---@field path string

---@class corev1.Lifecycle
---@field postStart corev1.LifecycleHandler
---@field preStop corev1.LifecycleHandler
---@field stopSignal string

---@class corev1.LifecycleHandler
---@field exec corev1.ExecAction
---@field httpGet corev1.HTTPGetAction
---@field sleep corev1.SleepAction
---@field tcpSocket corev1.TCPSocketAction

---@class corev1.LinuxContainerUser
---@field gid number
---@field supplementalGroups number[]
---@field uid number

---@class corev1.LoadBalancerIngress
---@field hostname string
---@field ip string
---@field ipMode string
---@field ports corev1.PortStatus[]

---@class corev1.LoadBalancerStatus
---@field ingress corev1.LoadBalancerIngress[]

---@class corev1.LocalObjectReference
---@field name string

---@class corev1.LocalVolumeSource
---@field fsType string
---@field path string

---@class corev1.ModifyVolumeStatus
---@field status string
---@field targetVolumeAttributesClassName string

---@class corev1.NFSVolumeSource
---@field path string
---@field readOnly boolean
---@field server string

---@class corev1.Namespace
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.NamespaceSpec
---@field status corev1.NamespaceStatus

---@class corev1.NamespaceCondition
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class corev1.NamespaceList
---@field TypeMeta v1.TypeMeta
---@field items corev1.Namespace[]
---@field metadata v1.ListMeta

---@class corev1.NamespaceSpec
---@field finalizers string[]

---@class corev1.NamespaceStatus
---@field conditions corev1.NamespaceCondition[]
---@field phase string

---@class corev1.Node
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.NodeSpec
---@field status corev1.NodeStatus

---@class corev1.NodeAddress
---@field address string
---@field type string

---@class corev1.NodeAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.PreferredSchedulingTerm[]
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.NodeSelector

---@class corev1.NodeCondition
---@field lastHeartbeatTime v1.Time
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class corev1.NodeConfigSource
---@field configMap corev1.ConfigMapNodeConfigSource

---@class corev1.NodeConfigStatus
---@field active corev1.NodeConfigSource
---@field assigned corev1.NodeConfigSource
---@field error string
---@field lastKnownGood corev1.NodeConfigSource

---@class corev1.NodeDaemonEndpoints
---@field kubeletEndpoint corev1.DaemonEndpoint

---@class corev1.NodeFeatures
---@field supplementalGroupsPolicy boolean

---@class corev1.NodeList
---@field TypeMeta v1.TypeMeta
---@field items corev1.Node[]
---@field metadata v1.ListMeta

---@class corev1.NodeRuntimeHandler
---@field features corev1.NodeRuntimeHandlerFeatures
---@field name string

---@class corev1.NodeRuntimeHandlerFeatures
---@field recursiveReadOnlyMounts boolean
---@field userNamespaces boolean

---@class corev1.NodeSelector
---@field nodeSelectorTerms corev1.NodeSelectorTerm[]

---@class corev1.NodeSelectorRequirement
---@field key string
---@field operator string
---@field values string[]

---@class corev1.NodeSelectorTerm
---@field matchExpressions corev1.NodeSelectorRequirement[]
---@field matchFields corev1.NodeSelectorRequirement[]

---@class corev1.NodeSpec
---@field configSource corev1.NodeConfigSource
---@field externalID string
---@field podCIDR string
---@field podCIDRs string[]
---@field providerID string
---@field taints corev1.Taint[]
---@field unschedulable boolean

---@class corev1.NodeStatus
---@field addresses corev1.NodeAddress[]
---@field allocatable table<string, resource.Quantity>
---@field capacity table<string, resource.Quantity>
---@field conditions corev1.NodeCondition[]
---@field config corev1.NodeConfigStatus
---@field daemonEndpoints corev1.NodeDaemonEndpoints
---@field features corev1.NodeFeatures
---@field images corev1.ContainerImage[]
---@field nodeInfo corev1.NodeSystemInfo
---@field phase string
---@field runtimeHandlers corev1.NodeRuntimeHandler[]
---@field volumesAttached corev1.AttachedVolume[]
---@field volumesInUse string[]

---@class corev1.NodeSwapStatus
---@field capacity number

---@class corev1.NodeSystemInfo
---@field architecture string
---@field bootID string
---@field containerRuntimeVersion string
---@field kernelVersion string
---@field kubeProxyVersion string
---@field kubeletVersion string
---@field machineID string
---@field operatingSystem string
---@field osImage string
---@field swap corev1.NodeSwapStatus
---@field systemUUID string

---@class corev1.ObjectFieldSelector
---@field apiVersion string
---@field fieldPath string

---@class corev1.ObjectReference
---@field apiVersion string
---@field fieldPath string
---@field kind string
---@field name string
---@field namespace string
---@field resourceVersion string
---@field uid string

---@class corev1.PersistentVolume
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.PersistentVolumeSpec
---@field status corev1.PersistentVolumeStatus

---@class corev1.PersistentVolumeClaim
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.PersistentVolumeClaimSpec
---@field status corev1.PersistentVolumeClaimStatus

---@class corev1.PersistentVolumeClaimCondition
---@field lastProbeTime v1.Time
---@field lastTransitionTime v1.Time
---@field message string
---@field reason string
---@field status string
---@field type string

---@class corev1.PersistentVolumeClaimList
---@field TypeMeta v1.TypeMeta
---@field items corev1.PersistentVolumeClaim[]
---@field metadata v1.ListMeta

---@class corev1.PersistentVolumeClaimSpec
---@field accessModes string[]
---@field dataSource corev1.TypedLocalObjectReference
---@field dataSourceRef corev1.TypedObjectReference
---@field resources corev1.VolumeResourceRequirements
---@field selector v1.LabelSelector
---@field storageClassName string
---@field volumeAttributesClassName string
---@field volumeMode string
---@field volumeName string

---@class corev1.PersistentVolumeClaimStatus
---@field accessModes string[]
---@field allocatedResourceStatuses table<string, string>
---@field allocatedResources table<string, resource.Quantity>
---@field capacity table<string, resource.Quantity>
---@field conditions corev1.PersistentVolumeClaimCondition[]
---@field currentVolumeAttributesClassName string
---@field modifyVolumeStatus corev1.ModifyVolumeStatus
---@field phase string

---@class corev1.PersistentVolumeClaimTemplate
---@field metadata v1.ObjectMeta
---@field spec corev1.PersistentVolumeClaimSpec

---@class corev1.PersistentVolumeClaimVolumeSource
---@field claimName string
---@field readOnly boolean

---@class corev1.PersistentVolumeList
---@field TypeMeta v1.TypeMeta
---@field items corev1.PersistentVolume[]
---@field metadata v1.ListMeta

---@class corev1.PersistentVolumeSource
---@field awsElasticBlockStore corev1.AWSElasticBlockStoreVolumeSource
---@field azureDisk corev1.AzureDiskVolumeSource
---@field azureFile corev1.AzureFilePersistentVolumeSource
---@field cephfs corev1.CephFSPersistentVolumeSource
---@field cinder corev1.CinderPersistentVolumeSource
---@field csi corev1.CSIPersistentVolumeSource
---@field fc corev1.FCVolumeSource
---@field flexVolume corev1.FlexPersistentVolumeSource
---@field flocker corev1.FlockerVolumeSource
---@field gcePersistentDisk corev1.GCEPersistentDiskVolumeSource
---@field glusterfs corev1.GlusterfsPersistentVolumeSource
---@field hostPath corev1.HostPathVolumeSource
---@field iscsi corev1.ISCSIPersistentVolumeSource
---@field local corev1.LocalVolumeSource
---@field nfs corev1.NFSVolumeSource
---@field photonPersistentDisk corev1.PhotonPersistentDiskVolumeSource
---@field portworxVolume corev1.PortworxVolumeSource
---@field quobyte corev1.QuobyteVolumeSource
---@field rbd corev1.RBDPersistentVolumeSource
---@field scaleIO corev1.ScaleIOPersistentVolumeSource
---@field storageos corev1.StorageOSPersistentVolumeSource
---@field vsphereVolume corev1.VsphereVirtualDiskVolumeSource

---@class corev1.PersistentVolumeSpec
---@field PersistentVolumeSource corev1.PersistentVolumeSource
---@field accessModes string[]
---@field capacity table<string, resource.Quantity>
---@field claimRef corev1.ObjectReference
---@field mountOptions string[]
---@field nodeAffinity corev1.VolumeNodeAffinity
---@field persistentVolumeReclaimPolicy string
---@field storageClassName string
---@field volumeAttributesClassName string
---@field volumeMode string

---@class corev1.PersistentVolumeStatus
---@field lastPhaseTransitionTime v1.Time
---@field message string
---@field phase string
---@field reason string

---@class corev1.PhotonPersistentDiskVolumeSource
---@field fsType string
---@field pdID string

---@class corev1.Pod
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.PodSpec
---@field status corev1.PodStatus

---@class corev1.PodAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.WeightedPodAffinityTerm[]
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.PodAffinityTerm[]

---@class corev1.PodAffinityTerm
---@field labelSelector v1.LabelSelector
---@field matchLabelKeys string[]
---@field mismatchLabelKeys string[]
---@field namespaceSelector v1.LabelSelector
---@field namespaces string[]
---@field topologyKey string

---@class corev1.PodAntiAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.WeightedPodAffinityTerm[]
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.PodAffinityTerm[]

---@class corev1.PodCertificateProjection
---@field certificateChainPath string
---@field credentialBundlePath string
---@field keyPath string
---@field keyType string
---@field maxExpirationSeconds number
---@field signerName string

---@class corev1.PodCondition
---@field lastProbeTime v1.Time
---@field lastTransitionTime v1.Time
---@field message string
---@field observedGeneration number
---@field reason string
---@field status string
---@field type string

---@class corev1.PodDNSConfig
---@field nameservers string[]
---@field options corev1.PodDNSConfigOption[]
---@field searches string[]

---@class corev1.PodDNSConfigOption
---@field name string
---@field value string

---@class corev1.PodExtendedResourceClaimStatus
---@field requestMappings corev1.ContainerExtendedResourceRequest[]
---@field resourceClaimName string

---@class corev1.PodIP
---@field ip string

---@class corev1.PodList
---@field TypeMeta v1.TypeMeta
---@field items corev1.Pod[]
---@field metadata v1.ListMeta

---@class corev1.PodOS
---@field name string

---@class corev1.PodReadinessGate
---@field conditionType string

---@class corev1.PodResourceClaim
---@field name string
---@field resourceClaimName string
---@field resourceClaimTemplateName string

---@class corev1.PodResourceClaimStatus
---@field name string
---@field resourceClaimName string

---@class corev1.PodSchedulingGate
---@field name string

---@class corev1.PodSecurityContext
---@field appArmorProfile corev1.AppArmorProfile
---@field fsGroup number
---@field fsGroupChangePolicy string
---@field runAsGroup number
---@field runAsNonRoot boolean
---@field runAsUser number
---@field seLinuxChangePolicy string
---@field seLinuxOptions corev1.SELinuxOptions
---@field seccompProfile corev1.SeccompProfile
---@field supplementalGroups number[]
---@field supplementalGroupsPolicy string
---@field sysctls corev1.Sysctl[]
---@field windowsOptions corev1.WindowsSecurityContextOptions

---@class corev1.PodSpec
---@field activeDeadlineSeconds number
---@field affinity corev1.Affinity
---@field automountServiceAccountToken boolean
---@field containers corev1.Container[]
---@field dnsConfig corev1.PodDNSConfig
---@field dnsPolicy string
---@field enableServiceLinks boolean
---@field ephemeralContainers corev1.EphemeralContainer[]
---@field hostAliases corev1.HostAlias[]
---@field hostIPC boolean
---@field hostNetwork boolean
---@field hostPID boolean
---@field hostUsers boolean
---@field hostname string
---@field hostnameOverride string
---@field imagePullSecrets corev1.LocalObjectReference[]
---@field initContainers corev1.Container[]
---@field nodeName string
---@field nodeSelector table<string, string>
---@field os corev1.PodOS
---@field overhead table<string, resource.Quantity>
---@field preemptionPolicy string
---@field priority number
---@field priorityClassName string
---@field readinessGates corev1.PodReadinessGate[]
---@field resourceClaims corev1.PodResourceClaim[]
---@field resources corev1.ResourceRequirements
---@field restartPolicy string
---@field runtimeClassName string
---@field schedulerName string
---@field schedulingGates corev1.PodSchedulingGate[]
---@field securityContext corev1.PodSecurityContext
---@field serviceAccount string
---@field serviceAccountName string
---@field setHostnameAsFQDN boolean
---@field shareProcessNamespace boolean
---@field subdomain string
---@field terminationGracePeriodSeconds number
---@field tolerations corev1.Toleration[]
---@field topologySpreadConstraints corev1.TopologySpreadConstraint[]
---@field volumes corev1.Volume[]

---@class corev1.PodStatus
---@field conditions corev1.PodCondition[]
---@field containerStatuses corev1.ContainerStatus[]
---@field ephemeralContainerStatuses corev1.ContainerStatus[]
---@field extendedResourceClaimStatus corev1.PodExtendedResourceClaimStatus
---@field hostIP string
---@field hostIPs corev1.HostIP[]
---@field initContainerStatuses corev1.ContainerStatus[]
---@field message string
---@field nominatedNodeName string
---@field observedGeneration number
---@field phase string
---@field podIP string
---@field podIPs corev1.PodIP[]
---@field qosClass string
---@field reason string
---@field resize string
---@field resourceClaimStatuses corev1.PodResourceClaimStatus[]
---@field startTime v1.Time

---@class corev1.PodTemplateSpec
---@field metadata v1.ObjectMeta
---@field spec corev1.PodSpec

---@class corev1.PortStatus
---@field error string
---@field port number
---@field protocol string

---@class corev1.PortworxVolumeSource
---@field fsType string
---@field readOnly boolean
---@field volumeID string

---@class corev1.PreferredSchedulingTerm
---@field preference corev1.NodeSelectorTerm
---@field weight number

---@class corev1.Probe
---@field ProbeHandler corev1.ProbeHandler
---@field failureThreshold number
---@field initialDelaySeconds number
---@field periodSeconds number
---@field successThreshold number
---@field terminationGracePeriodSeconds number
---@field timeoutSeconds number

---@class corev1.ProbeHandler
---@field exec corev1.ExecAction
---@field grpc corev1.GRPCAction
---@field httpGet corev1.HTTPGetAction
---@field tcpSocket corev1.TCPSocketAction

---@class corev1.ProjectedVolumeSource
---@field defaultMode number
---@field sources corev1.VolumeProjection[]

---@class corev1.QuobyteVolumeSource
---@field group string
---@field readOnly boolean
---@field registry string
---@field tenant string
---@field user string
---@field volume string

---@class corev1.RBDPersistentVolumeSource
---@field fsType string
---@field image string
---@field keyring string
---@field monitors string[]
---@field pool string
---@field readOnly boolean
---@field secretRef corev1.SecretReference
---@field user string

---@class corev1.RBDVolumeSource
---@field fsType string
---@field image string
---@field keyring string
---@field monitors string[]
---@field pool string
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field user string

---@class corev1.ResourceClaim
---@field name string
---@field request string

---@class corev1.ResourceFieldSelector
---@field containerName string
---@field divisor resource.Quantity
---@field resource string

---@class corev1.ResourceHealth
---@field health string
---@field resourceID string

---@class corev1.ResourceRequirements
---@field claims corev1.ResourceClaim[]
---@field limits table<string, resource.Quantity>
---@field requests table<string, resource.Quantity>

---@class corev1.ResourceStatus
---@field name string
---@field resources corev1.ResourceHealth[]

---@class corev1.SELinuxOptions
---@field level string
---@field role string
---@field type string
---@field user string

---@class corev1.ScaleIOPersistentVolumeSource
---@field fsType string
---@field gateway string
---@field protectionDomain string
---@field readOnly boolean
---@field secretRef corev1.SecretReference
---@field sslEnabled boolean
---@field storageMode string
---@field storagePool string
---@field system string
---@field volumeName string

---@class corev1.ScaleIOVolumeSource
---@field fsType string
---@field gateway string
---@field protectionDomain string
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field sslEnabled boolean
---@field storageMode string
---@field storagePool string
---@field system string
---@field volumeName string

---@class corev1.SeccompProfile
---@field localhostProfile string
---@field type string

---@class corev1.Secret
---@field TypeMeta v1.TypeMeta
---@field data table<string, number[]>
---@field immutable boolean
---@field metadata v1.ObjectMeta
---@field stringData table<string, string>
---@field type string

---@class corev1.SecretEnvSource
---@field LocalObjectReference corev1.LocalObjectReference
---@field optional boolean

---@class corev1.SecretKeySelector
---@field LocalObjectReference corev1.LocalObjectReference
---@field key string
---@field optional boolean

---@class corev1.SecretList
---@field TypeMeta v1.TypeMeta
---@field items corev1.Secret[]
---@field metadata v1.ListMeta

---@class corev1.SecretProjection
---@field LocalObjectReference corev1.LocalObjectReference
---@field items corev1.KeyToPath[]
---@field optional boolean

---@class corev1.SecretReference
---@field name string
---@field namespace string

---@class corev1.SecretVolumeSource
---@field defaultMode number
---@field items corev1.KeyToPath[]
---@field optional boolean
---@field secretName string

---@class corev1.SecurityContext
---@field allowPrivilegeEscalation boolean
---@field appArmorProfile corev1.AppArmorProfile
---@field capabilities corev1.Capabilities
---@field privileged boolean
---@field procMount string
---@field readOnlyRootFilesystem boolean
---@field runAsGroup number
---@field runAsNonRoot boolean
---@field runAsUser number
---@field seLinuxOptions corev1.SELinuxOptions
---@field seccompProfile corev1.SeccompProfile
---@field windowsOptions corev1.WindowsSecurityContextOptions

---@class corev1.Service
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec corev1.ServiceSpec
---@field status corev1.ServiceStatus

---@class corev1.ServiceAccount
---@field TypeMeta v1.TypeMeta
---@field automountServiceAccountToken boolean
---@field imagePullSecrets corev1.LocalObjectReference[]
---@field metadata v1.ObjectMeta
---@field secrets corev1.ObjectReference[]

---@class corev1.ServiceAccountList
---@field TypeMeta v1.TypeMeta
---@field items corev1.ServiceAccount[]
---@field metadata v1.ListMeta

---@class corev1.ServiceAccountTokenProjection
---@field audience string
---@field expirationSeconds number
---@field path string

---@class corev1.ServiceList
---@field TypeMeta v1.TypeMeta
---@field items corev1.Service[]
---@field metadata v1.ListMeta

---@class corev1.ServicePort
---@field appProtocol string
---@field name string
---@field nodePort number
---@field port number
---@field protocol string
---@field targetPort intstr.IntOrString

---@class corev1.ServiceSpec
---@field allocateLoadBalancerNodePorts boolean
---@field clusterIP string
---@field clusterIPs string[]
---@field externalIPs string[]
---@field externalName string
---@field externalTrafficPolicy string
---@field healthCheckNodePort number
---@field internalTrafficPolicy string
---@field ipFamilies string[]
---@field ipFamilyPolicy string
---@field loadBalancerClass string
---@field loadBalancerIP string
---@field loadBalancerSourceRanges string[]
---@field ports corev1.ServicePort[]
---@field publishNotReadyAddresses boolean
---@field selector table<string, string>
---@field sessionAffinity string
---@field sessionAffinityConfig corev1.SessionAffinityConfig
---@field trafficDistribution string
---@field type string

---@class corev1.ServiceStatus
---@field conditions v1.Condition[]
---@field loadBalancer corev1.LoadBalancerStatus

---@class corev1.SessionAffinityConfig
---@field clientIP corev1.ClientIPConfig

---@class corev1.SleepAction
---@field seconds number

---@class corev1.StorageOSPersistentVolumeSource
---@field fsType string
---@field readOnly boolean
---@field secretRef corev1.ObjectReference
---@field volumeName string
---@field volumeNamespace string

---@class corev1.StorageOSVolumeSource
---@field fsType string
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field volumeName string
---@field volumeNamespace string

---@class corev1.Sysctl
---@field name string
---@field value string

---@class corev1.TCPSocketAction
---@field host string
---@field port intstr.IntOrString

---@class corev1.Taint
---@field effect string
---@field key string
---@field timeAdded v1.Time
---@field value string

---@class corev1.Toleration
---@field effect string
---@field key string
---@field operator string
---@field tolerationSeconds number
---@field value string

---@class corev1.TopologySpreadConstraint
---@field labelSelector v1.LabelSelector
---@field matchLabelKeys string[]
---@field maxSkew number
---@field minDomains number
---@field nodeAffinityPolicy string
---@field nodeTaintsPolicy string
---@field topologyKey string
---@field whenUnsatisfiable string

---@class corev1.TypedLocalObjectReference
---@field apiGroup string
---@field kind string
---@field name string

---@class corev1.TypedObjectReference
---@field apiGroup string
---@field kind string
---@field name string
---@field namespace string

---@class corev1.Volume
---@field VolumeSource corev1.VolumeSource
---@field name string

---@class corev1.VolumeDevice
---@field devicePath string
---@field name string

---@class corev1.VolumeMount
---@field mountPath string
---@field mountPropagation string
---@field name string
---@field readOnly boolean
---@field recursiveReadOnly string
---@field subPath string
---@field subPathExpr string

---@class corev1.VolumeMountStatus
---@field mountPath string
---@field name string
---@field readOnly boolean
---@field recursiveReadOnly string

---@class corev1.VolumeNodeAffinity
---@field required corev1.NodeSelector

---@class corev1.VolumeProjection
---@field clusterTrustBundle corev1.ClusterTrustBundleProjection
---@field configMap corev1.ConfigMapProjection
---@field downwardAPI corev1.DownwardAPIProjection
---@field podCertificate corev1.PodCertificateProjection
---@field secret corev1.SecretProjection
---@field serviceAccountToken corev1.ServiceAccountTokenProjection

---@class corev1.VolumeResourceRequirements
---@field limits table<string, resource.Quantity>
---@field requests table<string, resource.Quantity>

---@class corev1.VolumeSource
---@field awsElasticBlockStore corev1.AWSElasticBlockStoreVolumeSource
---@field azureDisk corev1.AzureDiskVolumeSource
---@field azureFile corev1.AzureFileVolumeSource
---@field cephfs corev1.CephFSVolumeSource
---@field cinder corev1.CinderVolumeSource
---@field configMap corev1.ConfigMapVolumeSource
---@field csi corev1.CSIVolumeSource
---@field downwardAPI corev1.DownwardAPIVolumeSource
---@field emptyDir corev1.EmptyDirVolumeSource
---@field ephemeral corev1.EphemeralVolumeSource
---@field fc corev1.FCVolumeSource
---@field flexVolume corev1.FlexVolumeSource
---@field flocker corev1.FlockerVolumeSource
---@field gcePersistentDisk corev1.GCEPersistentDiskVolumeSource
---@field gitRepo corev1.GitRepoVolumeSource
---@field glusterfs corev1.GlusterfsVolumeSource
---@field hostPath corev1.HostPathVolumeSource
---@field image corev1.ImageVolumeSource
---@field iscsi corev1.ISCSIVolumeSource
---@field nfs corev1.NFSVolumeSource
---@field persistentVolumeClaim corev1.PersistentVolumeClaimVolumeSource
---@field photonPersistentDisk corev1.PhotonPersistentDiskVolumeSource
---@field portworxVolume corev1.PortworxVolumeSource
---@field projected corev1.ProjectedVolumeSource
---@field quobyte corev1.QuobyteVolumeSource
---@field rbd corev1.RBDVolumeSource
---@field scaleIO corev1.ScaleIOVolumeSource
---@field secret corev1.SecretVolumeSource
---@field storageos corev1.StorageOSVolumeSource
---@field vsphereVolume corev1.VsphereVirtualDiskVolumeSource

---@class corev1.VsphereVirtualDiskVolumeSource
---@field fsType string
---@field storagePolicyID string
---@field storagePolicyName string
---@field volumePath string

---@class corev1.WeightedPodAffinityTerm
---@field podAffinityTerm corev1.PodAffinityTerm
---@field weight number

---@class corev1.WindowsSecurityContextOptions
---@field gmsaCredentialSpec string
---@field gmsaCredentialSpecName string
---@field hostProcess boolean
---@field runAsUserName string

---@class networkingv1.HTTPIngressPath
---@field backend networkingv1.IngressBackend
---@field path string
---@field pathType string

---@class networkingv1.HTTPIngressRuleValue
---@field paths networkingv1.HTTPIngressPath[]

---@class networkingv1.IPBlock
---@field cidr string
---@field except string[]

---@class networkingv1.Ingress
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec networkingv1.IngressSpec
---@field status networkingv1.IngressStatus

---@class networkingv1.IngressBackend
---@field resource corev1.TypedLocalObjectReference
---@field service networkingv1.IngressServiceBackend

---@class networkingv1.IngressList
---@field TypeMeta v1.TypeMeta
---@field items networkingv1.Ingress[]
---@field metadata v1.ListMeta

---@class networkingv1.IngressLoadBalancerIngress
---@field hostname string
---@field ip string
---@field ports networkingv1.IngressPortStatus[]

---@class networkingv1.IngressLoadBalancerStatus
---@field ingress networkingv1.IngressLoadBalancerIngress[]

---@class networkingv1.IngressPortStatus
---@field error string
---@field port number
---@field protocol string

---@class networkingv1.IngressRule
---@field IngressRuleValue networkingv1.IngressRuleValue
---@field host string

---@class networkingv1.IngressRuleValue
---@field http networkingv1.HTTPIngressRuleValue

---@class networkingv1.IngressServiceBackend
---@field name string
---@field port networkingv1.ServiceBackendPort

---@class networkingv1.IngressSpec
---@field defaultBackend networkingv1.IngressBackend
---@field ingressClassName string
---@field rules networkingv1.IngressRule[]
---@field tls networkingv1.IngressTLS[]

---@class networkingv1.IngressStatus
---@field loadBalancer networkingv1.IngressLoadBalancerStatus

---@class networkingv1.IngressTLS
---@field hosts string[]
---@field secretName string

---@class networkingv1.NetworkPolicy
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field spec networkingv1.NetworkPolicySpec

---@class networkingv1.NetworkPolicyEgressRule
---@field ports networkingv1.NetworkPolicyPort[]
---@field to networkingv1.NetworkPolicyPeer[]

---@class networkingv1.NetworkPolicyIngressRule
---@field from networkingv1.NetworkPolicyPeer[]
---@field ports networkingv1.NetworkPolicyPort[]

---@class networkingv1.NetworkPolicyList
---@field TypeMeta v1.TypeMeta
---@field items networkingv1.NetworkPolicy[]
---@field metadata v1.ListMeta

---@class networkingv1.NetworkPolicyPeer
---@field ipBlock networkingv1.IPBlock
---@field namespaceSelector v1.LabelSelector
---@field podSelector v1.LabelSelector

---@class networkingv1.NetworkPolicyPort
---@field endPort number
---@field port intstr.IntOrString
---@field protocol string

---@class networkingv1.NetworkPolicySpec
---@field egress networkingv1.NetworkPolicyEgressRule[]
---@field ingress networkingv1.NetworkPolicyIngressRule[]
---@field podSelector v1.LabelSelector
---@field policyTypes string[]

---@class networkingv1.ServiceBackendPort
---@field name string
---@field number number

---@class rbacv1.AggregationRule
---@field clusterRoleSelectors v1.LabelSelector[]

---@class rbacv1.ClusterRole
---@field TypeMeta v1.TypeMeta
---@field aggregationRule rbacv1.AggregationRule
---@field metadata v1.ObjectMeta
---@field rules rbacv1.PolicyRule[]

---@class rbacv1.ClusterRoleBinding
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field roleRef rbacv1.RoleRef
---@field subjects rbacv1.Subject[]

---@class rbacv1.ClusterRoleBindingList
---@field TypeMeta v1.TypeMeta
---@field items rbacv1.ClusterRoleBinding[]
---@field metadata v1.ListMeta

---@class rbacv1.ClusterRoleList
---@field TypeMeta v1.TypeMeta
---@field items rbacv1.ClusterRole[]
---@field metadata v1.ListMeta

---@class rbacv1.PolicyRule
---@field apiGroups string[]
---@field nonResourceURLs string[]
---@field resourceNames string[]
---@field resources string[]
---@field verbs string[]

---@class rbacv1.Role
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field rules rbacv1.PolicyRule[]

---@class rbacv1.RoleBinding
---@field TypeMeta v1.TypeMeta
---@field metadata v1.ObjectMeta
---@field roleRef rbacv1.RoleRef
---@field subjects rbacv1.Subject[]

---@class rbacv1.RoleBindingList
---@field TypeMeta v1.TypeMeta
---@field items rbacv1.RoleBinding[]
---@field metadata v1.ListMeta

---@class rbacv1.RoleList
---@field TypeMeta v1.TypeMeta
---@field items rbacv1.Role[]
---@field metadata v1.ListMeta

---@class rbacv1.RoleRef
---@field apiGroup string
---@field kind string
---@field name string

---@class rbacv1.Subject
---@field apiGroup string
---@field kind string
---@field name string
---@field namespace string

---@class v1.Condition
---@field lastTransitionTime v1.Time
---@field message string
---@field observedGeneration number
---@field reason string
---@field status string
---@field type string

---@class v1.LabelSelector
---@field matchExpressions v1.LabelSelectorRequirement[]
---@field matchLabels table<string, string>

---@class v1.LabelSelectorRequirement
---@field key string
---@field operator string
---@field values string[]

---@class v1.ListMeta
---@field continue string
---@field remainingItemCount number
---@field resourceVersion string
---@field selfLink string

---@class v1.ManagedFieldsEntry
---@field apiVersion string
---@field fieldsType string
---@field fieldsV1 v1.FieldsV1
---@field manager string
---@field operation string
---@field subresource string
---@field time v1.Time

---@class v1.ObjectMeta
---@field annotations table<string, string>
---@field creationTimestamp v1.Time
---@field deletionGracePeriodSeconds number
---@field deletionTimestamp v1.Time
---@field finalizers string[]
---@field generateName string
---@field generation number
---@field labels table<string, string>
---@field managedFields v1.ManagedFieldsEntry[]
---@field name string
---@field namespace string
---@field ownerReferences v1.OwnerReference[]
---@field resourceVersion string
---@field selfLink string
---@field uid string

---@class v1.OwnerReference
---@field apiVersion string
---@field blockOwnerDeletion boolean
---@field controller boolean
---@field kind string
---@field name string
---@field uid string

---@class v1.Status
---@field TypeMeta v1.TypeMeta
---@field code number
---@field details v1.StatusDetails
---@field message string
---@field metadata v1.ListMeta
---@field reason string
---@field status string

---@class v1.StatusCause
---@field field string
---@field message string
---@field reason string

---@class v1.StatusDetails
---@field causes v1.StatusCause[]
---@field group string
---@field kind string
---@field name string
---@field retryAfterSeconds number
---@field uid string

---@class v1.TypeMeta
---@field apiVersion string
---@field kind string

---@meta kubernetes

---@class kubernetes
local kubernetes = {}

---@param quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
---@return number The memory value in bytes, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_memory(quantity) end

---@param quantity string The CPU quantity to parse (e.g., "100m", "1", "2000m")
---@return number The CPU value in millicores, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_cpu(quantity) end

---@param timestr string The time string in RFC3339 format (e.g., "2025-10-03T16:39:00Z")
---@return number The Unix timestamp, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_time(timestr) end

---@param timestamp number The Unix timestamp to convert
---@return string The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_time(timestamp) end

---@param obj table The Kubernetes object (must have a metadata field)
---@return table The same object with initialized defaults (modified in-place)
function kubernetes.init_defaults(obj) end

---@param duration string The duration string to parse (e.g., "5s", "10m", "2h")
---@return number The duration value in seconds, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_duration(duration) end

---@param seconds number The duration in seconds to convert
---@return string The duration string (e.g., "5m0s", "1h30m0s"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_duration(seconds) end

---@param obj table The Kubernetes object to check
---@param matcher kubernetes.GVKMatcher The GVK matcher with group, version, and kind fields
---@return boolean true if the GVK matches
function kubernetes.match_gvk(obj, matcher) end

---@param obj table The Kubernetes object
---@return table The same object with initialized metadata (modified in-place)
function kubernetes.ensure_metadata(obj) end

---@param obj table The Kubernetes object
---@param key string The label key
---@param value string The label value
---@return table The modified object (for chaining)
function kubernetes.add_label(obj, key, value) end

---@param obj table The Kubernetes object
---@param labels table A table of key-value pairs to add as labels
---@return table The modified object (for chaining)
function kubernetes.add_labels(obj, labels) end

---@param obj table The Kubernetes object
---@param key string The label key to remove
---@return table The modified object (for chaining)
function kubernetes.remove_label(obj, key) end

---@param obj table The Kubernetes object
---@param key string The label key to check
---@return boolean true if the label exists
function kubernetes.has_label(obj, key) end

---@param obj table The Kubernetes object
---@param key string The label key
---@return string|nil The label value, or nil if not found
function kubernetes.get_label(obj, key) end

---@param obj table The Kubernetes object
---@param key string The annotation key
---@param value string The annotation value
---@return table The modified object (for chaining)
function kubernetes.add_annotation(obj, key, value) end

---@param obj table The Kubernetes object
---@param annotations table A table of key-value pairs to add as annotations
---@return table The modified object (for chaining)
function kubernetes.add_annotations(obj, annotations) end

---@param obj table The Kubernetes object
---@param key string The annotation key to remove
---@return table The modified object (for chaining)
function kubernetes.remove_annotation(obj, key) end

---@param obj table The Kubernetes object
---@param key string The annotation key to check
---@return boolean true if the annotation exists
function kubernetes.has_annotation(obj, key) end

---@param obj table The Kubernetes object
---@param key string The annotation key
---@return string|nil The annotation value, or nil if not found
function kubernetes.get_annotation(obj, key) end

return kubernetes
