package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorPoolInstance represent CStorPoolInstance interface.
type CStorPoolInstance interface {
	Meta
	GetNode(probeID string) report.Node
	GetBlockDeviceList() string
}

// cStorPoolInstance represent the cStorPoolInstance CRD of Kubernetes.
type cStorPoolInstance struct {
	*mayav1alpha1.CStorPoolInstance
	Meta
}

// NewCStorPoolInstance return new NewCStorPoolInstance type.
func NewCStorPoolInstance(c *mayav1alpha1.CStorPoolInstance) CStorPoolInstance {
	return &cStorPoolInstance{CStorPoolInstance: c, Meta: meta{c.ObjectMeta}}
}

func (c *cStorPoolInstance) GetBlockDeviceList() string {
	blockDeviceList := []string{}
	for _, raidGroup := range c.Spec.RaidGroups {
		for _, blockDevice := range raidGroup.BlockDevices {
			blockDeviceList = append(blockDeviceList, blockDevice.BlockDeviceName)
		}
	}
	return strings.Join(blockDeviceList, report.ScopeDelim)
}

// GetNode returns CStorPoolCluster as Node
func (c *cStorPoolInstance) GetNode(probeID string) report.Node {
	return c.MetaNode(report.MakeCStorPoolInstanceNodeID(c.UID())).WithLatests(map[string]string{
		NodeType:              "CStor Pool",
		Status:                string(c.Status.Phase),
		TotalSize:             c.Status.Capacity.Total,
		FreeSize:              c.Status.Capacity.Free,
		UsedSize:              c.Status.Capacity.Used,
		BlockDeviceList:       c.GetBlockDeviceList(),
		StoragePoolClaimName:  c.GetLabels()["openebs.io/cstor-pool-cluster"],
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
