./build.sh
kubectl delete po/metrifuge-test -n metrifuge
sleep 1
mftagtbl=$(docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}" | head -n 3)
mftagname=$(echo "$mftagtbl" | grep metrifuge | grep -v "latest" | awk '{print $1}')
mftagtime=$(echo "$mftagtbl" | grep metrifuge | grep -v "latest" | awk '{print $2}')
minikube image load "$mftagname:$mftagtime"
yq eval --inplace ".spec.containers[0].image = \"$mftagname:$mftagtime\"" mf-pod.yaml
kubectl apply -f mf-pod.yaml
