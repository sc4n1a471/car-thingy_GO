def dockerImage
def version
def buildNumber
def branchName

pipeline {
    agent any
    environment {
        DOCKER_REGISTRY = 'sc4n1a471'

        CAR_THINGY_GO_DB_USERNAME = credentials('CAR-THINGY_GO_DB_USERNAME')
        CAR_THINGY_GO_DB_PASSWORD = credentials('CAR-THINGY_GO_DB_PASSWORD')
        CAR_THINGY_GO_DB_IP = credentials('CAR-THINGY_GO_DB_IP')
        CAR_THINGY_GO_DB_PORT = credentials('CAR-THINGY_GO_DB_PORT')
        
        CAR_THINGY_GO_DB_NAME_DEV = credentials('CAR-THINGY_GO_DB_NAME_DEV')
        CAR_THINGY_GO_DB_NAME_PROD = credentials('CAR-THINGY_GO_DB_NAME_PROD')

        CAR_THINGY_GO_API_SECRET = credentials('CAR-THINGY_GO_API_SECRET')

        CAR_THINGY_PYTHON_GRAYLOG_HOST_DEV = credentials('CAR_THINGY_PYTHON_GRAYLOG_HOST_DEV')
        CAR_THINGY_PYTHON_GRAYLOG_HOST_PROD = credentials('CAR_THINGY_PYTHON_GRAYLOG_HOST_PROD')
    }
    stages {
        stage('Checkout') {
            parallel {
                stage('Checkout the branch') {
                    when {
                        not {
                            changeRequest()
                        }
                    }
                    steps {
                        echo "Checking out ${env.BRANCH_NAME} branch..."
                        git branch: env.BRANCH_NAME, credentialsId: 'Home-VM_jenkins', url: 'git@github.com:sc4n1a471/car-thingy_GO.git'
                    }
                }
            }
        }

        // MARK: Read Version
        stage('Read Version') {
            when {
                not {
                    changeRequest()
                }
            }
            steps {
                script {
                    version = readFile('version').trim()
                    echo "Building version ${version}"

                    buildNumber = env.BUILD_NUMBER
                    echo "Build number: ${buildNumber}"

                    branchName = env.BRANCH_NAME.split('-')[0]
                    echo "Build branch: ${branchName}"
                }
            }
        }

        // MARK: Build and Push Docker Image
        stage('Build and Push') {
            parallel {
                stage('Push production docker image') {
                    when {
                        branch 'main'
                    }
                    steps {
                        script {
                            dockerImage = docker.build("sc4n1a471/car-thingy_go:${version}-${buildNumber}")
                            docker.withRegistry('https://registry.hub.docker.com', 'DOCKER_HUB') {
                                dockerImage.push("latest")
                                dockerImage.push("${version}-${buildNumber}")
                            }
                        }
                    }
                }

                stage('Push not production docker image') {
                    when {
                        not {
                            changeRequest()
                        }
                        not {
                            branch 'main'
                        }
                    }
                    steps {
                        script {
                            dockerImage = docker.build("sc4n1a471/car-thingy_go:${version}-${branchName}-${buildNumber}")
                            docker.withRegistry('https://registry.hub.docker.com', 'DOCKER_HUB') {
                                dockerImage.push("latest-${branchName}")
                                dockerImage.push("${version}-${branchName}-${buildNumber}")
                            }
                        }
                    }
                }
            }
        }

        // MARK: Deploy to Development
        stage('Deploy development') {
            when {
                not {
                    changeRequest()
                }
                not {
                    branch 'main'
                }
            }

            steps {
                script {
                    echo "Deploying version ${version}, build ${buildNumber} to ${branchName} branch"

                    sh """
                    if [ \$(docker ps -a -q -f name=car-thingy_go_$branchName) ]; then
                        docker rm -f car-thingy_go_$branchName
                        echo "Container removed"
                    fi
                        
                    if [ \$(docker images -q sc4n1a471/car-thingy_go:$version-$branchName-$buildNumber) ]; then
                        docker rmi -f sc4n1a471/car-thingy_go:$version-$branchName-$buildNumber
                        echo "Image removed"
                    fi
                    """

                    sh """
                    terraform init

                    terraform apply \
                        -var="container_name=car-thingy_go_$branchName" \
                        -var="container_version=$version-$branchName-$buildNumber" \
                        -var="env=$branchName" \
                        -var="db_username=\$CAR_THINGY_GO_DB_USERNAME" \
                        -var="db_password=\$CAR_THINGY_GO_DB_PASSWORD" \
                        -var="db_ip=\$CAR_THINGY_GO_DB_IP" \
                        -var="db_port=\$CAR_THINGY_GO_DB_PORT" \
                        -var="db_name=\$CAR_THINGY_GO_DB_NAME_DEV" \
                        -var="api_secret=\$CAR_THINGY_GO_API_SECRET" \
                        -var="graylog_host=\$CAR_THINGY_PYTHON_GRAYLOG_HOST_DEV" \
                        -auto-approve
                    """
                }
            }
        }
        // MARK: Deploy to Production
        stage('Deploy production') {
            when {
                branch 'main'
            }

            steps {
                script {
                    echo "Deploying version ${version}, build ${buildNumber} to PROD"

                    sh """
                    docker rm -f car-thingy_go
                    echo "Container removed"

                    if [ \$(docker images -q sc4n1a471/car-thingy_go:$version-$buildNumber) ]; then
                        docker rmi -f sc4n1a471/car-thingy_go:$version-$buildNumber
                        echo "Image removed"
                    fi
                    """

                    sh """
                    terraform init

                    terraform apply \
                        -var="container_name=car-thingy_go" \
                        -var="container_version=$version-$buildNumber" \
                        -var="env=prod" \
                        -var="db_username=\$CAR_THINGY_GO_DB_USERNAME" \
                        -var="db_password=\$CAR_THINGY_GO_DB_PASSWORD" \
                        -var="db_ip=\$CAR_THINGY_GO_DB_IP" \
                        -var="db_port=\$CAR_THINGY_GO_DB_PORT" \
                        -var="db_name=\$CAR_THINGY_GO_DB_NAME_PROD" \
                        -var="api_secret=\$CAR_THINGY_GO_API_SECRET" \
                        -var="graylog_host=\$CAR_THINGY_PYTHON_GRAYLOG_HOST_PROD" \
                        -auto-approve
                    """
                }
            }
        }
    }
    post {
        success {
            echo 'Build and deployment successful!'
        }
        failure {
            echo 'Build or deployment failed.'
        }
    }
}