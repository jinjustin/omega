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

        stage('Unit Test: classroomCreatorController') {
            steps {
                    echo 'classroomCreatorController'
                    sh 'cd classroomCreatorController && go test -v'
            }
        }
        stage('Unit Test: classroomDeleterController') {
            steps {
                    echo 'classroomDeleterController'
                    sh 'cd classroomDeleterController && go test -v'
            }
        }
        stage('Unit Test: classroomListController') {
            steps {
                    echo 'classroomListController'
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
