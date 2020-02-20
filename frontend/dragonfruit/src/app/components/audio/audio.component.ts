import { Component, OnInit, Input as AngularInput } from "@angular/core";

import { RoomRef } from "../../services/bff.service";
import {
  ControlGroup,
  AudioDevice,
  AudioGroup,
  Input
} from "../../../../../objects/control";

@Component({
  selector: "app-audio",
  templateUrl: "./audio.component.html",
  styleUrls: ["./audio.component.scss"]
})
export class AudioComponent implements OnInit {
  @AngularInput() cg: ControlGroup;
  @AngularInput() private _roomRef: RoomRef;

  constructor() {}

  ngOnInit() {}

  setVolume = (level: number, device: any) => {
    const audioDevice = device as AudioDevice;
    this._roomRef.setVolume(audioDevice.id, level);
  };

  setMute = (muted: boolean, device: any) => {
    const audioDevice = device as AudioDevice;
    this._roomRef.setMuted(audioDevice.id, muted);
  };

  // if there is at least one that is not muted
  // then mute everything
  // if all of them are muted, unmute everything
  muteAll = (ag: AudioGroup) => {
    const muted = ag.audioDevices.some(ad => !ad.muted);
    // muted = true if there is at least one that is not muted
    // muted = false if there are no devices that are not muted
    //                  all devices are muted

    for (const ad of ag.audioDevices) {
      this._roomRef.setMuted(ad.id, muted);
    }
  };
}
