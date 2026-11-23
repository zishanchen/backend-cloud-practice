pipeline {
    agent any

    /*
     * Declarative Pipeline normally performs an automatic
     * SCM checkout.
     *
     * We disable it here because this exercise explicitly
     * practices "checkout scm" in its own stage.
     */
    options {
        skipDefaultCheckout()
    }

    /*
     * Global environment variables.
     * These are available to all stages.
     */
    environment {
        APP_NAME = 'task-api'
        BUILD_DIR = 'bin'
        RUN_EXTRA_CHECKS = 'true'
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'

                checkout scm
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

        /*
         * "when" controls whether this stage executes.
         *
         * RUN_EXTRA_CHECKS is defined in the global
         * environment block above.
         */
        stage('Extra Checks') {
            when {
                environment name: 'RUN_EXTRA_CHECKS', value: 'true'
            }

            steps {
                echo 'RUN_EXTRA_CHECKS=true, running go vet...'

                sh 'go vet ./...'
            }
        }

        stage('Build') {
            /*
             * Stage-level environment variable.
             * CGO_ENABLED only applies to this stage.
             */
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

    /*
     * post runs after the Pipeline finishes.
     */
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