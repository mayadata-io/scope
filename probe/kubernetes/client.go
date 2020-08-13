package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/weaveworks/common/backoff"

	snapshotv1 "github.com/openebs/k8s-snapshot-client/snapshot/pkg/apis/volumesnapshot/v1"
	snapshot "github.com/openebs/k8s-snapshot-client/snapshot/pkg/client/clientset/versioned"

	csisnapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	csisnapshot "github.com/kubernetes-csi/external-snapshotter/v2/pkg/client/clientset/versioned"
	mayav1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	mayaclient "github.com/openebs/maya/pkg/client/clientset/versioned"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	apiappsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	apibatchv1 "k8s.io/api/batch/v1"
	apibatchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubectldescribe "k8s.io/kubectl/pkg/describe"
)

// Client keeps track of running kubernetes pods and services
type Client interface {
	Stop()
	WalkPods(f func(Pod) error) error
	WalkServices(f func(Service) error) error
	WalkDeployments(f func(Deployment) error) error
	WalkDaemonSets(f func(DaemonSet) error) error
	WalkStatefulSets(f func(StatefulSet) error) error
	WalkCronJobs(f func(CronJob) error) error
	WalkNamespaces(f func(NamespaceResource) error) error
	WalkPersistentVolumes(f func(PersistentVolume) error) error
	WalkPersistentVolumeClaims(f func(PersistentVolumeClaim) error) error
	WalkStorageClasses(f func(StorageClass) error) error
	WalkVolumeSnapshots(f func(VolumeSnapshot) error) error
	WalkVolumeSnapshotData(f func(VolumeSnapshotData) error) error
	WalkJobs(f func(Job) error) error
	WalkDisks(f func(Disk) error) error
	WalkStoragePoolClaims(f func(StoragePoolClaim) error) error
	WalkCStorVolumes(f func(CStorVolume) error) error
	WalkCStorVolumeReplicas(f func(CStorVolumeReplica) error) error
	WalkCStorPools(f func(CStorPool) error) error
	WalkBlockDevices(f func(BlockDevice) error) error
	WalkBlockDeviceClaims(f func(BlockDeviceClaim) error) error
	WalkCStorPoolClusters(f func(CStorPoolCluster) error) error
	WalkCStorPoolInstances(f func(CStorPoolInstance) error) error
	WalkCsiVolumeSnapshots(f func(CsiVolumeSnapshot) error) error
	WalkVolumeSnapshotClasses(f func(VolumeSnapshotClass) error) error
	WalkVolumeSnapshotContents(f func(VolumeSnapshotContent) error) error
	WatchPods(f func(Event, Pod))

	CloneVolumeSnapshot(namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity string, ctx context.Context) error
	CloneCsiVolumeSnapshot(namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity, driver string, ctx context.Context) error
	CreateVolumeSnapshot(namespaceID, persistentVolumeClaimID, capacity, driver string, ctx context.Context) error
	GetLogs(namespaceID, podID string, containerNames []string, ctx context.Context) (io.ReadCloser, error)
	Describe(namespaceID, resourceID string, groupKind schema.GroupKind, restMapping apimeta.RESTMapping) (io.ReadCloser, error)
	DeletePod(namespaceID, podID string, ctx context.Context) error
	DeleteVolumeSnapshot(namespaceID, volumeSnapshotID string, ctx context.Context) error
	DeleteCsiVolumeSnapshot(namespaceID, volumeSnapshotID string, ctx context.Context) error
	ScaleUp(namespaceID, id string, ctx context.Context) error
	ScaleDown(namespaceID, id string, ctx context.Context) error
}

// ResourceMap is the mapping of resource and their GroupKind
var ResourceMap = map[string]schema.GroupKind{
	"Pod":                   {Group: apiv1.GroupName, Kind: "Pod"},
	"Service":               {Group: apiv1.GroupName, Kind: "Service"},
	"Deployment":            {Group: apiappsv1.GroupName, Kind: "Deployment"},
	"DaemonSet":             {Group: apiappsv1.GroupName, Kind: "DaemonSet"},
	"StatefulSet":           {Group: apiappsv1.GroupName, Kind: "StatefulSet"},
	"Job":                   {Group: apibatchv1.GroupName, Kind: "Job"},
	"CronJob":               {Group: apibatchv1.GroupName, Kind: "CronJob"},
	"Node":                  {Group: apiv1.GroupName, Kind: "Node"},
	"PersistentVolume":      {Group: apiv1.GroupName, Kind: "PersistentVolume"},
	"PersistentVolumeClaim": {Group: apiv1.GroupName, Kind: "PersistentVolumeClaim"},
	"StorageClass":          {Group: storagev1.GroupName, Kind: "StorageClass"},
}

