import { Component, OnInit, Input as AngularInput, AfterViewInit, ViewChild } from '@angular/core';
import { RoomRef } from '../../../services/bff.service';
import { ControlGroup } from '../../../../../objects/control';
import { MatTabGroup } from '@angular/material';

@Component({
  selector: 'audiocontrol',
  templateUrl: './audiocontrol.component.html',
  styleUrls: ['./audiocontrol.component.scss']
})
export class AudioControlComponent implements OnInit, AfterViewInit {

  @ViewChild("tabs", {static: true})
  tabs: MatTabGroup;
  @AngularInput()
  roomRef: RoomRef;
  cg: ControlGroup;
  constructor() { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          this.cg = r.controlGroups[r.selectedControlGroup];
        }
      })
    }
  }

  ngAfterViewInit() {
    // this is disgusting. :(
    // but, it moves the second line of tabs to be left aligned
    this.tabs._elementRef.nativeElement.getElementsByClassName(
      "mat-tab-labels"
    )[0].style.justifyContent = "flex-start";
  }

}
