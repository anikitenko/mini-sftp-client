let localPathInputSearch = "",
    localPathHistoryBack = [],
    localPathSeparator = "/";

$(function () {
    $("#localPath").select2({
        ajax: {
            type: 'POST',
            url: '/getLocalPathCompletion',
            dataType: 'json',
            delay: 200,
            data: function (params) {
                return {
                    path: params.term
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
        let _this = this,
            postPath = $(_this).val(),
            notify = $.notify("Loading local files...", {
                type: 'success',
                allow_dismiss: false,
                showProgressbar: true,
                delay: 0
            });
        notify.update('progress', '35');
        $.post("/localPathGoTo", {
            path: postPath
        }, function (response) {
            if (response["result"]) {
                let htmlBlockLocal = "",
                    localFiles = response["local_files"];
                localPathInputSearch = postPath;
                localPathSeparator = response["local_path_separator"];
                if (localPathHistoryBack[0] !== postPath) {
                    localPathHistoryBack.unshift(postPath);
                }
                $.each(localFiles, function (i, val) {
                    let path = val.Path,
                        dir = val.Directory,
                        symlink = val.Symlink;
                    htmlBlockLocal += '<div><span data-dir="' + dir + '" data-symlink="' + symlink + '" data-name="' + path + '">' + path + '</span></div>';
                });
                $(".localFilesBlock").html(htmlBlockLocal).find("div").each(function () {
                    if ($(this).outerWidth() < $(this).find("span").outerWidth()) {
                        $(this).css("border-right-color", "initial").css("border-right-style", "initial").css("border-right-width", "initial");
                    }
                });
            } else {
                $(_this).select2("trigger", "select", {
                    data: {
                        id: localPathInputSearch,
                        text: localPathInputSearch
                    }
                });
                sendNotify(response["message"], "danger");
            }
        }, 'json').always(function () {
            notify.update('progress', '100');
            notify.close();
        });
    });

    $(document).on("click", "div.localFilesBlock div > span[data-dir='true']", function () {
        const localPathInput = $('#localPath');
        let dirName = $(this).attr("data-name"),
            newPathName = "",
            separator = localPathSeparator;

        if (localPathInput.val() === localPathSeparator) {
            separator = "";
        }

        if ($(this).attr("data-symlink") === "true") {
            let dirNameSplit = dirName.split("->");
            dirName = $.trim(dirNameSplit[dirNameSplit.length - 1])
        }
        if (dirName.charAt(0) === localPathSeparator) {
            newPathName = dirName;
        } else {
            newPathName = localPathInput.val() + separator + dirName;
        }
        localPathInput.select2("trigger", "select", {
            data: {
                id: newPathName,
                text: newPathName
            }
        });
    });

    $(".localGoBack").on("click", function () {
        if (localPathHistoryBack.length < 2) {
            return;
        }
        const localPathInput = $('#localPath');
        localPathInput.select2("trigger", "select", {
            data: {
                id: localPathHistoryBack[1],
                text: localPathHistoryBack[1]
            }
        });
        localPathHistoryBack.splice(0, 2);
    });
});