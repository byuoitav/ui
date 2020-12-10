import { Component, OnInit, Input as AngularInput, Output, ɵConsole } from '@angular/core';
import { MobileControlComponent } from "../../dialogs/mobilecontrol/mobilecontrol.component";
import { MatDialog } from "@angular/material";
import { RoomRef, BFFService } from '../../../services/bff.service';
import { ControlGroup, DisplayGroup, Input, Room } from '../../../../../objects/control';



@Component({
  selector: 'display',
  templateUrl: './display.component.html',
  styleUrls: ['./display.component.scss']
})
export class DisplayComponent implements OnInit {

  @AngularInput()
  roomRef: RoomRef
  cg: ControlGroup;
  selectedOutput: number;
  selectedInput: Input;
  // blanked: boolean;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
        if (this.cg.displayGroups.length > 0) {
          if (this.selectedOutput == undefined) {
            this.selectedOutput = 0;
          }
          this.selectedInput = this.cg.displayGroups[this.selectedOutput].inputs.find((input) => input.name === this.cg.displayGroups[this.selectedOutput].input)
        }
      }
    })
  }

  public changeInput(display: DisplayGroup, input: Input) {
    if (display.input != input.name) {
      document.getElementById("input" + input.name).classList.toggle("feedback");
      this.roomRef.setInput(display.name, input.name);
    }

    if (display.blanked) {
      document.getElementById("input" + input.name).classList.toggle("feedback");
      this.setBlank(display, false);
    }
  }

  public openMobileControlDialog() {
    console.log(this.cg.controlInfo.url);
    const dialogRef = this.dialog.open(MobileControlComponent, {
      width: "70vw",
      data: {
        url: this.cg.controlInfo.url,
        key: this.cg.controlInfo.key
      }
    });
  }

  public getInputForOutput(d: DisplayGroup) {
    this.selectedInput = d.inputs.find((i) => i.name === d.input)
    console.log("input", this.selectedInput)
  }

  public getInputIcon(d: DisplayGroup) {
    const input = d.inputs.find((i) => i.name === d.input);
    if (input == undefined || d.blanked) {
      return "crop_landscape";
    }
    return input.icon;
  }

  public getInputName(d: DisplayGroup) {
    if (d.blanked) {
      return "Blank"
    }
    const input = d.inputs.find((i) => i.name === d.input);
    if (input == undefined) {
      return "unknown";
    }
    return input.name;
  }

  public setBlank(d: DisplayGroup, blanked: boolean) {
    if (blanked) {
      document.getElementById("blank").classList.toggle("feedback");
      this.roomRef.setBlank(d.name, blanked);
      setTimeout(() => {
        console.log(this.cg.displayGroups[this.selectedOutput].blanked)
      }, 2000);
    } else {
      this.roomRef.setBlank(d.name, blanked)
    }
  }
}
