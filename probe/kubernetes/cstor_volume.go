package kubernetes

import (
	"strconv"
	"strings"

	openebsapiv1 "github.com/openebs/api/pkg/apis/cstor/v1"
	"github.com/weaveworks/scope/report"
)

// CStorVolume interface
type CStorVolume interface {
	Meta
	GetNode(probeID string) report.Node
	GetPersistentVolumeName() string
	GetStatus() string
	GetNodeTagOnStatus(status string) string
	GetConsistencyFactor() string
	GetReplicationFactor() string
	GetIQNInfo() string
}

// cStorVolume represents cStor Volume CR
type cStorVolume struct {
	*openebsapiv1.CStorVolume
	Meta
}

// NewCStorVolume returns fresh CStorVolume instance
func NewCStorVolume(p *openebsapiv1.CStorVolume) CStorVolume {
	return &cStorVolume{CStorVolume: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorVolume) GetNode(probeID string) report.Node {
	status := p.GetStatus()
	latests := map[string]string{
		NodeType:              "CStor Volume",
		APIVersion:            p.APIVersion,
		report.ControlProbeID: probeID,
	}

	if p.GetPersistentVolumeName() != "" {
		latests[VolumeName] = p.GetPersistentVolumeName()
	}

	if status != "" {
		latests[Status] = status
	}

	if p.GetConsistencyFactor() != "" {
		latests[CStorVolumeConsistencyFactor] = p.GetConsistencyFactor()
	}

	if p.GetReplicationFactor() != "" {
		latests[CStorVolumeReplicationFactor] = p.GetReplicationFactor()
	}

	if p.GetIQNInfo() != "" {
		latests[CStorVolumeIQN] = p.GetIQNInfo()
	}

	return p.MetaNode(report.MakeCStorVolumeNodeID(p.Name())).
		WithLatests(latests).
		WithNodeTag(p.GetNodeTagOnStatus(strings.ToLower(status))).
		WithLatestActiveControls(Describe)
}

func (p *cStorVolume) GetPersistentVolumeName() string {
	persistentVolumeName := p.Labels()["openebs.io/persistent-volume"]
	return persistentVolumeName
}

func (p *cStorVolume) GetStatus() string {
	status := p.Status.Phase
	return string(status)
}

func (p *cStorVolume) GetNodeTagOnStatus(status string) string {
	if status != "" {
		return CStorVolumeStatusMap[status]
	}
	return "unknown"
}

func (p *cStorVolume) GetConsistencyFactor() string {
	return strconv.Itoa(p.Spec.ConsistencyFactor)
}

func (p *cStorVolume) GetReplicationFactor() string {
	return strconv.Itoa(p.Spec.ReplicationFactor)
}

func (p *cStorVolume) GetIQNInfo() string {
	return p.Spec.Iqn
}
