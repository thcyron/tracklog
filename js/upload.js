"use strict";

export function uploadLog({file}) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (event) => {
      window.fetch("/logs", {
        method: "POST",
        credentials: "same-origin",
        headers: {
          "Content-Type": "application/json; charset=utf-8",
          "X-CSRF-Token": window.tracklog.csrfToken,
        },
        body: JSON.stringify({
          "filename": file.name,
          "gpx": reader.result,
        }),
      })
      .then((data) => {
        return data.json();
      })
      .then((json) => {
        if (json.id) {
          resolve({
            id: json.id,
          });
        } else {
          reject(new Error("bad response from server"));
        }
      })
      .catch((err) => {
        reject(err);
      });
    };

    reader.onerror = () => {
      reject(reader.error);
    };

    reader.readAsText(file);
  });
}

export function uploadLogs({files}) {
  return Promise.all(files.map((file) => {
    return uploadLog({
      file: file,
    });
  }));
}
