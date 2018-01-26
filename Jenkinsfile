#!groovy

notifyBuildDetails = ""
p2pCommitId = ""
cdnHost = ""
dhtHost = ""
gitcmd = ""

switch (env.BRANCH_NAME) {
	case ~/master/: 
		cdnHost = "mastercdn.subut.ai"; 
		dhtHost = "54.93.172.70:6881"
		break;
	case ~/dev/:
		cdnHost = "devcdn.subut.ai";
		dhtHost = "18.195.169.215:6881";
        gitcmd = "git checkout -B dev && git pull origin dev"
        break;
	case ~/sysnet/:
		cdnHost = "devcdn.subut.ai";
		dhtHost = "18.195.169.215:6881";
        gitcmd = "git checkout -B sysnet && git pull origin sysnet "
        break;
	default: 
		cdnHost = "devcdn.subut.ai";
		dhtHost = "mdht.subut.ai:6881"
		break;
}

try {
	notifyBuild('STARTED')

	/* Building agent binary.
	Node block used to separate agent and subos code.
	*/
	node() {
		String goenvDir = ".goenv"
		deleteDir()

		stage("Checkout source")
		/* checkout agent repo */
		notifyBuildDetails = "\nFailed on Stage - Checkout source"

		checkout scm

		p2pCommitId = sh (script: "git rev-parse HEAD", returnStdout: true)

		stage("Prepare GOENV")
		/* Creating GOENV path
		Recreating GOENV path to catch possible issues with external libraries.
		*/
		notifyBuildDetails = "\nFailed on Stage - Prepare GOENV"

		sh """
			if test ! -d ${goenvDir}; then mkdir -p ${goenvDir}/src/github.com/subutai-io/; fi
			ln -s ${workspace} ${workspace}/${goenvDir}/src/github.com/subutai-io/p2p
		"""

		stage("Build p2p")
		/* Build subutai binary */
		notifyBuildDetails = "\nFailed on Stage - Build p2p"

		sh """
            ./configure --dht=${dhtHost} --branch=${env.BRANCH_NAME}
		"""

		sh """
			export GOPATH=${workspace}/${goenvDir}
			export GOBIN=${workspace}/${goenvDir}/bin
			go get
			go get golang.org/x/sys/windows
			make all
		"""

		/* stash p2p binary to use it in next node() */
		stash includes: 'p2p', name: 'p2p'
		stash includes: 'p2p.exe', name: 'p2p.exe'
		stash includes: 'p2p_osx', name: 'p2p_osx'

		if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {
			for (codeName in [ 'trusty', 'vivid', 'xenial', 'yakkety', 'zesty']) {
				sh """
					find ../ -maxdepth 1 -type f -name subutai-p2p*.dsc -delete
					find ../ -maxdepth 1 -type f -name subutai-p2p*.build -delete
					find ../ -maxdepth 1 -type f -name subutai-p2p*.tar.gz -delete
					find ../ -maxdepth 1 -type f -name subutai-p2p*.changes -delete
					find ../ -maxdepth 1 -type f -name subutai-p2p*.ppa.upload -delete
					#./configure --maintainer='Jenkins Admin' --maintainer-email='jenkins@subut.ai' --debian-release=${codeName} --scheme=${env.BRANCH_NAME} --version-postfix=${env.BUILD_NUMBER}
					#make debian-source
					#dput ppa:subutai-social/subutai \$(ls ../subutai-p2p*changes)
				"""			
			}
		}
	}

	// if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {
	// 	node() {
	// 		/* Checkout subos repo and push new subutai binary */
	// 		deleteDir()

	// 		stage("Push new p2p binary to subos repo")
	// 		/* Get subutai binary from stage and push it to same branch of subos repo
	// 		*/
	// 		notifyBuildDetails = "\nFailed on Stage - Push new subutai binary to subos repo"

	// 		String subosRepoName = "github.com/subutai-io/subos.git"

	// 		git branch: "${env.BRANCH_NAME}", changelog: false, credentialsId: 'hub-optdyn-github-auth', poll: false, url: "https://${subosRepoName}"

	// 		dir("p2p/bin") {
	// 			unstash 'p2p'
	// 		}

	// 		withCredentials([[$class: 'UsernamePasswordMultiBinding', 
	// 			credentialsId: 'hub-optdyn-github-auth', 
	// 			passwordVariable: 'GIT_PASSWORD', 
	// 			usernameVariable: 'GIT_USER']]) {
	// 			sh """
	// 				git config user.email jenkins@subut.ai
	// 				git config user.name 'Jenkins Admin'
	// 				git commit p2p/bin/p2p -m 'Push subutai version from subutai-io/p2p@${p2pCommitId}'
	// 				git push https://${env.GIT_USER}:'${env.GIT_PASSWORD}'@${subosRepoName} ${env.BRANCH_NAME}
	// 			"""
	// 		}
	// 	}
	// }
	
	/*
	** Trigger subutai-io/snap build on commit to p2p/dev
	*/

    /*
	if (env.BRANCH_NAME == 'dev') {
		build job: 'snap.subutai-io.pipeline/dev/', propagate: false, wait: false
	}
	
	if (env.BRANCH_NAME == 'master') {
		build job: 'snap.subutai-io.pipeline/master/', propagate: false, wait: false
	}
    */

	if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {
		node() {
			/* Upload builed p2p artifacts to kurjun */
			deleteDir()

			stage("Upload p2p binaries to kurjun")
			/* Get subutai binary from stage and push it to same branch of subos repo
			*/
			notifyBuildDetails = "\nFailed on Stage - Upload p2p binaries to kurjun"

			/* cdn auth creadentials */
			String url = "https://${cdnHost}:8338/kurjun/rest"
			String user = "jenkins"
			def authID = sh (script: """
				set +x
				curl -s -k ${url}/auth/token?user=${user} | gpg --clearsign --no-tty
				""", returnStdout: true)
			def token = sh (script: """
				set +x
				curl -s -k -Fmessage=\"${authID}\" -Fuser=${user} ${url}/auth/token
				""", returnStdout: true)

			/* upload p2p */
			unstash 'p2p'
			/* get p2p version */
			String p2pVersion = sh (script: """
				set +x
				./p2p -v | cut -d " " -f 3 | tr -d '\n'
				""", returnStdout: true)
			String responseP2P = sh (script: """
				set +x
				curl -s -k ${url}/raw/info?name=p2p
				""", returnStdout: true)
			sh """
				set +x
				curl -s -k -H "token: ${token}" -Ffile=@p2p -Fversion=${p2pVersion} ${url}/raw/upload
			"""
			/* delete old p2p */
			if (responseP2P != "Not found") {
				def jsonp2p = jsonParse(responseP2P)
				sh """
					set +x
					curl -s -k -X DELETE ${url}/raw/delete?id=${jsonp2p[0]["id"]}'&'token=${token}
				"""
			}

			/* upload p2p.exe */
			unstash 'p2p.exe'
			String responseP2Pexe = sh (script: """
				set +x
				curl -s -k ${url}/raw/info?name=p2p.exe
				""", returnStdout: true)
			sh """
				set +x
				curl -s -k -H "token: ${token}" -Ffile=@p2p.exe -Fversion=${p2pVersion} ${url}/raw/upload
			"""
			/* delete old p2p.exe */
			if (responseP2Pexe != "Not found") {
				def jsonp2pexe = jsonParse(responseP2Pexe)
				sh """
					set +x
					curl -s -k -X DELETE ${url}/raw/delete?id=${jsonp2pexe[0]["id"]}'&'token=${token}
				"""
			}

			/* upload p2p_osx */
			unstash 'p2p_osx'
			String responseP2Posx = sh (script: """
				set +x
				curl -s -k ${url}/raw/info?name=p2p_osx
				""", returnStdout: true)
			sh """
				set +x
				curl -s -k -H "token: ${token}" -Ffile=@p2p_osx -Fversion=${p2pVersion} ${url}/raw/upload
			"""
			/* delete old p2p */
			if (responseP2Posx != "Not found") {
				def jsonp2posx = jsonParse(responseP2Posx)
				sh """
					set +x
					curl -s -k -X DELETE ${url}/raw/delete?id=${jsonp2posx[0]["id"]}'&'token=${token}
				"""
			}
		}
	}
/*
    node("windows") {
        Stage("Packaging for Windows")
        notifyBuildDetails = "\nFailed on stage - Starting Windows Packaging"
    }
    */

    node("debian") {
        notifyBuild('INFO', "Packaging P2P for Debian")
        stage("Packaging for Debian")
        notifyBuildDetails = "\nFailed on stage - Starting Debian Packaging"

        sh """
            set -x
            rm -rf /tmp/devops
            git clone git@github.com:optdyn/devops.git /tmp/devops
            cd /tmp/devops
            ${gitcmd}
            cd /tmp/devops/p2p
            wget https://eu0.${env.BRANCH_NAME}cdn.subut.ai:8338/kurjun/rest/raw/get?name=p2p -O /tmp/devops/p2p/linux/debian/p2p
            chmod +x /tmp/devops/p2p/linux/debian/p2p
            ./configure --debian --branch=${env.BRANCH_NAME}
            cd linux
            debuild -B -d
        """

        stage("Uploading debian packag")
        notifyBuildDetails = "\nFailed on stage - Uploading Debian Package"

		String debfile = sh (script: """
			set +x
			ls /tmp/devops/p2p | grep .deb | tr -d '\n'
			""", returnStdout: true)

        sh """
            /tmp/devops/p2p/upload.sh debian ${env.BRANCH_NAME} ${debfile}
        """
    }

} catch (e) { 
	currentBuild.result = "FAILED"
	throw e
} finally {
	// Success or failure, always send notifications
	notifyBuild(currentBuild.result, notifyBuildDetails)
}

