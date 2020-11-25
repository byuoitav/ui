import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef, BFFService } from 'src/app/services/bff.service';
import { ControlGroup } from '../../../../../objects/control';


@Component({
  selector: 'app-minion',
  templateUrl: './minion.component.html',
  styleUrls: ['./minion.component.scss']
})
export class MinionComponent implements OnInit {
  public cg: ControlGroup;

  constructor(
    public ref: MatDialogRef<MinionComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {
      roomRef: RoomRef;
    },
    public bff: BFFService) {
      this.bff.dialogCloser.subscribe((closeEvent) => {
        this.ref.close();
      })
     }

  ngOnInit() {
    this.cg = this.data.roomRef.room.controlGroups[this.data.roomRef.room.selectedControlGroup];
    console.log(this.cg);
  }

  cancel = () => {
    this.data.roomRef.stopSharing(this.cg.displayGroups[0].name);
  }

}
