import { Component, OnInit, AfterViewInit, Input as AngularInput, ViewChild } from '@angular/core';
import {MatTabGroup} from "@angular/material/tabs";
import {Http} from "@angular/http";

import { ControlGroup, Camera, Room } from '../../../../../objects/control';
import { RoomRef } from '../../../services/bff.service';

@Component({
  selector: 'camera-control',
  templateUrl: './camera-control.component.html',
  styleUrls: ['./camera-control.component.scss']
})
export class CameraControlComponent implements OnInit, AfterViewInit {
  @AngularInput()
  roomRef: RoomRef;


  @ViewChild("tabs", {static: true})
  tabs: MatTabGroup;
  code: string;
  cg: ControlGroup;

  constructor(private http: Http) {}

  ngOnInit() {
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
      }
    })
    console.log("cameras", this.cg.cameras);

    // this.getControlKey();
    // setInterval(() => {
    //   this.getControlKey();
    // }, 120000)
  }

  ngAfterViewInit() {
    // this is disgusting. :(
    // but, it moves the second line of tabs to be left aligned
    this.tabs._elementRef.nativeElement.getElementsByClassName(
      "mat-tab-labels"
    )[0].style.justifyContent = "flex-start";
  }

  tiltUp = (cam: Camera) => {
    console.log("tilting up");
    this.roomRef.tiltUp(cam.name)
  }

  tiltDown = (cam: Camera) => {
    console.log("tilting down");
    this.roomRef.tiltDown(cam.name)
  }

  panLeft = (cam: Camera) => {
    console.log("panning left");
    this.roomRef.panLeft(cam.name)
  }

  panRight = (cam: Camera) => {
    console.log("panning right");
    this.roomRef.panRight(cam.name)
  }

  panTiltStop = (cam: Camera) => {
    console.log("stopping tilt");
    this.roomRef.panTiltStop(cam.name)
  }

  zoomIn = (cam: Camera) => {
    console.log("zooming in");
    this.roomRef.zoomIn(cam.name)
  }

  zoomOut = (cam: Camera) => {
    console.log("zooming out");
    this.roomRef.zoomOut(cam.name)
  }

  zoomStop = (cam: Camera) => {
    console.log("stopping zoom");
    this.roomRef.zoomStop(cam.name)
  }

  selectPreset = (cam: Camera, preset: string) => {
    console.log("selecting preset", preset);
    this.roomRef.setPreset(cam.name, preset)
  }

  // getControlKey = () => {
  //   this.http
  //   .get(window.location.protocol + "//" + window.location.host +"/control-key/" + this.room + "/" + this.cg.name)
  //   .map(response => response.json()).subscribe(
  //     data => {
  //       this.code = data.ControlKey;
  //     }
  //   )
  // }
}
