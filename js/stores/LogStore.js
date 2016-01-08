"use strict";

import {EventEmitter} from "events";

import Dispatcher from "../Dispatcher";

class LogStore extends EventEmitter {
  init(log) {
    this._log = log;
  }

  get log() {
    return this._log;
  }

  constructor() {
    super();

    Dispatcher.register((action) => {
      switch (action.type) {
      case "log-set-name":
        this._log = this._log.set("name", action.name);
        break;
      case "log-set-tags":
        this._log = this._log.set("tags", action.tags);
        break;
      default:
        return; // do not emit change event if no action was triggered
      }

      this.emit("change");
    });
  }
}

export default new LogStore();
