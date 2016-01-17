"use strict";

import React from "react";
import ReactDOM from "react-dom";
import Immutable from "immutable";

import LogMap from "./LogMap";
import LogDetails from "./LogDetails";
import LogName from "./LogName";
import LogCharts from "./LogCharts";

import Dispatcher from "../Dispatcher";

import LogStore from "../stores/LogStore";

export default class Log extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      log: props.log,
      editing: false,
    };

    LogStore.init(this.props.log);

    LogStore.on("change", () => {
      this.setState({
        log: LogStore.log,
      });
    });
  }

  onEdit() {
    if (!this.state.editing) {
      this.setState({
        editing: true,
        oldLog: this.state.log,
      });
    }
  }

  onSave(event) {
    this.setState({
      editing: false,
      oldLog: null,
    });

    const tags = this.state.log.get("tags").filter(tag => tag.length > 0);

    window.fetch(`/logs/${this.state.log.get("id")}`, {
      method: "PATCH",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json; charset=utf-8",
        "X-CSRF-Token": Tracklog.csrfToken,
      },
      body: JSON.stringify({
        "name": this.state.log.get("name"),
        "tags": tags.toJSON(),
      }),
    })
    .then((data) => {
      if (data.status != 204) {
        alert("Failed to save log");
        this.setState({
          editing: true,
          oldLog: this.state.log,
        });
      }
    })
    .catch((err) => {
      alert(err);
    });
  }

  onCancel(event) {
    this.setState({
      editing: false,
      log: this.state.oldLog,
      oldLog: null,
    });
  }

  get topRow() {
    if (this.state.editing) {
      return (
        <div className="row">
          <div className="col-md-9">
            <LogName log={this.state.log} editing={this.state.editing} />
          </div>
          <div className="col-md-3">
            <div className="row">
              <div className="col-sm-6">
                <button className="btn btn-block btn-success" onClick={this.onSave.bind(this)}>Save</button>
              </div>
              <div className="col-sm-6">
                <button className="btn btn-block btn-danger" onClick={this.onCancel.bind(this)}>Cancel</button>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return (
      <div className="row">
        <div className="col-md-12">
          <LogName log={this.state.log} editing={this.state.editing} />
        </div>
      </div>
    );
  }

  render() {
    return (
      <div className="log">
        {this.topRow}
        <div className="row">
          <div className="col-md-9">
            <LogMap log={this.state.log} mapboxAccessToken={Tracklog.mapboxAccessToken} />
            <LogCharts log={this.state.log} />
          </div>
          <div className="col-md-3">
            <LogDetails log={this.state.log} onEdit={this.onEdit.bind(this)} editing={this.state.editing} />
          </div>
        </div>
      </div>
    );
  }
}

export function renderLog(container, log) {
  const immutableLog = Immutable.fromJS(log);
  ReactDOM.render(<Log log={immutableLog} />, container);
}
