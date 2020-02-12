import { Component, OnInit, Input } from "@angular/core";
import { MatDialog } from "@angular/material";

import { RoomRef } from "../../services/bff.service";
import { HelpInfoComponent } from "./help-info/help-info.component";
import { ControlGroup } from '../../../../../objects/control';

@Component({
  selector: "app-help",
  templateUrl: "./help.component.html",
  styleUrls: ["./help.component.scss"]
})
export class HelpComponent implements OnInit {
  helpHasBeenSent = false;
  @Input() cg: ControlGroup;
  @Input() private _roomRef: RoomRef;

  // TODO change the way help has been sent is tracked
  // TODO: figure out how to track the status of the help request and stuff

  constructor(private dialog: MatDialog) {}

  ngOnInit() {}

  sendForHelp = () => {
    this.dialog
      .open(HelpInfoComponent, { data: this.cg })
      .afterClosed()
      .subscribe(info => {
        if (info) {
          this._roomRef.requestHelp(info);
        }
      });
  };
}
