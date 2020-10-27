# Authentication and Authorization
istio1.4 AuthorizationPolicy only supports “ALLOW” action. This means that if multiple authorization policies apply to the same workload, the effect is additive.
The v1beta1 policy is not backward compatible and requires a one time conversion. A tool is provided to automate this process. The previous configuration resources ClusterRbacConfig, ServiceRole, and ServiceRoleBinding will not be supported from Istio 1.6 onwards.

 istio1.7+ is better
- authentication policy
- mutual TLS authentication
- sso integretion

# Scenario Setup

```
# case setup
oc adm policy add-scc-to-user anyuid -z httpbin
kubectl create ns foo
kubectl apply -f httpbin.yaml -n foo
kubectl apply -f sleep.yaml -n foo

oc adm policy add-scc-to-user anyuid -z httpbin
kubectl create ns bar
kubectl apply -f httpbin.yaml -n bar
kubectl apply -f sleep.yaml -n bar

oc adm policy add-scc-to-user anyuid -z httpbin
kubectl create ns legacy
kubectl apply -f httpbin.yaml -n legacy
kubectl apply -f sleep.yaml -n legacy

#verify setup by sending an HTTP request with curl from any sleep pod in the namespace foo, bar or legacy

for i in foo bar legacy;
do 
kubectl exec $(kubectl get pod -l app=sleep -n bar -o jsonpath={.items..metadata.name}) -c sleep -n bar -- curl http://httpbin.${i}:8000/ip -s -o /dev/null -w "%{http_code}\n"
done

for from in "foo" "bar" "legacy"; do for to in "foo" "bar" "legacy"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

#  verify mesh authentication policy kubectl get policies.authentication.istio.io --all-namespaces
oc get smmr --all-namespaces

# review crd for openshift 4.5 
oc get crd |grep -i  istio
adapters.config.istio.io                                    2020-08-11T11:20:03Z
attributemanifests.config.istio.io                          2020-08-11T11:20:04Z
authorizationpolicies.security.istio.io                     2020-08-11T11:20:04Z
destinationrules.networking.istio.io                        2020-08-11T11:20:04Z
envoyfilters.networking.istio.io                            2020-08-11T11:20:04Z
gateways.networking.istio.io                                2020-08-11T11:20:04Z
handlers.config.istio.io                                    2020-08-11T11:20:04Z
httpapispecbindings.config.istio.io                         2020-08-11T11:20:04Z
httpapispecs.config.istio.io                                2020-08-11T11:20:04Z
instances.config.istio.io                                   2020-08-11T11:20:04Z
policies.authentication.istio.io                            2020-08-11T11:20:04Z
quotaspecbindings.config.istio.io                           2020-08-11T11:20:03Z
quotaspecs.config.istio.io                                  2020-08-11T11:20:03Z
rbacconfigs.rbac.istio.io                                   2020-08-11T11:20:03Z
rules.config.istio.io                                       2020-08-11T11:20:03Z
serviceentries.networking.istio.io                          2020-08-11T11:20:04Z
servicerolebindings.rbac.istio.io                           2020-08-11T11:20:03Z
serviceroles.rbac.istio.io                                  2020-08-11T11:20:04Z
sidecars.networking.istio.io                                2020-08-11T11:20:04Z
templates.config.istio.io                                   2020-08-11T11:20:04Z
virtualservices.networking.istio.io                         2020-08-11T11:20:04Z

# verify that there are no destination rules
kubectl get destinationrules.networking.istio.io --all-namespaces -o yaml | grep "host:"

# deployment inject 

oc patch deployment httpbin -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n foo
oc patch deployment httpbin -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n bar
oc patch deployment httpbin -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n legacy

oc patch deployment sleep -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n foo
oc patch deployment sleep -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n bar
oc patch deployment sleep -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n legacy

# enable mtls 
for i in foo bar legacy;do oc apply -f httpbin-mtls.yaml -n $i;done
for i in foo bar legacy;do oc apply -f sleep-mtls.yaml -n $i;done

#  add a destination rule to overwrite the TLS setting for httpbin.legacy,this item can't overwrite the policy 
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
 name: "httpbin-legacy"
 namespace: "legacy"
spec:
 host: "httpbin.legacy.svc.cluster.local"
 trafficPolicy:
   tls:
     mode: DISABLE
EOF
# delete httpbin-mtls policy 
oc delete policy httpbin-mtls  -n legacy
for from in "foo" "bar"; do for to in "legacy"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Request from Istio services to Kubernetes API server
TOKEN=$(kubectl describe secret $(kubectl get secrets | grep default-token | cut -f1 -d ' ' | head -1) | grep -E '^token' | cut -f2 -d':' | tr -d ' \t')
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl https://kubernetes.default/api --header "Authorization: Bearer $TOKEN" --insecure -s -o /dev/null -w "%{http_code}\n"

# Namespace-wide policy

kubectl apply -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "default"
  namespace: "legacy"
spec:
  peers:
  - mtls: {}
EOF

# Add corresponding destination rule

 kubectl apply -f - <<EOF
apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "default"
  namespace: "legacy"
spec:
  host: "*.legacy.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
EOF

# Service-specific policy

 kubectl apply -f - <<EOF
apiVersion: authentication.istio.io/v1alpha1
kind: Policy
metadata:
  name: httpbin-mtls
spec:
  peers:
  - mtls:
      mode: STRICT
  targets:
  - name: httpbin
EOF

# Policy precedence

kubectl apply -n legacy -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "overwrite-example"
spec:
  targets:
  - name: httpbin
EOF

# destination rule
kubectl apply -n legacy -f - <<EOF
apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "overwrite-example"
spec:
  host: httpbin.legacy.svc.cluster.local
  trafficPolicy:
    tls:
      mode: DISABLE
EOF

# confirming service-specific policy overrides the namespace-wide policy.

kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo  -- curl http://httpbin.legacy:8000/ip -s -o /dev/null -w "%{http_code}\n"

```
# End-user authentication

