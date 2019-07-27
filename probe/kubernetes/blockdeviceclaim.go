package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// BlockDeviceClaim represent NDM BlockDeviceClaim interface
type BlockDeviceClaim interface {
	Meta
	GetNode(probeID string) report.Node
}

type blockDeviceClaim struct {
	*mayav1alpha1.BlockDeviceClaim
	Meta
}

// NewBlockDeviceClaim returns new block device claim type
func NewBlockDeviceClaim(b *mayav1alpha1.BlockDeviceClaim) BlockDeviceClaim {
	return &blockDeviceClaim{BlockDeviceClaim: b, Meta: meta{b.ObjectMeta}}
}

// GetNode returns Block Device Claim as Node
func (b *blockDeviceClaim) GetNode(probeID string) report.Node {
	return b.MetaNode(report.MakeBlockDeviceClaimNodeID(b.UID())).WithLatests(map[string]string{
		NodeType:              "Block Device Claim",
		BlockDeviceName:       b.Spec.BlockDeviceName,
		HostName:              b.Spec.HostName,
		Status:                string(b.Status.Phase),
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
