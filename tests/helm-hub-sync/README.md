# JFrog helm-hub-sync Helm Chart

A tool to synchronize [Helm Hub](https://github.com/helm/hub) repositories with [JFrog Artifactory](https://jfrog.com/artifactory/)

## Why do I need this

That's a really good question to begin with! [Helm Hub](https://hub.helm.sh) with the new UI is super awesome, but it only can be used as distributed public repository to search for charts in UI.
You might still want to have a single central location where you can find the Helm charts for your organization. `helm-hub-sync` helps you maintain a virtual repository in Artifactory that can be that single source of truth, using the configuration from Helm Hub.

## Prerequisites Details

* Kubernetes 1.10+

## Chart Details

This chart will do the following:

* Deploy JFrog Helm-hub-sync

## Requirements

- A running Kubernetes cluster
- A running Artifactory
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and setup to use the cluster
- [Helm](https://helm.sh/) installed and setup to use the cluster (helm init) or [Tillerless Helm](https://github.com/rimusz/helm-tiller)

### TODO: add install examples



## Remove

Removing a **helm** release is done with

```
helm delete --purge helm-hub-sync
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the helm-hub-sync chart and their default values.

|         Parameter            |                    Description                   |           Default                  |
|------------------------------|--------------------------------------------------|------------------------------------|
| `nameOverride`               | Chart name override                              | ``                                 |
| `fullNameOverride`           | Chart full name override                         | ``                                 |
| `existingSecret`             | An existing secret holding secrets               | ``                                 |
| `existingConfigMap`          | An existing configmap holding env vars           | ``                                 |
| `replicaCount`               | Replica count                                    | `1`                                |
| `image.repository`           | Image repository                                 | `docker.bintray.io/helm-hub-sync`  |
| `image.PullPolicy`           | Container pull policy                            | `IfNotPresent`                     |
| `imagePullSecrets`           | List of imagePullSecrets                         | ``                                 |
| `securityContext.enabled`    | Enables Security Context                         | `true`                             |
| `securityContext.userId`     | Security UserId                                  | `1000`                             |
| `securityContext.groupId`    | Security GroupId                                 | `1000`                             |
| `env.timeInterval`           | The time in seconds between two successive runs  | `14400`                            |
| `env.logLevel`               | Logs level                                       | `info`                             |
| `env.consoleLog`             | To create human-friendly, colorized output       | `true`                             |
| `env.artifactory.host`       | The hostname of JFrog Artifactory to connect to  | ``                                 |
| `env.artifactory.helmRepo`   | The Helm Virtual Repository to use               | `helmhub`                          |
| `env.artifactory.authType`   | The authentication type to use                   | `basic`                            |
| `env.artifactory.authData`   | The authentication data to use                   | ``                                 |
| `env.artifactory.keepList`   | List containing Helm Remote repos that will never be removed | ``                     |
| `env.artifactory.keepDeletedRepos`| Whether to keep repos that have been removed from the Helm Hub | `true`          |
| `env.githubIgnoreList`       | A comma separated list containing Helm repos that should never be created | `stable`  | 
| `resources`                  | Specify resources                                | `{}`                               |
| `nodeSelector`               | kubexray micro-service node selector             | `{}`                               |
| `tolerations`                | kubexray micro-service node tolerations          | `[]`                               |
| `affinity`                   | kubexray micro-service node affinity             | `{}`                               |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install/upgrade`.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example

```
helm upgrade --install helm-hub-sync --namespace helm-hub-sync jfrog/helm-hub-sync \
    --set existingSecret="helm-hub-sync" -f override-values.yaml 
```
