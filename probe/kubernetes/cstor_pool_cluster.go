package kubernetes

import (
	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
	"github.com/weaveworks/scope/report"
	"strconv"
)

// CStorPoolCluster represent CStorPoolCluster interface.
type CStorPoolCluster interface {
	Meta
	GetNode(probeID string) report.Node
}

// cStorPoolCluster represent the cStorPoolCluster CRD of Kubernetes.
type cStorPoolCluster struct {
	*cstorv1.CStorPoolCluster
	Meta
}

// NewCStorPoolCluster return new CStorPoolCluster type.
func NewCStorPoolCluster(c *cstorv1.CStorPoolCluster) CStorPoolCluster {
	return &cStorPoolCluster{CStorPoolCluster: c, Meta: meta{c.ObjectMeta}}
}

// GetNode returns CStorPoolCluster as Node
func (c *cStorPoolCluster) GetNode(probeID string) report.Node {
	return c.MetaNode(report.MakeCStorPoolClusterNodeID(c.UID())).WithLatests(map[string]string{
		NodeType:              "CStor Pool Cluster",
		ProvisionedInstances:  strconv.Itoa(int(c.Status.ProvisionedInstances)),
		DesiredInstances:      strconv.Itoa(int(c.Status.DesiredInstances)),
		HealthyInstances:      strconv.Itoa(int(c.Status.HealthyInstances)),
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
