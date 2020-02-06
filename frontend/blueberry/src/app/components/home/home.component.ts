import { Component, OnInit, Input, ViewChild } from '@angular/core';
import { BFFService, RoomRef } from '../../services/bff.service';
import { ControlGroup } from 'src/app/objects/control';
import { AudioComponent } from '../audio/audio.component';
import { ProjectorComponent } from '../projector/projector.component';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss', '../../colorscheme.scss']
})
export class HomeComponent implements OnInit {
  @Input() roomRef: RoomRef;
  cg: ControlGroup

  @ViewChild(AudioComponent, {static: false}) public audio: AudioComponent;
  @ViewChild(ProjectorComponent, {static: false}) public screen: ProjectorComponent;

  constructor(public bff: BFFService) { }

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
  
  public turnOff() {
    console.log("turning off the room");
    this.bff.locked = true;
  }
}
