package tests_test

import (
	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/loader"
	"sigs.k8s.io/kustomize/v3/pkg/plugins"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/target"
	"sigs.k8s.io/kustomize/v3/pkg/validators"
	"testing"
)

func writeKubebenchOverlaysIstio(th *KustTestHarness) {
	th.writeF("/manifests/kubebench/overlays/istio/virtual-service.yaml", `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: kubebench-dashboard
spec:
  gateways:
  - kubeflow-gateway
  hosts:
  - '*'
  http:
  - match:
    - uri:
        prefix: /dashboard/
    rewrite:
      uri: /dashboard/
    route:
    - destination:
        host: kubebench-dashboard.$(namespace).svc.$(clusterDomain)
        port:
          number: 80
`)
	th.writeF("/manifests/kubebench/overlays/istio/params.yaml", `
varReference:
- path: spec/http/route/destination/host
  kind: VirtualService
`)
	th.writeK("/manifests/kubebench/overlays/istio", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- virtual-service.yaml
configurations:
- params.yaml
`)
	th.writeF("/manifests/kubebench/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kubebench-operator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubebench-operator
subjects:
- kind: ServiceAccount
  name: default
`)
	th.writeF("/manifests/kubebench/base/cluster-role.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: kubebench-operator
rules:
- apiGroups:
  - kubebench.operator
  resources:
  - kubebenchjobs.kubebench.operator
  - kubebenchjobs
  verbs:
  - create
  - update
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  - pods
  - pods/exec
  - services
  - endpoints
  - persistentvolumeclaims
  - events
  - secrets
  verbs:
  - '*'
- apiGroups:
  - kubeflow.org
  resources:
  - tfjobs
  - pytorchjobs
  verbs:
  - '*'
- apiGroups:
  - argoproj.io
  resources:
  - workflows
  verbs:
  - '*'
`)
	th.writeF("/manifests/kubebench/base/config-map.yaml", `
apiVersion: v1
kind: ConfigMap
metadata:
  name: kubebench-config
data:
  kubebenchconfig.yaml: |
    """
    defaultWorkflowAgent:
      container:
        name: kubebench-workflow-agent
        image: gcr.io/kubeflow-images-public/kubebench/workflow-agent:v0.5.0-11-gea53ad5
    defaultManagedVolumes:
      experimentVolume:
        name: kubebench-experiment-volume
        emptyDir: {}
      workflowVolume:
        name: kubebench-workflow-volume
        emptyDir: {}
    """
`)
	th.writeF("/manifests/kubebench/base/crd.yaml", `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kubebenchjobs.kubebench.operator
spec:
  group: kubebench.operator
  names:
    kind: KubebenchJob
    plural: kubebenchjobs
  scope: Namespaced
  version: v1
`)
	th.writeF("/manifests/kubebench/base/deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubebench-operator
spec:
  selector:
    matchLabels:
      app: kubebench-operator
  volumes:
    - name: kubebench-config
      configMap:
        name: kubebench-config
  template:
    metadata:
      labels:
        app: kubebench-operator
    spec:
      containers:
      - image: gcr.io/kubeflow-images-public/kubebench/kubebench-operator-v1alpha2:v0.5.0-11-gea53ad5
        name: kubebench-operator
      serviceAccountName: kubebench-operator
`)
	th.writeF("/manifests/kubebench/base/params.yaml", `
varReference:
- path: metadata/annotations/getambassador.io\/config
  kind: Service
`)
	th.writeF("/manifests/kubebench/base/params.env", `
namespace=
clusterDomain=cluster.local
`)
	th.writeK("/manifests/kubebench/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- cluster-role-binding.yaml
- cluster-role.yaml
- config-map.yaml
- crd.yaml
- deployment.yaml
namespace: kubeflow
commonLabels:
  kustomize.component: kubebench
configMapGenerator:
- name: parameters
  env: params.env
images:
  - name:  gcr.io/kubeflow-images-public/kubebench/kubebench-operator-v1alpha2
    newName: gcr.io/kubeflow-images-public/kubebench/kubebench-operator-v1alpha2
    newTag: v0.5.0-11-gea53ad5
  - name: gcr.io/kubeflow-images-public/kubebench/kubebench-controller
    newName: gcr.io/kubeflow-images-public/kubebench/kubebench-controller
    newTag: v0.4.0-13-g262c593
  - name: gcr.io/kubeflow-images-public/kubebench/kubebench-example-tf-cnn-post-processor
    newName: gcr.io/kubeflow-images-public/kubebench/kubebench-example-tf-cnn-post-processor
    newTag: v0.4.0-13-g262c593
vars:
- name: clusterDomain
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.clusterDomain
configurations:
- params.yaml
`)
}

func TestKubebenchOverlaysIstio(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/kubebench/overlays/istio")
	writeKubebenchOverlaysIstio(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../kubebench/overlays/istio"
	fsys := fs.MakeRealFS()
	lrc := loader.RestrictionRootOnly
	_loader, loaderErr := loader.NewLoader(lrc, validators.MakeFakeValidator(), targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()), transformer.NewFactoryImpl())
	pc := plugins.DefaultPluginConfig()
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl(), plugins.NewLoader(pc, rf))
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	actual, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(actual, string(expected))
}
