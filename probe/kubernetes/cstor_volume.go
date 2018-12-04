package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorVolume interface
type CStorVolume interface {
	Meta
	GetNode() report.Node
	GetPersistentVolumeName() string
	GetStatus() string
	GetNodeTagOnStatus(status string) string
}

// cStorVolume represents cStor Volume CR
type cStorVolume struct {
	*mayav1alpha1.CStorVolume
	Meta
}

// NewCStorVolume returns fresh CStorVolume instance
func NewCStorVolume(p *mayav1alpha1.CStorVolume) CStorVolume {
	return &cStorVolume{CStorVolume: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorVolume) GetNode() report.Node {
	status := p.GetStatus()
	latests := map[string]string{
		NodeType:   "CStor Volume",
		APIVersion: p.APIVersion,
	}

	if p.GetPersistentVolumeName() != "" {
		latests[VolumeName] = p.GetPersistentVolumeName()
	}

	if status != "" {
		latests[CStorVolumeStatus] = status
	}
	return p.MetaNode(report.MakeCStorVolumeNodeID(p.Name())).
		WithLatests(latests).
		WithNodeTag(p.GetNodeTagOnStatus(strings.ToLower(status)))
}

func (p *cStorVolume) GetPersistentVolumeName() string {
	persistentVolumeName := p.Labels()["openebs.io/persistent-volume"]
	return persistentVolumeName
}

func (p *cStorVolume) GetStatus() string {
	status := p.Status.Phase
	return string(status)
}

func (p *cStorVolume) GetNodeTagOnStatus(status string) string {
	if status != "" {
		return CStorVolumeStatusMap[status]
	}
	return "unknown"
}
