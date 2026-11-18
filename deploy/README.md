# Deployment

This directory contains deployment configuration for the practice project.

```text
deploy/
└── k8s/
    ├── 01-configmap.yaml
    ├── 02-deployment.yaml
    └── 03-service.yaml
```

The standalone Pod manifest remains under the Kubernetes learning directory because production-style workloads should normally be managed through a higher-level controller such as a Deployment.

# Kubernetes Desired State

Kubernetes uses a declarative model.

The manifests describe the desired state of the system.

For example:

```yaml
spec:
  replicas: 2
```

means:

```text
Desired State
→ two Pods should exist
```

It does not describe how many Pods currently exist.

The current state of the cluster is the actual state.

```text
Manifest
↓
Desired State

Cluster
↓
Actual State
```

Kubernetes continuously compares them:

```text
Desired State
        ↓
   Kubernetes
   Controllers
        ↑
Actual State
```

If they differ, Kubernetes attempts to reconcile the actual state with the desired state.

Example:

```text
Desired = 2 Pods
Actual  = 1 Pod

→ create another Pod
```

Another example:

```text
Desired = 2 Pods
Actual  = 3 Pods

→ remove one Pod
```

The central idea is:

```text
declare desired state
↓
Kubernetes observes actual state
↓
controllers reconcile differences
```

# Deployment Resources

## ConfigMap

```text
nginx-demo-config
```

stores the HTML configuration used by nginx.

## Deployment

```text
nginx-demo
```

defines:

```text
replicas = 2
image = nginx:alpine
Pod label = app=nginx-demo
```

The Deployment manages the Pods.

```text
Deployment
↓
ReplicaSet
↓
Pods
```

## Service

```text
nginx-demo-service
```

selects Pods using:

```text
app=nginx-demo
```

The relationship is:

```text
Deployment
    │
    ▼
Pods
app=nginx-demo
    ▲
    │
Service
```

# CI / CD Boundary

The current Jenkins pipeline performs Continuous Integration.

```text
Git Push
↓
CI
↓
Format Check
↓
Go Tests
↓
Deployment Manifest Check
```

The current pipeline stops here.

```text
==============================
       CI / CD BOUNDARY
==============================
```

Deployment to a Kubernetes cluster belongs to CD.

A future deployment stage could perform:

```bash
kubectl apply -f deploy/k8s/
```

followed by:

```bash
kubectl rollout status \
  deployment/nginx-demo
```

Conceptually:

```text
Code
↓
CI
├── format
├── test
└── validate/build artifacts
        │
        │ CI / CD boundary
        ▼
CD
├── deploy
├── rollout
└── verify running application
```

The current repository prepares deployment manifests but does not automatically deploy them to a real Kubernetes cluster.

# Apply Manifests

When a Kubernetes cluster is available:

```bash
kubectl apply \
  -f deploy/k8s/
```

Inspect actual state:

```bash
kubectl get deployments
kubectl get pods
kubectl get services
kubectl get configmaps
```

More specifically:

```bash
kubectl get pods \
  -l app=nginx-demo
```

The Deployment declares:

```text
desired replicas = 2
```

so the actual cluster should eventually report two running Pods.

# Inspect Desired and Actual State

Desired state can be inspected with:

```bash
kubectl get deployment nginx-demo \
  -o yaml
```

Actual Pods can be inspected with:

```bash
kubectl get pods \
  -l app=nginx-demo
```

Deployment status:

```bash
kubectl get deployment nginx-demo
```

A healthy result should eventually show:

```text
READY   UP-TO-DATE   AVAILABLE
2/2     2            2
```

# Test the Service

Forward the Service to the local machine:

```bash
kubectl port-forward \
  service/nginx-demo-service \
  8080:80
```

Then:

```bash
curl http://localhost:8080
```

The expected response contains:

```html
<h1>Hello from Kubernetes</h1>
```

# Cleanup

Delete all deployment resources:

```bash
kubectl delete \
  -f deploy/k8s/
```

# Key Concepts

```text
Manifest
→ declares desired state

Desired State
→ what should exist

Actual State
→ what currently exists

Controller
→ compares desired and actual state

Reconciliation
→ changes actual state toward desired state

CI
→ integrate, check, test

CD
→ deliver/deploy validated changes

kubectl apply
→ sends desired state to Kubernetes
```

The key mental model is:

```text
Git
↓
CI
↓
tested change
================
CI / CD boundary
================
↓
CD
↓
Kubernetes manifest
↓
Desired State
↓
Controller
↓
Actual State
```