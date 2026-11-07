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
    }
}