"use strict";

import React from "react";

export default class LogTags extends React.Component {
  constructor(props) {
    super(props);

    let tags = this.props.log.get("tags");
    this.state = {
      tags: tags,
    };
  }

  render() {
    const tags = this.state.tags.map((tag, i) => {
      return (
        <li key={i} className="list-group-item">
          <span className="label label-primary">{tag}</span>
        </li>
      );
    }).toJS();

    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h4 className="panel-title">Tags</h4>
        </div>
        <ul className="list-group">
          {tags}
        </ul>
      </div>
    );
  }
}
