import React from "react"

import "../../../sass/gallery.scss";

import AlbumCard from "./albumCard";

function AlbumsPage(props){
    return(
        <div>
            <nav className="breadcrumb is-large">
                <ul>
                    <li className="is-active" key={"g-"+props.gallery.id}><a>{props.gallery.title}</a></li>
                </ul>
            </nav>
            <div className="columns is-multiline">
                {props.albums.map(a=>{
                    return (
                        <div className="column is-one-third" key={"a-"+a.id}>
                            <AlbumCard gallery={props.gallery}
                                       album={a}
                                       onClick={()=>props.loadImages(a.id)} />
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

export default AlbumsPage;