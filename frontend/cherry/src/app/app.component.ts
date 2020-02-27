import { Component, ViewEncapsulation, OnInit, ViewChild } from "@angular/core";
import { MatDialog } from '@angular/material';
import { trigger, transition, animate } from "@angular/animations";
import { Http } from "@angular/http";
import { Output } from '../objects/status.objects';
import { BFFService, RoomRef } from '../services/bff.service';
import { HelpDialog } from "./dialogs/help.dialog";
import { MobileControlComponent } from "./dialogs/mobilecontrol/mobilecontrol.component";
import { ControlGroup } from "../../../objects/control";
import { LockScreenAudioComponent } from "./components/lockscreenaudio/lockscreenaudio.component";
import { LockScreenScreenControlComponent } from "./components/lockscreenscreencontrol/lockscreenscreencontrol.component";


const HIDDEN = "hidden";
const QUERY = "query";
const LOADING = "indeterminate";
const BUFFER = "buffer";

@Component({
  selector: "cherry",
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [
    trigger("delay", [
      transition(":enter", [animate(500)]),
      transition(":leave", [animate(500)])
    ])
  ]
})
export class AppComponent implements OnInit {
  public power: boolean;
  public roomRef: RoomRef;
  public cg: ControlGroup;

  @ViewChild(LockScreenAudioComponent, {static: true})
  public lockAudio: LockScreenAudioComponent;

  @ViewChild(LockScreenScreenControlComponent, {static: true})
  public screen: LockScreenScreenControlComponent;

  constructor(
    public bff: BFFService,
    public dialog: MatDialog
  ) {
    this.roomRef = this.bff.getRoom();
    this.power = false;
    console.log(this.bff);
    console.log(this.roomRef);
  }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          this.cg = r.controlGroups[r.selectedControlGroup];
          // if (this.cg.poweredOn == true) {
          //   this.roomRef.loading = false;
          // }
          // else {
          //   this.roomRef.loading = false;
          // }
        }
      })
    }
  }

  public openHelpDialog() {
    const dialogRef = this.dialog.open(HelpDialog, {
      data: this.roomRef,
      width: "70vw",
      disableClose: true
    });
  }

  public openMobileControlDialog() {
    const dialogRef = this.dialog.open(MobileControlComponent, {
      width: "70vw",
      height: "52.5vw"
    });
  }

  public togglePower() {

    if (this.cg.poweredOn == true) {
      //probably have to do a check to see if all the displays should turn off
      this.roomRef.setPower(false);
    } else {
      this.roomRef.setPower(true);
    }
  }

  public showManagement() {
    if (this.screen.isShowing() || this.lockAudio.isShowing()) {
      return false;
    }

    if (this.cg.poweredOn == true) {
      return false;
    }
    return true;
  }
}
