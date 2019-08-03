pipeline {
  agent any
  stages {
    stage('Go Vet/Staticcheck') {
      agent any
      environment {
        GOROOT = '/usr/local/go'
        PATH = '$GOROOT/bin:$PATH'
      }
      steps {
        sh 'make checks'
      }
    }
    stage('Unit Tests') {
      steps {
        sh 'export GOROOT=/usr/local/go'
        sh 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'
        sh 'make test'
      }
    }
  }
}