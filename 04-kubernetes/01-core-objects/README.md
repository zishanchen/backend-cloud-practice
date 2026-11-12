# Kubernetes Core Objects Practice

This exercise introduces four basic Kubernetes resources:

- Pod
- Deployment
- Service
- ConfigMap

## Architecture

```text
                  ConfigMap
             nginx-demo-config
                      │
                      │ mounted as volume
                      ▼
                 Deployment
                  nginx-demo
                      │
                      │ manages
                      ▼
                  ReplicaSet
                   /      \
                  /        \
                 ▼          ▼
              Pod A       Pod B
                 ▲          ▲
                  \        /
                   \      /
                    Service
             nginx-demo-service
                  selector:
                app=nginx-demo
```

A standalone Pod is also included for practicing the basic Pod definition:

```text
ConfigMap
   │
   ▼
nginx-demo-pod
```

## Pod

A Pod is the basic workload unit used to run containers.

```yaml
apiVersion: v1
kind: Pod
```

The example Pod runs:

```text
nginx:alpine
```

and mounts the ConfigMap as:

```text
/usr/share/nginx/html
```

## ConfigMap

The ConfigMap stores application configuration separately from the container image.

In this example it stores:

```text
index.html
```

The file is mounted into the Nginx container.

Flow:

```text
ConfigMap
↓
Volume
↓
Container
```

## Deployment

A Deployment manages Pods.

The example declares:

```yaml
replicas: 2
```

so Kubernetes attempts to maintain two matching Pods.

Conceptually:

```text
Deployment
↓
ReplicaSet
↓
Pods
```

The Deployment selector:

```yaml
selector:
  matchLabels:
    app: nginx-demo
```

must match the Pod template label:

```yaml
labels:
  app: nginx-demo
```

## Service

A Service provides a stable network endpoint for matching Pods.

The example uses:

```yaml
type: ClusterIP
```

and selects Pods with:

```yaml
selector:
  app: nginx-demo
```

Flow:

```text
Client
↓
Service
↓
Pod
```

The Service does not select the Deployment directly.

Instead:

```text
Deployment
→ manages Pods

Service
→ selects Pods
```

## Key Relationships

```text
ConfigMap
→ configuration

Pod
→ runs containers

Deployment
→ manages Pods

Service
→ provides stable access to Pods
```

## Apply

Create the ConfigMap first:

```bash
kubectl apply \
  -f 04-kubernetes/01-core-objects/01-configmap.yaml
```

Then create the standalone Pod:

```bash
kubectl apply \
  -f 04-kubernetes/01-core-objects/02-pod.yaml
```

Create the Deployment:

```bash
kubectl apply \
  -f 04-kubernetes/01-core-objects/03-deployment.yaml
```

Create the Service:

```bash
kubectl apply \
  -f 04-kubernetes/01-core-objects/04-service.yaml
```

## Inspect Resources

```bash
kubectl get configmaps
kubectl get pods
kubectl get deployments
kubectl get services
```

More detailed inspection:

```bash
kubectl describe pod nginx-demo-pod

kubectl describe deployment nginx-demo

kubectl describe service nginx-demo-service
```

## Check Deployment Pods

```bash
kubectl get pods \
  -l app=nginx-demo
```

Expected:

```text
2 Pods
```

because:

```yaml
replicas: 2
```

## Port Forward

Forward the Service locally:

```bash
kubectl port-forward \
  service/nginx-demo-service \
  8080:80
```

Then open:

```text
http://localhost:8080
```

The page content comes from the ConfigMap.

## Cleanup

```bash
kubectl delete \
  -f 04-kubernetes/01-core-objects/04-service.yaml

kubectl delete \
  -f 04-kubernetes/01-core-objects/03-deployment.yaml

kubectl delete \
  -f 04-kubernetes/01-core-objects/02-pod.yaml

kubectl delete \
  -f 04-kubernetes/01-core-objects/01-configmap.yaml
```

## Key Takeaway

```text
ConfigMap
→ provides configuration

Deployment
→ declares and manages Pods

Pod
→ runs containers

Service
→ finds Pods using labels and exposes them
```