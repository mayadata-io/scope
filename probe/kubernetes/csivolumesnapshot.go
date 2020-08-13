package kubernetes

import (
	csisnapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	"github.com/weaveworks/scope/report"
)

const (
	DriverAnnotation = "driver"
)

// CsiVolumeSnapshot represent kubernetes VolumeSnapshot interface
type CsiVolumeSnapshot interface {
	Meta
	GetNode(probeID string) report.Node
	GetVolumeName() string
	GetCapacity() string
	GetDriver() string
}

// csiVolumeSnapshot represents kubernetes volume snapshots
type csiVolumeSnapshot struct {
	*csisnapshotv1beta1.VolumeSnapshot
	Meta
}

// NewCsiVolumeSnapshot returns new Volume Snapshot type
func NewCsiVolumeSnapshot(p *csisnapshotv1beta1.VolumeSnapshot) CsiVolumeSnapshot {
	return &csiVolumeSnapshot{VolumeSnapshot: p, Meta: meta{p.ObjectMeta}}
}

// GetVolumeName returns the PVC name for volume snapshot
func (p *csiVolumeSnapshot) GetVolumeName() string {
	return *p.Spec.Source.PersistentVolumeClaimName
}

// GetCapacity returns the capacity of the source PVC stored in annotation
func (p *csiVolumeSnapshot) GetCapacity() string {
	capacity := p.GetAnnotations()[Capacity]
	if capacity != "" {
		return capacity
	}
	return ""
}

// GetCapacity returns the capacity of the source PVC stored in annotation
func (p *csiVolumeSnapshot) GetDriver() string {
	driver := p.GetAnnotations()[DriverAnnotation]
	if driver != "" {
		return driver
	}
	return ""
}

// GetNode returns CsiVolumeSnapshot as Node
func (p *csiVolumeSnapshot) GetNode(probeID string) report.Node {
	return p.MetaNode(report.MakeCsiVolumeSnapshotNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID: probeID,
		NodeType:              "Volume Snapshot",
		VolumeClaim:           p.GetVolumeName(),
		SnapshotClass:         *p.Spec.VolumeSnapshotClassName,
		SnapshotData:          *p.Status.BoundVolumeSnapshotContentName,
	}).WithLatestActiveControls(CloneCsiVolumeSnapshot, DeleteCsiVolumeSnapshot, Describe)
}
