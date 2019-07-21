import React, { useState } from "react"
import { apiURl } from "../api"

const Login = ({ hisory }) => {
    const [state, setState] = useState({email:"", password:""})
    const { email, password } = state
    
    const handleChange = (e) => {
        const {name, value} = e.target
        setState({...state, [name] :value })
        console.log(email, password)
    }

    const handleSubmit = async (e) =>{
        const {email, password} = state

        const res =await fetch(`${apiURl}/login`, {
            method: "POST",
            body: JSON.stringify({
                email, 
                password
            }),
            headers: {
                "Content-Type": "application/json",
            }
        }).then((res) => res.json())
        const { token } = res
        
        localStorage.setItem("auth",`Bearer ${token}`)

        // history.push('/session')
    }

    return(
        <div>
            <input type="text" placeholder="email" value={email} name="email" onChange={(e)=>{handleChange(e)}}/>
            <input type="password" placeholder="password" value={password} name="password" onChange={(e)=>{handleChange(e)}}/>
            <button onClick={()=>handleSubmit()} >submit</button>
        </div>
    );
}

export default Login