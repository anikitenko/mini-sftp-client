let localPathSeparator = "/",
    remoteHome = "",
    localHome = "",
    connectionName = "New Connection";

window.onbeforeunload = function () {
    const sshIPValue = $.trim($("#sshIp").val());
    if (sshIPValue !== "" && sshIPValue !== undefined) {
        if (sshIPValue !== "sftp-mock-test") {
            return "Are you sure you wish to leave the page?";
        }
    }
};

$.notifyDefaults({
    newest_on_top: true
});

$.notifyClose('top-right');

$('[data-toggle="tooltip"]').tooltip({
    trigger: 'hover'
});

function sendNotify(text, type) {
    let icon = "";
    if (type === "warning" || type === "danger") {
        icon = "glyphicon glyphicon-warning-sign"
    }
    $.notify({
            message: text,
            icon: icon
        }, {
            type: type,
            timer: 50
        });
}

$(document).ajaxError(function (event, jqxhr) {
    if (jqxhr.status === 500 || jqxhr.status === 502) {
        sendNotify("Something went wrong!", "danger");
    } else if (jqxhr.status === 404) {
        sendNotify("Requested resource not found!", "danger");
    } else if (jqxhr.status === 403) {
        sendNotify("You don't have permission to resource you are trying to access!", "danger");
    }
});

$(function () {
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

    $(".mainForm").on("submit", function (e) {
        e.preventDefault();
        $("#sshConnect").trigger("click");
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
                    remotePath = response["remote_path"];
                localPathSeparator = response["local_path_separator"];
                remoteHome = response["remote_path"];
                localHome = response["local_path"];
                if (response["errors"] !== null) {
                    sendNotify("We found the following errors: " + response["errors"].join(", "), "warning");
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
                sendNotify(response["message"], "danger");
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