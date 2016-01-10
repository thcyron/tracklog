"use strict";

import React from "react";

import LogTags from "./LogTags";
import LogHeartrateZones from "./LogHeartrateZones";

export default class LogDetails extends React.Component {
  constructor(props) {
    super(props);
  }

  onEditClick(event) {
    event.preventDefault();

    if (this.props.onEdit) {
      this.props.onEdit();
    }
  }

  render() {
    let tags = "", hrZones = "";

    if (this.props.log.get("tags").size > 0 || this.props.editing) {
      tags = <LogTags log={this.props.log} editing={this.props.editing} />;
    }

    if (this.props.log.get("hr")) {
      const zones = this.props.log.get("hrzones");
      hrZones = <LogHeartrateZones zones={zones} />;
    }

    let details = [
      ["Start", this.props.log.get("start")],
      ["End", this.props.log.get("end")],
      ["Duration", this.props.log.get("duration")],
      ["Distance", this.props.log.get("distance")],
      ["∅ Speed", this.props.log.get("speed")],
      ["∅ Pace", this.props.log.get("pace")],
    ];

    if (this.props.log.get("hr")) {
      details.push(["∅ HR", this.props.log.get("hr")]);
    }

    let dlElements = [];
    details.forEach((detail, i) => {
      dlElements.push(<dt key={2*i}>{detail[0]}</dt>);
      dlElements.push(<dd key={2*i+1}>{detail[1]}</dd>);
    });

    return (
      <div className="log-details">
        <div className="panel panel-default">
          <div className="panel-heading">
            <h4 className="panel-title">Details</h4>
          </div>
          <div className="panel-body kill-bottom-margin">
            <dl className="dl-horizontal-small">
              {dlElements}
            </dl>
          </div>
        </div>
        {tags}
        {hrZones}
        <ul className="list-group">
          <li className="list-group-item"><a href={`/logs/${this.props.log.get("id")}/download`}>Download .gpx file</a></li>
          {(() => {
            if (!this.props.editing) {
              return [
                <li key="edit" className="list-group-item"><a href="#edit" onClick={this.onEditClick.bind(this)}>Edit</a></li>,
                <li key="delete" className="list-group-item"><a href={`/logs/${this.props.log.get("id")}`} data-method="delete">Delete</a></li>,
              ];
            }
          })()}
        </ul>
      </div>
    );
  }
}
