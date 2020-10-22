package kubernetes

import (
	"strconv"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// BlockDevice represent NDM BlockDevice interface
type BlockDevice interface {
	Meta
	GetNode(probeID string) report.Node
}

// blockDevice represents NDM blockDevices
type blockDevice struct {
	*mayav1alpha1.BlockDevice
	Meta
}

// NewBlockDevice returns new block device type
func NewBlockDevice(b *mayav1alpha1.BlockDevice) BlockDevice {
	return &blockDevice{BlockDevice: b, Meta: meta{b.ObjectMeta}}
}

// GetNode returns Block Device as Node
func (b *blockDevice) GetNode(probeID string) report.Node {
	return b.MetaNode(report.MakeBlockDeviceNodeID(b.UID())).WithLatests(map[string]string{
		NodeType:              "Block Device",
		LogicalSectorSize:     strconv.Itoa(int(b.Spec.Capacity.LogicalSectorSize)),
		Storage:               strconv.Itoa(int(b.Spec.Capacity.Storage / (1024 * 1024 * 1024))),
		FirmwareRevision:      b.Spec.Details.FirmwareRevision,
		Model:                 b.Spec.Details.Model,
		Serial:                b.Spec.Details.Serial,
		Vendor:                b.Spec.Details.Vendor,
		HostName:              b.GetLabels()["kubernetes.io/hostname"],
		Path:                  b.Spec.Path,
		report.ControlProbeID: probeID,
	}).WithLatestActiveControls(Describe)
}
