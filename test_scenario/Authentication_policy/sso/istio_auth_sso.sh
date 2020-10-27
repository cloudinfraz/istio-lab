#!/bin/bash
NS=foo
oc -n -$NS  apply -f httpbin.yaml

oc -n -$NS apply -f sleep.yaml

oc -n -$NS patch deployment httpbin -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n foo
oc -n -$NS patch deployment sleep -p '{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject": "true"}}}}}' -n foo

oc -n -$NS apply -f httpbin-gateway.yaml
oc -n -$NS apply -f foo-dr.yaml
oc -n -$NS apply -f httpbin-mtls.yaml
oc -n -$NS apply -f sleep-mtls.yaml
oc -n -$NS apply -f jwt-keycloak-mtls.yaml
oc -n -$NS apply -f deny-all.yaml
oc -n -$NS apply -f httpbin-AuthorizationPolicy.yaml

oc -n istio-system -f httpbin-route.yaml

export INGRESS_ROUTE=$(oc get route -n istio-system httpbin-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
export TOKEN=$(curl -sk -d "username=test&password=test123&grant_type=password&client_id=istio&client_secret=2a2dcc91-4637-47f8-96cf-6a7eb7125613"   https://sso.apps.ocp4.example.com/auth/realms/istio/protocol/openid-connect/token   | jq  ".access_token")

curl $INGRESS_ROUTE/headers -s -o /dev/null -w "%{http_code}\n" --header "Authorization: Bearer $TOKEN"
