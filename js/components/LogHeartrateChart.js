"use strict";

import React from "react";
import Highcharts from "highcharts";
import Leaflet from "leaflet";

export default class LogHeartrateChart extends React.Component {
  _createChart(container) {
    if (container == null) {
      return;
    }
    Highcharts.chart(container, {
      chart: {
        style: {
          fontFamily: `"Helvetica Neue", Helvetica, Arial, sans-serif`,
          fontSize: "12px",
        },
      },
      title: {
        text: null,
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
          text: "Heartrate"
        },
        labels: {
          format: "{value} bpm"
        },
      },
      legend: {
        enabled: false,
      },
      series: [
        {
          name: "Heartrate",
          color: "rgb(30, 179, 0)",
          data: this._dataFromLog(this.props.log),
        },
      ],
    });
  }

  _dataFromLog(log) {
    let data = [];
    let distance = 0;

    log.get("tracks").forEach((track) => {
      track.forEach((point, i) => {
        if (i > 0) {
          const lastPoint = track.get(i - 1);
          const p0 = new Leaflet.latLng(lastPoint.get("lat"), lastPoint.get("lon"));
          const p1 = new Leaflet.latLng(point.get("lat"), point.get("lon"));
          distance += p0.distanceTo(p1);
        }
        const hr = point.get("hr");
        if (hr) {
          data.push([distance, hr]);
        }
      });
    });

    return data;
  }

  render() {
    return (
      <div className="log-chart-chart log-chart-heartrate" ref={this._createChart.bind(this)}></div>
    );
  }
}
