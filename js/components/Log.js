"use strict";

import React from "react";
import ReactDOM from "react-dom";

import LogMap from "./LogMap";
import LogDetails from "./LogDetails";

class LogName extends React.Component {
  constructor(props) {
    super(props);
    this.state = { name: props.log.name };
  }

  onChange(event) {
    const name = event.target.value;
    this.setState({ name: name });

    if (this.props.onChange) {
      this.props.onChange(name);
    }
  }

  render() {
    if (this.props.editing) {
      return (
        <input
          type="text"
          className="log-name-edit-field form-control"
          value={this.state.name}
          onChange={this.onChange.bind(this)} />
      );
    }

    return (
      <h1 className="log-name">{this.props.log.name}</h1>
    );
  }
}

export default class Log extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      editing: false,
    };
  }

  onEdit() {
    if (this.state.editing) {
      return;
    }

    this.setState({
      editing: true,
      name: this.props.log.name,
    });
  }

  onNameChange(name) {
    this.state.name = name;
  }

  onSave(event) {
    let newLog = this.props.log;
    newLog.name = this.state.name;

    this.props = { log: newLog };
    this.setState({ editing: false });

    window.fetch(`/logs/${this.props.log.id}`, {
      method: "PATCH",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json; charset=utf-8",
        "X-CSRF-Token": window.tracklog.csrfToken,
      },
      body: JSON.stringify({
        "name": newLog.name,
      }),
    })
    .then((data) => {
      if (data.status != 204) {
        alert("Failed to save log");
        this.setState({ editing: true });
      }
    })
    .catch((err) => {
      alert(err);
    });
  }

  onCancel(event) {
    this.setState({
      editing: false,
    });
  }

  get topRow() {
    if (this.state.editing) {
      return (
        <div className="row">
          <div className="col-md-9">
            <LogName log={this.props.log} editing={this.state.editing} onChange={this.onNameChange.bind(this)} />
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
          <LogName log={this.props.log} editing={this.state.editing} />
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
            <LogMap log={this.props.log} />
          </div>
          <div className="col-md-3">
            <LogDetails log={this.props.log} onEdit={this.onEdit.bind(this)} />
          </div>
        </div>
      </div>
    );
  }
}

export function renderLog(container, log) {
  ReactDOM.render(<Log log={log} />, container);
}
