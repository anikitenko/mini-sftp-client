let localPathSeparator = "/",
    remoteHome = "",
    localHome = "",
    connectionName = "New Connection";

window.onbeforeunload = function () {
    return "Are you sure you wish to leave the page?";
};

$.notifyDefaults({
    newest_on_top: true
});

$('[data-toggle="tooltip"]').tooltip({
    trigger: 'hover'
});

$(document).ajaxError(function (event, jqxhr) {
    $.notifyClose('top-right');
    if (jqxhr.status === 500 || jqxhr.status === 502) {
        $.notify(
            {
                message: "Something went wrong!",
                icon: 'glyphicon glyphicon-warning-sign'
            }, {type: 'danger', timer: 50});
    } else if (jqxhr.status === 404) {
        $.notify(
            {
                message: "Requested resource not found!",
                icon: 'glyphicon glyphicon-warning-sign'
            }, {type: 'danger', timer: 50});
    } else if (jqxhr.status === 403) {
        $.notify(
            {
                message: "You don't have permission to resource you are trying to access!",
                icon: 'glyphicon glyphicon-warning-sign'
            }, {type: 'danger', timer: 50});
    }
});

$(function () {
    $('#connectionNameDisplay').editable({
        type: 'text',
        mode: 'inline',
        display: function (value) {
            $(this).text(connectionName);
        },
        success: function (response, val) {
            connectionName = $.trim(val);
            document.title = connectionName;
            $.notify("Connection name set to '" + connectionName + "'", {type: 'success', timer: 50});
        }
    });

    $("#sshConnect").on("click", function (e) {
        e.preventDefault();
        let sshIP = $("#sshIp").val(),
            sshUser = $("#sshUser").val(),
            sshPassword = $("#sshPassword").val(),
            sshPort = $("#sshPort").val(),
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
                    htmlBlockLocal = "";
                localPathSeparator = response["local_path_separator"];
                remoteHome = response["remote_path"];
                localHome = response["local_path"];
                if (response["errors"] !== null) {
                    $.notify(
                        {
                            message: "We found the following errors: " + response["errors"].join(", "),
                            icon: 'glyphicon glyphicon-warning-sign'
                        },
                        {
                            type: 'warning',
                            timer: 50
                        }
                    );
                }
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
                $.notify(
                    {
                        message: response["message"],
                        icon: 'glyphicon glyphicon-warning-sign'
                    },
                    {
                        type: 'danger',
                        timer: 50
                    }
                );
            }
        }, 'json').always(function () {
            l.stop();
        });
    });

    $("#testSSHConnection").on("click", function () {
        let sshIP = $("#sshIp").val(),
            sshUser = $("#sshUser").val(),
            sshPassword = $("#sshPassword").val(),
            sshPort = $("#sshPort").val(),
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
                $.notify("SSH connection was established successfully to '" + sshIP + ":" + sshPort + "'", {
                    type: 'success',
                    timer: 50
                });
            } else {
                $.notify(
                    {
                        message: response["message"],
                        icon: 'glyphicon glyphicon-warning-sign'
                    },
                    {
                        type: 'danger',
                        timer: 50
                    }
                );
            }
        }, 'json').always(function () {
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
                $.notify(
                    {
                        message: response["message"],
                        icon: 'glyphicon glyphicon-warning-sign'
                    },
                    {
                        type: 'danger',
                        timer: 50
                    }
                );
            }
        }, 'json');
    });

    $(".localCreateNewDire").on("click", function () {
        let newDirName = prompt("Please enter new directory name");
        newDirName = $.trim(newDirName);
        if (newDirName === "") {
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
                $.notify(
                    {
                        message: response["message"],
                        icon: 'glyphicon glyphicon-warning-sign'
                    },
                    {
                        type: 'danger',
                        timer: 50
                    }
                );
            }
        }, 'json');
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