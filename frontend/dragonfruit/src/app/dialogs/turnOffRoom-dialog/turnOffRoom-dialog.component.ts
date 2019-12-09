import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

export class TurnOffRoomDialogData {

}

@Component({
  selector: 'turnOffRoom-dialog',
  templateUrl: './turnOffRoom-dialog.component.html',
  styleUrls: ['./turnOffRoom-dialog.component.scss']
})
export class TurnOffRoomDialogComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<TurnOffRoomDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: TurnOffRoomDialogData) { }

  ngOnInit() {
  }

  onCloseClick(): void {
    this.dialogRef.close();
  }

}
