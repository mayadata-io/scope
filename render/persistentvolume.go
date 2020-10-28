package render

import (
	"context"
	"strings"

	"github.com/weaveworks/scope/report"
)

// KubernetesVolumesRenderer is a Renderer which combines all Kubernetes
// volumes components such as stateful Pods, Persistent Volume, Persistent Volume Claim, Storage Class.
var KubernetesVolumesRenderer = MakeReduce(
	CStorVolumeRenderer,
	VolumesRenderer,
	PodToVolumeRenderer,
	PVCToStorageClassRenderer,
	PVToControllerRenderer,
	VolumeSnapshotRenderer,
	CSPToBdOrDiskRenderer,
	BlockDeviceClaimToBlockDeviceRenderer,
	CSPIToBDRenderer,
	BlockDeviceToDiskRenderer,
	CSiVolumeSnapshotRenderer,
	VolumeSnapshotClassRenderer,
	MakeFilter(
		func(n report.Node) bool {
			value, _ := n.Latest.Lookup(report.KubernetesVolumePod)
			if value == "true" {
				return true
			}
			return false
		},
		PodRenderer,
	),
)

// VolumesRenderer is a Renderer which produces a renderable kubernetes PV & PVC
// graph by merging the pods graph and the Persistent Volume topology.
var VolumesRenderer = volumesRenderer{}

// volumesRenderer is a Renderer to render PV & PVC nodes.
type volumesRenderer struct{}

