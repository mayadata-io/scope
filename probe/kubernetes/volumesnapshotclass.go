package kubernetes

import (
	csisnapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	"github.com/weaveworks/scope/report"
)

// VolumeSnapshotClass represent kubernetes VolumeSnapshotClass interface
type VolumeSnapshotClass interface {
	Meta
	GetNode(probeID string) report.Node
}

// volumeSnapshot represents kubernetes volume snapshots class
type volumeSnapshotClass struct {
	*csisnapshotv1beta1.VolumeSnapshotClass
	Meta
}

// NewVolumeSnapshotClass returns new Volume Snapshot Class type
func NewVolumeSnapshotClass(p *csisnapshotv1beta1.VolumeSnapshotClass) VolumeSnapshotClass {
	return &volumeSnapshotClass{VolumeSnapshotClass: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns VolumeSnapshotClass as Node
func (p *volumeSnapshotClass) GetNode(probeID string) report.Node {
	return p.MetaNode(report.MakeVolumeSnapshotClassNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID: probeID,
		NodeType:              "Volume Snapshot Class",
		Driver:                p.Driver,
		DeletionPolicy:        string(p.DeletionPolicy),
	}).WithLatestActiveControls(Describe)
}
