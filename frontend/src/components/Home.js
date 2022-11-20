import React from 'react';
import '../styles/home.css';

function Home () {
  return (
    <div className="home">
      <div className="home-container">
        <div className='text1'>
          <img
            src="home.png"
            className="home-image"
            width={'400px'}
            height={'400px'}
          >
          </img>
          <div className='text-container'>
            <h3>Had a new adventure?</h3>
            <p>Its time to store your wonderful pictures to later remember!</p>
            <p></p>
          </div>
        </div>
        <div className='text1'>
          <div className='text-container'>
            <h3>My Journey</h3>
            <p>A map of all the best memories for all over the world</p>
            <p>Choose your favorite ones, and visualize them over location or date</p>
          </div>
          <img
            src="home.png"
            className="home-image"
            width={'400px'}
            height={'400px'}
          >
          </img>
        </div>
      </div>
    </div>
  );
}

export default Home;
