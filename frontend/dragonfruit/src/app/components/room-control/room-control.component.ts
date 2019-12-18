import { Component, OnInit, HostListener } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { MatTabChangeEvent, MatTab } from "@angular/material";

import { BFFService, RoomRef } from "src/app/services/bff.service";
import {
  Room,
  ControlGroup,
  CONTROL_TAB,
  AUDIO_TAB,
  PRESENT_TAB,
  HELP_TAB
} from "src/app/objects/control";

@Component({
  selector: "app-room-control",
  templateUrl: "./room-control.component.html",
  styleUrls: ["./room-control.component.scss"]
})
export class RoomControlComponent implements OnInit {
  // to use in the template
  public objectKeys = Object.keys;

  public _roomRef: RoomRef;
  get room(): Room {
    if (this._roomRef) {
      return this._roomRef.room;
    }

    return undefined;
  }

  public _controlGroupID: string;
  get controlGroup(): ControlGroup {
    if (this.room && this._controlGroupID) {
      return this.room.controlGroups[this._controlGroupID];
    }

    return undefined;
  }

  tabPosition = "below";
  selectedTab: number | string;

  @HostListener("window:resize", ["$event"])
  onResize(event) {
    if (window.innerWidth >= 768) {
      this.tabPosition = "above";
    } else {
      this.tabPosition = "below";
    }
  }

  constructor(
    public bff: BFFService,
    public route: ActivatedRoute,
    private router: Router
  ) {
    this.route.data.subscribe(data => {
      this._roomRef = data.roomRef;
    });

    this.route.params.subscribe(params => {
      this._controlGroupID = params["groupid"];
      this.selectedTab = +params["tab"];

      // TODO make sure the room has this group, if not, redirect up?
    });

    /*
    this.route.params.subscribe(params => {
      this.selectedTab = +params["tabName"];
      if (this.bff.room === undefined) {
        this.bff.getRoom(this.controlKey);
        // this.bff.connectToRoom(this.controlKey);

        /*
        this.bff.done.subscribe(e => {
          this.controlGroup = this.bff.room.controlGroups[this.groupIndex];
          if (this.controlGroup.id === "Third") {
          }
        });
        // *
      } else {
        // this.controlGroup = this.bff.room.controlGroups[this.groupIndex];
        if (this.controlGroup.id === "Third") {
        }
      }
    });
      */
  }

  ngOnInit() {
    if (window.innerWidth >= 768) {
      this.tabPosition = "above";
    } else {
      this.tabPosition = "below";
    }
  }

  goBack = () => {
    if (this.room && Object.keys(this.room.controlGroups).length == 1) {
      this._roomRef.logout();
    } else {
      this.router.navigate(["../"], { relativeTo: this.route });
    }
  };

  tabChange(index: number | string) {
    this.selectedTab = index;
    const currentURL = decodeURI(window.location.pathname);
    const newURL =
      currentURL.substr(0, currentURL.lastIndexOf("/") + 1) + this.selectedTab;
    this.router.navigate([newURL]);
  }
}
