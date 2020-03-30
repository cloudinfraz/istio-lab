#!/bin/bash
#CA

Set_Probe() {
   oc get smcp/service-mesh-installation -n $1  -o yaml |sed -e 's/rewriteAppHTTPProbe: false/rewriteAppHTTPProbe: true/g' | oc  apply -f -
}

if [ $# != 1 ];then
    echo "Please input the right istio namespace!"
    exit 1
else 
   Set_Probe $1
fi
