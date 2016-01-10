"use strict";

import React from "react";
import Chart from "chart.js";

export default class LogHeartrateZones extends React.Component {
  constructor(props) {
    super(props);
  }

  _createChart(canvas) {
    if (canvas != null) {
      new Chart(canvas.getContext("2d")).Doughnut(this.data, {
        animation: false,
        tooltipTemplate: "<%= label %>: <%= value %>%",
      });
    }
  }

  get data() {
    return [
      {
        value: Math.round(this.props.zones.get("red")),
        label: "Red",
        color: "#d9534f",
      },
      {
        value: Math.round(this.props.zones.get("anaerobic")),
        label: "Anaerobic",
        color: "#f0ad4e",
      },
      {
        value: Math.round(this.props.zones.get("aerobic")),
        label: "Aerobic",
        color: "#80c780",
      },
      {
        value: Math.round(this.props.zones.get("fatburning")),
        label: "Fat Burning",
        color: "#449d44",
      },
      {
        value: Math.round(this.props.zones.get("easy")),
        label: "Easy",
        color: "#5bc0de",
      },
      {
        value: Math.round(this.props.zones.get("none")),
        label: "None",
        color: "#777777",
      },
    ];
  }

  render() {
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h4 className="panel-title">Heart Rate Zones</h4>
        </div>
        <div className="panel-body text-center">
          <canvas width="220" height="220" ref={this._createChart.bind(this)}></canvas>
        </div>
      </div>
    );
  }
}
