/* eslint-disable no-unused-vars */
import React, { useEffect, useState } from 'react';
import axiosPrivate from '../api/axios';
import '../styles/my-images.css';

export default function MyImages () {
  const [images, setImages] = useState([]);

  useEffect(() => {
    const getImages = async () => {
      try {
        const res = await axiosPrivate({
          method: 'get',
          url: '/images'
        });

        setImages(res.data);
      } catch (err) {
        // error handler
      }
    };
    getImages();
  }, []);

  /*   const allImages = [
    { id: 1, imageName: 'img.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '250px', width: '200px' },
    { id: 2, imageName: 'hori.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '150px', width: '300px' },
    { id: 3, imageName: 'img.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '250px', width: '200px' },
    { id: 4, imageName: 'img.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '250px', width: '200px' },
    { id: 5, imageName: 'img.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '250px', width: '200px' },
    { id: 6, imageName: 'img.jpg', date: '13-08-2022', place: 'Quarteira, Algarve', height: '250px', width: '200px' }
  ]; */

  return (
    <div className='image-gallery'>
      {images.map(image => (
        <div key={image.id} className="image-card">
          <img className="image" height={image.height} width={image.width} src={image.imageName} alt="" />
          <div className='image-info'>
            <h3 className='image-date'>{image.date}</h3>
            <h3 className='image-place'>{image.place}</h3>
          </div>
        </div>
      ))}
    </div>
  );
}
