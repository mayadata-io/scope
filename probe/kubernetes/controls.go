package kubernetes

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/weaveworks/scope/common/xfer"
	"github.com/weaveworks/scope/probe/controls"
	"github.com/weaveworks/scope/report"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Control IDs used by the kubernetes integration.
const (
	CloneVolumeSnapshot     = report.KubernetesCloneVolumeSnapshot
	CloneCsiVolumeSnapshot  = report.KubernetesCloneCsiVolumeSnapshot
	CreateVolumeSnapshot    = report.KubernetesCreateVolumeSnapshot
	GetLogs                 = report.KubernetesGetLogs
	Describe                = report.KubernetesDescribe
	DeletePod               = report.KubernetesDeletePod
	DeleteVolumeSnapshot    = report.KubernetesDeleteVolumeSnapshot
	DeleteCsiVolumeSnapshot = report.KubernetesDeleteCsiVolumeSnapshot
	ScaleUp                 = report.KubernetesScaleUp
	ScaleDown               = report.KubernetesScaleDown
)

// GroupName and version used by CRDs
const (
	SnapshotGroupName    = "volumesnapshot.external-storage.k8s.io"
	SnapshotVersion      = "v1"
	CsiSnapshotGroupName = "snapshot.storage.k8s.io"
	CsiSnapshotVersion   = "v1beta1"
	OpenEBSGroupName     = "openebs.io"
	OpenEBSVersion       = "v1alpha1"
)

var ctx = context.TODO()

// GetLogs is the control to get the logs for a kubernetes pod
func (r *Reporter) GetLogs(req xfer.Request, namespaceID, podID string, containerNames []string) xfer.Response {
	readCloser, err := r.client.GetLogs(ctx, namespaceID, podID, containerNames)
	if err != nil {
		return xfer.ResponseError(err)
	}

	readWriter := struct {
		io.Reader
		io.Writer
	}{
		readCloser,
		ioutil.Discard,
	}
	id, pipe, err := controls.NewPipeFromEnds(nil, readWriter, r.pipes, req.AppID)
	if err != nil {
		return xfer.ResponseError(err)
	}
	pipe.OnClose(func() {
		readCloser.Close()
	})
	return xfer.Response{
		Pipe: id,
	}
}

func (r *Reporter) describePod(req xfer.Request, namespaceID, podID string, _ []string) xfer.Response {
	return r.describe(req, namespaceID, podID, ResourceMap["Pod"], apimeta.RESTMapping{})
}

func (r *Reporter) describePVC(req xfer.Request, namespaceID, pvcID, _, _ string) xfer.Response {
	return r.describe(req, namespaceID, pvcID, ResourceMap["PersistentVolumeClaim"], apimeta.RESTMapping{})
}

func (r *Reporter) describeDeployment(req xfer.Request, namespaceID, deploymentID string) xfer.Response {
	return r.describe(req, namespaceID, deploymentID, ResourceMap["Deployment"], apimeta.RESTMapping{})
}

func (r *Reporter) describeService(req xfer.Request, namespaceID, serviceID string) xfer.Response {
	return r.describe(req, namespaceID, serviceID, ResourceMap["Service"], apimeta.RESTMapping{})
}

func (r *Reporter) describeCronJob(req xfer.Request, namespaceID, cronJobID string) xfer.Response {
	return r.describe(req, namespaceID, cronJobID, ResourceMap["CronJob"], apimeta.RESTMapping{})
}

func (r *Reporter) describePV(req xfer.Request, PVID string) xfer.Response {
	return r.describe(req, "", PVID, ResourceMap["PersistentVolume"], apimeta.RESTMapping{})
}

func (r *Reporter) describeDaemonSet(req xfer.Request, namespaceID, daemonSetID string) xfer.Response {
	return r.describe(req, namespaceID, daemonSetID, ResourceMap["DaemonSet"], apimeta.RESTMapping{})
}

func (r *Reporter) describeStatefulSet(req xfer.Request, namespaceID, statefulSetID string) xfer.Response {
	return r.describe(req, namespaceID, statefulSetID, ResourceMap["StatefulSet"], apimeta.RESTMapping{})
}

