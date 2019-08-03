pipeline {
  agent any
  stages {
    stage('GoLangChecks') {
      steps {
        sh 'export GOROOT=/usr/local/go'
        sh 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'
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