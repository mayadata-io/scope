package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/common/mtime"
	"github.com/weaveworks/scope/probe"
	"github.com/weaveworks/scope/probe/controls"
	"github.com/weaveworks/scope/probe/docker"
	"github.com/weaveworks/scope/report"
)

// These constants are keys used in node metadata
const (
	IP                           = report.KubernetesIP
	ObservedGeneration           = report.KubernetesObservedGeneration
	Replicas                     = report.KubernetesReplicas
	DesiredReplicas              = report.KubernetesDesiredReplicas
	NodeType                     = report.KubernetesNodeType
	Type                         = report.KubernetesType
	Ports                        = report.KubernetesPorts
	VolumeClaim                  = report.KubernetesVolumeClaim
	Storage                      = report.KubernetesStorage
	StorageClassName             = report.KubernetesStorageClassName
	AccessModes                  = report.KubernetesAccessModes
	ReclaimPolicy                = report.KubernetesReclaimPolicy
	Status                       = report.KubernetesStatus
	Message                      = report.KubernetesMessage
	VolumeName                   = report.KubernetesVolumeName
	Provisioner                  = report.KubernetesProvisioner
	StorageDriver                = report.KubernetesStorageDriver
	VolumeSnapshotName           = report.KubernetesVolumeSnapshotName
	VolumeSnapshotNamespace      = report.KubernetesVolumeSnapshotNamespace
	SnapshotData                 = report.KubernetesSnapshotData
	SnapshotClass                = report.KubernetesSnapshotClass
	VolumeCapacity               = report.KubernetesVolumeCapacity
	Model                        = report.KubernetesModel
	LogicalSectorSize            = report.KubernetesLogicalSectorSize
	FirmwareRevision             = report.KubernetesFirmwareRevision
	Serial                       = report.KubernetesSerial
	Vendor                       = report.KubernetesVendor
	DiskList                     = report.KubernetesDiskList
	MaxPools                     = report.KubernetesMaxPools
	APIVersion                   = report.KubernetesAPIVersion
	Value                        = report.KubernetesValue
	StoragePoolClaimName         = report.KubernetesStoragePoolClaimName
	DiskName                     = report.KubernetesDiskName
	PoolName                     = report.KubernetesPoolName
	PoolClaim                    = report.KubernetesPoolClaim
	HostName                     = report.KubernetesHostName
	VolumePod                    = report.KubernetesVolumePod
	CStorVolumeName              = report.KubernetesCStorVolumeName
	CStorVolumeReplicaName       = report.KubernetesCStorVolumeReplicaName
	CStorPoolName                = report.KubernetesCStorPoolName
	CStorPoolUID                 = report.KubernetesCStorPoolUID
	CStorVolumeConsistencyFactor = report.KubernetesCStorVolumeConsistencyFactor
	CStorVolumeReplicationFactor = report.KubernetesCStorVolumeReplicationFactor
	CStorVolumeIQN               = report.KubernetesCStorVolumeIQN
	PhysicalSectorSize           = report.KubernetesPhysicalSectorSize
	RotationRate                 = report.KubernetesRotationRate
	CurrentTemperature           = report.KubernetesCurrentTemperature
	HighestTemperature           = report.KubernetesHighestTemperature
	LowestTemperature            = report.KubernetesLowestTemperature
	TotalBytesRead               = report.KubernetesTotalBytesRead
	TotalBytesWritten            = report.KubernetesTotalBytesWritten
	DeviceUtilizationRate        = report.KubernetesDeviceUtilizationRate
	PercentEnduranceUsed         = report.KubernetesPercentEnduranceUsed
	BlockDeviceList              = report.KubernetesBlockDeviceList
	Path                         = report.KubernetesPath
	BlockDeviceName              = report.KubernetesBlockDeviceName
	BlockDeviceClaimName         = report.KubernetesBlockDeviceClaimName
	CASType                      = report.KubernetesCASType
	TotalSize                    = report.KubernetesTotalSize
	FreeSize                     = report.KubernetesFreeSize
	UsedSize                     = report.KubernetesUsedSize
	LogicalUsed                  = report.KubernetesLogicalUsed
	ProvisionedInstances         = report.KubernetesProvisionedInstances
	DesiredInstances             = report.KubernetesDesiredInstances
	HealthyInstances             = report.KubernetesHealthyInstances
	ReadOnly                     = report.KubernetesReadOnly
	ProvisionedReplicas          = report.KubernetesProvisionedReplicas
	HealthyReplicas              = report.KubernetesHealthyReplicas
	CStorPoolInstanceUID         = report.KubernetesCStorPoolInstanceUID
	Driver                       = report.KubernetesDriver
	DeletionPolicy               = report.KubernetesDeletionPolicy
)

var (
	// CStorVolumeStatusMap is map of status and node tag
	CStorVolumeStatusMap = map[string]string{
		"degraded":   "degraded",
		"error":      "failed",
		"healthy":    "",
		"init":       "pending",
		"invalid":    "notpermitted",
		"offline":    "offline",
		"online":     "",
		"rebuilding": "reload",
		"inactive":   "offline",
		"active":     "",
	}
)

