import React, {Component} from "react"
import ReactDOM from "react-dom"

import AlbumsPage from "./components/albumsPage";
import ImagesPage from "./components/imagesPage";

const container = document.getElementById("app");

class App extends Component {
    constructor(props) {
        super(props);

        this.loadAlbums = this.loadAlbums.bind(this);
        this.loadImages = this.loadImages.bind(this);

        this.state = {
            page: 0, // 0->Albums, 1...->Images (albumId)
            isLoading: true,
            gallery: null,
            albums: null,
            images: new Map(),
            error: null,
        };
    }

    componentDidMount() {
        window.addEventListener("hashchange", () => {
            this.loadFromHref();
        });

        fetch("/api/gallery/" + container.dataset.gid)
            .then(resp => resp.json())
            .then(
                (json) => {
                    this.setState({gallery: json});

                    fetch("/api/gallery/" + this.state.gallery.id + "/albums")
                        .then(resp => resp.json())
                        .then(
                            (json) => {
                                this.setState({albums: json});
                                this.loadFromHref();
                            },
                            (error) => {
                                this.setState({isLoading: false, error: error.message});
                            }
                        )
                },
                (error) => {
                    this.setState({isLoading: false, error: error.message});
                }
            )
    }

    loadFromHref() {
        let href = location.hash.split("/").filter(v => !v.includes("#")).join("/");
        if (href === "") {
            this.loadAlbums();
        } else {
            const aid = parseInt(href);
            if (aid) this.loadImages(aid);
            else this.setState({isLoading: false, error: {message: "Error parsing href"}});
        }
    }

    loadAlbums() {
        if (!this.state.gallery) return;
        if (!this.state.albums) return;

        this.setState({page: 0, isLoading: false});
    }

    loadImages(aid) {
        if (!this.state.gallery) return;
        if (!this.state.albums) return;

        if (this.state.images.has(aid)) {
            this.setState({page: aid});
            return;
        }

        this.setState({isLoading: true, error: null});
        fetch("/api/gallery/" + this.state.gallery.id + "/album/" + aid + "/images")
            .then(resp => resp.json())
            .then(
                (json) => {
                    this.setState({isLoading: false, images: this.state.images.set(aid, json), page: aid});
                },
                (error) => {
                    this.setState({isLoading: false, error: error.message});
                }
            )
    }

    render() {
        if (this.state.isLoading) {
            return <h1>Loading...</h1>;

        } else if (this.state.error) {
            return (
                <div>
                    <h1>Error</h1>
                    <pre>{this.state.error}</pre>
                </div>
            );

        } else if (this.state.page === 0) {
            return <AlbumsPage gallery={this.state.gallery}
                               albums={this.state.albums}
                               loadImages={this.loadImages}/>;

        } else if (this.state.page > 0) {
            return <ImagesPage gallery={this.state.gallery}
                               album={this.state.albums.filter(a => a.id === this.state.page)[0]}
                               images={this.state.images.get(this.state.page)}
                               loadAlbums={this.loadAlbums}/>;

        } else {
            return <h1>Not Found</h1>;
        }
    }
}

ReactDOM.render(<App/>, container);
