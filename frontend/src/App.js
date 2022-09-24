import React, { useRef, useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Home from './components/Home'
import MyMap from './components/MyMap'
import Navbar from './components/Navbar'
import Login from './components/Login'
import Register from './components/Register'
import UploadImage from './components/UploadImage'
import MyImages from './components/MyImages'
import NotFound from './components/NotFound'
 
 
export default function App() {
    
    return (
      <div className="App">
        <Router>
          <Navbar />
          <Routes>
              <Route path='/' element={<Home />}/>
              <Route path='/home' element={<Home /> }/>
              <Route path='/myimages' element={<MyImages /> }/>
              <Route path='/map' element={<MyMap /> }/>
              <Route path='/upload' element={<UploadImage /> }/>
              <Route path='/login' element={<Login /> }/>
              <Route path='/register' element={<Register /> }/>
              <Route path='*' element={<NotFound />}/>
          </Routes>
        </Router>
      </div>
    );
}
