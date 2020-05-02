import React, {Component} from "react";
import Lightbox from "react-image-lightbox";
import "react-image-lightbox/style.css";

class ImageCard extends Component {
    constructor(props){
        super(props);

        this.state={
            currentIndex:this.props.index,
            isLightboxOpen: false,
        };

        this.image=this.props.images[this.props.index]
    }

    toImgSrc(id) {
        return `/api/gallery/${this.props.gallery.id}/album/${this.props.album.id}/image/${id}`
    }

    getPrevIndex() {
        return (this.state.currentIndex + this.props.images.length - 1) % this.props.images.length;
    }

    getNextIndex() {
        return (this.state.currentIndex + 1) % this.props.images.length;
    }

    render() {
        return (
            <div>
                <a className="image-card">
                    <figure className="image is-1by1 img"
                            data-description={this.image.description === "" ? null : this.image.description}
                            style={{backgroundImage: `url(${this.toImgSrc(this.image.id)}?thumb=1)`}}
                            onClick={()=>{this.setState({isLightboxOpen: true, currentIndex: this.props.index})}}
                    />
                </a>
                {this.state.isLightboxOpen &&
                    <Lightbox
                        mainSrc={this.toImgSrc(this.props.images[this.state.currentIndex].id)}
                        mainSrcThumbnail={this.toImgSrc(this.props.images[this.state.currentIndex].id) + "?thumb=1"}
                        imageTitle={this.props.images[this.state.currentIndex].description}
                        prevSrc={this.toImgSrc(this.props.images[this.getPrevIndex()].id)}
                        prevSrcThumbnail={this.toImgSrc(this.props.images[this.getPrevIndex()].id) + "?thumb=1"}
                        onMovePrevRequest={() => {
                            this.setState({currentIndex: this.getPrevIndex()})
                        }}
                        nextSrc={this.toImgSrc(this.props.images[this.getNextIndex()].id)}
                        nextSrcThumbnail={this.toImgSrc(this.props.images[this.getNextIndex()].id) + "?thumb=1"}
                        onMoveNextRequest={() => {
                            this.setState({currentIndex: this.getNextIndex()})
                        }}
                        onCloseRequest={() => {
                            this.setState({isLightboxOpen: false, currentIndex: this.props.index})
                        }}
                    />
                }
            </div>
        );
    }
}

export default ImageCard;