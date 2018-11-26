package kubernetes

import (
	"strconv"
	"strings"

	maya1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// Disk represent NDM Disk interface
type Disk interface {
	Meta
	GetNode() report.Node
	GetPath() []string
}

// disk represents NDM Disks
type disk struct {
	*maya1alpha1.Disk
	Meta
}

// NewDisk returns new Disk type
func NewDisk(p *maya1alpha1.Disk) Disk {
	return &disk{Disk: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns Disk as Node
func (p *disk) GetNode() report.Node {
	return p.MetaNode(report.MakeDiskNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:          "Disk",
		LogicalSectorSize: strconv.Itoa(int(p.Spec.Capacity.LogicalSectorSize)),
		Storage:           strconv.Itoa(int(p.Spec.Capacity.Storage/(1024*1024*1024))) + " GB",
		FirmwareRevision:  p.Spec.Details.FirmwareRevision,
		Model:             p.Spec.Details.Model,
		Serial:            p.Spec.Details.Serial,
		Vendor:            p.Spec.Details.Vendor,
		HostName:          p.GetLabels()["kubernetes.io/hostname"],
		DiskList:          strings.Join(p.GetPath(), "~p$"),
	})
}

func (p *disk) GetPath() []string {
	if len(p.Spec.DevLinks) > 0 {
		return p.Spec.DevLinks[0].Links
	}
	diskList := []string{p.Spec.Path}
	return diskList
}
