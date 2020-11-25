import { Component, OnInit } from '@angular/core';
import { BFFService } from 'src/app/services/bff.service';
import { ControlGroup, AudioGroup } from '../../../../../objects/control';

@Component({
  selector: 'app-audio',
  templateUrl: './audio.component.html',
  styleUrls: ['./audio.component.scss', '../../colorscheme.scss']
})
export class AudioComponent implements OnInit {
  cg: ControlGroup;

  pages: Map<string, number[]>;
  curPageNumbers: Map<string, number>;

  _show: boolean;

  constructor(public bff: BFFService) {
    this._show = false;
  }

  ngOnInit() {
  }

  show = (group: ControlGroup) => {
    this.cg = group;
    this.bff.roomRef.subject().subscribe((r) => {
      if (r) {
        this.applyChanges(r.controlGroups[r.selectedControlGroup]);
      }
    })

    this.pages = new Map();
    this.curPageNumbers = new Map();

    for (const ag of this.cg.audioGroups) {
      const pages = Math.ceil(ag.audioDevices.length / 4);
      this.pages.set(ag.name, new Array(pages).fill(undefined).map((x, i) => i));
      this.curPageNumbers.set(ag.name, 0);
    }

    this._show = true;
  }

  hide = () => {
    this._show = false;
  }

  applyChanges(cg: ControlGroup) {
    for (let i = 0; i < this.cg.audioGroups.length; i++) {
      for (let x = 0; x < this.cg.audioGroups[i].audioDevices.length; x++) {
        if (this.cg.audioGroups[i].audioDevices[x].muted != cg.audioGroups[i].audioDevices[x].muted) {
          this.cg.audioGroups[i].audioDevices[x].muted = cg.audioGroups[i].audioDevices[x].muted;
        }
      }
    }
  }

  isShowing = () => {
    return this._show;
  }

  pageLeft = (ag: AudioGroup) => {
    if (this.canPageLeft(ag)) {
      let n = this.curPageNumbers.get(ag.name);
      let newN = --n
      this.curPageNumbers.set(ag.name, newN);

      const idx = 4 * this.curPageNumbers.get(ag.name);
      document.querySelector("#" + ag.name + idx).scrollIntoView({
        behavior: "smooth",
        block: "nearest",
        inline: "start"
      });
    }
  }

  pageRight = (ag: AudioGroup) => {
    if (this.canPageRight(ag)) {
      let n = this.curPageNumbers.get(ag.name);
      let newN = ++n
      this.curPageNumbers.set(ag.name, newN);

      const idx = 4 * this.curPageNumbers.get(ag.name);
      document.querySelector("#" + ag.name + idx).scrollIntoView({
        behavior: "smooth",
        block: "nearest",
        inline: "start"
      });
    }    
  }

  canPageLeft = (ag: AudioGroup):boolean => {
    return this.curPageNumbers.get(ag.name) > 0
  }

  canPageRight = (ag: AudioGroup):boolean => {
    return this.curPageNumbers.get(ag.name) + 1 < this.pages.get(ag.name).length
  }

  selectPage = (ag: AudioGroup, pageNum: number) => {
    this.curPageNumbers.set(ag.name, pageNum);
  }

  setMute(mute: boolean, id: string) {
    console.log(this.cg.audioGroups[0].muted)
    this.bff.roomRef.setMuted(mute, id);
  }
}
