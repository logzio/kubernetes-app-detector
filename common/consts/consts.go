package consts

import "errors"

const (
	AppDetectorContainerAnnotationKey = "logzio/app-detection-pod"
	CurrentNamespaceEnvVar            = "CURRENT_NS"
	DefaultMonitoringNamespace        = "monitoring"
	KubeSystemNamespace               = "kube-system"
	GateKeeperSystemNamespace         = "gatekeeper-system"
	LocalPathStorageNamespace         = "local-path-storage"
)

var (
	PodsNotFoundErr = errors.New("could not find a ready pod")
)
