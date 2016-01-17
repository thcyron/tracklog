"use strict";

import React from "react";
import classNames from "classnames";

import LogElevationChart from "./LogElevationChart";
import LogHeartrateChart from "./LogHeartrateChart";

export default class LogCharts extends React.Component {
  constructor(props) {
    super(props);
    this.state = { tab: "elevation" };
  }

  static tabs() {
    return [
      { name: "Elevation", key: "elevation" },
      { name: "Heartrate", key: "heartrate" },
    ];
  }

  get _chart() {
    switch (this.state.tab) {
    case "elevation":
      return <LogElevationChart log={this.props.log} />;
    case "heartrate":
      return <LogHeartrateChart log={this.props.log} />;
    default:
      return null;
    }
  }

  _onTabClick(tab, event) {
    event.preventDefault();

    if (this.state.tab != tab) {
      this.setState({ tab: tab });
    }
  }

  render() {
    return (
      <div className="log-charts">
        <ul className="nav nav-tabs log-charts-tabs">
          {LogCharts.tabs().map((tab) => {
            return (
              <li key={tab.key} className={classNames({ active: tab.key == this.state.tab })}>
                <a href={`#${tab.key}`} onClick={(event) => { this._onTabClick(tab.key, event) }}>{tab.name}</a>
              </li>
            );
          })}
        </ul>
        <div className="panel panel-default">
          <div className="panel-body">
            {this._chart}
          </div>
        </div>
      </div>
    );
  }
}
