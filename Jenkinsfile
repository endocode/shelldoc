#!groovy​

pipeline {

    agent {
        dockerfile {
            filename 'ci/Dockerfile'
            args  '-u root'
            label 'linux'
        }
    }

    triggers {
        pollSCM('H/10 * * * *')
    }

    stages {
        stage('Build') {
            steps {
                sh 'cd /shelldoc && make'
            }
        }
        stage('Test') {
            steps {
                sh 'cd /shelldoc && go test ./... 2>&1 | go-junit-report > testresults-gotest.xml'
            }
        }
        stage('DocTest') {
            steps {
                sh 'cd /shelldoc && ./cmd/shelldoc/shelldoc run README.md --xml testresults-shelldoc.xml'
            }
            post {
                always {
		    sh 'rm -f "${WORKSPACE}"/testresults*.xml'
		    sh 'cp /shelldoc/testresults*.xml "${WORKSPACE}"/'
		    sh 'ls -la "${WORKSPACE}"/'
                    junit '**/testresults*.xml'
                    cleanWs()
                }
            }
        }
    }
}
