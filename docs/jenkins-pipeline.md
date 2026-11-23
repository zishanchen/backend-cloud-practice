# Jenkins CI Pipeline

This repository uses a Declarative Jenkins Pipeline defined in the root `Jenkinsfile`.

The current CI flow is:

```text
Checkout
↓
Environment Check
↓
Format Check
↓
Go Test
↓
Conditional Extra Checks
↓
Go Build
↓
Post Actions
```

## Pipeline as Code

The Jenkins Pipeline is stored in:

```text
Jenkinsfile
```

Keeping the pipeline in the repository allows the CI configuration to be versioned together with the source code.

## Checkout

The Pipeline explicitly checks out the source code with:

```groovy
checkout scm
```

Because the checkout is performed explicitly, automatic Declarative Pipeline checkout is disabled with:

```groovy
options {
    skipDefaultCheckout()
}
```

Conceptually:

```text
Git Repository
↓
checkout scm
↓
Jenkins Workspace
```

All later commands operate on this workspace.

## Environment Variables

Global environment variables are defined using:

```groovy
environment {
    APP_NAME = 'task-api'
    BUILD_DIR = 'bin'
    RUN_EXTRA_CHECKS = 'true'
}
```

These variables are available throughout the Pipeline.

A stage can also define its own environment:

```groovy
stage('Build') {
    environment {
        CGO_ENABLED = '0'
    }
}
```

A stage-level variable only applies inside that stage.

## `environment` vs `when`

`environment` defines values:

```text
environment
→ define variables
```

`when` controls whether a stage executes:

```text
when
→ evaluate a condition
→ run or skip stage
```

Example:

```groovy
environment {
    RUN_EXTRA_CHECKS = 'true'
}
```

defines a value.

Then:

```groovy
when {
    environment name: 'RUN_EXTRA_CHECKS', value: 'true'
}
```

uses that value as a condition.

Conceptually:

```text
environment
↓
RUN_EXTRA_CHECKS=true
↓
when
↓
true?
├── yes → execute stage
└── no  → skip stage
```

## Format Check

The Pipeline checks Go formatting using:

```bash
gofmt -l 01-go
```

The CI Pipeline does not automatically modify incorrectly formatted source code.

Instead:

```text
formatted correctly
→ continue

incorrectly formatted
→ fail CI
```

Developers should format code locally with:

```bash
go fmt ./...
```

before pushing.

## Tests

All Go tests are run with:

```bash
go test ./...
```

This includes the unit tests in the repository.

Integration tests that require MySQL may skip when the required `MYSQL_DSN` environment variable is not configured.

## Conditional Checks

The `Extra Checks` stage uses:

```groovy
when {
    environment name: 'RUN_EXTRA_CHECKS', value: 'true'
}
```

When enabled, it runs:

```bash
go vet ./...
```

This demonstrates conditional stage execution.

## Build

The Pipeline builds the MySQL repository API:

```bash
go build \
    -o bin/task-api \
    ./01-go/06-database/02-mysql-repository
```

The result is:

```text
bin/task-api
```

The build flow is:

```text
Go source
↓
go build
↓
binary
```

## Post Actions

The Pipeline defines:

```groovy
post {
    always {
        ...
    }

    success {
        ...
    }

    failure {
        ...
    }
}
```

Their meanings are:

```text
always
→ runs after every Pipeline result

success
→ runs only after successful Pipeline execution

failure
→ runs only after failed Pipeline execution
```

## Current CI Boundary

The current Jenkinsfile performs Continuous Integration:

```text
checkout
↓
format check
↓
test
↓
static check
↓
build
```

It does not currently deploy the application.

The boundary is therefore:

```text
Git Push
↓
Jenkins CI
├── checkout
├── format
├── test
├── vet
└── build
        │
====================
    CI / CD boundary
====================
        │
        ▼
future deployment
```

A future CD stage could deploy Kubernetes manifests from:

```text
deploy/k8s/
```

but deployment is intentionally not part of the current Jenkins Pipeline.

## Key Concepts

```text
pipeline
→ entire Jenkins workflow

agent
→ machine/executor where Pipeline runs

stage
→ logical phase of the Pipeline

steps
→ actual commands executed inside a stage

environment
→ defines environment variables

when
→ decides whether a stage executes

checkout scm
→ retrieves the source revision

go test
→ verifies Go code

go build
→ produces the application binary

post
→ actions performed after Pipeline/stage completion
```

## Main Mental Model

```text
Source Code
↓
Checkout
↓
Check
↓
Test
↓
Build
↓
Post
```

and:

```text
environment
→ provides values

when
→ makes execution decisions using conditions
```