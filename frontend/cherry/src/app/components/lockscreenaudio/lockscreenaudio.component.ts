import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { RoomRef } from '../../../services/bff.service';
import { ControlGroup } from '../../../../../objects/control';
import { MatTabGroup } from '@angular/material';

@Component({
  selector: 'lock-screen-audio',
  templateUrl: './lockscreenaudio.component.html',
  styleUrls: ['./lockscreenaudio.component.scss']
})
export class LockScreenAudioComponent implements OnInit {
  @AngularInput()
  roomRef: RoomRef;
  cg: ControlGroup;
  tabs: MatTabGroup;
  public _show: boolean;

  groupPages: Map<string, number[]> = new Map();
  groupCurPage: Map<string, number> = new Map();
  constructor() { }

  ngOnInit() {
    this._show = false;
  }

  show(roomRef: RoomRef) {
    this._show = true;
    this.roomRef = roomRef; 
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        if (!this.cg || this.cg.audioGroups.length != r.controlGroups[r.selectedControlGroup].audioGroups.length) {
          this.cg = r.controlGroups[r.selectedControlGroup];
          if (this.cg.audioGroups.length > 0) {
            if (this.groupPages.size != this.cg.audioGroups.length) {
              this.cg.audioGroups.forEach(group => {
                if (!this.groupPages.get(group.id)) {
                  const numPages = Math.ceil(group.audioDevices.length / 4);
                  const tempPages = new Array(numPages).fill(undefined).map((x, i) => i);
                  this.groupPages.set(group.id, tempPages);
                }
                if (!this.groupCurPage.get(group.id)) {
                  this.groupCurPage.set(group.id, 0);
                }
                console.log(
                  group.id, ":",
                  group.audioDevices.length,
                  "pages:",
                  this.groupPages.get(group.id)
                )
              });
            }
          }
        }
      }
    })
  }

  hide() {
    this._show = false;
  }

  isShowing() {
    return this._show;
  }

  groupPageLeft(groupID: string) {
    if (this.groupCanPageLeft(groupID)) {
      let pageNum = this.groupCurPage.get(groupID);
      pageNum--;
      this.groupCurPage.set(groupID, pageNum);
    }

    const idx = 4 * this.groupCurPage.get(groupID);
    // probably have to look at exactly what needs to be selected
    document.querySelector("#" + groupID + idx).scrollIntoView({ 
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  }

  groupPageRight(groupID: string) {
    if (this.groupCanPageRight(groupID)) {
      let pageNum = this.groupCurPage.get(groupID);
      pageNum++;
      this.groupCurPage.set(groupID, pageNum);
    }

    const idx = 4 * this.groupCurPage.get(groupID);
    document.querySelector("#" + groupID + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  }

  groupCanPageLeft(groupID: string) {
    if (this.groupCurPage.get(groupID) <= 0) {
      return false;
    }
    return true;
  }

  groupCanPageRight(groupID: string) {
    if (this.groupCurPage.get(groupID) + 1 >= this.groupPages.get(groupID).length) {
      return false;
    }
    return true;
  }
}
