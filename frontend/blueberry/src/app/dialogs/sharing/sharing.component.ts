import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { ControlGroup, Room } from 'src/app/objects/control';

@Component({
  selector: 'app-sharing',
  templateUrl: './sharing.component.html',
  styleUrls: ['./sharing.component.scss']
})
export class SharingComponent implements OnInit {
  // cg: ControlGroup
  constructor(
    public ref: MatDialogRef<SharingComponent>,
    @Inject(MAT_DIALOG_DATA) public cg: ControlGroup
    ) {
      this.cg.displays[0].shareOptions = ["Station 1", "Station 2", "Station 3"];
    }

  ngOnInit() {
  }

  getSelectedInput() {
    if (this.cg) {
      const x = this.cg.inputs.find((i) => {
        return i.id === this.cg.displays[0].input;
      });
      return x.name;
    } else {
      return "";
    }
  }

  cancel() {
    this.ref.close();
  }
}
