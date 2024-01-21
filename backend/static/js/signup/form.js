let signup = document.querySelector("form#signup")
signup.addEventListener("submit", function (event) {
    event.preventDefault()

    let body = {
        Name: signup.querySelector("#name").value,
        Password: signup.querySelector("#password").value,
        Profile: signup.querySelector("#profile").value
    }

    console.log(body)

    // set post to /signup with form body. if response.json contains error, show in alert. else redirect to "/"
    fetch("/signup", {
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
