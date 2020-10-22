package kubernetes

import (
	csisnapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	"github.com/weaveworks/scope/report"
)

// VolumeSnapshotContent represent kubernetes VolumeSnapshotContent interface
type VolumeSnapshotContent interface {
	Meta
	GetNode(probeID string) report.Node
}

// volumeSnapshotContent represents kubernetes volume snapshot content
type volumeSnapshotContent struct {
	*csisnapshotv1beta1.VolumeSnapshotContent
	Meta
}

// NewVolumeSnapshotContent returns new Volume Snapshot Content type
func NewVolumeSnapshotContent(p *csisnapshotv1beta1.VolumeSnapshotContent) VolumeSnapshotContent {
	return &volumeSnapshotContent{VolumeSnapshotContent: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns VolumeSnapshotContent as Node
func (p *volumeSnapshotContent) GetNode(probeID string) report.Node {
	return p.MetaNode(report.MakeVolumeSnapshotContentNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID:   probeID,
		NodeType:                "Volume Snapshot Content",
		VolumeSnapshotName:      p.Spec.VolumeSnapshotRef.Name,
		VolumeSnapshotNamespace: p.Spec.VolumeSnapshotRef.Namespace,
	}).WithLatestActiveControls(Describe)
}
