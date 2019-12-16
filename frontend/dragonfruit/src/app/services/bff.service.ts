import { Injectable, EventEmitter } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { BehaviorSubject, Observable, throwError } from "rxjs";

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

export class RoomRef {
  private _room: BehaviorSubject<Room>;
  private _logout;

  get room() {
    if (this._room) {
      return this._room.value;
    }

    return undefined;
  }

  constructor(room: BehaviorSubject<Room>, logout: () => void) {
    this._room = room;
    this._logout = logout;
  }

  logout = () => {
    if (this._logout) {
      return this._logout();
    }

    return undefined;
  };

  subject = (): BehaviorSubject<Room> => {
    return this._room;
  };
}

@Injectable({
  providedIn: "root"
})
export class BFFService {
  constructor(private router: Router) {}

  getRoom = (key: string | number): RoomRef => {
    const room = new BehaviorSubject<Room>(undefined);

    // use ws for http, wss for https
    let protocol = "ws:";
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    }

    const endpoint = protocol + "//" + window.location.host + "/ws/" + key;
    const ws = new WebSocket(endpoint);

    const roomRef = new RoomRef(room, () => {
      console.log("closing room connection", room.value.id);

      // close the websocket
      ws.close();

      // say that we are done with sending rooms
      room.complete();

      // route back to login page since we are gonna need a new code
      this.router.navigate(["/login"], { replaceUrl: true });
    });

    // handle incoming messages from bff
    ws.onmessage = msg => {
      const data = JSON.parse(msg.data);

      for (const k in data) {
        switch (k) {
          case "room":
            console.log("new room", data[k]);
            room.next(data[k]);

            break;
          default:
            console.warn(
              "got key '" + k + "', not sure how to handle that message"
            );
        }
      }
    };

    ws.onerror = err => {
      console.error("websocket error", err);
      room.error(err);
    };

    return roomRef;
  };

  setInput(display: Display, input: Input) {
    const kv = {
      setInput: {
        display: display.id,
        input: input.id
      }
    };

    console.log(JSON.stringify(kv));
    // this.ws.send(JSON.stringify(kv));
  }

  setVolume(ad: AudioDevice, level: number) {
    const kv = {
      setVolume: {
        audioDevice: ad.id,
        level: level
      }
    };

    console.log(JSON.stringify(kv));
    // this.ws.send(JSON.stringify(kv));
  }

  setMuted(ad: AudioDevice, m: boolean) {
    const kv = {
      setMuted: {
        audioDevice: ad.id,
        muted: m
      }
    };

    console.log(JSON.stringify(kv));
    // this.ws.send(JSON.stringify(kv));
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
    // this.ws.send(JSON.stringify(kv));
  }

  turnOffRoom() {
    const kv = {
      turnOffRoom: {}
    };

    console.log(JSON.stringify(kv));
    // this.ws.send(JSON.stringify(kv));
  }
}
