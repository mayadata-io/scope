package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorVolume interface
type CStorVolume interface {
	Meta
	GetNode() report.Node
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
	return p.MetaNode(report.MakeCStorVolumeNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:   "CStor Volume",
		APIVersion: p.APIVersion,
	})
}
