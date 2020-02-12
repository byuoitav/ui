import { Component, OnInit, Input as AngularInput, Output } from '@angular/core';
import { MobileControlComponent } from "../../dialogs/mobilecontrol/mobilecontrol.component";
import { MatDialog } from "@angular/material";
import { RoomRef } from '../../../services/bff.service';
import { ControlGroup, Display, Input } from '../../../objects/control';



@Component({
  selector: 'display',
  templateUrl: './display.component.html',
  styleUrls: ['./display.component.scss']
})
export class DisplayComponent implements OnInit {

  @AngularInput()
  roomRef: RoomRef;
  cg: ControlGroup;
  selectedOutput: Display;
  selectedInput: Input;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          if (!this.cg) {
            this.cg = r.controlGroups[r.selectedControlGroup];
            if (this.cg.displays.length > 0) {
              this.selectedOutput = this.cg.displays[0];
            }
          }
        }
      })
    }
  }

  public changeInput(display: string, input: Input) {
    this.roomRef.setInput(display, input.id);
    this.selectedInput = input;
    this.roomRef.setVolume(this.selectedOutput.id, 10)
  }

  public openMobileControlDialog() {
    const dialogRef = this.dialog.open(MobileControlComponent, {
      width: "70vw"
    });
  }

}
