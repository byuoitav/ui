import { Component, OnInit, Input as AngularInput, ViewChild } from '@angular/core';
import { BFFService, RoomRef } from '../../services/bff.service';
import { ControlGroup, Input } from '../../../../../objects/control';
import { AudioComponent } from '../audio/audio.component';
import { ProjectorComponent } from '../projector/projector.component';
import { MatDialog, MatDialogRef } from '@angular/material';
import { HelpComponent } from 'src/app/dialogs/help/help.component';
import { SharingComponent } from 'src/app/dialogs/sharing/sharing.component';
import { MinionComponent } from 'src/app/dialogs/minion/minion.component';
import { MobileComponent } from 'src/app/dialogs/mobile/mobile.component';
import { PowerOffComponent } from 'src/app/dialogs/power-off/power-off.component';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss', '../../colorscheme.scss']
})
export class HomeComponent implements OnInit {
  @AngularInput() roomRef: RoomRef;
  cg: ControlGroup;
  mirrorMaster: Input;
  minionRef: MatDialogRef<MinionComponent>;
  sharingRef: MatDialogRef<SharingComponent>;

  @ViewChild(AudioComponent, {static: false}) public audio: AudioComponent;
  @ViewChild(ProjectorComponent, {static: false}) public screen: ProjectorComponent;

  constructor(public bff: BFFService, public dialog: MatDialog) { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          if (!this.cg) {
            this.cg = r.controlGroups[r.selectedControlGroup];
            if (this.cg.displayGroups[0].shareInfo.state == 3 && !this.dialog.openDialogs.includes(this.minionRef)) {
              this.minionRef = this.dialog.open(MinionComponent, {
                width: "70vw",
                data: {
                  roomRef: this.roomRef
                },
                disableClose: true
              });
            }
          } else {
            this.applyChanges(r.controlGroups[r.selectedControlGroup])
          }
        }
      })
    }
  }

  applyChanges(tempCG: ControlGroup) {
    if (this.cg.poweredOn == true && tempCG.poweredOn == false) {
      this.dialog.closeAll();
    }
    this.cg.displayGroups[0].shareInfo.state = tempCG.displayGroups[0].shareInfo.state;
    if (this.cg.displayGroups[0].shareInfo.state == 3 && !this.dialog.openDialogs.includes(this.minionRef)) {
      this.minionRef = this.dialog.open(MinionComponent, {
        width: "70vw",
        data: {
          roomRef: this.roomRef
        },
        disableClose: true
      });
    }
    if (this.cg.displayGroups[0].shareInfo.state == 1 && this.dialog.openDialogs.includes(this.minionRef)) {
      this.minionRef.close();
    }
    if (this.cg.displayGroups[0].shareInfo.state != 1 && this.dialog.openDialogs.includes(this.sharingRef)) {
      this.sharingRef.close();
    }
  }
  
  public turnOff() {
    if (this.roomRef.room) {
      var size = Object.keys(this.roomRef.room.controlGroups).length;
      if (size > 1) {
        this.dialog.open(PowerOffComponent, {data: this.roomRef, disableClose: true}).afterClosed().subscribe((turnAllOff) => {
          if (turnAllOff !== "cancel") {
            console.log("turning off the room");
            if (turnAllOff === "one") {
              this.roomRef.setPower(false, false);
            } else if (turnAllOff === "all") {
              this.roomRef.setPower(false, true);
            }
            this.bff.locked = true;
          }
        });
      } else {
        this.roomRef.setPower(false, true);
        this.bff.locked = true;
      }
    }
  }

  openHelp = () => {
    this.dialog.open(HelpComponent, {data: this.roomRef, disableClose: true});
  }

  openSharing = () => {
    this.sharingRef = this.dialog.open(SharingComponent, {data: this.roomRef});
  }

  openMobileControl = () => {
    this.dialog.open(MobileComponent, {data: this.cg});
  }

  canShare = () => {
    if (this.cg) {
      return this.cg.displayGroups[0].shareInfo.state === 1;
    }
  }

  currentlySharing = () => {
    if (this.cg) {
      return this.cg.displayGroups[0].shareInfo.state === 2;
    }
  }

  stopSharing = () => {
    if (this.roomRef) {
      this.roomRef.stopSharing(this.cg.displayGroups[0].id);
    }
  }

  hasScreens = () => {
    if (this.cg) {
      return this.cg.screens && this.cg.screens.length > 0
    }
  }
}
