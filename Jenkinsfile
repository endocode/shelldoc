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
                sh 'cd /shelldoc && go test ./...'
            }
        }
        stage('DocTest') {
            steps {
                sh 'cd /shelldoc && ./cmd/shelldoc/shelldoc run README.md --xml testresults.xml'
            }
            post {
                always {
                    sh 'cp /shelldoc/testresults*.xml "${WORKSPACE}"/'
                    junit '**/testresults*.xml'
                }
            }
        }
    }
}
