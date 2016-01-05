"use strict";

import $ from "../node_modules/jquery/dist/jquery";
import "whatwg-fetch";

import Log from "./log";
import Map from "./map";
import {uploadLogs} from "./upload";

$(() => {
  $("a[data-method]").on("click", (event) => {
    event.preventDefault();
    const $a = $(event.target);
    $("<form>")
      .attr("method", "POST")
      .attr("action", $a.attr("href"))
      .append($("<input>").attr({
        type: "hidden",
        name: "_csrf",
        value: tracklog.csrfToken,
      }))
      .append($("<input>").attr({
        type: "hidden",
        name: "_method",
        value: $a.attr("data-method").toUpperCase(),
      }))
      .hide()
      .appendTo(document.body)
      .submit();
  });

  $(".logs-upload-button").on("click", (event) => {
    event.preventDefault();

    const $input = $("<input>")
      .attr("type", "file")
      .prop("multiple", true)
      .hide()
      .appendTo(document.body)
      .trigger("click");

    $input.on("change", (event) => {
      const input = event.target;
      const files = input.files;

      let filesArray = [];
      for (let i = 0; i < files.length; i++) {
        filesArray.push(files[i]);
      }

      uploadLogs({
        files: filesArray,
      })
      .then((results) => {
        if (results.length == 1) {
          const id = results[0].id;
          window.location = `/logs/${id}`;
        } else {
          window.location = "/logs";
        }
      })
      .catch((err) => {
        alert(err);
      });
    });
  });

  $("form").on("submit", (event) => {
    const $input = $("<input>").attr({
      type: "hidden",
      name: "_csrf",
      value: tracklog.csrfToken,
    });
    $(event.target).append($input);
    return true;
  });
});

window.tracklog = {
  Log: Log,
  Map: Map,
};
