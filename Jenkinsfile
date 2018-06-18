#!groovy

notifyBuildDetails = ""
p2pCommitId = ""
cdnHost = ""
dhtHost = ""
gitcmd = ""
p2p_log_level = "INFO"

switch (env.BRANCH_NAME) {
	case ~/master/: 
		cdnHost = "mastercdn.subutai.io"; 
		dhtHost = "eu0.mastercdn.subutai.io:6881"
		p2p_log_level = "DEBUG"
		break;
	case ~/dev/:
		cdnHost = "devcdn.subutai.io";
		dhtHost = "eu0.devcdn.subutai.io:6881";
        gitcmd = "git checkout -B dev && git pull origin dev"
		p2p_log_level = "DEBUG"
        break;
	case ~/sysnet/:
		cdnHost = "sysnetcdn.subutai.io";
		dhtHost = "eu0.sysnetcdn.subutai.io:6881";
        gitcmd = "git checkout -B sysnet && git pull origin sysnet "
		p2p_log_level = "TRACE"
        break;
	default: 
		cdnHost = "cdn.subutai.io";
		dhtHost = "eu0.cdn.subutai.io:6881"
		break;
}

try {
	notifyBuild('STARTED')

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
			go get
			go get golang.org/x/sys/windows
			./configure --dht=${dhtHost} --branch=${env.BRANCH_NAME}
			make all
		"""

		/* stash p2p binary to use it in next node() */
		stash includes: 'bin/p2p', name: 'p2p'
		stash includes: 'bin/p2p.exe', name: 'p2p.exe'
		stash includes: 'bin/p2p_osx', name: 'p2p_osx'
	}
	
	/*
	** Trigger subutai-io/snap build on commit to p2p/dev


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
				./bin/p2p -v | cut -d " " -f 3 | tr -d '\n'
				""", returnStdout: true)
			if (env.BRANCH_NAME == 'master') {
				String responseP2P = sh (script: """
					set +x
					curl -s -k ${url}/raw/info?name=p2p
					""", returnStdout: true)
				sh """
					set +x
					curl -s -k -H "token: ${token}" -Ffile=@bin/p2p -Fversion=${p2pVersion} ${url}/raw/upload
				"""
				/* delete old p2p */
				
				if (responseP2P != "Not found") {
					def jsonp2p = jsonParse(responseP2P)
					sh """
						set +x
						curl -s -k -X DELETE ${url}/raw/delete?id=${jsonp2p[0]["id"]}'&'token=${token}
					"""
				}
			}

			/* upload p2p.exe */
			unstash 'p2p.exe'
			String responseP2Pexe = sh (script: """
				set +x
				curl -s -k ${url}/raw/info?name=p2p.exe
				""", returnStdout: true)
			sh """
				set +x
				curl -s -k -H "token: ${token}" -Ffile=@bin/p2p.exe -Fversion=${p2pVersion} ${url}/raw/upload
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
			if (env.BRANCH_NAME == 'master') {
				String responseP2Posx = sh (script: """
					set +x
					curl -s -k ${url}/raw/info?name=p2p_osx
					""", returnStdout: true)
				sh """
					set +x
					curl -s -k -H "token: ${token}" -Ffile=@bin/p2p_osx -Fversion=${p2pVersion} ${url}/raw/upload
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

		node("deb") {
			notifyBuild('INFO', "Building Debian Package")
			stage("Building Debian")
			notifyBuildDetails = "\nFailed on stage - Building Debian Package"

			String date = new Date().format( 'yyyyMMddHHMMSS' )
			
			def CWD = pwd()

			sh """
			rm -rf ${CWD}/p2p
			git clone https://github.com/subutai-io/p2p
			"""

			if (env.BRANCH_NAME != 'master') {
				sh """
				cd ${CWD}/p2p
				git checkout --track origin/${env.BRANCH_NAME} && rm -rf .git*
				"""
			}

			String plain_version = sh (script: """
					cat ${CWD}/p2p/VERSION | tr -d '\n'
					""", returnStdout: true)
			def p2p_version = "${plain_version}+${date}"

			sh """
			cd ${CWD}/p2p
			sed -i 's/quilt/native/' debian/source/format
			sed -i 's/DHT_ENDPOINT/${dhtHost}/' debian/rules
			sed -i 's/DEFAULT_LOG_LEVEL/${p2p_log_level}/' debian/rules
			dch -v '${p2p_version}' -D stable 'Test build for ${p2p_version}' 1>/dev/null 2>/dev/null
			"""

			stage("Build P2P package")
			notifyBuildDetails = "\nFailed on Stage - Build package"
			sh """
			cd p2p
			dpkg-buildpackage -rfakeroot -us -uc
			cd ${CWD} || exit 1
			for i in *.deb; do
				echo '\$i:';
				dpkg -c \$i;
			done
			"""

			stage("Upload Packages")
			notifyBuildDetails = "\nFailed on Stage - Upload"
			sh """
			cd ${CWD}
			touch uploading_agent
			scp uploading_agent subutai*.deb dak@deb.subutai.io:incoming/${env.BRANCH_NAME}/
			ssh dak@deb.subutai.io sh /var/reprepro/scripts/scan-incoming.sh ${env.BRANCH_NAME} agent
			"""

		}

		if (env.BRANCH_NAME == 'master') {
			node("debian") {
				notifyBuild('INFO', "Packaging P2P for Debian")
				stage("Packaging for Debian")
				notifyBuildDetails = "\nFailed on stage - Starting Debian Packaging"

				sh """
					set -x
					rm -rf /tmp/p2p-packaging
					git clone git@github.com:optdyn/p2p-packaging.git /tmp/p2p-packaging
					cd /tmp/p2p-packaging
					${gitcmd}
					wget --no-check-certificate https://eu0.${env.BRANCH_NAME}cdn.subutai.io:8338/kurjun/rest/raw/get?name=p2p -O /tmp/p2p-packaging/linux/debian/p2p
					chmod +x /tmp/p2p-packaging/linux/debian/p2p
					./configure --debian --branch=${env.BRANCH_NAME}
					cd linux
					debuild -B -d
				"""

				notifyBuildDetails = "\nFailed on stage - Uploading Debian Package"

				String debfile = sh (script: """
					set +x
					ls /tmp/p2p-packaging | grep .deb | tr -d '\n'
					""", returnStdout: true)

				sh """
					/tmp/p2p-packaging/upload.sh debian ${env.BRANCH_NAME} /tmp/p2p-packaging/${debfile}
				"""
			}
			
			node("mac") {
				notifyBuild('INFO', "Packaging P2P for Darwin")
				stage("Packaging for Darwin")
				notifyBuildDetails = "\nFailed on stage - Starting Darwin Packaging"

				sh """
					set -x
					rm -rf /tmp/p2p-packaging
					git clone git@github.com:optdyn/p2p-packaging.git /tmp/p2p-packaging
					cd /tmp/p2p-packaging
					curl -fsSLk https://eu0.${env.BRANCH_NAME}cdn.subutai.io:8338/kurjun/rest/raw/get?name=p2p_osx -o /tmp/p2p-packaging/darwin/p2p_osx
					chmod +x /tmp/p2p-packaging/darwin/p2p_osx
					/tmp/p2p-packaging/darwin/pack.sh /tmp/p2p-packaging/darwin/p2p_osx ${env.BRANCH_NAME}
				"""

				notifyBuildDetails = "\nFailed on stage - Uploading Darwin Package"

				sh """
					/tmp/p2p-packaging/upload.sh darwin ${env.BRANCH_NAME} /tmp/p2p-packaging/darwin/p2p.pkg
				"""
			}
		} // If branch == master

		node("windows") {
			notifyBuild('INFO', "Packaging P2P for Windows")
			stage("Packaging for Windows")
			notifyBuildDetails = "\nFailed on stage - Starting Windows Packaging"

			bat """
				if exist "C:\\tmp" RD /S /Q "c:\\tmp"
				if not exist "C:\\tmp" mkdir "C:\\tmp"
				echo rm -rf /c/tmp/p2p-packaging > c:\\tmp\\p2p-win.do
				echo git clone git@github.com:optdyn/p2p-packaging.git /c/tmp/p2p-packaging >> c:\\tmp\\p2p-win.do
				echo cd /c/tmp/p2p-packaging >> c:\\tmp\\p2p-win.do
				echo curl -fsSLk https://eu0.${env.BRANCH_NAME}cdn.subutai.io:8338/kurjun/rest/raw/get?name=p2p.exe -o /c/tmp/p2p-packaging/p2p.exe >> c:\\tmp\\p2p-win.do
				echo curl -fsSLk https://eu0.cdn.subutai.io:8338/kurjun/rest/raw/get?name=tap-windows-9.21.2.exe -o /c/tmp/p2p-packaging/tap-windows-9.21.2.exe >> c:\\tmp\\p2p-win.do

				echo /c/tmp/p2p-packaging/upload.sh windows ${env.BRANCH_NAME} /c/tmp/p2p-packaging/windows/P2PInstaller/Release/P2PInstaller.msi > c:\\tmp\\p2p-win-upload.do

				echo call "C:\\Program Files (x86)\\Microsoft Visual Studio\\2017\\Community\\Common7\\Tools\\VsDevCmd.bat" > c:\\tmp\\p2p-pack.bat
				echo signtool.exe sign /tr http://timestamp.comodoca.com/authenticode /f "c:\\users\\tray\\od.p12" /p testpassword "c:\\tmp\\p2p-packaging\\p2p.exe" >> c:\\tmp\\p2p-pack.bat
				echo devenv.com c:\\tmp\\p2p-packaging\\windows\\win.sln /Rebuild Release >> c:\\tmp\\p2p-pack.bat
				echo signtool.exe sign /tr http://timestamp.comodoca.com/authenticode /f "c:\\users\\tray\\od.p12" /p testpassword "c:\\tmp\\p2p-packaging\\windows\\P2PInstaller\\Release\\P2PInstaller.msi" >> c:\\tmp\\p2p-pack.bat
			"""

			notifyBuildDetails = "\nFailed on stage - Deploying DevOps"
			bat "c:\\tmp\\p2p-win.do"

			notifyBuildDetails = "\nFailed on stage - Building package"
			bat "c:\\tmp\\p2p-pack.bat"

			notifyBuildDetails = "\nFailed on stage - Uploading Windows package"
			bat "c:\\tmp\\p2p-win-upload.do"
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
  def slackToken = getSlackToken('p2p-bots')
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
