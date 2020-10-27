# Traffic Management

- Request Routing
- Fault Injection
- Traffic Shifting
- TCP Traffic Shifting
- Request Timeouts
- Circuit Breaking

# Request Routing
```
# istio gateway access
export GATEWAY_URL=$INGRESS_HOST
curl -s "http://${GATEWAY_URL}/productpage" | grep -o "<title>.*</title>"

```
# Fault Injection
```
# istio gateway access
export GATEWAY_URL=$INGRESS_HOST
curl -s "http://${GATEWAY_URL}/productpage" | grep -o "<title>.*</title>"

```
# Circuit Breaking
```
# perform circuit breaking testing

export FORTIO_POD=$(kubectl get pods -lapp=fortio -o 'jsonpath={.items[0].metadata.name}')
kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio curl -quiet http://httpbin:8000/get
# Tripping the circuit breaker
# Call the service with two concurrent connections (-c 2) and send 20 requests (-n 20):
kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio load -c 2 -qps 0 -n 20 -loglevel Warning http://httpbin:8000/get

# Bring the number of concurrent connections up to 3:
kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio load -c 3 -qps 0 -n 30 -loglevel Warning http://httpbin:8000/get

# Query the istio-proxy stats to see more:
kubectl exec "$FORTIO_POD" -c istio-proxy -- pilot-agent request GET stats | grep httpbin | grep pending

```

# Mirroring

```
# Creating a default routing policy
1. Create a default route rule to route all traffic to v1 of the service:

2. Send some traffic to the service:
export SLEEP_POD=$(kubectl get pod -l app=sleep -o jsonpath={.items..metadata.name})
kubectl exec "${SLEEP_POD}" -c sleep -- curl -s http://httpbin:8000/headers 

3. Check the logs for v1 and v2 of the httpbin pods. You should see access log entries for v1 and none for v2:
 export V1_POD=$(kubectl get pod -l app=httpbin,version=v1 -o jsonpath={.items..metadata.name})
 kubectl logs "$V1_POD" -c httpbin

 export V2_POD=$(kubectl get pod -l app=httpbin,version=v2 -o jsonpath={.items..metadata.name})
 kubectl logs "$V2_POD" -c httpbin

# Mirroring traffic to v2
1. Change the route rule to mirror traffic to v2:

2. Send in traffic:

 kubectl exec "${SLEEP_POD}" -c sleep -- curl -s http://httpbin:8000/headers
 kubectl logs "$V1_POD" -c httpbin
 kubectl logs "$V2_POD" -c httpbin


=======
路由自动创建验证
```

 

```

负载均衡验证
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

超时/重试/故障注入
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

限流/熔断
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

应用监控指标
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

调用链追踪
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

蓝绿和灰度发布
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

流量镜像
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```

压力测试
```

curl -kv https://productpage-service.apps.cluster-2b7a.2b7a.sandbox1314.opentlc.com/productpage

```