func (r *Reporter) describeStoragelass(req xfer.Request, storageClassID string) xfer.Response {
	return r.describe(req, "", storageClassID, ResourceMap["StorageClass"], apimeta.RESTMapping{})
}

func (r *Reporter) describeJob(req xfer.Request, namespaceID, jobID string) xfer.Response {
	return r.describe(req, namespaceID, jobID, ResourceMap["Job"], apimeta.RESTMapping{})
}

func (r *Reporter) describeVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, _, _ string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    SnapshotGroupName,
			Version:  SnapshotVersion,
			Resource: "volumesnapshots",
		},
	}
	return r.describe(req, namespaceID, volumeSnapshotID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeVolumeSnapshotData(req xfer.Request, volumeSnapshotID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    SnapshotGroupName,
			Version:  SnapshotVersion,
			Resource: "volumesnapshotdatas",
		},
	}
	return r.describe(req, "", volumeSnapshotID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeCsiVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, _, _, _ string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    CsiSnapshotGroupName,
			Version:  CsiSnapshotVersion,
			Resource: "volumesnapshots",
		},
	}
	return r.describe(req, namespaceID, volumeSnapshotID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeVolumeSnapshotClass(req xfer.Request, volumeSnapshotClassID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    CsiSnapshotGroupName,
			Version:  CsiSnapshotVersion,
			Resource: "volumesnapshotclasses",
		},
	}
	return r.describe(req, "", volumeSnapshotClassID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeVolumeSnapshotContent(req xfer.Request, volumeSnapshotContentID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    CsiSnapshotGroupName,
			Version:  CsiSnapshotVersion,
			Resource: "volumesnapshotcontents",
		},
	}
	return r.describe(req, "", volumeSnapshotContentID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeCV(req xfer.Request, namespaceID, CVID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "cstorvolumes",
		},
	}
	return r.describe(req, namespaceID, CVID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeCVR(req xfer.Request, namespaceID, CVRID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "cstorvolumereplicas",
		},
	}
	return r.describe(req, namespaceID, CVRID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeCSP(req xfer.Request, CSPID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "cstorpools",
		},
	}
	return r.describe(req, "", CSPID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeSPC(req xfer.Request, SPCID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "storagepoolclaims",
		},
	}
	return r.describe(req, "", SPCID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeDisk(req xfer.Request, diskID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "disks",
		},
	}
	return r.describe(req, "", diskID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeBlockDevice(req xfer.Request, namespaceID, blockDeviceID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "blockdevices",
		},
	}
	return r.describe(req, namespaceID, blockDeviceID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeBlockDeviceClaim(req xfer.Request, namespaceID, blockDeviceClaimID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "blockdeviceclaims",
		},
	}
	return r.describe(req, namespaceID, blockDeviceClaimID, schema.GroupKind{}, restMapping)
}
func (r *Reporter) describeCStorPoolCluster(req xfer.Request, namespaceID, cStorPoolClusterID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "cstorpoolclusters",
		},
	}
	return r.describe(req, namespaceID, cStorPoolClusterID, schema.GroupKind{}, restMapping)
}

func (r *Reporter) describeCStorPoolInstance(req xfer.Request, namespaceID, cStorPoolInstanceID string) xfer.Response {
	restMapping := apimeta.RESTMapping{
		Resource: schema.GroupVersionResource{
			Group:    OpenEBSGroupName,
			Version:  OpenEBSVersion,
			Resource: "cstorpoolinstances",
		},
	}
	return r.describe(req, namespaceID, cStorPoolInstanceID, schema.GroupKind{}, restMapping)
}

// GetLogs is the control to get the logs for a kubernetes pod
func (r *Reporter) describe(req xfer.Request, namespaceID, resourceID string, groupKind schema.GroupKind, restMapping apimeta.RESTMapping) xfer.Response {
	readCloser, err := r.client.Describe(namespaceID, resourceID, groupKind, restMapping)
	if err != nil {
		return xfer.ResponseError(err)
	}

	readWriter := struct {
		io.Reader
		io.Writer
	}{
		readCloser,
		ioutil.Discard,
	}
	id, pipe, err := controls.NewPipeFromEnds(nil, readWriter, r.pipes, req.AppID)
	if err != nil {
		return xfer.ResponseError(err)
	}
	pipe.OnClose(func() {
		readCloser.Close()
	})
	return xfer.Response{
		Pipe: id,
	}
}

