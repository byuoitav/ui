import { Component, OnInit, ViewEncapsulation, Input, ElementRef } from '@angular/core';


class SquareButtonBase {
  constructor(public _elementRef: ElementRef) {}
}

export type ButtonAction = (data: any) => Promise<boolean>;

@Component({
  selector: 'square-button',
  encapsulation: ViewEncapsulation.None,
  templateUrl: './square-button.component.html',
  styleUrls: ['./square-button.component.scss']
})
export class SquareButtonComponent extends SquareButtonBase implements OnInit {
  @Input() data: any;
  @Input() id: string;
  @Input() click: ButtonAction;
  @Input() press: ButtonAction;
  @Input() selected = false;
  @Input() icon: string;
  @Input() title: string;
  @Input() subIcon: string;
  @Input() subTitle: string;
  @Input() showIcon = true;
  @Input() empty = false;
  @Input() curDisplayID: string;

  constructor(elementRef: ElementRef) {
    super(elementRef);
  }

  ngOnInit() {
  }

  ngOnChanges() {
    if (this.id) {
      let btn = document.getElementById(this.id);
      if (btn && this.selected && btn.classList.contains("feedback")) {
        btn.classList.remove("feedback")
      }
    }
    
  }

  do(f: ButtonAction) {
    if (this.id && !this.selected) {
      document.getElementById(this.id).classList.toggle("feedback")
    }
    if (!f) {
      console.warn('no function for this action has been defined');
      return;
    }

    f(this.data);
    // setTimeout(() => {
      

    // }, 700);
  }
}
