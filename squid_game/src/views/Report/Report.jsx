import React, { Component } from "react";
import ColumnChart from "../../components/ColummChart/ColumnChart";
import PieChart from "../../components/PieChart/PieChart";
import NavBar2 from "../../components/NavBar2/NavBar2";
import { io } from "socket.io-client";
import "./Report.css";

class Report extends Component {
  constructor(props) {
    super(props);
    this.state = {
      dataworks: [
        { label: "worker Kafka", y: 30000 },
        { label: "worker RabbitMQ", y: 35000 },
        { label: "worker PubSub", y: 25000 },
      ],
      datagame: [
        { label: "worker Kafka", y: 30000 },
        { label: "worker RabbitMQ", y: 35000 },
        { label: "worker PubSub", y: 25000 },
      ],
    };
  }
  componentDidMount() {
    const ENDPOINT = process.env.REACT_APP_API_URL;
    this.socket = io(ENDPOINT);
    this.socket.emit("graficas");
    this.socket.on("worker", (data) => {
      console.log(data);
      this.setState({ dataworks: data });
    });

    this.socket.on("topgame", (data) => {
      console.log(data);
      this.setState({ datagame: data });
    });
  }
  componentWillUnmount() {
    this.socket.disconnect();
  }
  render() {
    return (
      <>
        <NavBar2></NavBar2>
        <div className="col-8">
          <ColumnChart datapoints={this.state.dataworks}></ColumnChart>
        </div>
        <br />
        <br />
        <div>
          <PieChart datapoints={this.state.datagame}></PieChart>
        </div>
      </>
    );
  }
}

export default Report;
