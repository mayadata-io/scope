package render

import (
	"strings"

	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// KubernetesVolumesRenderer is a Renderer which combines all Kubernetes
// volumes components such as stateful Pods, Persistent Volume, Persistent Volume Claim, Storage Class.
var KubernetesVolumesRenderer = MakeReduce(
	VolumesRenderer,
	PodToVolumeRenderer,
	PVCToStorageClassRenderer,
	PVToControllerRenderer,
	SPCToSPRenderer,
	SPToDiskRenderer,
	VolumeSnapshotRenderer,
)

// VolumesRenderer is a Renderer which produces a renderable kubernetes PV & PVC
// graph by merging the pods graph and the Persistent Volume topology.
var VolumesRenderer = volumesRenderer{}

// volumesRenderer is a Renderer to render PV & PVC nodes.
type volumesRenderer struct{}

// Render renders PV & PVC nodes along with adjacency
func (v volumesRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for id, n := range rpt.PersistentVolumeClaim.Nodes {
		volume, _ := n.Latest.Lookup(kubernetes.VolumeName)
		for pvNodeID, p := range rpt.PersistentVolume.Nodes {
			volumeName, _ := p.Latest.Lookup(kubernetes.Name)
			if volume == volumeName {
				n.Adjacency = n.Adjacency.Add(p.ID)
				n.Children = n.Children.Add(p)
			}
			nodes[pvNodeID] = p
		}
		nodes[id] = n
	}
	return Nodes{Nodes: nodes}
}

// PodToVolumeRenderer is a Renderer which produces a renderable kubernetes Pod
// graph by merging the pods graph and the Persistent Volume Claim topology.
// Pods having persistent volumes are rendered.
var PodToVolumeRenderer = podToVolumesRenderer{}

// VolumesRenderer is a Renderer to render Pods & PVCs.
type podToVolumesRenderer struct{}

// Render renders the Pod nodes having volumes adjacency.
func (v podToVolumesRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for podID, podNode := range rpt.Pod.Nodes {
		ClaimName, _ := podNode.Latest.Lookup(kubernetes.VolumeClaim)
		for _, pvcNode := range rpt.PersistentVolumeClaim.Nodes {
			pvcName, _ := pvcNode.Latest.Lookup(kubernetes.Name)
			if pvcName == ClaimName {
				podNode.Adjacency = podNode.Adjacency.Add(pvcNode.ID)
				podNode.Children = podNode.Children.Add(pvcNode)
			}
		}
		nodes[podID] = podNode
	}
	return Nodes{Nodes: nodes}
}

// PVCToStorageClassRenderer is a Renderer which produces a renderable kubernetes PVC
// & Storage class graph.
var PVCToStorageClassRenderer = pvcToStorageClassRenderer{}

// pvcToStorageClassRenderer is a Renderer to render PVC & StorageClass.
type pvcToStorageClassRenderer struct{}

// Render renders the PVC & Storage Class nodes with adjacency.
func (v pvcToStorageClassRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for scID, scNode := range rpt.StorageClass.Nodes {
		storageClass, _ := scNode.Latest.Lookup(kubernetes.Name)
		spcNameFromValue, _ := scNode.Latest.Lookup(kubernetes.Value)
		for _, pvcNode := range rpt.PersistentVolumeClaim.Nodes {
			storageClassName, _ := pvcNode.Latest.Lookup(kubernetes.StorageClassName)
			if storageClassName == storageClass {
				scNode.Adjacency = scNode.Adjacency.Add(pvcNode.ID)
				scNode.Children = scNode.Children.Add(pvcNode)
			}
		}

		// Expecting spcName from sc instead obtained a string i.e  - name: StoragePoolClaim value: "spcName" .
		// Hence we are spliting it to get spcName.
		if strings.Contains(spcNameFromValue, "\"") {
			storageValue := strings.Split(spcNameFromValue, "\"")
			spcNameFromValue = storageValue[1]
			for _, spcNode := range rpt.StoragePoolClaim.Nodes {
				spcName, _ := spcNode.Latest.Lookup(kubernetes.Name)
				if spcName == spcNameFromValue {
					scNode.Adjacency = scNode.Adjacency.Add(spcNode.ID)
					scNode.Children = scNode.Children.Add(spcNode)
				}
			}
		}
		nodes[scID] = scNode
	}
	return Nodes{Nodes: nodes}
}

//PVToControllerRenderer is a Renderer which produces a renderable kubernetes PVC
var PVToControllerRenderer = pvTocontrollerRenderer{}

//pvTocontrollerRenderer is a Renderer to render PV & Controller.
type pvTocontrollerRenderer struct{}

//Render renders the PV & Controller nodes with adjacency.
func (v pvTocontrollerRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for pvNodeID, p := range rpt.PersistentVolume.Nodes {
		volumeName, _ := p.Latest.Lookup(kubernetes.Name)
		for _, podNode := range rpt.Pod.Nodes {
			Controller, _ := podNode.Latest.Lookup(kubernetes.Name)
			if strings.Contains(Controller, "-ctrl-") {
				podName := strings.Split(Controller, "-ctrl-")
				Controller = podName[0]
				if volumeName == Controller {
					p.Adjacency = p.Adjacency.Add(podNode.ID)
					p.Children = p.Children.Add(podNode)
				}
			}
		}
		nodes[pvNodeID] = p
	}
	return Nodes{Nodes: nodes}
}

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
			spcNameFromLabel, _ := spNode.Latest.Lookup(kubernetes.Label)
			if spcName == spcNameFromLabel {
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
// First checking if sp and spc's are present if present we are obtaining the disk names from spc, adjacency is created by matching it with disk names.
func (v spToDiskRenderer) Render(rpt report.Report) Nodes {
	var disks []string
	nodes := make(report.Nodes)
	for spID, spNode := range rpt.StoragePool.Nodes {
		spcNameFromLabel, _ := spNode.Latest.Lookup(kubernetes.Label)
		for _, spcNode := range rpt.StoragePoolClaim.Nodes {
			spcName, _ := spcNode.Latest.Lookup(kubernetes.Name)
			if spcName == spcNameFromLabel {
				disk, _ := spcNode.Latest.Lookup(kubernetes.DiskList)
				if strings.Contains(disk, "/") {
					disks = strings.Split(disk, "/")
				} else {
					disks = []string{disk}
				}
				break
			}
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

// VolumeSnapshotRenderer is a renderer which produces a renderable Kubernetes Volume Snapshot and Volume Snapshot Data
var VolumeSnapshotRenderer = volumeSnapshotRenderer{}

// volumeSnapshotRenderer is a render to volume snapshot & volume snapshot data
type volumeSnapshotRenderer struct{}

// Render renders the volumeSnapshots & volumeSnapshotDatas with adjacency
// It checks for the volumeSnapshotData name in volumeSnapshot, adjacency is created by matching the volumeSnapshotData name.
func (v volumeSnapshotRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for volumeSnapshotID, volumeSnapshotNode := range rpt.VolumeSnapshot.Nodes {
		snapshotData, _ := volumeSnapshotNode.Latest.Lookup(kubernetes.SnapshotData)
		for volumeSnapshotDataID, volumeSnapshotDataNode := range rpt.VolumeSnapshotData.Nodes {
			snapshotDataName, _ := volumeSnapshotDataNode.Latest.Lookup(kubernetes.Name)
			if snapshotDataName == snapshotData {
				volumeSnapshotNode.Adjacency = volumeSnapshotNode.Adjacency.Add(volumeSnapshotDataNode.ID)
				volumeSnapshotNode.Children = volumeSnapshotNode.Children.Add(volumeSnapshotDataNode)
			}
			nodes[volumeSnapshotDataID] = volumeSnapshotDataNode
		}
		nodes[volumeSnapshotID] = volumeSnapshotNode
	}
	return Nodes{Nodes: nodes}
}
