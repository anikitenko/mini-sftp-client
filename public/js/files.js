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
                confirmOverwrite = true,
                sourcePath = $("#remotePath").val(),
                localPath = $('#localPath').val(),
                isDir = $(ui.draggable).find("span").attr("data-dir"),
                localFileExists = false,
                localRemoteFile = localFilesBlock.find("span[data-name='" + draggableName + "']");
            notify.update('progress', '35');
            if (localRemoteFile.length) {
                if (localRemoteFile.attr("data-dir") === "true") {
                    notify.update('progress', '100');
                    notify.close();
                    $.notify(
                        {
                            message: "Since directory already exists we won't download it",
                            icon: 'glyphicon glyphicon-warning-sign'
                        },
                        {
                            type: 'warning',
                            timer: 50
                        }
                    );
                    return;
                }
                localFileExists = true;
                confirmOverwrite = confirm("File already exists. File will be overwritten, is it okay?")
            }
            if (!confirmOverwrite) {
                notify.update('progress', '100');
                notify.close();
                $.notify(
                    {
                        message: "Download was canceled"
                    },
                    {
                        type: 'success',
                        timer: 50
                    }
                );
                return;
            }
            $.post("/download", {
                ssh_ip: sshIP,
                ssh_user: sshUser,
                ssh_password: sshPassword,
                ssh_port: sshPort,
                source_path: sourcePath,
                file_name: draggablePath,
                local_path: localPath,
                is_dir: isDir
            }, function (response) {
                if (response["result"]) {
                    if (!localFileExists) {
                        localFilesBlock.append($(ui.draggable).clone());
                        $("#localPathGo").trigger("click");
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
            }, 'json').always(function () {
                notify.update('progress', '100');
                notify.close();
            });
        }
    });
});