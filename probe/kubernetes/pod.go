package kubernetes

import (
	"strconv"
	"strings"

	"github.com/weaveworks/scope/report"

	apiv1 "k8s.io/api/core/v1"
)

// These constants are keys used in node metadata
const (
	State           = report.KubernetesState
	IsInHostNetwork = report.KubernetesIsInHostNetwork
	RestartCount    = report.KubernetesRestartCount
)

// Pod states we handle specially
const (
	StateDeleted = "deleted"
	StateFailed  = "Failed"
)

// Pod labels to get pv name, if it is a controller/target or replica pod
const (
	PersistentVolumeLabel = "openebs.io/persistent-volume"
	VSMLabel              = "vsm"
	PVLabel               = "openebs.io/pv"
)

// Pod label to distinguish replica pod or cstor pool pod
const (
	AppLabel         = "app"
	AppValue         = "cstor-pool"
	ReplicaPodLabel  = "openebs.io/replica"
	JivaReplicaValue = "jiva-replica"
)

// Pod represents a Kubernetes pod
type Pod interface {
	Meta
	AddParent(topology, id string)
	NodeName() string
	GetNode(probeID string) report.Node
	GetVolumeName() string
	IsReplicaOrPoolPod() bool
	VolumeClaimName() string
	RestartCount() uint
	ContainerNames() []string
}

type pod struct {
	*apiv1.Pod
	Meta
	parents report.Sets
	Node    *apiv1.Node
}

// NewPod creates a new Pod
func NewPod(p *apiv1.Pod) Pod {
	return &pod{
		Pod:     p,
		Meta:    meta{p.ObjectMeta},
		parents: report.MakeSets(),
	}
}

func (p *pod) UID() string {
	// Work around for master pod not reporting the right UID.
	if hash, ok := p.ObjectMeta.Annotations["kubernetes.io/config.hash"]; ok {
		return hash
	}
	return p.Meta.UID()
}

func (p *pod) AddParent(topology, id string) {
	p.parents = p.parents.AddString(topology, id)
}

func (p *pod) State() string {
	return string(p.Status.Phase)
}

func (p *pod) NodeName() string {
	return p.Spec.NodeName
}

func (p *pod) RestartCount() uint {
	count := uint(0)
	for _, cs := range p.Status.ContainerStatuses {
		count += uint(cs.RestartCount)
	}
	return count
}

func (p *pod) IsReplicaOrPoolPod() bool {
	replicaPod, _ := p.GetLabels()[ReplicaPodLabel]
	cstorPoolPod, _ := p.GetLabels()[AppLabel]
	if replicaPod == JivaReplicaValue || cstorPoolPod == AppValue {
		return true
	}
	return false
}

func (p *pod) GetVolumeName() string {
	if strings.Contains(p.GetName(), "-rep-") {
		return ""
	}

	var volumeName string
	var ok bool

	if volumeName, ok = p.GetLabels()[VSMLabel]; ok {
		return volumeName
	} else if volumeName, ok = p.GetLabels()[PersistentVolumeLabel]; ok {
		return volumeName
	} else if volumeName, ok = p.GetLabels()[PVLabel]; ok {
		return volumeName
	}
	return ""
}

func (p *pod) VolumeClaimName() string {
	var claimName string
	for _, volume := range p.Spec.Volumes {
		if volume.VolumeSource.PersistentVolumeClaim != nil {
			claimName = volume.VolumeSource.PersistentVolumeClaim.ClaimName
			break
		}
	}
	return claimName
}

func (p *pod) GetNode(probeID string) report.Node {
	latests := map[string]string{
		State: p.State(),
		IP:    p.Status.PodIP,
		report.ControlProbeID: probeID,
		RestartCount:          strconv.FormatUint(uint64(p.RestartCount()), 10),
	}

	if p.VolumeClaimName() != "" {
		latests[VolumeClaim] = p.VolumeClaimName()
		latests[VolumePod] = "true"
	}

	if p.GetVolumeName() != "" {
		latests[VolumeName] = p.GetVolumeName()
		latests[VolumePod] = "true"
	}

	if p.IsReplicaOrPoolPod() {
		latests[VolumePod] = "true"
	}

	if p.Pod.Spec.HostNetwork {
		latests[IsInHostNetwork] = "true"
	}

	return p.MetaNode(report.MakePodNodeID(p.UID())).WithLatests(latests).
		WithParents(p.parents).
		WithLatestActiveControls(GetLogs, DeletePod)
}

func (p *pod) ContainerNames() []string {
	containerNames := make([]string, 0, len(p.Pod.Spec.Containers))
	for _, c := range p.Pod.Spec.Containers {
		containerNames = append(containerNames, c.Name)
	}
	return containerNames
}