var csiDriverMap = map[string]bool{
	"diskplugin.csi.alibabacloud.com":          true,
	"nasplugin.csi.alibabacloud.com":           true,
	"ossplugin.csi.alibabacloud.com":           true,
	"arstor.csi.huayun.io":                     true,
	"ebs.csi.aws.com":                          true,
	"efs.csi.aws.com":                          true,
	"fsx.csi.aws.com":                          true,
	"disk.csi.azure.com":                       true,
	"file.csi.azure.com":                       true,
	"csi.block.bigtera.com":                    true,
	"csi.fs.bigtera.com":                       true,
	"cephfs.csi.ceph.com":                      true,
	"rbd.csi.ceph.com":                         true,
	"csi.chubaofs.com":                         true,
	"cinder.csi.openstack.org":                 true,
	"csi.cloudscale.ch":                        true,
	"csi-infiblock-plugin":                     true,
	"csi-infifs-plugin":                        true,
	"dsp.csi.daterainc.io":                     true,
	"csi-isilon.dellemc.com":                   true,
	"csi-powermax.dellemc.com":                 true,
	"csi-powerstore.dellemc.com":               true,
	"csi-unity.dellemc.com":                    true,
	"csi-vxflexos.dellemc.com":                 true,
	"csi-xtremio.dellemc.com":                  true,
	"org.democratic-csi.v1.0":                  true,
	"org.democratic-csi.v1.1":                  true,
	"org.democratic-csi.v1.2":                  true,
	"dcx.csi.diamanti.com":                     true,
	"dobs.csi.digitalocean.com":                true,
	"csi.drivescale.com":                       true,
	"v0.2.ember-csi.io":                        true,
	"v0.3.ember-csi.io":                        true,
	"v1.0.ember-csi.io":                        true,
	"pd.csi.storage.gke.io":                    true,
	"com.google.csi.filestore":                 true,
	"gcs.csi.ofek.dev":                         true,
	"org.gluster.glusterfs":                    true,
	"org.gluster.glustervirtblock":             true,
	"com.hammerspace.csi":                      true,
	"io.hedvig.csi":                            true,
	"csi.hetzner.cloud":                        true,
	"com.hitachi.hspc.csi":                     true,
	"csi.hpe.com":                              true,
	"csi.huawei.com":                           true,
	"eu.zetanova.csi.hyperv":                   true,
	"block.csi.ibm.com":                        true,
	"spectrumscale.csi.ibm.com":                true,
	"vpc.block.csi.ibm.io":                     true,
	"infinibox-csi-driver":                     true,
	"csi-instorage":                            true,
	"pmem-csi.intel.com":                       true,
	"csi.juicefs.com":                          true,
	"org.kadalu.gluster":                       true,
	"linodebs.csi.linode.com":                  true,
	"io.drbd.linstor-csi":                      true,
	"driver.longhorn.io":                       true,
	"csi-macrosan":                             true,
	"manila.csi.openstack.org":                 true,
	"com.mapr.csi-kdf":                         true,
	"com.tuxera.csi.moosefs":                   true,
	"csi.trident.netapp.io":                    true,
	"nexentastor-csi-driver.nexenta.com":       true,
	"nexentastor-block-csi-driver.nexenta.com": true,
	"com.nutanix.csi":                          true,
	"cstor.csi.openebs.io":                     true,
	"csi-opensdsplugin":                        true,
	"com.open-e.joviandss.csi":                 true,
	"pxd.openstorage.org":                      true,
	"pure-csi":                                 true,
	"disk.csi.qingcloud.com":                   true,
	"csi-neonsan":                              true,
	"quobyte-csi":                              true,
	"robin":                                    true,
	"csi-sandstone-plugin":                     true,
	"eds.csi.sangfor.com":                      true,
	"seaweedfs-csi-driver":                     true,
	"secrets-store.csi.k8s.io":                 true,
	"csi-smtx-plugin":                          true,
	"csi.spdk.io":                              true,
	"storageos":                                true,
	"com.tencent.cloud.csi.cbs":                true,
	"com.tencent.cloud.csi.cfs":                true,
	"com.tencent.cloud.csi.cosfs":              true,
	"topolvm.cybozu.com":                       true,
	"csi.vastdata.com":                         true,
	"csi.block.xsky.com":                       true,
	"csi.fs.xsky.com":                          true,
	"secrets.csi.kubevault.com":                true,
	"csi.vsphere.vmware.com":                   true,
	"csi.weka.io":                              true,
	"yandex.csi.flant.com":                     true,
	"csi.zadara.com":                           true,
}

