import React, { Component } from "react";
import { CanvasJSChart } from "canvasjs-react-charts";
import "./PieChart.css";
class PieChart extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    const options = {
      exportEnabled: true,
      animationEnabled: true,
      title: {
        text: "TOP JUEGOS",
      },
      data: [
        {
          type: "pie",
          startAngle: 75,
          toolTipContent: "<b>{label}</b>: {y}",
          showInLegend: "true",
          legendText: "{label}",
          indexLabelFontSize: 16,
          indexLabel: "{label} - {y}",
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

export default PieChart;
