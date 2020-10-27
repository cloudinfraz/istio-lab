# install verification 
 kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"

# ingress env
export INGRESS_HOST=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')

# adopt for the native k8s

```
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')
export TCP_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].port}')
export INGRESS_PORT=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.port.targetPort}')

# istio gateway access
export GATEWAY_URL=$INGRESS_HOST
curl -s "http://${GATEWAY_URL}/productpage" | grep -o "<title>.*</title>"

```
