import React, {Component} from "react";

import {Lightbox} from "react-modal-image";

import "../../../sass/gallery.scss";

class ImageCard extends Component {
    constructor(props) {
        super(props);
        this.imgSrc = `/api/gallery/${props.gallery.id}/album/${props.album.id}/image/${props.image.id}`;
        this.state = {
            isLightboxOpen: false
        }
    }

    render() {
        return (
            <div>
                <a className="image-card">
                    <figure className="image is-1by1 img"
                            data-description={this.props.image.description === "" ? null : this.props.image.description}
                            style={{backgroundImage: `url(${this.imgSrc}?thumb=1)`}}
                            onClick={() => {this.setState({isLightboxOpen: true})}}
                    />
                </a>
                {this.state.isLightboxOpen && (
                    <Lightbox small={this.imgSrc + "?thumb=1"}
                              large={this.imgSrc}
                              alt={this.props.image.description}
                              onClose={()=>{this.setState({isLightboxOpen: false})}}
                    />
                )}
            </div>
        );
    }
}

export default ImageCard;