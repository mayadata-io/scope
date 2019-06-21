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
	CSPToBdOrDiskRenderer,
	BlockDeviceToDiskRenderer,
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

// CSPToBdOrDiskRenderer is a Renderer which produces a renderable kubernetes CRD Disk
var CSPToBdOrDiskRenderer = cspToBdOrDiskRenderer{}

// cspToBdOrDiskRenderer is a Renderer to render CSP & Disk .
type cspToBdOrDiskRenderer struct{}

// Render renders the CSP & Disk nodes with adjacency.
func (v cspToBdOrDiskRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cspID, cspNode := range rpt.CStorPool.Nodes {
		cspDiskPaths, _ := cspNode.Latest.Lookup(kubernetes.DiskList)
		cspBlockDevice, _ := cspNode.Latest.Lookup(kubernetes.BlockDeviceList)
		cspHostname, _ := cspNode.Latest.Lookup(kubernetes.HostName)

		if cspDiskPaths != "" {
			diskList := strings.Split(cspDiskPaths, report.ScopeDelim)
			for diskNodeID, diskNode := range rpt.Disk.Nodes {
				diskHostname, _ := diskNode.Latest.Lookup(kubernetes.HostName)
				diskPaths, _ := diskNode.Latest.Lookup(kubernetes.DiskList)
				diskPathList := strings.Split(diskPaths, report.ScopeDelim)
				for _, cspDiskPath := range diskList {
					for _, diskPath := range diskPathList {
						if (cspDiskPath == diskPath) && (cspHostname == diskHostname) {
							cspNode.Adjacency = cspNode.Adjacency.Add(diskNodeID)
							cspNode.Children = cspNode.Children.Add(diskNode)
						}
					}
				}
				nodes[diskNodeID] = diskNode
			}
		}

		if cspBlockDevice != "" {
			cspBlockDeviceList := strings.Split(cspBlockDevice, report.ScopeDelim)
			for blockDeviceID, blockDeviceNode := range rpt.BlockDevice.Nodes {
				blockDeviceName, _ := blockDeviceNode.Latest.Lookup(kubernetes.Name)
				for _, cspBlockDeviceName := range cspBlockDeviceList {
					if blockDeviceName == cspBlockDeviceName {
						cspNode.Adjacency = cspNode.Adjacency.Add(blockDeviceID)
						cspNode.Children = cspNode.Children.Add(blockDeviceNode)
					}
				}
			}
		}
		nodes[cspID] = cspNode
	}
	return Nodes{Nodes: nodes}
}

// BlockDeviceToDiskRenderer is a renderer which produces a renderable kubernetes block device and disk.
var BlockDeviceToDiskRenderer = blockDeviceToDiskRenderer{}

// blockDeviceToDiskRenderer is a renderer to render block device and disk.
type blockDeviceToDiskRenderer struct{}

func (b blockDeviceToDiskRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for blockDeviceNodeID, blockDeviceNode := range rpt.BlockDevice.Nodes {
		blockDevicePath, _ := blockDeviceNode.Latest.Lookup(kubernetes.Path)
		blockDeviceHost, _ := blockDeviceNode.Latest.Lookup(kubernetes.HostName)
		blockDeviceName, _ := blockDeviceNode.Latest.Lookup(kubernetes.Name)

		for diskNodeID, diskNode := range rpt.Disk.Nodes {
			diskPath, _ := diskNode.Latest.Lookup(kubernetes.Path)
			diskHost, _ := diskNode.Latest.Lookup(kubernetes.HostName)
			diskName, _ := diskNode.Latest.Lookup(kubernetes.Name)

			if blockDevicePath == diskPath && blockDeviceHost == diskHost &&
				strings.Split(blockDeviceName, "-")[1] == strings.Split(diskName, "-")[1] {
				blockDeviceNode.Adjacency = blockDeviceNode.Adjacency.Add(diskNodeID)
				blockDeviceNode.Children = blockDeviceNode.Children.Add(diskNode)
			}

			nodes[diskNodeID] = diskNode
		}
		nodes[blockDeviceNodeID] = blockDeviceNode
	}
	return Nodes{Nodes: nodes}
}
