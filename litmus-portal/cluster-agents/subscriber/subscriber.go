package main

import (
	"encoding/json"
	"flag"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/requests"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/workflow"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/k8s"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/types"
	"github.com/sirupsen/logrus"
)

var (
	clusterData = map[string]string{
		"ACCESS_KEY":           os.Getenv("ACCESS_KEY"),
		"CLUSTER_ID":           os.Getenv("CLUSTER_ID"),
		"SERVER_ADDR":          os.Getenv("SERVER_ADDR"),
		"IS_CLUSTER_CONFIRMED": os.Getenv("IS_CLUSTER_CONFIRMED"),
		"AGENT_SCOPE":          os.Getenv("AGENT_SCOPE"),
		"COMPONENTS":           os.Getenv("COMPONENTS"),
	}

	err error
)

func init() {
	logrus.Info("Go Version: ", runtime.Version())
	logrus.Info("Go OS/Arch: ", runtime.GOOS, "/", runtime.GOARCH)

	k8s.KubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	// check agent component status
	err := k8s.CheckComponentStatus(clusterData["COMPONENTS"])
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("all components live...starting up subscriber")

	isConfirmed, newKey, err := k8s.IsClusterConfirmed()
	if err != nil {
		logrus.Fatal(err)
	}

	if isConfirmed == true {
		clusterData["ACCESS_KEY"] = newKey
	} else if isConfirmed == false {
		clusterConfirmByte, err := k8s.ClusterConfirm(clusterData)
		if err != nil {
			logrus.Fatal(err)
		}

		var clusterConfirmInterface types.Payload
		err = json.Unmarshal(clusterConfirmByte, &clusterConfirmInterface)
		if err != nil {
			logrus.Fatal(err)
		}

		if clusterConfirmInterface.Data.ClusterConfirm.IsClusterConfirmed == true {
			clusterData["ACCESS_KEY"] = clusterConfirmInterface.Data.ClusterConfirm.NewAccessKey
			clusterData["IS_CLUSTER_CONFIRMED"] = "true"
			_, err = k8s.ClusterRegister(clusterData)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Info(clusterData["CLUSTER_ID"] + " has been confirmed")
		} else {
			logrus.Info(clusterData["CLUSTER_ID"] + " hasn't been confirmed")
		}
	}

}

func main() {
	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal)
	stream := make(chan types.WorkflowEvent, 10)

	//start workflow event watcher
	workflow.WorkflowEventWatcher(stopCh, stream)

	//streams the event data to graphql server
	go workflow.WorkflowUpdates(clusterData, stream)

	// listen for cluster actions
	go requests.ClusterConnect(clusterData)

	signal.Notify(sigCh, os.Kill, os.Interrupt)
	<-sigCh
	close(stopCh)
	close(stream)
}
