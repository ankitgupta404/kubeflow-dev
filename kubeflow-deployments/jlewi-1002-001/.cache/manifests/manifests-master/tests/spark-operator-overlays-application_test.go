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

func writeSparkOperatorOverlaysApplication(th *KustTestHarness) {
	th.writeF("/manifests/spark/spark-operator/overlays/application/application.yaml", `
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: $(generateName)
spec:
  selector:
    matchLabels:
      app.kubernetes.io/instance: $(generateName)
  componentKinds:
  - group: core
    kind: Service
  - group: apps
    kind: Deployment
  - group: core
    kind: ConfigMap
  - group: core
    kind: ServiceAccount
  - group: kubeflow.org
    kind: SparkOperator
  descriptor:
    type: "spark-operator"
    version: "v1"
    description: "Spark-operator allows users to create and manage the \"SparkApplication\" custom resource."
    maintainers:
    - name: Holden Karau
      email: holden@pigscanfly.ca
    owners:
    - name: Holden Karau
      email: holden@pigscanfly.ca
    keywords:
    - "spark"
  addOwnerRef: true
`)
	th.writeF("/manifests/spark/spark-operator/overlays/application/params.yaml", `
varReference:
- path: metadata/name
  kind: Application
- path: spec/selector/matchLabels/app.kubernetes.io\/instance
  kind: Application
- path: spec/selector/app.kubernetes.io\/instance
  kind: Service
- path: spec/selector/matchLabels/app.kubernetes.io\/instance
  kind: StatefulSet
- path: spec/template/metadata/labels/app.kubernetes.io\/instance
  kind: StatefulSet
`)
	th.writeF("/manifests/spark/spark-operator/overlays/application/params.env", `
generateName=
`)
	th.writeK("/manifests/spark/spark-operator/overlays/application", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- application.yaml
configMapGenerator:
- name: spark-operator-app-parameters
  env: params.env
vars:
- name: generateName
  objref:
    kind: ConfigMap
    name: spark-operator-app-parameters 
    apiVersion: v1
  fieldref:
    fieldpath: data.generateName
configurations:
- params.yaml
commonLabels:
  app.kubernetes.io/name: sparkoperator
  app.kubernetes.io/instance: $(generateName)
  app.kubernetes.io/managed-by: kfctl
  app.kubernetes.io/component: spark-operator
  app.kubernetes.io/part-of: kubeflow
  app.kubernetes.io/version: v0.6
`)
	th.writeF("/manifests/spark/spark-operator/base/spark-sa.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: spark-sa
`)
	th.writeF("/manifests/spark/spark-operator/base/cr-clusterrole.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: operator-cr
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - services
  - configmaps
  - secrets
  verbs:
  - create
  - get
  - delete
  - update
- apiGroups:
  - extensions
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - create
  - get
  - delete
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - update
  - patch
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - create
  - get
  - update
  - delete
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  verbs:
  - create
  - get
  - update
  - delete
- apiGroups:
  - sparkoperator.k8s.io
  resources:
  - sparkapplications
  - scheduledsparkapplications
  verbs:
  - '*'
`)
	th.writeF("/manifests/spark/spark-operator/base/crb.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: sparkoperator-crb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: operator-cr
subjects:
- kind: ServiceAccount
  name: operator-sa
`)
	th.writeF("/manifests/spark/spark-operator/base/crd-cleanup-job.yaml", `
apiVersion: batch/v1
kind: Job
metadata:
  name: crd-cleanup
  namespace: default
spec:
  template:
    spec:
      containers:
      - command:
        - /bin/sh
        - -c
        - 'curl -ik -X DELETE -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
          -H "Accept: application/json" -H "Content-Type: application/json" https://kubernetes.default.svc/apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/sparkapplications.sparkoperator.k8s.io'
        image: gcr.io/spark-operator/spark-operator:v2.4.0-v1beta1-0.8.2
        imagePullPolicy: IfNotPresent
        name: delete-sparkapp-crd
      - command:
        - /bin/sh
        - -c
        - 'curl -ik -X DELETE -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
          -H "Accept: application/json" -H "Content-Type: application/json" https://kubernetes.default.svc/apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/scheduledsparkapplications.sparkoperator.k8s.io'
        image: gcr.io/spark-operator/spark-operator:v2.4.0-v1beta1-0.8.2
        imagePullPolicy: IfNotPresent
        name: delete-scheduledsparkapp-crd
      restartPolicy: OnFailure
      serviceAccountName: operator-sa
`)
	th.writeF("/manifests/spark/spark-operator/base/deploy.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sparkoperator
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: sparkoperator
      app.kubernetes.io/version: v2.4.0-v1beta1-0.8.2
  strategy:
    type: Recreate
  template:
    metadata:
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      initializers:
        pending: []
      labels:
        app.kubernetes.io/name: sparkoperator
        app.kubernetes.io/version: v2.4.0-v1beta1-0.8.2
    spec:
      containers:
      - args:
        - -v=2
        - -namespace=
        - -ingress-url-format=
        - -install-crds=true
        - -controller-threads=10
        - -resync-interval=30
        - -logtostderr
        - -enable-metrics=true
        - -metrics-labels=app_type
        - -metrics-port=10254
        - -metrics-endpoint=/metrics
        - -metrics-prefix=
        image: gcr.io/spark-operator/spark-operator:v2.4.0-v1beta1-0.8.2
        imagePullPolicy: IfNotPresent
        name: sparkoperator
        ports:
        - containerPort: 10254
      serviceAccountName: operator-sa
`)
	th.writeF("/manifests/spark/spark-operator/base/operator-sa.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: operator-sa
`)
	th.writeK("/manifests/spark/spark-operator/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: kubeflow
# Labels to add to all resources and selectors.
commonLabels:
  kustomize.component: spark-operator
  app.kubernetes.io/instance: spark-operator
  app.kubernetes.io/name: sparkoperator

# Images modify the tags for images without
# creating patches.
images:
- name: gcr.io/spark-operator/spark-operator
  newTag: v2.4.0-v1beta1-0.8.2

# Value of this field is prepended to the
# names of all resources
namePrefix: spark-operator

# List of resource files that kustomize reads, modifies
# and emits as a YAML string
resources:
- spark-sa.yaml
- cr-clusterrole.yaml
- crb.yaml
- crd-cleanup-job.yaml
- deploy.yaml
- operator-sa.yaml
`)
}

func TestSparkOperatorOverlaysApplication(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/spark/spark-operator/overlays/application")
	writeSparkOperatorOverlaysApplication(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../spark/spark-operator/overlays/application"
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
