package kubernetes

import (
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"

	"github.com/weaveworks/scope/report"
)

// StoragePool represent StoragePool interface
type StoragePool interface {
	Meta
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

// GetNode returns StoragePool as Node
func (p *storagePool) GetNode() report.Node {
	return p.MetaNode(report.MakeStoragePoolNodeID(p.UID())).WithLatests(map[string]string{
		NodeType:   "Storage Pool",
		APIVersion: p.APIVersion,
		Label:      p.GetLabels()["openebs.io/storagepoolclaim"],
	})
}
