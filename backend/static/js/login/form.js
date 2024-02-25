let signupForm = document.querySelector("form#login")
signupForm.addEventListener("submit", function (event) {
    event.preventDefault()

    let body = {
        Name: signupForm.querySelector("#name").value,
        Password: signupForm.querySelector("#password").value,
    }

    console.log(body)

    // set post to /signup with form body. if response.json contains error, show in alert. else redirect to "/"
    fetch("/login", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body)
    })
        .then(response => response.json())
        .then(data => {
            if (data.error != null) {
                alert(data.error)
                return
            }

            window.location.href = "/"
        })
})
