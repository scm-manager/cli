#!groovy
pipeline {

  options {
    buildDiscarder(logRotator(numToKeepStr: '10'))
    disableConcurrentBuilds()
  }

  agent {
    node {
      label 'docker'
    }
  }

  environment {
    HOME = "${env.WORKSPACE}"
   	LANGUAGE = "en"
  }

  stages {

    stage('Compute Version') {
      steps {
        script {
          version = computeVersion()
        }
      }
    }

    stage('Tests') {
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        sh 'go test ./...'
      }
    }

    stage('Build') {
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        sh 'go build -a -tags netgo -ldflags \'-w -extldflags "-static"\' -o scm scm.go'
      }
    }

	stage('Publish') {
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        withPublishEnvironment {
       	  sh 'curl -sL https://git.io/goreleaser | bash -s -- release --rm-dist --skip-publish --skip-validate'
        }
      }
	}

  }

  post {
    failure {
      mail to: "scm-team@cloudogu.com",
        subject: "${JOB_NAME} - Build #${BUILD_NUMBER} - ${currentBuild.currentResult}!",
        body: "Check console output at ${BUILD_URL} to view the results."
    }
  }

}

String version

String computeVersion() {
  def commitHashShort = sh(returnStdout: true, script: 'git rev-parse --short HEAD')
  return "${new Date().format('yyyyMMddHHmm')}-${commitHashShort}".trim()
}


void withPublishEnvironment(Closure<Void> closure) {
  withCredentials([
    usernamePassword(credentialsId: 'maven.scm-manager.org', usernameVariable: 'ORG_GRADLE_PROJECT_packagesScmManagerUsername', passwordVariable: 'ORG_GRADLE_PROJECT_packagesScmManagerPassword'),
    file(credentialsId: 'oss-gpg-secring', variable: 'GPG_KEY_PATH'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'GPG_KEY_ID', passwordVariable: 'GPG_PASSWORD'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'NFPM_RPM_KEY_ID', passwordVariable: 'NFPM_RPM_PASSPHRASE'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'NFPM_DEB_KEY_ID', passwordVariable: 'NFPM_DEB_PASSPHRASE'),
  ]) {
      sh 'gpg --no-tty --batch --yes --import $GPG_KEY_PATH'
  	  closure.call()
  }
}
