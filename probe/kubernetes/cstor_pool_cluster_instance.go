package kubernetes

import (
	"strconv"
	"strings"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
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
	*cstorv1.CStorPoolInstance
	Meta
}

// NewCStorPoolInstance return new NewCStorPoolInstance type.
func NewCStorPoolInstance(c *cstorv1.CStorPoolInstance) CStorPoolInstance {
	return &cStorPoolInstance{CStorPoolInstance: c, Meta: meta{c.ObjectMeta}}
}

func (c *cStorPoolInstance) GetBlockDeviceList() string {
	blockDeviceList := []string{}
	for _, raidGroup := range c.Spec.DataRaidGroups {
		for _, blockDevice := range raidGroup.CStorPoolInstanceBlockDevices {
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
		TotalSize:             c.Status.Capacity.Total.String(),
		FreeSize:              c.Status.Capacity.Free.String(),
		UsedSize:              c.Status.Capacity.Used.String(),
		LogicalUsed:           c.Status.Capacity.ZFS.LogicalUsed.String(),
		ReadOnly:              strconv.FormatBool(c.Status.ReadOnly),
		ProvisionedReplicas:   strconv.Itoa(int(c.Status.ProvisionedReplicas)),
		HealthyReplicas:       strconv.Itoa(int(c.Status.HealthyReplicas)),
		BlockDeviceList:       c.GetBlockDeviceList(),
		StoragePoolClaimName:  c.GetLabels()["openebs.io/cstor-pool-cluster"],
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
