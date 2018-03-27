package kubernetes

import (
	"strconv"

	"github.com/weaveworks/scope/report"

	apiv1 "k8s.io/api/core/v1"
)

// ApplicationPod represents a Kubernetes pod
type ApplicationPod interface {
	Meta
	GetApp(probeID string) report.Node
	RestartCount() uint
}

type applicationPod struct {
	*apiv1.Pod
	Meta
	parents report.Sets
	Node    *apiv1.Node
}

func (p *applicationPod) State() string {
	return string(p.Status.Phase)
}

func (p *applicationPod) RestartCount() uint {
	count := uint(0)
	for _, cs := range p.Status.ContainerStatuses {
		count += uint(cs.RestartCount)
	}
	return count
}

func (p *applicationPod) EmptyMetaNode() report.Node {
	return report.MakeNode("")
}

func (p *applicationPod) AddParent(topology, id string) {
	p.parents = p.parents.Add(topology, report.MakeStringSet(id))
}

// NewApp creates a new Pod
func NewApp(p *apiv1.Pod) ApplicationPod {
	return &applicationPod{
		Pod:     p,
		Meta:    meta{p.ObjectMeta},
		parents: report.MakeSets(),
	}
}

func (p *applicationPod) GetApp(probeID string) report.Node {
	for _, v := range p.Spec.Volumes {
		if v.VolumeSource.PersistentVolumeClaim != nil {
			latests := map[string]string{
				State: p.State(),
				IP:    p.Status.PodIP,
				report.ControlProbeID: probeID,
				RestartCount:          strconv.FormatUint(uint64(p.RestartCount()), 10),
				VolumeClaimName:       v.VolumeSource.PersistentVolumeClaim.ClaimName,
			}

			if p.Pod.Spec.HostNetwork {
				latests[IsInHostNetwork] = "true"
			}

			return p.MetaNode(report.MakePodNodeID(p.UID())).WithLatests(latests).
				WithParents(p.parents)
		}
	}
	return p.EmptyMetaNode()
}
