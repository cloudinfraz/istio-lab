# Enabling Policy Enforcement
The mixer policy is deprecated in Istio 1.5 and not recommended for production usage.
Rate limiting: Consider using Envoy native rate limiting instead of mixer rate limiting.
  Istio will add support for native rate limiting API through the Istio extensions API.

https://github.com/envoyproxy/ratelimit

Control headers and routing: Consider using Envoy ext_authz filter, lua filter, or write a filter using the Envoy-wasm sandbox.

Denials and White/Black Listing: Please use the Authorization Policy for enforcing access control to a workload.

- Enabling Policy Enforcement
- Enabling Rate Limits
- Headers and Routing
- Denials and White/Black Listing

# Enabling Policy Enforcement

```
1. check policy enforcement
# native k8s/istioctl manifest apply --set values.global.disablePolicyChecks=false
kubectl -n istio-system get cm istio -o jsonpath="{@.data.mesh}" | grep disablePolicyChecks
oc -n istio-system edit cm istio
 kubectl -n istio-system get cm istio -o jsonpath="{@.data.mesh}" | grep disablePolicyChecks
disablePolicyChecks: false

```
# Install redis for Rate Limit
- note redisquota do not support authentication 

```
 oc -n istio-redis get pods
NAME             READY   STATUS      RESTARTS   AGE
redis-1-8mbhl    1/1     Running     0          49m
redis-1-deploy   0/1     Completed   0          49m
    
    获取密码
        redis-cli
        config get requirepass
    设置密码
        config set requirepass redhat
    当有密码的时候登录时需要密码登录
        auth 密码
    取消密码
        config set requirepass ''

```
# Enabling Rate Limit(Deprecated)

```
configure Istio to rate limit traffic to productpage based on the IP address of the originating client. You will use X-Forwarded-For request header as the client IP address
1. memory quota (memquota) adapter 
2. Redis quota (redisquota) adapter(production recommendation)

5 extenstion resource objects

Client Side

* QuotaSpec defines quota name and amount that the client should request.
* QuotaSpecBinding conditionally associates QuotaSpec with one or more services.

Mixer Side

* quota instance defines how quota is dimensioned by Mixer.
* memquota handler defines memquota adapter configuration.
* quota rule defines when quota instance is dispatched to the memquota adapter.

1. enable rate limits using redisquota:
   oc apply -f mixer-rule-productpage-redis-quota-rolling-window.yaml
2. Confirm the quota instance 
   oc -n istio-system get instance requestcountquota -o yaml
3. Confirm the quota rule
   oc -n istio-system get rule quota
4. Confirm the QuotaSpec was created
   oc -n istio-system get QuotaSpec request-count
5. Confirm the QuotaSpecBinding was created
   oc -n istio-system get QuotaSpecBinding request-count 
6. Refresh the product page in your browser

   export INGRESS_ROUTE=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
for i in `seq 20`;do echo -----$i------; curl -s "http://$INGRESS_ROUTE/productpage" ;done   
-----1------
<!DOCTYPE html>
<html>
...............
<!-- Latest compiled and minified JavaScript -->
<script src="static/jquery.min.js"></script>

<!-- Latest compiled and minified JavaScript -->
<script src="static/bootstrap/js/bootstrap.min.js"></script>

<script type="text/javascript">
  $('#login-modal').on('shown.bs.modal', function () {
    $('#username').focus();
  });
</script>

  </body>
</html>
-----5------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----6------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----7------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----8------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----9------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----10------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----11------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----12------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----13------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----14------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----15------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----16------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----17------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----18------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----19------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota-----20------
RESOURCE_EXHAUSTED:Quota is exhausted for: requestcountquota[core@support Enforcement_policy]$

```
# Control Headers and Routing (Deprecated)

```
1. configure an ingress using a gateway

   kubectl apply -f httpbin-gateway.yaml
   kubectl apply -f httpbin-route.yaml

2. Output-producing adapters 
  # Deploy the demo adapter:

  kubectl run keyval --image=gcr.io/istio-testing/keyval:release-1.1 --namespace istio-system --port 9070 --expose

  # Enable the keyval adapter by deploying its template and configuration descriptors:

  kubectl apply -f keyval-template.yaml
  kubectl apply -f keyval.yaml

  # Create a handler for the demo adapter with a fixed lookup table:

   kubectl apply -f demo_handler_adapter.yaml

 # Create an instance for the handler with the user 
  
   kubectl apply -f handler_instance.yaml

3. Request header operations

 # Get ingress route

  export INGRESS_ROUTE=$(oc get route -n istio-system httpbin-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')  
  curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"
  
 # Create a rule for the demo adapter:
  
   kubectl apply -f demo_adapter_rule.yaml

 # Issue a new request to the ingress gateway   

   curl http://httpbin-foo.apps.ocp4.example.com/headers
   <error_detail>key "" not found</error_detail>
   
   curl -Huser:jason http://httpbin-foo.apps.ocp4.example.com/headers
{
  "headers": {
    "Accept": "*/*",
    "Content-Length": "0",
    "Forwarded": "for=192.168.11.70;host=httpbin-foo.apps.ocp4.example.com;proto=http",
    "Host": "httpbin-foo.apps.ocp4.example.com",
    "User": "jason",
    "User-Agent": "curl/7.61.1",
    "User-Group": "admin",
    "X-B3-Parentspanid": "835a1b26f3d9361e",
    "X-B3-Sampled": "1",
    "X-B3-Spanid": "c1b49c70bc05cf42",
    "X-B3-Traceid": "6948a2cdab4d1c03835a1b26f3d9361e",
    "X-Envoy-External-Address": "10.254.0.1",
    "X-Forwarded-Host": "httpbin-foo.apps.ocp4.example.com"
  }
}

4. Modify the rule to rewrite the URI path 

  oc apply -f demo_adapter_rule_status.yaml
  curl -Huser:jason http://httpbin-foo.apps.ocp4.example.com/headers

    -=[ teapot ]=-

       _...._
     .'  _ _ `.
    | ."` ^ `". _,
    \_;`"---"`|//
      |       ;/
      \_     _/
        `"""`
  
6. Cleanup

   kubectl delete rule/keyval handler/keyval instance/keyval adapter/keyval template/keyval -n istio-system
   kubectl delete service keyval -n istio-system
   kubectl delete deployment keyval -n istio-system

```
# Denials and White/Black Listing

```
1. Simple denials


2. Attribute-based whitelists or blacklists


3. IP-based whitelists or blacklists


4. Cleanup

   oc delete -f Denials_White_Black_list

```
