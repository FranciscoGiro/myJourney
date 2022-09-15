import React, { useState } from 'react';
import {Link} from "react-router-dom";
import "../styles/Navbar.css"

function Navbar() {
  const [clicked, setClicked] = useState(false);

  const handleMenuClick = () => {
    setClicked(!clicked)
  };


  return (
    <nav className="navbar">
      <h1 className='nav-logo'>MyJourney<i className='fab fa-react'></i></h1>
      <div className='menu-icon' onClick={handleMenuClick}>
        <i className={clicked ? 'fas fa-times' : 'fas fa-bars'}></i>
      </div>
      <ul className={clicked ? 'nav-menu active' : 'nav-menu'}>
        <li><Link className='nav-link' to='/home'>Home</Link></li>
        <li><Link className='nav-link' to='/map'>Map</Link></li>
        <li><Link className='nav-link' to='/myimages'>My Images</Link></li>
        <li><Link className='nav-link' to='/upload'>Upload</Link></li>
      </ul>
    </nav>
  )
}

export default Navbar