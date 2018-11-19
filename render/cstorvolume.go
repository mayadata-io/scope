package render

import (
	"github.com/weaveworks/scope/report"
)

//CStorVolumeRenderer is a Renderer which produces a renderable openebs CV.
var CStorVolumeRenderer = cStorVolumeRenderer{}

//cStorVolumeRenderer is a Renderer to render CStor Volumes.
type cStorVolumeRenderer struct{}

//Render renders the CV.
func (v cStorVolumeRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cvID, cvNode := range rpt.CStorVolume.Nodes {
		nodes[cvID] = cvNode
	}
	return Nodes{Nodes: nodes}
}

//CStorVolumeReplicaRenderer is a Renderer which produces a renderable openebs CVR.
var CStorVolumeReplicaRenderer = cStorVolumeReplicaRenderer{}

//cStorVolumeReplicaRenderer is a Renderer to render CStor Volume Replica.
type cStorVolumeReplicaRenderer struct{}

//Render renders the CVR.
func (v cStorVolumeReplicaRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cvrID, cvrNode := range rpt.CStorVolumeReplica.Nodes {
		nodes[cvrID] = cvrNode
	}
	return Nodes{Nodes: nodes}
}
