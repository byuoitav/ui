import { Injectable, EventEmitter } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { BehaviorSubject, Observable, throwError } from "rxjs";
import { Router, Event, ActivationEnd } from "@angular/router";
import { MatDialog } from "@angular/material";

import {
  Room,
  ControlGroup,
  DisplayGroup,
  Input,
  AudioDevice,
  AudioGroup,
  PresentGroup
} from "../../../../objects/control";
// import { ErrorDialog } from "../dialogs/error/error.dialog";
// import { TurnOffRoomDialogComponent } from "../dialogs/turnOffRoom-dialog/turnOffRoom-dialog.component";

export class RoomRef {
  private _room: BehaviorSubject<Room>;
  private _ws: WebSocket;
  private _logout: () => void;
  loading: boolean;

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
        displayGroup: displayID,
        input: inputID
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  setVolume = (level: number, audioDeviceID?: string) => {
    const kv = {
      setVolume: {
        audioDevice: audioDeviceID,
        level: level
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  setMuted = (muted: boolean, audioDeviceID?: string) => {
    const kv = {
      setMuted: {
        audioDevice: audioDeviceID,
        muted: muted
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  setBlanked = (displayID: string, blanked: boolean) => {
    const kv = {
      setBlanked: {
        displayGroup: displayID,
        blanked: blanked
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  }

  setPower = (power: boolean) => {
    const kv = {
      setPower: {
        poweredOn: power
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  requestHelp = (msg: string) => {
    console.log("requesting help:", msg);
    const req = {
      helpRequest: {
        msg: msg
      }
    };

    this._ws.send(JSON.stringify(req));
  };

  raiseProjectorScreen = (screen: string) => {
    const kv = {
      raiseProjectorScreen: screen
    }

    this._ws.send(JSON.stringify(kv));
  }

  lowerProjectorScreen = (screen: string) => {
    const kv = {
      lowerProjectorScreen: screen
    }

    this._ws.send(JSON.stringify(kv));
  }

  stopProjectorScreen = (screen: string) => {
    const kv = {
      stopProjectorScreen: screen
    }

    this._ws.send(JSON.stringify(kv));
  }

  getControlKey = (cgID: string) => {
    const kv = {
      getControlKey: {
        controlGroupID: cgID
      }
    }

    this._ws.send(JSON.stringify(kv));
  }

  startSharing = (masterID: string, optionsIDs: string[]) => {
    const kv = {
      setSharing: {
        status: true,
        master: masterID,
        options: optionsIDs
      }
    }

    this._ws.send(JSON.stringify(kv));
  }

  stopSharing = (masterID: string) => {
    const kv = {
      setSharing: {
        status: false,
        master: masterID
      }
    }

    this._ws.send(JSON.stringify(kv));
  }
}

@Injectable({
  providedIn: "root"
})
export class BFFService {
  locked = true;
  loaded = false;
  controlKey: string;
  roomControlUrl: string;

  roomRef: RoomRef;

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

  getRoom = (): RoomRef => {
    const room = new BehaviorSubject<Room>(undefined);

    // use ws for http, wss for https
    let protocol = "ws:";
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    }

    const endpoint = protocol + "//" + window.location.host + "/ws";
    const ws = new WebSocket(endpoint);

    const roomRef = new RoomRef(room, ws, () => {
      console.log('closing room connection', room.value.id);
    });

    this.loaded = false;

    // handle incoming messages from bff
    ws.onmessage = msg => {
      const data = JSON.parse(msg.data);

      for (const k in data) {
        switch (k) {
          case "room":
            console.log("new room", data[k]);
            room.next(data[k]);
            this.loaded = true;

            break;

          case "mobileControl":
            console.log("mobile control info", data[k]);
            
            console.log(this.roomControlUrl, this.controlKey);
            break;
          default:
            console.warn(
              "got key '" + k + "', not sure how to handle that message"
            );
        }
      }
    };

    ws.onclose = event => {
      console.warn("websocket close", event);
      room.error(event);
    };

    this.roomRef = roomRef;
    return roomRef;
  };

  error = (msg: string) => {
    console.log("showing error", msg);
    const dialogs = this.dialog.openDialogs.filter(dialog => {
      // return dialog.componentInstance instanceof ErrorDialog;
    });

    if (dialogs.length > 0) {
      // do something?
      // either open another one on top of the current ones,
      // change the text on the current one?
    } else {
      // open a new dialog
      // const ref = this.dialog.open(ErrorDialog, {
      //   width: "80vw",
      //   data: {
      //     msg: msg
      //   }
      // });

      // ref.afterClosed().subscribe(result => {
      //   this.router.navigate([], {
      //     queryParams: { error: null },
      //     queryParamsHandling: "merge",
      //     preserveFragment: true
      //   });
      // });
    }
  };
}
