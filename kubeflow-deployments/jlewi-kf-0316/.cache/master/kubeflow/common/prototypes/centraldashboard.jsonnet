// @apiVersion 0.1
// @name io.ksonnet.pkg.centraldashboard
// @description centraldashboard Component
// @shortDescription centraldashboard
// @param name string Name
// @optionalParam image string gcr.io/kubeflow-images-public/centraldashboard:v20190315-v0.4.0-rc.1-222-gb42734af Image for the central dashboard

local centraldashboard = import "kubeflow/common/centraldashboard.libsonnet";
local instance = centraldashboard.new(env, params);
instance.list(instance.all)
