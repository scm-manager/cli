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

    stage('Set Version') {
      when {
          branch pattern: 'release/*', comparator: 'GLOB'
      }
      steps {
        // fetch all remotes from origin
        sh 'git config --replace-all "remote.origin.fetch" "+refs/heads/*:refs/remotes/origin/*"'
        sh 'git fetch --all'

		// checkout, reset and merge
		sh 'git checkout main'
		sh 'git reset --hard origin/main'
		sh "git merge --ff-only ${env.BRANCH_NAME}"

        // set tag
        tag releaseVersion
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
        sh 'make build'
      }
    }

	stage('Publish') {
      when {
        branch pattern: 'release/*', comparator: 'GLOB'
        expression { return isBuildSuccess() }
      }
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        withPublishEnvironment {
		  ansiColor('xterm') {
       	    sh 'VERSION=v1.7.0 curl -sL https://git.io/goreleaser | bash -s -- release --rm-dist'
		  }
		  sh "go run pkg/build/upload/app.go dist/scm-cli.json scoop-bucket main scoops/scm-cli.json \"Update scoop scm-cli to ${releaseVersion}\""
		  sh "go run pkg/build/upload/app.go dist/scm-cli.rb homebrew-tap master Formula/scm-cli.rb \"Update brew scm-cli to ${releaseVersion}\""
		  sh "go run pkg/build/descriptor/app.go dist > dist/release.yaml"
		  sh "go run pkg/build/upload/app.go dist/release.yaml website master content/cli/releases/${hyphenatedReleaseVersion}.yaml \"Release cli version ${releaseVersion}\""
        }
      }
	}

    stage('Update Repository') {
	  when {
		branch pattern: 'release/*', comparator: 'GLOB'
	  }
	  steps {
		// merge main in to develop
		sh 'git checkout develop'
		sh 'git merge main'

		// push changes back to remote repository
		authGit 'cesmarvin-github', 'push origin main --tags'
		authGit 'cesmarvin-github', 'push origin develop --tags'
		authGit 'cesmarvin-github', "push origin :${env.BRANCH_NAME}"
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

void withPublishEnvironment(Closure<Void> closure) {
  withCredentials([
    usernamePassword(credentialsId: 'maven.scm-manager.org', usernameVariable: 'UPLOAD_DEFAULT_USERNAME', passwordVariable: 'UPLOAD_DEFAULT_SECRET'),
    usernamePassword(credentialsId: 'maven.scm-manager.org', usernameVariable: 'UPLOAD_RPM_USERNAME', passwordVariable: 'UPLOAD_RPM_SECRET'),
    usernamePassword(credentialsId: 'maven.scm-manager.org', usernameVariable: 'UPLOAD_DEB_USERNAME', passwordVariable: 'UPLOAD_DEB_SECRET'),
    file(credentialsId: 'oss-gpg-secring', variable: 'GPG_KEY_PATH'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'GPG_KEY_ID', passwordVariable: 'GPG_PASSWORD'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'NFPM_RPM_KEY_ID', passwordVariable: 'NFPM_RPM_PASSPHRASE'),
    usernamePassword(credentialsId: 'oss-keyid-and-passphrase', usernameVariable: 'NFPM_DEB_KEY_ID', passwordVariable: 'NFPM_DEB_PASSPHRASE'),
    usernamePassword(credentialsId: 'cesmarvin-github', usernameVariable: 'GITHUB_USERNAME', passwordVariable: 'GITHUB_API_TOKEN'),
  ]) {
      sh 'gpg --no-tty --batch --yes --import $GPG_KEY_PATH'
  	  closure.call()
  }
}

String getHyphenatedReleaseVersion() {
  return getReleaseVersion().replace('.', '-')
}

String getReleaseVersion() {
  return env.BRANCH_NAME.substring("release/".length());
}

void commit(String message) {
  sh "git -c user.name='CES Marvin' -c user.email='cesmarvin@cloudogu.com' commit -m '${message}'"
}

void tag(String version) {
  String message = "Release version ${version}"
  sh "git -c user.name='CES Marvin' -c user.email='cesmarvin@cloudogu.com' tag -m '${message}' ${version}"
}

boolean isBuildSuccess() {
  return currentBuild.result == null || currentBuild.result == 'SUCCESS'
}

void authGit(String credentials, String command) {
  withCredentials([
    usernamePassword(credentialsId: credentials, usernameVariable: 'AUTH_USR', passwordVariable: 'AUTH_PSW')
  ]) {
    sh "git -c credential.helper=\"!f() { echo username='\$AUTH_USR'; echo password='\$AUTH_PSW'; }; f\" ${command}"
  }
}
