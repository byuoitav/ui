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
  blank: Input;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.blank = {
      id: "blank",
      icon: "crop_landscape",
      name: "Blank",
      subInputs: null,
    }
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
        if (this.cg.displayGroups.length > 0) {
          if (this.selectedOutput == undefined) {
            this.selectedOutput = 0;
          }
          if (this.cg.displayGroups[this.selectedOutput].blanked == true) {
            this.selectedInput = this.blank;
            
            // for some reason the spinny edge of blank doesn't work the same as the inputs
            // so we need to do this...
            let btn = document.getElementById("input" + this.blank.id);
            btn.classList.remove("feedback")
          } else {
            this.selectedInput = this.cg.inputs.find((i) => i.id === this.cg.displayGroups[this.selectedOutput].input)
          }
        }
      }
    })
  }

  public changeInput(display: DisplayGroup, input: Input) {
    if (display.input != input.id) {
      document.getElementById("input" + input.id).classList.toggle("feedback");
      this.roomRef.setInput(display.id, input.id);
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
    if (d.blanked == true) {
      this.selectedInput = this.blank;
    } else {
      this.selectedInput = this.cg.inputs.find((i) => i.id === d.input)
      if (this.selectedInput == undefined) {
        this.selectedInput = this.blank;
      }
    }
  }

  public setBlank(d: DisplayGroup) {
    if (!d.blanked) {
      document.getElementById("input" + this.blank.id).classList.toggle("feedback");
      this.roomRef.setBlanked(d.id, true);
    }
  }

  public getInputIcon(d: DisplayGroup) {
    if (d.blanked == true) {
      return this.blank.icon;
    } else {
      const input = this.cg.inputs.find((i) => i.id === d.input);
      if (input == undefined) {
        return "crop_landscape";
      }
      return input.icon;
    }
  }

  public getInputName(d: DisplayGroup) {
    if (d.blanked == true) {
      return this.blank.name;
    } else {
      const input = this.cg.inputs.find((i) => i.id === d.input);
      if (input == undefined) {
        return "unknown";
      }
      return input.name;
    }
  }
}
