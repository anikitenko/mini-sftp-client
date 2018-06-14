#!/usr/bin/env bash

rm -rf node_modules
revel build github.com/anikitenko/mini-sftp-client mini-sftp-osx
GOOS=linux revel build github.com/anikitenko/mini-sftp-client mini-sftp-linux
GOOS=windows revel build github.com/anikitenko/mini-sftp-client mini-sftp-windows

rm -rf mini-sftp-osx/src mini-sftp-osx/run.sh mini-sftp-osx/run.bat
rm -rf mini-sftp-linux/src mini-sftp-linux/run.sh mini-sftp-linux/run.bat
rm -rf mini-sftp-windows/src mini-sftp-windows/run.sh mini-sftp-windows/run.bat

go build -o ./mini-sftp-osx/run ./run && chmod +x ./mini-sftp-osx/run
GOOS=linux go build -o ./mini-sftp-linux/run ./run && chmod +x ./mini-sftp-linux/run
GOOS-windows go build -o ./mini-sftp-windows/run.exe ./run

mv mini-sftp-osx/mini-sftp-client mini-sftp-osx/mini-sftp-client-darwin
mv mini-sftp-linux/mini-sftp-client mini-sftp-linux/mini-sftp-client-linux
mv mini-sftp-windows/mini-sftp-client.exe mini-sftp-windows/mini-sftp-client-windows.exe

cp mini-sftp-osx/mini-sftp-client-darwin artifacts/
cp mini-sftp-linux/mini-sftp-client-linux artifacts/
cp mini-sftp-windows/mini-sftp-client-windows.exe artifacts/

zip mini-sftp-osx.zip mini-sftp-osx
zip mini-sftp-linux.zip mini-sftp-linux
zip mini-sftp-windows.zip mini-sftp-windows

mv mini-sftp-osx.zip artifacts/
mv mini-sftp-linux.zip artifacts/
mv mini-sftp-windows.zip artifacts/