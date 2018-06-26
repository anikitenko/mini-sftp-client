#!/usr/bin/env bash

cd "$TRAVIS_BUILD_DIR"

rm -rf node_modules
CGO_ENABLED=0 revel build github.com/anikitenko/mini-sftp-client ../mini-sftp-osx
CGO_ENABLED=0 GOOS=linux revel build github.com/anikitenko/mini-sftp-client ../mini-sftp-linux
CGO_ENABLED=0 GOOS=windows revel build github.com/anikitenko/mini-sftp-client ../mini-sftp-windows

cd run
govendor sync

CGO_ENABLED=0 go build -o ../run-osx && chmod +x ../run-osx
CGO_ENABLED=0 GOOS=linux go build -o ../run-linux && chmod +x ../run-linux
CGO_ENABLED=0 GOOS=windows go build -o ../run.exe

cd ../..
mkdir artifacts

rm -f mini-sftp-osx/run.sh mini-sftp-osx/run.bat
rm -f mini-sftp-linux/run.sh mini-sftp-linux/run.bat
rm -f mini-sftp-windows/run.sh mini-sftp-windows/run.bat

find 'mini-sftp-osx/src/github.com/anikitenko/mini-sftp-client' \
 -maxdepth 1 ! -path 'mini-sftp-osx/src/github.com/anikitenko/mini-sftp-client' \
 -not -name 'app' \
 -not -name 'conf' \
 -not -name 'public' \
 -exec rm -rf {} +

find 'mini-sftp-linux/src/github.com/anikitenko/mini-sftp-client' \
 -maxdepth 1 ! -path 'mini-sftp-linux/src/github.com/anikitenko/mini-sftp-client' \
 -not -name 'app' \
 -not -name 'conf' \
 -not -name 'public' \
 -exec rm -rf {} +

find 'mini-sftp-windows/src/github.com/anikitenko/mini-sftp-client' \
 -maxdepth 1 ! -path 'mini-sftp-windows/src/github.com/anikitenko/mini-sftp-client' \
 -not -name 'app' \
 -not -name 'conf' \
 -not -name 'public' \
 -exec rm -rf {} +

mv run-osx mini-sftp-osx/run
mv run-linux mini-sftp-linux/run
mv run.exe mini-sftp-windows/run.exe

mv mini-sftp-osx/mini-sftp-client mini-sftp-osx/mini-sftp-client-darwin
mv mini-sftp-linux/mini-sftp-client mini-sftp-linux/mini-sftp-client-linux
mv mini-sftp-windows/mini-sftp-client.exe mini-sftp-windows/mini-sftp-client-windows.exe

cp mini-sftp-osx/mini-sftp-client-darwin artifacts/
cp mini-sftp-linux/mini-sftp-client-linux artifacts/
cp mini-sftp-windows/mini-sftp-client-windows.exe artifacts/

version="$(echo $TRAVIS_TAG | sed -e 's/^v//g')"
echo $version > mini-sftp-osx/.version
echo $version > mini-sftp-linux/.version
echo $version > mini-sftp-windows/.version

zip -r mini-sftp-osx.zip mini-sftp-osx
zip -r mini-sftp-linux.zip mini-sftp-linux
zip -r mini-sftp-windows.zip mini-sftp-windows

mv mini-sftp-osx.zip artifacts/
mv mini-sftp-linux.zip artifacts/
mv mini-sftp-windows.zip artifacts/