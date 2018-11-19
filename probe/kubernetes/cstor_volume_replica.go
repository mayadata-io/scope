package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorVolumeReplica interface
type CStorVolumeReplica interface {
	Meta
	GetNode() report.Node
}

// cStorVolume represents cStor Volume Replica CR
type cStorVolumeReplica struct {
	*mayav1alpha1.CStorVolumeReplica
	Meta
}

// NewCStorVolumeReplica returns fresh CStorVolumeReplica instance
func NewCStorVolumeReplica(p *mayav1alpha1.CStorVolumeReplica) CStorVolumeReplica {
	return &cStorVolumeReplica{CStorVolumeReplica: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorVolumeReplica) GetNode() report.Node {
	return p.MetaNode(report.MakeCStorVolumeNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:   "CStor Volume Replica",
		APIVersion: p.APIVersion,
	})
}
