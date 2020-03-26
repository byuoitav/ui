import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef } from 'src/app/services/bff.service';

@Component({
  selector: 'turnOffRoom-dialog',
  templateUrl: './turnOffRoom-dialog.component.html',
  styleUrls: ['./turnOffRoom-dialog.component.scss']
})
export class TurnOffRoomDialogComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<TurnOffRoomDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public roomRef: RoomRef) { }

  ngOnInit() {
  }

  cancel() {
    this.dialogRef.close();
  }

  leaveRoom(turnThingsOff: boolean) {
    if (turnThingsOff) {
      this.roomRef.setPower(false, true);
    } 
    this.roomRef.logout();
    this.dialogRef.close();
  }
}
