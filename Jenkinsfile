// DO NOT EDIT
pipeline {
    agent { docker { image 'golang' } }
    environment { GOCACHE = '/tmp/.cache' }
    stages {
        stage('Pre Test') {
            steps{
                echo 'Pulling Dependencies'
                echo 'mkdir /tmp/.cahce'
                echo 'chmod 777 /tmp/.cache'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
                sh 'go get github.com/tebeka/go2xunit'  
            }  
        }
        stage('Test'){
            steps{
                echo 'Vetting'
                sh 'cd $GOPATH/src/github.com/goat-systems/go-tezos && go tool vet ./...'

                echo 'Linting'
                sh 'cd $GOPATH/src/github.com/goat-systems/go-tezos && golint ./...'

                echo 'Testing'
                sh 'cd $GOPATH/src/github.com/goat-systems/go-tezos && go test -race -cover ./...'
                
            }  
        }
    }
}