package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// NewTestCStorPool represent NewTestCStorPool interface.
type NewTestCStorPool interface {
	Meta
	GetNode(probeID string) report.Node
	GetBlockDeviceList() string
}

// newTestCStorPool represent the newTestCStorPool CRD of Kubernetes.
type newTestCStorPool struct {
	*mayav1alpha1.NewTestCStorPool
	Meta
}

// NewNewTestCStorPool return new NewNewTestCStorPool type.
func NewNewTestCStorPool(c *mayav1alpha1.NewTestCStorPool) NewTestCStorPool {
	return &newTestCStorPool{NewTestCStorPool: c, Meta: meta{c.ObjectMeta}}
}

func (c *newTestCStorPool) GetBlockDeviceList() string {
	blockDeviceList := []string{}
	for _, raidGroup := range c.Spec.RaidGroups {
		for _, blockDevice := range raidGroup.BlockDevices {
			blockDeviceList = append(blockDeviceList, blockDevice.BlockDeviceName)
		}
	}
	return strings.Join(blockDeviceList, report.ScopeDelim)
}

// GetNode returns CStorPoolCluster as Node
func (c *newTestCStorPool) GetNode(probeID string) report.Node {
	return c.MetaNode(report.MakeNewTestCStorPoolNodeID(c.UID())).WithLatests(map[string]string{
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