```
# create ingress gateway for Scenario
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
  namespace: legacy
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
EOF
# create virtual service

kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
  namespace: legacy
spec:
  hosts:
  - "*"
  gateways:
  - httpbin-gateway
  http:
  - route:
    - destination:
        port:
          number: 8000
        host: httpbin.legacy.svc.cluster.local
EOF
# k8s native export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
# Get ingress route
export INGRESS_ROUTE=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"

1. Enable end-user JWT for httpbin.foo /request authentication policy

 kubectl apply -n legacy -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "jwt-example"
spec:
  targets:
  - name: httpbin
  origins:
  - jwt:
      issuer: "testing@secure.istio.io"
      jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/jwks.json"
  principalBinding: USE_ORIGIN
EOF

$ curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"
401

# Enable the valid JWT token 

export INGRESS_ROUTE=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
TOKEN=$(curl https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/demo.jwt -s)
curl --header "Authorization: Bearer $TOKEN" $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"

# use the script gen-jwt.py to generate new tokens to test with different issuer, audiences, expiry date

wget https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/gen-jwt.py
chmod +x gen-jwt.py
wget https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/key.pem

# perform a test

export INGRESS_ROUTE=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
TOKEN=$(./gen-jwt.py ./key.pem --expire 5)
for i in `seq 1 10`; do curl --header "Authorization: Bearer $TOKEN" $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"; sleep 1; done

# Disable End-user authentication for specific paths

 kubectl apply -n legacy -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "jwt-example"
spec:
  targets:
  - name: httpbin
  origins:
  - jwt:
      issuer: "testing@secure.istio.io"
      jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/jwks.json"
      trigger_rules:
      - excluded_paths:
        - exact: /user-agent
  principalBinding: USE_ORIGIN
EOF

curl $INGRESS_ROUTE/user-agent -s -o /dev/null -w "%{http_code}\n"
curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n"


#. Enable End-user authentication for specific paths

 kubectl apply -n legacy -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "jwt-example"
spec:
  targets:
  - name: httpbin
  origins:
  - jwt:
      issuer: "testing@secure.istio.io"
      jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/jwks.json"
      trigger_rules:
      - included_paths:
        - exact: /ip
  principalBinding: USE_ORIGIN
EOF

curl $INGRESS_ROUTE/user-agent -s -o /dev/null -w "%{http_code}\n"
curl $INGRESS_ROUTE/ip -s -o /dev/null -w "%{http_code}\n"

# Confirm it’s allowed to access the path /ip with a valid JWT token

TOKEN=$(curl https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/demo.jwt -s)
curl --header "Authorization: Bearer $TOKEN" $INGRESS_ROUTE/ip -s -o /dev/null -w "%{http_code}\n"

# End-user authentication with mutual TLS

 kubectl apply -n legacy -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "jwt-example"
spec:
  targets:
  - name: httpbin
  peers:
  - mtls: {}
  origins:
  - jwt:
      issuer: "testing@secure.istio.io"
      jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/jwks.json"
  principalBinding: USE_ORIGIN
EOF

# add a destination rule

kubectl apply -f - <<EOF
apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "httpbin"
  namespace: "legacy"
spec:
  host: "httpbin.legacy.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
EOF

TOKEN=$(curl https://raw.githubusercontent.com/istio/istio/release-1.4/security/tools/jwt/samples/demo.jwt -s)
kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl http://httpbin.legacy:8000/ip -s -o /dev/null -w "%{http_code}\n" --header "Authorization: Bearer $TOKEN"

# Cleanup part 3

kubectl -n legacy delete policy jwt-example
kubectl -n legacy delete destinationrule httpbin
kubectl delete ns foo bar legacy

```
# Authorization for HTTP traffic

