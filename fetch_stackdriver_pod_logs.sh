#!/bin/bash
# Fetch some events from stackdriver
set -ex
PROJECT=$1
POD=$2
FRESHNESS="${FRESHNESS:-24h}"
OUTFILE=~/tmp/${POD}.log.txt
gcloud logging read --format="table(timestamp, resource.labels.container_name, textPayload)" \
	--project=${PROJECT} \
	--freshness=${FRESHNESS} \
	--order asc  \
	"resource.type=\"k8s_container\" resource.labels.pod_name=\"${POD}\"  " > ${OUTFILE}

cat ${OUTFILE}
echo Logs written to ${OUTFILE}
