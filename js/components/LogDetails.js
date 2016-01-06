"use strict";

import React from "react";

export default class LogDetails extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let hrZones = "";

    if (this.props.log.hr) {
      hrZones = (
        <div className="panel panel-default">
          <div className="panel-heading">
            <h4 className="panel-title">Heart Rate</h4>
          </div>
          <ul className="list-group">
            <li className="list-group-item">
              <div className="progress log-heart-rate-bar">
                {["none", "easy", "fatburning", "aerobic", "anaerobic", "red"].map((zone) => {
                  return <div key={zone} className={`progress-bar heart-rate-${zone}`} style={{ width: `${this.props.log.hrzones[zone]}%` }}></div>;
                })}
              </div>
            </li>
            <li className="list-group-item text-heart-rate-red">
              <span title="≥175">Red</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.red)}%</span>
            </li>
            <li className="list-group-item text-heart-rate-anaerobic">
              <span title="164–175">Anaerobic</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.anaerobic)}%</span>
            </li>
            <li className="list-group-item text-heart-rate-aerobic">
              <span title="153–164">Aerobic</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.aerobic)}%</span>
            </li>
            <li className="list-group-item text-heart-rate-fatburning">
              <span title="142–153">Fat Burning</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.fatburning)}%</span>
            </li>
            <li className="list-group-item text-heart-rate-easy">
              <span title="131–142">Easy</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.easy)}%</span>
            </li>
            <li className="list-group-item text-heart-rate-none">
              <span>None</span>
              <span className="pull-right">{Math.round(this.props.log.hrzones.none)}%</span>
            </li>
          </ul>
        </div>
      );
    }

    let details = [
      ["Start", this.props.log.start],
      ["End", this.props.log.end],
      ["Duration", this.props.log.duration],
      ["Distance", this.props.log.distance],
      ["∅ Speed", this.props.log.speed],
      ["∅ Pace", this.props.log.pace],
    ];

    if (this.props.log.hr) {
      details.push(["∅ HR", this.props.log.hr]);
    }

    let dlElements = [];
    details.forEach((detail) => {
      dlElements.push(<dt>{detail[0]}</dt>);
      dlElements.push(<dd>{detail[1]}</dd>);
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
        {hrZones}
        <ul className="list-group">
          <li className="list-group-item"><a href={`/logs/${this.props.log.id}/download`}>Download .gpx file</a></li>
          <li className="list-group-item"><a href={`/logs/${this.props.log.id}`} data-method="delete">Delete</a></li>
        </ul>
      </div>
    );
  }
}

export function renderLog(container, log) {
  ReactDOM.render(<Log log={log} />, container);
}
