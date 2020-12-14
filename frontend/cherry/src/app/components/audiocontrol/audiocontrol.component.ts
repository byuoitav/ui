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

  groupPages: Map<string, number[]> = new Map();
  groupCurPage: Map<string, number> = new Map();
  constructor() { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          if (!this.cg || this.cg.audioGroups.length != r.controlGroups[r.selectedControlGroup].audioGroups.length) {
            this.cg = r.controlGroups[r.selectedControlGroup];
            if (this.cg.audioGroups.length > 0) {
              if (this.groupPages.size != this.cg.audioGroups.length) {
                this.cg.audioGroups.forEach(group => {
                  if (!this.groupPages.get(group.name)) {
                    const numPages = Math.ceil(group.audioDevices.length / 4);
                    const tempPages = new Array(numPages).fill(undefined).map((x, i) => i);
                    this.groupPages.set(group.name, tempPages);
                  }
                  if (!this.groupCurPage.get(group.name)) {
                    this.groupCurPage.set(group.name, 0);
                  }
                  console.log(
                    group.name, ":",
                    group.audioDevices.length,
                    "pages:",
                    this.groupPages.get(group.name)
                  )
                });
              }
            }
          }
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

  groupPageLeft(groupName: string) {
    if (this.groupCanPageLeft(groupName)) {
      let pageNum = this.groupCurPage.get(groupName);
      pageNum--;
      this.groupCurPage.set(groupName, pageNum);
    }

    const idx = 4 * this.groupCurPage.get(groupName);
    // probably have to look at exactly what needs to be selected
    document.querySelector("#" + groupName + idx).scrollIntoView({ 
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  }

  groupPageRight(groupName: string) {
    if (this.groupCanPageRight(groupName)) {
      let pageNum = this.groupCurPage.get(groupName);
      pageNum++;
      this.groupCurPage.set(groupName, pageNum);
    }

    const idx = 4 * this.groupCurPage.get(groupName);
    document.querySelector("#" + groupName + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  }

  groupCanPageLeft(groupName: string) {
    if (this.groupCurPage.get(groupName) <= 0) {
      return false;
    }
    return true;
  }

  groupCanPageRight(groupName: string) {
    if (this.groupCurPage.get(groupName) + 1 >= this.groupPages.get(groupName).length) {
      return false;
    }
    return true;
  }
}
