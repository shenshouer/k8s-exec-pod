package main

import (
	"fmt"
	"log"
	"os"

	"k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	coreclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	kubeConfig := "/Users/shenshouer/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err)
	}

	// config.GroupVersion = &unversioned.GroupVersion{Group: "", Version: "v1"}
	// config.NegotiatedSerializer = api.Codecs

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podClient := kubeClient.Core()
	podlist, err := podClient.Pods("test").List(api.ListOptions{})
	if err != nil {
		panic(err)
	}

	// rClient, err := restclient.RESTClientFor(config)
	// if err != nil {
	// 	panic(err)
	// }

	// rClient := podClient.RESTClient()
	for _, pod := range podlist.Items {
		fmt.Println(pod.Name)
		for _, c := range pod.Spec.Containers {
			log.Println("ContainerName:", c.Name)

			execPod(pod, c, config, podClient)

			// req := rClient.Post().
			// 	Resource("pods").
			// 	Name(pod.Name).
			// 	Namespace(pod.Namespace).
			// 	SubResource("exec").
			// 	Param("container", c.Name)

			// req = req.SetHeader("X-Stream-Protocol-Version", "v2.channel.k8s.io").SetHeader("X-Stream-Protocol-Version", "channel.k8s.io")
			// req.VersionedParams(&api.PodExecOptions{
			// 	Container: c.Name,
			// 	Command:   []string{"nginx", "-t"},
			// 	Stdin:     true,
			// 	Stdout:    true,
			// 	Stderr:    true,
			// 	TTY:       false,
			// }, api.ParameterCodec)

			// res := req.Do()
			// ro, err := res.Get()
			// if err != nil {
			// 	log.Println("[ERROR] ", err)
			// 	continue
			// }
			// log.Println(ro)
		}
	}

	// time.Sleep(30 * time.Second)
}

func execPod(pod api.Pod, container api.Container, config *restclient.Config, podClient coreclient.CoreInterface) {
	options := &ExecOptions{
		StreamOptions: StreamOptions{
			In: os.Stdin, Out: os.Stdout, Err: os.Stdout,
			TTY: false, Stdin: false,
			PodName:       pod.Name,
			ContainerName: container.Name,
			Namespace:     pod.Namespace,
		},
		PodClient: podClient,
		Config:    config,
		Executor:  &DefaultRemoteExecutor{},
	}

	// options.Complete(cmdutil.NewFactory(nil), nil, []string{"for i in {1..30};do echo $i;sleep 1;done"}, 1)
	options.Complete([]string{"nginx", "-t"})
	if err := options.Validate(); err != nil {
		log.Println("[ERROR] ", err)
	}

	if err := options.Run(); err != nil {
		log.Println("[ERROR] ", err.Error(), (err == nil))
	}
}
