package render

import (
	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// KubernetesVolumesRenderer renders Kubernetes volume components
var KubernetesVolumesRenderer = MakeReduce(
	VolumesRenderer(),
	PodToVolumeRenderer(),
	PVCToStorageClassRenderer(),
)

// VolumesRenderer returns renderer
func VolumesRenderer() Renderer {
	return volumesRenderer{}
}

// volumesRenderer is the renderer to render PVC & PV
type volumesRenderer struct{}

// Render renders the nodes
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

//PodToVolumeRenderer renders Pod and PVC resources
func PodToVolumeRenderer() Renderer {
	return podToVolumesRenderer{}
}

// VolumesRenderer is the renderer to render volumes
type podToVolumesRenderer struct{}

// Render renders the nodes
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

// PVCToStorageClassRenderer renders PVC and Storage class objects.
func PVCToStorageClassRenderer() Renderer {
	return pvcToStorageClassRenderer{}
}

// pvcToStorageClassRenderer is the renderer to render PVC & StorageClass
type pvcToStorageClassRenderer struct{}

// Render renders the nodes
func (v pvcToStorageClassRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for scID, scNode := range rpt.StorageClass.Nodes {
		storageClass, _ := scNode.Latest.Lookup(kubernetes.Name)
		for _, pvcNode := range rpt.PersistentVolumeClaim.Nodes {
			storageClassName, _ := pvcNode.Latest.Lookup(kubernetes.StorageClassName)
			if storageClassName == storageClass {
				scNode.Adjacency = scNode.Adjacency.Add(pvcNode.ID)
				scNode.Children = scNode.Children.Add(pvcNode)
			}
		}
		nodes[scID] = scNode
	}
	return Nodes{Nodes: nodes}
}
