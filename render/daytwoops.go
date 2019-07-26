package render

import (
	"context"
	"strings"

	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/report"
)

// CSPCToNCSPRenderer is a Renderer which produces a renderable kubernetes CRD CSPC
var CSPCToNCSPRenderer = cspcToNCSPRenderer{}

// cspcToNCSPRenderer is a Renderer to render CSPC & NCSP nodes.
type cspcToNCSPRenderer struct{}

// Render renders the SPC & CSP nodes with adjacency.
// Here we are obtaining the spc name from csp and adjacency is created by matching it with spc name.
func (v cspcToNCSPRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for cspcID, cspcNode := range rpt.CStorPoolCluster.Nodes {
		cspcName, _ := cspcNode.Latest.Lookup(kubernetes.Name)
		cspcNamespace, _ := cspcNode.Latest.Lookup(kubernetes.Namespace)
		for ncspID, ncspNode := range rpt.NewTestCStorPool.Nodes {
			spcName, _ := ncspNode.Latest.Lookup(kubernetes.StoragePoolClaimName)
			ncspNamespace, _ := ncspNode.Latest.Lookup(kubernetes.Namespace)
			if cspcName == spcName && cspcNamespace == ncspNamespace {
				cspcNode.Adjacency = cspcNode.Adjacency.Add(ncspID)
				cspcNode.Children = cspcNode.Children.Add(ncspNode)
			}
		}
		nodes[cspcID] = cspcNode
	}
	return Nodes{Nodes: nodes}
}

// NCSPToBDRenderer is a renderer which produces a renderable CRD NCSP.
var NCSPToBDRenderer = ncspToBDRenderer{}

// ncspToBDRenderer is a Renderer to render NCSP & BD nodes.
type ncspToBDRenderer struct{}

func (n ncspToBDRenderer) Render(ctx context.Context, rpt report.Report) Nodes {
	nodes := make(report.Nodes)
	for ncspID, ncspNode := range rpt.NewTestCStorPool.Nodes {
		ncspNamespace, _ := ncspNode.Latest.Lookup(kubernetes.Namespace)
		blockDeviceList, _ := ncspNode.Latest.Lookup(kubernetes.BlockDeviceList)
		blockDevices := strings.Split(blockDeviceList, report.ScopeDelim)
		for _, blockDevice := range blockDevices {
			for bdID, bdNode := range rpt.BlockDevice.Nodes {
				bdName, _ := bdNode.Latest.Lookup(kubernetes.Name)
				bdNamespace, _ := bdNode.Latest.Lookup(kubernetes.Namespace)
				if bdName == blockDevice && bdNamespace == ncspNamespace {
					ncspNode.Adjacency = ncspNode.Adjacency.Add(bdID)
					ncspNode.Children = ncspNode.Children.Add(bdNode)
					break
				}
			}
		}
		nodes[ncspID] = ncspNode
	}
	return Nodes{Nodes: nodes}
}
