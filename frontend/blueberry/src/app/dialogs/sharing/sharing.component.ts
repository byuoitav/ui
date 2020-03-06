import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { ControlGroup, Room } from '../../../../../objects/control';
import { RoomRef, BFFService } from 'src/app/services/bff.service';

@Component({
  selector: 'app-sharing',
  templateUrl: './sharing.component.html',
  styleUrls: ['./sharing.component.scss']
})
export class SharingComponent implements OnInit {
  chosenOptions: string[];
  cg: ControlGroup;
  sharingSpin = false;

  constructor(
    public ref: MatDialogRef<SharingComponent>,
    @Inject(MAT_DIALOG_DATA) public roomRef: RoomRef,
    public bff: BFFService) {
      this.cg = this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup];
      this.chosenOptions = [];
      // this.cg.displayGroups[0].shareOptions = ["Station 1", "Station 2", "Station 3"];
      for (let i = 0; i < this.cg.displayGroups[0].shareInfo.opts.length; i++) {
        let g = this.cg.displayGroups[0].shareInfo.opts[i];
        this.chosenOptions.push(g);
      }
      this.chosenOptions.sort();
      console.log(this.chosenOptions);
      this.bff.dialogCloser.subscribe((shouldClose) => {
        if (shouldClose === "sharing") {
          this.cancel();
        }
      });
    }

  ngOnInit() {
  }

  getSelectedInput() {
    if (this.cg) {
      const x = this.cg.inputs.find((i) => {
        return i.id === this.cg.displayGroups[0].input;
      });
      return x.name;
    } else {
      return "";
    }
  }

  cancel() {
    this.ref.close();
  }

  toggle(group: string) {
    if (this.chosenOptions.includes(group)) {
      this.chosenOptions.splice(this.chosenOptions.indexOf(group), 1);
    } else {
      this.chosenOptions.push(group);
      this.chosenOptions.sort();
    }
  }

  startShare = () => {
    this.roomRef.startSharing(this.cg.displayGroups[0].id, this.chosenOptions);
    this.sharingSpin = true;
  }
}
