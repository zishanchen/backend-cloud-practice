pipeline {
    agent any

    stages {
        stage('Environment') {
            steps {
                sh 'go version'
            }
        }

        stage('Format Check') {
            steps {
                sh 'test -z "$(gofmt -l 01-go)"'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Deployment Manifest Check') {
            steps {
                sh '''
                    test -f deploy/k8s/01-configmap.yaml
                    test -f deploy/k8s/02-deployment.yaml
                    test -f deploy/k8s/03-service.yaml
                '''
            }
        }

    }
}