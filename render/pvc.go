package render

import (
	"github.com/weaveworks/scope/report"
)

// PVCRenderer to create a renderer for PVC objects.
var PVCRenderer = MakeReduce(

	MapEndpoints(endpoint2PVC, report.PersistentVolumeClaim),
)

// endpoint2PVC returns pvc node ID
func endpoint2PVC(n report.Node) string {
	if pvcNodeID, ok := n.Latest.Lookup(report.MakePersistentVolumeClaimNodeID(n.ID)); ok {
		return pvcNodeID
	}
	return ""
}
