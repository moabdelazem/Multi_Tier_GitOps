pipeline {
    agent any

    // ================================
    // Environment Variables
    // ================================
    environment {
        // Go Docker image for running tests/lint
        GO_IMAGE = 'golang:1.25-alpine'
        
        // Docker configuration
        DOCKER_REGISTRY = credentials('docker-registry-url')
        DOCKER_CREDENTIALS = credentials('docker-registry-credentials')
        IMAGE_NAME = 'multi-tier-api'
        IMAGE_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT?.take(7) ?: 'latest'}"
        
        // Application configuration
        APP_NAME = 'multi-tier-api'
        API_DIR = 'api'
    }

    // ================================
    // Build Parameters
    // ================================
    parameters {
        booleanParam(name: 'SKIP_TESTS', defaultValue: false, description: 'Skip running tests')
        booleanParam(name: 'SKIP_LINT', defaultValue: false, description: 'Skip linting')
        booleanParam(name: 'PUSH_IMAGE', defaultValue: false, description: 'Push Docker image to registry')
        string(name: 'DOCKER_TAG', defaultValue: '', description: 'Custom Docker tag (optional)')
    }

    // ================================
    // Pipeline Options
    // ================================
    options {
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
        buildDiscarder(logRotator(numToKeepStr: '10'))
        skipDefaultCheckout(false)
    }

    // ================================
    // Pipeline Stages
    // ================================
    stages {
        // ---------------------------
        // Stage: Checkout
        // ---------------------------
        stage('Checkout') {
            steps {
                echo '📥 Checking out source code...'
                checkout scm
                
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    
                    env.GIT_BRANCH_NAME = sh(
                        script: 'git rev-parse --abbrev-ref HEAD',
                        returnStdout: true
                    ).trim()
                    
                    echo "Branch: ${env.GIT_BRANCH_NAME}"
                    echo "Commit: ${env.GIT_COMMIT_SHORT}"
                }
            }
        }

        // ---------------------------
        // Stage: Dependencies
        // ---------------------------
        stage('Dependencies') {
            steps {
                echo 'Downloading Go dependencies...'
                script {
                    docker.image("${GO_IMAGE}").inside("-v ${WORKSPACE}:/workspace -w /workspace/${API_DIR}") {
                        sh '''
                            go version
                            go mod download
                            go mod verify
                        '''
                    }
                }
            }
        }

        // ---------------------------
        // Stage: Code Quality (Parallel)
        // ---------------------------
        stage('Code Quality') {
            parallel {
                // Lint Stage
                stage('Lint') {
                    when {
                        expression { return !params.SKIP_LINT }
                    }
                    steps {
                        echo 'Running linter...'
                        script {
                            docker.image('golangci/golangci-lint:v1.61.0-alpine').inside("-v ${WORKSPACE}:/workspace -w /workspace/${API_DIR}") {
                                sh 'golangci-lint run ./... --timeout=5m'
                            }
                        }
                    }
                }

                // Format Check
                stage('Format Check') {
                    steps {
                        echo 'Checking code formatting...'
                        script {
                            docker.image("${GO_IMAGE}").inside("-v ${WORKSPACE}:/workspace -w /workspace/${API_DIR}") {
                                sh '''
                                    # Check if code is properly formatted
                                    gofmt_output=$(gofmt -l .)
                                    if [ -n "$gofmt_output" ]; then
                                        echo "The following files are not properly formatted:"
                                        echo "$gofmt_output"
                                        exit 1
                                    fi
                                    echo "All files are properly formatted"
                                '''
                            }
                        }
                    }
                }
            }
        }

        // ---------------------------
        // Stage: Test
        // ---------------------------
        stage('Test') {
            when {
                expression { return !params.SKIP_TESTS }
            }
            steps {
                echo 'Running tests...'
                script {
                    docker.image("${GO_IMAGE}").inside("-v ${WORKSPACE}:/workspace -w /workspace/${API_DIR}") {
                        sh '''
                            # Run tests with coverage
                            go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
                            
                            # Generate coverage report
                            go tool cover -func=coverage.out
                            
                            # Generate HTML coverage report
                            go tool cover -html=coverage.out -o coverage.html
                        '''
                    }
                }
            }
            post {
                always {
                    // Archive test coverage reports
                    archiveArtifacts artifacts: "${API_DIR}/coverage.*", allowEmptyArchive: true
                    
                    // Publish coverage results if plugin available
                    script {
                        if (fileExists("${API_DIR}/coverage.out")) {
                            echo "Coverage report generated"
                        }
                    }
                }
            }
        }

        // ---------------------------
        // Stage: Docker Build
        // ---------------------------
        stage('Docker Build') {
            steps {
                echo 'Building Docker image...'
                dir("${API_DIR}") {
                    script {
                        def imageTag = params.DOCKER_TAG ?: env.IMAGE_TAG
                        
                        sh """
                            docker build \
                                --build-arg BUILD_NUMBER=${BUILD_NUMBER} \
                                --build-arg GIT_COMMIT=${GIT_COMMIT_SHORT} \
                                -t ${IMAGE_NAME}:${imageTag} \
                                -t ${IMAGE_NAME}:latest \
                                .
                            
                            # List built images
                            docker images | grep ${IMAGE_NAME}
                        """
                        
                        env.BUILT_IMAGE_TAG = imageTag
                    }
                }
            }
        }

        // ---------------------------
        // Stage: Docker Push
        // ---------------------------
        stage('Docker Push') {
            when {
                expression { return params.PUSH_IMAGE }
            }
            steps {
                echo 'Pushing Docker image to registry...'
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'docker-registry-credentials',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh """
                            echo \${DOCKER_PASS} | docker login ${DOCKER_REGISTRY} -u \${DOCKER_USER} --password-stdin
                            
                            docker tag ${IMAGE_NAME}:${BUILT_IMAGE_TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${BUILT_IMAGE_TAG}
                            docker tag ${IMAGE_NAME}:latest ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest
                            
                            docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${BUILT_IMAGE_TAG}
                            docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest
                            
                            docker logout ${DOCKER_REGISTRY}
                        """
                    }
                }
            }
        }

    }

    // ================================
    // Post-Build Actions
    // ================================
    post {
        always {
            echo 'Cleaning up workspace...'
            
            // Clean up Docker images to save space
            sh '''
                docker image prune -f || true
            '''
            
            // Clean workspace
            cleanWs(notFailBuild: true)
        }
        
        success {
            echo 'Pipeline completed successfully!'
            
            script {
                // Send success notification (configure as needed)
                if (env.BRANCH_NAME == 'main' || env.BRANCH_NAME == 'master') {
                    echo "Main branch build #${BUILD_NUMBER} succeeded!"
                    // Uncomment to enable Slack notification
                    // slackSend(color: 'good', message: "Build #${BUILD_NUMBER} succeeded for ${APP_NAME}")
                }
            }
        }
        
        failure {
            echo 'Pipeline failed!'
            
            script {
                // Send failure notification
                echo "Build #${BUILD_NUMBER} failed!"
                // Uncomment to enable Slack notification
                // slackSend(color: 'danger', message: "Build #${BUILD_NUMBER} failed for ${APP_NAME}")
            }
        }
        
        unstable {
            echo 'Pipeline is unstable!'
        }
    }
}
