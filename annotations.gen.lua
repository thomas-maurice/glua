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

---@class corev1.AzureDiskVolumeSource
---@field cachingMode string
---@field diskName string
---@field diskURI string
---@field fsType string
---@field kind string
---@field readOnly boolean

---@class corev1.AzureFileVolumeSource
---@field readOnly boolean
---@field secretName string
---@field shareName string

---@class corev1.CSIVolumeSource
---@field driver string
---@field fsType string
---@field nodePublishSecretRef corev1.LocalObjectReference
---@field readOnly boolean
---@field volumeAttributes table<string, string>

---@class corev1.Capabilities
---@field add string[]
---@field drop string[]

---@class corev1.CephFSVolumeSource
---@field monitors string[]
---@field path string
---@field readOnly boolean
---@field secretFile string
---@field secretRef corev1.LocalObjectReference
---@field user string

---@class corev1.CinderVolumeSource
---@field fsType string
---@field readOnly boolean
---@field secretRef corev1.LocalObjectReference
---@field volumeID string

---@class corev1.ClusterTrustBundleProjection
---@field labelSelector v1.LabelSelector
---@field name string
---@field optional boolean
---@field path string
---@field signerName string

---@class corev1.ConfigMapEnvSource
---@field LocalObjectReference corev1.LocalObjectReference
---@field optional boolean

---@class corev1.ConfigMapKeySelector
---@field LocalObjectReference corev1.LocalObjectReference
---@field key string
---@field optional boolean

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

---@class corev1.LocalObjectReference
---@field name string

---@class corev1.NFSVolumeSource
---@field path string
---@field readOnly boolean
---@field server string

---@class corev1.NodeAffinity
---@field preferredDuringSchedulingIgnoredDuringExecution corev1.PreferredSchedulingTerm[]
---@field requiredDuringSchedulingIgnoredDuringExecution corev1.NodeSelector

---@class corev1.NodeSelector
---@field nodeSelectorTerms corev1.NodeSelectorTerm[]

---@class corev1.NodeSelectorRequirement
---@field key string
---@field operator string
---@field values string[]

---@class corev1.NodeSelectorTerm
---@field matchExpressions corev1.NodeSelectorRequirement[]
---@field matchFields corev1.NodeSelectorRequirement[]

---@class corev1.ObjectFieldSelector
---@field apiVersion string
---@field fieldPath string

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

---@class corev1.PersistentVolumeClaimTemplate
---@field metadata v1.ObjectMeta
---@field spec corev1.PersistentVolumeClaimSpec

---@class corev1.PersistentVolumeClaimVolumeSource
---@field claimName string
---@field readOnly boolean

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

---@class corev1.SecretEnvSource
---@field LocalObjectReference corev1.LocalObjectReference
---@field optional boolean

---@class corev1.SecretKeySelector
---@field LocalObjectReference corev1.LocalObjectReference
---@field key string
---@field optional boolean

---@class corev1.SecretProjection
---@field LocalObjectReference corev1.LocalObjectReference
---@field items corev1.KeyToPath[]
---@field optional boolean

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

---@class corev1.ServiceAccountTokenProjection
---@field audience string
---@field expirationSeconds number
---@field path string

---@class corev1.SleepAction
---@field seconds number

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

---@class resource.Quantity

---@class v1.FieldsV1

---@class v1.LabelSelector
---@field matchExpressions v1.LabelSelectorRequirement[]
---@field matchLabels table<string, string>

---@class v1.LabelSelectorRequirement
---@field key string
---@field operator string
---@field values string[]

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

---@class v1.Time

---@class v1.TypeMeta
---@field apiVersion string
---@field kind string

---@class intstr.IntOrString

