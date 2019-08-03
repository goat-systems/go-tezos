pipeline {
  agent any
  stages {
    stage('GoLangChecks') {
      steps {
        sh 'make checks'
      }
    }
    stage('Unit Tests') {
      steps {
        sh 'make test'
      }
    }
  }
}