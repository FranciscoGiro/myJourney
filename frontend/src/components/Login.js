import React from 'react';
import { useRef, useState, useEffect, useContext } from 'react';
import "../styles/login.css"

export default function Login() {

  const [user, setUser] = useState('');
  const [pwd, setPwd] = useState('');
  const [errMsg, setErrMsg] = useState('');

  useEffect(() => {
    setErrMsg('');
}, [user, pwd])
  
  return(
    <div className='login-container'>
      <form className='login-form'>
        <h3 className='login-title'>Login Here</h3>
        <div className='form-entry'>
          <label className='form-label'>
            Username
          </label>
          <input
            type="text"
            id="username"
            className='form-input'
            onChange={(e) => setUser(e.target.value)}
            value={user}
            required
          ></input>
        </div>
        <div className='form-entry'>
          <label className='form-label'>
            Password
          </label>
          <input
            type="password"
            id="password"
            className='form-input'
            onChange={(e) => setPwd(e.target.value)}
            value={pwd}
            required
          ></input>
        </div>
        <button className="login-button" type="submit">Login</button>
      </form>
    </div>
  )
}