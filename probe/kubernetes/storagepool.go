package kubernetes

import (
	"strings"

	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"

	"github.com/weaveworks/scope/report"
)

// Labels key for StoragePool
const (
	StoragePoolClaimLabel    = "openebs.io/storage-pool-claim"
	OldStoragePoolClaimLabel = "openebs.io/storagepoolclaim"
)

// StoragePool represent StoragePool interface
type StoragePool interface {
	Meta
	GetStoragePoolClaimName() string
	GetNode() report.Node
}

// storagePool represents the StoragePool CRD of Kubernetes
type storagePool struct {
	*mayav1alpha1.StoragePool
	Meta
}

// NewStoragePool returns new StoragePool type
func NewStoragePool(p *mayav1alpha1.StoragePool) StoragePool {
	return &storagePool{StoragePool: p, Meta: meta{p.ObjectMeta}}
}

func (p *storagePool) GetStoragePoolClaimName() string {
	labels := p.GetLabels()
	if storagePoolClaimName, ok := labels[StoragePoolClaimLabel]; ok {
		return storagePoolClaimName
	}
	return labels[OldStoragePoolClaimLabel]
}

// GetNode returns StoragePool as Node
func (p *storagePool) GetNode() report.Node {
	return p.MetaNode(report.MakeStoragePoolNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:             "Storage Pool",
		APIVersion:           p.APIVersion,
		DiskList:             strings.Join(p.Spec.Disks.DiskList, "/"),
		StoragePoolClaimName: p.GetStoragePoolClaimName(),
	})
}
