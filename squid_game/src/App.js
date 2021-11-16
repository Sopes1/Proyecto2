import { BrowserRouter, Route, Switch, Redirect } from "react-router-dom";
import { Component } from "react";
import "./App.css";
import Topjuegos from "./views/Topjuegos/Topjuegos";
import PlayerStatus from "./views/PlayerStatus/PlayerStatus";
import Report from "./views/Report/Report";
import DataMongo from "./views/DataMongo/DataMongo";

export default class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <Switch>
          <Route exact path="/">
            <Redirect to="topreport" />
          </Route>
          <Route exact path="/topreport">
            <Topjuegos></Topjuegos>
          </Route>
          <Route exact path="/status">
            <PlayerStatus></PlayerStatus>
          </Route>
          <Route exact path="/report">
            <Report></Report>
          </Route>
          <Route exact path="/datamongo">
            <DataMongo></DataMongo>
          </Route>
        </Switch>
      </BrowserRouter>
    );
  }
}
