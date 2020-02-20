import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { ControlGroup } from '../../../../../objects/control'
import { RoomRef } from '../../../services/bff.service'

@Component({
  selector: 'screencontrol',
  templateUrl: './screencontrol.component.html',
  styleUrls: ['./screencontrol.component.scss']
})
export class ScreenControlComponent implements OnInit {

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

  public screenUp(screen: string) {
    this.roomRef.raiseProjectorScreen(screen);
  }

  public screenStop(screen: string) {
    this.roomRef.stopProjectorScreen(screen);
  }

  public screenDown(screen: string) {
    this.roomRef.lowerProjectorScreen(screen);
  }
}
