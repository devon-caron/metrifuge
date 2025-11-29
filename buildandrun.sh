#!/bin/bash

# Continue execution even if commands fail (e.g., pod deletion when pod doesn't exist)
# Only stop on Ctrl+C (SIGINT)
set +e

# Trap Ctrl+C for clean exit
trap 'echo -e "\n==== Interrupted by user ===="; exit 130' INT

echo "==== Building..."
./build.sh

echo "==== Deleting old pods..."
kubectl delete po/metrifuge-test -n metrifuge
kubectl delete po/mf-log-pod -n metrifuge

echo "==== Loading image..."
mftagtbl=$(docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}" | head -n 2)
mftagname=$(echo "$mftagtbl" | grep metrifuge | grep -v "latest" | awk '{print $1}')
mftagtime=$(echo "$mftagtbl" | grep metrifuge | grep -v "latest" | awk '{print $2}')
minikube image load "$mftagname:$mftagtime"

echo "==== Updating pods..."
yq eval --inplace ".spec.containers[0].image = \"$mftagname:$mftagtime\"" mf-pod.yaml
kubectl apply -f mf-log-pod.yaml
kubectl apply -f mf-pod.yaml
sleep 2

echo "==== Following logs..."
if [ -n "$1" ]; then
  echo "Applying filter: $1"
  eval "kubectl logs -n metrifuge po/metrifuge-test -f $1"
else
  kubectl logs -n metrifuge po/metrifuge-test -f
fi