// Exposed for testing
var (
	PodMetadataTemplates = report.MetadataTemplates{
		State:            {ID: State, Label: "State", From: report.FromLatest, Priority: 2},
		IP:               {ID: IP, Label: "IP", From: report.FromLatest, Datatype: report.IP, Priority: 3},
		report.Container: {ID: report.Container, Label: "# Containers", From: report.FromCounters, Datatype: report.Number, Priority: 4},
		Namespace:        {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 5},
		Created:          {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 6},
		RestartCount:     {ID: RestartCount, Label: "Restart #", From: report.FromLatest, Priority: 7},
	}

	PodMetricTemplates = docker.ContainerMetricTemplates

	ServiceMetadataTemplates = report.MetadataTemplates{
		Namespace:  {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:    {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		PublicIP:   {ID: PublicIP, Label: "Public IP", From: report.FromLatest, Datatype: report.IP, Priority: 4},
		IP:         {ID: IP, Label: "Internal IP", From: report.FromLatest, Datatype: report.IP, Priority: 5},
		report.Pod: {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 6},
		Type:       {ID: Type, Label: "Type", From: report.FromLatest, Priority: 7},
		Ports:      {ID: Ports, Label: "Ports", From: report.FromLatest, Priority: 8},
	}

	ServiceMetricTemplates = PodMetricTemplates

	DeploymentMetadataTemplates = report.MetadataTemplates{
		NodeType:           {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:          {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:            {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		ObservedGeneration: {ID: ObservedGeneration, Label: "Observed gen.", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		DesiredReplicas:    {ID: DesiredReplicas, Label: "Desired replicas", From: report.FromLatest, Datatype: report.Number, Priority: 5},
		report.Pod:         {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 6},
		Strategy:           {ID: Strategy, Label: "Strategy", From: report.FromLatest, Priority: 7},
	}

	DeploymentMetricTemplates = PodMetricTemplates

	DaemonSetMetadataTemplates = report.MetadataTemplates{
		NodeType:        {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:       {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:         {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		DesiredReplicas: {ID: DesiredReplicas, Label: "Desired replicas", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		report.Pod:      {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 5},
	}

	DaemonSetMetricTemplates = PodMetricTemplates

	StatefulSetMetadataTemplates = report.MetadataTemplates{
		NodeType:           {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:          {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:            {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		ObservedGeneration: {ID: ObservedGeneration, Label: "Observed gen.", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		DesiredReplicas:    {ID: DesiredReplicas, Label: "Desired replicas", From: report.FromLatest, Datatype: report.Number, Priority: 5},
		report.Pod:         {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 6},
	}

	StatefulSetMetricTemplates = PodMetricTemplates

	CronJobMetadataTemplates = report.MetadataTemplates{
		NodeType:      {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:     {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:       {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		Schedule:      {ID: Schedule, Label: "Schedule", From: report.FromLatest, Priority: 4},
		LastScheduled: {ID: LastScheduled, Label: "Last scheduled", From: report.FromLatest, Datatype: report.DateTime, Priority: 5},
		Suspended:     {ID: Suspended, Label: "Suspended", From: report.FromLatest, Priority: 6},
		ActiveJobs:    {ID: ActiveJobs, Label: "# Jobs", From: report.FromLatest, Datatype: report.Number, Priority: 7},
		report.Pod:    {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 8},
	}

	CronJobMetricTemplates = PodMetricTemplates

	PersistentVolumeMetadataTemplates = report.MetadataTemplates{
		NodeType:         {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeClaim:      {ID: VolumeClaim, Label: "Volume claim", From: report.FromLatest, Priority: 2},
		StorageClassName: {ID: StorageClassName, Label: "Storage class", From: report.FromLatest, Priority: 3},
		AccessModes:      {ID: AccessModes, Label: "Access modes", From: report.FromLatest, Priority: 5},
		Status:           {ID: Status, Label: "Status", From: report.FromLatest, Priority: 6},
		StorageDriver:    {ID: StorageDriver, Label: "Storage driver", From: report.FromLatest, Priority: 7},
	}

	PersistentVolumeClaimMetadataTemplates = report.MetadataTemplates{
		NodeType:         {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:        {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Status:           {ID: Status, Label: "Status", From: report.FromLatest, Priority: 3},
		VolumeName:       {ID: VolumeName, Label: "Volume", From: report.FromLatest, Priority: 4},
		StorageClassName: {ID: StorageClassName, Label: "Storage class", From: report.FromLatest, Priority: 5},
		VolumeCapacity:   {ID: VolumeCapacity, Label: "Capacity", From: report.FromLatest, Priority: 6},
	}

	StorageClassMetadataTemplates = report.MetadataTemplates{
		NodeType:    {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Provisioner: {ID: Provisioner, Label: "Provisioner", From: report.FromLatest, Priority: 2},
	}

	VolumeSnapshotMetadataTemplates = report.MetadataTemplates{
		NodeType:     {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:    {ID: Namespace, Label: "Name", From: report.FromLatest, Priority: 2},
		VolumeClaim:  {ID: VolumeClaim, Label: "Persistent volume claim", From: report.FromLatest, Priority: 3},
		SnapshotData: {ID: SnapshotData, Label: "Volume snapshot data", From: report.FromLatest, Priority: 4},
	}

	VolumeSnapshotDataMetadataTemplates = report.MetadataTemplates{
		NodeType:           {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeName:         {ID: VolumeName, Label: "Persistent volume", From: report.FromLatest, Priority: 2},
		VolumeSnapshotName: {ID: VolumeSnapshotName, Label: "Volume snapshot", From: report.FromLatest, Priority: 3},
	}

	JobMetadataTemplates = report.MetadataTemplates{
		NodeType:   {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Name:       {ID: Name, Label: "Name", From: report.FromLatest, Priority: 2},
		Namespace:  {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 3},
		Created:    {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 4},
		report.Pod: {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 5},
	}

	JobMetricTemplates = PodMetricTemplates

	DiskMetadataTemplates = report.MetadataTemplates{
		NodeType:              {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Model:                 {ID: Model, Label: "Model", From: report.FromLatest, Priority: 2},
		Serial:                {ID: Serial, Label: "Serial", From: report.FromLatest, Priority: 3},
		Vendor:                {ID: Vendor, Label: "Vendor", From: report.FromLatest, Priority: 4},
		FirmwareRevision:      {ID: FirmwareRevision, Label: "Firmware Revision", From: report.FromLatest, Priority: 5},
		LogicalSectorSize:     {ID: LogicalSectorSize, Label: "Logical Sector Size", From: report.FromLatest, Priority: 6},
		PhysicalSectorSize:    {ID: PhysicalSectorSize, Label: "Physical Sector Size", From: report.FromLatest, Priority: 7},
		VolumeCapacity:        {ID: VolumeCapacity, Label: "Capacity", From: report.FromLatest, Priority: 8},
		Status:                {ID: Status, Label: "Status", From: report.FromLatest, Priority: 9},
		Created:               {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 10},
		RotationRate:          {ID: RotationRate, Label: "Rotation Rate", From: report.FromLatest, Priority: 11},
		CurrentTemperature:    {ID: CurrentTemperature, Label: "Current Temperature", From: report.FromLatest, Priority: 12},
		HighestTemperature:    {ID: HighestTemperature, Label: "Highest Temperature", From: report.FromLatest, Priority: 13},
		LowestTemperature:     {ID: LowestTemperature, Label: "Lowest Temperature", From: report.FromLatest, Priority: 14},
		TotalBytesRead:        {ID: TotalBytesRead, Label: "Total Bytes Read", From: report.FromLatest, Priority: 15},
		TotalBytesWritten:     {ID: TotalBytesWritten, Label: "Total Bytes Written", From: report.FromLatest, Priority: 16},
		DeviceUtilizationRate: {ID: DeviceUtilizationRate, Label: "Device Utilization Rate", From: report.FromLatest, Priority: 17},
		PercentEnduranceUsed:  {ID: PercentEnduranceUsed, Label: "Percent Endurance Used", From: report.FromLatest, Priority: 18},
	}

	BlockDeviceMetadataTemplates = report.MetadataTemplates{
		NodeType:          {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Model:             {ID: Model, Label: "Model", From: report.FromLatest, Priority: 2},
		Serial:            {ID: Serial, Label: "Serial", From: report.FromLatest, Priority: 3},
		Vendor:            {ID: Vendor, Label: "Vendor", From: report.FromLatest, Priority: 4},
		FirmwareRevision:  {ID: FirmwareRevision, Label: "Firmware Revision", From: report.FromLatest, Priority: 5},
		LogicalSectorSize: {ID: LogicalSectorSize, Label: "Logical Sector Size", From: report.FromLatest, Priority: 6},
	}

	StoragePoolClaimMetadataTemplates = report.MetadataTemplates{
		NodeType:   {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		APIVersion: {ID: APIVersion, Label: "API Version", From: report.FromLatest, Priority: 2},
		Status:     {ID: Status, Label: "Status", From: report.FromLatest, Priority: 3},
		MaxPools:   {ID: MaxPools, Label: "MaxPools", From: report.FromLatest, Priority: 4},
	}

	CStorVolumeMetadataTemplates = report.MetadataTemplates{
		NodeType:                     {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeName:                   {ID: CStorVolumeName, Label: "CStor Volume", From: report.FromLatest, Priority: 2},
		Status:                       {ID: Status, Label: "Status", From: report.FromLatest, Priority: 3},
		CStorVolumeConsistencyFactor: {ID: CStorVolumeConsistencyFactor, Label: "Conistency Factor", From: report.FromLatest, Priority: 4},
		CStorVolumeReplicationFactor: {ID: CStorVolumeReplicationFactor, Label: "Replication Factor", From: report.FromLatest, Priority: 5},
		CStorVolumeIQN:               {ID: CStorVolumeIQN, Label: "Iqn", From: report.FromLatest, Priority: 6},
	}
	CStorVolumeReplicaMetadataTemplates = report.MetadataTemplates{
		NodeType:   {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeName: {ID: CStorVolumeReplicaName, Label: "CStor Volume Replica", From: report.FromLatest, Priority: 2},
		Status:     {ID: Status, Label: "Status", From: report.FromLatest, Priority: 3},
	}

	CStorPoolMetadataTemplates = report.MetadataTemplates{
		NodeType:   {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeName: {ID: CStorPoolName, Label: "CStor Pool", From: report.FromLatest, Priority: 2},
		Status:     {ID: Status, Label: "Status", From: report.FromLatest, Priority: 3},
	}

	BlockDeviceClaimMetadataTemplates = report.MetadataTemplates{
		NodeType:        {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		BlockDeviceName: {ID: BlockDeviceName, Label: "Block device name", From: report.FromLatest, Priority: 2},
		HostName:        {ID: HostName, Label: "Host", From: report.FromLatest, Priority: 3},
		Status:          {ID: Status, Label: "Status", From: report.FromLatest, Priority: 4},
	}

	CStorPoolClusterMetadataTemplates = report.MetadataTemplates{
		NodeType:             {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		ProvisionedInstances: {ID: ProvisionedInstances, Label: "Provisioned Instances", From: report.FromLatest, Priority: 2},
		DesiredInstances:     {ID: DesiredInstances, Label: "Desired Instances", From: report.FromLatest, Priority: 3},
		HealthyInstances:     {ID: HealthyInstances, Label: "Healthy Instances", From: report.FromLatest, Priority: 4},
	}

	CStorPoolInstanceMetadataTemplates = report.MetadataTemplates{
		NodeType:             {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Status:               {ID: Status, Label: "Status", From: report.FromLatest, Priority: 2},
		TotalSize:            {ID: TotalSize, Label: "Total size", From: report.FromLatest, Priority: 3},
		FreeSize:             {ID: FreeSize, Label: "Free size", From: report.FromLatest, Priority: 4},
		UsedSize:             {ID: UsedSize, Label: "Used size", From: report.FromLatest, Priority: 5},
		LogicalUsed:          {ID: LogicalUsed, Label: "Logical Used size", From: report.FromLatest, Priority: 6},
		ReadOnly:             {ID: ReadOnly, Label: "Read Only", From: report.FromLatest, Priority: 7},
		ProvisionedReplicas:  {ID: ProvisionedReplicas, Label: "Provisioned Replicas", From: report.FromLatest, Priority: 8},
		HealthyReplicas:      {ID: HealthyReplicas, Label: "Healthy Replicas", From: report.FromLatest, Priority: 9},
	}

	CsiVolumeSnapshotMetadataTemplates = report.MetadataTemplates{
		NodeType:      {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:     {ID: Namespace, Label: "Name", From: report.FromLatest, Priority: 2},
		VolumeClaim:   {ID: VolumeClaim, Label: "Persistent volume claim", From: report.FromLatest, Priority: 3},
		SnapshotData:  {ID: SnapshotData, Label: "Volume snapshot data", From: report.FromLatest, Priority: 4},
		SnapshotClass: {ID: SnapshotClass, Label: "Volume snapshot class", From: report.FromLatest, Priority: 5},
	}

	VolumeSnapshotClassMetadataTemplates = report.MetadataTemplates{
		NodeType:       {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Driver:         {ID: Driver, Label: "Driver", From: report.FromLatest, Priority: 2},
		DeletionPolicy: {ID: DeletionPolicy, Label: "Deletion policy", From: report.FromLatest, Priority: 3},
	}

	VolumeSnapshotContentMetadataTemplates = report.MetadataTemplates{
		NodeType:                {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		VolumeSnapshotName:      {ID: VolumeSnapshotName, Label: "Volume snapshot name", From: report.FromLatest, Priority: 2},
		VolumeSnapshotNamespace: {ID: VolumeSnapshotNamespace, Label: "Volume snapshot namespace", From: report.FromLatest, Priority: 3},
	}

	TableTemplates = report.TableTemplates{
		LabelPrefix: {
			ID:     LabelPrefix,
			Label:  "Kubernetes labels",
			Type:   report.PropertyListType,
			Prefix: LabelPrefix,
		},
	}

	ScalingControls = []report.Control{
		{
			ID:       ScaleDown,
			Human:    "Scale down",
			Category: report.AdminControl,
			Icon:     "fa fa-minus",
			Rank:     0,
		},
		{
			ID:       ScaleUp,
			Human:    "Scale up",
			Category: report.AdminControl,
			Icon:     "fa fa-plus",
			Rank:     1,
		},
	}

	DescribeControl = report.Control{
		ID:       Describe,
		Human:    "Describe",
		Category: report.ReadOnlyControl,
		Icon:     "fa fa-file-text",
		Rank:     2,
	}
)

// Reporter generate Reports containing Container and ContainerImage topologies
type Reporter struct {
	client          Client
	pipes           controls.PipeClient
	probeID         string
	probe           *probe.Probe
	hostID          string
	handlerRegistry *controls.HandlerRegistry
	nodeName        string
	kubeletPort     uint
}

// NewReporter makes a new Reporter
func NewReporter(client Client, pipes controls.PipeClient, probeID string, hostID string, probe *probe.Probe, handlerRegistry *controls.HandlerRegistry, nodeName string, kubeletPort uint) *Reporter {
	reporter := &Reporter{
		client:          client,
		pipes:           pipes,
		probeID:         probeID,
		probe:           probe,
		hostID:          hostID,
		handlerRegistry: handlerRegistry,
		nodeName:        nodeName,
		kubeletPort:     kubeletPort,
	}
	reporter.registerControls()
	client.WatchPods(reporter.podEvent)
	return reporter
}

// Stop unregisters controls.
func (r *Reporter) Stop() {
	r.deregisterControls()
}

// Name of this reporter, for metrics gathering
func (Reporter) Name() string { return "K8s" }

func (r *Reporter) podEvent(e Event, pod Pod) {
	// filter out non-local pods, if we have been given a node name to report on
	if r.nodeName != "" && pod.NodeName() != r.nodeName {
		return
	}
	switch e {
	case ADD:
		rpt := report.MakeReport()
		rpt.Shortcut = true
		rpt.Pod.AddNode(pod.GetNode(r.probeID))
		r.probe.Publish(rpt)
	case DELETE:
		rpt := report.MakeReport()
		rpt.Shortcut = true
		rpt.Pod.AddNode(
			report.MakeNodeWith(
				report.MakePodNodeID(pod.UID()),
				map[string]string{State: report.StateDeleted},
			),
		)
		r.probe.Publish(rpt)
	}
}

func isPauseContainer(n report.Node, rpt report.Report) bool {
	containerImageIDs, ok := n.Parents.Lookup(report.ContainerImage)
	if !ok {
		return false
	}
	for _, imageNodeID := range containerImageIDs {
		imageNode, ok := rpt.ContainerImage.Nodes[imageNodeID]
		if !ok {
			continue
		}
		imageName, ok := imageNode.Latest.Lookup(docker.ImageName)
		if !ok {
			continue
		}
		return report.IsPauseImageName(imageName)
	}
	return false
}

// Tagger adds pod parents to container nodes.
type Tagger struct {
}

// Name of this tagger, for metrics gathering
func (Tagger) Name() string { return "K8s" }

// Tag adds pod parents to container nodes.
func (r *Tagger) Tag(rpt report.Report) (report.Report, error) {
	for id, n := range rpt.Container.Nodes {
		uid, ok := n.Latest.Lookup(docker.LabelPrefix + "io.kubernetes.pod.uid")
		if !ok {
			continue
		}

		// Tag the pause containers with "does-not-make-connections"
		if isPauseContainer(n, rpt) {
			n = n.WithLatest(report.DoesNotMakeConnections, mtime.Now(), "")
		}

		rpt.Container.Nodes[id] = n.WithParent(report.Pod, report.MakePodNodeID(uid))
	}
	return rpt, nil
}

// Report generates a Report containing Container and ContainerImage topologies
func (r *Reporter) Report() (report.Report, error) {
	result := report.MakeReport()
	serviceTopology, services, err := r.serviceTopology()
	if err != nil {
		return result, err
	}
	daemonSetTopology, daemonSets, err := r.daemonSetTopology()
	if err != nil {
		return result, err
	}
	statefulSetTopology, statefulSets, err := r.statefulSetTopology()
	if err != nil {
		return result, err
	}
	cronJobTopology, cronJobs, err := r.cronJobTopology()
	if err != nil {
		return result, err
	}
	deploymentTopology, deployments, err := r.deploymentTopology()
	if err != nil {
		return result, err
	}
	jobTopology, jobs, err := r.jobTopology()
	if err != nil {
		return result, err
	}
	podTopology, err := r.podTopology(services, deployments, daemonSets, statefulSets, cronJobs, jobs)
	if err != nil {
		return result, err
	}
	namespaceTopology, err := r.namespaceTopology()
	if err != nil {
		return result, err
	}
	persistentVolumeTopology, _, err := r.persistentVolumeTopology()
	if err != nil {
		return result, err
	}
	persistentVolumeClaimTopology, _, err := r.persistentVolumeClaimTopology()
	if err != nil {
		return result, err
	}
	storageClassTopology, _, err := r.storageClassTopology()
	if err != nil {
		return result, err
	}
	volumeSnapshotTopology, _, err := r.volumeSnapshotTopology()
	if err != nil {
		return result, err
	}
	volumeSnapshotDataTopology, _, err := r.volumeSnapshotDataTopology()
	if err != nil {
		return result, err
	}
	diskTopology, _, err := r.diskTopology()
	if err != nil {
		return result, err
	}

	blockDeviceTopology, _, err := r.blockDeviceTopology()
	if err != nil {
		return result, err
	}

	storagePoolClaimTopology, _, err := r.storagePoolClaimTopology()
	if err != nil {
		return result, err
	}

	cStorVolumeTopology, _, err := r.cStorVolumeTopology()
	if err != nil {
		return result, err
	}
	cStorVolumeReplicaTopology, _, err := r.cStorVolumeReplicaTopology()
	if err != nil {
		return result, err
	}
	cStorPoolTopology, _, err := r.cStorPoolTopology()
	if err != nil {
		return result, err
	}
	blockDeviceClaimTopology, _, err := r.blockDeviceClaimTopology()
	if err != nil {
		return result, err
	}
	cStorPoolClusterTopology, _, err := r.cStorPoolClusterTopology()
	if err != nil {
		return result, err
	}
	cStorPoolInstanceTopology, _, err := r.cStorPoolInstanceTopology()
	if err != nil {
		return result, err
	}
	csiVolumeSnapshotTopology, _, err := r.csiVolumeSnapshotTopology()
	if err != nil {
		return result, nil
	}
	volumeSnapshotClassTopology, _, err := r.volumeSnapshotClassTopology()
	if err != nil {
		return result, nil
	}
	volumeSnapshotContentTopology, _, err := r.volumeSnapshotContentTopology()
	if err != nil {
		return result, nil
	}
	result.Pod = result.Pod.Merge(podTopology)
	result.Service = result.Service.Merge(serviceTopology)
	result.DaemonSet = result.DaemonSet.Merge(daemonSetTopology)
	result.StatefulSet = result.StatefulSet.Merge(statefulSetTopology)
	result.CronJob = result.CronJob.Merge(cronJobTopology)
	result.Deployment = result.Deployment.Merge(deploymentTopology)
	result.Namespace = result.Namespace.Merge(namespaceTopology)
	result.PersistentVolume = result.PersistentVolume.Merge(persistentVolumeTopology)
	result.PersistentVolumeClaim = result.PersistentVolumeClaim.Merge(persistentVolumeClaimTopology)
	result.StorageClass = result.StorageClass.Merge(storageClassTopology)
	result.VolumeSnapshot = result.VolumeSnapshot.Merge(volumeSnapshotTopology)
	result.VolumeSnapshotData = result.VolumeSnapshotData.Merge(volumeSnapshotDataTopology)
	result.Job = result.Job.Merge(jobTopology)
	result.Disk = result.Disk.Merge(diskTopology)
	result.StoragePoolClaim = result.StoragePoolClaim.Merge(storagePoolClaimTopology)
	result.CStorVolume = result.CStorVolume.Merge(cStorVolumeTopology)
	result.CStorVolumeReplica = result.CStorVolumeReplica.Merge(cStorVolumeReplicaTopology)
	result.CStorPool = result.CStorPool.Merge(cStorPoolTopology)
	result.BlockDevice = result.BlockDevice.Merge(blockDeviceTopology)
	result.BlockDeviceClaim = result.BlockDeviceClaim.Merge(blockDeviceClaimTopology)
	result.CStorPoolCluster = result.CStorPoolCluster.Merge(cStorPoolClusterTopology)
	result.CStorPoolInstance = result.CStorPoolInstance.Merge(cStorPoolInstanceTopology)
	result.CsiVolumeSnapshot = result.CsiVolumeSnapshot.Merge(csiVolumeSnapshotTopology)
	result.VolumeSnapshotClass = result.VolumeSnapshotClass.Merge(volumeSnapshotClassTopology)
	result.VolumeSnapshotContent = result.VolumeSnapshotContent.Merge(volumeSnapshotContentTopology)
	return result, nil
}

func (r *Reporter) serviceTopology() (report.Topology, []Service, error) {
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(ServiceMetadataTemplates).
			WithMetricTemplates(ServiceMetricTemplates).
			WithTableTemplates(TableTemplates)
		services = []Service{}
	)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkServices(func(s Service) error {
		result.AddNode(s.GetNode(r.probeID))
		services = append(services, s)
		return nil
	})
	return result, services, err
}

func (r *Reporter) deploymentTopology() (report.Topology, []Deployment, error) {
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(DeploymentMetadataTemplates).
			WithMetricTemplates(DeploymentMetricTemplates).
			WithTableTemplates(TableTemplates)
		deployments = []Deployment{}
	)
	result.Controls.AddControls(ScalingControls)
	result.Controls.AddControl(DescribeControl)

	err := r.client.WalkDeployments(func(d Deployment) error {
		result.AddNode(d.GetNode(r.probeID))
		deployments = append(deployments, d)
		return nil
	})
	return result, deployments, err
}

func (r *Reporter) daemonSetTopology() (report.Topology, []DaemonSet, error) {
	daemonSets := []DaemonSet{}
	result := report.MakeTopology().
		WithMetadataTemplates(DaemonSetMetadataTemplates).
		WithMetricTemplates(DaemonSetMetricTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkDaemonSets(func(d DaemonSet) error {
		result.AddNode(d.GetNode(r.probeID))
		daemonSets = append(daemonSets, d)
		return nil
	})
	return result, daemonSets, err
}

func (r *Reporter) statefulSetTopology() (report.Topology, []StatefulSet, error) {
	statefulSets := []StatefulSet{}
	result := report.MakeTopology().
		WithMetadataTemplates(StatefulSetMetadataTemplates).
		WithMetricTemplates(StatefulSetMetricTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkStatefulSets(func(s StatefulSet) error {
		result.AddNode(s.GetNode(r.probeID))
		statefulSets = append(statefulSets, s)
		return nil
	})
	return result, statefulSets, err
}

func (r *Reporter) cronJobTopology() (report.Topology, []CronJob, error) {
	cronJobs := []CronJob{}
	result := report.MakeTopology().
		WithMetadataTemplates(CronJobMetadataTemplates).
		WithMetricTemplates(CronJobMetricTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCronJobs(func(c CronJob) error {
		result.AddNode(c.GetNode(r.probeID))
		cronJobs = append(cronJobs, c)
		return nil
	})
	return result, cronJobs, err
}

func (r *Reporter) persistentVolumeTopology() (report.Topology, []PersistentVolume, error) {
	persistentVolumes := []PersistentVolume{}
	result := report.MakeTopology().
		WithMetadataTemplates(PersistentVolumeMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkPersistentVolumes(func(p PersistentVolume) error {
		result.AddNode(p.GetNode(r.probeID))
		persistentVolumes = append(persistentVolumes, p)
		return nil
	})
	return result, persistentVolumes, err
}

func (r *Reporter) persistentVolumeClaimTopology() (report.Topology, []PersistentVolumeClaim, error) {
	persistentVolumeClaims := []PersistentVolumeClaim{}
	result := report.MakeTopology().
		WithMetadataTemplates(PersistentVolumeClaimMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(report.Control{
		ID:       CreateVolumeSnapshot,
		Human:    "Create snapshot",
		Category: report.AdminControl,
		Icon:     "fa fa-camera",
		Rank:     0,
	})
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkPersistentVolumeClaims(func(p PersistentVolumeClaim) error {
		result.AddNode(p.GetNode(r.probeID))
		persistentVolumeClaims = append(persistentVolumeClaims, p)
		return nil
	})
	return result, persistentVolumeClaims, err
}

func (r *Reporter) storageClassTopology() (report.Topology, []StorageClass, error) {
	storageClasses := []StorageClass{}
	result := report.MakeTopology().
		WithMetadataTemplates(StorageClassMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkStorageClasses(func(p StorageClass) error {
		result.AddNode(p.GetNode(r.probeID))
		storageClasses = append(storageClasses, p)
		return nil
	})
	return result, storageClasses, err
}

func (r *Reporter) volumeSnapshotTopology() (report.Topology, []VolumeSnapshot, error) {
	volumeSnapshots := []VolumeSnapshot{}
	result := report.MakeTopology().
		WithMetadataTemplates(VolumeSnapshotMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(report.Control{
		ID:       CloneVolumeSnapshot,
		Human:    "Clone snapshot",
		Category: report.AdminControl,
		Icon:     "far fa-clone",
		Rank:     0,
	})
	result.Controls.AddControl(report.Control{
		ID:       DeleteVolumeSnapshot,
		Human:    "Delete",
		Category: report.AdminControl,
		Icon:     "far fa-trash-alt",
		Rank:     3,
	})
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkVolumeSnapshots(func(p VolumeSnapshot) error {
		result.AddNode(p.GetNode(r.probeID))
		volumeSnapshots = append(volumeSnapshots, p)
		return nil
	})
	return result, volumeSnapshots, err
}

func (r *Reporter) volumeSnapshotDataTopology() (report.Topology, []VolumeSnapshotData, error) {
	volumeSnapshotData := []VolumeSnapshotData{}
	result := report.MakeTopology().
		WithMetadataTemplates(VolumeSnapshotDataMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkVolumeSnapshotData(func(p VolumeSnapshotData) error {
		result.AddNode(p.GetNode(r.probeID))
		volumeSnapshotData = append(volumeSnapshotData, p)
		return nil
	})
	return result, volumeSnapshotData, err
}

func (r *Reporter) csiVolumeSnapshotTopology() (report.Topology, []CsiVolumeSnapshot, error) {
	volumeSnapshots := []CsiVolumeSnapshot{}
	result := report.MakeTopology().
		WithMetadataTemplates(CsiVolumeSnapshotMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(report.Control{
		ID:       CloneCsiVolumeSnapshot,
		Human:    "Clone snapshot",
		Category: report.AdminControl,
		Icon:     "far fa-clone",
		Rank:     0,
	})
	result.Controls.AddControl(report.Control{
		ID:       DeleteCsiVolumeSnapshot,
		Human:    "Delete",
		Category: report.AdminControl,
		Icon:     "far fa-trash-alt",
		Rank:     3,
	})
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCsiVolumeSnapshots(func(p CsiVolumeSnapshot) error {
		result.AddNode(p.GetNode(r.probeID))
		volumeSnapshots = append(volumeSnapshots, p)
		return nil
	})
	return result, volumeSnapshots, err
}

func (r *Reporter) volumeSnapshotClassTopology() (report.Topology, []VolumeSnapshotClass, error) {
	vsc := []VolumeSnapshotClass{}
	result := report.MakeTopology().
		WithMetadataTemplates(VolumeSnapshotClassMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkVolumeSnapshotClasses(func(p VolumeSnapshotClass) error {
		result.AddNode(p.GetNode(r.probeID))
		vsc = append(vsc, p)
		return nil
	})
	return result, vsc, err
}

func (r *Reporter) volumeSnapshotContentTopology() (report.Topology, []VolumeSnapshotContent, error) {
	vsc := []VolumeSnapshotContent{}
	result := report.MakeTopology().
		WithMetadataTemplates(VolumeSnapshotContentMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkVolumeSnapshotContents(func(p VolumeSnapshotContent) error {
		result.AddNode(p.GetNode(r.probeID))
		vsc = append(vsc, p)
		return nil
	})
	return result, vsc, err
}

func (r *Reporter) jobTopology() (report.Topology, []Job, error) {
	jobs := []Job{}
	result := report.MakeTopology().
		WithMetadataTemplates(JobMetadataTemplates).
		WithMetricTemplates(JobMetricTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkJobs(func(c Job) error {
		result.AddNode(c.GetNode(r.probeID))
		jobs = append(jobs, c)
		return nil
	})
	return result, jobs, err
}

func (r *Reporter) diskTopology() (report.Topology, []Disk, error) {
	disks := []Disk{}
	result := report.MakeTopology().
		WithMetadataTemplates(DiskMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkDisks(func(p Disk) error {
		result.AddNode(p.GetNode(r.probeID))
		disks = append(disks, p)
		return nil
	})
	return result, disks, err
}

func (r *Reporter) blockDeviceTopology() (report.Topology, []BlockDevice, error) {
	blockDevices := []BlockDevice{}
	result := report.MakeTopology().
		WithMetadataTemplates(BlockDeviceMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkBlockDevices(func(p BlockDevice) error {
		result.AddNode(p.GetNode(r.probeID))
		blockDevices = append(blockDevices, p)
		return nil
	})
	return result, blockDevices, err
}

func (r *Reporter) storagePoolClaimTopology() (report.Topology, []StoragePoolClaim, error) {
	storagePoolClaims := []StoragePoolClaim{}
	result := report.MakeTopology().
		WithMetadataTemplates(StoragePoolClaimMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkStoragePoolClaims(func(p StoragePoolClaim) error {
		result.AddNode(p.GetNode(r.probeID))
		storagePoolClaims = append(storagePoolClaims, p)
		return nil
	})
	return result, storagePoolClaims, err
}

func (r *Reporter) cStorVolumeTopology() (report.Topology, []CStorVolume, error) {
	cStorVolumes := []CStorVolume{}
	result := report.MakeTopology().
		WithMetadataTemplates(CStorVolumeMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCStorVolumes(func(p CStorVolume) error {
		result.AddNode(p.GetNode(r.probeID))
		cStorVolumes = append(cStorVolumes, p)
		return nil
	})
	return result, cStorVolumes, err
}

func (r *Reporter) cStorVolumeReplicaTopology() (report.Topology, []CStorVolumeReplica, error) {
	cStorVolumeReplicas := []CStorVolumeReplica{}
	result := report.MakeTopology().
		WithMetadataTemplates(CStorVolumeReplicaMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCStorVolumeReplicas(func(p CStorVolumeReplica) error {
		result.AddNode(p.GetNode(r.probeID))
		cStorVolumeReplicas = append(cStorVolumeReplicas, p)
		return nil
	})
	return result, cStorVolumeReplicas, err
}

func (r *Reporter) cStorPoolTopology() (report.Topology, []CStorPool, error) {
	cStorPool := []CStorPool{}
	result := report.MakeTopology().
		WithMetadataTemplates(CStorPoolMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCStorPools(func(p CStorPool) error {
		result.AddNode(p.GetNode(r.probeID))
		cStorPool = append(cStorPool, p)
		return nil
	})
	return result, cStorPool, err
}

func (r *Reporter) blockDeviceClaimTopology() (report.Topology, []BlockDeviceClaim, error) {
	blockDeviceClaims := []BlockDeviceClaim{}
	result := report.MakeTopology().
		WithMetadataTemplates(BlockDeviceClaimMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkBlockDeviceClaims(func(p BlockDeviceClaim) error {
		result.AddNode(p.GetNode(r.probeID))
		blockDeviceClaims = append(blockDeviceClaims, p)
		return nil
	})
	return result, blockDeviceClaims, err
}

func (r *Reporter) cStorPoolClusterTopology() (report.Topology, []CStorPoolCluster, error) {
	cStorPoolCluster := []CStorPoolCluster{}
	result := report.MakeTopology().
		WithMetadataTemplates(CStorPoolClusterMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCStorPoolClusters(func(p CStorPoolCluster) error {
		result.AddNode(p.GetNode(r.probeID))
		cStorPoolCluster = append(cStorPoolCluster, p)
		return nil
	})
	return result, cStorPoolCluster, err
}

func (r *Reporter) cStorPoolInstanceTopology() (report.Topology, []CStorPoolInstance, error) {
	cStorPoolInstance := []CStorPoolInstance{}
	result := report.MakeTopology().
		WithMetadataTemplates(CStorPoolInstanceMetadataTemplates).
		WithTableTemplates(TableTemplates)
	result.Controls.AddControl(DescribeControl)
	err := r.client.WalkCStorPoolInstances(func(p CStorPoolInstance) error {
		result.AddNode(p.GetNode(r.probeID))
		cStorPoolInstance = append(cStorPoolInstance, p)
		return nil
	})
	return result, cStorPoolInstance, err
}

type labelledChild interface {
	Labels() map[string]string
	AddParent(string, string)
	Namespace() string
}

// Match parses the selectors and adds the target as a parent if the selector matches.
func match(namespace string, selector labels.Selector, topology, id string) func(labelledChild) {
	return func(c labelledChild) {
		if namespace == c.Namespace() && selector.Matches(labels.Set(c.Labels())) {
			c.AddParent(topology, id)
		}
	}
}

func (r *Reporter) podTopology(services []Service, deployments []Deployment, daemonSets []DaemonSet, statefulSets []StatefulSet, cronJobs []CronJob, jobs []Job) (report.Topology, error) {
	var (
		pods = report.MakeTopology().
			WithMetadataTemplates(PodMetadataTemplates).
			WithMetricTemplates(PodMetricTemplates).
			WithTableTemplates(TableTemplates)
		selectors = []func(labelledChild){}
	)
	pods.Controls.AddControl(report.Control{
		ID:       GetLogs,
		Human:    "Get logs",
		Category: report.ReadOnlyControl,
		Icon:     "fa fa-desktop",
		Rank:     0,
	})
	pods.Controls.AddControl(report.Control{
		ID:           DeletePod,
		Human:        "Delete",
		Category:     report.AdminControl,
		Icon:         "far fa-trash-alt",
		Confirmation: "Are you sure you want to delete this pod?",
		Rank:         3,
	})
	pods.Controls.AddControl(DescribeControl)
	for _, service := range services {
		selectors = append(selectors, match(
			service.Namespace(),
			service.Selector(),
			report.Service,
			report.MakeServiceNodeID(service.UID()),
		))
	}
	for _, deployment := range deployments {
		selector, err := deployment.Selector()
		if err != nil {
			return pods, err
		}
		selectors = append(selectors, match(
			deployment.Namespace(),
			selector,
			report.Deployment,
			report.MakeDeploymentNodeID(deployment.UID()),
		))
	}
	for _, daemonSet := range daemonSets {
		selector, err := daemonSet.Selector()
		if err != nil {
			return pods, err
		}
		selectors = append(selectors, match(
			daemonSet.Namespace(),
			selector,
			report.DaemonSet,
			report.MakeDaemonSetNodeID(daemonSet.UID()),
		))
	}
	for _, statefulSet := range statefulSets {
		selector, err := statefulSet.Selector()
		if err != nil {
			return pods, err
		}
		selectors = append(selectors, match(
			statefulSet.Namespace(),
			selector,
			report.StatefulSet,
			report.MakeStatefulSetNodeID(statefulSet.UID()),
		))
	}
	for _, cronJob := range cronJobs {
		cronJobSelectors, err := cronJob.Selectors()
		if err != nil {
			return pods, err
		}
		for _, selector := range cronJobSelectors {
			selectors = append(selectors, match(
				cronJob.Namespace(),
				selector,
				report.CronJob,
				report.MakeCronJobNodeID(cronJob.UID()),
			))
		}
		for _, job := range jobs {
			selector, err := job.Selector()
			if err != nil {
				return pods, err
			}
			selectors = append(selectors, match(
				job.Namespace(),
				selector,
				report.Job,
				report.MakeJobNodeID(job.UID()),
			))
		}
	}

	var localPodUIDs map[string]struct{}
	if r.nodeName == "" && r.kubeletPort != 0 {
		// We don't know the node name: fall back to obtaining the local pods from kubelet
		var err error
		localPodUIDs, err = GetLocalPodUIDs(fmt.Sprintf("127.0.0.1:%d", r.kubeletPort))
		if err != nil {
			log.Debugf("No node name and cannot obtain local pods, reporting all (which may impact performance): %v", err)
		}
	}
	err := r.client.WalkPods(func(p Pod) error {
		// filter out non-local pods: we only want to report local ones for performance reasons.
		if r.nodeName != "" {
			if p.NodeName() != r.nodeName {
				return nil
			}
		} else if localPodUIDs != nil {
			if _, ok := localPodUIDs[p.UID()]; !ok {
				return nil
			}
		}
		for _, selector := range selectors {
			selector(p)
		}
		pods.AddNode(p.GetNode(r.probeID))
		return nil
	})
	return pods, err
}

func (r *Reporter) namespaceTopology() (report.Topology, error) {
	result := report.MakeTopology()
	err := r.client.WalkNamespaces(func(ns NamespaceResource) error {
		result.AddNode(ns.GetNode())
		return nil
	})
	return result, err
}