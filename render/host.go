package render

import (
	"strings"

	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// HostRenderer is a Renderer which produces a renderable host
// graph from the host topology.
//
// not memoised
var HostRenderer = MakeReduce(
	CustomRenderer{RenderFunc: nodes2Hosts, Renderer: ProcessRenderer},
	CustomRenderer{RenderFunc: nodes2Hosts, Renderer: ContainerRenderer},
	CustomRenderer{RenderFunc: nodes2Hosts, Renderer: ContainerImageRenderer},
	CustomRenderer{RenderFunc: nodes2Hosts, Renderer: PodRenderer},
	MapEndpoints(endpoint2Host, report.Host),
	NdmRenderer,
)

// nodes2Hosts maps any Nodes to host Nodes.
//
// If this function is given a node without a hostname
// (including other pseudo nodes), it will drop the node.
//
// Otherwise, this function will produce nodes with the correct ID
// format for a host, but without any Major or Minor labels.  It does
// not have enough info to do that, and the resulting graph must be
// merged with a host graph to get that info.
func nodes2Hosts(nodes Nodes) Nodes {
	ret := newJoinResults(nil)

	for _, n := range nodes.Nodes {
		if n.Topology == Pseudo {
			continue // Don't propagate pseudo nodes - we do this in endpoints2Hosts
		}
		isImage := n.Topology == report.ContainerImage
		hostIDs, _ := n.Parents.Lookup(report.Host)
		for _, id := range hostIDs {
			if isImage {
				// We need to treat image nodes specially because they
				// aggregate adjacencies of containers across multiple
				// hosts, and hence mapping these adjacencies to host
				// adjacencies would produce edges that aren't present
				// in reality.
				ret.addUnmappedChild(n, id, report.Host)
			} else {
				ret.addChild(n, id, report.Host)
			}
		}
	}
	return ret.result(nodes)
}

func endpoint2Host(n report.Node) string {
	if hostNodeID, ok := n.Latest.Lookup(report.HostNodeID); ok {
		return hostNodeID
	}
	return ""
}

// NdmRenderer is a Renderer which produces a renderable NDM Disk
var NdmRenderer = ndmRenderer{}

//ndmRenderer is a Renderer to render disk nodes.
type ndmRenderer struct{}

//Render will render the Disk and add Adjacency in host nodes i.e Host->Disk.
func (v ndmRenderer) Render(rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for hostNodeID, h := range rpt.Host.Nodes {
		hostName, _ := h.Latest.Lookup(report.HostNodeID)
		if strings.Contains(hostName, ";") {
			hostid := strings.Split(hostName, ";")
			hostName = hostid[0]
			for diskNode, d := range rpt.Disk.Nodes {
				Label, _ := d.Latest.Lookup(kubernetes.HostName)
				if strings.ToLower(hostName) == Label {
					h.Adjacency = h.Adjacency.Add(d.ID)
					h.Children = h.Children.Add(d)
				}
				nodes[diskNode] = d
			}
			nodes[hostNodeID] = h
		}
	}
	return Nodes{Nodes: nodes}
}
