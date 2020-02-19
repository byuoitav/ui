import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { RoomRef } from '../../../services/bff.service';
import { ControlGroup } from '../../../objects/control';

@Component({
  selector: 'audiocontrol',
  templateUrl: './audiocontrol.component.html',
  styleUrls: ['./audiocontrol.component.scss']
})
export class AudioControlComponent implements OnInit {

  @AngularInput()
  roomRef: RoomRef;
  cg: ControlGroup;
  constructor() { }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          if (!this.cg) {
            this.cg = r.controlGroups[r.selectedControlGroup];
          }
        }
      })
    }
  }

}
