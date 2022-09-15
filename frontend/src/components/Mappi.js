import "mapbox-gl/dist/mapbox-gl.css";
import "react-map-gl-geocoder/dist/mapbox-gl-geocoder.css";
import React, { useEffect, useRef, useState } from "react";
import Map from "react-map-gl";
import "../App.css"

const token = process.env.REACT_APP_MAPBOX_TOKEN;

const Mappi = () => {
  const [viewport, setViewPort] = useState({
    latitude: 0,
    longitude: 0,
    zoom: 1,
    transitionDuration: 100,
  });
  const mapRef = useRef();
  useEffect(() => {
    console.log({ viewport });
  }, [viewport]);

  return (
    <div className="map-container">
      <h1>Use the search bar to find a location on the map</h1>
      <Map
        ref={mapRef}
        {...viewport}
        mapStyle="mapbox://styles/mapbox/streets-v9"
        width="100%"
        height="70vh"
        onViewportChange={setViewPort}
        mapboxAccessToken={token}
      >
      </Map>
    </div>
  );
};
export default Mappi;