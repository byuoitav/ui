import { Component, OnInit, ViewEncapsulation, Input as AngularInput, Output as AngularOutput, EventEmitter, ViewChild } from '@angular/core';
import { MatSlider } from '@angular/material';
import { AudioDevice, AudioGroup, ControlGroup } from '../../../../../objects/control';
import { RoomRef } from '../../../services/bff.service';

@Component({
  selector: 'volume',
  templateUrl: './volume.component.html',
  styleUrls: ['./volume.component.scss'],
  encapsulation: ViewEncapsulation.Emulated

})
export class VolumeComponent implements OnInit {
  mute: boolean;
  @AngularInput()
  level: number;

  @AngularInput()
  roomRef: RoomRef;

  @AngularInput()
  audioGroupName: string;

  @AngularInput()
  audioDevice: AudioDevice;

  @AngularOutput()
  levelChange: EventEmitter<number> = new EventEmitter();
  @AngularOutput()
  muteChange: EventEmitter<boolean> = new EventEmitter();

  @ViewChild("slider", {static: true})
  slider: MatSlider;

  cg: ControlGroup
  constructor() { }

  ngOnInit() {
    if (this.audioDevice) {
      this.mute = this.audioDevice.muted
    }
    if (this.audioGroupName) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          this.cg = r.controlGroups[r.selectedControlGroup]
          if (this.audioGroupName == "MediaAudio") {
            //do it for media audio
            this.mute = this.cg.mediaAudio.muted;
          } else {
            //do it for the audiodevice
            //find the group
            let ag = this.cg.audioGroups.find((ag) => ag.name === this.audioGroupName)
            //find the actual device and get the mute state
            let dev = ag.audioDevices.find((dev) => dev.name === this.audioDevice.name)
            this.mute = dev.muted
          }
        }
      })
    }
  }

  toggleMute() {
    if (this.audioGroupName == "MediaAudio") {
      this.roomRef.setMuted(!this.mute);
      this.roomRef.buttonPress("master mute set to", String(!this.cg.mediaAudio.muted));
    } else {
      this.roomRef.setMuted(!this.mute, this.audioGroupName, this.audioDevice.name)
      this.roomRef.buttonPress("mute set on " + this.audioDevice.name + " to", String(!this.audioDevice.muted))
    }
  }

  public closeThumb() {
    setTimeout(() => {
      this.slider._elementRef.nativeElement.blur();
    }, 2000);
  }

  
}
