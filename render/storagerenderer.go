package render

import (
	"strings"

	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// KubernetesStorageRenderer is a Renderer which combines all Kubernetes
// storage components such as storage pools, storage pool claims and disks.
var KubernetesStorageRenderer = MakeReduce(
	PVCToStorageClassRenderer,
	SPCToSPRenderer,
	SPToDiskRenderer,
)

// SPCToSPRenderer is a Renderer which produces a renderable kubernetes CRD SPC
var SPCToSPRenderer = spcToSpRenderer{}

// spcToSpRenderer is a Renderer to render SPC & SP nodes.
type spcToSpRenderer struct{}

// Render renders the SPC & SP nodes with adjacency.
// Here we are obtaining the spc name from sp and adjacency is created by matching it with spc name.
func (v spcToSpRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for spcID, spcNode := range rpt.StoragePoolClaim.Nodes {
		spcName, _ := spcNode.Latest.Lookup(kubernetes.Name)
		for _, spNode := range rpt.StoragePool.Nodes {
			spcNameFromLatest, _ := spNode.Latest.Lookup(kubernetes.StoragePoolClaimName)
			if spcName == spcNameFromLatest {
				spcNode.Adjacency = spcNode.Adjacency.Add(spNode.ID)
				spcNode.Children = spcNode.Children.Add(spNode)
			}
		}
		nodes[spcID] = spcNode
	}
	return Nodes{Nodes: nodes}
}

// SPToDiskRenderer is a Renderer which produces a renderable kubernetes CRD Disk
var SPToDiskRenderer = spToDiskRenderer{}

// spToDiskRenderer is a Renderer to render SP & Disk .
type spToDiskRenderer struct{}

// Render renders the SP & Disk nodes with adjacency.
func (v spToDiskRenderer) Render(rpt report.Report) Nodes {
	var disks []string
	nodes := make(report.Nodes)
	for spID, spNode := range rpt.StoragePool.Nodes {
		disk, _ := spNode.Latest.Lookup(kubernetes.DiskList)
		if strings.Contains(disk, "/") {
			disks = strings.Split(disk, "/")
		} else {
			disks = []string{disk}
		}

		diskList := make(map[string]string)
		for _, disk := range disks {
			diskList[disk] = disk
		}

		for diskID, diskNode := range rpt.Disk.Nodes {
			diskName, _ := diskNode.Latest.Lookup(kubernetes.Name)
			if diskName == diskList[diskName] {
				spNode.Adjacency = spNode.Adjacency.Add(diskNode.ID)
				spNode.Children = spNode.Children.Add(diskNode)
			}
			nodes[diskID] = diskNode
		}
		nodes[spID] = spNode
	}
	return Nodes{Nodes: nodes}
}