// https://jenkins.io/blog/2016/07/18/pipline-notifications/
def notifyBuild(String buildStatus = 'STARTED', String details = '') {
  // build status of null means successful
  buildStatus = buildStatus ?: 'SUCCESSFUL'

  // Default values
  def colorName = 'RED'
  def colorCode = '#FF0000'
  def subject = "${buildStatus}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'"  	
  def summary = "${subject} (${env.BUILD_URL})"

  // Override default values based on build status
  if (buildStatus == 'STARTED') {
    color = 'YELLOW'
    colorCode = '#FFFF00'  
  } else if (buildStatus == 'INFO') {
    color = 'GREY'
    colorCode = '#555555'
    summary = "${subject}: ${details}"
  } else if (buildStatus == 'SUCCESSFUL') {
    color = 'GREEN'
    colorCode = '#00FF00'
  } else {
    color = 'RED'
    colorCode = '#FF0000'
	summary = "${subject} (${env.BUILD_URL})${details}"
  }
  // Get token
  def slackToken = getSlackToken('sysnet')
  // Send notifications
  slackSend (color: colorCode, message: summary, teamDomain: 'optdyn', token: "${slackToken}")
}

// get slack token from global jenkins credentials store
@NonCPS
def getSlackToken(String slackCredentialsId){
	// id is ID of creadentials
	def jenkins_creds = Jenkins.instance.getExtensionList('com.cloudbees.plugins.credentials.SystemCredentialsProvider')[0]

	String found_slack_token = jenkins_creds.getStore().getDomains().findResult { domain ->
	  jenkins_creds.getCredentials(domain).findResult { credential ->
	    if(slackCredentialsId.equals(credential.id)) {
	      credential.getSecret()
	    }
	  }
	}
	return found_slack_token
}

@NonCPS
def jsonParse(def json) {
    new groovy.json.JsonSlurperClassic().parseText(json)
}