// Render renders PV & PVC nodes along with adjacency
func (v volumesRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for id, n := range rpt.PersistentVolumeClaim.Nodes {
		volume, _ := n.Latest.Lookup(report.KubernetesVolumeName)
		for _, p := range rpt.PersistentVolume.Nodes {
			volumeName, _ := p.Latest.Lookup(report.KubernetesName)
			if volume == volumeName {
				n.Adjacency = n.Adjacency.Add(p.ID)
				n.Children = n.Children.Add(p)
			}
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
func (v podToVolumesRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for podID, podNode := range rpt.Pod.Nodes {
		claimNames, found := podNode.Latest.Lookup(report.KubernetesVolumeClaim)
		if !found {
			continue
		}
		podNamespace, _ := podNode.Latest.Lookup(report.KubernetesNamespace)
		claimNameList := strings.Split(claimNames, report.ScopeDelim)
		for _, ClaimName := range claimNameList {
			for _, pvcNode := range rpt.PersistentVolumeClaim.Nodes {
				pvcName, _ := pvcNode.Latest.Lookup(report.KubernetesName)
				pvcNamespace, _ := pvcNode.Latest.Lookup(report.KubernetesNamespace)
				if (pvcName == ClaimName) && (podNamespace == pvcNamespace) {
					podNode.Adjacency = podNode.Adjacency.Add(pvcNode.ID)
					podNode.Children = podNode.Children.Add(pvcNode)
					break
				}
			}
		}
		if found {
			nodes[podID] = podNode
		}
	}
	return Nodes{Nodes: nodes}
}

// PVCToStorageClassRenderer is a Renderer which produces a renderable kubernetes PVC
// & Storage class graph.
var PVCToStorageClassRenderer = pvcToStorageClassRenderer{}

// pvcToStorageClassRenderer is a Renderer to render PVC & StorageClass.
type pvcToStorageClassRenderer struct{}

// Render renders the PVC & Storage Class nodes with adjacency.
func (v pvcToStorageClassRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for scID, scNode := range rpt.StorageClass.Nodes {
		storageClass, _ := scNode.Latest.Lookup(report.KubernetesName)
		spcNameFromValue, _ := scNode.Latest.Lookup(report.KubernetesValue)
		for _, pvcNode := range rpt.PersistentVolumeClaim.Nodes {
			storageClassName, _ := pvcNode.Latest.Lookup(report.KubernetesStorageClassName)
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
				spcName, _ := spcNode.Latest.Lookup(report.KubernetesName)
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
var PVToControllerRenderer = pvToControllerRenderer{}

//pvTocontrollerRenderer is a Renderer to render PV & Controller.
type pvToControllerRenderer struct{}

//Render renders the PV & Controller nodes with adjacency.
func (v pvToControllerRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for pvNodeID, p := range rpt.PersistentVolume.Nodes {
		volumeName, _ := p.Latest.Lookup(report.KubernetesName)
		volumeClaimName, _ := p.Latest.Lookup(report.KubernetesVolumeClaim)
		for _, podNode := range rpt.Pod.Nodes {
			podVolumeName, _ := podNode.Latest.Lookup(report.KubernetesVolumeName)
			if volumeName == podVolumeName {
				p.Adjacency = p.Adjacency.Add(podNode.ID)
				p.Children = p.Children.Add(podNode)
			}
		}

		for _, volumeSnapshotNode := range rpt.VolumeSnapshot.Nodes {
			snapshotPVName, _ := volumeSnapshotNode.Latest.Lookup(report.KubernetesVolumeName)
			if volumeName == snapshotPVName {
				p.Adjacency = p.Adjacency.Add(volumeSnapshotNode.ID)
				p.Children = p.Children.Add(volumeSnapshotNode)
			}
		}

		for _, csiVolumeSnapshotNode := range rpt.CsiVolumeSnapshot.Nodes {
			snapshotPVCName, _ := csiVolumeSnapshotNode.Latest.Lookup(report.KubernetesVolumeClaim)
			if volumeClaimName == snapshotPVCName {
				p.Adjacency = p.Adjacency.Add(csiVolumeSnapshotNode.ID)
				p.Children = p.Children.Add(csiVolumeSnapshotNode)
			}
		}

		for cvID, cvNode := range rpt.CStorVolume.Nodes {
			pvName, _ := cvNode.Latest.Lookup(report.KubernetesVolumeName)
			if pvName == volumeName {
				p.Adjacency = p.Adjacency.Add(cvID)
				p.Children = p.Children.Add(cvNode)
			}
		}

		_, casOk := p.Latest.Lookup(report.KubernetesCASType)
		bdcNameFromPV, bdcOk := p.Latest.Lookup(report.KubernetesBlockDeviceClaimName)
		if casOk && bdcOk {
			for bdcID, bdcNode := range rpt.BlockDeviceClaim.Nodes {
				bdcName, _ := bdcNode.Latest.Lookup(report.KubernetesName)
				if bdcName == bdcNameFromPV {
					p.Adjacency = p.Adjacency.Add(bdcID)
					p.Children = p.Children.Add(bdcNode)
					break
				}
			}
		}

		if p.ID != "" {
			nodes[pvNodeID] = p
		}
	}
	return Nodes{Nodes: nodes}
}

// VolumeSnapshotRenderer is a renderer which produces a renderable Kubernetes Volume Snapshot and Volume Snapshot Data
var VolumeSnapshotRenderer = volumeSnapshotRenderer{}

// volumeSnapshotRenderer is a render to volume snapshot & volume snapshot data
type volumeSnapshotRenderer struct{}

// Render renders the volumeSnapshots & volumeSnapshotData with adjacency
// It checks for the volumeSnapshotData name in volumeSnapshot, adjacency is created by matching the volumeSnapshotData name.
func (v volumeSnapshotRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for volumeSnapshotID, volumeSnapshotNode := range rpt.VolumeSnapshot.Nodes {
		volumeSnapshotName, _ := volumeSnapshotNode.Latest.Lookup(report.KubernetesName)
		snapshotData, _ := volumeSnapshotNode.Latest.Lookup(report.KubernetesSnapshotData)
		for volumeSnapshotDataID, volumeSnapshotDataNode := range rpt.VolumeSnapshotData.Nodes {
			snapshotDataName, _ := volumeSnapshotDataNode.Latest.Lookup(report.KubernetesName)
			if snapshotDataName == snapshotData {
				volumeSnapshotNode.Adjacency = volumeSnapshotNode.Adjacency.Add(volumeSnapshotDataNode.ID)
				volumeSnapshotNode.Children = volumeSnapshotNode.Children.Add(volumeSnapshotDataNode)
			}
			nodes[volumeSnapshotDataID] = volumeSnapshotDataNode
		}

		for persistentVolumeClaimID, persistentVolumeClaimNode := range rpt.PersistentVolumeClaim.Nodes {
			vsName, ok := persistentVolumeClaimNode.Latest.Lookup(report.KubernetesVolumeSnapshotName)
			if !ok {
				continue
			}
			if vsName == volumeSnapshotName {
				volumeSnapshotNode.Adjacency = volumeSnapshotNode.Adjacency.Add(persistentVolumeClaimID)
				volumeSnapshotNode.Children = volumeSnapshotNode.Children.Add(persistentVolumeClaimNode)
			}
		}
		nodes[volumeSnapshotID] = volumeSnapshotNode
	}
	return Nodes{Nodes: nodes}
}

// CSiVolumeSnapshotRenderer is a renderer which produces a renderable Kubernetes Volume Snapshot and Volume Snapshot Data
var CSiVolumeSnapshotRenderer = csiVolumeSnapshotRenderer{}

// csiVolumeSnapshotRenderer is a render to volume snapshot & volume snapshot data
type csiVolumeSnapshotRenderer struct{}

// Render renders the volumeSnapshots & volumeSnapshotData with adjacency
// It checks for the volumeSnapshotData name in volumeSnapshot, adjacency is created by matching the volumeSnapshotData name.
func (v csiVolumeSnapshotRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for volumeSnapshotID, volumeSnapshotNode := range rpt.CsiVolumeSnapshot.Nodes {
		volumeSnapshotName, _ := volumeSnapshotNode.Latest.Lookup(report.KubernetesName)
		volumeSnapshotNamespace, _ := volumeSnapshotNode.Latest.Lookup(report.KubernetesNamespace)
		for persistentVolumeClaimID, persistentVolumeClaimNode := range rpt.PersistentVolumeClaim.Nodes {
			vsName, ok := persistentVolumeClaimNode.Latest.Lookup(report.KubernetesVolumeSnapshotName)
			if !ok {
				continue
			}
			if vsName == volumeSnapshotName {
				volumeSnapshotNode.Adjacency = volumeSnapshotNode.Adjacency.Add(persistentVolumeClaimID)
				volumeSnapshotNode.Children = volumeSnapshotNode.Children.Add(persistentVolumeClaimNode)
			}
		}

		for volumeSnapshotContentID, volumeSnapshotContentNode := range rpt.VolumeSnapshotContent.Nodes {
			vsName, _ := volumeSnapshotContentNode.Latest.Lookup(report.KubernetesVolumeSnapshotName)
			vsNamespace, _ := volumeSnapshotContentNode.Latest.Lookup(report.KubernetesVolumeSnapshotNamespace)
			if vsName == volumeSnapshotName && vsNamespace == volumeSnapshotNamespace {
				volumeSnapshotNode.Adjacency = volumeSnapshotNode.Adjacency.Add(volumeSnapshotContentNode.ID)
				volumeSnapshotNode.Children = volumeSnapshotNode.Children.Add(volumeSnapshotContentNode)
			}
			nodes[volumeSnapshotContentID] = volumeSnapshotContentNode
		}

		nodes[volumeSnapshotID] = volumeSnapshotNode
	}
	return Nodes{Nodes: nodes}
}

// VolumeSnapshotClassRenderer is a renderer which produces a renderable Kubernetes Volume Snapshot class and Volume Snapshot
var VolumeSnapshotClassRenderer = volumeSnapshotClassRenderer{}

// csiVolumeSnapshotRenderer is a render to volume snapshot class & volume snapshot
type volumeSnapshotClassRenderer struct{}

// Render renders the csiVolumeSnapshots & volumeSnapshotClass with adjacency
// It checks for the volumeSnapshotClass name in volumeSnapshot, adjacency is created by matching the volumeSnapshotClass name.
func (v volumeSnapshotClassRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for volumeSnapshotClassID, volumeSnapshotClassNode := range rpt.VolumeSnapshotClass.Nodes {
		volumeSnapshotClassName, _ := volumeSnapshotClassNode.Latest.Lookup(report.KubernetesName)
		for _, csiVolumeSnapshotNode := range rpt.CsiVolumeSnapshot.Nodes {
			snsapshotClassName, _ := csiVolumeSnapshotNode.Latest.Lookup(report.KubernetesSnapshotClass)
			if volumeSnapshotClassName == snsapshotClassName {
				volumeSnapshotClassNode.Adjacency = volumeSnapshotClassNode.Adjacency.Add(csiVolumeSnapshotNode.ID)
				volumeSnapshotClassNode.Children = volumeSnapshotClassNode.Children.Add(csiVolumeSnapshotNode)
			}
		}
		nodes[volumeSnapshotClassID] = volumeSnapshotClassNode
	}
	return Nodes{Nodes: nodes}
}
