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

        stage('Test Creator Controller') {
            steps {
                    echo 'Running test'
                    sh 'cd classroomCreatorController && go test -v'
            }
        }
        stage('Test Deleter Controller') {
            steps {
                    echo 'Running test'
                    sh 'cd classroomDeleterController && go test -v'
            }
        }
        stage('Test List Controller') {
            steps {
                    echo 'Running test'
                    sh 'cd classroomListController && go test -v'
            }
        }
    }
        post {
            always {
                echo 'Finish Pipeline'
            }
        }  
    }
