import React, { Component } from "react";
import DirTable from "../../components/DirTable/DirTable";
import NavBar2 from "../../components/NavBar2/NavBar2";
import { io } from "socket.io-client";
import "./PlayerStatus.css";

class PlayerStatus extends Component {
  constructor(props) {
    super(props);
    this.state = {
      player: "0001",
      titlE: ["GAME", "GAME NAME", "STATUS"],
      dataplayer: [],
    };
  }
  componentDidMount() {
    const ENDPOINT = process.env.REACT_APP_API_URL;
    this.socket = io(ENDPOINT);
  }
  componentWillUnmount() {
    this.socket.disconnect();
  }
  playerHandler = (e) => {
    this.setState({
      player: e.target.value,
    });
  };
  statusplayer() {
    if (this.state.player !== "") {
      this.socket.emit("status", this.state.player);
      this.socket.on("player", (data) => {
        console.log(data);
        this.setState({ dataplayer: data });
      });
    }
  }
  render() {
    return (
      <>
        <NavBar2></NavBar2>
        <h1 className="text-center text-info">REALTIME GAMER STATUS</h1>
        <br />
        <div className="row">
          <div className="col-4 wrap-input100 validate-input">
            <span className="label-input100">Player</span>
            <input
              className="input100"
              id="player"
              placeholder="Type your player"
              value={this.state.player}
              onChange={this.playerHandler}
              required
            />
            <span className="fas focus-input100" data-symbol="&#xf2bd;"></span>
          </div>
          <div className="col-2 container-login100-form-btn">
            <br />
            <br />
            <div className="wrap-login100-form-btn">
              <div className="login100-form-bgbtn"> </div>
              <button
                type="submit"
                id="search"
                onClick={this.statusplayer.bind(this)}
                className="login100-form-btn"
              >
                Search
              </button>
            </div>
          </div>
        </div>
        <div className="col-8">
          <DirTable
            dirData={this.state.dataplayer}
            dirTitle={this.state.titlE}
          ></DirTable>
        </div>
      </>
    );
  }
}

export default PlayerStatus;
