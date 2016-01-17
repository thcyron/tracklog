"use strict";

import "whatwg-fetch";

import {uploadLogs} from "./upload";
import {renderLog} from "./components/Log";

document.addEventListener("DOMContentLoaded", () => {
  const nodes = document.querySelectorAll("a[data-method]");
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i];
    node.addEventListener("click", (event) => {
      event.preventDefault();
      const a = event.target;
      const method = a.attributes.getNamedItem("data-method").value;
      if (method) {
        createForm(a.href, method).submit();
      }
    });
  }
});

function createForm(href, method) {
  const form = document.createElement("form");
  form.method = "POST";
  form.action = href;

  const methodField = document.createElement("input");
  methodField.type = "hidden";
  methodField.name = "_method";
  methodField.value = method.toUpperCase();
  form.appendChild(methodField);

  const csrfField = document.createElement("input");
  csrfField.type = "hidden";
  csrfField.name = "_csrf";
  csrfField.value = Tracklog.csrfToken;
  form.appendChild(csrfField);

  return form;
}

document.addEventListener("DOMContentLoaded", () => {
  const nodes = document.querySelectorAll("form");
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i];
    node.addEventListener("submit", (event) => {
      const form = event.target;
      const csrfField = document.createElement("input");
      csrfField.type = "hidden";
      csrfField.name = "_csrf";
      csrfField.value = Tracklog.csrfToken;
      form.appendChild(csrfField);
    });
  }
});

document.addEventListener("DOMContentLoaded", () => {
  const nodes = document.querySelectorAll(".logs-upload-button");
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i];
    node.addEventListener("click", () => {
      const fileInput = document.createElement("input");
      fileInput.type = "file";
      fileInput.multiple = true;

      fileInput.addEventListener("change", (event) => {
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

      fileInput.click();
    });
  }
});

window.Tracklog = {
  renderLog: renderLog,
};
