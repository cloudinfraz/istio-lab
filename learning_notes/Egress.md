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

2. Envoy passthrough to external services
  
 kubectl -n istio-system get cm istio -o yaml |grep -iA 3 outboundTrafficPolicy |head -2
    outboundTrafficPolicy:
      mode: ALLOW_ANY
    localityLbSetting:
      enabled: true

3. Test the access 

 kubectl exec "$SOURCE_POD" -c sleep -- curl -sI https://www.baidu.com | grep "HTTP/"
 kubectl exec "$SOURCE_POD" -c sleep -- curl -sI https://edition.cnn.com | grep "HTTP/"

4. Change to the blocking-by-default policy
  kubectl -n istio-system get cm istio -o yaml |grep -iA 3 outboundTrafficPolicy |head -2
    outboundTrafficPolicy:
      mode: REGISTRY_ONLY

 kubectl exec "$SOURCE_POD" -c sleep -- curl -sI https://www.baidu.com | grep "HTTP/"
command terminated with exit code 35

```
# Access an external HTTP service

```
1. Create a ServiceEntry to allow access to an external HTTP service
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: httpbin-ext
spec:
  hosts:
  - httpbin.org
  ports:
  - number: 80
    name: http
    protocol: HTTP
  resolution: DNS
  location: MESH_EXTERNAL
EOF

2. Make a request to the external HTTP service from SOURCE_POD

 kubectl exec "$SOURCE_POD" -c sleep -- curl -s http://httpbin.org/headers

```
# Direct access to external services(Red Hat service mesh ?)

```
values.global.proxy.includeIPRanges="10.96.0.0/12"

```
# Egress TLS Origination

```
1. Configuring access to an external service
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: edition-cnn-com
spec:
  hosts:
  - edition.cnn.com
  ports:
  - number: 80
    name: http-port
    protocol: HTTP
  - number: 443
    name: https-port
    protocol: HTTPS
  resolution: DNS
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: edition-cnn-com
spec:
  hosts:
  - edition.cnn.com
  tls:
  - match:
    - port: 443
      sniHosts:
      - edition.cnn.com
    route:
    - destination:
        host: edition.cnn.com
        port:
          number: 443
      weight: 100
EOF

Make a request to the external HTTP service:

$ kubectl exec "${SOURCE_POD}" -c sleep -- curl -sL -o /dev/null -D - http://edition.cnn.com/politics
HTTP/1.1 301 Moved Permanently
...
location: https://edition.cnn.com/politics
...

HTTP/2 200
...

The output should be similar to the above (some details replaced by ellipsis).

Notice the -L flag of curl which instructs curl to follow redirects. In this case, the server returned a redirect response (301 Moved Permanently) for the HTTP request to http://edition.cnn.com/politics. The redirect response instructs the client to send an additional request, this time using HTTPS, to https://edition.cnn.com/politics. For the second request, the server returned the requested content and a 200 OK status code.

Although the curl command handled the redirection transparently, there are two issues here. The first issue is the redundant request, which doubles the latency of fetching the content of http://edition.cnn.com/politics. The second issue is that the path of the URL, politics in this case, is sent in clear text. If there is an attacker who sniffs the communication between your application and edition.cnn.com, the attacker would know which specific topics of edition.cnn.com the application fetched. For privacy reasons, you might want to prevent such disclosure.

Both of these issues can be resolved by configuring Istio to perform TLS origination.

TLS origination for egress traffic

Redefine your VirtualService from the previous section to rewrite the HTTP request port and add a DestinationRule to perform TLS origination.

$ kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: edition-cnn-com
spec:
  hosts:
  - edition.cnn.com
  http:
  - match:
    - port: 80
    route:
    - destination:
        host: edition.cnn.com
        subset: tls-origination
        port:
          number: 443
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: edition-cnn-com
spec:
  host: edition.cnn.com
  subsets:
  - name: tls-origination
    trafficPolicy:
      loadBalancer:
        simple: ROUND_ROBIN
      portLevelSettings:
      - port:
          number: 443
        tls:
          mode: SIMPLE # initiates HTTPS when accessing edition.cnn.com
EOF

2. Send an HTTP request to http://edition.cnn.com/politics

 kubectl exec "${SOURCE_POD}" -c sleep -- curl -sL -o /dev/null -D - http://edition.cnn.com/politics

3. Note that the applications that used HTTPS to access the external service 

  kubectl exec "${SOURCE_POD}" -c sleep -- curl -sL -o /dev/null -D - https://edition.cnn.com/politics
HTTP/2 200
```
# Egress Gateways

```
# Scenrario setup

1. Set the SOURCE_POD environmen

 export SOURCE_POD=$(kubectl get pod -l app=sleep -o jsonpath={.items..metadata.name})

2. verify Egressgateway

  kubectl get pod -l istio=egressgateway -n istio-system
  

```
