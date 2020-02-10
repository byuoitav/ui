import { Component, OnInit, Input as AngularInput } from "@angular/core";

import { RoomRef } from "src/app/services/bff.service";
import {
  ControlGroup,
  Display,
  Input,
  IconPair
} from "src/app/objects/control";
import { IControlTab } from "../control-tab/icontrol-tab";

class Page {
  pageOption: string;
  weight: number;
  displays: Display[];

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
export class DisplayComponent implements OnInit, IControlTab {
  @AngularInput() cg: ControlGroup;
  @AngularInput() private _roomRef: RoomRef;

  selectedDisplayIdx: number = 0;
  get selectedDisplay() {
    if (this.cg && this.cg.displayBlocks && this.cg.displayBlocks.length > 0) {
      return this.cg.displayBlocks[this.selectedDisplayIdx];
    }

    return undefined;
  }

  displayPages: Page[];
  curDisplayPage = 0;
  inputPages: number[] = [];
  curInputPage = 0;

  constructor() {
    this.displayPages = [];
  }

  ngOnInit() {}

  ngOnChanges() {
    if (this.cg) {
      this.generatePages();
    }
  }

  generatePages = () => {
    if (this.cg === undefined || this.cg.displayBlocks === undefined) {
      console.error("uninitialized control group");
      return;
    }

    this.displayPages = [];
    this.inputPages = [];

    let displayIndex = 0;

    let p = new Page();
    p.displays = [];

    while (displayIndex < this.cg.displayBlocks.length) {
      if (
        p.weight > 0 &&
        p.weight + this.cg.displayBlocks[displayIndex].outputs.length >= 5
      ) {
        this.displayPages.push(p);
        p = new Page();
      }

      // set the length of the outputs to the weight of the page
      p.weight += this.cg.displayBlocks[displayIndex].outputs.length;
      p.displays.push(this.cg.displayBlocks[displayIndex]);
      if (p.weight > 4) {
        p.pageOption = "4";
      } else {
        p.pageOption += "" + this.cg.displayBlocks[displayIndex].outputs.length;
      }

      // check to see if the weight is less than the max
      if (p.weight >= 4) {
        // assign the page and move on to the next one
        this.displayPages.push(p);
        p = new Page();
      } else {
        if (displayIndex === this.cg.displayBlocks.length - 1) {
          this.displayPages.push(p);
        }
      }

      displayIndex++;
    }

    // set up the input pages
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
    // this.selectedDisplay.input = input.id;
    this._roomRef.setInput(this.selectedDisplay.id, input.id);
  };

  setVolume = (level: number) => {
    // set the volume in some way
    this._roomRef.setVolume(this.cg.audioGroups[0].audioDevices[0].id, level);
  };

  setMute = (muted: boolean) => {
    // mute the volume in some way
    this._roomRef.setMuted(this.cg.audioGroups[0].audioDevices[0].id, muted);
  };
}
