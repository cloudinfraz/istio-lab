#!/bin/sh #

#!/bin/bash
#CA

apps_subdoman=$(oc get route productpage --template '{{ .spec.host }}'|awk -F'.' '{for(i=2;i<=NF;++i) printf $i "."}')
Sub_make_ca() {
sslpath=/tmp/ssl
mkdir -p $sslpath
echo "
[ req ]
req_extensions     = req_ext
distinguished_name = req_distinguished_name
prompt             = no
 
[req_distinguished_name]
commonName=$apps_subdoman
 
[req_ext]
subjectAltName   = @alt_names
 
[alt_names]
DNS.2  =  *.$apps_subdoman
" > $sslpath/cert.cfg

openssl req -x509 -config $sslpath/cert.cfg -extensions req_ext -nodes -days 730 -newkey rsa:2048 -sha256 -keyout $sslpath/tls.key -out $sslpath/tls.crt
oc delete secret tls istio-ingressgateway-certs -n $1  > /dev/null 2>&1
oc create secret tls istio-ingressgateway-certs --cert $sslpath/tls.crt --key $sslpath/tls.key -n $1
rm -rf $sslpath
}

if [ $# != 1 ];then
    echo "Please input the right istio namespace!"
    exit 1
else
   Sub_make_ca $1
   exit 0
fi
