package main

import (
"context"
"log"
"os"
"fmt"

v1alpha3Spec "istio.io/api/networking/v1alpha3"
"istio.io/client-go/pkg/apis/networking/v1alpha3"
versionedclient "istio.io/client-go/pkg/clientset/versioned"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/client-go/tools/clientcmd"
)

func main() {

kubeconfig := os.Getenv("KUBECONFIG")
namespace := os.Getenv("NAMESPACE")

if len(kubeconfig) == 0 || len(namespace) == 0 {
log.Fatalf("Environment variables KUBECONFIG and NAMESPACE need to be set")
}

restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
log.Fatalf("Failed to create k8s rest client: %s", err)
}

ic, err := versionedclient.NewForConfig(restConfig)
if err != nil {
log.Fatalf("Failed to create istio client: %s", err)
}
var host, gateway []string
if len(host) == 0 {
	host = append(host,"abc.com")
}

if len(gateway) == 0 {
	gateway = append(gateway, "test")
}
	GateWayCrd := &v1alpha3.Gateway{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "Gateway",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
		},
		// attach to use the istioingressgateawy (default)
		Spec: v1alpha3Spec.Gateway{
			Selector: map[string]string{"istio": "ingressgateway"},

			Servers: []*v1alpha3Spec.Server{{
				Port: &v1alpha3Spec.Port{
					Number:   80,
					Name:     "http",
					Protocol: "HTTP",
				},
				Hosts: []string{"*"},
			}},
		},
	}
	// Test Gateway
	gw, err := ic.NetworkingV1alpha3().Gateways(namespace).Create(context.TODO(), GateWayCrd, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create Gateway in %s namespace: %s", namespace, err)
	}
	fmt.Println(gw)

	virtualServiceCrd := &v1alpha3.VirtualService{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "Virtualservice",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: v1alpha3Spec.VirtualService{
			Hosts: host,
			Gateways: gateway,
			Http: []*v1alpha3Spec.HTTPRoute{{
				Match: []*v1alpha3Spec.HTTPMatchRequest{{
					Uri: &v1alpha3Spec.StringMatch{
						MatchType: &v1alpha3Spec.StringMatch_Prefix{
							Prefix: "/",
						},
					},
				}},
				Route: []*v1alpha3Spec.HTTPRouteDestination{{
					Destination: &v1alpha3Spec.Destination{
						Port: &v1alpha3Spec.PortSelector{
							Number: 9080,
						},
						Host: "productpage",
					},
				}},
			}},
		},
	}

    vs, err := ic.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), virtualServiceCrd, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create VirtualService in %s namespace: %s", namespace, err)
	}
	fmt.Println(vs)
}