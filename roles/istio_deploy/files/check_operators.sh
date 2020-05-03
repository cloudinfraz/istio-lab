#!/bin/bash
ok_num=$(oc get pods -n openshift-operators |grep 'operator'|grep Running | wc -l)
echo "$ok_num operators are ready"
