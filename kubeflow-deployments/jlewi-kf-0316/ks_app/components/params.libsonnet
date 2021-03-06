{
  global: {},
  components: {
    // Component-level parameters, defined initially from 'ks prototype use ...'
    // Each object below should correspond to a component in the components/ directory
    ambassador: {
      ambassadorImage: 'quay.io/datawire/ambassador:0.37.0',
      ambassadorNodePort: 0,
      ambassadorServiceType: 'NodePort',
      name: 'ambassador',
      platform: 'none',
      replicas: 3,
    },
    application: {
      components: ['ambassador', 'argo', 'basic-auth', 'basic-auth-ingress', 'centraldashboard', 'cert-manager', 'cloud-endpoints', 'jupyter', 'jupyter-web-app', 'katib', 'notebook-controller', 'pipeline', 'profiles', 'pytorch-operator', 'tensorboard', 'tf-job-operator'],
      extendedInfo: 'false',
      name: 'jlewi-kf-0316',
      type: 'kubeflow',
      version: '0.4',
    },
    argo: {
      artifactRepositoryAccessKeySecretKey: 'accesskey',
      artifactRepositoryAccessKeySecretName: 'mlpipeline-minio-artifact',
      artifactRepositoryBucket: 'mlpipeline',
      artifactRepositoryEndpoint: 'minio-service.kubeflow:9000',
      artifactRepositoryInsecure: 'true',
      artifactRepositoryKeyPrefix: 'artifacts',
      artifactRepositorySecretKeySecretKey: 'secretkey',
      artifactRepositorySecretKeySecretName: 'mlpipeline-minio-artifact',
      executorImage: 'argoproj/argoexec:v2.2.0',
      name: 'argo',
      uiImage: 'argoproj/argoui:v2.2.0',
      workflowControllerImage: 'argoproj/workflow-controller:v2.2.0',
    },
    "basic-auth": {
      authSecretName: 'kubeflow-login',
      image: 'gcr.io/kubeflow-images-public/gatekeeper:v20190211-v0.4.0-rc.1-119-g5098995b-e3b0c4',
      imageui: 'gcr.io/kubeflow-images-public/kflogin-ui:v20190123-v0.4.0-rc.1-73-g38ad5f77',
      name: 'basic-auth',
    },
    "basic-auth-ingress": {
      hostname: 'jlewi-kf-0316.endpoints.kubeflow-ci-deployment.cloud.goog',
      ingressSetupImage: 'gcr.io/kubeflow-images-public/ingress-setup:latest',
      ipName: 'jlewi-kf-0316-ip',
      issuer: 'letsencrypt-prod',
      name: 'basic-auth-ingress',
      privateGKECluster: 'false',
      secretName: 'envoy-ingress-tls',
    },
    centraldashboard: {
      image: 'gcr.io/kubeflow-images-public/centraldashboard:v20190315-v0.4.0-rc.1-222-gb42734af',
      name: 'centraldashboard',
    },
    "cert-manager": {
      acmeEmail: 'jlewi@google.com',
      acmeUrl: 'https://acme-v02.api.letsencrypt.org/directory',
      certManagerImage: 'quay.io/jetstack/cert-manager-controller:v0.4.0',
      name: 'cert-manager',
    },
    "cloud-endpoints": {
      name: 'cloud-endpoints',
      secretKey: 'admin-gcp-sa.json',
      secretName: 'admin-gcp-sa',
    },
    jupyter: {
      accessLocalFs: 'false',
      disks: 'null',
      gcpSecretName: 'user-gcp-sa',
      image: 'gcr.io/kubeflow/jupyterhub-k8s:v20180531-3bb991b1',
      jupyterHubAuthenticator: 'null',
      name: 'jupyter',
      notebookGid: '-1',
      notebookUid: '-1',
      platform: 'gke',
      rokSecretName: 'secret-rok-{username}',
      serviceType: 'ClusterIP',
      storageClass: 'null',
      ui: 'default',
      useJupyterLabAsDefault: 'false',
    },
    "jupyter-web-app": {
      image: 'gcr.io/kubeflow-images-public/jupyter-web-app:v20190228-v0.4.0-rc.1-173-g3ea53cc2',
      name: 'jupyter-web-app',
      policy: 'IfNotPresent',
      port: '80',
      prefix: 'jupyter',
      ui: 'default',
    },
    katib: {
      katibUIImage: 'gcr.io/kubeflow-images-public/katib/katib-ui:v0.1.2-alpha-138-g26da3ea',
      metricsCollectorImage: 'gcr.io/kubeflow-images-public/katib/metrics-collector:v0.1.2-alpha-138-g26da3ea',
      name: 'katib',
      studyJobControllerImage: 'gcr.io/kubeflow-images-public/katib/studyjob-controller:v0.1.2-alpha-138-g26da3ea',
      suggestionBayesianOptimizationImage: 'gcr.io/kubeflow-images-public/katib/suggestion-bayesianoptimization:v0.1.2-alpha-138-g26da3ea',
      suggestionGridImage: 'gcr.io/kubeflow-images-public/katib/suggestion-grid:v0.1.2-alpha-138-g26da3ea',
      suggestionHyperbandImage: 'gcr.io/kubeflow-images-public/katib/suggestion-hyperband:v0.1.2-alpha-138-g26da3ea',
      suggestionRandomImage: 'gcr.io/kubeflow-images-public/katib/suggestion-random:v0.1.2-alpha-138-g26da3ea',
      vizierCoreImage: 'gcr.io/kubeflow-images-public/katib/vizier-core:v0.1.2-alpha-138-g26da3ea',
      vizierCoreRestImage: 'gcr.io/kubeflow-images-public/katib/vizier-core-rest:v0.1.2-alpha-138-g26da3ea',
      vizierDbImage: 'mysql:8.0.3',
    },
    metacontroller: {
      image: 'metacontroller/metacontroller:v0.3.0',
      name: 'metacontroller',
    },
    "notebook-controller": {
      controllerImage: 'gcr.io/kubeflow-images-public/notebook-controller:v20190304-v0.4.0-rc.1-189-g157a3399-dirty-47ed08',
      name: 'notebook-controller',
    },
    pipeline: {
      minioPd: 'jlewi-kf-0316-storage-artifact-store',
      mysqlPd: 'jlewi-kf-0316-storage-metadata-store',
      name: 'pipeline',
    },
    profiles: {
      image: 'gcr.io/kubeflow-images-public/profile-controller:v20190228-v0.4.0-rc.1-192-g1a802656-dirty-f95773',
      name: 'profiles',
    },
    "pytorch-operator": {
      cloud: 'null',
      deploymentNamespace: 'null',
      deploymentScope: 'cluster',
      disks: 'null',
      name: 'pytorch-operator',
      pytorchDefaultImage: 'null',
      pytorchJobImage: 'gcr.io/kubeflow-images-public/pytorch-operator:v0.5.0-rc.1-6-gd09eb65',
    },
    tensorboard: {
      defaultTbImage: 'tensorflow/tensorflow:1.8.0',
      logDir: 'logs',
      name: 'tensorboard',
      servicePort: 9000,
      serviceType: 'ClusterIP',
      targetPort: 6006,
    },
    "tf-job-operator": {
      cloud: 'null',
      deploymentNamespace: 'null',
      deploymentScope: 'cluster',
      name: 'tf-job-operator',
      tfDefaultImage: 'null',
      tfJobImage: 'gcr.io/kubeflow-images-public/tf_operator:kubeflow-tf-operator-postsubmit-v1beta1-c284947-309-b42b',
      tfJobUiServiceType: 'ClusterIP',
    },
  },
}