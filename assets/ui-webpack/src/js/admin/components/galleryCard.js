import React, {Component} from "react"
import {Link} from 'react-router-dom';

import "../../../sass/admin-card.scss"

const popups = require('../popup');

class GalleryCard extends Component {
    constructor(props) {
        super(props);
        this.deleteGallery = this.deleteGallery.bind(this);
        this.renameGallery = this.renameGallery.bind(this);
    }

    deleteGallery() {
        popups.confirm("Delete Gallery", `Delete gallery ${this.props.gallery.title}?`)
            .then(ok => {
                if (ok) {
                    fetch("/api/gallery/" + this.props.gallery.id, {
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

    renameGallery() {
        popups.prompt("Rename Gallery", `Rename gallery ${this.props.gallery.title} to:`)
            .then(resp => {
                if (resp) {
                    fetch("/api/gallery/" + this.props.gallery.id, {
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
                    <div className="card-header">
                        <button className="btn btn-link float-right tooltip" data-tooltip="Delete" onClick={this.deleteGallery}><i className="fa fa-trash-alt"/></button>
                        <button className="btn btn-link float-right tooltip" data-tooltip="Rename" onClick={this.renameGallery}><i className="fa fa-edit"/></button>
                        <Link className="card-title h5" to={{
                            pathname: "/admin/gallery/" + this.props.gallery.id,
                            state: {
                                gallery: this.props.gallery
                            }
                        }}>{this.props.gallery.title}</Link>
                    </div>
                </div>
            </div>
        )
    }
}

export default GalleryCard;