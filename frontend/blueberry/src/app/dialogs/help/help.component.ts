import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog } from '@angular/material';
import { ControlGroup } from '../../../../../objects/control';
import { RoomRef } from 'src/app/services/bff.service';
import { ConfirmHelpDialog } from './confirmhelp';

@Component({
  selector: 'app-help',
  templateUrl: './help.component.html',
  styleUrls: ['./help.component.scss']
})
export class HelpComponent implements OnInit {
  public cg: ControlGroup

  constructor(
    public ref: MatDialogRef<HelpComponent>,
    @Inject(MAT_DIALOG_DATA) public roomRef: RoomRef,
    public dialog: MatDialog
  ) {
    if (this.roomRef && this.roomRef.room) {
      this.cg = this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup];
    }
  }

  ngOnInit() {
  }

  public cancel() {
    this.ref.close();
}

public requestHelp() {
  const dialogRef = this.dialog.open(ConfirmHelpDialog, {
    width: "70vw",
    disableClose: true,
    data: this.roomRef,
  });
  this.cancel();
}

public isAfterHours(): boolean {
    let date = new Date();
    let DayOfTheWeek = date.getDay();
    let CurrentHour = date.getHours();

    switch(DayOfTheWeek) {
      // Sunday
      case 0: { return true; }
      // Monday
      case 1: {
        if(CurrentHour < 7 || CurrentHour >= 19) { return true; }
        else { return false; }
      }
      // Tuesday
      case 2: {
        if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
        else { return false; }
      }
      // Wednesday
      case 3: {
        if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
        else { return false; }
      }
      // Thursday
      case 4: {
        if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
        else { return false; }
      }
      // Friday
      case 5: {
        if(CurrentHour < 7 || CurrentHour >= 20) { return true; }
        else { return false; }
      }
      // Saturday
      case 6: {
        if(CurrentHour < 8 || CurrentHour >= 12) { return true; }
        else { return false; }
      }
      default: { return false; }
    }
  }
}
