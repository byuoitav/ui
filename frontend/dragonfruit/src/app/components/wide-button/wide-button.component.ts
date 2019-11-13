import { Component, OnInit, ElementRef, ViewEncapsulation, Input } from '@angular/core';
import { ButtonAction } from '../square-button/square-button.component';
import { IconPair } from 'src/app/objects/control';

class WideButtonBase {
  constructor(public _elementRef: ElementRef) {}
}

@Component({
  selector: 'wide-button',
  encapsulation: ViewEncapsulation.None,
  templateUrl: './wide-button.component.html',
  styleUrls: ['./wide-button.component.scss']
})
export class WideButtonComponent extends WideButtonBase implements OnInit {
  @Input() data: any;
  @Input() click: ButtonAction;
  @Input() press: ButtonAction;
  @Input() selected = false;
  @Input() outputs: IconPair[];
  @Input() input: IconPair;
  @Input() showIcons = true;
  @Input() multiple = false;
  @Input() full = false;

  constructor(elementRef: ElementRef) {
    super(elementRef);
  }

  ngOnInit() {
  }

  do(f: ButtonAction) {
    if (!f) {
      console.warn('no function for this action has been defined');
      return;
    }

    f(this.data);
  }
}
