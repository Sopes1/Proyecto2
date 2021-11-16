import React, { Component } from "react";
import DirTable from "../../components/DirTable/DirTable";
import NavBar2 from "../../components/NavBar2/NavBar2";
import { io } from "socket.io-client";
import "./Topjuegos.css";

class Topjuegos extends Component {
  constructor(props) {
    super(props);
    this.state = {
      dataE: [
        { data1: "data1", data2: "data2", data3: "data3", data4: "data4" },
        { obj1: "obj1", obj2: "obj2", obj3: "obj3", obj4: "obj4" },
        { arr1: "arr1", arr2: "arr2", arr3: "arr3", arr4: "arr4" },
        { vect1: "vect1", vect2: "vect2", vect3: "vect3", vect4: "vect4" },
      ],
      titlE: ["title1", "title2", "title3", "title4"],
      datatop: [],
      datalast: [],
      titletop: ["#PLAYER", "WINS"],
      titlelast: ["GAME", "GAME NAME", "WINNER", "#PLAYERS"],
    };
  }
  componentDidMount() {
    const ENDPOINT = process.env.REACT_APP_API_URL;
    this.socket = io(ENDPOINT);
    this.socket.emit("top");
    this.socket.on("last", (data) => {
      console.log(data);
      this.setState({ datalast: data });
    });
    this.socket.on("top", (data) => {
      console.log(data);
      this.setState({ datatop: data });
    });
  }
  componentWillUnmount() {
    this.socket.disconnect();
  }

  render() {
    return (
      <>
        <NavBar2></NavBar2>
        <h1 className="text-center text-warning">LAST 10 GAME</h1>
        <div className="row justify-content-md-center">
          <div className="col-8" style={{ paddingLeft: 20 }}>
            <DirTable
              dirData={this.state.datalast}
              dirTitle={this.state.titlelast}
            ></DirTable>
          </div>
        </div>
        <br />
        <br />
        <h1 className="text-center text-success">TOP 10 PLAYERS</h1>
        <div className="row justify-content-md-center">
          <div className="col-8" style={{ paddingLeft: 20 }}>
            <DirTable
              dirData={this.state.datatop}
              dirTitle={this.state.titletop}
            ></DirTable>
          </div>
        </div>
      </>
    );
  }
}

export default Topjuegos;
