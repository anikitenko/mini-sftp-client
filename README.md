# Mini sFTP client

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Cypress.io tests](https://img.shields.io/badge/cypress.io-tests-green.svg?style=flat-square)](https://cypress.io)
[![Build Status](https://travis-ci.org/anikitenko/mini-sftp-client.svg?branch=staging)](https://travis-ci.org/anikitenko/mini-sftp-client)

This is a mini web based sFTP client written on Go using Revel Framework

## Quick Start

### Download and run

* Access [releases page](https://github.com/anikitenko/mini-sftp-client/releases)
and pickup the latest version for your OS
* Download and unzip archive locally
* Run run.exe for Windows OR ./run for Linux/OS X
* The run file will check for any updates
* When prompted enter port to listen on (ex. 9000)
* Access http://127.0.0.1:9000 (if you choosed port 9000) for your browser

### Run from sources

Prerequisite:

* Go 1.6+
* govendor (https://github.com/kardianos/govendor)

Install Revel:

    go get -u github.com/revel/cmd/revel

Get mini sftp client:

    git clone git@github.com:anikitenko/mini-sftp-client.git
    
Resolve dependencies:

    cd mini-sftp-client
    govendor sync
    
Run app:

    cd ..
    revel run mini-sftp-client
    
## Benefits and Key features
- [x] Nothing to install: unzip and run. Use different tabs for different connections
- [x] Runs on Linux, OS X, Windows
- [x] Run client for all interfaces and access client from mobile device and manage files
- [x] Run client on your file server (possible Linux based, Windows, OS X) and access from your desktop or mobile

## Security
Trying to access client from public network (not from localhost) is only possible with pin code. Pin code is shown in the top right corner of the page and in client logs. Pin code is generated each time you start the client and is stored in memory.
    
## Usage

Once you navigate to http://127.0.0.1:9000 you should see the following screen:

![first screen](https://github.com/anikitenko/mini-sftp-client/blob/master/doc-images/first-screen.png)

##### Notes:
* If you are able to authenticate without password on your server, you may ignore that field
* During establishing SSH connection client will try to use .ssh/id_rsa and .ssh/id_dsa if client finds them
* Unsure about connection? Use Test button
* Changing connection name also changes title of the page. 
Open client in a couple of tabs, change connection name and 
you will be able to distinguish different connections

##### Establishing connection:

![connecting](https://github.com/anikitenko/mini-sftp-client/blob/master/doc-images/connecting.png)

##### Notes:
* After you successfully established connection, client will try to detect home directories no matter what is your local OS
* Button to Test connection is disabled after successful connection
    * This is because you can enter credentials to 1 server and if you test 
    for another, input data remains and button to ReConnect also remains,
     so silly click on it will cause all data to load from your another server
    
![like double tab](https://github.com/anikitenko/mini-sftp-client/blob/master/doc-images/like-double-tab.gif)


![like double tab local](https://github.com/anikitenko/mini-sftp-client/blob/master/doc-images/like-double-tab-local.gif)

##### Notes:
* "Like double tab" works on Windows, OS X and Linux OS

##### Downloading files and using search:

![download and search](https://github.com/anikitenko/mini-sftp-client/blob/master/doc-images/download-search.gif)

##### Notes:
* Search works the same for remote files as for local
* Search will not search for files globally, it's only sorting files which are exist

##### Quick buttons
* Quick buttons for remote files:
    * Go Back: every time you navigate, client will save paths and on click button will return you to previous path
    * Go Home: button will navigate you to initial directory which was opened during connection
    * Go UP: navigates you to parent directory
    * Refresh: refresh current directory
* Quick buttons for local files:
    * Go Back: every time you navigate, client will save paths and on click button will return you to previous path
    * Go Home: button will navigate you to initial directory which was opened during connection
    * Go UP: navigates you to parent directory
    * Create New Directory: create new empty directory and navigate to it
    * Refresh: refresh current directory
