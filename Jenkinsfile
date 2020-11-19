def notifyLINE(status) {
    def token = "NRZ7KismUFdqDx7ctyBMdPruWjEEVsIlVASPxg7yT4f"
    def jobName = env.JOB_NAME
    def buildNo = env.BUILD_NUMBER
      
    def url = 'https://notify-api.line.me/api/notify'
    def message = "${jobName} Build #${buildNo} ${status} \r\n"
    sh "curl ${url} -H 'Authorization: Bearer ${token}' -F 'message=${message}'"
}

pipeline {
    agent any
    tools {
        go 'go1.15'
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Dependencies'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
            }
        }
        
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }

        stage('Unit Test: classroom') {
            steps {
                    echo 'classroom'
                    sh 'cd classroom && go test -v'
            }
        }

        stage('Unit Test') {
            steps {
                    echo 'classroomCreatorController'
                    sh 'cd classroomCreatorController && go test -v'
            }
            steps {
                    echo 'classroomDeleterController'
                    sh 'cd classroomDeleterController && go test -v'
            }
            steps {
                    echo 'classroomListController'
                    sh 'cd classroomListController && go test -v'
            }
        }
        stage('Integration Test') {
            steps {
                    echo 'Integration with API'
            }
        }
    }
        post {
            success{
                notifyLINE("succeed")
            }
    failure{
                notifyLINE("failed")
            }
        }  
    }
