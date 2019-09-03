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
	BlockDeviceClaimToBlockDeviceRenderer,
	CSPCToCSPIRenderer,
	CSPIToBDRenderer,
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

// BlockDeviceClaimToBlockDeviceRenderer is a renderer which produces renderable k8s object of block device claim.
var BlockDeviceClaimToBlockDeviceRenderer blockDeviceClaimToBlockDeviceRenderer

// blockDeviceClaimToBlockDeviceRenderer is a renderer to render BDC and BD.
type blockDeviceClaimToBlockDeviceRenderer struct{}

func (b blockDeviceClaimToBlockDeviceRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for bdcID, bdcNode := range rpt.BlockDeviceClaim.Nodes {
		bdcNamespace, _ := bdcNode.Latest.Lookup(kubernetes.Namespace)
		bdNameFromBdc, _ := bdcNode.Latest.Lookup(kubernetes.BlockDeviceName)
		for bdID, bdNode := range rpt.BlockDevice.Nodes {
			bdName, _ := bdNode.Latest.Lookup(kubernetes.Name)
			bdNamespace, _ := bdNode.Latest.Lookup(kubernetes.Namespace)
			if bdName == bdNameFromBdc && bdNamespace == bdcNamespace {
				bdcNode.Adjacency = bdcNode.Adjacency.Add(bdID)
				bdcNode.Children = bdcNode.Children.Add(bdNode)
				break
			}
		}
		nodes[bdcID] = bdcNode
	}
	return Nodes{Nodes: nodes}
}

// CSPCToCSPIRenderer is a Renderer which produces a renderable kubernetes CRD CSPC
var CSPCToCSPIRenderer = cspcToCSPIRenderer{}

// cspcToCSPIRenderer is a Renderer to render CSPC & CSPI nodes.
type cspcToCSPIRenderer struct{}

// Render renders the SPC & CSP nodes with adjacency.
// Here we are obtaining the spc name from csp and adjacency is created by matching it with spc name.
func (v cspcToCSPIRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cspcID, cspcNode := range rpt.CStorPoolCluster.Nodes {
		cspcName, _ := cspcNode.Latest.Lookup(kubernetes.Name)
		cspcNamespace, _ := cspcNode.Latest.Lookup(kubernetes.Namespace)
		for cspiID, cspiNode := range rpt.CStorPoolInstance.Nodes {
			spcName, _ := cspiNode.Latest.Lookup(kubernetes.StoragePoolClaimName)
			cspiNamespace, _ := cspiNode.Latest.Lookup(kubernetes.Namespace)
			if cspcName == spcName && cspcNamespace == cspiNamespace {
				cspcNode.Adjacency = cspcNode.Adjacency.Add(cspiID)
				cspcNode.Children = cspcNode.Children.Add(cspiNode)
			}
		}
		nodes[cspcID] = cspcNode
	}
	return Nodes{Nodes: nodes}
}

// CSPIToBDRenderer is a renderer which produces a renderable CRD CSPI.
var CSPIToBDRenderer = cspiToBDRenderer{}

// cspiToBDRenderer is a Renderer to render CSPI & BD nodes.
type cspiToBDRenderer struct{}

func (n cspiToBDRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cspiID, cspiNode := range rpt.CStorPoolInstance.Nodes {
		cspiNamespace, _ := cspiNode.Latest.Lookup(kubernetes.Namespace)
		blockDeviceList, _ := cspiNode.Latest.Lookup(kubernetes.BlockDeviceList)
		blockDevices := strings.Split(blockDeviceList, report.ScopeDelim)
		for _, blockDevice := range blockDevices {
			for bdID, bdNode := range rpt.BlockDevice.Nodes {
				bdName, _ := bdNode.Latest.Lookup(kubernetes.Name)
				bdNamespace, _ := bdNode.Latest.Lookup(kubernetes.Namespace)
				if bdName == blockDevice && bdNamespace == cspiNamespace {
					cspiNode.Adjacency = cspiNode.Adjacency.Add(bdID)
					cspiNode.Children = cspiNode.Children.Add(bdNode)
					break
				}
			}
		}
		nodes[cspiID] = cspiNode
	}
	return Nodes{Nodes: nodes}
}
