package main

import (
	"flag"
	"fmt"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newResourcesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_, resources, err := clientset.DiscoveryClient.ServerGroupsAndResources()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}

		b := &htmlbuilder{}
		b.html(func(b *htmlbuilder) {
			b.body(func(b *htmlbuilder) {
				b.table(func(b *htmlbuilder) {
					for _, resource := range resources {
						gv := resource.GroupVersion
						for _, apiResource := range resource.APIResources {
							b.tr(func(b *htmlbuilder) {
								k := apiResource.Kind
								txt := fmt.Sprintf("<td>%s</td> <td>%s</td>", gv, k)
								b.WriteString(txt)
							})
						}
					}
				})
			})
		})
		w.Write([]byte(b.String()))
	}
}

func main() {
	kubeconfig := flag.String("kubeconfig", ".kcp/admin.kubeconfig", "Pass a kubeconfig")
	hostport := flag.String("hostport", ":8090", "Pass a host:port on which to listen")
	flag.Parse()

	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Couldn't get in-cluster config: %s", err.Error())
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset := kubernetes.NewForConfigOrDie(config)
	http.HandleFunc("/", newResourcesHandler(clientset))
	http.ListenAndServe(*hostport, nil)
}
