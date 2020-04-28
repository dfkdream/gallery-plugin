import React, {Component} from "react"

import "../../../sass/admin-card.scss"

const popups = require('../popup');

class ImageCard extends Component {
    constructor(props) {
        super(props);
        this.deleteImage = this.deleteImage.bind(this);
        this.setDescription = this.setDescription.bind(this);
    }

    deleteImage() {
        popups.confirm("Delete Image", `Delete Image?`)
            .then(ok => {
                if (ok) {
                    fetch("/api/gallery/" + this.props.gallery.id + "/album/" + this.props.album.id + "/image/" + this.props.image.id, {
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

    setDescription() {
        popups.prompt("Set Image Description", `Set image description to:`)
            .then(resp => {
                if (resp) {
                    fetch("/api/gallery/" + this.props.gallery.id + "/album/" + this.props.album.id + "/image/" + this.props.image.id, {
                        method: "POST",
                        body: JSON.stringify({description: resp})
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
                    <div className="card-image img"
                         style={{
                             backgroundImage: `url(/api/gallery/${this.props.gallery.id}/album/${this.props.album.id}/image/${this.props.image.id}?thumb=1)`
                         }}/>
                    <div className="card-body">
                        <button className="btn btn-link float-right tooltip" data-tooltip="Delete"
                                onClick={this.deleteImage}><i className="fa fa-trash-alt"/></button>
                        <button className="btn btn-link float-right tooltip" data-tooltip="Edit description"
                                onClick={this.setDescription}><i className="fa fa-edit"/></button>
                        {this.props.image.description}
                    </div>
                </div>
            </div>
        )
    }
}

export default ImageCard