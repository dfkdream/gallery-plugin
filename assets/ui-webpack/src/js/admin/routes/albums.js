import React, {Component} from "react"
import AlbumsEmpty from "../components/albumsEmpty";

import "../../../sass/admin-location-header.scss";
import AlbumCard from "../components/albumCard";

const popups = require('../popup');

class AlbumsHeader extends Component {
    render() {
        return (
            <div className="location-header">
                <h1>{this.props.gallery.title}</h1>
                <div>
                    <button className="btn btn-primary" onClick={this.props.onClick}>
                        <i className="fa fa-plus"/>
                        Add Album
                    </button>
                </div>
            </div>
        )
    }
}

class Albums extends Component {
    constructor(props) {
        super(props);

        this.createAlbum = this.createAlbum.bind(this);

        this.state = {
            gallery: null,
            isGalleryLoaded: false,
            isLoaded: false,
            error: null,
            albums: []
        };

        this.gid = this.props.match.params.gid;
    }

    createAlbum() {
        popups.prompt("Create Album", "Enter album title:")
            .then(resp => {
                if (resp) {
                    fetch("/api/gallery/" + this.gid + "/albums", {
                        method: "POST",
                        body: JSON.stringify({title: resp})
                    })
                        .then((resp) => {
                            if (!resp.ok) {
                                popups.alert("Error", `${resp.status} ${resp.statusText}`);
                                return
                            }

                            this.loadAlbums();
                        }, (error) => {
                            popups.alert("Error", error.message);
                        })
                }
            })
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

    loadAlbums() {
        fetch("/api/gallery/" + this.gid + "/albums")
            .then(res => res.json())
            .then(
                (json) => {
                    this.setState({isLoaded: true, albums: json})
                },
                (error) => {
                    this.setState({isLoaded: true, error: error})
                }
            )
    }

    componentDidMount() {
        this.loadAlbums();
        if (!this.state.isGalleryLoaded) {
            this.loadGallery();
        }
    }

    render() {
        const {gallery, isGalleryLoaded, isLoaded, error, albums} = this.state;

        if (error) {
            return (
                <div>
                    <h1>Error</h1>
                    <pre>{error.message}</pre>
                </div>
            )
        } else if (!(isLoaded && isGalleryLoaded)) {
            return <h1>Loading...</h1>
        } else {
            if (albums.length > 0) {
                return (
                    <div>
                        <AlbumsHeader gallery={gallery} onClick={this.createAlbum}/>
                        <div className="columns">
                            {albums.map(a => {
                                return <AlbumCard gallery={gallery} album={a} onChange={() => {
                                    this.loadAlbums()
                                }}/>
                            })}
                        </div>
                    </div>
                )
            } else {
                return (
                    <div>
                        <AlbumsHeader gallery={gallery} onClick={this.createAlbum}/>
                        <AlbumsEmpty/>
                    </div>
                )
            }
        }
    }
}

export default Albums;
