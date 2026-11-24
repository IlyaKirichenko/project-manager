const form = document.getElementById("loginForm");

form.addEventListener("submit", async function (event) {
    event.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const response = await fetch("/login", {
        method: "POST",
        headers:{
            "Content-type": "application/json"
        },
        body: JSON.stringify({
            email: email,
            password: password
        })
    })

    const errorBox = document.getElementById("errorMessage");

    if (response.ok){
        errorBox.textContent = "";
        window.location.href = "/static/Projects/projects.html"
    }else{
        const data = await response.json();
        errorBox.textContent = data.error;
        
    }
})