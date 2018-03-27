package render

import (
	"github.com/weaveworks/scope/report"
)

// PVCRenderer to create a renderer for PVC objects.
var PVCRenderer = MakeReduce(

	MapEndpoints(endpoint2PVC, report.PersistentVolumeClaim),
	MapEndpoints(endpoint2PV, report.PersistentVolume),
	MapEndpoints(endpoint2StorageClass, report.StorageClass),
	MapEndpoints(endpoint2PVC, report.ApplicationPod),
)

// endpoint2PVC returns pvc node ID
func endpoint2PVC(n report.Node) string {
	if pvcNodeID, ok := n.Latest.Lookup(report.MakePersistentVolumeClaimNodeID(n.ID)); ok {
		return pvcNodeID
	}
	return ""
}

// endpoint2PV returns pv node ID
func endpoint2PV(n report.Node) string {
	if pvNodeID, ok := n.Latest.Lookup(report.MakePersistentVolumeNodeID(n.ID)); ok {
		return pvNodeID
	}
	return ""
}

// endpoint2StorageClass returns StorageClass node ID
func endpoint2StorageClass(n report.Node) string {
	if storageclassNodeID, ok := n.Latest.Lookup(report.MakeStorageClassNodeID(n.ID)); ok {
		return storageclassNodeID
	}
	return ""
}
