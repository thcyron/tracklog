"use strict";

import React from "react";
import ReactDOM from "react-dom";

import LogMap from "./LogMap";
import LogDetails from "./LogDetails";

export default class Log extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div className="log">
        <div className="row">
          <div className="col-md-12">
            <h1 className="log-name">{this.props.log.name}</h1>
          </div>
        </div>
        <div className="row">
          <div className="col-md-9">
            <LogMap log={this.props.log} />
          </div>
          <div className="col-md-3">
            <LogDetails log={this.props.log} />
          </div>
        </div>
      </div>
    );
  }
}

export function renderLog(container, log) {
  ReactDOM.render(<Log log={log} />, container);
}
