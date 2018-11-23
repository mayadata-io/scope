package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorPool interface
type CStorPool interface {
	Meta
	GetNode() report.Node
	GetStatus() string
}

// cStorPool represents cStor Volume CSP
type cStorPool struct {
	*mayav1alpha1.CStorPool
	Meta
}

// NewCStorPool returns fresh CStorPool instance
func NewCStorPool(p *mayav1alpha1.CStorPool) CStorPool {
	return &cStorPool{CStorPool: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorPool) GetNode() report.Node {
	latests := map[string]string{
		NodeType:   "CStor Pool",
		APIVersion: p.APIVersion,
	}

	if p.GetStatus() != "" {
		latests[CStorPoolStatus] = p.GetStatus()
	}
	return p.MetaNode(report.MakeCStorPoolNodeID(p.UID())).WithLatests(latests)
}

func (p *cStorPool) GetStatus() string {
	status := p.Status.Phase
	return string(status)
}
