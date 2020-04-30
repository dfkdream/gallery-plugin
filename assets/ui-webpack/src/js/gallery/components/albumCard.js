import React from "react";

import "../../../sass/gallery.scss";

function AlbumCard(props){
    return(
        <a className="album-card" href={"#!/"+props.album.id} onClick={props.onClick}>
            {props.album.cover!==0?
                <figure className="image is-1by1 img"
                        style={{backgroundImage: `url(/api/gallery/${props.gallery.id}/album/${props.album.id}/image/${props.album.cover}?thumb=1)`}}
                />:
                <figure className="image is-1by1 placeholder"/>
            }
            <h3>{props.album.title}</h3>
        </a>
    )
}

export default AlbumCard;
