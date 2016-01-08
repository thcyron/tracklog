"use strict";

import React from "react";

import LogTags from "./LogTags";

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

      hrZones = (
        <div className="panel panel-default">
          <div className="panel-heading">
            <h4 className="panel-title">Heart Rate</h4>
          </div>
          <ul className="list-group">
            <li className="list-group-item">
              <div className="progress log-heart-rate-bar">
                {["none", "easy", "fatburning", "aerobic", "anaerobic", "red"].map((zone) => {
                  return <div key={zone} className={`progress-bar heart-rate-${zone}`} style={{ width: `${zones.get(zone)}%` }}></div>;
                })}
              </div>
            </li>
            <li className="list-group-item text-heart-rate-red">
              <span title="≥175">Red</span>
              <span className="pull-right">{Math.round(zones.get("red"))}%</span>
            </li>
            <li className="list-group-item text-heart-rate-anaerobic">
              <span title="164–175">Anaerobic</span>
              <span className="pull-right">{Math.round(zones.get("anaerobic"))}%</span>
            </li>
            <li className="list-group-item text-heart-rate-aerobic">
              <span title="153–164">Aerobic</span>
              <span className="pull-right">{Math.round(zones.get("aerobic"))}%</span>
            </li>
            <li className="list-group-item text-heart-rate-fatburning">
              <span title="142–153">Fat Burning</span>
              <span className="pull-right">{Math.round(zones.get("fatburning"))}%</span>
            </li>
            <li className="list-group-item text-heart-rate-easy">
              <span title="131–142">Easy</span>
              <span className="pull-right">{Math.round(zones.get("easy"))}%</span>
            </li>
            <li className="list-group-item text-heart-rate-none">
              <span>None</span>
              <span className="pull-right">{Math.round(zones.get("none"))}%</span>
            </li>
          </ul>
        </div>
      );
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
