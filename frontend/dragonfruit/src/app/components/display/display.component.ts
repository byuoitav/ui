import { Component, OnInit, Input as AngularInput } from "@angular/core";

import { RoomRef } from "src/app/services/bff.service";
import {
  ControlGroup,
  DisplayGroup,
  Input,
  IconPair
} from "../../../../../objects/control";
import { MatDialogRef, MatDialog } from '@angular/material';
import { MirrorComponent } from 'src/app/dialogs/mirror/mirror.component';

class Page {
  pageOption: string;
  weight: number;
  displays: DisplayGroup[];

  constructor() {
    this.displays = [];
    this.weight = 0;
    this.pageOption = "";
  }
}

@Component({
  selector: "app-display",
  templateUrl: "./display.component.html",
  styleUrls: ["./display.component.scss"]
})
export class DisplayComponent implements OnInit {
  @AngularInput() cg: ControlGroup;
  @AngularInput() private _roomRef: RoomRef;

  selectedDisplayIdx: number = 0;
  get selectedDisplay() {
    if (this.cg && this.cg.displayGroups && this.cg.displayGroups.length > 0) {
      return this.cg.displayGroups[this.selectedDisplayIdx];
    }

    return undefined;
  }

  displayPages: Page[];
  curDisplayPage = 0;
  inputPages: number[] = [];
  curInputPage = 0;

  mirrorRef: MatDialogRef<MirrorComponent>;

  blankInput: Input = {
    id: "blank",
    name: "Blank",
    icon: "crop_landscape",
    subInputs: []
  }

  constructor(public dialog: MatDialog) {
    this.displayPages = [];
  }

  ngOnInit() {}

  ngOnChanges() {
    if (this.cg) {
      console.log(this.cg);
      this.generatePages();
      if (this.cg.displayGroups[0].shareInfo.state === 3 && !this.dialog.openDialogs.includes(this.mirrorRef)) {
        this.mirrorRef = this.dialog.open(MirrorComponent, {
          data: {
            roomRef: this._roomRef
          },
          disableClose: true
        })
      }
    }
  }

  generatePages = () => {
    if (this.cg === undefined || this.cg.displayGroups === undefined) {
      console.error("uninitialized control group");
      return;
    }

    this.displayPages = [];
    this.inputPages = [];

    let displayIndex = 0;

    let p = new Page();
    p.displays = [];

    for (let i = 0; i < this.cg.displayGroups.length; i++) {
      p.weight += this.cg.displayGroups[i].displays.length;
      p.displays.push(this.cg.displayGroups[i]);

      if (p.weight > 4) {
        p.pageOption = "4";
      } else {
        p.pageOption += "" + this.cg.displayGroups[i].displays.length;
      }

      if (p.weight >= 3) {
        this.displayPages.push(p);
        p = new Page();
      } else {
        if (i === this.cg.displayGroups.length - 1) {
          this.displayPages.push(p);
        }
      }
    }

    console.log(this.displayPages);

    // set up the input pages
    this.cg.inputs.unshift(this.blankInput);
    const fullPages = Math.floor(this.cg.inputs.length / 6);
    const remainderPage = this.cg.inputs.length % 6;

    for (let i = 0; i < fullPages; i++) {
      this.inputPages.push(6);
    }

    if (remainderPage !== 0) {
      this.inputPages.push(remainderPage);
    }

    this.curInputPage = 0;
  };

  getInputInfo(inputID: string): IconPair {
    const i = this.cg.inputs.find(x => {
      return x.id.includes(inputID);
    });

    const pair = {
      id: i.id,
      name: i.name,
      icon: i.icon
    };

    return pair;
  }

  onSwipe = (event, section: string) => {
    const x =
      Math.abs(event.deltaX) > 40 ? (event.deltaX > 0 ? "right" : "left") : "";
    const y =
      Math.abs(event.deltaY) > 40 ? (event.deltaY > 0 ? "down" : "up") : "";

    if (x === "right" && this.canPageLeft(section)) {
      this.pageLeft(section);
    }
    if (x === "left" && this.canPageRight(section)) {
      this.pageRight(section);
    }
  };

  canPageLeft = (section: string): boolean => {
    if (section === "display") {
      if (this.curDisplayPage <= 0) {
        return false;
      }
      return true;
    }
    if (section === "input") {
      if (this.curInputPage <= 0) {
        return false;
      }
      return true;
    }
  };

  canPageRight = (section: string): boolean => {
    if (section === "display") {
      if (this.curDisplayPage + 1 >= this.displayPages.length) {
        return false;
      }
      return true;
    }
    if (section === "input") {
      if (this.curInputPage + 1 >= this.inputPages.length) {
        return false;
      }
      return true;
    }
  };

  pageLeft = (section: string) => {
    let idx = 0;
    if (this.canPageLeft(section)) {
      if (section === "display") {
        this.curDisplayPage--;
        idx = this.curDisplayPage;
      }
      if (section === "input") {
        this.curInputPage--;
        idx = this.curInputPage;
      }
    } else {
      return;
    }

    // scroll the page into view
    document.querySelector("#" + section + "-page" + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "nearest"
    });
  };

  pageRight = (section: string) => {
    let idx = 0;
    if (this.canPageRight(section)) {
      if (section === "display") {
        this.curDisplayPage++;
        idx = this.curDisplayPage;
      }
      if (section === "input") {
        this.curInputPage++;
        idx = this.curInputPage;
      }
    } else {
      return;
    }

    // scroll the page into view
    document.querySelector("#" + section + "-page" + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "nearest"
    });
  };

  pageToNumber = (section: string, pageNum: number) => {
    if (section === "display") {
      this.curDisplayPage = pageNum;
    }
    if (section === "input") {
      this.curInputPage = pageNum;
    }

    // scroll the page into view
    document.querySelector("#" + section + "-page" + pageNum).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "nearest"
    });
  };

  setInput = (input: Input) => {
    if (this.isSelected(input.id)) {
      return;
    }
    if (input.id === "blank" && !this.selectedDisplay.blanked) {
      this._roomRef.setBlanked(this.selectedDisplay.id, true);
    } else {
      this._roomRef.setInput(this.selectedDisplay.id, input.id);
    }
  };

  setVolume = (level: number) => {
    // set the volume in some way
    this._roomRef.setVolume(level);
  };

  setMute = (muted: boolean) => {
    // mute the volume in some way
    this._roomRef.setMuted(muted);
  };

  isSelected = (id: string): boolean => {
    if (id === "blank" && this.selectedDisplay.blanked){
      return true;
    } else if (id === this.selectedDisplay.input && !this.selectedDisplay.blanked) {
      return true; 
    }
    return false;
  };
}