func (r *Reporter) cloneVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity string) xfer.Response {
	err := r.client.CloneVolumeSnapshot(ctx, namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity)
	if err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{}
}

func (r *Reporter) cloneCsiVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity, driver string) xfer.Response {
	err := r.client.CloneCsiVolumeSnapshot(ctx, namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity, driver)
	if err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{}
}

func (r *Reporter) createVolumeSnapshot(req xfer.Request, namespaceID, persistentVolumeClaimID, capacity, driver string) xfer.Response {
	err := r.client.CreateVolumeSnapshot(ctx, namespaceID, persistentVolumeClaimID, capacity, driver)
	if err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{}
}

func (r *Reporter) deletePod(req xfer.Request, namespaceID, podID string, _ []string) xfer.Response {
	if err := r.client.DeletePod(ctx, namespaceID, podID); err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{
		RemovedNode: req.NodeID,
	}
}

func (r *Reporter) deleteVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, _, _ string) xfer.Response {
	if err := r.client.DeleteVolumeSnapshot(ctx, namespaceID, volumeSnapshotID); err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{
		RemovedNode: req.NodeID,
	}
}

func (r *Reporter) deleteCsiVolumeSnapshot(req xfer.Request, namespaceID, volumeSnapshotID, _, _, _ string) xfer.Response {
	if err := r.client.DeleteCsiVolumeSnapshot(ctx, namespaceID, volumeSnapshotID); err != nil {
		return xfer.ResponseError(err)
	}
	return xfer.Response{
		RemovedNode: req.NodeID,
	}
}

// Describe will parse the nodeID and return response according to the node (resource) type.
func (r *Reporter) Describe() func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		var f func(req xfer.Request) xfer.Response
		_, tag, ok := report.ParseNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		switch tag {
		case "<pod>":
			f = r.CapturePod(r.describePod)
		case "<service>":
			f = r.CaptureService(r.describeService)
		case "<cronjob>":
			f = r.CaptureCronJob(r.describeCronJob)
		case "<deployment>":
			f = r.CaptureDeployment(r.describeDeployment)
		case "<daemonset>":
			f = r.CaptureDaemonSet(r.describeDaemonSet)
		case "<persistent_volume>":
			f = r.CapturePersistentVolume(r.describePV)
		case "<persistent_volume_claim>":
			f = r.CapturePersistentVolumeClaim(r.describePVC)
		case "<storage_class>":
			f = r.CaptureStorageClass(r.describeStoragelass)
		case "<statefulset>":
			f = r.CaptureStatefulSet(r.describeStatefulSet)
		case "<volume_snapshot>":
			f = r.CaptureVolumeSnapshot(r.describeVolumeSnapshot)
		case "<volume_snapshot_data>":
			f = r.CaptureVolumeSnapshotData(r.describeVolumeSnapshotData)
		case "<job>":
			f = r.CaptureJob(r.describeJob)
		case "<cstor_volume>":
			f = r.CaptureCStorVolume(r.describeCV)
		case "<cstor_volume_replica>":
			f = r.CaptureCStorVolumeReplica(r.describeCVR)
		case "<cstor_pool>":
			f = r.CaptureCStorPool(r.describeCSP)
		case "<storage_pool_claim>":
			f = r.CaptureStoragePoolClaim(r.describeSPC)
		case "<disk>":
			f = r.CaptureDisk(r.describeDisk)
		case "<block_device>":
			f = r.CaptureBlockDevice(r.describeBlockDevice)
		case "<block_device_claim>":
			f = r.CaptureBlockDeviceClaim(r.describeBlockDeviceClaim)
		case "<cstor_pool_cluster>":
			f = r.CaptureCStorPoolCluster(r.describeCStorPoolCluster)
		case "<cstor_pool_instance>":
			f = r.CaptureCStorPoolInstance(r.describeCStorPoolInstance)
		case "<csi_volume_snapshot>":
			f = r.CaptureCsiVolumeSnapshot(r.describeCsiVolumeSnapshot)
		case "<volume_snapshot_class>":
			f = r.CaptureVolumeSnapshotClass(r.describeVolumeSnapshotClass)
		case "<volume_snapshot_content>":
			f = r.CaptureVolumeSnapshotContent(r.describeVolumeSnapshotContent)
		default:
			return xfer.ResponseErrorf("Node not found: %s", req.NodeID)
		}
		return f(req)
	}
}

