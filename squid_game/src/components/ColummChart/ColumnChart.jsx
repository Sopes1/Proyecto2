import React, { Component } from "react";
import { CanvasJSChart } from "canvasjs-react-charts";

class ColumnChart extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    const options = {
      title: {
        text: "workers",
      },
      data: [
        {
          // Change type to "doughnut", "line", "splineArea", etc2.
          type: "column",
          dataPoints: this.props.datapoints,
        },
      ],
    };
    return (
      <div>
        <CanvasJSChart options={options} />
      </div>
    );
  }
}

export default ColumnChart;