type client struct {
	quit                       chan struct{}
	client                     *kubernetes.Clientset
	snapshotClient             *snapshot.Clientset
	mayaClient                 *mayaclient.Clientset
	csiSnapshotClient          *csisnapshot.Clientset
	podStore                   cache.Store
	serviceStore               cache.Store
	deploymentStore            cache.Store
	daemonSetStore             cache.Store
	statefulSetStore           cache.Store
	jobStore                   cache.Store
	cronJobStore               cache.Store
	nodeStore                  cache.Store
	namespaceStore             cache.Store
	persistentVolumeStore      cache.Store
	persistentVolumeClaimStore cache.Store
	storageClassStore          cache.Store
	volumeSnapshotStore        cache.Store
	volumeSnapshotDataStore    cache.Store
	diskStore                  cache.Store
	storagePoolClaimStore      cache.Store
	cStorvolumeStore           cache.Store
	cStorvolumeReplicaStore    cache.Store
	cStorPoolStore             cache.Store
	blockDeviceStore           cache.Store
	blockDeviceClaimStore      cache.Store
	cStorPoolClusterStore      cache.Store
	cStorPoolInstanceStore     cache.Store
	csiVolumeSnapshotStore     cache.Store
	volumeSnapshotClassStore   cache.Store
	volumeSnapshotContentStore cache.Store

	podWatchesMutex sync.Mutex
	podWatches      []func(Event, Pod)
}

// ClientConfig establishes the configuration for the kubernetes client
type ClientConfig struct {
	CertificateAuthority string
	ClientCertificate    string
	ClientKey            string
	Cluster              string
	Context              string
	Insecure             bool
	Kubeconfig           string
	Password             string
	Server               string
	Token                string
	User                 string
	Username             string
}

// NewClient returns a usable Client. Don't forget to Stop it.
func NewClient(config ClientConfig) (Client, error) {
	var restConfig *rest.Config
	if config.Server == "" && config.Kubeconfig == "" {
		// If no API server address or kubeconfig was provided, assume we are running
		// inside a pod. Try to connect to the API server through its
		// Service environment variables, using the default Service
		// Account Token.
		var err error
		if restConfig, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	} else {
		var err error
		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
			&clientcmd.ConfigOverrides{
				AuthInfo: clientcmdapi.AuthInfo{
					ClientCertificate: config.ClientCertificate,
					ClientKey:         config.ClientKey,
					Token:             config.Token,
					Username:          config.Username,
					Password:          config.Password,
				},
				ClusterInfo: clientcmdapi.Cluster{
					Server:                config.Server,
					InsecureSkipTLSVerify: config.Insecure,
					CertificateAuthority:  config.CertificateAuthority,
				},
				Context: clientcmdapi.Context{
					Cluster:  config.Cluster,
					AuthInfo: config.User,
				},
				CurrentContext: config.Context,
			},
		).ClientConfig()
		if err != nil {
			return nil, err
		}
	}
	log.Infof("kubernetes: targeting api server %s", restConfig.Host)

	c, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	sc, err := snapshot.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	mc, err := mayaclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	csc, err := csisnapshot.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	result := &client{
		quit:              make(chan struct{}),
		client:            c,
		snapshotClient:    sc,
		mayaClient:        mc,
		csiSnapshotClient: csc,
	}

	result.podStore = NewEventStore(result.triggerPodWatches, cache.MetaNamespaceKeyFunc)
	result.runReflectorUntil("pods", result.podStore)

	result.serviceStore = result.setupStore("services")
	result.nodeStore = result.setupStore("nodes")
	result.namespaceStore = result.setupStore("namespaces")
	result.deploymentStore = result.setupStore("deployments")
	result.daemonSetStore = result.setupStore("daemonsets")
	result.jobStore = result.setupStore("jobs")
	result.statefulSetStore = result.setupStore("statefulsets")
	result.cronJobStore = result.setupStore("cronjobs")
	result.persistentVolumeStore = result.setupStore("persistentvolumes")
	result.persistentVolumeClaimStore = result.setupStore("persistentvolumeclaims")
	result.storageClassStore = result.setupStore("storageclasses")
	result.volumeSnapshotStore = result.setupStore("volumesnapshots")
	result.volumeSnapshotDataStore = result.setupStore("volumesnapshotdatas")
	result.diskStore = result.setupStore("disks")
	result.storagePoolClaimStore = result.setupStore("storagepoolclaims")
	result.cStorvolumeStore = result.setupStore("cstorvolumes")
	result.cStorvolumeReplicaStore = result.setupStore("cstorvolumereplicas")
	result.cStorPoolStore = result.setupStore("cstorpools")
	result.blockDeviceStore = result.setupStore("blockdevices")
	result.blockDeviceClaimStore = result.setupStore("blockdeviceclaims")
	result.cStorPoolClusterStore = result.setupStore("cstorpoolclusters")
	result.cStorPoolInstanceStore = result.setupStore("cstorpoolinstances")
	result.csiVolumeSnapshotStore = result.setupStore("csivolumesnapshots")
	result.volumeSnapshotClassStore = result.setupStore("volumesnapshotclasses")
	result.volumeSnapshotContentStore = result.setupStore("volumesnapshotcontents")

	return result, nil
}

