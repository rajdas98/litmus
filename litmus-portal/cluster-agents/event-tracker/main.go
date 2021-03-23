/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/event-tracker/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var (
	gvrdc = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pod",
	}
)

var KubeConfig = os.Getenv("KUBECONFIG")
func GetKubeConfig() (*rest.Config, error) {
	// Use in-cluster config if kubeconfig path is not specified
	if KubeConfig == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", KubeConfig)
}


func main() {
	restConfig, err := GetKubeConfig()
	if err != nil {
			log.Print(err)
	}

	clientSet, err := dynamic.NewForConfig(restConfig)
	if err != nil {
			log.Print(err)
	}

	deploymentRes := schema.GroupVersionResource{Group: "eventtracker.litmuschaos.io", Version: "v1alpha1", Resource: "eventtrackerpolicies"}

	//dc := clientSet.Resource(gvrdc)
	deploymentConfigList, err := clientSet.Resource(deploymentRes).Namespace("default").List(metav1.ListOptions{})

	if err != nil {
		log.Print(err)
	}

	//log.Print(deploymentConfigList)
	var etpl v1alpha1.EventTrackerPolicyList
	data, err := json.Marshal(deploymentConfigList.Items)
	if err != nil {
		log.Print(err)
	}

	err = json.Unmarshal(data, &etpl)
	if err != nil {
		log.Print(err)
	}

	for _, ep := range etpl.Items {
		eventTrackerPolicy, err := clientSet.Resource(deploymentRes).Namespace("default").Get( ep.Name,metav1.GetOptions{})
		if err != nil {
			log.Print(err)
		}

		var etp v1alpha1.EventTrackerPolicy
		data, err := json.Marshal(eventTrackerPolicy.Object)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal(data, &etp)
		if err != nil {
			log.Print(err)
		}

		log.Print(e)
	}
	////_ = clientgoscheme.AddToScheme(scheme)
	//
	//_ = eventtrackerv1alpha1.AddToScheme(scheme)
	//// +kubebuilder:scaffold:scheme
}

//func main() {
//	var metricsAddr string
//	var enableLeaderElection bool
//	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
//	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
//		"Enable leader election for controller manager. "+
//			"Enabling this will ensure there is only one active controller manager.")
//	flag.Parse()
//
//	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
//
//	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
//		Scheme:             scheme,
//		MetricsBindAddress: metricsAddr,
//		Port:               9443,
//		LeaderElection:     enableLeaderElection,
//		LeaderElectionID:   "2b79cec3.litmuschaos.io",
//	})
//	if err != nil {
//		setupLog.Error(err, "unable to start manager")
//		os.Exit(1)
//	}
//
//	log.Print("hi")
//	pod := &eventtrackerv1alpha1.EventTrackerPolicy{}
//	// c is a created client.
//	 err = c.Get(context.Background(), client.ObjectKey{
//		Namespace: "default",
//		Name:      "eventtrackerpolicy-sample",
//	}, pod)
//	if err != nil {
//		log.Print(err)
//	}
//	log.Print(pod)
//
//	if err = (&controllers.EventTrackerPolicyReconciler{
//		Client: mgr.GetClient(),
//		Log:    ctrl.Log.WithName("controllers").WithName("EventTrackerPolicy"),
//		Scheme: mgr.GetScheme(),
//	}).SetupWithManager(mgr); err != nil {
//		setupLog.Error(err, "unable to create controller", "controller", "EventTrackerPolicy")
//		os.Exit(1)
//	}
//	// +kubebuilder:scaffold:builder
//	setupLog.Info("starting manager")
//	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
//		setupLog.Error(err, "problem running manager")
//		os.Exit(1)
//	}
//
//}