// CapturePod is exported for testing
func (r *Reporter) CapturePod(f func(xfer.Request, string, string, []string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParsePodNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find pod by UID
		var pod Pod
		r.client.WalkPods(func(p Pod) error {
			if p.UID() == uid {
				pod = p
			}
			return nil
		})
		if pod == nil {
			return xfer.ResponseErrorf("Pod not found: %s", uid)
		}
		return f(req, pod.Namespace(), pod.Name(), pod.ContainerNames())
	}
}

// CaptureDeployment is exported for testing
func (r *Reporter) CaptureDeployment(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseDeploymentNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var deployment Deployment
		r.client.WalkDeployments(func(d Deployment) error {
			if d.UID() == uid {
				deployment = d
			}
			return nil
		})
		if deployment == nil {
			return xfer.ResponseErrorf("Deployment not found: %s", uid)
		}
		return f(req, deployment.Namespace(), deployment.Name())
	}
}

// CapturePersistentVolumeClaim will return name, namespace and capacity of PVC
func (r *Reporter) CapturePersistentVolumeClaim(f func(xfer.Request, string, string, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParsePersistentVolumeClaimNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find persistentVolumeClaim by UID
		var persistentVolumeClaim PersistentVolumeClaim
		r.client.WalkPersistentVolumeClaims(func(p PersistentVolumeClaim) error {
			if p.UID() == uid {
				persistentVolumeClaim = p
			}
			return nil
		})
		if persistentVolumeClaim == nil {
			return xfer.ResponseErrorf("Persistent volume claim not found: %s", uid)
		}

		// find provisioner from storage class
		var storageClass StorageClass
		r.client.WalkStorageClasses(func(p StorageClass) error {
			if p.Name() == persistentVolumeClaim.GetStorageClass() {
				storageClass = p
			}
			return nil
		})
		return f(req, persistentVolumeClaim.Namespace(), persistentVolumeClaim.Name(), persistentVolumeClaim.GetCapacity(), storageClass.GetProvisioner())
	}
}

// CaptureVolumeSnapshot will return name, pvc name, namespace and capacity of volume snapshot
func (r *Reporter) CaptureVolumeSnapshot(f func(xfer.Request, string, string, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseVolumeSnapshotNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find volume snapshot by UID
		var volumeSnapshot VolumeSnapshot
		r.client.WalkVolumeSnapshots(func(p VolumeSnapshot) error {
			if p.UID() == uid {
				volumeSnapshot = p
			}
			return nil
		})
		if volumeSnapshot == nil {
			return xfer.ResponseErrorf("Volume snapshot not found: %s", uid)
		}
		return f(req, volumeSnapshot.Namespace(), volumeSnapshot.Name(), volumeSnapshot.GetVolumeName(), volumeSnapshot.GetCapacity())
	}
}

// CaptureCsiVolumeSnapshot will return name, pvc name, namespace and capacity of volume snapshot
func (r *Reporter) CaptureCsiVolumeSnapshot(f func(xfer.Request, string, string, string, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCsiVolumeSnapshotNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find volume snapshot by UID
		var volumeSnapshot CsiVolumeSnapshot
		r.client.WalkCsiVolumeSnapshots(func(p CsiVolumeSnapshot) error {
			if p.UID() == uid {
				volumeSnapshot = p
			}
			return nil
		})
		if volumeSnapshot == nil {
			return xfer.ResponseErrorf("Volume snapshot not found: %s", uid)
		}
		return f(req, volumeSnapshot.Namespace(), volumeSnapshot.Name(), volumeSnapshot.GetVolumeName(), volumeSnapshot.GetCapacity(), volumeSnapshot.GetDriver())
	}
}

