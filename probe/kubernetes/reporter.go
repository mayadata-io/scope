package kubernetes

import (
	"fmt"
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/labels"

	log "github.com/Sirupsen/logrus"
	"github.com/weaveworks/common/mtime"
	"github.com/weaveworks/scope/probe"
	"github.com/weaveworks/scope/probe/controls"
	"github.com/weaveworks/scope/probe/docker"
	"github.com/weaveworks/scope/probe/host"
	"github.com/weaveworks/scope/report"
)

// These constants are keys used in node metadata
const (
	IP                 = report.KubernetesIP
	ObservedGeneration = report.KubernetesObservedGeneration
	Replicas           = report.KubernetesReplicas
	DesiredReplicas    = report.KubernetesDesiredReplicas
	NodeType           = report.KubernetesNodeType
	Type               = report.KubernetesType
	Ports              = report.KubernetesPorts
	Claim              = report.KubernetesClaim
	StorageClassName   = report.KubernetesStorageClassName
	AccessModes        = report.KubernetesAccessModes
	ReclaimPolicy      = report.KubernetesReclaimPolicy
	Status             = report.KubernetesStatus
	Message            = report.KubernetesMessage
	VolumeName         = report.KubernetesVolumeName
	Provisioner        = report.KubernetesProvisioner
	APIVersion         = report.KubernetesAPIVersion
	UID                = report.KubernetesUID
	ResourceVersion    = report.KubernetesResourceVersion
	SelfLink           = report.KubernetesSelfLink
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
		ObservedGeneration: {ID: ObservedGeneration, Label: "Observed Gen.", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		DesiredReplicas:    {ID: DesiredReplicas, Label: "Desired Replicas", From: report.FromLatest, Datatype: report.Number, Priority: 5},
		report.Pod:         {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 6},
		Strategy:           {ID: Strategy, Label: "Strategy", From: report.FromLatest, Priority: 7},
	}

	DeploymentMetricTemplates = PodMetricTemplates

	DaemonSetMetadataTemplates = report.MetadataTemplates{
		NodeType:        {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:       {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:         {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		DesiredReplicas: {ID: DesiredReplicas, Label: "Desired Replicas", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		report.Pod:      {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 5},
	}

	DaemonSetMetricTemplates = PodMetricTemplates

	StatefulSetMetadataTemplates = report.MetadataTemplates{
		NodeType:           {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:          {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:            {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		ObservedGeneration: {ID: ObservedGeneration, Label: "Observed Gen.", From: report.FromLatest, Datatype: report.Number, Priority: 4},
		DesiredReplicas:    {ID: DesiredReplicas, Label: "Desired Replicas", From: report.FromLatest, Datatype: report.Number, Priority: 5},
		report.Pod:         {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 6},
	}

	StatefulSetMetricTemplates = PodMetricTemplates

	CronJobMetadataTemplates = report.MetadataTemplates{
		NodeType:      {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:     {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		Created:       {ID: Created, Label: "Created", From: report.FromLatest, Datatype: report.DateTime, Priority: 3},
		Schedule:      {ID: Schedule, Label: "Schedule", From: report.FromLatest, Priority: 4},
		LastScheduled: {ID: LastScheduled, Label: "Last Scheduled", From: report.FromLatest, Datatype: report.DateTime, Priority: 5},
		Suspended:     {ID: Suspended, Label: "Suspended", From: report.FromLatest, Priority: 6},
		ActiveJobs:    {ID: ActiveJobs, Label: "# Jobs", From: report.FromLatest, Datatype: report.Number, Priority: 7},
		report.Pod:    {ID: report.Pod, Label: "# Pods", From: report.FromCounters, Datatype: report.Number, Priority: 8},
	}

	CronJobMetricTemplates = PodMetricTemplates

	PersistentVolumeClaimMetadataTemplates = report.MetadataTemplates{
		NodeType:         {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Namespace:        {ID: Namespace, Label: "Namespace", From: report.FromLatest, Priority: 2},
		StorageClassName: {ID: StorageClassName, Label: "Storage Class", From: report.FromLatest, Priority: 3},
		Status:           {ID: Status, Label: "Status", From: report.FromLatest, Priority: 4},
		VolumeName:       {ID: VolumeName, Label: "Volume", From: report.FromLatest, Priority: 5},
		AccessModes:      {ID: AccessModes, Label: "Access Modes", From: report.FromLatest, Priority: 6},
	}

	PersistentVolumeClaimMetricTemplates = PodMetricTemplates

	PersistentVolumeMetadataTemplates = report.MetadataTemplates{
		NodeType:         {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Claim:            {ID: Claim, Label: "Claim", From: report.FromLatest, Priority: 2},
		StorageClassName: {ID: StorageClassName, Label: "Storage Class", From: report.FromLatest, Priority: 3},
		ReclaimPolicy:    {ID: ReclaimPolicy, Label: "Reclaim Policy", From: report.FromLatest, Priority: 4},
		AccessModes:      {ID: AccessModes, Label: "Access Modes", From: report.FromLatest, Priority: 5},
		Status:           {ID: Status, Label: "Status", From: report.FromLatest, Priority: 6},
		Message:          {ID: Message, Label: "Message", From: report.FromLatest, Priority: 7},
	}

	PersistentVolumeMetricTemplates = PodMetricTemplates

	StorageClassMetadataTemplates = report.MetadataTemplates{
		NodeType:        {ID: NodeType, Label: "Type", From: report.FromLatest, Priority: 1},
		Name:            {ID: Name, Label: "Name", From: report.FromLatest, Priority: 2},
		Provisioner:     {ID: Provisioner, Label: "Provisioner", From: report.FromLatest, Priority: 3},
		APIVersion:      {ID: APIVersion, Label: "APIVersion", From: report.FromLatest, Priority: 4},
		UID:             {ID: UID, Label: "UID", From: report.FromLatest, Priority: 5},
		ResourceVersion: {ID: ResourceVersion, Label: "ResourceVersion", From: report.FromLatest, Priority: 6},
		SelfLink:        {ID: SelfLink, Label: "SelfLink", From: report.FromLatest, Priority: 7},
	}

	StorageClassMetricTemplates = PodMetricTemplates

	TableTemplates = report.TableTemplates{
		LabelPrefix: {
			ID:     LabelPrefix,
			Label:  "Kubernetes Labels",
			Type:   report.PropertyListType,
			Prefix: LabelPrefix,
		},
	}

	ScalingControls = []report.Control{
		{
			ID:    ScaleDown,
			Human: "Scale Down",
			Icon:  "fa-minus",
			Rank:  0,
		},
		{
			ID:    ScaleUp,
			Human: "Scale Up",
			Icon:  "fa-plus",
			Rank:  1,
		},
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
				map[string]string{State: StateDeleted},
			),
		)
		r.probe.Publish(rpt)
	}
}

// IsPauseImageName indicates whether an image name corresponds to a
// kubernetes pause container image.
func IsPauseImageName(imageName string) bool {
	return strings.Contains(imageName, "google_containers/pause")

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
		return IsPauseImageName(imageName)
	}
	return false
}

// Tag adds pod parents to container nodes.
func (r *Reporter) Tag(rpt report.Report) (report.Report, error) {
	for id, n := range rpt.Container.Nodes {
		uid, ok := n.Latest.Lookup(docker.LabelPrefix + "io.kubernetes.pod.uid")
		if !ok {
			continue
		}

		// Tag the pause containers with "does-not-make-connections"
		if isPauseContainer(n, rpt) {
			n = n.WithLatest(report.DoesNotMakeConnections, mtime.Now(), "")
		}

		rpt.Container.Nodes[id] = n.WithParents(report.MakeSets().Add(
			report.Pod,
			report.MakeStringSet(report.MakePodNodeID(uid)),
		))
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
	hostTopology := r.hostTopology(services)
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
	deploymentTopology, deployments, err := r.deploymentTopology(r.probeID)
	if err != nil {
		return result, err
	}
	podTopology, err := r.podTopology(services, deployments, daemonSets, statefulSets, cronJobs)
	if err != nil {
		return result, err
	}
	namespaceTopology, err := r.namespaceTopology()
	if err != nil {
		return result, err
	}
	persistentVolumeClaimTopology, _, err := r.persistentVolumeClaimTopology(r.probeID)
	if err != nil {
		return result, err
	}
	persistentVolumeTopology, _, err := r.persistentVolumeTopology(r.probeID)
	if err != nil {
		return result, err
	}
	storageClassTopology, _, err := r.storageClassTopology(r.probeID)
	if err != nil {
		return result, err
	}
	applicationPodTopology, _, err := r.applicationPodTopology(r.probeID)
	if err != nil {
		return result, err
	}
	result.Pod = result.Pod.Merge(podTopology)
	result.Service = result.Service.Merge(serviceTopology)
	result.Host = result.Host.Merge(hostTopology)
	result.DaemonSet = result.DaemonSet.Merge(daemonSetTopology)
	result.StatefulSet = result.StatefulSet.Merge(statefulSetTopology)
	result.CronJob = result.CronJob.Merge(cronJobTopology)
	result.Deployment = result.Deployment.Merge(deploymentTopology)
	result.Namespace = result.Namespace.Merge(namespaceTopology)
	result.PersistentVolumeClaim = result.PersistentVolumeClaim.Merge(persistentVolumeClaimTopology)
	result.PersistentVolume = result.PersistentVolume.Merge(persistentVolumeTopology)
	result.StorageClass = result.StorageClass.Merge(storageClassTopology)
	result.ApplicationPod = result.ApplicationPod.Merge(applicationPodTopology)

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
	err := r.client.WalkServices(func(s Service) error {
		result.AddNode(s.GetNode())
		services = append(services, s)
		return nil
	})
	return result, services, err
}

// FIXME: Hideous hack to remove persistent-connection edges to
// virtual service IPs attributed to the internet. The global
// service-cluster-ip-range is not exposed by the API server (see
// https://github.com/kubernetes/kubernetes/issues/25533), so instead
// we synthesise it by computing the smallest network that contains
// all service IPs. That network may be smaller than the actual range
// but that is ok, since in the end all we care about is that it
// contains all the service IPs.
//
// The right way of fixing this is performing DNAT mapping on
// persistent connections for which we don't have a robust solution
// (see https://github.com/weaveworks/scope/issues/1491).
func (r *Reporter) hostTopology(services []Service) report.Topology {
	serviceIPs := make([]net.IP, 0, len(services))
	for _, service := range services {
		if ip := net.ParseIP(service.ClusterIP()).To4(); ip != nil {
			serviceIPs = append(serviceIPs, ip)
		}
	}
	serviceNetwork := report.ContainingIPv4Network(serviceIPs)
	if serviceNetwork == nil {
		return report.MakeTopology()
	}
	t := report.MakeTopology()
	t.AddNode(
		report.MakeNode(report.MakeHostNodeID(r.hostID)).
			WithSets(report.MakeSets().Add(host.LocalNetworks, report.MakeStringSet(serviceNetwork.String()))))
	return t
}

func (r *Reporter) deploymentTopology(probeID string) (report.Topology, []Deployment, error) {
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(DeploymentMetadataTemplates).
			WithMetricTemplates(DeploymentMetricTemplates).
			WithTableTemplates(TableTemplates)
		deployments = []Deployment{}
	)
	result.Controls.AddControls(ScalingControls)

	err := r.client.WalkDeployments(func(d Deployment) error {
		result.AddNode(d.GetNode(probeID))
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
	err := r.client.WalkDaemonSets(func(d DaemonSet) error {
		result.AddNode(d.GetNode())
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
	err := r.client.WalkStatefulSets(func(s StatefulSet) error {
		result.AddNode(s.GetNode())
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
	err := r.client.WalkCronJobs(func(c CronJob) error {
		result.AddNode(c.GetNode())
		cronJobs = append(cronJobs, c)
		return nil
	})
	return result, cronJobs, err
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

func (r *Reporter) podTopology(services []Service, deployments []Deployment, daemonSets []DaemonSet, statefulSets []StatefulSet, cronJobs []CronJob) (report.Topology, error) {
	var (
		pods = report.MakeTopology().
			WithMetadataTemplates(PodMetadataTemplates).
			WithMetricTemplates(PodMetricTemplates).
			WithTableTemplates(TableTemplates)
		selectors = []func(labelledChild){}
	)
	pods.Controls.AddControl(report.Control{
		ID:    GetLogs,
		Human: "Get logs",
		Icon:  "fa-desktop",
		Rank:  0,
	})
	pods.Controls.AddControl(report.Control{
		ID:    DeletePod,
		Human: "Delete",
		Icon:  "fa-trash-o",
		Rank:  1,
	})
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
	}

	var localPodUIDs map[string]struct{}
	if r.nodeName == "" {
		// We don't know the node name: fall back to obtaining the local pods from kubelet
		var err error
		localPodUIDs, err = GetLocalPodUIDs(fmt.Sprintf("127.0.0.1:%d", r.kubeletPort))
		if err != nil {
			log.Warnf("No node name and cannot obtain local pods, reporting all (which may impact performance): %v", err)
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

func (r *Reporter) persistentVolumeClaimTopology(probeID string) (report.Topology, []PersistentVolumeClaim, error) {
	//TODO: Need to improve logic for topology.
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(PersistentVolumeClaimMetadataTemplates).
			WithMetricTemplates(PersistentVolumeClaimMetricTemplates).
			WithTableTemplates(TableTemplates)
		pvcs = []PersistentVolumeClaim{}
	)

	err := r.client.WalkPersistentVolumeClaims(func(p PersistentVolumeClaim) error {
		result.AddNode(p.GetNode(probeID))
		pvcs = append(pvcs, p)
		return nil
	})
	return result, pvcs, err
}

func (r *Reporter) persistentVolumeTopology(probeID string) (report.Topology, []PersistentVolume, error) {
	//TODO: Need to improve logic for topology.
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(PersistentVolumeMetadataTemplates).
			WithMetricTemplates(PersistentVolumeMetricTemplates).
			WithTableTemplates(TableTemplates)
		pvcs = []PersistentVolume{}
	)

	err := r.client.WalkPersistentVolumes(func(p PersistentVolume) error {
		result.AddNode(p.GetNode(probeID))
		pvcs = append(pvcs, p)
		return nil
	})
	return result, pvcs, err
}

func (r *Reporter) storageClassTopology(probeID string) (report.Topology, []StorageClass, error) {
	//TODO: Need to improve logic for topology.
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(StorageClassMetadataTemplates).
			WithMetricTemplates(StorageClassMetricTemplates).
			WithTableTemplates(TableTemplates)
		sc = []StorageClass{}
	)

	err := r.client.WalkStorageClasses(func(p StorageClass) error {
		result.AddNode(p.GetNode(probeID))
		sc = append(sc, p)
		return nil
	})
	return result, sc, err
}

func (r *Reporter) applicationPodTopology(probeID string) (report.Topology, []ApplicationPod, error) {
	//TODO: Need to improve logic for topology.
	var (
		result = report.MakeTopology().
			WithMetadataTemplates(PodMetadataTemplates).
			WithMetricTemplates(PodMetricTemplates).
			WithTableTemplates(TableTemplates)
		ap = []ApplicationPod{}
	)

	err := r.client.WalkApplicationPods(func(p ApplicationPod) error {
		if p.GetApp(probeID).ID != "" {
			result.AddNode(p.GetApp(probeID))
			ap = append(ap, p)
			return nil
		}
		return nil
	})
	return result, ap, err
}
