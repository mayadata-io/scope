package kubernetes

import (
	"github.com/weaveworks/scope/report"
	apiv1 "k8s.io/api/core/v1"
)

// PersistentVolume represent kubernetes PV interface
type PersistentVolume interface {
	Meta
	GetNode(probeID string) report.Node
}

// persistentVolume represents kubernetes PV
type persistentVolume struct {
	*apiv1.PersistentVolume
	Meta
}

// NewPV returns new PV type
func NewPV(p *apiv1.PersistentVolume) PersistentVolume {
	return &persistentVolume{PersistentVolume: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns PV as Node
func (p *persistentVolume) GetNode(probeID string) report.Node {
	return p.MetaNode(report.MakePersistentVolumeNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID: probeID,
		NodeType:              "PV",
		Namespace:             p.GetNamespace(),
	})
}