// CaptureVolumeSnapshotClass will return name
func (r *Reporter) CaptureVolumeSnapshotClass(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseVolumeSnapshotClassNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find volume snapshot class by UID
		var volumeSnapshotClass VolumeSnapshotClass
		r.client.WalkVolumeSnapshotClasses(func(p VolumeSnapshotClass) error {
			if p.UID() == uid {
				volumeSnapshotClass = p
			}
			return nil
		})
		if volumeSnapshotClass == nil {
			return xfer.ResponseErrorf("Volume snapshot class not found: %s", uid)
		}
		return f(req, volumeSnapshotClass.Name())
	}
}

// CaptureVolumeSnapshotClass will return name
func (r *Reporter) CaptureVolumeSnapshotContent(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseVolumeSnapshotContentNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find volume snapshot by UID
		var volumeSnapshotContent VolumeSnapshotContent
		r.client.WalkVolumeSnapshotContents(func(p VolumeSnapshotContent) error {
			if p.UID() == uid {
				volumeSnapshotContent = p
			}
			return nil
		})
		if volumeSnapshotContent == nil {
			return xfer.ResponseErrorf("Volume snapshot content not found: %s", uid)
		}
		return f(req, volumeSnapshotContent.Name())
	}
}

// CaptureService is exported for testing
func (r *Reporter) CaptureService(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseServiceNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var service Service
		r.client.WalkServices(func(s Service) error {
			if s.UID() == uid {
				service = s
			}
			return nil
		})
		if service == nil {
			return xfer.ResponseErrorf("Service not found: %s", uid)
		}
		return f(req, service.Namespace(), service.Name())
	}
}

// CaptureDaemonSet is exported for testing
func (r *Reporter) CaptureDaemonSet(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseDaemonSetNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var daemonSet DaemonSet
		r.client.WalkDaemonSets(func(d DaemonSet) error {
			if d.UID() == uid {
				daemonSet = d
			}
			return nil
		})
		if daemonSet == nil {
			return xfer.ResponseErrorf("Daemon Set not found: %s", uid)
		}
		return f(req, daemonSet.Namespace(), daemonSet.Name())
	}
}

// CaptureCronJob is exported for testing
func (r *Reporter) CaptureCronJob(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCronJobNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var cronJob CronJob
		r.client.WalkCronJobs(func(c CronJob) error {
			if c.UID() == uid {
				cronJob = c
			}
			return nil
		})
		if cronJob == nil {
			return xfer.ResponseErrorf("Cron Job not found: %s", uid)
		}
		return f(req, cronJob.Namespace(), cronJob.Name())
	}
}

// CaptureStatefulSet is exported for testing
func (r *Reporter) CaptureStatefulSet(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseStatefulSetNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var statefulSet StatefulSet
		r.client.WalkStatefulSets(func(s StatefulSet) error {
			if s.UID() == uid {
				statefulSet = s
			}
			return nil
		})
		if statefulSet == nil {
			return xfer.ResponseErrorf("Stateful Set not found: %s", uid)
		}
		return f(req, statefulSet.Namespace(), statefulSet.Name())
	}
}

// CaptureStorageClass is exported for testing
func (r *Reporter) CaptureStorageClass(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseStorageClassNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var storageClass StorageClass
		r.client.WalkStorageClasses(func(s StorageClass) error {
			if s.UID() == uid {
				storageClass = s
			}
			return nil
		})
		if storageClass == nil {
			return xfer.ResponseErrorf("StorageClass not found: %s", uid)
		}
		return f(req, storageClass.Name())
	}
}

// CapturePersistentVolume will return name of PV
func (r *Reporter) CapturePersistentVolume(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParsePersistentVolumeNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find persistentVolume by UID
		var persistentVolume PersistentVolume
		r.client.WalkPersistentVolumes(func(p PersistentVolume) error {
			if p.UID() == uid {
				persistentVolume = p
			}
			return nil
		})
		if persistentVolume == nil {
			return xfer.ResponseErrorf("Persistent volume  not found: %s", uid)
		}
		return f(req, persistentVolume.Name())
	}
}

