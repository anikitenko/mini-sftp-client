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

$(document).ajaxError(function (event, jqxhr) {
    if (jqxhr.status === 500 || jqxhr.status === 502) {
        sendNotify("Something went wrong!", "danger");
    } else if (jqxhr.status === 404) {
        sendNotify("Requested resource not found!", "danger");
    } else if (jqxhr.status === 403) {
        sendNotify("You don't have permission to resource you are trying to access!", "danger");
    }
});

$.fn.contextMenu = function (settings) {

    return this.each(function () {

        // Open context menu
        $(this).on("contextmenu", function (e) {
            // return native menu if pressing control
            if (e.ctrlKey) return;

            let $eTarget = $(e.target);

            if (e.target.nodeName === "DIV") {
                $eTarget = $(e.target).find("span")
            }

            //open menu
            let $menu = $(settings.menuSelector)
                .data("invokedOn", $eTarget)
                .show()
                .css({
                    position: "absolute",
                    left: getMenuPosition(e.clientX, 'width', 'scrollLeft'),
                    top: getMenuPosition(e.clientY, 'height', 'scrollTop')
                })
                .off("click")
                .on('click', 'a', function (e) {
                    $menu.hide();

                    let $invokedOn = $menu.data("invokedOn");

                    settings.menuSelected.call(this, $invokedOn, $(e.target));
                });
            $menu.find(".invokedOn").text(truncate(50, $eTarget.text()));

            if ($eTarget.attr("data-dir") === "true") {
                $menu.find("[data-action='showFileContent']").hide();
            } else {
                $menu.find("[data-action='showFileContent']").show();
            }

            return false;
        });

        //make sure menu closes on any click
        $(document).on("click", function () {
            $(settings.menuSelector).hide();
        });
    });

    function getMenuPosition(mouse, direction, scrollDir) {
        let win = $(window)[direction](),
            scroll = $(window)[scrollDir](),
            menu = $(settings.menuSelector)[direction](),
            position = mouse + scroll;

        // opening menu would pass the side of the page
        if (mouse + menu > win && menu < mouse)
            position -= menu;

        return position;
    }

};

function enableContextMenuLocal(selector) {
    selector.contextMenu({
        menuSelector: "#contextMenu",
        menuSelected: function (invokedOn, selectedMenu) {
            if (selectedMenu.hasClass("invokedOn")) {
                return
            }
            switch (selectedMenu.attr("data-action")) {
                case "showFileContent": {
                    const fileContentsModal = $("#fileContentsModal");
                    $.post("/showFileContent", {
                        path: $('#localPath').val(),
                        name: invokedOn.text()
                    }, function (response) {
                        fileContentsModal.find(".modal-title").text(invokedOn.text());
                        if (response["result"]) {
                            let fileContent = escapeHtml(response["contents"]);
                            let htmlBlock = '<pre class="pre-scrollable" style="max-height:430px"><code>' +
                                fileContent + "</code></pre>";
                            fileContentsModal.find(".modal-body").html(htmlBlock);
                        } else {
                            fileContentsModal.find(".modal-body").html('<span class="text-danger">' + response["message"] + '</span>');
                        }
                    }).always(function () {
                        fileContentsModal.modal("show");
                    });
                    break;
                }
                case "deleteItem":
                    bootbox.confirm(
                        "Are you sure you want to remove <strong>" +
                        invokedOn.text() +
                        "</strong> PERMANENTLY?", function (result) {
                            if (result) {
                                $.post("/removeLocalItem", {path: $('#localPath').val(), name: invokedOn.text()}, function (response) {
                                    if (response["result"]) {
                                        invokedOn.parent().detach();
                                        sendNotify(invokedOn.text()+" was successfully removed!", "success");
                                    } else {
                                        sendNotify(response["message"], "danger");
                                    }
                                }, 'json')
                            }
                        });
                    break;
            }
        }
    });
}

function truncate(length, string){
    if (string.length > length)
        return string.substring(0,length)+'...';
    else
        return string;
}

function escapeHtml(string) {
    let entityMap = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;',
        '/': '&#x2F;',
        '`': '&#x60;',
        '=': '&#x3D;'
    };
    return String(string).replace(/[&<>"'`=/]/g, function (s) {
        return entityMap[s];
    });
}

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