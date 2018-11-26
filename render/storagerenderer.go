package render

import (
	"context"
	"strings"

	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// KubernetesStorageRenderer is a Renderer which combines all Kubernetes
// storage components such as CstorPools, storage pool claims and disks.
var KubernetesStorageRenderer = MakeReduce(
	SPCToCSPRenderer,
	CSPToDiskRenderer,
)

// SPCToCSPRenderer is a Renderer which produces a renderable kubernetes CRD SPC
var SPCToCSPRenderer = spcToCSPRenderer{}

// spcToCSPRenderer is a Renderer to render SPC & CSP nodes.
type spcToCSPRenderer struct{}

// Render renders the SPC & CSP nodes with adjacency.
// Here we are obtaining the spc name from csp and adjacency is created by matching it with spc name.
func (v spcToCSPRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for spcID, spcNode := range rpt.StoragePoolClaim.Nodes {
		spcName, _ := spcNode.Latest.Lookup(kubernetes.Name)
		for cspID, cspNode := range rpt.CStorPool.Nodes {
			storagePoolCaimName, _ := cspNode.Latest.Lookup(kubernetes.StoragePoolClaimName)
			if storagePoolCaimName == spcName {
				spcNode.Adjacency = spcNode.Adjacency.Add(cspID)
				spcNode.Children = spcNode.Children.Add(cspNode)
			}
			nodes[spcID] = spcNode
		}

	}
	return Nodes{Nodes: nodes}
}

// CSPToDiskRenderer is a Renderer which produces a renderable kubernetes CRD Disk
var CSPToDiskRenderer = cspToDiskRenderer{}

// cspToDiskRenderer is a Renderer to render CSP & Disk .
type cspToDiskRenderer struct{}

// Render renders the CSP & Disk nodes with adjacency.
func (v cspToDiskRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	// var disks []string

	nodes := make(report.Nodes)
	for cspID, cspNode := range rpt.CStorPool.Nodes {
		cspDiskPaths, _ := cspNode.Latest.Lookup(kubernetes.DiskList)
		cspHostname, _ := cspNode.Latest.Lookup(kubernetes.HostName)

		diskList := strings.Split(cspDiskPaths, "~p$")

		for diskID, diskNode := range rpt.Disk.Nodes {
			diskHostname, _ := diskNode.Latest.Lookup(kubernetes.HostName)
			diskPaths, _ := diskNode.Latest.Lookup(kubernetes.DiskList)
			diskPathList := strings.Split(diskPaths, "~p$")
			for _, cspDiskPath := range diskList {
				for _, diskPath := range diskPathList {
					if (cspDiskPath == diskPath) && (cspHostname == diskHostname) {
						cspNode.Adjacency = cspNode.Adjacency.Add(diskID)
						cspNode.Children = cspNode.Children.Add(diskNode)
					}
				}
			}
			nodes[diskID] = diskNode
		}
		nodes[cspID] = cspNode
	}

	return Nodes{Nodes: nodes}
}
