const f = document.querySelector("form")
        
        f.addEventListener("submit", (e)=>{
            e.preventDefault()
            console.log(e.target.email.value)
            fetch("/db", 
                {method:"POST", 
                headers: {
                    "Content-Type" : "application/json",
                },
                body: JSON.stringify(
                    {
                    ok: "asdasdasda",
                    email: e.target.email.value, 
                    password: e.target.password.value
                })})
        })