package render

import (
	"context"

	"github.com/weaveworks/scope/report"
)

// KubernetesStorageRenderer is a Renderer which combines all Kubernetes
// storage components such as storage pools, storage pool claims and disks.
var KubernetesStorageRenderer = MakeReduce(
	SPCToSPRenderer,
	SPToDiskRenderer,
)

// SPCToSPRenderer is a Renderer which produces a renderable kubernetes CRD SPC
var SPCToSPRenderer = spcToSpRenderer{}

// spcToSpRenderer is a Renderer to render SPC & SP nodes.
type spcToSpRenderer struct{}

// Render renders the SPC & SP nodes with adjacency.
// Here we are obtaining the spc name from sp and adjacency is created by matching it with spc name.
func (v spcToSpRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for spcID, spcNode := range rpt.StoragePoolClaim.Nodes {
		nodes[spcID] = spcNode
	}
	return Nodes{Nodes: nodes}
}

// SPToDiskRenderer is a Renderer which produces a renderable kubernetes CRD Disk
var SPToDiskRenderer = spToDiskRenderer{}

// spToDiskRenderer is a Renderer to render SP & Disk .
type spToDiskRenderer struct{}

// Render renders the SP & Disk nodes with adjacency.
func (v spToDiskRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for diskID, diskNode := range rpt.Disk.Nodes {
		nodes[diskID] = diskNode
	}
	return Nodes{Nodes: nodes}
}
