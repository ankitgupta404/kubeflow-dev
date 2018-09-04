{
  global: {},
  components: {
    // Component-level parameters, defined initially from 'ks prototype use ...'
    // Each object below should correspond to a component in the components/ directory
    "pytorch-operator": {
      cloud: 'null',
      disks: 'null',
      name: 'pytorch-operator',
      namespace: 'null',
      pytorchDefaultImage: 'null',
      pytorchJobImage: 'gcr.io/kubeflow-images-public/pytorch-operator:v20180619-2e19016',
    },
    ambassador: {
      ambassadorImage: 'quay.io/datawire/ambassador:0.37.0',
      ambassadorServiceType: 'ClusterIP',
      cloud: 'gke',
      name: 'ambassador',
      statsdExporterImage: 'prom/statsd-exporter:v0.6.0',
      statsdImage: 'quay.io/datawire/statsd:0.37.0',
    },
    jupyterhub: {
      accessLocalFs: 'false',
      cloud: 'gke',
      disks: 'null',
      gcpSecretName: 'user-gcp-sa',
      image: 'gcr.io/kubeflow/jupyterhub-k8s:v20180531-3bb991b1',
      jupyterHubAuthenticator: 'iap',
      name: 'jupyterhub',
      namespace: 'null',
      notebookGid: '-1',
      notebookPVCMount: '/home/jovyan',
      notebookUid: '-1',
      registry: 'gcr.io',
      repoName: 'kubeflow-images-public',
      serviceType: 'ClusterIP',
    },
    centraldashboard: {
      image: 'gcr.io/kubeflow-images-public/centraldashboard:v0.2.1',
      name: 'centraldashboard',
    },
    "tf-job-operator": {
      cloud: 'null',
      deploymentNamespace: 'null',
      deploymentScope: 'cluster',
      name: 'tf-job-operator',
      namespace: 'null',
      tfDefaultImage: 'null',
      tfJobImage: 'gcr.io/kubeflow-images-public/tf_operator:v20180901-454c5c1a',
      tfJobUiServiceType: 'ClusterIP',
      tfJobVersion: 'v1alpha2',
    },
    argo: {
      executorImage: 'argoproj/argoexec:v2.1.1',
      name: 'argo',
      namespace: 'null',
      uiImage: 'argoproj/argoui:v2.1.1',
      workflowControllerImage: 'argoproj/workflow-controller:v2.1.1',
    },
    spartakus: {
      name: 'spartakus',
      reportUsage: 'true',
      usageId: '279236202',
    },
    "cloud-endpoints": {
      name: 'cloud-endpoints',
      namespace: 'null',
      secretKey: 'admin-gcp-sa.json',
      secretName: 'admin-gcp-sa',
    },
    "cert-manager": {
      acmeEmail: 'jlewi@google.com',
      acmeUrl: 'https://acme-v02.api.letsencrypt.org/directory',
      certManagerImage: 'quay.io/jetstack/cert-manager-controller:v0.4.0',
      name: 'cert-manager',
      namespace: 'null',
    },
    "iap-ingress": {
      disableJwtChecking: 'false',
      envoyImage: 'gcr.io/kubeflow-images-public/envoy:v20180309-0fb4886b463698702b6a08955045731903a18738',
      hostname: 'jlewi-0825.endpoints.cloud-ml-dev.cloud.goog',
      ingressSetupImage: 'gcr.io/kubeflow-images-public/ingress-setup:latest',
      ipName: 'jlewi-0825-ip',
      issuer: 'letsencrypt-prod',
      name: 'iap-ingress',
      namespace: 'null',
      oauthSecretName: 'kubeflow-oauth',
      privateGKECluster: 'false',
      secretName: 'envoy-ingress-tls',
    },
  },
}