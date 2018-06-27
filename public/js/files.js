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
                    allow_dismiss: false,
                    showProgressbar: true,
                    delay: 0
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
                }
            });
        }
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
                    localFilesBlock.append($(ui.draggable).clone());
                }
                if (fileToBackup) {
                    $(".localRefresh").trigger("click");
                }
            } else {
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function () {
            notify.update('progress', '100');
            notify.close();
        });
    }
});