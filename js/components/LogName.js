"use strict";

import React from "react";

import Dispatcher from "../Dispatcher";

export default class LogName extends React.Component {
  _onChange(event) {
    Dispatcher.dispatch({
      type: "log-set-name",
      name: event.target.value,
    });
  }

  render() {
    if (this.props.editing) {
      return (
        <input
          type="text"
          className="log-name-edit-field form-control"
          value={this.props.log.get("name")}
          onChange={this._onChange.bind(this)} />
      );
    }

    return (
      <h1 className="log-name">{this.props.log.get("name")}</h1>
    );
  }
}
