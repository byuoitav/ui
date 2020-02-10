import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'lock-screen-audio',
  templateUrl: './lockscreenaudio.component.html',
  styleUrls: ['./lockscreenaudio.component.scss']
})
export class LockScreenAudioComponent implements OnInit {
  public _show: boolean;
  constructor() { }

  ngOnInit() {
    this._show = false;
  }

  show = () => {
    this._show = true;
  }

  hide = () => {
    this._show = false;
  }

}
