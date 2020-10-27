# Egress traffic

- Accessing External Services
- Egress TLS Origination
- Egress Gateways
- Egress Gateways with TLS Origination
- Egress using Wildcard Hosts
- Kubernetes Services for Egress Traffic
- Using an External HTTPS Proxy

# Accessing External Services

```
# Scenrario setup

1. Set the SOURCE_POD environmen

 export SOURCE_POD=$(kubectl get pod -l app=sleep -o jsonpath={.items..metadata.name})

2. Enable Envoy’s access logging
  
 oc -n istio-system edit smcp basic-install
 spec:
  global:
    proxy:
      accessLogFile: /dev/stdout

3. Test the access log

 kubectl exec "$SOURCE_POD" -c sleep -- curl -s httpbin:8000/status/418
 kubectl logs -l app=httpbin -c istio-proxy
 kubectl logs -l app=sleep -c istio-proxy


```
