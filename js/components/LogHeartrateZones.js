"use strict";

import React from "react";
import Highcharts from "highcharts";

export default class LogHeartrateZones extends React.Component {
 _createChart(container) {
    if (container == null) {
      return;
    }
    Highcharts.chart(container, {
      chart: {
        type: "pie",
        animation: false,
        style: {
          fontFamily: `"Helvetica Neue", Helvetica, Arial, sans-serif`,
          fontSize: "12px",
        },
      },
      title: {
        text: null,
      },
      tooltip: {
        pointFormat: "<b>{point.y}</b>",
        valueSuffix: "%",
      },
      legend: {
        enabled: false,
      },
      series: [
        {
          name: "Heartrate",
          color: "rgb(30, 179, 0)",
          data: this.data,
          dataLabels: {
            enabled: false,
          },
          animation: false,
        },
      ],
    });
  }

  get data() {
    return [
      {
        y: Math.round(this.props.zones.get("red")),
        name: "Red",
        color: "#d9534f",
      },
      {
        y: Math.round(this.props.zones.get("anaerobic")),
        name: "Anaerobic",
        color: "#f0ad4e",
      },
      {
        y: Math.round(this.props.zones.get("aerobic")),
        name: "Aerobic",
        color: "#80c780",
      },
      {
        y: Math.round(this.props.zones.get("fatburning")),
        name: "Fat Burning",
        color: "#449d44",
      },
      {
        y: Math.round(this.props.zones.get("easy")),
        name: "Easy",
        color: "#5bc0de",
      },
      {
        y: Math.round(this.props.zones.get("none")),
        name: "None",
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
          <div style={{height: "230px"}} ref={this._createChart.bind(this)}></div>
        </div>
      </div>
    );
  }
}
