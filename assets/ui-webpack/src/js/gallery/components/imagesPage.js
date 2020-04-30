import React from "react"

import "../../../sass/gallery.scss";

import ImageCard from "./imageCard";

function ImagesPage(props) {
    return (
        <div>
            <nav className="breadcrumb is-large">
                <ul>
                    <li>
                        <a href="#!/" onClick={()=>{props.loadAlbums();}}>{props.gallery.title}</a>
                    </li>
                    <li className="is-active"><a>{props.album.title}</a></li>
                </ul>
            </nav>
            <div className="columns is-multiline">
                {props.images.map(i => {
                    return (
                        <div className="column is-one-third">
                            <ImageCard gallery={props.gallery} album={props.album} image={i}/>
                        </div>
                    );
                })}
            </div>
        </div>
    )
}

export default ImagesPage;
