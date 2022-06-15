/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"errors"

	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v3/apis/volumesnapshot/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	monitoringdashboardv1alpha1 "kubesphere.io/monitoring-dashboard/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/hzhhong/kubesphere/pkg/api"
	clusterv1alpha1 "github.com/hzhhong/kubesphere/pkg/apis/cluster/v1alpha1"
	devopsv1alpha3 "github.com/hzhhong/kubesphere/pkg/apis/devops/v1alpha3"
	iamv1alpha2 "github.com/hzhhong/kubesphere/pkg/apis/iam/v1alpha2"
	networkv1alpha1 "github.com/hzhhong/kubesphere/pkg/apis/network/v1alpha1"
	notificationv2beta1 "github.com/hzhhong/kubesphere/pkg/apis/notification/v2beta1"
	tenantv1alpha1 "github.com/hzhhong/kubesphere/pkg/apis/tenant/v1alpha1"
	tenantv1alpha2 "github.com/hzhhong/kubesphere/pkg/apis/tenant/v1alpha2"
	typesv1beta1 "github.com/hzhhong/kubesphere/pkg/apis/types/v1beta1"
	"github.com/hzhhong/kubesphere/pkg/apiserver/query"
	"github.com/hzhhong/kubesphere/pkg/informers"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/application"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/cluster"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/clusterdashboard"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/clusterrole"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/clusterrolebinding"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/configmap"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/customresourcedefinition"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/daemonset"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/dashboard"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/deployment"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/devops"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedapplication"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedconfigmap"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federateddeployment"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedingress"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatednamespace"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedpersistentvolumeclaim"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedsecret"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedservice"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/federatedstatefulset"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/globalrole"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/globalrolebinding"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/group"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/groupbinding"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/ingress"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/ippool"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/job"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/loginrecord"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/namespace"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/networkpolicy"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/node"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/notification"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/persistentvolumeclaim"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/pod"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/role"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/rolebinding"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/secret"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/service"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/serviceaccount"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/statefulset"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/user"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/volumesnapshot"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/workspace"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/workspacerole"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/workspacerolebinding"
	"github.com/hzhhong/kubesphere/pkg/models/resources/v1alpha3/workspacetemplate"
)

var ErrResourceNotSupported = errors.New("resource is not supported")

type ResourceGetter struct {
	clusterResourceGetters    map[schema.GroupVersionResource]v1alpha3.Interface
	namespacedResourceGetters map[schema.GroupVersionResource]v1alpha3.Interface
}

func NewResourceGetter(factory informers.InformerFactory, cache cache.Cache) *ResourceGetter {
	namespacedResourceGetters := make(map[schema.GroupVersionResource]v1alpha3.Interface)
	clusterResourceGetters := make(map[schema.GroupVersionResource]v1alpha3.Interface)

	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}] = deployment.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}] = daemonset.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}] = statefulset.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}] = service.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}] = namespace.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}] = configmap.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}] = secret.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}] = pod.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}] = node.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "serviceaccounts"}] = serviceaccount.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"}] = ingress.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}] = networkpolicy.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}] = job.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "app.k8s.io", Version: "v1beta1", Resource: "applications"}] = application.New(cache)
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumeclaims"}] = persistentvolumeclaim.New(factory.KubernetesSharedInformerFactory(), factory.SnapshotSharedInformerFactory())
	namespacedResourceGetters[snapshotv1beta1.SchemeGroupVersion.WithResource("volumesnapshots")] = volumesnapshot.New(factory.SnapshotSharedInformerFactory())
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralRoleBinding)] = rolebinding.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralRole)] = role.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}] = customresourcedefinition.New(factory.ApiExtensionSharedInformerFactory())

	// kubesphere resources
	namespacedResourceGetters[networkv1alpha1.SchemeGroupVersion.WithResource(networkv1alpha1.ResourcePluralIPPool)] = ippool.New(factory.KubeSphereSharedInformerFactory(), factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[devopsv1alpha3.SchemeGroupVersion.WithResource(devopsv1alpha3.ResourcePluralDevOpsProject)] = devops.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[tenantv1alpha1.SchemeGroupVersion.WithResource(tenantv1alpha1.ResourcePluralWorkspace)] = workspace.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[tenantv1alpha1.SchemeGroupVersion.WithResource(tenantv1alpha2.ResourcePluralWorkspaceTemplate)] = workspacetemplate.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralGlobalRole)] = globalrole.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralWorkspaceRole)] = workspacerole.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralUser)] = user.New(factory.KubeSphereSharedInformerFactory(), factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralGlobalRoleBinding)] = globalrolebinding.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralWorkspaceRoleBinding)] = workspacerolebinding.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralLoginRecord)] = loginrecord.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralGroup)] = group.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralGroupBinding)] = groupbinding.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralClusterRole)] = clusterrole.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralClusterRoleBinding)] = clusterrolebinding.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[clusterv1alpha1.SchemeGroupVersion.WithResource(clusterv1alpha1.ResourcesPluralCluster)] = cluster.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[notificationv2beta1.SchemeGroupVersion.WithResource(notificationv2beta1.ResourcesPluralConfig)] = notification.NewNotificationConfigGetter(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[notificationv2beta1.SchemeGroupVersion.WithResource(notificationv2beta1.ResourcesPluralReceiver)] = notification.NewNotificationReceiverGetter(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[monitoringdashboardv1alpha1.GroupVersion.WithResource("clusterdashboards")] = clusterdashboard.New(cache)

	// federated resources
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedNamespace)] = federatednamespace.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedDeployment)] = federateddeployment.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedSecret)] = federatedsecret.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedConfigmap)] = federatedconfigmap.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedService)] = federatedservice.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedApplication)] = federatedapplication.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedPersistentVolumeClaim)] = federatedpersistentvolumeclaim.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedStatefulSet)] = federatedstatefulset.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[typesv1beta1.SchemeGroupVersion.WithResource(typesv1beta1.ResourcePluralFederatedIngress)] = federatedingress.New(factory.KubeSphereSharedInformerFactory())
	namespacedResourceGetters[monitoringdashboardv1alpha1.GroupVersion.WithResource("dashboards")] = dashboard.New(cache)

	return &ResourceGetter{
		namespacedResourceGetters: namespacedResourceGetters,
		clusterResourceGetters:    clusterResourceGetters,
	}
}

// TryResource will retrieve a getter with resource name, it doesn't guarantee find resource with correct group version
// need to refactor this use schema.GroupVersionResource
func (r *ResourceGetter) TryResource(clusterScope bool, resource string) v1alpha3.Interface {
	if clusterScope {
		for k, v := range r.clusterResourceGetters {
			if k.Resource == resource {
				return v
			}
		}
	}
	for k, v := range r.namespacedResourceGetters {
		if k.Resource == resource {
			return v
		}
	}
	return nil
}

func (r *ResourceGetter) Get(resource, namespace, name string) (runtime.Object, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.Get(namespace, name)
}

func (r *ResourceGetter) List(resource, namespace string, query *query.Query) (*api.ListResult, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.List(namespace, query)
}
