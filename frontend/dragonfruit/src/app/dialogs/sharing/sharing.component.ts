import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef } from 'src/app/services/bff.service';
import { ControlGroup, DisplayGroup } from '../../../../../objects/control';

@Component({
  selector: 'app-sharing',
  templateUrl: './sharing.component.html',
  styleUrls: ['./sharing.component.scss']
})
export class SharingComponent implements OnInit {
  roomRef: RoomRef;
  cg: ControlGroup;
  selectedDisplay: DisplayGroup;


  constructor(public dialogRef: MatDialogRef<SharingComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any) {
      this.roomRef = this.data.roomRef as RoomRef;
      this.cg = this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup];
      this.selectedDisplay = this.data.display;
     }

  ngOnInit() {
  }

}
