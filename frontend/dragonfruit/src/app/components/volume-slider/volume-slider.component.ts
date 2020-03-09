import {
  Component,
  OnInit,
  Input as AngularInput,
  ViewChild,
  ViewEncapsulation
} from "@angular/core";
import { MatSlider } from "@angular/material";
import { AudioDevice } from "../../../../../objects/control";

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

  @AngularInput() audioID: string;

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

    if (this.audioID) {
      f(newLevel, this.audioID);
    } else {
      f(newLevel);
    }
  }

  setMute(f: SliderAction, mute: boolean) {
    if (!f) {
      console.warn("no function for this action has been defined");
      return;
    }

    if (this.audioID) {
      f(mute, this.audioID);
    } else {
      f(mute);
    }
  }
}
