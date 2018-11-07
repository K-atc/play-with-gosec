pipeline {
    agent { docker { image 'golang' } }
    // NOTE: docker image 'gosec' is built by `build.sh`
    // agent { docker { image 'gosec' } }

    stages {
        stage('setup') {
            steps {
                withEnv(["GOPATH=${env.WORKSPACE}"]) {
                    git 'https://github.com/K-atc/play-with-gosec.git'
                    sh 'go get github.com/securego/gosec/cmd/gosec/...'
                    sh 'go get github.com/google/go-github/github'
                    sh 'go get github.com/securego/gosec'
                    sh 'go get golang.org/x/oauth2'
                }
            }
        }

        stage('build') {
            steps {
                withEnv(["GOPATH=${env.WORKSPACE}"]) {
                    withCredentials([string(credentialsId: 'GITHUB_ACCESS_TOKEN', variable: 'GITHUB_ACCESS_TOKEN')]) {
                        sh 'go build comment_on_github'
                        // NOTE: exit code is 1 when there're issues.
                        sh './bin/gosec -fmt=json -out gosec.json src/... || echo ok'
                        sh 'GITHUB_ACCESS_TOKEN=$GITHUB_ACCESS_TOKEN ./comment_on_github gosec.json'
                    }
                }
            }
        }
    }
}