import React, { useEffect, useState } from 'react';
import {getAllImages} from "../services/RemoteServices"
import "../styles/my-images.css"

export default function MyImages() {

    const [filters, setFilters] = useState({})

    const [images, setImages] = useState([])

    useEffect(() => {
        let data = getAllImages(filters)
        setImages(data)
    }, [filters]);


    //TODO delete this
    const allImages = [
        {id:1, imageName: "img.jpg", date: "18-02-2022", place: "Quarteira, Algarve", height:"250px", width:"200px"},
        {id:2, imageName: "hori.jpg", date: "18-02-2022", place: "Quarteira, Algarve", height:"150px", width:"300px"},
        {id:3, imageName: "img.jpg", date: "18-02-2022", place: "Quarteira, Algarve", height:"250px", width:"200px"},
        {id:4, imageName: "img.jpg", date: "18-02-2022", place: "Quarteira, Algarve", height:"250px", width:"200px"},
        {id:5, imageName: "img.jpg", date: "18-02-2022", place: "Quarteira, Algarve", height:"250px", width:"200px"}
    ]


    return (
        <div className='image-gallery'>
            {allImages.map(image => (
                <div key={image.id} className="image-card">
                    <img className="image" height={image.height} width={image.width} src={image.imageName} alt="" />
                    <h3 className='image-date'>{image.date}</h3>
                    <h3 className='image-place'>{image.place}</h3>
                </div>
            ))}
        </div>
    );
}