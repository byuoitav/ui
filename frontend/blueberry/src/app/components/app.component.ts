import { Component, ViewChild } from '@angular/core';
import { BFFService, RoomRef } from '../services/bff.service';
import { AudioComponent } from './audio/audio.component';
import { ProjectorComponent } from './projector/projector.component';
import { MobileComponent } from '../dialogs/mobile/mobile.component';
import { MatDialog } from '@angular/material';
import { ControlGroup } from '../../../../objects/control';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  public roomRef: RoomRef;
  public cg: ControlGroup;

  @ViewChild(AudioComponent, {static: false}) public audio: AudioComponent;
  @ViewChild(ProjectorComponent, {static: false}) public screen: ProjectorComponent;

  constructor(public bff: BFFService, public dialog: MatDialog) {
    this.roomRef = this.bff.getRoom();
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
      }
    })
  }

  unlock = () => {
    this.bff.locked = false;
    this.roomRef.setPower(true);
  }

  powerIsOff = ():boolean => {
    if (this.roomRef && this.roomRef.room) {
      if (this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup].power == "standby") {
        return true;
      } else {
        // this.unlock()
        return false;
      }
    }
  }

  hasScreens() {
    // return true;
    if (this.roomRef && this.roomRef.room) {
      if (this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup].screens) {
        return this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup].screens.length > 0;
      }
    }
    return false;
  }

  haveControlKey() {
    // TODO: do this thing
    if (this.bff.controlKey) {
      return true;
    }
  }

  hasAudioGroups() {
    // return true;
    if (this.roomRef && this.roomRef.room) {
      if (this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup].audioGroups) {
        return this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup].audioGroups.length > 1;
      }
    }
  }

  showScreenControl() {
    if (this.roomRef && this.roomRef.room) {
      this.screen.show(this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup]);
    }
  }

  showAudioControl() {
    if (this.roomRef && this.roomRef.room) {
      this.audio.show(this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup]);
    }
  }

  showMobileControl() {
    if (this.roomRef && this.roomRef.room) {
      this.dialog.open(MobileComponent, {data: this.roomRef});
    }
  }
}
