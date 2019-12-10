import { Injectable, EventEmitter } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import {
  $WebSocket,
  WebSocketConfig
} from "angular2-websocket/angular2-websocket";
import { JsonConvert } from "json2typescript";
import {
  Device,
  UIConfig,
  IOConfiguration,
  DBRoom,
  Preset
} from "../objects/database";
import {
  Room,
  ControlGroup,
  Display,
  Input,
  AudioDevice,
  AudioGroup,
  PresentGroup
} from "../objects/control";
import { Router } from "@angular/router";

@Injectable({
  providedIn: "root"
})
export class BFFService {
  room: Room;
  done: EventEmitter<boolean>;
  ws: WebSocket;

  constructor(private router: Router) {
    this.done = new EventEmitter();
    // this.room = new Room();
  }

  connectToRoom(controlKey: string) {
    // use ws for http, wss for https
    let protocol = "ws:";
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    }

    const endpoint =
      protocol + "//" + window.location.host + "/ws/" + controlKey;
    this.ws = new WebSocket(endpoint);

    this.ws.onmessage = event => {
      console.log("ws event", event);
      this.room = JSON.parse(event.data);
      // this.room = Object.assign(new Room(), JSON.parse(event.data));

      console.log("Websocket data:", this.room);

      this.done.emit(true);
    };

    this.ws.onerror = event => {
      console.error("Websocket error", event);
      this.router.navigate(["/login"]);
    };
  }

  setInput(display: Display, input: Input) {
    const kv = {
      setInput: {
        display: display.id,
        input: input.id
      }
    };

    console.log(JSON.stringify(kv));
    this.ws.send(JSON.stringify(kv));
  }

  setVolume(ad: AudioDevice, level: number) {
    const kv = {
      setVolume: {
        audioDevice: ad.id,
        level: level
      }
    };

    console.log(JSON.stringify(kv));
    this.ws.send(JSON.stringify(kv));
  }

  setMuted(ad: AudioDevice, m: boolean) {
    const kv = {
      setMuted: {
        audioDevice: ad.id,
        muted: m
      }
    };

    console.log(JSON.stringify(kv));
    this.ws.send(JSON.stringify(kv));
  }

  setPower(displays: Display[], s: string) {
    const kv = {
      setPower: {
        display: [],
        status: s
      }
    };
    if (displays !== null) {
      for (const disp of displays) {
        kv.setPower.display.push(disp.id);
      }
    }

    console.log(JSON.stringify(kv));
    this.ws.send(JSON.stringify(kv));
  }
}
