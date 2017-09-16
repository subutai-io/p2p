#!/bin/bash

p2p_binary=$1
if [ ! -f "$p2p_binary" ]; then
    echo "Couldn't find specified file: $p2p_binary"
    exit 1
fi
# clean
rm -rf flat
rm -rf root
rm -f *.pkg
version=`cat ../VERSION`
# Copy files
mkdir -p flat/Resources/en.lproj
mkdir -p flat/base.pkg
mkdir -p root/bin
mkdir -p root/Library/LaunchDaemons
mkdir -p root/etc/newsyslog.d
mkdir -p root/Applications/SubutaiP2P.app/Contents/MacOS
mkdir -p root/Applications/SubutaiP2P.app/Contents/PlugIns
mkdir -p root/Applications/SubutaiP2P.app/Contents/Resources

cp $p2p_binary root/Applications/SubutaiP2P.app/Contents/MacOS/SubutaiP2P
cp $p2p_binary root/bin/p2p
cp io.subutai.p2p.daemon.plist.tmpl root/Library/LaunchDaemons/io.subutai.p2p.daemon.plist
cp p2p.conf.tmpl root/etc/newsyslog.d/p2p.conf

cp PkgInfo.tmpl root/Applications/SubutaiP2P.app/Contents/PkgInfo
cp Info.plist.tmpl root/Applications/SubutaiP2P.app/Contents/Info.plist

# Determine sizes and modify PackageInfo
rootfiles=`find root | wc -l`
rootsize=`du -b -s root | awk '{print $1}'`
mbsize=$(( ${rootsize%% *} / 1024 ))

echo "Size: $rootsize"
echo "MBSize: $mbsize"

cp ./PackageInfo.tmpl ./flat/base.pkg/PackageInfo
sed -i -e "s/{VERSION_PLACEHOLDER}/$version/g" ./flat/base.pkg/PackageInfo
sed -i -e "s/{SIZE_PLACEHOLDER}/$mbsize/g" ./flat/base.pkg/PackageInfo
sed -i -e "s/{FILES_PLACEHOLDER}/$rootfiles/g" ./flat/base.pkg/PackageInfo

# modify Distribution
cp ./Distribution.tmpl ./flat/Distribution
sed -i -e "s/{VERSION_PLACEHOLDER}/$version/g" ./flat/Distribution
sed -i -e "s/{SIZE_PLACEHOLDER}/$mbsize/g" ./flat/Distribution

# Pack and bom
( cd root && find . | cpio -o --format odc --owner 0:80 | gzip -c ) > flat/base.pkg/Payload
( cd scripts && find . | cpio -o --format odc --owner 0:80 | gzip -c ) > flat/base.pkg/Scripts
mkbom -u 0 -g 80 root flat/base.pkg/Bom
( cd flat && xar --compression none -cf "../SubutaiP2P-$version-Installer.pkg" * )
