import { Component, OnInit, AfterViewInit, Input as AngularInput, ViewChild } from '@angular/core';
import {MatTabGroup} from "@angular/material/tabs";
import {Http} from "@angular/http";

import { ControlGroup, Camera, CameraPreset, Room } from '../../../../../objects/control';
import { RoomRef } from '../../../services/bff.service';

@Component({
  selector: 'camera-control',
  templateUrl: './camera-control.component.html',
  styleUrls: ['./camera-control.component.scss']
})
export class CameraControlComponent implements OnInit, AfterViewInit {
  @AngularInput()
  roomRef: RoomRef;


  @ViewChild(MatTabGroup)
  // @ViewChild(MatTabGroup, null)

  private _tabs: MatTabGroup;
  code: string;
  cg: ControlGroup;
  room: Room
  // I think we can get this from a roomRef?
  // room = APIService.building + "-" + APIService.roomName;

  constructor(private http: Http) {}

  ngOnInit() {
    console.log("cameras", this.cg.cameras);
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.room = r
        this.cg = r.controlGroups[r.selectedControlGroup];
      }
    })
    this.getControlKey();
    setInterval(() => {
      this.getControlKey();
    }, 120000)
  }

  ngAfterViewInit() {
    // this is disgusting. :(
    // but, it moves the second line of tabs to be left aligned
    this._tabs._elementRef.nativeElement.getElementsByClassName(
      "mat-tab-labels"
    )[0].style.justifyContent = "flex-start";
  }

  tiltUp = (cam: Camera) => {
    console.log("tilting up", cam.tiltUp);
    if (!cam.tiltUp) {
      return;
    }

    this.http.get(cam.tiltUp).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  tiltDown = (cam: Camera) => {
    console.log("tilting down", cam.tiltDown);
    if (!cam.tiltDown) {
      return;
    }

    this.http.get(cam.tiltDown).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panLeft = (cam: Camera) => {
    console.log("panning left", cam.panLeft);
    if (!cam.panLeft) {
      return;
    }

    this.http.get(cam.panLeft).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panRight = (cam: Camera) => {
    console.log("panning right", cam.panRight);
    if (!cam.panRight) {
      return;
    }

    this.http.get(cam.panRight).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  panTiltStop = (cam: Camera) => {
    console.log("stopping pan", cam.panTiltStop);
    if (!cam.panTiltStop) {
      return;
    }

    this.http.get(cam.panTiltStop).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  zoomIn = (cam: Camera) => {
    console.log("zooming in", cam.zoomIn);
    if (!cam.zoomIn) {
      return;
    }

    this.http.get(cam.zoomIn).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  zoomOut = (cam: Camera) => {
    console.log("zooming out", cam.zoomOut);
    if (!cam.zoomOut) {
      return;
    }

    this.http.get(cam.zoomOut).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  zoomStop = (cam: Camera) => {
    console.log("stopping zoom", cam.zoomStop);
    if (!cam.zoomStop) {
      return;
    }

    this.http.get(cam.zoomStop).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  selectPreset = (preset: CameraPreset) => {
    console.log("selecting preset", preset.displayName, preset.setPreset);
    if (!preset.setPreset) {
      return;
    }

    this.http.get(preset.setPreset).subscribe(resp => {
      console.log("resp", resp);
    }, err => {
      console.warn("err", err);
    });
  }

  getControlKey = () => {
    this.http
    .get(window.location.protocol + "//" + window.location.host +"/control-key/" + this.room + "/" + this.cg.name)
    .map(response => response.json()).subscribe(
      data => {
        this.code = data.ControlKey;
      }
    )
  }
}
