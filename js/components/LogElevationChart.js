"use strict";

import React from "react";
import Highcharts from "highcharts";

export default class LogElevationChart extends React.Component {
  _createChart(container) {
    if (container == null) {
      return;
    }
    Highcharts.chart(container, {
      chart: {
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
        formatter: function() { return `<b>${this.y} m</b>`; },
      },
      xAxis: {
        title: {
          text: "Distance",
        },
        labels: {
          formatter: function() { return `${this.value / 1000} km`; },
        },
      },
      yAxis: {
        title: {
          text: "Elevation"
        },
        labels: {
          format: "{value} m"
        },
      },
      legend: {
        enabled: false,
      },
      series: [
        {
          name: "Elevation",
          color: "rgb(30, 179, 0)",
          data: this._dataFromLog(this.props.log),
          animation: false,
        },
      ],
    });
  }

  _dataFromLog(log) {
    let data = [];
    let distance = 0;

    log.get("tracks").forEach((track) => {
      track.forEach((point, i) => {
        const ele = point.get("ele");
        if (ele) {
          const distance = point.get("cumulated_distance");
          data.push([Math.round(distance), Math.round(ele)]);
        }
      });
    });

    return data;
  }

  render() {
    return (
      <div className="log-chart-chart log-chart-elevation" ref={this._createChart.bind(this)}></div>
    );
  }
}
