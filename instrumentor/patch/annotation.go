package patch

import (
	"context"
	"github.com/logzio/app-type-detector/api/v1alpha1"
	"github.com/logzio/app-type-detector/common/consts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	goclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const LogzioApplicationTypeAnnotation = "logzio/application_type"
const SkipAnnotation = "logzio/skip_app_detection"

var PodOwnedLabels = []string{
	"app",
	"app.kubernetes.io/name",
}

var IgnoredNamespaces = []string{"kube-system", "local-path-storage", "gatekeeper-system", consts.DefaultMonitoringNamespace}
var clusterClient *goclient.CoreV1Client
var annotationPatcher = &AnnotationPatcher{}

type AnnotationPatcher struct{}

func (d *AnnotationPatcher) Patch(ctx context.Context, detected *v1alpha1.AppDetector, object client.Object) error {

	if d.shouldPatch(object.GetAnnotations(), object.GetNamespace()) {
		kubeClient, err := getKubeClient()
		if err != nil {
			return err
		}
		podClient := kubeClient.Pods(object.GetNamespace())
		childPods, err := podClient.List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, pod := range childPods.Items {
			if podOwnedByObject(pod.GetLabels(), object.GetName()) && d.shouldPatch(pod.Annotations, object.GetNamespace()) {
				if pod.Annotations == nil {
					pod.Annotations = make(map[string]string)
				}
				pod.Annotations[LogzioApplicationTypeAnnotation] = string(detected.Spec.Applications[0].Application)
				_, err := podClient.Update(ctx, &pod, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getKubeClient() (*goclient.CoreV1Client, error) {
	if clusterClient == nil {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err := kubeConfig.ClientConfig()
		if err != nil {
			log.Println("Error creating kubernetes client creating config")
			return nil, err
		}
		clusterClient = goclient.NewForConfigOrDie(config)
		return clusterClient, err
	}
	return clusterClient, nil
}

func podOwnedByObject(labels map[string]string, name string) bool {
	for _, podAppLabel := range PodOwnedLabels {
		if val, exists := labels[podAppLabel]; exists {
			if val == name {
				return true
			}
		}
	}
	return false
}

func (d *AnnotationPatcher) shouldPatch(annotations map[string]string, namespace string) bool {
	for k, v := range annotations {
		if (k == SkipAnnotation && v == "true") || k == LogzioApplicationTypeAnnotation {
			return false
		}
	}
	for _, ns := range IgnoredNamespaces {
		if namespace == ns {
			return false
		}
	}
	return true
}
