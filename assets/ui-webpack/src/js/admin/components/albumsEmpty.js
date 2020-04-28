import React, {Component} from "react"
import "spectre.css/dist/spectre-icons.css"

class AlbumsEmpty extends Component{
    render(){
        return(
            <div className="empty">
                <div className="empty-icon">
                    <i className="icon icon-3x icon-photo" />
                </div>
                <p className="empty-title h5">You have not created any album yet</p>
                <p className="empty-subtitle">Click <strong>Add Album</strong> to create new one.</p>
            </div>
        )
    }
}

export default AlbumsEmpty