// CaptureVolumeSnapshotData will return name of volume snapshot data
func (r *Reporter) CaptureVolumeSnapshotData(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseVolumeSnapshotDataNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find volume snapshotData by UID
		var volumeSnapshotData VolumeSnapshotData
		r.client.WalkVolumeSnapshotData(func(v VolumeSnapshotData) error {
			if v.UID() == uid {
				volumeSnapshotData = v
			}
			return nil
		})
		if volumeSnapshotData == nil {
			return xfer.ResponseErrorf("Volume snapshot data not found: %s", uid)
		}
		return f(req, volumeSnapshotData.Name())
	}
}

// CaptureJob is exported for testing
func (r *Reporter) CaptureJob(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseJobNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		var job Job
		r.client.WalkJobs(func(c Job) error {
			if c.UID() == uid {
				job = c
			}
			return nil
		})
		if job == nil {
			return xfer.ResponseErrorf("Job not found: %s", uid)
		}
		return f(req, job.Namespace(), job.Name())
	}
}

// CaptureCStorVolume will return name and namespace of cstor volume
func (r *Reporter) CaptureCStorVolume(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCStorVolumeNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find cv by UID
		var cstorVolume CStorVolume
		r.client.WalkCStorVolumes(func(c CStorVolume) error {
			if c.Name() == uid {
				cstorVolume = c
			}
			return nil
		})
		if cstorVolume == nil {
			return xfer.ResponseErrorf("CStor volume  not found: %s", uid)
		}
		return f(req, cstorVolume.Namespace(), cstorVolume.Name())
	}
}

// CaptureCStorVolumeReplica will return name and namespace of cstor volume replica
func (r *Reporter) CaptureCStorVolumeReplica(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCStorVolumeReplicaNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find cvr by UID
		var cstorVolumeReplica CStorVolumeReplica
		r.client.WalkCStorVolumeReplicas(func(c CStorVolumeReplica) error {
			if c.UID() == uid {
				cstorVolumeReplica = c
			}
			return nil
		})
		if cstorVolumeReplica == nil {
			return xfer.ResponseErrorf("CStor volume replica  not found: %s", uid)
		}
		return f(req, cstorVolumeReplica.Namespace(), cstorVolumeReplica.Name())
	}
}

// CaptureCStorPool will return name of cstor pool
func (r *Reporter) CaptureCStorPool(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCStorPoolNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find csp by UID
		var cstorPool CStorPool
		r.client.WalkCStorPools(func(c CStorPool) error {
			if c.UID() == uid {
				cstorPool = c
			}
			return nil
		})
		if cstorPool == nil {
			return xfer.ResponseErrorf("CStor pool not found: %s", uid)
		}
		return f(req, cstorPool.Name())
	}
}

// CaptureStoragePoolClaim will return name of spc
func (r *Reporter) CaptureStoragePoolClaim(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseStoragePoolClaimNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find spc by UID
		var spc StoragePoolClaim
		r.client.WalkStoragePoolClaims(func(s StoragePoolClaim) error {
			if s.UID() == uid {
				spc = s
			}
			return nil
		})
		if spc == nil {
			return xfer.ResponseErrorf("Storage pool claim not found: %s", uid)
		}
		return f(req, spc.Name())
	}
}

// CaptureDisk will return name of disk
func (r *Reporter) CaptureDisk(f func(xfer.Request, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseDiskNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find disk by UID
		var disk Disk
		r.client.WalkDisks(func(d Disk) error {
			if d.UID() == uid {
				disk = d
			}
			return nil
		})
		if disk == nil {
			return xfer.ResponseErrorf("Disk  not found: %s", uid)
		}
		return f(req, disk.Name())
	}
}

// CaptureBlockDevice will return name and namespace of block device
func (r *Reporter) CaptureBlockDevice(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseBlockDeviceNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find blockDevice by UID
		var blockDevice BlockDevice
		r.client.WalkBlockDevices(func(b BlockDevice) error {
			if b.UID() == uid {
				blockDevice = b
			}
			return nil
		})
		if blockDevice == nil {
			return xfer.ResponseErrorf("Block Device  not found: %s", uid)
		}
		return f(req, blockDevice.Namespace(), blockDevice.Name())
	}
}

