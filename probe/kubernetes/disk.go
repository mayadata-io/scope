package kubernetes

import (
	"fmt"
	"strconv"
	"strings"

	maya1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// Disk represent NDM Disk interface
type Disk interface {
	Meta
	GetNode(probeID string) report.Node
	GetPath() []string
	GetNodeTagStatus(string) string
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
func (p *disk) GetNode(probeID string) report.Node {
	var diskStatus string
	diskStatus = p.Status.State

	latests := map[string]string{
		NodeType:              "Disk",
		PhysicalSectorSize:    strconv.Itoa(int(p.Spec.Capacity.PhysicalSectorSize)),
		LogicalSectorSize:     strconv.Itoa(int(p.Spec.Capacity.LogicalSectorSize)),
		VolumeCapacity:        strconv.Itoa(int(p.Spec.Capacity.Storage/(1024*1024*1024))) + " GB",
		FirmwareRevision:      p.Spec.Details.FirmwareRevision,
		Model:                 p.Spec.Details.Model,
		RotationRate:          strconv.Itoa(int(p.Spec.Details.RotationRate)),
		Serial:                p.Spec.Details.Serial,
		Vendor:                p.Spec.Details.Vendor,
		HostName:              p.GetLabels()["kubernetes.io/hostname"],
		DiskList:              strings.Join(p.GetPath(), report.ScopeDelim),
		Status:                diskStatus,
		CurrentTemperature:    strconv.Itoa(int(p.Stats.TempInfo.CurrentTemperature)),
		HighestTemperature:    strconv.Itoa(int(p.Stats.TempInfo.HighestTemperature)),
		LowestTemperature:     strconv.Itoa(int(p.Stats.TempInfo.LowestTemperature)),
		TotalBytesRead:        strconv.Itoa(int(p.Stats.TotalBytesRead)),
		TotalBytesWritten:     strconv.Itoa(int(p.Stats.TotalBytesWritten)),
		DeviceUtilizationRate: fmt.Sprintf("%.2f", p.Stats.DeviceUtilizationRate),
		PercentEnduranceUsed:  fmt.Sprintf("%.2f", p.Stats.PercentEnduranceUsed),
		Created:               p.ObjectMeta.CreationTimestamp.String(),
		Path:                  p.Spec.Path,
		report.ControlProbeID: probeID,
	}

	return p.MetaNode(report.MakeDiskNodeID(p.UID())).
		WithLatests(latests).
		WithNodeTag(p.GetNodeTagStatus(diskStatus)).
		WithLatestActiveControls(Describe)
}

func (p *disk) GetPath() []string {
	if len(p.Spec.DevLinks) > 0 {
		return p.Spec.DevLinks[0].Links
	}
	diskList := []string{p.Spec.Path}
	return diskList
}

func (p *disk) GetNodeTagStatus(status string) string {
	if len(status) > 0 {
		return CStorVolumeStatusMap[strings.ToLower(status)]
	}
	return ""
}
