import { Injectable, EventEmitter } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { BehaviorSubject, Observable, throwError } from "rxjs";
import { Router, ActivationEnd } from "@angular/router";
import { MatDialog } from "@angular/material";

import {
  Room,
  ControlGroup,
  DisplayGroup,
  Input,
  AudioDevice,
  AudioGroup,
  PresentGroup
} from "../../../objects/control";
// import { ErrorDialog } from "../dialogs/error/error.dialog";
// import { TurnOffRoomDialogComponent } from '../dialogs/turnOffRoom-dialog/turnOffRoom-dialog.component';

export class RoomRef {
  private _room: BehaviorSubject<Room>;
  private _ws: WebSocket;
  private _logout: () => void;
  loadingHome = false;
  loadingLock = false;
  commandInProgress: boolean;

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

    console.log(kv)

    this.commandInProgress = true;
    this._ws.send(JSON.stringify(kv));
  };

  setBlank = (displayID: string, blanked: boolean) => {
    const kv = {
      setBlank: {
        displayGroup: displayID,
        blanked: blanked,
      }
    };

    console.log(kv)

    this.commandInProgress = true;
    this._ws.send(JSON.stringify(kv))
  }

  setVolume = (level: number, audioGroupdID?: string, audioDeviceID?: string) => {
    const kv = {
      setVolume: {
        volume: level,
        audioGroup: audioGroupdID,
        audioDevice: audioDeviceID,
      }
    };

    console.log(kv)

    this.commandInProgress = true;
    this._ws.send(JSON.stringify(kv));
  };

  setMuted = (muted: boolean, audioGroupID?: string, audioDeviceID?: string) => {
    const kv = {
      setMute: {
        mute: muted,
        audioGroup: audioGroupID,
        audioDevice: audioDeviceID
      }
    };
    console.log(kv)
    this.commandInProgress = true;
    this._ws.send(JSON.stringify(kv));
  };

  setPower = (power: boolean, all?: string) => {
    const kv = {
      setPower: {
        on: power,
        all: all
      }
    };

    if (power == true) {
      this.loadingHome = true;
    } else {
      this.loadingLock = true;
    }
    console.log(kv)

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

  getControlInfo = (cgID: string) => {
    const kv = {
      getControlKey: {
        controlGroupID: cgID
      }
    }

    this._ws.send(JSON.stringify(kv));
  }

  tiltUp = (cam: string) => {
    const kv = {
      tiltUp: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  tiltDown = (cam: string) => {
    const kv = {
      tiltDown: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  panLeft = (cam: string) => {
    const kv = {
      panLeft: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  panRight = (cam: string) => {
    const kv = {
      panRight: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  panTiltStop = (cam: string) => {
    const kv = {
      panTiltStop: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  zoomIn = (cam: string) => {
    const kv = {
      zoomIn: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  zoomOut = (cam: string) => {
    const kv = {
      zoomOut: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  zoomStop = (cam: string) => {
    const kv = {
      zoomStop: {
        camera: cam
      }
    }

    this._ws.send(JSON.stringify(kv))
  }

  setPreset = (cam: string, preset: string) => {
    const kv = {
      setPreset: {
        camera: cam,
        preset: preset
      }
    }

    this._ws.send(JSON.stringify(kv))
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
  closeEmitter: EventEmitter<string>;
  
  roomRef: RoomRef;
  
  constructor(private router: Router, private dialog: MatDialog) {
    // do things based on route changes
    this.closeEmitter = new EventEmitter();
    this.router.events.subscribe(event => {
      if (event instanceof ActivationEnd) {
        const snapshot = event.snapshot;

        // show error if the "error" query paramater is present
        // if (snapshot && snapshot.queryParams && snapshot.queryParams.error) {
        //   this.error(snapshot.queryParams.error);
        // }
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

    const endpoint = protocol + "//" + window.location.host + "/api/v1/ws";
    const ws = new WebSocket(endpoint);

    this.roomRef = new RoomRef(room, ws, () => {
      console.log("closing room connection", room.value.id);
    });


    // handle incoming messages from bff
    ws.onmessage = msg => {
      const data = JSON.parse(msg.data);

      for (const k in data) {
        switch (k) {
          case "room":
            console.log("new room", data[k]);

            room.next(data[k]);
            if (this.roomRef.loadingHome == true && data[k].controlGroups[data[k].selectedControlGroup].poweredOn == true) {
              this.roomRef.loadingHome = false;
            }
            if (this.roomRef.loadingLock == true && data[k].controlGroups[data[k].selectedControlGroup].poweredOn == false) {
              this.roomRef.loadingLock = false;
            }
            this.roomRef.commandInProgress = false;
            break;

          case "refresh":
            console.log("refreshing!");
            window.location.reload()
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
      this.closeEmitter.emit("conn closed");
      room.error(event);
      // this.roomRef = this.getRoom();
    };

    return this.roomRef;
  };

  error = (msg: string) => {
    console.log("showing error", msg);
    // const dialogs = this.dialog.openDialogs.filter(dialog => {
    //   return dialog.componentInstance instanceof ErrorDialog;
    // });

  //   if (dialogs.length > 0) {
  //     // do something?
  //     // either open another one on top of the current ones,
  //     // change the text on the current one?
  //   } else {
  //     // open a new dialog
  //     // const ref = this.dialog.open(ErrorDialog, {
  //     //   width: "80vw",
  //     //   data: {
  //     //     msg: msg
  //     //   }
  //     // });

  //     ref.afterClosed().subscribe(result => {
  //       this.router.navigate([], {
  //         queryParams: { error: null },
  //         queryParamsHandling: "merge",
  //         preserveFragment: true
  //       });
  //     });
  //   }
  // };
  }
}
