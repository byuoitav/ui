import { Component, OnInit, ViewChild, ElementRef, AfterViewInit, EventEmitter, Inject } from '@angular/core';
import Keyboard from 'simple-keyboard';
import { MatBottomSheetRef, MAT_BOTTOM_SHEET_DATA } from '@angular/material';
import { ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-numpad-keyboard',
  encapsulation: ViewEncapsulation.None,
  templateUrl: './numpad.component.html',
  styleUrls: [
    './numpad.component.scss',
    '../../../../node_modules/simple-keyboard/build/css/index.css']
})
export class NumpadComponent implements OnInit, AfterViewInit {
  private keyboard: Keyboard;
  roomCodeValue = '';

  constructor(
    private bottomSheetRef: MatBottomSheetRef<NumpadComponent>,
    @Inject(MAT_BOTTOM_SHEET_DATA) public data: EventEmitter<string>) {
   }

  ngOnInit() {
  }

  ngAfterViewInit() {
    this.keyboard = new Keyboard({
      onChange: input => this.onChange(input),
      onKeyPress: button => this.onKeyPress(button),
      layout: {
        default: [
          '1 2 3',
          '4 5 6',
          '7 8 9',
          '{bksp} 0 {enter}'
        ]
      },
      display: {
        '{bksp}': '⌫',
        '{enter}': 'OK'
      },
      maxLength: {
        default: 6
      }
    });

    this.keyboard.addButtonTheme('{enter}', 'kb-done');
  }

  onChange = (input: string) => {
    this.roomCodeValue = input;
    this.data.emit(this.roomCodeValue);
  }

  onKeyPress = (button: string) => {
    // this.roomCode.nativeElement.focus();
    if (button === '{bksp}') {
      this.roomCodeValue = this.roomCodeValue.substring(0, this.roomCodeValue.length - 1);
    }
    if (button === '{enter}') {
      this.bottomSheetRef.dismiss('all good in the hood');
    }
  }
}
