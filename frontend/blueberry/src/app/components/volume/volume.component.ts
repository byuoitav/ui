import {
  Component,
  OnInit,
  Input,
  Output,
  ViewChild,
  EventEmitter
} from "@angular/core";
import { MatSlider } from "@angular/material";

@Component({
  selector: "volume",
  templateUrl: "./volume.component.html",
  styleUrls: ["./volume.component.scss"]
})
export class VolumeComponent implements OnInit {
  @Input()
  level: number;
  @Input()
  mute: boolean;
  @Input() name: string;

  @Input()
  muteType: string;

  @Output()
  levelChange: EventEmitter<number> = new EventEmitter();
  @Output()
  muteChange: EventEmitter<boolean> = new EventEmitter();

  @ViewChild("slider", {static: false})
  slider: MatSlider;

  constructor() {}
  ngOnInit() {}

  public closeThumb() {
    setTimeout(() => {
      this.slider._elementRef.nativeElement.blur();
    }, 2000);
  }
}