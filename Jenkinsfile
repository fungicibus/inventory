pipeline {
    agent {
        label 'RemoteAgentTest'
    }

    environment { 
        IMAGE_NAME = 'fungicibus/inventory' 
        CONTAINER_NAME = 'inventory' 
        REPO_URL = 'https://github.com/fungicibus/inventory.git' 
        DOCKER_HUB_REPO = 'shifter1703/fungicibus' 
        DOCKER_CREDENTIALS_ID = 'ce217c82-26b0-4acb-b57f-71a11965e25d' 
        ENV_FILE = 'dev.env' 
    }

    stages { 

        stage('Checkout') {
            steps {
                script {
                    git branch: 'v0', url: env.REPO_URL 
                    env.IMAGE_TAG = sh(script: "git describe --tags --abbrev=0", returnStdout: true).trim()
                    echo "Используемый IMAGE_TAG: ${env.IMAGE_TAG}" 
                }
            }
        }

        stage('Parse .env File and Fetch Secrets') {
            steps {
                script {
                    sh "grep -v '^#' ${env.ENV_FILE} | awk 'NF' > parsed_env"
                    echo "Parsed environment variables from ${env.ENV_FILE}"
                    
                    def envVars = readFile('parsed_env').readLines()
                    def updatedVars = []
                    
                    envVars.each { line ->
                        if (line.contains("=\$")) {
                            def key = line.split("=")[0]
                            def secretName = "my-secret"
                            echo "Fetching secret for ${key} from path: kv/data/${secretName}"
                            def secretValue = vault path: "kv/data/${secretName}", key: "value"
                            if (secretValue != null) {
                                updatedVars.add("${key}=${secretValue}")
                            } else {
                                error "Secret ${secretName} does not exist in the vault"
                            }
                        } else {
                            updatedVars.add(line)
                        }
                    }
                    
                    writeFile file: 'parsed_env', text: updatedVars.join("\n")
                }
            }
        }

        stage('Stop and remove container') {
            steps {
                script {
                    sh "docker rm -f ${env.CONTAINER_NAME} || true"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    sh """
                        docker build -f ./docker/Dockerfile -t ${env.IMAGE_NAME}:${env.IMAGE_TAG} . 
                        docker tag ${env.IMAGE_NAME}:${env.IMAGE_TAG} ${env.IMAGE_NAME}:latest
                    """
                }
            }
        }

        stage('Start container') {
            steps {
                script {
                    sh "docker run -d --name ${env.CONTAINER_NAME} --env-file parsed_env ${env.IMAGE_NAME}:${env.IMAGE_TAG}" 
                }
            }
        }

        stage('Push in Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: env.DOCKER_CREDENTIALS_ID, 
                                                 usernameVariable: 'DOCKER_HUB_USERNAME', 
                                                 passwordVariable: 'DOCKER_HUB_PASSWORD')]) { 
                    script {
                        sh """
                            echo \$DOCKER_HUB_PASSWORD | docker login -u \$DOCKER_HUB_USERNAME --password-stdin
                            docker tag ${env.IMAGE_NAME}:${env.IMAGE_TAG} ${env.DOCKER_HUB_REPO}:${env.IMAGE_TAG} 
                            docker tag ${env.IMAGE_NAME}:${env.IMAGE_TAG} ${env.DOCKER_HUB_REPO}:latest 
                            docker push ${env.DOCKER_HUB_REPO}:${env.IMAGE_TAG} 
                            docker push ${env.DOCKER_HUB_REPO}:latest 
                        """
                    }
                }
            }
        }
    }

    post { 
        success {
            echo "The pipeline has been completed successfully!" 
        }
        failure {
            echo "The pipeline ended with an error. Check the logs for details."
        }
        always {
            echo "Pipeline completion completed" 
            cleanWs() 
        }
    }
}