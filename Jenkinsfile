node {    
    def app     
    stage('Clone repository') {               
        checkout scm    
    }
    stage('Unit test'){
        sh 'go version'
    }
    stage('Integration test'){
        sh 'go version'
    }     
    stage('Build image') {         
        app = docker.build("jinjustin/omega_api")    
    }          
    stage('Push image') {          
    app.push("${env.BUILD_NUMBER}")            
    app.push("latest")          
    }
}
