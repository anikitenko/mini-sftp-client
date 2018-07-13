$(function () {
    let remoteHome = "",
        localHome = "",
        connectionName = "New Connection";

    $("#fileContentsModal").on('shown.bs.modal', function () {
        let preCodeBlock = $(this).find("pre code").get(0);
        if (preCodeBlock !== undefined) {
            hljs.highlightBlock($(this).find("pre code").get(0));
        }
    });

    $('#connectionNameDisplay').editable({
        type: 'text',
        mode: 'inline',
        onblur: 'submit',
        display: function () {
            $(this).text(connectionName);
        },
        success: function (response, val) {
            connectionName = $.trim(val);
            document.title = connectionName;
            if (connectionName !== "") {
                sendNotify("Connection name set to '" + connectionName + "'", "success");
            }
        }
    });

    $(".saveCredentialsInCookies").on("click", function (e) {
        let _this = this;
        if ($(this).is(":checked")) {
            e.preventDefault();
            bootbox.confirm({
                message: "<div><span>Please note that your passwords will be <ins>stored unencrypted in cookies</ins>!</span></div>"+
                "<div><span>This is a <ins>serious security issue</ins>!<br>But.. it will be easier for you to connect..</span></div>"+
                "<span>A quick note:</span><ul>"+
                "<li>Passwords which we also keep are ONLY in memory. After you stop client they will disappear (you may check yourself) :)</li>"+
                "<li>So... if this is a 'one-time' connection OR if you are connecting to serious machine better skip this option</li>"+
                "<li>Again.. until you stop your client your passwords will be also saved.. but in more secure way :)</li></ul>",
                title: "Security Warning!",
                buttons: {
                    confirm: {
                        label: 'I totally agree',
                        className: 'btn-danger'
                    },
                    cancel: {
                        label: 'Keep my passwords safe',
                        className: 'btn-success'
                    }
                },
                callback: function (result) {
                    if (!result) {
                        $(_this).prop("checked", false);
                        return
                    }
                    $(_this).prop("checked", true);
                }
            });
        }
    });

    $(".showStoredConnections").on("click", function (e) {
        let _this = this,
            l = Ladda.create(_this);

        if ($(this).parent().find("ul").is(":visible")) {
            return
        }

        e.stopPropagation();
        e.preventDefault();
        l.start();
        $.post("/getStoredConnections", function (response) {
            const dropDownMenu = $(_this).parent().find(".dropdown-menu");
            if (response["result"]) {
                let connections = response["connections"],
                    localConnections = Cookies.getJSON('stored_connections');
                dropDownMenu.empty();
                if (localConnections === undefined) {
                    if ($.isEmptyObject(connections)) {
                        let htmlBlock = "<li><a href='javascript:void(0)'>No stored connections</li>";
                        dropDownMenu.html(htmlBlock);
                    }
                    $.each(connections, function (ip, val) {
                        $.each(val, function (i, ports) {
                            $.each(ports, function (port, users) {
                                $.each(users, function () {
                                    let htmlBlock = "<li data-username='"+this.User+
                                        "' data-password='"+this.Password+
                                        "' data-port='"+port+
                                        "' data-host='"+ip+"'><a class='storedConnection' href='javascript:void(0)'>"+this.User+"@"+ip+":"+port+"</a></li>";
                                    dropDownMenu.append(htmlBlock);
                                });
                            })
                        });
                    })
                } else {
                    $.each(connections, function (ip, val) {
                        $.each(val, function (i, ports) {
                            $.each(ports, function (port, users) {
                                $.each(users, function (i, user) {
                                    let credentialsFound = false,
                                        connectionInformation = {};
                                    $.each(localConnections, function (i, val) {
                                        if (val.Host === ip && val.User === user.User && val.Port === port) {
                                            credentialsFound = true;
                                            return false
                                        }
                                    });
                                    if (!credentialsFound) {
                                        connectionInformation.Host = ip;
                                        connectionInformation.User = user.User;
                                        connectionInformation.Password = user.Password;
                                        connectionInformation.Port = port;
                                        localConnections.push(connectionInformation);
                                    }
                                });
                            });
                        });
                    });
                    $.each(localConnections, function (i, val) {
                        let htmlBlock = "<li data-username='"+val.User+
                            "' data-password='"+val.Password+
                            "' data-port='"+val.Port+
                            "' data-host='"+val.Host+"'><a class='storedConnection' href='javascript:void(0)'>"+val.User+"@"+val.Host+":"+val.Port+"</a></li>";
                        dropDownMenu.append(htmlBlock);
                    });
                }
            } else {
                sendNotify(response["message"], "danger");
                let htmlBlock = "<li><a href='javascript:void(0)'>Error</li>";
                dropDownMenu.html(htmlBlock);
            }
        }, 'json').always(function () {
            l.stop();
            $(_this).dropdown("toggle");
        });
    });

    $(document).on("click", ".storedConnection", function () {
        let sshIp = $(this).parent().attr("data-host"),
            sshUser = $(this).parent().attr("data-username"),
            sshPass = $(this).parent().attr("data-password"),
            sshPort = $(this).parent().attr("data-port");
        $("#sshIp").val(sshIp);
        $("#sshUser").val(sshUser);
        $("#sshPassword").val(sshPass);
        $("#sshPort").val(sshPort);
    });

    $(".mainForm").on("submit", function (e) {
        e.preventDefault();
        $("#sshConnect").trigger("click");
    }).keypress(function(e) {
        if(e.which === 13) {
            e.preventDefault();
            $("#sshConnect").trigger("click");
        }
    });

    $("#sshConnect").on("click", function (e) {
        e.preventDefault();
        let sshIP = $.trim($("#sshIp").val()),
            sshUser = $.trim($("#sshUser").val()),
            sshPassword = $.trim($("#sshPassword").val()),
            sshPort = $.trim($("#sshPort").val()),
            l = Ladda.create(this),
            _this = this;
        l.start();
        if (sshPort === "") {
            sshPort = "22"
        }

        $.post("/connectViaSSH", {
            ssh_ip: sshIP,
            ssh_user: sshUser,
            ssh_password: sshPassword,
            ssh_port: sshPort
        }, function (response) {
            if (response["result"]) {
                let localPath = response["local_path"],
                    remotePath = response["remote_path"],
                    pinCode = response["pin_code"],
                    connectionInformation = {},
                    localConnections = Cookies.getJSON('stored_connections'),
                    localConnectionExists = false;
                connectionInformation.Host = sshIP;
                connectionInformation.User = sshUser;
                connectionInformation.Password = sshPassword;
                connectionInformation.Port = sshPort;
                if (localConnections === undefined) {
                    localConnections = [];
                } else {
                    $.each(localConnections, function (i, val) {
                        if (val.Host === sshIP && val.User === sshUser && val.Port === sshPort) {
                            localConnectionExists = true;
                            return false
                        }
                    });
                }

                if (!localConnectionExists && $(".saveCredentialsInCookies").is(":checked")) {
                    localConnections.push(connectionInformation);
                    Cookies.set("stored_connections", localConnections, {expires: 7});
                }

                remoteHome = response["remote_path"];
                localHome = response["local_path"];
                if (response["errors"] !== null) {
                    sendNotify("We found the following errors: " + response["errors"].join(", "), "warning");
                }

                $("#pinCode").text(pinCode);

                $("#remoteConnectionName").text(sshUser + "@" + sshIP + ":").parent().attr("title", sshUser + "@" + sshIP).tooltip({
                    title: sshUser + "@" + sshIP
                });

                if (remotePath !== "") {
                    $("#remotePath").select2("trigger", "select", {
                        data: {
                            id: remotePath,
                            text: remotePath
                        }
                    });
                }
                if (localPath !== "") {
                    $("#localPath").select2("trigger", "select", {
                        data: {
                            id: localPath,
                            text: localPath
                        }
                    });
                }

                $(".remoteFilesNavigationBlock").css("visibility", "visible");
                $(".localFilesNavigationBlock").css("visibility", "visible");
                remotePathHistoryBack = [];
                localPathHistoryBack = [];
                $("#testSSHConnection").prop("disabled", true);
                $(_this).find("span.ladda-label").text("ReConnect!");
            } else {
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function () {
            l.stop();
        });
    });

    $("#testSSHConnection").on("click", function () {
        let sshIP = $.trim($("#sshIp").val()),
            sshUser = $.trim($("#sshUser").val()),
            sshPassword = $.trim($("#sshPassword").val()),
            sshPort = $.trim($("#sshPort").val()),
            l = Ladda.create(this);
        l.start();
        if (sshPort === "") {
            sshPort = "22"
        }

        $.post("/testSSHConnection", {
            ssh_ip: sshIP,
            ssh_user: sshUser,
            ssh_password: sshPassword,
            ssh_port: sshPort
        }, function (response) {
            if (response["result"]) {
                let pinCode = response["pin_code"],
                    connectionInformation = {},
                    localConnections = Cookies.getJSON('stored_connections'),
                    localConnectionExists = false;
                connectionInformation.Host = sshIP;
                connectionInformation.User = sshUser;
                connectionInformation.Password = sshPassword;
                connectionInformation.Port = sshPort;
                if (localConnections === undefined) {
                    localConnections = [];
                } else {
                    $.each(localConnections, function (i, val) {
                        if (val.Host === sshIP && val.User === sshUser && val.Port === sshPort) {
                            localConnectionExists = true;
                            return false
                        }
                    });
                }

                if (!localConnectionExists && $(".saveCredentialsInCookies").is(":checked")) {
                    localConnections.push(connectionInformation);
                    Cookies.set("stored_connections", localConnections, {expires: 7});
                }

                $("#pinCode").text(pinCode);
                sendNotify("SSH connection was established successfully to '" + sshIP + ":" + sshPort + "'", "success");
            } else {
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function() {
            l.stop();
        });
    });

    $(".remoteGoHome").on("click", function () {
        const remotePathInput = $('#remotePath');
        remotePathInput.select2("trigger", "select", {
            data: {
                id: remoteHome,
                text: remoteHome
            }
        });
    });

    $(".localGoHome").on("click", function () {
        const localPathInput = $('#localPath');
        localPathInput.select2("trigger", "select", {
            data: {
                id: localHome,
                text: localHome
            }
        });
    });

    $(".remoteRefresh").on("click", function () {
        const remotePathInput = $('#remotePath');
        remotePathInput.select2("trigger", "select", {
            data: {
                id: remotePathInput.val(),
                text: remotePathInput.val()
            }
        });
    });

    $(".localRefresh").on("click", function () {
        const localPathInput = $('#localPath');
        localPathInput.select2("trigger", "select", {
            data: {
                id: localPathInput.val(),
                text: localPathInput.val()
            }
        });
    });

    $(".glyphicon-arrow-up").on("click", function () {
        let postPath = $('#localPath').val(),
            remote = false;
        if ($(this).hasClass("remoteGoUp")) {
            postPath = $('#remotePath').val();
            remote = true;
        }

        $.post("/getPath", {path: postPath, remote: remote}, function (response) {
            if (response["result"]) {
                let basePath = response["path"];
                if (remote) {
                    $('#remotePath').select2("trigger", "select", {
                        data: {
                            id: basePath,
                            text: basePath
                        }
                    });
                } else {
                    $('#localPath').select2("trigger", "select", {
                        data: {
                            id: basePath,
                            text: basePath
                        }
                    });
                }
            } else {
                sendNotify(response["message"], "danger");
            }
        }, 'json');
    });

    $(".localCreateNewDir").on("click", function () {
        bootbox.prompt({
            title: "Please enter new directory name",
            callback: function (newDirName) {
                if ($.trim(newDirName) === "") {
                    return
                }
                $.post("/createNewLocalDirectory", {path: $('#localPath').val(), name: newDirName}, function (response) {
                    if (response["result"]) {
                        let newPath = response["new_path"];
                        $('#localPath').select2("trigger", "select", {
                            data: {
                                id: newPath,
                                text: newPath
                            }
                        });
                    } else {
                        sendNotify(response["message"], "danger");
                    }
                }, 'json')
            }
        });
    });

    $("#searchRemoteFiles").on("keyup", function () {
        let input = $(this).val();
        $(".remoteFilesBlock").find("div > span").each(function (index, element) {
            let regex = new RegExp($.trim(input), "gi");
            if ($(element).attr("data-name").match(regex) !== null) {
                $(element).parent().show();
            } else {
                $(element).parent().hide();
            }
        });
    });

    $("#searchLocalFiles").on("keyup", function () {
        let input = $(this).val();
        $(".localFilesBlock").find("div > span").each(function (index, element) {
            let regex = new RegExp($.trim(input), "gi");
            if ($(element).attr("data-name").match(regex) !== null) {
                $(element).parent().show();
            } else {
                $(element).parent().hide();
            }
        });
    });
});