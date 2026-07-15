import { useState } from "react";
import api from "../api/axios";
import { useNavigate } from "react-router-dom";


function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();

  const handleLogin=async () =>{
    console.log("Login button clicked");
    try{
        const response=await api.post("/login",{
            email,
            password,
        });
    localStorage.setItem("token", response.data.token);
    navigate("/dashboard");
    console.log("Token saved:", response.data.token);
    
    }catch(error){
        console.error(error);
    }
    };  

  return (
    <div>
      <h1>PlayLedger Login</h1>

      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />

      <br />
      <br />

      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />

      <br />
      <br />

      <button onClick={handleLogin}>Login</button>
    </div>
  );
}

export default Login;