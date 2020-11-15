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

        stage('Test') {
            steps {
                sh 'cd classroomCreatorController'
                sh 'go test'
                sh 'cd ..'
            }
        }
        
    }
    post {
        always {
            echo 'Finish Pipeline'
        }
    }  
}