import React, {Component} from "react"
import ImagesEmpty from "../components/imagesEmpty";
import ImageCard from "../components/imageCard";

import "../../../sass/admin-location-header.scss";

const popups = require('../popup');

class ImagesHeader extends Component {
    render() {
        return (
            <div className="location-header">
                <h1>{this.props.gallery.title} > {this.props.album.title}</h1>
                <div>
                    <button className="btn btn-primary" onClick={this.props.onClick}>
                        <i className="fa fa-file-upload"/>
                        Upload Images
                    </button>
                </div>
            </div>
        )
    }
}

class Images extends Component {
    constructor(props) {
        super(props);

        this.uploadImages = this.uploadImages.bind(this);

        this.state = {
            gallery: null,
            isGalleryLoaded: false,
            album: null,
            isAlbumLoaded: false,
            isLoaded: false,
            error: null,
            images: []
        };

        this.gid = this.props.match.params.gid;
        this.aid = this.props.match.params.aid;
    }

    loadGallery() {
        if (this.props.location.state) {
            this.setState({
                gallery: this.props.location.state.gallery,
                isGalleryLoaded: true
            });
            return
        }

        fetch("/api/gallery/" + this.gid)
            .then(res => res.json())
            .then(
                (json) => {
                    this.setState({
                        gallery: json,
                        isGalleryLoaded: true,
                    });
                },
                (error) => {
                    popups.alert("Error", error.message);
                }
            )
    }

    loadAlbum() {
        if (this.props.location.state) {
            this.setState({
                album: this.props.location.state.album,
                isAlbumLoaded: true
            });
            return
        }

        fetch("/api/gallery/" + this.gid + "/album/" + this.aid)
            .then(res => res.json())
            .then(
                (json) => {
                    this.setState({
                        album: json,
                        isAlbumLoaded: true,
                    });
                },
                (error) => {
                    popups.alert("Error", error.message);
                }
            )
    }

    loadImages() {
        fetch("/api/gallery/" + this.gid + "/album/" + this.aid + "/images")
            .then(res => res.json())
            .then(
                (json) => {
                    this.setState({isLoaded: true, images: json})
                },
                (error) => {
                    this.setState({isLoaded: true, error: error})
                }
            )
    }

    uploadImages(){
        popups.upload(`/api/gallery/${this.gid}/album/${this.aid}/images`)
            .then(()=>{
                this.loadImages();
            });
    }

    componentDidMount() {
        this.loadGallery();
        this.loadAlbum();
        this.loadImages();
    }

    render() {
        const {gallery, isGalleryLoaded, album, isAlbumLoaded, isLoaded, error, images} = this.state;

        if (error) {
            return (
                <div>
                    <h1>Error</h1>
                    <pre>{error.message}</pre>
                </div>
            )
        } else if (!(isLoaded && isGalleryLoaded && isAlbumLoaded)) {
            return <h1>Loading...</h1>
        } else {
            if (images.length > 0) {
                return (
                    <div>
                        <ImagesHeader gallery={gallery} album={album} onClick={this.uploadImages}/>
                        <div className="columns">
                            {images.map(i => {
                                return <ImageCard gallery={gallery} album={album} image={i} onChange={() => {
                                    this.loadImages();
                                }}/>
                            })}
                        </div>
                    </div>
                )
            } else {
                return (
                    <div>
                        <ImagesHeader gallery={gallery} album={album} onClick={this.uploadImages}/>
                        <ImagesEmpty/>
                    </div>
                )
            }
        }
    }
}

export default Images;
