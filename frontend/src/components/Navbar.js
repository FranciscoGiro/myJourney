import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from '../api/axios';
import useAuth from '../context/useAuth';
import '../styles/nav-bar.css';

function Navbar () {
  // eslint-disable-next-line no-unused-vars
  const { auth, setAuth } = useAuth();
  const [clicked, setClicked] = useState(false);

  const navigate = useNavigate();

  const handleMenuClick = () => {
    setClicked(!clicked);
  };

  const HandleLogout = async () => {
    try {
      await axios({
        method: 'get',
        url: '/auth/logout'
      });

      setAuth({ role: null, accessToken: null });
      navigate('/');
    } catch (err) {
      // ignore error
    }
  };

  return (
    <nav className="navbar">
      <div className='nav-info'>
        <Link className='nav-name' to='/home'>MyJourney</Link>
        <div className='menu-icon' onClick={handleMenuClick}>
          <i className={clicked ? 'fas fa-times' : 'fas fa-bars'}></i>
        </div>
      </div>
      <ul className={clicked ? 'nav-menu active' : 'nav-menu'}>
        <li><Link className='nav-link' to='/home'>Home</Link></li>
        <li><Link className='nav-link' to='/map'>Map</Link></li>
        <li><Link className='nav-link' to='/myimages'>My Images</Link></li>
        <li><Link className='nav-link' to='/upload'>Upload</Link></li>
        {console.log('hi')}
        {console.log(localStorage.getItem('Authorization'))}
        {auth.accessToken
          ? <li><button className='nav-link' onClick={HandleLogout}>Logout</button></li>
          : <li><Link className='nav-link' to='/login'>Login</Link></li>
        }
      </ul>
    </nav>
  );
}

export default Navbar;
