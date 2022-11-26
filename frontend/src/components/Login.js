import React, { useRef, useState, useEffect } from 'react';

import { Link, useNavigate, useLocation } from 'react-router-dom';
import axios from '../api/axios';
import useAuth from '../context/useAuth';
import '../styles/login.css';

export default function Login () {
  const { setAuth } = useAuth();

  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || '/';

  const errRef = useRef();

  const [user, setUser] = useState('');
  const [pwd, setPwd] = useState('');
  const [errMsg, setErrMsg] = useState('');

  useEffect(() => {
    setErrMsg('');
  }, [user, pwd]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const res = await axios({
        method: 'post',
        url: '/auth/login',
        withCredentials: true,
        headers: { 'Content-Type': 'application/json' },
        data: JSON.stringify({ username: user, password: pwd })
      });

      const accessToken = res.data?.accessToken;
      const role = res.data?.role;

      console.log(accessToken);

      setUser('');
      setPwd('');
      setAuth({ role, accessToken });
      navigate(from);
    } catch (err) {
      if (!err.response) {
        setErrMsg('No Server Response');
      } else if (err.response.status === 400) {
        setErrMsg('Username or password missing');
      } else {
        setErrMsg('Registration Failed');
      }
      errRef.current.focus();
    }
  };

  return (
    <div className='login-container'>
      <form className='login-form'>
        <p ref={errRef} className={errMsg ? 'errmsg' : 'hide'}>{errMsg}</p>
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

        <button className="login-button" onClick={handleSubmit}>Login</button>
        <p className="not-signed-in">
                Not registered?<br />
          <span className="line">
            <Link to="/register" className="sign-in">Register now!</Link>
          </span>
        </p>
      </form>
    </div>
  );
}
