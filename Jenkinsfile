pipeline {
    agent any

    options {
        // We explicitly practice checkout scm below.
        skipDefaultCheckout()
    }

    environment {
        APP_NAME = 'task-api'
        BUILD_DIR = 'bin'
        RUN_EXTRA_CHECKS = 'true'
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'

                /*
                 * In a Multibranch Pipeline, checkout scm
                 * checks out the revision that belongs to
                 * the current branch / pull request.
                 */
                checkout scm
            }
        }

        stage('Branch Context') {
            steps {
                script {
                    echo "Branch: ${env.BRANCH_NAME ?: 'not-multibranch'}"
                    echo "Change Request: ${env.CHANGE_ID ?: 'none'}"
                }
            }
        }

        stage('Environment') {
            steps {
                echo 'Checking build environment...'

                sh 'go version'
                sh 'git --version'

                echo "Application: ${APP_NAME}"
                echo "Build directory: ${BUILD_DIR}"
            }
        }

        stage('Format Check') {
            steps {
                echo 'Checking Go formatting...'

                sh '''
                    unformatted="$(gofmt -l 01-go)"

                    if [ -n "$unformatted" ]; then
                        echo "The following files are not formatted:"
                        echo "$unformatted"
                        exit 1
                    fi
                '''
            }
        }

        stage('Test') {
            steps {
                echo 'Running Go tests...'

                sh 'go test ./...'
            }
        }

        stage('Extra Checks') {
            when {
                environment(
                    name: 'RUN_EXTRA_CHECKS',
                    value: 'true'
                )
            }

            steps {
                echo 'Running go vet...'

                sh 'go vet ./...'
            }
        }

        /*
         * Example of branch-based conditional execution.
         *
         * In a Multibranch Pipeline, BRANCH_NAME is
         * provided by Jenkins.
         */
        stage('Main Branch Check') {
            when {
                branch 'main'
            }

            steps {
                echo 'Running checks specific to the main branch...'

                sh 'go test ./...'
            }
        }

        stage('Build') {
            environment {
                CGO_ENABLED = '0'
            }

            steps {
                echo 'Building Go application...'

                sh '''
                    mkdir -p "$BUILD_DIR"

                    go build \
                        -o "$BUILD_DIR/$APP_NAME" \
                        ./01-go/06-database/02-mysql-repository
                '''

                sh 'ls -lh "$BUILD_DIR/$APP_NAME"'
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
        }

        success {
            echo 'CI pipeline succeeded.'
        }

        failure {
            echo 'CI pipeline failed.'
        }
    }
}