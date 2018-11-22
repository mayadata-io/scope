package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"

	"strconv"

	"github.com/weaveworks/scope/report"
)

// StoragePoolClaim represent StoragePoolClaim interface
type StoragePoolClaim interface {
	Meta
	GetNode() report.Node
}

// storagePoolClaim represent the StoragePoolClaims CRD of Kubernetes
type storagePoolClaim struct {
	*mayav1alpha1.StoragePoolClaim
	Meta
}

// NewStoragePoolClaim returns new StoragePoolClaim type
func NewStoragePoolClaim(p *mayav1alpha1.StoragePoolClaim) StoragePoolClaim {
	return &storagePoolClaim{StoragePoolClaim: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns StoragePoolClaim as Node
func (p *storagePoolClaim) GetNode() report.Node {
	return p.MetaNode(report.MakeStoragePoolClaimNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:   "Storage Pool Claim",
		APIVersion: p.APIVersion,
		MaxPools:   strconv.Itoa(int(p.Spec.MaxPools)),
		Status:     p.Status.Phase,
	})
}
