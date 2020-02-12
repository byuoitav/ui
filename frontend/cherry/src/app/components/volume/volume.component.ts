import { Component, OnInit, ViewEncapsulation, Input as AngularInput, Output as AngularOutput, EventEmitter, ViewChild } from '@angular/core';
import { MatSlider } from '@angular/material';

@Component({
  selector: 'volume',
  templateUrl: './volume.component.html',
  styleUrls: ['./volume.component.scss'],
  encapsulation: ViewEncapsulation.Emulated

})
export class VolumeComponent implements OnInit {
  @AngularInput()
  mute: boolean;
  @AngularInput()
  level: number;

  @AngularOutput()
  levelChange: EventEmitter<number> = new EventEmitter();
  @AngularOutput()
  muteChange: EventEmitter<boolean> = new EventEmitter();

  @ViewChild("slider")
  slider: MatSlider;
  constructor() { }

  ngOnInit() {
  }

  toggleMute() {
    let emit: boolean;
    if (this.mute) {
      emit = false;
    } else {
      emit = true;
    }
    this.muteChange.emit(emit);
    this.mute = !this.mute;
  }

  public closeThumb() {
    setTimeout(() => {
      this.slider._elementRef.nativeElement.blur();
    }, 2000);
  }
}