func (c *client) isResourceSupported(groupVersion schema.GroupVersion, resource string) (bool, error) {
	resourceList, err := c.client.Discovery().ServerResourcesForGroupVersion(groupVersion.String())
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	for _, v := range resourceList.APIResources {
		if v.Name == resource {
			return true, nil
		}
	}

	return false, nil
}

func (c *client) setupStore(resource string) cache.Store {
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	c.runReflectorUntil(resource, store)
	return store
}

func (c *client) clientAndType(resource string) (rest.Interface, interface{}, error) {
	switch resource {
	case "pods":
		return c.client.CoreV1().RESTClient(), &apiv1.Pod{}, nil
	case "services":
		return c.client.CoreV1().RESTClient(), &apiv1.Service{}, nil
	case "nodes":
		return c.client.CoreV1().RESTClient(), &apiv1.Node{}, nil
	case "namespaces":
		return c.client.CoreV1().RESTClient(), &apiv1.Namespace{}, nil
	case "persistentvolumes":
		return c.client.CoreV1().RESTClient(), &apiv1.PersistentVolume{}, nil
	case "persistentvolumeclaims":
		return c.client.CoreV1().RESTClient(), &apiv1.PersistentVolumeClaim{}, nil
	case "storageclasses":
		return c.client.StorageV1().RESTClient(), &storagev1.StorageClass{}, nil
	case "deployments":
		return c.client.AppsV1().RESTClient(), &apiappsv1.Deployment{}, nil
	case "daemonsets":
		return c.client.AppsV1().RESTClient(), &apiappsv1.DaemonSet{}, nil
	case "jobs":
		return c.client.BatchV1().RESTClient(), &apibatchv1.Job{}, nil
	case "statefulsets":
		return c.client.AppsV1().RESTClient(), &apiappsv1.StatefulSet{}, nil
	case "volumesnapshots":
		return c.snapshotClient.VolumesnapshotV1().RESTClient(), &snapshotv1.VolumeSnapshot{}, nil
	case "volumesnapshotdatas":
		return c.snapshotClient.VolumesnapshotV1().RESTClient(), &snapshotv1.VolumeSnapshotData{}, nil
	case "disks":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.Disk{}, nil
	case "storagepoolclaims":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.StoragePoolClaim{}, nil
	case "cstorvolumes":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.CStorVolume{}, nil
	case "cstorvolumereplicas":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.CStorVolumeReplica{}, nil
	case "cstorpools":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.CStorPool{}, nil
	case "blockdevices":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.BlockDevice{}, nil
	case "blockdeviceclaims":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.BlockDeviceClaim{}, nil
	case "cstorpoolclusters":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.CStorPoolCluster{}, nil
	case "cstorpoolinstances":
		return c.mayaClient.OpenebsV1alpha1().RESTClient(), &mayav1alpha1.CStorPoolInstance{}, nil
	case "cronjobs":
		return c.client.BatchV1beta1().RESTClient(), &apibatchv1beta1.CronJob{}, nil
	case "csivolumesnapshots":
		return c.csiSnapshotClient.SnapshotV1beta1().RESTClient(), &csisnapshotv1beta1.VolumeSnapshot{}, nil
	case "volumesnapshotclasses":
		return c.csiSnapshotClient.SnapshotV1beta1().RESTClient(), &csisnapshotv1beta1.VolumeSnapshotClass{}, nil
	case "volumesnapshotcontents":
		return c.csiSnapshotClient.SnapshotV1beta1().RESTClient(), &csisnapshotv1beta1.VolumeSnapshotContent{}, nil
	}
	return nil, nil, fmt.Errorf("Invalid resource: %v", resource)
}