// CaptureBlockDeviceClaim will return name and namespace of block device claim
func (r *Reporter) CaptureBlockDeviceClaim(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseBlockDeviceClaimNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find blockDeviceClaim by UID
		var blockDeviceClaim BlockDeviceClaim
		r.client.WalkBlockDeviceClaims(func(b BlockDeviceClaim) error {
			if b.UID() == uid {
				blockDeviceClaim = b
			}
			return nil
		})
		if blockDeviceClaim == nil {
			return xfer.ResponseErrorf("Block Device Claim  not found: %s", uid)
		}
		return f(req, blockDeviceClaim.Namespace(), blockDeviceClaim.Name())
	}
}

// CaptureCStorPoolCluster will return name and namespace of cstor pool cluster
func (r *Reporter) CaptureCStorPoolCluster(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCStorPoolClusterNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find cStorPoolCluster by UID
		var cStorPoolCluster CStorPoolCluster
		r.client.WalkCStorPoolClusters(func(c CStorPoolCluster) error {
			if c.UID() == uid {
				cStorPoolCluster = c
			}
			return nil
		})
		if cStorPoolCluster == nil {
			return xfer.ResponseErrorf("CStor Pool Cluster not found: %s", uid)
		}
		return f(req, cStorPoolCluster.Namespace(), cStorPoolCluster.Name())
	}
}

// CaptureCStorPoolInstance will return the name and namespace of new test cstor pool
func (r *Reporter) CaptureCStorPoolInstance(f func(xfer.Request, string, string) xfer.Response) func(xfer.Request) xfer.Response {
	return func(req xfer.Request) xfer.Response {
		uid, ok := report.ParseCStorPoolInstanceNodeID(req.NodeID)
		if !ok {
			return xfer.ResponseErrorf("Invalid ID: %s", req.NodeID)
		}
		// find cStorPoolInstance by UID
		var cStorPoolInstance CStorPoolInstance
		r.client.WalkCStorPoolInstances(func(c CStorPoolInstance) error {
			if c.UID() == uid {
				cStorPoolInstance = c
			}
			return nil
		})
		if cStorPoolInstance == nil {
			return xfer.ResponseErrorf("CStor Pool Cluster not found: %s", uid)
		}
		return f(req, cStorPoolInstance.Namespace(), cStorPoolInstance.Name())
	}
}

// ScaleUp is the control to scale up a deployment
func (r *Reporter) ScaleUp(req xfer.Request, namespace, id string) xfer.Response {
	return xfer.ResponseError(r.client.ScaleUp(ctx, namespace, id))
}

// ScaleDown is the control to scale up a deployment
func (r *Reporter) ScaleDown(req xfer.Request, namespace, id string) xfer.Response {
	return xfer.ResponseError(r.client.ScaleDown(ctx, namespace, id))
}

func (r *Reporter) registerControls() {
	controls := map[string]xfer.ControlHandlerFunc{
		CloneVolumeSnapshot:     r.CaptureVolumeSnapshot(r.cloneVolumeSnapshot),
		CloneCsiVolumeSnapshot:  r.CaptureCsiVolumeSnapshot(r.cloneCsiVolumeSnapshot),
		CreateVolumeSnapshot:    r.CapturePersistentVolumeClaim(r.createVolumeSnapshot),
		GetLogs:                 r.CapturePod(r.GetLogs),
		Describe:                r.Describe(),
		DeletePod:               r.CapturePod(r.deletePod),
		DeleteVolumeSnapshot:    r.CaptureVolumeSnapshot(r.deleteVolumeSnapshot),
		DeleteCsiVolumeSnapshot: r.CaptureCsiVolumeSnapshot(r.deleteCsiVolumeSnapshot),
		ScaleUp:                 r.CaptureDeployment(r.ScaleUp),
		ScaleDown:               r.CaptureDeployment(r.ScaleDown),
	}
	r.handlerRegistry.Batch(nil, controls)
}

func (r *Reporter) deregisterControls() {
	controls := []string{
		CloneVolumeSnapshot,
		CloneCsiVolumeSnapshot,
		CreateVolumeSnapshot,
		GetLogs,
		Describe,
		DeletePod,
		DeleteVolumeSnapshot,
		DeleteCsiVolumeSnapshot,
		ScaleUp,
		ScaleDown,
	}
	r.handlerRegistry.Batch(controls, nil)
}
