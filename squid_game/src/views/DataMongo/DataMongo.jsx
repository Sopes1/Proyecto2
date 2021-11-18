import React, { Component } from "react";
import NavBar2 from "../../components/NavBar2/NavBar2";
import DirTable from "../../components/DirTable/DirTable";
import { io } from "socket.io-client";
import "./DataMongo.css";

class DataMongo extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      logs: [],
      titledata: ["#GAME", "GAME NAME", "WINNER", "PLAYERS"],
      titlelogs: [
        "REQUEST GAME",
        "#GAME",
        "GAME NAME",
        "WINNER",
        "PLAYERS",
        "WORKER",
      ],
    };
  }

  componentDidMount() {
    const ENDPOINT = process.env.REACT_APP_API_URL;
    this.socket = io(ENDPOINT);
    this.socket.emit("datamongo");
    this.socket.on("data", (data) => {
      console.log(data);
      this.setState({ data: data });
    });
    this.socket.on("logs", (data) => {
      console.log(data);
      this.setState({ logs: data });
    });
  }

  componentWillUnmount() {
    this.socket.disconnect();
  }

  render() {
    return (
      <>
        <NavBar2></NavBar2>
        <h1 className="text-center text-warning">DATA MONGO</h1>
        <div className="row justify-content-md-center">
          <div className="col-8 altura" id="altura" style={{ paddingLeft: 20 }}>
            <DirTable
              dirData={this.state.data}
              dirTitle={this.state.titledata}
            ></DirTable>
          </div>
        </div>
        <br />
        <br />
        <h1 className="text-center text-success">LOGS MONGO</h1>
        <div className="row justify-content-md-center">
          <div className="col-8 altura" id="altura" style={{ paddingLeft: 20 }}>
            <DirTable
              dirData={this.state.logs}
              dirTitle={this.state.titlelogs}
            ></DirTable>
          </div>
        </div>
      </>
    );
  }
}

export default DataMongo;
