// Prompt function is JS module for allerts, notifications and popup dialogs
function Prompt() {
    let toast = function (c) {
        const {
            title = "",
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            title = "",
            msg = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer,
            confirmButtonColor: "#163b65"
        })
    }

    let error = function (c) {
        const {
            title = "",
            msg = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer,
            confirmButtonColor: "#163b65"
        })
    }

    async function custom(c) {
        const {
            html = "",
            title = "",
            footer = "",
            icon = "",
            showCancelButton = true,
            showConfirmButton = true,
        } = c;
    
        const { value: result } = await Swal.fire({
            html: html,
            title: title,
            footer: footer,
            icon: icon,
            backdrop: false,
            focusConfirm: true,
            showCancelButton: showCancelButton,
            showConfirmButton: showConfirmButton,
            confirmButtonColor: "#163b65",
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start").value,
                    document.getElementById("end").value
                ]
            },
        })
    
        if (result) {
            if (result.dismiss != Swal.DismissReason.cancel) {
                if (result[0] !== "") {
                    console.log("here");
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

