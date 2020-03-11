import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef, BFFService } from 'src/app/services/bff.service';
import { ControlGroup, DisplayGroup, Input } from '../../../../../objects/control';

@Component({
  selector: 'app-sharing',
  templateUrl: './sharing.component.html',
  styleUrls: ['./sharing.component.scss']
})
export class SharingComponent implements OnInit {
  roomRef: RoomRef;
  cg: ControlGroup;
  selectedDisplay: DisplayGroup;
  chosenOptions: string[] = [];
  sharingSpin = false;


  constructor(public dialogRef: MatDialogRef<SharingComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any, public bff: BFFService) {
      this.roomRef = this.data.roomRef as RoomRef;
      this.cg = this.roomRef.room.controlGroups[this.roomRef.room.selectedControlGroup];
      this.selectedDisplay = this.data.display;

      this.bff.dialogCloser.subscribe((shouldClose) => {
        if (shouldClose === "sharing") {
          this.dialogRef.close();
        }
      });
     }

  ngOnInit() {
  }

  get selectedInput():Input {
    if (this.cg && this.selectedDisplay) {
      return this.cg.inputs.find((i) => {
        return this.selectedDisplay.input === i.id
      });
    }
  }

  trimOption(optionID: string): string {
    if (optionID.includes("-")) {
      let pieces = optionID.split("-");
      return pieces[2];
    }
    return optionID;
  }

  toggleDisplaySelection = (option: string) => {
    if (this.chosenOptions.includes(option)) {
      console.log("removing...");
      this.chosenOptions.splice(this.chosenOptions.indexOf(option), 1);
    } else {
      console.log("adding...");
      this.chosenOptions.push(option);
    }
  }

  startSharing() {
    console.log(this.chosenOptions);
    this.roomRef.startSharing(this.selectedDisplay.id, this.chosenOptions);
    this.sharingSpin = true;
    this.dialogRef.disableClose = true;
  }
}
