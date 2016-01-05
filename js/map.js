"use strict";

import Leaflet from "leaflet";

export default class Map {
  constructor({log, container}) {
    this.log = log;

    this.map = Leaflet.map(container);
    this.map.addLayer(Leaflet.tileLayer("http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"));

    const latlngs = this.log.tracks.map(track => {
      return track.map(point => {
        return [point.lat, point.lon];
      });
    });

    const multiPolyline = Leaflet.multiPolyline(latlngs, { color: "red" });
    multiPolyline.addTo(this.map);
    this.map.fitBounds(multiPolyline.getBounds());
  }
}
