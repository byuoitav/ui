import { Component, OnInit, Input as AngularInput, Output } from '@angular/core';
import { MobileControlComponent } from "../../dialogs/mobilecontrol/mobilecontrol.component";
import { MatDialog } from "@angular/material";
import { RoomRef, BFFService } from '../../../services/bff.service';
import { ControlGroup, DisplayBlock, Input, Room } from '../../../../../objects/control';



@Component({
  selector: 'display',
  templateUrl: './display.component.html',
  styleUrls: ['./display.component.scss']
})
export class DisplayComponent implements OnInit {

  @AngularInput()
  roomRef: RoomRef
  cg: ControlGroup;
  selectedOutput: DisplayBlock;
  selectedInput: Input;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
        if (this.cg.displayBlocks.length > 0) {
          this.selectedOutput = this.cg.displayBlocks[0];
          // Danny says I don't have to worry about the blanked stuff and change input will eventually do that
          // if (this.selectedOutput.blanked == true) {
          //   this.selectedInput = this.cg.inputs[0];
          // } else {
            for ( let input of this.cg.inputs) {
              if (this.selectedOutput.input == input.id) {
                this.selectedInput = input;
                break;
              }
            }
          // }
        }
      }
    })
  }

  public changeInput(display: string, input: Input) {
    document.getElementById("input" + input.id).classList.toggle("feedback");
    this.roomRef.setInput(display, input.id);
  }

  public openMobileControlDialog() {
    const dialogRef = this.dialog.open(MobileControlComponent, {
      width: "70vw"
    });
  }

}
