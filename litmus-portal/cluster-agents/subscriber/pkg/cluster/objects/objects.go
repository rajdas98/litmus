package objects

import (
	"encoding/json"
	"errors"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/k8s"
	"github.com/litmuschaos/litmus/litmus-portal/cluster-agents/subscriber/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"strings"
	v1 "k8s.io/api/apps/v1"
)

//GetKubernetesObjects is used to get the Kubernetes Object details according to the request type
func GetKubernetesObjects(request types.KubeObjRequest) ([]*types.KubeObject, error) {
	conf, err := k8s.GetKubeConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}

	resourceType := schema.GroupVersionResource{
		Group:    request.KubeGVRRequest.Group,
		Version:  request.KubeGVRRequest.Version,
		Resource: request.KubeGVRRequest.Resource,
	}
	_, dynamicClient, err := k8s.GetDynamicAndDiscoveryClient()
	if err != nil {
		return nil, err
	}
	var ObjData []*types.KubeObject
	namespace, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(namespace.Items) > 0 {
		for _, namespace := range namespace.Items {
			podList, err := getObjectDataByNamespace(namespace.GetName(), dynamicClient, resourceType)
			if err != nil {
				return nil, err
			}
			KubeObj := &types.KubeObject{
				Namespace: namespace.GetName(),
				Data:      podList,
			}
			ObjData = append(ObjData, KubeObj)
		}
		kubeData, _ := json.Marshal(ObjData)
		var kubeObjects []*types.KubeObject
		err := json.Unmarshal(kubeData, &kubeObjects)
		if err != nil {
			return nil, err
		}
		return kubeObjects, nil
	} else {
		return nil, errors.New("No namespace available")
	}
}

//GetObjectDataByNamespace uses dynamic client to fetch Kubernetes Objects data.
func getObjectDataByNamespace(namespace string, dynamicClient dynamic.Interface, resourceType schema.GroupVersionResource) ([]types.ObjectData, error) {
	list, err := dynamicClient.Resource(resourceType).Namespace(namespace).List(metav1.ListOptions{})
	var kubeObjects []types.ObjectData
	if err != nil {
		return kubeObjects, nil
	}
	var newXyz xyz
	newXyz.ResourceNamespace = namespace

	if err != nil {
		return nil, err
	}

	var tmpObject []unstructured.Unstructured
	for _, list := range list.Items {
		//obj, err:= k8s.ApplyRequest("get", &list)
		obj, err := dynamicClient.Resource(resourceType).Get(list.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		tmpObject = append(tmpObject, *obj)
	}

	newXyz.Resources = append(newXyz.Resources, resource{
		ResourceType: resourceType.Resource,
		Objects: tmpObject,
	})

	// graphql-server
	// 1. make struct of the response from the subscriber
	// 2. Unmarshal
	// 3. Follow the
	for _, xyz := range newXyz.Resources{
		for _, obj := range xyz.Objects{

			mar, err := json.Marshal(obj)
			if err != nil {

			}

			var in interface{}
			json.Unmarshal(mar, &in)

			if strings.ToLower(xyz.ResourceType) == "deployment" {
				newDep := in.(v1.StatefulSet)

			}
		}
	}

	//for _, list := range list.Items {
	//	listInfo := types.ObjectData{
	//		Name:                    list.GetName(),
	//		UID:                     list.GetUID(),
	//		Namespace:               list.GetNamespace(),
	//		APIVersion:              list.GetAPIVersion(),
	//		CreationTimestamp:       list.GetCreationTimestamp(),
	//		TerminationGracePeriods: list.GetDeletionGracePeriodSeconds(),
	//		Labels:                  list.GetLabels(),
	//	}
	//	kubeObjects = append(kubeObjects, listInfo)
	//}
	return kubeObjects, nil
}

// TODO- Resouce

type xyz struct {
	ResourceNamespace string
	Resources []resource
}


type resource struct {
	ResourceType string
	Objects []unstructured.Unstructured
}





/*
	one time request
 		resourcetype
 */