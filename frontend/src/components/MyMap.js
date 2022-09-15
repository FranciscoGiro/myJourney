import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import '../App.css';
import "mapbox-gl/dist/mapbox-gl.css";
 
mapboxgl.accessToken = process.env.REACT_APP_MAPBOX_TOKEN;


export default function MyMap() {

  
    const mapContainer = useRef(null);
    const map = useRef(null);
    const [lng, setLng] = useState(-70.9);
    const [lat, setLat] = useState(42.35);
    const [zoom, setZoom] = useState(3);
    
    useEffect(() => {
      if (map.current) return;
      map.current = new mapboxgl.Map({
          container: mapContainer.current,
          style: 'mapbox://styles/mapbox/streets-v9',
          center: [lng, lat],
          zoom: zoom
      });

      const marker = new mapboxgl.Marker()
        .setLngLat([30.5, 50.5])
        .addTo(map.current);
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
    <div ref={mapContainer} className="map-container" />
    );
}
