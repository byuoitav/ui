import { Injectable, EventEmitter } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { BehaviorSubject, Observable, throwError } from "rxjs";
import { Router, Event, ActivationEnd } from "@angular/router";
import { MatDialog } from "@angular/material";

import {
  Room,
  ControlGroup,
  Display,
  Input,
  AudioDevice,
  AudioGroup,
  PresentGroup
} from "../objects/control";
import { ErrorDialog } from "../dialogs/error/error.dialog";

export class RoomRef {
  private _room: BehaviorSubject<Room>;
  private _ws: WebSocket;
  private _logout: () => void;

  get room() {
    if (this._room) {
      return this._room.value;
    }

    return undefined;
  }

  constructor(room: BehaviorSubject<Room>, ws: WebSocket, logout: () => void) {
    this._room = room;
    this._ws = ws;
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

  /* control functions */
  setInput = (displayID: string, inputID: string) => {
    const kv = {
      setInput: {
        display: displayID,
        input: inputID
      }
    };

    this._ws.send(JSON.stringify(kv));
  };

  setVolume = (audioDeviceID: string, level: number) => {
    const kv = {
      setVolume: {
        audioDevice: audioDeviceID,
        level: level
      }
    };

    this._ws.send(JSON.stringify(kv));
  };

  setMuted = (audioDeviceID: string, muted: boolean) => {
    const kv = {
      setMuted: {
        audioDevice: audioDeviceID,
        muted: muted
      }
    };

    this._ws.send(JSON.stringify(kv));
  };

  setPower = (displays: Display[], power: string) => {
    const kv = {
      setPower: {
        display: [],
        status: power
      }
    };

    for (const disp of displays) {
      kv.setPower.display.push(disp.id);
    }

    this._ws.send(JSON.stringify(kv));
  };

  turnOff = () => {
    const kv = {
      turnOffRoom: {}
    };

    this._ws.send(JSON.stringify(kv));
  };
}

@Injectable({
  providedIn: "root"
})
export class BFFService {
  constructor(private router: Router, private dialog: MatDialog) {
    // do things based on route changes
    this.router.events.subscribe(event => {
      if (event instanceof ActivationEnd) {
        const snapshot = event.snapshot;

        // show error if the "error" query paramater is present
        if (snapshot && snapshot.queryParams && snapshot.queryParams.error) {
          this.error(snapshot.queryParams.error);
        }
      }
    });
  }

  getRoom = (key: string | number): RoomRef => {
    const room = new BehaviorSubject<Room>(undefined);

    // use ws for http, wss for https
    let protocol = "ws:";
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    }

    const endpoint = protocol + "//" + window.location.host + "/ws/" + key;
    const ws = new WebSocket(endpoint);

    const roomRef = new RoomRef(room, ws, () => {
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

  error = (msg: string) => {
    console.log("showing error", msg);
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      return dialog.componentInstance instanceof ErrorDialog;
    });

    if (dialogs.length > 0) {
      // do something?
      // either open another one on top of the current ones,
      // change the text on the current one?
    } else {
      // open a new dialog
      const ref = this.dialog.open(ErrorDialog, {
        width: "80vw",
        data: {
          msg: msg
        }
      });

      ref.afterClosed().subscribe(result => {
        this.router.navigate([], {
          queryParams: { error: null },
          queryParamsHandling: "merge",
          preserveFragment: true
        });
      });
    }
  };

  /*
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
  */
}
