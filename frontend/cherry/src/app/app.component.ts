import { Component, ViewEncapsulation, OnInit, ViewChild } from "@angular/core";
import { MatDialog, MatDialogRef, SELECT_PANEL_INDENT_PADDING_X } from '@angular/material';
import { trigger, transition, animate } from "@angular/animations";
import { Http } from "@angular/http";
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
  public helpRef: MatDialogRef<HelpDialog>
  public mobileRef: MatDialogRef<MobileControlComponent>

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
          if (this.cg) {
            if ((this.cg.poweredOn == false && r.controlGroups[r.selectedControlGroup].poweredOn == true)
            || (this.cg.poweredOn == true && r.controlGroups[r.selectedControlGroup].poweredOn == false)) {
            if (this.dialog.openDialogs.includes(this.mobileRef)) {
              this.mobileRef.close();
            }
          }
          }
          this.cg = r.controlGroups[r.selectedControlGroup];
        }
      })
    }
    if (this.bff) {
      this.bff.closeEmitter.subscribe((e) => {
        if (e) {
          this.roomRef = this.bff.getRoom();
          if (this.roomRef) {
            this.roomRef.subject().subscribe((r) => {
              if (r) {
                console.log(this.dialog.openDialogs.includes(this.mobileRef));
                if ((this.cg.poweredOn == false && r.controlGroups[r.selectedControlGroup].poweredOn == true)
                  || (this.cg.poweredOn == true && r.controlGroups[r.selectedControlGroup].poweredOn == false)) {
                  if (this.dialog.openDialogs.includes(this.mobileRef)) {
                    this.dialog.closeAll();
                  }
                }
                this.cg = r.controlGroups[r.selectedControlGroup];
              }
            })
          }
          console.log("emitted!");
        }
      })
    }
  }

  public openHelpDialog() {
    this.helpRef = this.dialog.open(HelpDialog, {
      data: this.roomRef,
      width: "70vw",
      disableClose: true
    });
  }

  public openMobileControlDialog() {
    this.mobileRef = this.dialog.open(MobileControlComponent, {
      width: "70vw",
      data: {
        url: this.cg.controlInfo.url,
        key: this.cg.controlInfo.key
      }
    });
  }

  public togglePower() {
    if (this.cg.poweredOn == true) {
      this.roomRef.setPower(false);
    } else {
      this.roomRef.setPower(true);
    }
  }

  public showManagement() {
    if (!this.cg) {
      return true;
    }
    if (this.dialog.openDialogs.includes(this.helpRef)) {
      return true;
    }
    if (this.screen.isShowing() || this.lockAudio.isShowing()) {
      return false;
    }

    if (this.cg.poweredOn == true) {
      return false;
    }
    return true;
  }
}
