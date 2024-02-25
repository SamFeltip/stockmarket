//@ts-check

/** @type {HTMLFormElement?} */
signupForm = document.querySelector("form#signup")

if (signupForm == null) {
    throw new Error("form#signup element is not present")
}

signupForm.addEventListener("submit", (event) => {
    event.preventDefault()

    console.log(event);
    console.log(event.target);

    /** @type {HTMLInputElement?} */
    // @ts-ignore
    const target = event.target

    if (target == null) {
        throw new Error("input.input-stock event listener failed ")
    }

    /** @type {HTMLInputElement?} */
    let form_name_elem = target.querySelector("input#name")

    /** @type {HTMLInputElement?} */
    let form_password_elem = target.querySelector("input#password")

    /** @type {HTMLSelectElement?} */
    let form_profile_elem = target.querySelector("select#profile")
    
    if (form_name_elem == null || form_password_elem == null || form_profile_elem == null) {
        throw new Error("form is missing elements")
    }

    let body = {
        Name: form_name_elem.value,
        Password: form_password_elem.value,
        Profile: form_profile_elem.value
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
