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
  setInput = (displayID: string, sourceName: string, subSourceName?: string) => {
    const kv = {
      setInput: {
        displayGroup: displayID,
        source: sourceName,
        subSource: subSourceName
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  setBlank = (displayID: string, blanked: boolean) => {
    const kv = {
      setBlank: {
        displayGroup: displayID,
        blanked: blanked,
      }
    }

    console.log(kv)

    this.loading = true;
    this._ws.send(JSON.stringify(kv))
  }

  setVolume = (level: number, audioGroupName?: string, audioDeviceName?: string) => {
    const kv = {
      setVolume: {
        volume: level,
        audioGroup: audioGroupName,
        audioDevice: audioDeviceName
      }
    };

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  };

  setMuted = (muted: boolean, audioGroupName?: string, audioDeviceName?: string) => {
    const kv = {
      setMute: {
        mute: muted,
        audioGroup: audioGroupName,
        audioDevice: audioDeviceName
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

  setPower = (power: boolean, doAll: boolean) => {
    const kv = {
      setPower: {
        on: power,
        all: doAll
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

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  }

  lowerProjectorScreen = (screen: string) => {
    const kv = {
      lowerProjectorScreen: screen
    }

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  }

  stopProjectorScreen = (screen: string) => {
    const kv = {
      stopProjectorScreen: screen
    }

    this.loading = true;
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
        group: masterID,
        opts: optionsIDs
      }
    }

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  }

  stopSharing = (masterID: string) => {
    const kv = {
      setSharing: {
        group: masterID
      }
    }

    this.loading = true;
    this._ws.send(JSON.stringify(kv));
  }

  buttonPress = (key: string, value?: string) => {    
    const kv = {
      event: {
        key: key,
        value: value
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

  dialogCloser: EventEmitter<string>;
  retryEmitter: EventEmitter<any>;

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
    this.dialogCloser = new EventEmitter();
    this.retryEmitter = new EventEmitter();
  }

  getRoom = (): RoomRef => {
    const room = new BehaviorSubject<Room>(undefined);

    // use ws for http, wss for https
    let protocol = "ws:";
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    }

    const endpoint = protocol + "//" + window.location.host + "/api/v1/ws";
    const ws = new WebSocket(endpoint);

    const roomRef = new RoomRef(room, ws, () => {
      console.log('closing room connection', room.value.id);
    });


    this.loaded = false;

    // handle incoming messages from bff
    ws.onmessage = msg => {
      const data = JSON.parse(msg.data);
      console.log("data", data)

      for (const k in data) {
        switch (k) {
          case "room":
            console.log("new room", data[k]);
            room.next(data[k]);
            this.loaded = true;
            roomRef.loading = false;

            break;

          case "mobileControl":
            console.log("mobile control info", data[k]);
            
            console.log(this.roomControlUrl, this.controlKey);
            break;
          case "shareStarted":
            this.dialogCloser.emit("sharing");
            break;
          case "shareEnded":
            console.log("The sharing session has ended.");
            break;
          case "refresh":
            console.log("refreshing!");
            window.location.reload()
            break;
          case "becameInactive":
            this.dialogCloser.emit("inactive");
            break;
          default:
            console.warn(
              "got key '" + k + "', not sure how to handle that message"
            );
        }
      }
    };

    ws.onclose = event => {
      setTimeout(() => {
        console.warn("websocket close", event);
        this.retryEmitter.emit();
        room.error(event);
      }, 3000);
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