// runReflectorUntil runs cache.Reflector#ListAndWatch in an endless loop, after checking that the resource is supported by kubernetes.
// Errors are logged and retried with exponential backoff.
func (c *client) runReflectorUntil(resource string, store cache.Store) {
	var r *cache.Reflector
	listAndWatch := func() (bool, error) {
		if r == nil {
			kclient, itemType, err := c.clientAndType(resource)
			if err != nil {
				return false, err
			}

			if resource == "csivolumesnapshots" {
				resource = "volumesnapshots"
			}

			ok, err := c.isResourceSupported(kclient.APIVersion(), resource)
			if err != nil {
				return false, err
			}
			if !ok {
				log.Infof("%v are not supported by this Kubernetes version", resource)
				return true, nil
			}
			lw := cache.NewListWatchFromClient(kclient, resource, metav1.NamespaceAll, fields.Everything())
			r = cache.NewReflector(lw, itemType, store, 0)
		}

		select {
		case <-c.quit:
			return true, nil
		default:
			err := r.ListAndWatch(c.quit)
			return false, err
		}
	}
	bo := backoff.New(listAndWatch, fmt.Sprintf("Kubernetes reflector (%s)", resource))
	bo.SetMaxBackoff(5 * time.Minute)
	go bo.Start()
}

func (c *client) WatchPods(f func(Event, Pod)) {
	c.podWatchesMutex.Lock()
	defer c.podWatchesMutex.Unlock()
	c.podWatches = append(c.podWatches, f)
}

func (c *client) triggerPodWatches(e Event, pod interface{}) {
	c.podWatchesMutex.Lock()
	defer c.podWatchesMutex.Unlock()
	for _, watch := range c.podWatches {
		watch(e, NewPod(pod.(*apiv1.Pod)))
	}
}

