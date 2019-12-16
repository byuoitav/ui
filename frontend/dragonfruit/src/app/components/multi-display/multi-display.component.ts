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
  selector: "app-multi-display",
  templateUrl: "./multi-display.component.html",
  styleUrls: ["./multi-display.component.scss"]
})
export class MultiDisplayComponent implements OnInit, IControlTab {
  @AngularInput() cg: ControlGroup;
  @AngularInput() private _roomRef: RoomRef;

  selectedDisplay: Display;
  displayPages: Page[];
  curDisplayPage = 0;
  inputPages: number[] = [];
  curInputPage = 0;

  constructor() {
    this.displayPages = [];
  }

  ngOnInit() {
    // select an initial display
    if (this.cg && this.cg.displays && this.cg.displays.length > 0) {
      this.selectedDisplay = this.cg.displays[0];
    }
  }

  ngOnChanges() {
    if (this.cg !== undefined) {
      this.generatePages();

      const fullPages = Math.floor(this.cg.inputs.length / 6);
      const remainderPage = this.cg.inputs.length % 6;

      for (let pageIndex = 0; pageIndex < fullPages; pageIndex++) {
        this.inputPages.push(6);
      }

      if (remainderPage !== 0) {
        this.inputPages.push(remainderPage);
      }

      this.curInputPage = 0;
    }
  }

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

  generatePages() {
    console.log("selected display", this.selectedDisplay);
    this.displayPages = [];
    console.log(this.cg);
    if (this.cg === undefined || this.cg.displays === undefined) {
      console.log("uninitialized control group");
      return;
    }
    let dispIndex = 0;

    let p = new Page();
    p.displays = [];

    while (dispIndex < this.cg.displays.length) {
      if (
        p.weight > 0 &&
        p.weight + this.cg.displays[dispIndex].outputs.length >= 5
      ) {
        this.displayPages.push(p);
        p = new Page();
      }

      // set the length of the outputs to the weight of the page
      p.weight += this.cg.displays[dispIndex].outputs.length;
      p.displays.push(this.cg.displays[dispIndex]);
      if (p.weight > 4) {
        p.pageOption = "4";
      } else {
        p.pageOption += "" + this.cg.displays[dispIndex].outputs.length;
      }

      // check to see if the weight is less than the max
      if (p.weight >= 4) {
        // assign the page and move on to the next one
        this.displayPages.push(p);
        p = new Page();
      } else {
        if (dispIndex === this.cg.displays.length - 1) {
          this.displayPages.push(p);
        }
      }

      dispIndex++;
    }

    console.log(this.displayPages);
  }

  onSwipe(evt, section: string) {
    const x =
      Math.abs(evt.deltaX) > 40 ? (evt.deltaX > 0 ? "right" : "left") : "";
    const y = Math.abs(evt.deltaY) > 40 ? (evt.deltaY > 0 ? "down" : "up") : "";

    // console.log(x, y);

    if (x === "right" && this.canPageLeft(section)) {
      // console.log('paging left...');
      this.pageLeft(section);
    }
    if (x === "left" && this.canPageRight(section)) {
      // console.log('paging right...');
      this.pageRight(section);
    }
  }

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
      // console.log('going to page ', this.curPage);
    } else {
      return;
    }

    // scroll to the bottom of the page
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
      // console.log('going to page ', this.curPage);
    } else {
      return;
    }

    // scroll to the bottom of the page
    document.querySelector("#" + section + "-page" + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "nearest"
    });
  };

  pageToNumber(section: string, pageNum: number) {
    if (section === "display") {
      this.curDisplayPage = pageNum;
    }
    if (section === "input") {
      this.curInputPage = pageNum;
    }

    document.querySelector("#" + section + "-page" + pageNum).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "nearest"
    });
  }

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

  setInput = (input: Input) => {
    // select the input now to make it look like it worked :)
    this.selectedDisplay.input = input.id;
    this._roomRef.setInput(this.selectedDisplay.id, input.id);
  };

  setVolume = (level: number) => {
    // this._roomRef.setVolume(this.displayAudio.id, level);
  };

  setMute = (muted: boolean) => {
    // this._roomRef.setMuted(this.displayAudio.id, muted);
  };
}
