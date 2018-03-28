package kubernetes

import (
	"encoding/json"

	"github.com/weaveworks/scope/report"
	storagev1 "k8s.io/api/storage/v1"
)

// StorageClass represent kubernetes StorageClass interface
type StorageClass interface {
	Meta
	GetNode(probeID string) report.Node
}

// storageClass represents kubernetes storage classes
type storageClass struct {
	*storagev1.StorageClass
	Meta
}

type apiversion struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

// NewStorageClass returns new Storage Class type
func NewStorageClass(p *storagev1.StorageClass) StorageClass {
	return &storageClass{StorageClass: p, Meta: meta{p.ObjectMeta}}
}

// GetNode returns StorageClass as Node
func (p *storageClass) GetNode(probeID string) report.Node {
	var version apiversion
	store := []byte(p.Annotations["kubectl.kubernetes.io/last-applied-configuration"])
	json.Unmarshal(store, &version)
	return p.MetaNode(report.MakeStorageClassNodeID(p.UID())).WithLatests(map[string]string{
		report.ControlProbeID: probeID,
		NodeType:              version.Kind,
		Name:                  p.GetName(),
		Provisioner:           p.Provisioner,
		APIVersion:            version.APIVersion,
		UID:                   p.UID(),
		ResourceVersion:       p.ResourceVersion,
		SelfLink:              p.SelfLink,
	})
}
