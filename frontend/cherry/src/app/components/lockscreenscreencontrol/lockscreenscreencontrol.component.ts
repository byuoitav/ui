import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'lock-screen-screen-control',
  templateUrl: './lockscreenscreencontrol.component.html',
  styleUrls: ['./lockscreenscreencontrol.component.scss']
})
export class LockScreenScreenControlComponent implements OnInit {
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
