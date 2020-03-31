Red Hat Service Mesh Lab 
-----------------------------------------------------------------------------

- ansible 2.6+
- ocp 4.3
- istio 1.0.9+

This is istio learning series

### Initial Site Setup

First we configure the host file for oc client 
inventory file, grouped by their purpose:
```
 [all:vars]
ansible_ssh_user=chunzhan
ansible_ssh_pass=1S6wC04NPlEaa

[istio-client]
clientvm.2b7a.internal

```
ansible playbook including 3 roles:
- bookinfo deployment
- istio-system deployment
- mtls enable

```
├── deploy-istio-bookinfo.yaml
├── deploy-mtls-site.yaml
├── group_vars
│   └── all
├── hosts
├── LICENSE.md
├── README.md
└── roles
    ├── bookinfo_deploy
    │   ├── files
    │   └── tasks
    ├── istio_deploy
    │   ├── files
    │   ├── tasks
    │   └── templates
    └── mtls_enable
        ├── files
        ├── tasks
        └── templates

```
deploy Red Hat Service Mesh and the example app: bookinfo 

		ansible-playbook -i hosts deploy-istio-bookinfo.yaml

enable mtls for bookinfo app

                ansible-playbook -i hosts deploy-mtls-site.yaml

mtls verification 
```
istioctl -n $SM_CP_NS -i $SM_CP_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} productpage.$ERDEMO_NS.svc.cluster.local
HOST:PORT                                       STATUS     SERVER     CLIENT     AUTHN POLICY                          DESTINATION RULE
productpage.bookinfo.svc.cluster.local:9080     OK         mTLS       mTLS       productpage-service-mtls/bookinfo     productpage-client-mtls/bookinfo

 istioctl -n $SM_CP_NS -i $SM_CP_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} reviews.$ERDEMO_NS.svc.cluster.local
HOST:PORT                                   STATUS     SERVER     CLIENT     AUTHN POLICY                      DESTINATION RULE
reviews.bookinfo.svc.cluster.local:9080     OK         mTLS       mTLS       reviews-service-mtls/bookinfo     reviews-client-mtls/bookinfo

istioctl -n $SM_CP_NS -i $SM_CP_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} ratings.$ERDEMO_NS.svc.cluster.local
HOST:PORT                                   STATUS     SERVER     CLIENT     AUTHN POLICY                      DESTINATION RULE
ratings.bookinfo.svc.cluster.local:9080     OK         mTLS       mTLS       ratings-service-mtls/bookinfo     ratings-client-mtls/bookinfo

istioctl -n $SM_CP_NS -i $SM_CP_NS authn tls-check ${ISTIO_INGRESSGATEWAY_POD} details.$ERDEMO_NS.svc.cluster.local
HOST:PORT                                   STATUS     SERVER     CLIENT     AUTHN POLICY                      DESTINATION RULE
details.bookinfo.svc.cluster.local:9080     OK         mTLS       mTLS       details-service-mtls/bookinfo     details-client-mtls/bookinfo


curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage|

```
### group var setting 

```
istio_app1: bookinfo
istio_ns:  bookretail-istio-system
#elasticsearch version
es_channel: "4.3"
#service mesh version
sm_channel: "1.0"
cluster_username: admin
cluster_password: r3dhxxxx4
cluster_url: https://api.cluster-2b7a.2b7a.sandbox1314.opentlc.com:6443
apps_subdomain: "apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com"

```
# istio-lab
