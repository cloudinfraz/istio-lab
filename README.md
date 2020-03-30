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
ansible_ssh_user=chunzhan-redhat.com
ansible_ssh_pass=1S6wC04NPlEh

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
    │   │   └── bookinfo.yaml
    │   └── tasks
    │       ├── add-control-plane-ips.yaml
    │       ├── authenticate.yaml
    │       ├── deploy-bookinfo.yaml
    │       └── main.yaml
    ├── istio_deploy
    │   ├── files
    │   │   ├── jo-sub.yaml
    │   │   └── ko-sub.yaml
    │   ├── tasks
    │   │   ├── add-control-plane-ips.yaml
    │   │   ├── authenticate.yaml
    │   │   ├── deploy-Elasticsearch-operator.yaml
    │   │   ├── deploy-Jaeger-operator.yaml
    │   │   ├── deploy-Kiali-operator.yaml
    │   │   ├── deploy-Mesh-operator.yaml
    │   │   ├── install-istio.yaml
    │   │   └── main.yaml
    │   └── templates
    │       ├── eo-sub.yaml.j2
    │       ├── service-mesh.yaml.j2
    │       └── sm-sub.yaml.j2
    └── mtls_enable
        ├── files
        │   ├── auto_inject.sh
        │   ├── eo-sub.yaml
        │   ├── jo-sub.yaml
        │   ├── ko-sub.yaml
        │   ├── self-wild-certificate.sh
        │   ├── set_probe.sh
        │   └── sm-sub.yaml
        ├── tasks
        │   ├── authenticate.yaml
        │   ├── configure-istio.yaml
        │   ├── deploy-bookinfo-injection.yaml
        │   ├── deploy-self-signed-certificate.yaml
        │   └── main.yaml
        └── templates
            ├── bookinfo-service-destinationrule.yaml.j2
            ├── bookinfo-service-gateway.yaml.j2
            ├── bookinfo-service-policy.yaml.j2
            ├── bookinfo-service-virtualservice.yaml.j2
            ├── service-mesh.yaml.j2
            └── wildcard-gateway.yaml.j2
```
deploy Red Hat Service Mesh and the example app: bookinfo 

		ansible-playbook -i hosts deploy-istio-bookinfo.yaml

enable mtls for bookinfo app

                ansible-playbook -i hosts deploy-mtls-site.yaml

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
