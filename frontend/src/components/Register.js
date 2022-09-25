import React, { useRef, useState, useEffect } from "react";
import { faCheck, faTimes, faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { register } from "../services/RemoteServices"
import "../styles/register.css"


const EMAIL_REGEX = /\S+@\S+\.\S+/;

const Register = () => {
    const userRef = useRef();
    const errRef = useRef();

    const [user, setUser] = useState('');
    const [validName, setValidName] = useState(false);
    const [userFocus, setUserFocus] = useState(false);

    const [email, setEmail] = useState('');
    const [validEmail, setValidEmail] = useState(false);
    const [emailFocus, setEmailFocus] = useState(false);

    const [pwd, setPwd] = useState('');
    const [validPwd, setValidPwd] = useState(false);
    const [pwdFocus, setPwdFocus] = useState(false);

    const [matchPwd, setMatchPwd] = useState('');
    const [validMatch, setValidMatch] = useState(false);
    const [matchFocus, setMatchFocus] = useState(false);

    const [errMsg, setErrMsg] = useState('');
    const [success, setSuccess] = useState(false);

    useEffect(() => {
        userRef.current.focus();
    }, [])

    useEffect(() => {
        setValidName(user.length >= 5);
    }, [user])

    useEffect(() => {
        setValidEmail(EMAIL_REGEX.test(email));
    }, [email])

    useEffect(() => {
        setValidPwd(pwd.length >= 8);
        setValidMatch(pwd === matchPwd);
    }, [pwd, matchPwd])

    useEffect(() => {
        setErrMsg('');
    }, [user, email, pwd, matchPwd])

    const handleSubmit = async (e) => {
        e.preventDefault();
        // if button enabled with JS hack
        const v1 = user.length >= 5;
        const v2 = pwd.length >= 8;
        if (!v1 || !v2) {
            setErrMsg("Invalid Entry");
            return;
        }
        try {
            const response = await register(user, email, pwd)
            // TODO: remove console.logs before deployment
            if(response.data){
                console.log(JSON.stringify(response.data));
            }
            setSuccess(true);
            //clear state and controlled inputs
            setUser('');
            setPwd('');
            setMatchPwd('');
        } catch (err) {
            if (!err.response) {
                setErrMsg('No Server Response');
            } else if (err.response.status === 409) {
                setErrMsg('Username Taken');
            } else {
                setErrMsg('Registration Failed')
            }
            errRef.current.focus();
        }
    }

    return (
        <section className="form-section">
            <p ref={errRef} className={errMsg ? "errmsg" : "hide"} aria-live="assertive">{errMsg}</p>
            <h1 className="title">Register</h1>
            <form onSubmit={handleSubmit} className='form'>
                <label htmlFor="username" className="input-name">
                    Username:
                    <FontAwesomeIcon icon={faCheck} className={validName ? "valid" : "hide"} />
                    <FontAwesomeIcon icon={faTimes} className={validName || !user ? "hide" : "invalid"} />
                </label>
                <input
                    type="text"
                    id="username"
                    ref={userRef}
                    autoComplete="off"
                    onChange={(e) => setUser(e.target.value)}
                    value={user}
                    required
                    aria-invalid={validName ? "false" : "true"}
                    aria-describedby="uidnote"
                    onFocus={() => setUserFocus(true)}
                    onBlur={() => setUserFocus(false)}
                />
                <p id="uidnote" className={userFocus && user && !validName ? "instructions" : "offscreen"}>
                    <FontAwesomeIcon icon={faInfoCircle} />
                    4 to 24 characters.<br />
                    Must begin with a letter.<br />
                    Letters, numbers, underscores, hyphens allowed.
                </p>

                <label htmlFor="email" className="input-name">
                    Email:
                    <FontAwesomeIcon icon={faCheck} className={validEmail ? "valid" : "hide"} />
                    <FontAwesomeIcon icon={faTimes} className={validEmail || !email ? "hide" : "invalid"} />
                </label>
                <input
                    type="email"
                    id="email"
                    onChange={(e) => setEmail(e.target.value)}
                    value={email}
                    required
                    aria-invalid={validEmail ? "false" : "true"}
                    aria-describedby="emailnote"
                    onFocus={() => setEmailFocus(true)}
                    onBlur={() => setEmailFocus(false)}
                />
                <p id="emailnote" className={emailFocus && email ? "instructions" : "offscreen"}>
                    <FontAwesomeIcon icon={faInfoCircle} />
                    Insert valid email account<br />
                </p>


                <label htmlFor="password" className="input-name">
                    Password:
                    <FontAwesomeIcon icon={faCheck} className={validPwd ? "valid" : "hide"} />
                    <FontAwesomeIcon icon={faTimes} className={validPwd || !pwd ? "hide" : "invalid"} />
                </label>
                <input
                    type="password"
                    id="password"
                    onChange={(e) => setPwd(e.target.value)}
                    value={pwd}
                    required
                    aria-invalid={validPwd ? "false" : "true"}
                    aria-describedby="pwdnote"
                    onFocus={() => setPwdFocus(true)}
                    onBlur={() => setPwdFocus(false)}
                />
                <p id="pwdnote" className={pwdFocus && !validPwd ? "instructions" : "offscreen"}>
                    <FontAwesomeIcon icon={faInfoCircle} />
                    8 to 24 characters.<br />
                    Must include uppercase and lowercase letters, a number and a special character.<br />
                    Allowed special characters: <span aria-label="exclamation mark">!</span> <span aria-label="at symbol">@</span> <span aria-label="hashtag">#</span> <span aria-label="dollar sign">$</span> <span aria-label="percent">%</span>
                </p>


                <label htmlFor="confirm_pwd" className="input-name">
                    Confirm Password:
                    <FontAwesomeIcon icon={faCheck} className={validMatch && matchPwd ? "valid" : "hide"} />
                    <FontAwesomeIcon icon={faTimes} className={validMatch || !matchPwd ? "hide" : "invalid"} />
                </label>
                <input
                    type="password"
                    id="confirm_pwd"
                    onChange={(e) => setMatchPwd(e.target.value)}
                    value={matchPwd}
                    required
                    aria-invalid={validMatch ? "false" : "true"}
                    aria-describedby="confirmnote"
                    onFocus={() => setMatchFocus(true)}
                    onBlur={() => setMatchFocus(false)}
                />
                <p id="confirmnote" className={matchFocus && !validMatch ? "instructions" : "offscreen"}>
                    <FontAwesomeIcon icon={faInfoCircle} />
                    Must match the first password input field.
                </p>

                <button disabled={!validName || !validEmail || !validPwd || !validMatch ? true : false}>Sign Up</button>
            </form>
            <p className="signed-in">
                Already registered?<br />
                <span className="line">
                    {/*put router link here*/}
                    <a href="#" className="sign-in">Sign In</a>
                </span>
            </p>
        </section>
    )
}

export default Register