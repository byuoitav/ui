import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { RoomRef } from '../../../services/bff.service';
import { ControlGroup } from '../../../../../objects/control';

@Component({
  selector: 'lock-screen-audio',
  templateUrl: './lockscreenaudio.component.html',
  styleUrls: ['./lockscreenaudio.component.scss']
})
export class LockScreenAudioComponent implements OnInit {
  @AngularInput()
  roomRef: RoomRef;
  cg: ControlGroup;
  public _show: boolean;
  constructor() { }

  ngOnInit() {
    this._show = false;
  }

  show(roomRef: RoomRef) {
    this._show = true;
    this.roomRef = roomRef; 
    this.roomRef.subject().subscribe((r) => {
      if (r) {
        this.cg = r.controlGroups[r.selectedControlGroup];
      }
    })
  }

  hide() {
    this._show = false;
  }

  isShowing() {
    return this._show;
  }
}
