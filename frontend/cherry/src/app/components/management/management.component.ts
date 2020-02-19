import { Component, OnChanges, Input as AngularInput, SimpleChanges } from '@angular/core';

const LOW= 3;
const REDIRECT: string = "http://" + window.location.hostname + ":10000/dashboard";

@Component({
  selector: 'management',
  templateUrl: './management.component.html',
  styleUrls: ['./management.component.scss']
})
export class ManagementComponent implements OnChanges {
  @AngularInput()
  enabled: boolean;
  defcon: number;

  constructor() {
    this.reset();
   }

  ngOnChanges(changes: SimpleChanges) {
    this.reset();
  }

  public update(level: number) {
    console.log("defcon", this.defcon);

    switch(level) {
      case LOW: //3
        if (this.defcon === LOW) {
          this.defcon--;
        } else {
          this.reset();
        }
        break;
      case LOW -1: //2
        if (this.defcon ===LOW - 1) {
          this.defcon--;
        } else {
          this.reset();
        }
        break;
      case LOW - 2: //1
        if (this.defcon === LOW - 2) {
          this.defcon--;
        } else {
          this.reset();
        }
        break;
      case LOW - 3: //0
        if (this.defcon === LOW - 3) {
          console.log("redirecting to dashboard")
          location.assign(REDIRECT);
        } else {
          this.reset();
        }
        break;
      default:
        this.reset();
        break;
    }
  }

  public reset() {
    this.defcon = LOW;
  }
}
