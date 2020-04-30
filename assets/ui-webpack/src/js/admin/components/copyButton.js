import React, {Component} from "react"
import copy from "copy-to-clipboard";

class CopyButton extends Component {
    constructor(props) {
        super(props);

        this.copy = this.copy.bind(this);

        this.state = {
            copied: false
        }
    }

    copy() {
        copy(this.props.text);
        this.setState({copied: true});
        window.setTimeout(() => {
            this.setState({copied: false});
        }, this.props.timeout);
    }

    render() {
        if (this.state.copied) {
            return (
                <button className="btn btn-success"><i className="fa fa-check"/>Copied!</button>
            );
        }
        return (
            <button className="btn btn-primary" onClick={this.copy}>
                <i className="fa fa-copy"/>{this.props.caption}
            </button>
        )
    }
}

export default CopyButton;
