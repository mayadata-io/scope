package render

import (
	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

//CStorVolumeRenderer is a Renderer which produces a renderable openebs CV.
var CStorVolumeRenderer = cStorVolumeRenderer{}

//cStorVolumeRenderer is a Renderer to render CStor Volumes.
type cStorVolumeRenderer struct{}

//Render renders the CV.
func (v cStorVolumeRenderer) Render(rpt report.Report) Nodes {
	cStorNodes := make(report.Nodes)

	for cvID, cvNode := range rpt.CStorVolume.Nodes {
		cStorNodes[cvID] = cvNode
	}

	for cspID, cspNode := range rpt.CStorPool.Nodes {
		cStorNodes[cspID] = cspNode
	}

	for cvrID, cvrNode := range rpt.CStorVolumeReplica.Nodes {
		cStorVolume, _ := cvrNode.Latest.Lookup(kubernetes.CStorVolumeName)
		cStorVolumeNodeID := report.MakeCStorVolumeNodeID(cStorVolume)

		if cvNode, ok := cStorNodes[cStorVolumeNodeID]; ok {
			cvNode.Adjacency = cvNode.Adjacency.Add(cvrID)
			cvNode.Children = cvNode.Children.Add(cvrNode)
			cStorNodes[cStorVolumeNodeID] = cvNode
		}

		cStorPoolUID, _ := cvrNode.Latest.Lookup(kubernetes.CStorPoolUID)
		cStorPoolNodeID := report.MakeCStorPoolNodeID(cStorPoolUID)
		if cStorPoolNode, ok := cStorNodes[cStorPoolNodeID]; ok {
			cvrNode.Children = cvrNode.Children.Add(cStorPoolNode)
		}

		cStorNodes[cvrID] = cvrNode
	}

	return Nodes{Nodes: cStorNodes}
}
