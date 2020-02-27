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
  blanked: Input;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.blanked = {
      id: "blank",
      icon: "crop_landscape",
      name: "Blank",
      subInputs: null,
      disabled: false
    }
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
        if (this.cg.displayGroups.length > 0) {
          if (this.selectedOutput == undefined) {
            this.selectedOutput = 0;
          }
          if (this.cg.displayGroups[this.selectedOutput].blanked == true) {
            this.selectedInput = this.blanked;
          } else {
            this.selectedInput = this.cg.inputs.find((i) => i.id === this.cg.displayGroups[this.selectedOutput].input)
          }
        }
      }
    })
  }

  public changeInput(display: string, input: Input) {
    document.getElementById("input" + input.id).classList.toggle("feedback");
    this.roomRef.setInput(display, input.id);
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
      this.selectedInput = this.blanked;
    } else {
      this.selectedInput = this.cg.inputs.find((i) => i.id === d.input)
      if (this.selectedInput == undefined) {
        this.selectedInput = this.blanked;
      }
    }
  }

  public toggleBlank(d: DisplayGroup) {
    if (d.blanked == true) {
      this.roomRef.setBlanked(d.id, false);
    } else {
      this.roomRef.setBlanked(d.id, true);
    }
  }

  public getInputIcon(d: DisplayGroup) {
    if (d.blanked == true) {
      return this.blanked.icon;
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
      return this.blanked.name;
    } else {
      const input = this.cg.inputs.find((i) => i.id === d.input);
      if (input == undefined) {
        return "unknown";
      }
      return input.name;
    }
  }
}
