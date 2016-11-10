#!groovy

notifyBuildDetails = ""
p2pCommitId = ""

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
			export GOPATH=${workspace}/${goenvDir}
			export GOBIN=${workspace}/${goenvDir}/bin
			make
		"""

		/* stash p2p binary to use it in next node() */
		stash includes: 'p2p', name: 'p2p'
	}

	node() {
		/* Checkout subos repo and push new subutai binary */
		deleteDir()

		stage("Push new p2p binary to subos repo")
		/* Get subutai binary from stage and push it to same branch of subos repo
		*/
		notifyBuildDetails = "\nFailed on Stage - Push new subutai binary to subos repo"

		String subosRepoName = "github.com/subutai-io/subos.git"

		git branch: "${env.BRANCH_NAME}", changelog: false, credentialsId: 'hub-optdyn-github-auth', poll: false, url: "https://${subosRepoName}"

		dir("p2p/bin") {
			unstash 'p2p'
		}

		withCredentials([[$class: 'UsernamePasswordMultiBinding', 
			credentialsId: 'hub-optdyn-github-auth', 
			passwordVariable: 'GIT_PASSWORD', 
			usernameVariable: 'GIT_USER']]) {
			sh """
				git config user.email jenkins@subut.ai
				git config user.name 'Jenkins Admin'
				git commit p2p/bin/p2p -m 'Push subutai version from subutai-io/p2p@${p2pCommitId}'
				git push https://${env.GIT_USER}:'${env.GIT_PASSWORD}'@${subosRepoName} ${env.BRANCH_NAME}
			"""
		}
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
  } else if (buildStatus == 'SUCCESSFUL') {
    color = 'GREEN'
    colorCode = '#00FF00'
  } else {
    color = 'RED'
    colorCode = '#FF0000'
	summary = "${subject} (${env.BUILD_URL})${details}"
  }
  // Get token
  def slackToken = getSlackToken('sysnet-bots-slack-token')
  // Send notifications
  slackSend (color: colorCode, message: summary, teamDomain: 'subutai-io', token: "${slackToken}")
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