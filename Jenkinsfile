#!groovy

notifyBuildDetails = ""
p2pCommitId = ""
cdnHost = ""
dhtHost = ""
gitcmd = ""
p2p_log_level = "INFO"
global_version = "1.0.0"
product_code = "101366BA-A375-46C1-8871-C46D29EE7C70"
dhtSrv = "dht"

switch (env.BRANCH_NAME) {
    case ~/master/: 
        dhtHost = "eu0.mastercdn.subutai.io:6881";
        dhtSrv = "masterdht";
        p2p_log_level = "DEBUG";
        break;
    case ~/dev/:
        dhtHost = "eu0.devcdn.subutai.io:6881";
        dhtSrv = "devdht";
        gitcmd = "git checkout -B dev && git pull origin dev";
        p2p_log_level = "DEBUG";
        break;
    default: 
        dhtHost = "eu0.cdn.subutai.io:6881";
        dhtSrv = "dht";
        break;
}

try {
    notifyBuild('STARTED')

        node("deb") {
            String goenvDir = ".goenv";
            deleteDir();

            stage("Checkout source") {
                notifyBuildDetails = "\nFailed on Stage - Checkout source";
                checkout scm;
                p2pCommitId = sh (script: "git rev-parse HEAD", returnStdout: true);
            }

            stage("Prepare GOENV") {
                /* Creating GOENV path
                   Recreating GOENV path to catch possible issues with external libraries.
                 */
                notifyBuildDetails = "\nFailed on Stage - Prepare GOENV";

                sh """
                    if test ! -d ${goenvDir}; then mkdir -p ${goenvDir}/src/github.com/subutai-io/; fi
                    ln -s ${workspace} ${workspace}/${goenvDir}/src/github.com/subutai-io/p2p
                """;
            }

            stage("Build p2p") {
                /* Build subutai binary */
                notifyBuildDetails = "\nFailed on Stage - Build p2p";

                /* go get golang.org/x/sys/windows */
                sh """
                    export GOPATH=${workspace}/${goenvDir}
                    export GOBIN=${workspace}/${goenvDir}/bin
                    go get
                    go get -u github.com/urfave/cli
                    ./configure --dht=${dhtSrv} --branch=${env.BRANCH_NAME}
                    make all
                """;

                /* stash p2p binary to use it in next node() */
                stash includes: 'bin/p2p.exe', name: 'p2p.exe';
                stash includes: 'bin/p2p_osx', name: 'p2p_osx';
                stash includes: 'upload-ipfs.sh', name: 'upload-ipfs.sh';
            }
        }

    if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {
        node("deb") {
            /* Upload builed p2p artifacts to CDN */
            deleteDir();

            stage("Upload p2p binaries to CDN") {
                /* Get subutai binary from stage and push it to same branch of subos repo
                 */
                notifyBuildDetails = "\nFailed on Stage - Upload p2p binaries to CDN";

                /* upload p2p */
                unstash 'p2p.exe';
                unstash 'p2p_osx';
                unstash 'upload-ipfs.sh';
                if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {
                    sh """
                        set +x
                        ./upload-ipfs.sh ${env.BRANCH_NAME} MSYS_NT-10.0
                        ./upload-ipfs.sh ${env.BRANCH_NAME} Darwin
                        """;
                }
            }
        }

        node("deb") {
            notifyBuild('INFO', "Building Debian Package");
            stage("Building Debian") {
                notifyBuildDetails = "\nFailed on stage - Building Debian Package";

                String date = new Date().format( 'yyyyMMddHHMMSS' );

                def CWD = pwd()

                sh """
                    rm -rf ${CWD}/p2p
                    git clone https://github.com/subutai-io/p2p
                    go get -u github.com/urfave/cli
                """;

                if (env.BRANCH_NAME != 'master') {
                    sh """
                        cd ${CWD}/p2p
                        git checkout --track origin/${env.BRANCH_NAME} && rm -rf .git*
                    """;
                }

                String plain_version = sh (script: """
                        cat ${CWD}/p2p/VERSION | tr -d '\n'
                        """, returnStdout: true);
                def p2p_version = "${plain_version}+${date}";
                global_version = plain_version;

                product_code = sh (script: """
                        cat /proc/sys/kernel/random/uuid | awk '{print toupper(\$0)}' | tr -d '\n'
                        """, returnStdout: true);

                sh """
                    cd ${CWD}/p2p
                    sed -i 's/quilt/native/' debian/source/format
                    sed -i 's/DHT_ENDPOINT/${dhtSrv}/' debian/rules
                    sed -i 's/DEFAULT_LOG_LEVEL/${p2p_log_level}/' debian/rules
                    dch -v '${p2p_version}' -D stable 'Test build for ${p2p_version}' 1>/dev/null 2>/dev/null
                    """;
            }

            stage("Build P2P package") {
                notifyBuildDetails = "\nFailed on Stage - Build package";
                sh """
                    cd p2p
                    dpkg-buildpackage -rfakeroot -us -uc
                    cd ${CWD} || exit 1
                    for i in *.deb; do
                        echo '\$i:';
                dpkg -c \$i;
                done
                    """;
            }

            stage("Upload Packages") {
                notifyBuildDetails = "\nFailed on Stage - Upload";
                sh """
                    cd ${CWD}
                ./upload-ipfs.sh ${env.BRANCH_NAME} Linux
                    touch uploading_agent
                    scp uploading_agent subutai*.deb dak@debup.subutai.io:incoming/${env.BRANCH_NAME}/
                ssh dak@debup.subutai.io sh /var/reprepro/scripts/scan-incoming.sh ${env.BRANCH_NAME} agent
                    """;

                sh """
                    set -x
                    rm -rf /tmp/p2p-packaging
                    git clone git@github.com:optdyn/p2p-packaging.git /tmp/p2p-packaging
                    cd /tmp/p2p-packaging/
                    ${gitcmd}
                cp ${CWD}/subutai*.deb . 
                    ./upload-ipfs.sh ${env.BRANCH_NAME}
                """;
            }
        }

        if (env.BRANCH_NAME == 'dev' || env.BRANCH_NAME == 'master') {

            node("mac") {
                notifyBuild('INFO', "Packaging P2P for Darwin");
                stage("Packaging for Darwin") {
                    notifyBuildDetails = "\nFailed on stage - Starting Darwin Packaging";

                    sh """
                        set -x
                        rm -rf /tmp/p2p-packaging
                        git clone git@github.com:optdyn/p2p-packaging.git /tmp/p2p-packaging
                        cd /tmp/p2p-packaging
                        curl -fsSLk 'https://${env.BRANCH_NAME}bazaar.subutai.io/rest/v1/cdn/raw?name=p2p-${env.BRANCH_NAME}_osx&download&latest' -o /tmp/p2p-packaging/darwin/p2p_osx
                        chmod +x /tmp/p2p-packaging/darwin/p2p_osx
                        /tmp/p2p-packaging/darwin/pack.sh /tmp/p2p-packaging/darwin/p2p_osx ${env.BRANCH_NAME}
                    """;

                    notifyBuildDetails = "\nFailed on stage - Uploading Darwin Package";

                    sh """
                        /tmp/p2p-packaging/./upload-ipfs.sh ${env.BRANCH_NAME}
                    """;
                }
            }
        } // If branch == master

        node("windows") {
            notifyBuild('INFO', "Packaging P2P for Windows");
            stage("Packaging for Windows") {
                notifyBuildDetails = "\nFailed on stage - Starting Windows Packaging";

                bat """
                    if exist "C:\\tmp" RD /S /Q "c:\\tmp"
                        if not exist "C:\\tmp" mkdir "C:\\tmp"
                            echo rm -rf /c/tmp/p2p-packaging > c:\\tmp\\p2p-win.do
                                echo git clone git@github.com:optdyn/p2p-packaging.git /c/tmp/p2p-packaging >> c:\\tmp\\p2p-win.do
                                echo cd /c/tmp/p2p-packaging >> c:\\tmp\\p2p-win.do
                                echo git reset --hard >> c:\\tmp\\p2p-win.do
                                echo git checkout -B ${env.BRANCH_NAME} >> c:\\tmp\\p2p-win.do
                                echo git pull origin ${env.BRANCH_NAME} >> c:\\tmp\\p2p-win.do
                                echo curl -fsSLk "https://${env.BRANCH_NAME}bazaar.subutai.io/rest/v1/cdn/raw?name=p2p-${env.BRANCH_NAME}.exe&download&latest" -o /c/tmp/p2p-packaging/p2p.exe >> c:\\tmp\\p2p-win.do
                                echo curl -fsSLk "https://bazaar.subutai.io/rest/v1/cdn/raw?name=tap-windows-9.21.2.exe&download&latest" -o /c/tmp/p2p-packaging/tap-windows-9.21.2.exe >> c:\\tmp\\p2p-win.do
                                echo sed -i -e "s/{VERSION_PLACEHOLDER}/${global_version}/g" /c/tmp/p2p-packaging/windows/P2PInstaller/P2PInstaller.vdproj >> c:\\tmp\\p2p-win.do
                                echo sed -i -e "s/PRODUCT_CODE_PLACEHOLDER/${product_code}/g" /c/tmp/p2p-packaging/windows/P2PInstaller/P2PInstaller.vdproj >> c:\\tmp\\p2p-win.do

                                echo /c/tmp/p2p-packaging/upload-ipfs.sh ${env.BRANCH_NAME} >> c:\\tmp\\p2p-win-upload.do

                                echo call "C:\\Program Files (x86)\\Microsoft Visual Studio\\2017\\Community\\Common7\\Tools\\VsDevCmd.bat" > c:\\tmp\\p2p-pack.bat
                                echo signtool.exe sign /tr http://timestamp.comodoca.com/authenticode /f "c:\\users\\tray\\od.p12" /p testpassword "c:\\tmp\\p2p-packaging\\p2p.exe" >> c:\\tmp\\p2p-pack.bat
                                echo devenv.com c:\\tmp\\p2p-packaging\\windows\\win.sln /Rebuild Release >> c:\\tmp\\p2p-pack.bat
                                echo signtool.exe sign /tr http://timestamp.comodoca.com/authenticode /f "c:\\users\\tray\\od.p12" /p testpassword "c:\\tmp\\p2p-packaging\\windows\\P2PInstaller\\Release\\P2PInstaller.msi" >> c:\\tmp\\p2p-pack.bat


                                """;

                notifyBuildDetails = "\nFailed on stage - Deploying DevOps";
                bat "c:\\tmp\\p2p-win.do";

                notifyBuildDetails = "\nFailed on stage - Building package";
                bat "c:\\tmp\\p2p-pack.bat";

                notifyBuildDetails = "\nFailed on stage - Uploading Windows package";
                bat "c:\\tmp\\p2p-win-upload.do";
            }
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
    //def slackToken = getSlackToken('p2p-bots')
    // Send notifications
    //slackSend (color: colorCode, message: summary, teamDomain: 'optdyn', token: "${slackToken}")
    def mattermost_rest = "https://mm.optdyn.com/hooks/bixecqjzujg498nyqp9kw8myja"
        mattermostSend(color: colorCode, icon: "https://jenkins.io/images/logos/jenkins/jenkins.png", message: summary, channel: "#p2p-bots", endpoint: "${mattermost_rest}" )
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
