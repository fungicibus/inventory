pipeline {
    agent {
        label 'RemoteAgentKolya'
    }

    environment { // Окружение
        IMAGE_NAME = 'fungicibus/inventory' // Имя образа
        CONTAINER_NAME = 'inventory' // Имя контейнера
        REPO_URL = 'https://github.com/fungicibus/inventory.git' // URL репозитория
        DOCKER_HUB_REPO = 'shifter1703/fungicibus' // Репозиторий Docker Hub
        DOCKER_CREDENTIALS_ID = 'ce217c82-26b0-4acb-b57f-71a11965e25d' // ID учетных данных Docker
        ENV_FILE = 'dev.env' // Файл окружения
    }

    stages { // Этапы

        stage('Получение кода') {
            steps {
                script {
                    git branch: 'v0', url: env.REPO_URL // Клонирование ветки v0
                    env.IMAGE_TAG = sh(script: "git describe --tags --abbrev=0", returnStdout: true).trim() // Получение тега
                    echo "Используемый IMAGE_TAG: ${env.IMAGE_TAG}" // Вывод тега
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
                            def key = line.split("=")[0]  // SECRET_SOURCE
                            def secretName = "my-secret"  // Фиксированное имя секрета
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

        stage('Остановка и удаление существующего контейнера') {
            steps {
                script {
                    sh "docker rm -f ${env.CONTAINER_NAME} || true" // Удаление контейнера
                }
            }
        }

        stage('Сборка Docker-образа') {
            steps {
                script {
                    sh """
                        docker build -f ./docker/Dockerfile -t ${env.IMAGE_NAME}:${env.IMAGE_TAG} . 
                        docker tag ${env.IMAGE_NAME}:${env.IMAGE_TAG} ${env.IMAGE_NAME}:latest
                    """
                }
            }
        }

        stage('Запуск нового контейнера') {
            steps {
                script {
                    sh "docker run -d --name ${env.CONTAINER_NAME} --env-file parsed_env ${env.IMAGE_NAME}:${env.IMAGE_TAG}" // Запуск контейнера
                }
            }
        }

        stage('Отправка в Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: env.DOCKER_CREDENTIALS_ID, 
                                                 usernameVariable: 'DOCKER_HUB_USERNAME', 
                                                 passwordVariable: 'DOCKER_HUB_PASSWORD')]) { // Учетные данные
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

    post { // Пост-условия
        success {
            echo "Пайплайн успешно выполнен!" // Успешное выполнение
        }
        failure {
            echo "Пайплайн завершился с ошибкой. Проверьте логи для деталей." // Ошибка
        }
        always {
            echo "Выполнение пайплайна завершено!" // Завершение
            cleanWs() // Очистка рабочего пространства
        }
    }
}