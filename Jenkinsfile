node {    
    def app     
    stage('Clone repository') {               
        checkout scm    
    }
    stage('Unit test'){
        //
    }
    stage('Integration test'){
        //
    }     
    stage('Build image') {         
        app = docker.build("jinjustin/omega_api")    
    }          
    stage('Push image') {          
    app.push("${env.BUILD_NUMBER}")            
    app.push("latest")          
    }
}
