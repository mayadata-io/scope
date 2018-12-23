package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/weaveworks/scope/report"
)

// CStorPool interface
type CStorPool interface {
	Meta
	GetNode(probeID string) report.Node
	GetStatus() string
	GetNodeTagOnStatus(status string) string
	GetHost() string
	GetStoragePoolClaim() string
}

// cStorPool represents cStor Volume CSP
type cStorPool struct {
	*mayav1alpha1.CStorPool
	Meta
}

// NewCStorPool returns fresh CStorPool instance
func NewCStorPool(p *mayav1alpha1.CStorPool) CStorPool {
	return &cStorPool{CStorPool: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns updated node with CStor Volume details
func (p *cStorPool) GetNode(probeID string) report.Node {
	status := p.GetStatus()
	latests := map[string]string{
		NodeType:              "CStor Pool",
		APIVersion:            p.APIVersion,
		DiskList:              strings.Join(p.Spec.Disks.DiskList, "~p$"),
		HostName:              p.GetHost(),
		StoragePoolClaimName:  p.GetStoragePoolClaim(),
		report.ControlProbeID: probeID,
	}

	if status != "" {
		latests[Status] = status
	}
	return p.MetaNode(report.MakeCStorPoolNodeID(p.UID())).
		WithLatests(latests).
		WithNodeTag(p.GetNodeTagOnStatus(strings.ToLower(status))).
		WithLatestActiveControls(Describe)

}

func (p *cStorPool) GetStatus() string {
	status := p.Status.Phase
	return string(status)
}

func (p *cStorPool) GetNodeTagOnStatus(status string) string {
	if status != "" {
		return CStorVolumeStatusMap[status]
	}
	return "unknown"
}

func (p *cStorPool) GetHost() string {
	host := p.Labels()["kubernetes.io/hostname"]
	return string(host)
}

func (p *cStorPool) GetStoragePoolClaim() string {
	host := p.Labels()["openebs.io/storage-pool-claim"]
	return string(host)
}
