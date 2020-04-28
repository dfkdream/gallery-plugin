import React, {Component} from "react"
import {Link} from 'react-router-dom';

import "../../../sass/admin-card.scss"

const popups = require('../popup');

function CardImage(props) {
    if (props.album.cover !== 0) {
        return <div className="card-image img"
                    style={{
                        backgroundImage: `url(/api/gallery/${props.gallery.id}/album/${props.album.id}/image/${props.album.cover}?thumb=1)`
                    }}/>
    } else {
        return <div className="card-image"><div className="image-placeholder img-responsive"/></div>
    }
}

class AlbumCard extends Component {
    constructor(props) {
        super(props);
        this.deleteAlbum = this.deleteAlbum.bind(this);
        this.renameAlbum = this.renameAlbum.bind(this);
    }

    deleteAlbum() {
        popups.confirm("Delete Album", `Delete album ${this.props.album.title}?`)
            .then(ok => {
                if (ok) {
                    fetch("/api/gallery/" + this.props.gallery.id + "/album/" + this.props.album.id, {
                        method: "DELETE"
                    })
                        .then(
                            resp => {
                                if (!resp.ok) {
                                    popups.alert("Error", `${resp.status} ${resp.statusText}`);
                                } else {
                                    this.props.onChange();
                                }
                            },
                            error => {
                                popups.alert("Error", error.message);
                            });
                }
            })
    }

    renameAlbum() {
        popups.prompt("Rename Album", `Rename album ${this.props.album.title} to:`)
            .then(resp => {
                if (resp) {
                    fetch("/api/gallery/" + this.props.gallery.id + "/album/" + this.props.album.id, {
                        method: "POST",
                        body: JSON.stringify({title: resp})
                    })
                        .then(
                            resp => {
                                if (!resp.ok) {
                                    popups.alert("Error", `${resp.status} ${resp.statusText}`);
                                } else {
                                    this.props.onChange();
                                }
                            },
                            error => {
                                popups.alert("Error", error.message);
                            }
                        )
                }
            })
    }

    render() {
        return (
            <div className="column col-4 col-sm-12 card-container">
                <div className="card">
                    <CardImage gallery={this.props.gallery} album={this.props.album}/>
                    <div className="card-header">
                        <button className="btn btn-link float-right tooltip" data-tooltip="Delete"
                                onClick={this.deleteAlbum}><i className="fa fa-trash-alt"/></button>
                        <button className="btn btn-link float-right tooltip" data-tooltip="Rename"
                                onClick={this.renameAlbum}><i className="fa fa-edit"/></button>
                        <Link className="card-title h5" to={{
                            pathname: "/admin/gallery/" + this.props.gallery.id + "/album/" + this.props.album.id,
                            state: {
                                gallery: this.props.gallery,
                                album: this.props.album
                            }
                        }}>{this.props.album.title}</Link>
                    </div>
                </div>
            </div>
        )
    }
}

export default AlbumCard;
