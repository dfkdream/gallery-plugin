import React, {Component} from "react"
import "spectre.css/dist/spectre-icons.css"

class ImagesEmpty extends Component{
    render(){
        return(
            <div className="empty">
                <div className="empty-icon">
                    <i className="icon icon-3x icon-photo" />
                </div>
                <p className="empty-title h5">You have not uploaded any image yet</p>
                <p className="empty-subtitle">Click <strong>Upload Images</strong> to upload images.</p>
            </div>
        )
    }
}

export default ImagesEmpty
