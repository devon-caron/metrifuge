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

# Start timer in background
start_time=$SECONDS
(
  while true; do
    elapsed=$((SECONDS - start_time))
    printf "\rElapsed time: %02d:%02d" $((elapsed / 60)) $((elapsed % 60))
    sleep 1
  done
) &
timer_pid=$!

# Load the image
minikube image load "$mftagname:$mftagtime"

# Stop the timer and show final time
kill $timer_pid 2>/dev/null
wait $timer_pid 2>/dev/null
final_elapsed=$((SECONDS - start_time))
printf "\rElapsed time: %02d:%02d\n" $((final_elapsed / 60)) $((final_elapsed % 60))

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
