Red Hat Service Mesh Lab 
-----------------------------------------------------------------------------

- ansible 2.6+
- ocp 4.3
- istio 1.0.9+

Istio learning series
----------------------------------------------------------------------------
- multi tenancy
- mTLS and authorization policies
- Istio Reliability: Retries, Timeouts and Circuit Breaker
- Observability
- mixer ? 

## Initial Site Setup

First we configure the host file for oc client 
inventory file, grouped by their purpose:
```
 [all:vars]
ansible_ssh_user=chunzhan
ansible_ssh_pass=1SxxxxxlEh
ansible_ssh_pass=1S6wC04NPlEaa

[istio-client]
clientvm.2b7a.internal

```
### group var setting
```
istio_app1: bookinfo
ISTIO_NS:  bookretail-istio-system
#elasticsearch version
es_channel: "4.3"
#service mesh version
sm_channel: "1.0"
cluster_username: admin
cluster_password: r3dhxxxx4
cluster_url: https://api.cluster-2b7a.2dasfds.sandbox1314.opentlc.com:6443
apps_subdomain: "apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com"

```

ansible playbook including 3 roles:
- bookinfo deployment
- istio-system deployment
- mtls enable

```
├── group_vars
├── learning_notes
├── roles
│   ├── bookinfo_deploy
│   │   ├── files
│   │   ├── tasks
│   │   └── templates
│   ├── istio_deploy
│   │   ├── files
│   │   ├── tasks
│   │   └── templates
│   ├── istio_test
│   │   ├── files
│   │   ├── tasks
│   │   └── templates
│   └── mtls_enable
│       ├── files
│       ├── tasks
│       └── templates
└── test_scenario
    ├── Authentication_policy
    ├── Circuit_breaking
    ├── Fault_injection
    ├── Mirroring
    ├── Request_timeout
    └── Traffic_shifting
        └── tcp

```
deploy Red Hat Service Mesh and the example app: bookinfo 

	ansible-playbook -i hosts deploy-istio.yaml
	ansible-playbook -i hosts deploy-bookinfo.yaml

enable mtls for bookinfo app

	ansible-playbook -i hosts deploy-mtls-site.yaml


bookinfo verification 

```
export INGRESS_GATEWAY=$(oc get route -n istio-system istio-ingressgateway -o jsonpath='{.items[*]}{.spec.host}')
curl -s "http://${INGRESS_GATEWAY}/productpage" | grep -o "<title>.*</title>"

```

mtls verification 
```
export ISTIO_NS=istio-system
export APP_NS=bookinfo
export INGRESS_ROUTE=$(oc get route -n istio-system bookinfo-service-gateway -o jsonpath='{.items[*]}{.spec.host}')
export ISTIO_INGRESSGATEWAY_POD=$(oc get pods -n istio-system |grep istio-egressgateway | awk '{print $1}')

curl -s -k "https://${INGRESS_ROUTE}/productpage" | grep -o "<title>.*</title>"

istioctl -n $ISTIO_NS -i $ISTIO_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} productpage.$APP_NS.svc.cluster.local
istioctl -n $ISTIO_NS -i $ISTIO_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} reviews.$APP_NS.svc.cluster.local
istioctl -n $ISTIO_NS -i $ISTIO_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} ratings.$APP_NS.svc.cluster.local
istioctl -n $ISTIO_NS -i $ISTIO_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} details.$APP_NS.svc.cluster.local

```
Istio testing 

```

```
# istio-lab
