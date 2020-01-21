import { Component, OnInit, AfterContentInit, Input as AngularInput, Output as AngularOutput, ViewChild, ElementRef, NgZone, } from '@angular/core';
import { RoomRef, BFFService } from 'src/app/services/bff.service';
import { ControlGroup, Input } from 'src/app/objects/control';

@Component({
  selector: 'app-wheel',
  templateUrl: './wheel.component.html',
  styleUrls: ['./wheel.component.scss', '../../colorscheme.scss']
})
export class WheelComponent implements AfterContentInit, OnInit {
  private static TITLE_ANGLE = 100;
  private static TITLE_ANGLE_ROTATE: number = WheelComponent.TITLE_ANGLE / 2;

  @AngularInput() blur: boolean;

  @AngularInput() room: RoomRef;

  @AngularInput()
  top: string;
  @AngularInput()
  right: string;

  cg: ControlGroup;

  circleOpen = true;
  arcpath: string;
  titlearcpath: string;
  rightoffset: string;
  topoffset: string;
  translate: string;
  thumbLabel = true;

  @ViewChild("wheel", {static: false}) wheel: ElementRef;

  constructor(public bff: BFFService) {
    
  }
  
  ngOnChanges() {
    console.log("this is on changes", this.room);
    
  }

  ngOnInit() {
    console.log("this is on init", this.room);
    
  }

  ngDoCheck() {
    console.log("this is do check", this.room);
  }

  ngAfterContentInit() {
    console.log("this is after content init", this.room);
      if (this.bff.roomRef) {
        this.bff.roomRef.subject().subscribe((room) => {
          if (room) {
            const tempCG = this.cg;
            this.cg = room.controlGroups[room.selectedControlGroup];
    
            if (tempCG == undefined || (this.cg.inputs.length != tempCG.inputs.length)) {
              setTimeout(() => {
                this.render();
              }, 0)
            }
             
          }
        });
    }
    
  }

  ngAfterContentChecked() {
    console.log("this is after content checked", this.room);
  }

  ngAfterViewInit() {
    console.log("this is after view init", this.room);
  }

  ngAfterViewChecked() {
    console.log("this is after view checked", this.room);
  }

  

  public render() {
    console.log("rendering...")
    this.setTranslate();

    const numOfChildren = this.cg.inputs.length;
    const children = this.wheel.nativeElement.children;
    const angle = (360 - WheelComponent.TITLE_ANGLE) / numOfChildren;
    // console.log("children", children.length);
    // console.log("angle", angle);

    this.arcpath = this.getArc(0.5, 0.5, 0.5, 0, angle);
    this.titlearcpath = this.getArc(
      0.5,
      0.5,
      0.5,
      0,
      WheelComponent.TITLE_ANGLE
    );

    let rotate =
      "rotate(" + String(-WheelComponent.TITLE_ANGLE_ROTATE) + "deg)";
    children[0].style.transform = rotate;
    // console.log("what number is this?", 0 + numOfChildren + 1);
    children[0 + numOfChildren + 1].style.transform = rotate; // rotate the line the corrosponds to this slice
    rotate = "rotate(" + String(WheelComponent.TITLE_ANGLE_ROTATE) + "deg)";
    children[0].firstElementChild.style.transform = rotate;

    for (let i = 1; i <= numOfChildren; ++i) {
      rotate =
        "rotate(" +
        String(angle * -i - WheelComponent.TITLE_ANGLE_ROTATE) +
        "deg)";
      children[i].style.transform = rotate;
      // console.log("what number is this?", i + numOfChildren + 1);
      children[i + numOfChildren + 1].style.transform = rotate; // rotate the line that corrosponds to this slice

      rotate =
        "rotate(" +
        String(angle * i + WheelComponent.TITLE_ANGLE_ROTATE) +
        "deg)";
      children[i].firstElementChild.style.transform = rotate;
    }

    this.setInputOffset();
  }

  private setTranslate() {
    const offsetX: number = parseInt(this.right, 10);
    const offsetY: number = parseInt(this.top, 10);

    const x = 50 - offsetX;
    const y = 50 - offsetY;

    this.translate = String("translate(" + x + "vw," + y + "vh)");
  }

  private setInputOffset() {
    let top: number;
    let right: number;

    switch (this.cg.inputs.length) {
      case 7:
        top = -0.6;
        right = 25.4;
        break;
      case 6:
        top = 0.8;
        right = 24;
        break;
      case 5:
        top = 2;
        right = 20.4;
        break;
      case 4:
        top = 4;
        right = 17.5;
        break;
      case 3:
        top = 9;
        right = 12;
        break;
      case 2:
        top = 20;
        right = 2;
        break;
      case 1:
        top = 63;
        right = 7;
        break;
      default:
        console.warn(
          "no configuration for",
          this.cg.inputs.length,
          "inputs"
        );
        break;
    }

    this.topoffset = String(top) + "%";
    this.rightoffset = String(right) + "%";
  }

  private getArc(x, y, radius, startAngle, endAngle): string {
    const start = this.polarToCart(x, y, radius, endAngle);
    const end = this.polarToCart(x, y, radius, startAngle);

    const largeArc = endAngle - startAngle <= 180 ? "0" : "1";

    const d = [
      "M",
      start.x,
      start.y,
      "A",
      radius,
      radius,
      0,
      largeArc,
      0,
      end.x,
      end.y,
      "L",
      x,
      y,
      "L",
      start.x,
      start.y
    ].join(" ");

    return d;
  }

  private polarToCart(cx, cy, r, angle) {
    const angleInRad = ((angle - 90) * Math.PI) / 180.0;

    return {
      x: cx + r * Math.cos(angleInRad),
      y: cy + r * Math.sin(angleInRad)
    };
  }

  setInput = (input: Input) => {
    this.room.setInput(this.cg.displays[0].id, input.id);
  }

  isSelected = (input: Input) => {
    if (this.cg) {
      console.log("checking the current input", this.cg.displays[0].input)
      return this.cg.displays[0].input === input.id
    }
  }
}
