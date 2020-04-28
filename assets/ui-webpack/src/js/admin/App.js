import React, {Component} from "react"
import ReactDOM from "react-dom"
import {BrowserRouter as Router, Route, Switch} from "react-router-dom"

import Galleries from "./routes/galleries";
import Albums from "./routes/albums";
import Images from "./routes/images";

function NoMatch(props){
    return <h1>Not Found</h1>
}

class App extends Component{
    render(){
        return(
            <Router>
                <div>
                    <Switch>
                        <Route exact path="/admin/gallery/" component={Galleries}/>
                        <Route exact path="/admin/gallery/:gid" component={Albums}/>
                        <Route exact path="/admin/gallery/:gid/album/:aid" component={Images}/>
                        <Route component={NoMatch}/>
                    </Switch>
                </div>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById("app"));
