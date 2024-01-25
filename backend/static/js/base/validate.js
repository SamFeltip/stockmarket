// GET /validate with fetch, returns json, if return current_user, set to local storage, if not, remove current_user from local storage

fetch("/validate")
    .then(response => response.json())
    .then(data => {
        if (data.current_user != null) {
            localStorage.setItem("current_user", JSON.stringify(data.current_user))
        } else {
            localStorage.removeItem("current_user")
        }
    })
