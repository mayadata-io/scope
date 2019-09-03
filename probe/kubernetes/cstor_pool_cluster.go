package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorPoolCluster represent CStorPoolCluster interface.
type CStorPoolCluster interface {
	Meta
	GetNode(probeID string) report.Node
}

// cStorPoolCluster represent the cStorPoolCluster CRD of Kubernetes.
type cStorPoolCluster struct {
	*mayav1alpha1.CStorPoolCluster
	Meta
}

// NewCStorPoolCluster return new CStorPoolCluster type.
func NewCStorPoolCluster(c *mayav1alpha1.CStorPoolCluster) CStorPoolCluster {
	return &cStorPoolCluster{CStorPoolCluster: c, Meta: meta{c.ObjectMeta}}
}

// GetNode returns CStorPoolCluster as Node
func (c *cStorPoolCluster) GetNode(probeID string) report.Node {
	return c.MetaNode(report.MakeCStorPoolClusterNodeID(c.UID())).WithLatests(map[string]string{
		NodeType:              "CStor Pool Cluster",
		Status:                c.Status.Phase,
		TotalSize:             c.Status.Capacity.Total,
		FreeSize:              c.Status.Capacity.Free,
		UsedSize:              c.Status.Capacity.Used,
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