func (c *client) WalkPods(f func(Pod) error) error {
	for _, m := range c.podStore.List() {
		pod := m.(*apiv1.Pod)
		if err := f(NewPod(pod)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkPersistentVolumes(f func(PersistentVolume) error) error {
	for _, m := range c.persistentVolumeStore.List() {
		pv := m.(*apiv1.PersistentVolume)
		if err := f(NewPersistentVolume(pv)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkPersistentVolumeClaims(f func(PersistentVolumeClaim) error) error {
	for _, m := range c.persistentVolumeClaimStore.List() {
		pvc := m.(*apiv1.PersistentVolumeClaim)
		if err := f(NewPersistentVolumeClaim(pvc)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkStorageClasses(f func(StorageClass) error) error {
	for _, m := range c.storageClassStore.List() {
		sc := m.(*storagev1.StorageClass)
		if err := f(NewStorageClass(sc)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkServices(f func(Service) error) error {
	for _, m := range c.serviceStore.List() {
		s := m.(*apiv1.Service)
		if err := f(NewService(s)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkDeployments(f func(Deployment) error) error {
	if c.deploymentStore == nil {
		return nil
	}
	for _, m := range c.deploymentStore.List() {
		d := m.(*apiappsv1.Deployment)
		if err := f(NewDeployment(d)); err != nil {
			return err
		}
	}
	return nil
}

// WalkDaemonSets calls f for each daemonset
func (c *client) WalkDaemonSets(f func(DaemonSet) error) error {
	if c.daemonSetStore == nil {
		return nil
	}
	for _, m := range c.daemonSetStore.List() {
		ds := m.(*apiappsv1.DaemonSet)
		if err := f(NewDaemonSet(ds)); err != nil {
			return err
		}
	}
	return nil
}

// WalkStatefulSets calls f for each statefulset
func (c *client) WalkStatefulSets(f func(StatefulSet) error) error {
	if c.statefulSetStore == nil {
		return nil
	}
	for _, m := range c.statefulSetStore.List() {
		s := m.(*apiappsv1.StatefulSet)
		if err := f(NewStatefulSet(s)); err != nil {
			return err
		}
	}
	return nil
}

// WalkCronJobs calls f for each cronjob
func (c *client) WalkCronJobs(f func(CronJob) error) error {
	if c.cronJobStore == nil {
		return nil
	}
	// We index jobs by id to make lookup for each cronjob more efficient
	jobs := map[types.UID]*apibatchv1.Job{}
	for _, m := range c.jobStore.List() {
		j := m.(*apibatchv1.Job)
		jobs[j.UID] = j
	}
	for _, m := range c.cronJobStore.List() {
		cj := m.(*apibatchv1beta1.CronJob)
		if err := f(NewCronJob(cj, jobs)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkNamespaces(f func(NamespaceResource) error) error {
	for _, m := range c.namespaceStore.List() {
		namespace := m.(*apiv1.Namespace)
		if err := f(NewNamespace(namespace)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkVolumeSnapshots(f func(VolumeSnapshot) error) error {
	for _, m := range c.volumeSnapshotStore.List() {
		volumeSnapshot := m.(*snapshotv1.VolumeSnapshot)
		if err := f(NewVolumeSnapshot(volumeSnapshot)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkVolumeSnapshotData(f func(VolumeSnapshotData) error) error {
	for _, m := range c.volumeSnapshotDataStore.List() {
		volumeSnapshotData := m.(*snapshotv1.VolumeSnapshotData)
		if err := f(NewVolumeSnapshotData(volumeSnapshotData)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCsiVolumeSnapshots(f func(CsiVolumeSnapshot) error) error {
	for _, m := range c.csiVolumeSnapshotStore.List() {
		cvs := m.(*csisnapshotv1beta1.VolumeSnapshot)
		if err := f(NewCsiVolumeSnapshot(cvs)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkVolumeSnapshotClasses(f func(VolumeSnapshotClass) error) error {
	for _, m := range c.volumeSnapshotClassStore.List() {
		vsc := m.(*csisnapshotv1beta1.VolumeSnapshotClass)
		if err := f(NewVolumeSnapshotClass(vsc)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkVolumeSnapshotContents(f func(VolumeSnapshotContent) error) error {
	for _, m := range c.volumeSnapshotContentStore.List() {
		vsc := m.(*csisnapshotv1beta1.VolumeSnapshotContent)
		if err := f(NewVolumeSnapshotContent(vsc)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkJobs(f func(Job) error) error {
	for _, m := range c.jobStore.List() {
		job := m.(*apibatchv1.Job)
		if err := f(NewJob(job)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkDisks(f func(Disk) error) error {
	for _, m := range c.diskStore.List() {
		disk := m.(*mayav1alpha1.Disk)
		if err := f(NewDisk(disk)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkBlockDevices(f func(BlockDevice) error) error {
	for _, m := range c.blockDeviceStore.List() {
		blockDevice := m.(*mayav1alpha1.BlockDevice)
		if err := f(NewBlockDevice(blockDevice)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkStoragePoolClaims(f func(StoragePoolClaim) error) error {
	for _, m := range c.storagePoolClaimStore.List() {
		spc := m.(*mayav1alpha1.StoragePoolClaim)
		if err := f(NewStoragePoolClaim(spc)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCStorVolumes(f func(CStorVolume) error) error {
	for _, m := range c.cStorvolumeStore.List() {
		cStorVolume := m.(*mayav1alpha1.CStorVolume)
		if err := f(NewCStorVolume(cStorVolume)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCStorVolumeReplicas(f func(CStorVolumeReplica) error) error {
	for _, m := range c.cStorvolumeReplicaStore.List() {
		cStorVolumeReplica := m.(*mayav1alpha1.CStorVolumeReplica)
		if err := f(NewCStorVolumeReplica(cStorVolumeReplica)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCStorPools(f func(CStorPool) error) error {
	for _, m := range c.cStorPoolStore.List() {
		cStorPool := m.(*mayav1alpha1.CStorPool)
		if err := f(NewCStorPool(cStorPool)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkBlockDeviceClaims(f func(BlockDeviceClaim) error) error {
	for _, m := range c.blockDeviceClaimStore.List() {
		blockDeviceClaim := m.(*mayav1alpha1.BlockDeviceClaim)
		if err := f(NewBlockDeviceClaim(blockDeviceClaim)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCStorPoolClusters(f func(CStorPoolCluster) error) error {
	for _, m := range c.cStorPoolClusterStore.List() {
		cStorPoolCluster := m.(*mayav1alpha1.CStorPoolCluster)
		if err := f(NewCStorPoolCluster(cStorPoolCluster)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) WalkCStorPoolInstances(f func(CStorPoolInstance) error) error {
	for _, m := range c.cStorPoolInstanceStore.List() {
		cStorPoolInstance := m.(*mayav1alpha1.CStorPoolInstance)
		if err := f(NewCStorPoolInstance(cStorPoolInstance)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) CloneVolumeSnapshot(namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity string, ctx context.Context) error {
	var scName string
	var claimSize string
	UID := strings.Split(uuid.New(), "-")
	scProvisionerName := "volumesnapshot.external-storage.k8s.io/snapshot-promoter"
	scList, err := c.client.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	// Retrieve the first snapshot-promoter storage class
	for _, sc := range scList.Items {
		if sc.Provisioner == scProvisionerName {
			scName = sc.Name
			break
		}
	}
	if scName == "" {
		return errors.New("snapshot-promoter storage class is not present")
	}
	volumeSnapshot, _ := c.snapshotClient.VolumesnapshotV1().VolumeSnapshots(namespaceID).Get(ctx, volumeSnapshotID, metav1.GetOptions{})
	if volumeSnapshot.Spec.PersistentVolumeClaimName != "" {
		persistentVolumeClaim, err := c.client.CoreV1().PersistentVolumeClaims(namespaceID).Get(ctx, volumeSnapshot.Spec.PersistentVolumeClaimName, metav1.GetOptions{})
		if err == nil {
			storage := persistentVolumeClaim.Spec.Resources.Requests[apiv1.ResourceStorage]
			if storage.String() != "" {
				claimSize = storage.String()
			}
		}
	}
	// Set default volume size to the one stored in volume snapshot annotation,
	// if unable to get PVC size.
	if claimSize == "" {
		claimSize = capacity
	}

	persistentVolumeClaim := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clone-" + persistentVolumeClaimID + "-" + UID[1],
			Namespace: namespaceID,
			Annotations: map[string]string{
				"snapshot.alpha.kubernetes.io/snapshot": volumeSnapshotID,
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse(claimSize),
				},
			},
		},
	}
	_, err = c.client.CoreV1().PersistentVolumeClaims(namespaceID).Create(ctx, persistentVolumeClaim, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateVolumeSnapshot(namespaceID, persistentVolumeClaimID, capacity, driver string, ctx context.Context) error {
	UID := strings.Split(uuid.New(), "-")
	snapshotName := "snapshot-" + time.Now().Format("20060102150405") + "-" + UID[1]
	var err error
	if _, ok := csiDriverMap[driver]; ok {
		err = c.createCsiVolumeSnapshot(ctx, snapshotName, namespaceID, persistentVolumeClaimID, capacity, driver)
	} else {
		err = c.createVolumeSnapshot(ctx, snapshotName, namespaceID, persistentVolumeClaimID, capacity)
	}
	return err
}

func (c *client) GetLogs(namespaceID, podID string, containerNames []string, ctx context.Context) (io.ReadCloser, error) {
	readClosersWithLabel := map[io.ReadCloser]string{}
	for _, container := range containerNames {
		req := c.client.CoreV1().Pods(namespaceID).GetLogs(
			podID,
			&apiv1.PodLogOptions{
				Follow:     true,
				Timestamps: true,
				Container:  container,
			},
		)
		readCloser, err := req.Stream(ctx)
		if err != nil {
			for rc := range readClosersWithLabel {
				rc.Close()
			}
			return nil, err
		}
		readClosersWithLabel[readCloser] = container
	}

	return NewLogReadCloser(readClosersWithLabel), nil
}

func (c *client) Describe(namespaceID, resourceID string, groupKind schema.GroupKind, restMapping apimeta.RESTMapping) (io.ReadCloser, error) {
	readClosersWithLabel := map[io.ReadCloser]string{}
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	describer, ok := kubectldescribe.DescriberFor(groupKind, restConfig)
	if !ok {
		describer, ok = kubectldescribe.GenericDescriberFor(&restMapping, restConfig)
		if !ok {
			return nil, errors.New("Resource not found")
		}
	}
	describerSetting := kubectldescribe.DescriberSettings{
		ShowEvents: true,
	}
	obj, err := describer.Describe(namespaceID, resourceID, describerSetting)
	if err != nil {
		return nil, err
	}
	formattedObj := ioutil.NopCloser(bytes.NewReader([]byte(obj)))
	readClosersWithLabel[formattedObj] = "describe"

	return NewLogReadCloser(readClosersWithLabel), nil
}

func (c *client) DeletePod(namespaceID, podID string, ctx context.Context) error {
	return c.client.CoreV1().Pods(namespaceID).Delete(ctx, podID, metav1.DeleteOptions{})
}

func (c *client) DeleteVolumeSnapshot(namespaceID, volumeSnapshotID string, ctx context.Context) error {
	return c.snapshotClient.VolumesnapshotV1().VolumeSnapshots(namespaceID).Delete(ctx, volumeSnapshotID, metav1.DeleteOptions{})
}

func (c *client) DeleteCsiVolumeSnapshot(namespaceID, volumeSnapshotID string, ctx context.Context) error {
	return c.csiSnapshotClient.SnapshotV1beta1().VolumeSnapshots(namespaceID).Delete(ctx, volumeSnapshotID, metav1.DeleteOptions{})
}

func (c *client) ScaleUp(namespaceID, id string, ctx context.Context) error {
	return c.modifyScale(namespaceID, id, ctx, func(scale *autoscalingv1.Scale) {
		scale.Spec.Replicas++
	})
}

func (c *client) ScaleDown(namespaceID, id string, ctx context.Context) error {
	return c.modifyScale(namespaceID, id, ctx, func(scale *autoscalingv1.Scale) {
		scale.Spec.Replicas--
	})
}

func (c *client) modifyScale(namespaceID, id string, ctx context.Context, f func(*autoscalingv1.Scale)) error {
	scaler := c.client.AppsV1().Deployments(namespaceID)
	scale, err := scaler.GetScale(ctx, id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	f(scale)
	_, err = scaler.UpdateScale(ctx, id, scale, metav1.UpdateOptions{})
	return err
}

func (c *client) Stop() {
	close(c.quit)
}

func (c *client) createVolumeSnapshot(ctx context.Context, name, namespaceID, persistentVolumeClaimID, capacity string) error {
	volumeSnapshot := &snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceID,
			Annotations: map[string]string{
				"capacity": capacity,
			},
		},
		Spec: snapshotv1.VolumeSnapshotSpec{
			PersistentVolumeClaimName: persistentVolumeClaimID,
		},
	}
	_, err := c.snapshotClient.VolumesnapshotV1().VolumeSnapshots(namespaceID).Create(ctx, volumeSnapshot, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) createCsiVolumeSnapshot(ctx context.Context, name, namespaceID, persistentVolumeClaimID, capacity, driver string) error {
	volumeSnapshotClassName := strings.ReplaceAll(driver, ".", "-")
	_, err := c.csiSnapshotClient.SnapshotV1beta1().VolumeSnapshotClasses().Get(ctx, volumeSnapshotClassName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			err = c.createVolumeSnapshotClass(ctx, volumeSnapshotClassName, driver)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	volumeSnapshot := &csisnapshotv1beta1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceID,
			Annotations: map[string]string{
				Capacity:         capacity,
				DriverAnnotation: driver,
			},
		},
		Spec: csisnapshotv1beta1.VolumeSnapshotSpec{
			VolumeSnapshotClassName: &volumeSnapshotClassName,
			Source: csisnapshotv1beta1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &persistentVolumeClaimID,
			},
		},
	}
	_, err = c.csiSnapshotClient.SnapshotV1beta1().VolumeSnapshots(namespaceID).Create(ctx, volumeSnapshot, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) createVolumeSnapshotClass(ctx context.Context, name, driver string) error {
	volumeSnapshotClass := &csisnapshotv1beta1.VolumeSnapshotClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Driver:         driver,
		DeletionPolicy: csisnapshotv1beta1.VolumeSnapshotContentDelete,
	}

	_, err := c.csiSnapshotClient.SnapshotV1beta1().VolumeSnapshotClasses().Create(ctx, volumeSnapshotClass, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CloneCsiVolumeSnapshot(namespaceID, volumeSnapshotID, persistentVolumeClaimID, capacity, driver string, ctx context.Context) error {
	var scName string
	var claimSize string
	UID := strings.Split(uuid.New(), "-")
	scList, err := c.client.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	// Retrieve the first snapshot-promoter storage class
	for _, sc := range scList.Items {
		if sc.Provisioner == driver {
			scName = sc.Name
			break
		}
	}
	if scName == "" {
		return errors.New("csi driver " + driver + " related storage class is not present")
	}
	persistentVolumeClaim, err := c.client.CoreV1().PersistentVolumeClaims(namespaceID).Get(ctx, persistentVolumeClaimID, metav1.GetOptions{})
	if err == nil {
		storage := persistentVolumeClaim.Spec.Resources.Requests[apiv1.ResourceStorage]
		if storage.String() != "" {
			claimSize = storage.String()
		}
	}
	// Set default volume size to the one stored in volume snapshot annotation,
	// if unable to get PVC size.
	if claimSize == "" {
		claimSize = capacity
	}

	csiAPIGroup := "snapshot.storage.k8s.io"
	clonePVC := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "clone-" + persistentVolumeClaimID + "-" + UID[1],
			Namespace: namespaceID,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse(claimSize),
				},
			},
			DataSource: &apiv1.TypedLocalObjectReference{
				APIGroup: &csiAPIGroup,
				Name:     volumeSnapshotID,
				Kind:     "VolumeSnapshot",
			},
		},
	}
	_, err = c.client.CoreV1().PersistentVolumeClaims(namespaceID).Create(ctx, clonePVC, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
