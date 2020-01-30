import { Component, OnInit, ViewEncapsulation, Input } from '@angular/core';

@Component({
  selector: 'volume',
  templateUrl: './volume.component.html',
  styleUrls: ['./volume.component.scss'],
  encapsulation: ViewEncapsulation.Emulated

})
export class VolumeComponent implements OnInit {
  @Input()
  mute: boolean;
  constructor() { }

  ngOnInit() {
  }

}
