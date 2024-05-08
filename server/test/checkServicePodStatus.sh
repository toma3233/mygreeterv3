#!/bin/bash

# Define the namespaces
NAMESPACES=("servicehub-mygreeterv3-server" "servicehub-mygreeterv3-client" "servicehub-mygreeterv3-demoserver")

# define wait time before exit
POD_WAIT_TIME="120s"
NS_WAIT_TIME="5s"

echo "Checking pod status..."

# Loop through each namespace
for NAMESPACE in "${NAMESPACES[@]}"
do
    echo "Checking namespace: $NAMESPACE"
    kubectl wait --for jsonpath='{.status.phase}=Active' --timeout="$NS_WAIT_TIME" namespace/$NAMESPACE > /dev/null
    if [ $? -ne 0 ]
    then
        echo "ERROR: $NAMESPACE does not exist."   
        exit 1
    fi
    PODS=$(kubectl get pods -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}')
    if [ -z "$PODS" ]; then
      echo "ERROR: No pods are in this namespace."
      exit 1
    fi
    kubectl wait --for=jsonpath='{.status.phase}'=Running pods --all --namespace $NAMESPACE --timeout=$POD_WAIT_TIME > /dev/null
    if [ $? -ne 0 ]
    then
        echo "ERROR: $NAMESPACE pods did not run successfully."   
        exit 1
    fi
    kubectl wait --for=condition=ready pods --all --namespace $NAMESPACE --timeout=$POD_WAIT_TIME > /dev/null
    if [ $? -ne 0 ]
    then
        echo "ERROR: $NAMESPACE pods are not ready."   
        exit 1
    fi
    echo "Pods in $NAMESPACE are running."
done
echo "All pods are running."