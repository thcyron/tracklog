"use strict";

import "whatwg-fetch";

import {uploadLogs} from "./upload";
import {renderLog} from "./components/Log";

document.addEventListener("DOMContentLoaded", () => {
  const nodes = document.querySelectorAll("a[data-method]");
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i];
    const method = node.attributes.getNamedItem("data-method").value;
    node.appendChild(createForm(node.href, method));
    node.addEventListener("click", (event) => {
      event.preventDefault();
      const a = event.target;
      const form = a.querySelector("form");
      form.submit();
    });
  }
});

function createForm(href, method) {
  const form = document.createElement("form");
  form.method = "POST";
  form.action = href;
  form.style = "display: none;"

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

    const changeButton = ({textAttribute, disabled}) => {
      node.textContent = node.attributes.getNamedItem(textAttribute).value;
      node.disabled = disabled;
    };

    node.addEventListener("click", () => {
      const fileInput = node.nextElementSibling;

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
          changeButton({
            textAttribute: "data-text-upload",
            disabled: false,
          });

          if (results.length == 1) {
            const id = results[0].id;
            window.location = `/logs/${id}`;
          } else {
            window.location = "/logs";
          }
        })
        .catch((err) => {
          changeButton({
            textAttribute: "data-text-upload",
            disabled: false,
          });

          alert(err);
        });

        changeButton({
          textAttribute: "data-text-uploading",
          disabled: true,
        });
      });

      fileInput.click();
    });
  }
});

window.Tracklog = {
  renderLog: renderLog,
};
