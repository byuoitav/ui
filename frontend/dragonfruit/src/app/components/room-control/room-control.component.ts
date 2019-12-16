import { Component, OnInit, HostListener } from "@angular/core";
import { BFFService } from "src/app/services/bff.service";
import { ActivatedRoute, Router } from "@angular/router";
import { Preset } from "src/app/objects/database";
import {
  ControlGroup,
  CONTROL_TAB,
  AUDIO_TAB,
  PRESENT_TAB,
  HELP_TAB
} from "src/app/objects/control";
import { MatTabChangeEvent, MatTab } from "@angular/material";

@Component({
  selector: "app-room-control",
  templateUrl: "./room-control.component.html",
  styleUrls: ["./room-control.component.scss"]
})
export class RoomControlComponent implements OnInit {
  controlGroup: ControlGroup;
  groupIndex: string;
  roomID: string;
  controlKey: string;

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
    this.route.params.subscribe(params => {
      this.controlKey = params["key"];
      this.roomID = params["id"];
      this.groupIndex = params["index"];
      this.selectedTab = +params["tabName"];
      /*
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
      */
    });
  }

  ngOnInit() {
    if (window.innerWidth >= 768) {
      this.tabPosition = "above";
    } else {
      this.tabPosition = "below";
    }
  }

  goBack = () => {
    this.router.navigate(["/key/" + this.controlKey + "/room/" + this.roomID]);
  };

  tabChange(index: number | string) {
    this.selectedTab = index;
    const currentURL = decodeURI(window.location.pathname);
    const newURL =
      currentURL.substr(0, currentURL.lastIndexOf("/") + 1) + this.selectedTab;
    this.router.navigate([newURL]);
  }
}
