import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { BFFService } from 'src/app/services/bff.service';
import { ControlGroup } from '../../../../../objects/control';

@Component({
  selector: 'app-mirror',
  templateUrl: './mirror.component.html',
  styleUrls: ['./mirror.component.scss']
})
export class MirrorComponent implements OnInit {
  public cg: ControlGroup;

  constructor(
    public dialogRef: MatDialogRef<MirrorComponent>, 
    @Inject(MAT_DIALOG_DATA) public data: any, 
    public bff: BFFService) {
      this.bff.dialogCloser.subscribe((closeEvent) => {
        if (closeEvent === "inactive") {
          this.dialogRef.close();
        }
      })
      
      this.cg = this.data.roomRef.room.controlGroups[this.data.roomRef.room.selectedControlGroup];
    }

  ngOnInit() {
  }

  cancel = () => {
    this.data.roomRef.stopSharing(this.cg.displayGroups[0].id);
    this.dialogRef.close();
  }
}
