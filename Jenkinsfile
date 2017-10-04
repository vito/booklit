pipeline {
    agent { docker 'golang' }
    environment {
      PATH = '/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin'
    }
    stages {
        stage('build') {
            steps {
                sh 'env'
            }
        }
    }
}
