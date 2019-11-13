import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { ControlGroup, AudioDevice, AudioGroup, Input } from 'src/app/objects/control';
import { BFFService } from 'src/app/services/bff.service';
import { IControlTab } from '../control-tab/icontrol-tab';

@Component({
  selector: 'app-audio',
  templateUrl: './audio.component.html',
  styleUrls: ['./audio.component.scss']
})
export class AudioComponent implements OnInit, IControlTab {
  @AngularInput() cg: ControlGroup;

  constructor(private bff: BFFService) {
  }

  ngOnInit() {
    console.log('I know who I am', this.cg);
  }

  setVolume = (level: number, device: any) => {
    const audioDevice = device as AudioDevice;
    this.bff.setVolume(audioDevice, level);
  }

  setMute = (muted: boolean, device: any) => {
    const audioDevice = device as AudioDevice;
    this.bff.setMuted(audioDevice, muted);
  }

  muteAll = (ag: AudioGroup) => {
    let muteState = true;
    // if (ag.allAreMuted()) {
    //   muteState = false;
    // }

    for (const ad of ag.audioDevices) {
      // this.bff.setMute(this.cg, muteState, ad.id);
    }
  }
}
