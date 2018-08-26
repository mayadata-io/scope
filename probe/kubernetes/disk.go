package kubernetes

import (
	"strconv"

	ndmv1alpha1 "github.com/openebs/node-disk-manager/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// Disk represent NDM Disk interface
type Disk interface {
	Meta
	GetNode() report.Node
}

// disk represents NDM Disks
type disk struct {
	*ndmv1alpha1.Disk
	Meta
}

// NewDisk returns new Disk type
func NewDisk(p *ndmv1alpha1.Disk) Disk {
	return &disk{Disk: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns Disk as Node
func (p *disk) GetNode() report.Node {
	return p.MetaNode(report.MakeDiskNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:          "Disk",
		Name:              p.GetName(),
		LogicalSectorSize: strconv.Itoa(int(p.Spec.Capacity.LogicalSectorSize)),
		Storage:           strconv.Itoa(int(p.Spec.Capacity.Storage / (1024 * 1024 * 1024))),
		FirmwareRevision:  p.Spec.Details.FirmwareRevision,
		Model:             p.Spec.Details.Model,
		Serial:            p.Spec.Details.Serial,
		SpcVersion:        p.Spec.Details.SPCVersion,
		Vendor:            p.Spec.Details.Vendor,
		HostName:          p.GetLabels()["kubernetes.io/hostname"],
	})
}
