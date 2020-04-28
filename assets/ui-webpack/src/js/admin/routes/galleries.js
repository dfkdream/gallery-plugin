import React, {Component} from "react"

import GalleriesEmpty from "../components/galleriesEmpty";

import GalleryCard from "../components/galleryCard";

import "../../../sass/admin-location-header.scss";

const popups = require('../popup');

class GalleriesHeader extends Component{
    constructor(props){
        super(props);
    }

    render() {
        return(
            <div className="location-header">
               <h1>Galleries</h1>
                <div>
                    <button className="btn btn-primary" onClick={this.props.onClick}>
                        <i className="fa fa-plus" />
                        Add Gallery
                    </button>
                </div>
            </div>
        )
    }
}

class Galleries extends Component {
    constructor(props) {
        super(props);
        this.createGallery = this.createGallery.bind(this);
        this.state = {
            isLoaded: false,
            error: null,
            galleries: []
        }
    }

    createGallery(){
        popups.prompt("Create Gallery","Enter gallery title:")
            .then(resp=>{
                if (resp){
                    fetch("/api/gallery/",{
                        method: "POST",
                        body: JSON.stringify({title:resp})
                    })
                        .then(resp=>{
                            if (!resp.ok){
                                popups.alert("Error",`${resp.status} ${resp.statusText}`);
                                return
                            }

                            this.loadGallery();
                        })
                }
            })
    }

    loadGallery(){
        fetch("/api/gallery/")
            .then(res => res.json())
            .then(
                (json) => {
                    this.setState({isLoaded: true, galleries: json});
                },
                (error) => {
                    this.setState({isLoaded: true, error: error});
                }
            )
    }

    componentDidMount() {
        this.loadGallery();
    }

    render() {
        const {isLoaded, error, galleries} = this.state;
        if (error){
            return (
                <div>
                    <h1>Error</h1>
                    <pre>{error.message}</pre>
                </div>
            )
        } else if (!isLoaded){
            return <h1>Loading...</h1>
        } else{
            if (galleries.length>0){
                return (
                    <div>
                        <GalleriesHeader onClick={this.createGallery}/>
                        <div className="columns">
                            {galleries.map(g=>{
                                return <GalleryCard gallery={g} onChange={()=>{this.loadGallery()}}/>
                            })}
                        </div>
                    </div>
                )
            }else{
                return (
                    <div>
                        <GalleriesHeader onClick={this.createGallery}/>
                        <GalleriesEmpty />
                    </div>
                )
            }
        }
    }
}

export default Galleries