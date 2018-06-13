let remotePathInputSearch = "",
    remotePathHistoryBack = [];

$(function () {
    $("#remotePath").select2({
        ajax: {
            type: 'POST',
            url: '/getRemotePathCompletion',
            dataType: 'json',
            delay: 500,
            data: function (params) {
                let sshIP = $("#sshIp").val(),
                    sshUser = $("#sshUser").val(),
                    sshPassword = $("#sshPassword").val(),
                    sshPort = $("#sshPort").val();
                return {
                    path: params.term,
                    ssh_ip: sshIP,
                    ssh_user: sshUser,
                    ssh_password: sshPassword,
                    ssh_port: sshPort
                }
            },
            processResults: function (data, params) {
                if (data["result"]) {
                    params.page = params.page || 1;

                    return {
                        results: data.items
                    }
                } else {
                    sendNotify(data["message"], "danger");
                    return {
                        results: ""
                    }
                }
            }
        },
        escapeMarkup: function (markup) {
            return markup;
        },
        minimumInputLength: 1
    }).on('select2:open', function () {
        $('.select2-search > input').val($(this).val()).trigger('input');
    }).on('select2:select', function () {
        let sshIP = $("#sshIp").val(),
            sshUser = $("#sshUser").val(),
            sshPassword = $("#sshPassword").val(),
            sshPort = $("#sshPort").val(),
            _this = this,
            postPath = $(_this).val(),
            notify = $.notify("Loading remote files...", {
                type: 'success',
                allow_dismiss: false,
                showProgressbar: true,
                delay: 0
            });
        notify.update('progress', '35');
        $.post("/remotePathGoTo", {
            ssh_ip: sshIP,
            ssh_user: sshUser,
            ssh_password: sshPassword,
            ssh_port: sshPort,
            path: postPath
        }, function (response) {
            if (response["result"]) {
                let htmlBlockRemote = "",
                    remoteFiles = response["remote_files"];
                remotePathInputSearch = postPath;
                if (remotePathHistoryBack[0] !== postPath) {
                    remotePathHistoryBack.unshift(postPath);
                }
                $.each(remoteFiles, function (i, val) {
                    let path = val.Path,
                        dir = val.Directory,
                        symlink = val.Symlink;
                    htmlBlockRemote += '<div><span data-dir="' + dir + '" data-symlink="' + symlink + '" data-path="' + postPath + "/" + path + '" data-name="' + path + '">' + path + '</span></div>';
                });
                $(".remoteFilesBlock").html(htmlBlockRemote).find("div").draggable({
                    revert: 'invalid',
                    helper: 'clone',
                    snap: true,
                    snapMode: "inner",
                    cursor: 'move',
                    appendTo: "body",
                    start: function (event, ui) {
                        $(ui.helper).addClass("fileIsMoving")
                    }
                }).each(function () {
                    if ($(this).outerWidth() < $(this).find("span").outerWidth()) {
                        $(this).css("border-right-color", "initial").css("border-right-style", "initial").css("border-right-width", "initial");
                    }
                });
            } else {
                $(_this).select2("trigger", "select", {
                    data: {
                        id: remotePathInputSearch,
                        text: remotePathInputSearch
                    }
                });
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function () {
            notify.update('progress', '100');
            notify.close();
        });
    });

    $(document).on("click", "div.remoteFilesBlock div > span[data-dir='true']", function () {
        const remotePathInput = $('#remotePath');
        let currentPath = remotePathInput.val(),
            separator = "/",
            dirName = $(this).attr("data-name"),
            newPathName = "";
        if (currentPath === "/") {
            separator = ""
        }
        if ($(this).attr("data-symlink") === "true") {
            let dirNameSplit = dirName.split("->");
            dirName = $.trim(dirNameSplit[dirNameSplit.length - 1]);
        }
        if (dirName.charAt(0) === "/") {
            newPathName = dirName;
        } else {
            newPathName = currentPath + separator + dirName;
        }
        remotePathInput.select2("trigger", "select", {
            data: {
                id: newPathName,
                text: newPathName
            }
        });
    });

    $(".remoteGoBack").on("click", function () {
        if (remotePathHistoryBack.length < 2) {
            return;
        }
        const remotePathInput = $('#remotePath');
        remotePathInput.select2("trigger", "select", {
            data: {
                id: remotePathHistoryBack[1],
                text: remotePathHistoryBack[1]
            }
        });
        remotePathHistoryBack.splice(0, 2);
    });
});