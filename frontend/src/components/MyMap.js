/* eslint-disable no-unused-vars */
import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import useAxiosPrivate from '../api/useAxiosPrivate';
import '../styles/my-map.css';
import 'mapbox-gl/dist/mapbox-gl.css';
import { useNavigate } from 'react-router';

mapboxgl.accessToken = process.env.REACT_APP_MAPBOX_TOKEN;

export default function MyMap () {
  const mapContainer = useRef(null);
  const map = useRef(null);
  const [lng, setLng] = useState(41.3);
  const [lat, setLat] = useState(12.57);
  const [zoom, setZoom] = useState(1);

  const axiosPrivate = useAxiosPrivate();
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    if (map.current) return;
    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v9',
      center: [lng, lat],
      zoom
    });
    const getImages = async () => {
      try {
        const res = await axiosPrivate({
          method: 'get',
          url: '/images'
        });

        setImages(res.data);
        res.data.map(m => {
          const popup = new mapboxgl.Popup({ offset: 25 }).setText(
            'Construction on the Washington Monument began in 1848.');

          const el = document.createElement('div');
          el.className = 'marker';
          el.style.backgroundImage = `url(${m.url})`;
          el.style.width = '50px';
          el.style.height = '50px';
          el.style.backgroundSize = '100%';

          return new mapboxgl.Marker(el)
            .setLngLat([m.lng, m.lat])
            .addTo(map.current);
        });
        setLoading(false);
      } catch (err) {
        if (err.status === 401) {
          navigate('/login');
        }
      }
    };
    getImages();
  });

  useEffect(() => {
    if (!map.current) return;
    map.current.on('move', () => {
      setLng(map.current.getCenter().lng.toFixed(4));
      setLat(map.current.getCenter().lat.toFixed(4));
      setZoom(map.current.getZoom().toFixed(2));
    });
  });

  return (
    <div className='map-background'>
      <div ref={mapContainer} className="map-container" />
    </div>
  );
}
