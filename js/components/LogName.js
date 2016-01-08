"use strict";

import React from "react";

export default class LogName extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      name: props.log.get("name"),
    };
  }

  componentWillReceiveProps(props) {
    this.state = {
      name: props.log.get("name"),
    };
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
      <h1 className="log-name">{this.props.log.get("name")}</h1>
    );
  }
}
