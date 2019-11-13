import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

export class DisplayDialogData {
  displayName: string;
  displayIcon: string;
}

@Component({
  selector: 'display-dialog',
  templateUrl: './display-dialog.component.html',
  styleUrls: ['./display-dialog.component.scss']
})
export class DisplayDialogComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<DisplayDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DisplayDialogData) { }

  ngOnInit() {
  }

  onCloseClick(): void {
    this.dialogRef.close();
  }

}
