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
  curDisplayGroup: DisplayGroup;
  constructor(
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
        // console.log("cg", this.cg)
        if (this.cg.displayGroups.length > 0) {
          if (this.selectedOutput == undefined) {
            this.selectedOutput = 0;
          }
          this.curDisplayGroup = this.cg.displayGroups[this.selectedOutput]
          // for (let i = 0; i < this.cg.displayGroups.length; i++) {
            this.selectedInput = this.cg.displayGroups[this.selectedOutput].inputs.find((input) => input.name === this.cg.displayGroups[this.selectedOutput].input)
            // this.selectedInput = this.cg.inputs.find((i) => i.id === this.cg.displayGroups[this.selectedOutput].input)
            // if (this.selectedInput != undefined) {
            //   break
            // }
          // } 
          this.selectedInput = this.cg.displayGroups[this.selectedOutput].inputs[0]
          // console.log("selected", this.selectedInput)
        }
      }
    })
  }

  public changeInput(display: DisplayGroup, input: Input) {
    if (display.input != input.name) {
      document.getElementById("input" + input.name).classList.toggle("feedback");
      this.roomRef.setInput(display.name, input.name);
    }
    
    // if (display.input != input.id) {
    //   document.getElementById("input" + input.id).classList.toggle("feedback");
    //   this.roomRef.setInput(display.name, input.id);
    // }
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
    // this.selectedInput = this.cg.inputs.find((i) => i.id === d.input)
    // if (this.selectedInput == undefined) {
    //   this.selectedInput = this.blank;
    // }
  }

  public getInputIcon(d: DisplayGroup) {
    const input = d.inputs.find((i) => i.name === d.input);
    if (input == undefined) {
      return "crop_landscape";
    }
    return input.icon;
  }

  public getInputName(d: DisplayGroup) {
    const input = d.inputs.find((i) => i.name === d.input);
    if (input == undefined) {
      return "unknown";
    }
    return input.name;
  }
}
