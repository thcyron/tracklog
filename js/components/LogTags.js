"use strict";

import React from "react";

import Dispatcher from "../Dispatcher";

export default class LogTags extends React.Component {
  constructor(props) {
    super(props);
    this.state = this._initialStateForProps(props);
  }

  componentWillReceiveProps(nextProps) {
    this.setState(this._initialStateForProps(nextProps));
  }

  _initialStateForProps(props) {
    let tags = props.log.get("tags");

    if (props.editing && tags.filter(tag => tag.length == 0).size == 0) {
      tags = tags.push("");
    }

    return {
      tags: tags,
    };
  }

  _onTagChange(index, tag) {
    this.setState({
      tags: this.state.tags.set(index, tag),
    }, () => {
      Dispatcher.dispatch({
        type: "log-set-tags",
        tags: this.state.tags.filter(tag => tag.length > 0),
      });
    });
  }

  render() {
    let items;

    if (this.props.editing) {
      items = this.state.tags.toJS().map((tag, i) => {
        const onChange = (event) => {
          return this._onTagChange(i, event.target.value);
        };
        return (
          <li key={i} className="list-group-item">
            <input type="text" className="form-control input-sm" value={tag} onChange={onChange} />
          </li>
        );
      });
    } else {
      items = this.state.tags.toJS().map((tag, i) => {
        return (
          <li key={i} className="list-group-item">
            <span className="label label-primary">{tag}</span>
          </li>
        );
      });
    }

    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h4 className="panel-title">Tags</h4>
        </div>
        <ul className="list-group">
          {items}
        </ul>
      </div>
    );
  }
}
