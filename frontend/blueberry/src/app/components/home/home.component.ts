import { Component, OnInit, Input as AngularInput, ViewChild } from '@angular/core';
import { BFFService, RoomRef } from '../../services/bff.service';
import { ControlGroup, Input } from '../../../../../objects/control';
import { AudioComponent } from '../audio/audio.component';
import { ProjectorComponent } from '../projector/projector.component';
import { MatDialog } from '@angular/material';
import { HelpComponent } from 'src/app/dialogs/help/help.component';
import { SharingComponent } from 'src/app/dialogs/sharing/sharing.component';
import { MinionComponent } from 'src/app/dialogs/minion/minion.component';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss', '../../colorscheme.scss']
})
export class HomeComponent implements OnInit {
  @AngularInput() roomRef: RoomRef;
  cg: ControlGroup;
  mirrorMaster: Input;

  @ViewChild(AudioComponent, {static: false}) public audio: AudioComponent;
  @ViewChild(ProjectorComponent, {static: false}) public screen: ProjectorComponent;

  constructor(public bff: BFFService, public dialog: MatDialog) { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          if (!this.cg) {
            this.cg = r.controlGroups[r.selectedControlGroup];
          } else {
            this.applyChanges(r.controlGroups[r.selectedControlGroup])
          }
        }
      })
    }
  }

  applyChanges(tempCG: ControlGroup) {
    this.cg.displayGroups[0].shareInfo.state = tempCG.displayGroups[0].shareInfo.state;
    if (this.cg.displayGroups[0].shareInfo.state == 3) {
      this.dialog.open(MinionComponent, {
        width: "70vw"
      });
    }
  }
  
  public turnOff() {
    console.log("turning off the room");
    this.bff.locked = true;
    this.roomRef.setPower(false);
  }

  openHelp = () => {
    this.dialog.open(HelpComponent, {data: this.cg}).afterClosed().subscribe((helpMe) => {
      if (helpMe) {
        this.roomRef.requestHelp("");
      }
    });
  }

  openSharing = () => {
    this.dialog.open(SharingComponent, {data: this.roomRef})
  }

  testMinion = () => {
    this.dialog.open(MinionComponent, {
      width: "70vw"
    });
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
