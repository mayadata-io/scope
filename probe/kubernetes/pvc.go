package kubernetes

import (
	"github.com/weaveworks/scope/report"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/labels"
)

// PersistentVolumeClaim represent kubernetes PVC interface
type PersistentVolumeClaim interface {
	Meta
	Selector() (labels.Selector, error)
	GetNode(probeID string) report.Node
}

// persistentVolumeClaim represents kubernetes PVC
type persistentVolumeClaim struct {
	*apiv1.PersistentVolumeClaim
	Meta
}

// NewPVC returns new PVC type
func NewPVC(p *apiv1.PersistentVolumeClaim) PersistentVolumeClaim {
	return &persistentVolumeClaim{PersistentVolumeClaim: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns PVC as Node
func (p *persistentVolumeClaim) GetNode(probeID string) report.Node {
	return p.MetaNode(report.MakePersistentVolumeClaimNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID: probeID,
		NodeType:              "PVC",
		Namespace:             p.GetNamespace(),
	})
}

// Selector returns all PVC selector
func (p *persistentVolumeClaim) Selector() (labels.Selector, error) {
	selector, err := metav1.LabelSelectorAsSelector(p.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return selector, nil
}
