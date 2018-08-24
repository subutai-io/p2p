BRANCH=$1
OS=$2
PKGNAME=""
BINNAME=""
FILE=""
URL=""
PKG_EXT=""
BINARY_EXT=""

upload_ipfs (){
    filename=$1
    user="jenkins@optimal-dynamics.com"
    fingerprint="877B586E74F170BC4CF6ECABB971E2AC63D23DC9"
    cdnHost=$2
    echo $filename
    extract_id()
        {
            id_src=$(echo $json | grep "id")
            id=${id_src:10:46}
        }       

    json=`curl -k -s -X GET ${cdnHost}/rest/v1/cdn/raw?name=$filename`
    echo "Received: $json"
    extract_id
    echo "Previous file ID is $id"

    authId="$(curl -s ${cdnHost}/rest/v1/cdn/token?fingerprint=${fingerprint})"
    echo "Auth id obtained and signed $authId"

    sign="$(echo ${authId} | gpg --clearsign -u ${user})"
    token="$(curl -s --data-urlencode "request=${sign}"  ${cdnHost}/rest/v1/cdn/token)"
    echo "Token obtained $token"

    echo "Uploading file..."
    curl -sk -H "token: ${token}" -Ffile=@$filename -Ftoken=${token} -X POST "${cdnHost}/rest/v1/cdn/uploadRaw"

    echo "Removing previous"
    if [[ -z "$id" ]]; then
        echo "File not found"
    else curl -k -s -X DELETE "$cdnHost/rest/v1/cdn/raw?token=${token}&id=$id"
    fi
    echo -e "\\nCompleted"
}

case $OS in
    Linux)
        BASENAME="p2p"
        BIN_EXT=""
        ;;
    MSYS_NT-10.0)
        BASENAME="p2p.exe"
        BIN_EXT=".exe"
        ;;
    Darwin)
        BASENAME="p2p_osx"
        BIN_EXT="_osx"
        ;;
esac

case $BRANCH in
    dev)
        BINNAME="p2p-dev$BIN_EXT"
        cd bin
	    cp $BASENAME $BINNAME
        IPFSURL=https://devbazaar.subutai.io
        upload_ipfs $BINNAME $IPFSURL
        ;;
    master)
        BINNAME="p2p-master$BIN_EXT"
        cd bin
	    cp $BASENAME $BINNAME
        IPFSURL=https://masterbazaar.subutai.io
        upload_ipfs $BINNAME $IPFSURL
        ;;
    head)
        BINNAME="p2p$BIN_EXT"
        if [ $OS = Linux ] || [$OS = MSYS_NT-10.0 ]
        then
        cd bin
	    cp $BASENAME $BINNAME
        IPFSURL=https://bazaar.subutai.io
        upload_ipfs $BINNAME $IPFSURL
        fi
        ;;
    HEAD)
        BINNAME="subutai-p2p$PKG_EXT"
        if [ $OS = Linux ] || [$OS = MSYS_NT-10.0 ]
        then
        cd bin
	    cp $BASENAME $BINNAME
        IPFSURL=https://bazaar.subutai.io
        upload_ipfs $BINNAME $IPFSURL
        fi
        ;;
esac

echo "---------"
echo $BINNAME
echo $OS
echo $BRANCH
echo $VERSION
echo "---------"