-- Export all types for convenience
local types = {
  corev1.AWSElasticBlockStoreVolumeSource = {},
  corev1.Affinity = {},
  corev1.AppArmorProfile = {},
  corev1.AzureDiskVolumeSource = {},
  corev1.AzureFileVolumeSource = {},
  corev1.CSIVolumeSource = {},
  corev1.Capabilities = {},
  corev1.CephFSVolumeSource = {},
  corev1.CinderVolumeSource = {},
  corev1.ClusterTrustBundleProjection = {},
  corev1.ConfigMapEnvSource = {},
  corev1.ConfigMapKeySelector = {},
  corev1.ConfigMapProjection = {},
  corev1.ConfigMapVolumeSource = {},
  corev1.Container = {},
  corev1.ContainerExtendedResourceRequest = {},
  corev1.ContainerPort = {},
  corev1.ContainerResizePolicy = {},
  corev1.ContainerRestartRule = {},
  corev1.ContainerRestartRuleOnExitCodes = {},
  corev1.ContainerState = {},
  corev1.ContainerStateRunning = {},
  corev1.ContainerStateTerminated = {},
  corev1.ContainerStateWaiting = {},
  corev1.ContainerStatus = {},
  corev1.ContainerUser = {},
  corev1.DownwardAPIProjection = {},
  corev1.DownwardAPIVolumeFile = {},
  corev1.DownwardAPIVolumeSource = {},
  corev1.EmptyDirVolumeSource = {},
  corev1.EnvFromSource = {},
  corev1.EnvVar = {},
  corev1.EnvVarSource = {},
  corev1.EphemeralContainer = {},
  corev1.EphemeralContainerCommon = {},
  corev1.EphemeralVolumeSource = {},
  corev1.ExecAction = {},
  corev1.FCVolumeSource = {},
  corev1.FileKeySelector = {},
  corev1.FlexVolumeSource = {},
  corev1.FlockerVolumeSource = {},
  corev1.GCEPersistentDiskVolumeSource = {},
  corev1.GRPCAction = {},
  corev1.GitRepoVolumeSource = {},
  corev1.GlusterfsVolumeSource = {},
  corev1.HTTPGetAction = {},
  corev1.HTTPHeader = {},
  corev1.HostAlias = {},
  corev1.HostIP = {},
  corev1.HostPathVolumeSource = {},
  corev1.ISCSIVolumeSource = {},
  corev1.ImageVolumeSource = {},
  corev1.KeyToPath = {},
  corev1.Lifecycle = {},
  corev1.LifecycleHandler = {},
  corev1.LinuxContainerUser = {},
  corev1.LocalObjectReference = {},
  corev1.NFSVolumeSource = {},
  corev1.NodeAffinity = {},
  corev1.NodeSelector = {},
  corev1.NodeSelectorRequirement = {},
  corev1.NodeSelectorTerm = {},
  corev1.ObjectFieldSelector = {},
  corev1.PersistentVolumeClaimSpec = {},
  corev1.PersistentVolumeClaimTemplate = {},
  corev1.PersistentVolumeClaimVolumeSource = {},
  corev1.PhotonPersistentDiskVolumeSource = {},
  corev1.Pod = {},
  corev1.PodAffinity = {},
  corev1.PodAffinityTerm = {},
  corev1.PodAntiAffinity = {},
  corev1.PodCertificateProjection = {},
  corev1.PodCondition = {},
  corev1.PodDNSConfig = {},
  corev1.PodDNSConfigOption = {},
  corev1.PodExtendedResourceClaimStatus = {},
  corev1.PodIP = {},
  corev1.PodOS = {},
  corev1.PodReadinessGate = {},
  corev1.PodResourceClaim = {},
  corev1.PodResourceClaimStatus = {},
  corev1.PodSchedulingGate = {},
  corev1.PodSecurityContext = {},
  corev1.PodSpec = {},
  corev1.PodStatus = {},
  corev1.PortworxVolumeSource = {},
  corev1.PreferredSchedulingTerm = {},
  corev1.Probe = {},
  corev1.ProbeHandler = {},
  corev1.ProjectedVolumeSource = {},
  corev1.QuobyteVolumeSource = {},
  corev1.RBDVolumeSource = {},
  corev1.ResourceClaim = {},
  corev1.ResourceFieldSelector = {},
  corev1.ResourceHealth = {},
  corev1.ResourceRequirements = {},
  corev1.ResourceStatus = {},
  corev1.SELinuxOptions = {},
  corev1.ScaleIOVolumeSource = {},
  corev1.SeccompProfile = {},
  corev1.SecretEnvSource = {},
  corev1.SecretKeySelector = {},
  corev1.SecretProjection = {},
  corev1.SecretVolumeSource = {},
  corev1.SecurityContext = {},
  corev1.ServiceAccountTokenProjection = {},
  corev1.SleepAction = {},
  corev1.StorageOSVolumeSource = {},
  corev1.Sysctl = {},
  corev1.TCPSocketAction = {},
  corev1.Toleration = {},
  corev1.TopologySpreadConstraint = {},
  corev1.TypedLocalObjectReference = {},
  corev1.TypedObjectReference = {},
  corev1.Volume = {},
  corev1.VolumeDevice = {},
  corev1.VolumeMount = {},
  corev1.VolumeMountStatus = {},
  corev1.VolumeProjection = {},
  corev1.VolumeResourceRequirements = {},
  corev1.VolumeSource = {},
  corev1.VsphereVirtualDiskVolumeSource = {},
  corev1.WeightedPodAffinityTerm = {},
  corev1.WindowsSecurityContextOptions = {},
  resource.Quantity = {},
  v1.FieldsV1 = {},
  v1.LabelSelector = {},
  v1.LabelSelectorRequirement = {},
  v1.ManagedFieldsEntry = {},
  v1.ObjectMeta = {},
  v1.OwnerReference = {},
  v1.Time = {},
  v1.TypeMeta = {},
  intstr.IntOrString = {},
}

return types