```
1. Configure access control for workloads using HTTP traffic

# Run the following command to create a deny-all policy in the default namespace

  kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
  namespace: bookinfo
spec:
  {}
EOF

# Run the following command to create a productpage-viewer policy to allow access with GET method t

 kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "productpage-viewer"
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: productpage
  rules:
  - to:
    - operation:
        methods: ["GET"]
EOF

# Run the following command to create the details-viewer policy to allow the productpage workload, which issues requests using
 the cluster.local/ns/bookinfo/sa/bookinfo-productpage service account, to access the details

  kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "details-viewer"
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: details
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/bookinfo/sa/bookinfo-productpage"]
    to:
    - operation:
        methods: ["GET"]
EOF

# Run the following command to create a policy reviews-viewer to allow the productpage workload

 kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "reviews-viewer"
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: reviews
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/bookinfo/sa/bookinfo-productpage"]
    to:
    - operation:
        methods: ["GET"]
EOF

# Run the following command to create the ratings-viewer policy to allow the reviews workload

  kubectl apply -f - <<EOF
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "ratings-viewer"
  namespace: bookinfo
spec:
  selector:
    matchLabels:
      app: ratings
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/bookinfo/sa/bookinfo-reviews"]
    to:
    - operation:
        methods: ["GET"]
EOF

# clean up

kubectl -n bookinfo delete authorizationpolicy.security.istio.io/deny-all
kubectl -n bookinfo delete authorizationpolicy.security.istio.io/productpage-viewer
kubectl -n bookinfo delete authorizationpolicy.security.istio.io/details-viewer
kubectl -n bookinfo delete authorizationpolicy.security.istio.io/reviews-viewer
kubectl -n bookinfo delete authorizationpolicy.security.istio.io/ratings-viewer

```
# Authorization for TCP traffic

```


```

# Keycloak/JWT Integretion

```
1. Keycloak setup

# Create a security realm

# Create istio client id

# Create a role

# Create a user to generate token

# Assign the user to role

# Mapper Type: Audiece included: istio,  open client istio and select the tab Mappers and press on the button Create

# Mapper Type: User Realm Role, Token Claim Name:roles, Claim JSON Type: String  open client istio and select the tab Mappers and press on the button Create
 - Add to ID token: on
 - Add to access token: on
 - Add to userinfo: on 

2.  End-user authentication with jwt token

 kubectl apply -n foo -f - <<EOF
apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "jwt-keycloak-mtls"
spec:
  targets:
    - name: httpbin
  peers:
    - mtls: {}
  origins:
  - jwt:
      audiences:
        - istio
      issuer: "https://sso.ocp4.example.com/auth/realms/istio"
      jwksUri: "https://sso.apps.ocp4.example.com/auth/realms/istio/protocol/openid-connect/certs"
  principalBinding: USE_ORIGIN
EOF

3. Call api for a valid token
export INGRESS_ROUTE=$(oc get route -n istio-system httpbin-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
export TOKEN=$(curl -sk -d "username=test&password=test123&grant_type=password&client_id=istio&client_secret=2a2dcc91-4637-47f8-96cf-6a7eb7125613"   https://sso.apps.ocp4.example.com/auth/realms/istio/protocol/openid-connect/token   | jq  ".access_token")

curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n" --header "Authorization: Bearer $TOKEN"

kubectl exec $(kubectl get pod -l app=sleep -n foo -o jsonpath={.items..metadata.name}) -c sleep -n foo -- curl http://httpbin.foo:8000/ip -s -o /dev/null -w "%{http_code}\n"

 4. Configure groups-based authorization

 kubectl apply -n foo -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
spec:
  {}
EOF

# allows users in role: users to access it with GET method:
cat <<EOF | kubectl apply -n foo -f -
apiVersion: "security.istio.io/v1beta1"
kind: "AuthorizationPolicy"
metadata:
  name: "httpbin-viewer"
spec:
  selector:
    matchLabels:
      app: httpbin
  rules:
  - to:
    - operation:
        methods: ["GET"]
        paths: ["/ip"]
    when:
    - key: request.auth.claims[preferred_username]
      values: ["test"]
EOF

5. Authorization Policy

 curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n" --header "Authorization: Bearer $TOKEN"
403
 curl $INGRESS_ROUTE/ip -s -o /dev/null -w "%{http_code}\n" --header "Authorization: Bearer $TOKEN"
200

6. troubleshooting for certification

# Download the wildcard certificate of your OpenShift cluster with the following command:

openssl s_client \
  -showcerts \
  -servername sso.ocp4.example.com \
  -connect sso.ocp4.example.com:443 </dev/null 2>/dev/null \
  | openssl x509 -outform PEM >openshift-wildcard.pem

# Create a secret with the certificate. The filename in the secret has to be extra.pem:

oc -n istio-system create secret generic openshift-wildcard \
--from-file=extra.pem=openshift-wildcard.pem


# mount the volume to  the Istio Pilot pod:

oc -n istio-system set volumes deployment/istio-pilot \
  --add \
  --name=extracacerts \
  --mount-path=/cacerts \
  --secret-name=openshift-wildcard \
  --containers=discovery

```
