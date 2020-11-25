import {
  Component,
  Input as AngularInput,
  Output as AngularOutput,
  AfterContentInit,
  ElementRef,
  ViewChild,
  EventEmitter
} from "@angular/core";
import { MatDialog } from "@angular/material";
import { RoomRef, BFFService } from 'src/app/services/bff.service';
import { ControlGroup, Input } from '../../../../../objects/control';
import { MinionComponent } from 'src/app/dialogs/minion/minion.component';


@Component({
  selector: "app-wheel",
  templateUrl: "./wheel.component.html",
  styleUrls: ["./wheel.component.scss", "../../colorscheme.scss"]
})
export class WheelComponent {
  private static TITLE_ANGLE = 100;
  private static TITLE_ANGLE_ROTATE: number = WheelComponent.TITLE_ANGLE / 2;

  @AngularInput()
  roomRef: RoomRef;

  cg: ControlGroup;

  top = "50vh";
  right = "50vw";

  arcpath: string;
  titlearcpath: string;
  rightoffset: string;
  topoffset: string;
  translate: string;
  circleOpen = true;
  thumbLabel = true;
  mirrorMaster: Input;
  blank: Input;
  // we need this because of how we chose to do blanking
  lastInputName: string;


  @ViewChild("wheel", {static: false}) wheel: ElementRef;

  constructor(
    public bff: BFFService,
    public dialog: MatDialog
  ) {    
  }

  ngOnInit() {
    if (this.roomRef) {
      this.roomRef.subject().subscribe((r) => {
        if (r) {
          console.log("r", r)
          // if (!this.cg || this.cg.displayGroups[0].inputs.length != r.controlGroups[r.selectedControlGroup].inputs.length) {
            // get the blank input from the bff
            this.blank = r.controlGroups[r.selectedControlGroup].displayGroups[0].inputs.find((input) => input.name === "Blank")
            this.cg = r.controlGroups[r.selectedControlGroup];

            // remove the blank input from the wheel
            this.cg.displayGroups[0].inputs.forEach((input, index) =>{
              if (input.name == "Blank") this.cg.displayGroups[0].inputs.splice(index, 1)
            })

            setTimeout(() => {
              this.render();
            }, 0);
          // } else {
          //   console.log("hello")
          //   this._applyChanges(r.controlGroups[r.selectedControlGroup]);
          // }
        }
      })
    }
  }

  // private _applyChanges(tempCG: ControlGroup) {
  //   this.cg.displayGroups[0].input = tempCG.displayGroups[0].input;
  //   // this.cg.displayGroups[0].blanked = tempCG.displayGroups[0].blanked;
  //   this.cg.mediaAudio = tempCG.mediaAudio; 
  //   // this.cg.audioGroups[0].audioDevices[0] = tempCG.audioGroups[0].audioDevices[0];
  // }

  public render() {
    this.setTranslate();

    const numOfChildren = this.cg.displayGroups[0].inputs.length;
    const children = this.wheel.nativeElement.children;
    const angle = (360 - WheelComponent.TITLE_ANGLE) / numOfChildren;

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
    children[0 + numOfChildren + 1].style.transform = rotate; // rotate the line the corrosponds to this slice
    rotate = "rotate(" + String(WheelComponent.TITLE_ANGLE_ROTATE) + "deg)";
    children[0].firstElementChild.style.transform = rotate;

    for (let i = 1; i <= numOfChildren; ++i) {
      rotate =
        "rotate(" +
        String(angle * -i - WheelComponent.TITLE_ANGLE_ROTATE) +
        "deg)";
      children[i].style.transform = rotate;
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

    switch (this.cg.displayGroups[0].inputs.length) {
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

  setInput = (input: string) => {
    console.log("setting input...", this.cg.displayGroups[0].name)
    this.roomRef.setInput(this.cg.displayGroups[0].name, input);
  }

  setVolume(level: number) {
    this.roomRef.setVolume(level);
  }

  setMute(muted: boolean) {
    this.roomRef.setMuted(muted);
  }

  switchBlanked() {
    if (this.cg.displayGroups[0].input == this.blank.name) {
      // we need to go back to the old input
      this.roomRef.setInput(this.cg.displayGroups[0].name, this.lastInputName);
    } else {
      // we need to save the old input and set the current to blank
      this.lastInputName = this.cg.displayGroups[0].input;
      this.roomRef.setInput(this.cg.displayGroups[0].name, this.blank.name)
    }
    this.roomRef.setBlanked(this.cg.displayGroups[0].name, !this.cg.displayGroups[0].blanked);
  }
}
