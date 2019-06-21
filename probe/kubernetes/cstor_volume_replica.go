package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorVolumeReplica interface
type CStorVolumeReplica interface {
	Meta
	GetNode(probeID string) report.Node
	GetCStorVolume() string
	GetCStorPool() string
	GetStatus() string
	GetNodeTagOnStatus(status string) string
}

// cStorVolume represents cStor Volume Replica CR
type cStorVolumeReplica struct {
	*mayav1alpha1.CStorVolumeReplica
	Meta
}

// NewCStorVolumeReplica returns fresh CStorVolumeReplica instance
func NewCStorVolumeReplica(p *mayav1alpha1.CStorVolumeReplica) CStorVolumeReplica {
	return &cStorVolumeReplica{CStorVolumeReplica: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorVolumeReplica) GetNode(probeID string) report.Node {
	var cStorPoolNodeID string
	status := p.GetStatus()
	latests := map[string]string{
		NodeType:              "CStor Volume Replica",
		APIVersion:            p.APIVersion,
		report.ControlProbeID: probeID,
	}
	if p.GetCStorVolume() != "" {
		latests[CStorVolumeName] = p.GetCStorVolume()
	}
	if p.GetCStorPool() != "" {
		cStorPoolNodeID = report.MakeCStorPoolNodeID(p.GetCStorPool())
	}

	if p.GetCStorPool() != "" {
		latests[CStorPoolUID] = p.GetCStorPool()
	}

	if status != "" {
		latests[Status] = status
	}

	return p.MetaNode(report.MakeCStorVolumeReplicaNodeID(p.UID())).
		WithLatests(latests).
		WithAdjacent(cStorPoolNodeID).
		WithNodeTag(p.GetNodeTagOnStatus(strings.ToLower(status))).
		WithLatestActiveControls(Describe)
}

func (p *cStorVolumeReplica) GetCStorVolume() string {
	cStorVolumeName := p.Labels()["cstorvolume.openebs.io/name"]
	return cStorVolumeName
}

func (p *cStorVolumeReplica) GetCStorPool() string {
	cStorVolumeUID := p.Labels()["cstorpool.openebs.io/uid"]
	return cStorVolumeUID
}

func (p *cStorVolumeReplica) GetStatus() string {
	status := p.Status.Phase
	return string(status)
}

func (p *cStorVolumeReplica) GetNodeTagOnStatus(status string) string {
	if status != "" {
		return CStorVolumeStatusMap[status]
	}
	return "unknown"
}
