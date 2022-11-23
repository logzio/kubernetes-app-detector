**This project is in Beta and subject to breaking changes**

Application and auto instrumentation microservice for kubernetes.  Inspired by odigos.io
This project is built for [logzio-fluentd kuberentes helm chart](https://github.com/logzio/logzio-helm/tree/master/charts/fluentd).
When deployed, a manager and a proxy pod are created.
If the exe or commandline of a pod contains one of the applications in the supported application list,
the application name will be added as a value to "logzio/application_type" annotation of the pod.

**Notes**
* Supported types for detection: Deployment, StatefulSet
* In order to skip detection for a specific Deployment/StatefulSet, the following annotation should be added:
  `logzio/skip_app_detection: true`.
* Ignored namespaces for app detection:
  * kube-system 
  * local-path-storage
  * gatekeeper-system
  * monitoring


#### Development
Build:
```
export TAG=<your-tag>>
make build-images
```


####Chnagelog

**v0.0.1**:
 * Initial release - v0.0.1

