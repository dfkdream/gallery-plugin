import React, {Component} from "react"

import ImageCard from "./imageCard";

import "../../../sass/gallery.scss";

class ImagesPage extends Component {
    render() {
        return (
            <div>
                <nav className="breadcrumb is-large">
                    <ul>
                        <li key={"g-"+this.props.gallery.id}>
                            <a href="#!/" onClick={() => {
                                this.props.loadAlbums();
                            }}>{this.props.gallery.title}</a>
                        </li>
                        <li className="is-active" key={"a-"+this.props.album.id}><a>{this.props.album.title}</a></li>
                    </ul>
                </nav>
                <div className="columns is-multiline">
                    {this.props.images.map((i, idx) => {
                        return (
                            <div className="column is-one-third" key={"i-"+i.id}>
                                <ImageCard gallery={this.props.gallery} album={this.props.album} images={this.props.images} index={idx}/>
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    }
}

export default ImagesPage;
