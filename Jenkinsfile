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

        stage('Unit Test') {
            steps {
                    echo 'course'
                    sh 'cd course && go test -v'
                    echo 'courseController'
                    sh 'cd coursecontroller && go test -v'
                    echo 'courseMemberController'
                    sh 'cd coursemembercontroller && go test -v'
            }
        }
        stage('Integration Test') {
            steps {
                    echo 'Integration with API'
            }
        }
    }

    /*stages {    
        def app     
        stage('Clone repository'){               
            checkout scm    
        }     
        stage('Build image'){         
            app = docker.build("omega:${env.BUILD_NUMBER}")    
        }     
        
        stage('Push image') {            
            app.push("${env.BUILD_NUMBER}")            
            app.push("latest")        
        }
    }*/

    post{
            success{
                notifyLINE("succeed")
            }
    failure{
                notifyLINE("failed")
            }
        }  
}

node {    
    def app     
    stage('Clone repository'){               
        checkout scm    
    }     
    stage('Build image'){         
        app = docker.build("omega:${env.BUILD_NUMBER}")    
    }     
    
    stage('Push image') {            
        app.push("${env.BUILD_NUMBER}")            
        app.push("latest")        
    }
}
