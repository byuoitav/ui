import { Component, OnInit } from '@angular/core';
import { ControlGroup } from '../../../../../objects/control';
import { BFFService } from 'src/app/services/bff.service';

@Component({
  selector: 'app-projector',
  templateUrl: './projector.component.html',
  styleUrls: ['./projector.component.scss']
})
export class ProjectorComponent implements OnInit {
  cg: ControlGroup;
  _show: boolean;

  pages: number[] = [];
  curPage: number;

  constructor(public bff: BFFService) { }

  ngOnInit() {
  }

  show = (cg: ControlGroup) => {
    this.cg = cg;

    // this.cg.screens = ["SCR1"];

    if (this.cg.screens) {
      const pages = Math.ceil(this.cg.screens.length / 4);
      this.pages = new Array(pages).fill(undefined).map((x, i) => i);
  
      console.log("devices:", this.cg.screens.length, "pages:", this.pages);
      this.curPage = 0;
    }
    

    this._show = true;
  }

  hide = () => {
    this._show = false;
  }

  isShowing = (): boolean => {
    return this._show;
  }

  pageLeft = () => {
    if (this.canPageLeft()) {
      this.curPage--;
    }

    // scroll to the bottom of the page
    const idx = 3 * this.curPage;
    console.log(document.querySelector("#device" + idx));
    document.querySelector("#device" + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  };

  pageRight = () => {
    if (this.canPageRight()) {
      this.curPage++;
    }

    // scroll to the top of the page
    const idx = 4 * this.curPage + 3;
    console.log(document.querySelector("#device" + idx));
    document.querySelector("#device" + idx).scrollIntoView({
      behavior: "smooth",
      block: "nearest",
      inline: "start"
    });
  };

  canPageLeft = (): boolean => {
    if (this.curPage <= 0) {
      return false;
    }

    return true;
  };

  canPageRight = (): boolean => {
    if (this.curPage + 1 >= this.pages.length) {
      return false;
    }

    return true;
  };

  projectorUp(screen: string) {
    this.bff.roomRef.raiseProjectorScreen(screen);
  }

  projectorDown(screen: string) {
    this.bff.roomRef.lowerProjectorScreen(screen);
  }

  projectorStop(screen: string) {
    this.bff.roomRef.lowerProjectorScreen(screen);
  }
}
