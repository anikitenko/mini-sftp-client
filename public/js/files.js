$(function () {
    const localFilesBlock = $(".localFilesBlock");
    localFilesBlock.droppable({
        drop: function (event, ui) {
            let draggablePath = $(ui.draggable).find("span").attr("data-path"),
                draggableName = $(ui.draggable).find("span").attr("data-name"),
                sshIP = $("#sshIp").val(),
                sshUser = $("#sshUser").val(),
                sshPassword = $("#sshPassword").val(),
                sshPort = $("#sshPort").val(),
                notify = $.notify("Downloading...", {
                    type: 'success',
                    showProgressbar: true,
                    delay: 0,
                    allow_dismiss: false,
                    template: '<div data-notify="container" class="col-xs-11 col-sm-3 alert alert-{0}" role="alert">' +
                    '<span data-notify="message">{2}</span>' +
                    '<div class="progress" data-notify="progressbar" style="margin-bottom: 10px;">' +
                    '<div class="progress-bar progress-bar-success progress-bar-striped active" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0;"></div>' +
                    '</div>' +
                    '<div class="text-right"><button data-path="'+
                    draggablePath+'/'+draggableName+
                    '" type="button" class="cancelDownload btn btn-danger ladda-button" data-style="zoom-out">Cancel</button></div>' +
                    '</div>'
                }),
                sourcePath = $("#remotePath").val(),
                localPath = $('#localPath').val(),
                isDir = $(ui.draggable).find("span").attr("data-dir"),
                localRemoteFile = localFilesBlock.find("span[data-name='" + draggableName + "']");
            notify.update('progress', '35');
            if (!localRemoteFile.length) {
                download(sshIP, sshUser, sshPassword, sshPort, sourcePath, draggablePath, localPath, isDir, false, false, ui, notify);
                return
            }

            if (localRemoteFile.attr("data-dir") === "true") {
                notify.update('progress', '100');
                notify.close();
                sendNotify("Since directory already exists we won't download it", "warning");
                return;
            }
            bootbox.dialog({
                title: 'Confirm overwrite',
                message: "<p>File already exists. File will be overwritten, is it okay?</p>",
                buttons: {
                    cancel: {
                        label: "Cancel",
                        className: 'btn-info',
                        callback: function () {
                            notify.update('progress', '100');
                            notify.close();
                            sendNotify("Download was canceled", "success");
                        }
                    },
                    backup: {
                        label: "Backup",
                        className: 'btn-success',
                        callback: function () {
                            download(sshIP, sshUser, sshPassword, sshPort, sourcePath, draggablePath, localPath, isDir, true, true, ui, notify)
                        }
                    },
                    ok: {
                        label: "Overwrite",
                        className: 'btn-warning',
                        callback: function () {
                            download(sshIP, sshUser, sshPassword, sshPort, sourcePath, draggablePath, localPath, isDir, false, true, ui, notify)
                        }
                    }
                },
                onEscape: function () {
                    notify.update('progress', '100');
                    notify.close();
                    sendNotify("Download was canceled", "success");
                }
            });
        }
    });

    $(document).on("click", ".cancelDownload", function () {
        let _this = this,
            filePath = $(_this).attr("data-path"),
            l = Ladda.create(_this);
        l.start();
        bootbox.confirm("Are you sure you want to stop download?", function(result){
            if (!result) {
                l.stop();
                return;
            }
            alert("ok!");
        });
    });

    function download(sshIP, sshUser, sshPassword, sshPort, sourcePath, draggablePath, localPath, isDir, fileToBackup, localFileExists, ui, notify) {
        $.post("/download", {
            ssh_ip: sshIP,
            ssh_user: sshUser,
            ssh_password: sshPassword,
            ssh_port: sshPort,
            source_path: sourcePath,
            file_name: draggablePath,
            local_path: localPath,
            is_dir: isDir,
            backup: fileToBackup
        }, function (response) {
            if (response["result"]) {
                if (!localFileExists) {
                    let newItem = $(ui.draggable).clone();
                    localFilesBlock.append(newItem);
                    if (contextMenuEnabled) {
                        enableContextMenuLocal(newItem);
                    }
                }
                if (fileToBackup) {
                    $(".localRefresh").trigger("click");
                }
            } else {
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function () {
            notify.update('progress', '100');
            //notify.close();
        });
    }
});