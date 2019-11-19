import {
  Component,
  OnInit,
  Input as AngularInput,
  ViewChild,
  ViewEncapsulation
} from "@angular/core";
import { MatSlider } from "@angular/material";
import { AudioDevice } from "src/app/objects/control";

export type SliderAction = (attribute: any, data?: any) => Promise<boolean>;

@Component({
  selector: "volume-slider",
  templateUrl: "./volume-slider.component.html",
  styleUrls: ["./volume-slider.component.scss"],
  encapsulation: ViewEncapsulation.Emulated
})
export class VolumeSliderComponent implements OnInit {
  @AngularInput() muted: boolean;
  @AngularInput() level: number;

  @AngularInput() levelChange: SliderAction;
  @AngularInput() muteChange: SliderAction;

  @AngularInput() audioDevice: AudioDevice;

  @AngularInput() master = false;

  @ViewChild("slider", null) slider: MatSlider;

  constructor() {}

  ngOnInit() {}

  // toggleMute() {
  //   this.muted = !this.muted;
  // }

  public closeThumb() {
    setTimeout(() => {
      this.slider._elementRef.nativeElement.blur();
    }, 1500);
  }

  setLevel(f: SliderAction, newLevel: number) {
    if (!f) {
      console.warn("no function for this action has been defined");
      return;
    }

    if (this.audioDevice) {
      f(newLevel, this.audioDevice);
    } else {
      f(newLevel);
    }
  }

  setMute(f: SliderAction, mute: boolean) {
    if (!f) {
      console.warn("no function for this action has been defined");
      return;
    }

    if (this.audioDevice) {
      f(mute, this.audioDevice);
    } else {
      f(mute);
    }
  }
}
