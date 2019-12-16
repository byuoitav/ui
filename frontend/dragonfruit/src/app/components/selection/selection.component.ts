import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material";
import { ActivatedRoute, Router } from "@angular/router";

import { BFFService, RoomRef } from "src/app/services/bff.service";
import { TurnOffRoomDialogComponent } from "src/app/dialogs/turnOffRoom-dialog/turnOffRoom-dialog.component";
import { ControlGroup, Display, Room } from "src/app/objects/control";

@Component({
  selector: "app-selection",
  templateUrl: "./selection.component.html",
  styleUrls: ["./selection.component.scss"]
})
export class SelectionComponent implements OnInit {
  private _roomRef: RoomRef;
  get room(): Room {
    if (this._roomRef) {
      return this._roomRef.room;
    }

    return undefined;
  }

  constructor(
    private route: ActivatedRoute,
    public bff: BFFService,
    private dialog: MatDialog,
    private router: Router
  ) {
    this.route.data.subscribe(data => {
      this._roomRef = data.roomRef;

      this._roomRef.subject().subscribe(room => {
        switch (Object.keys(room.controlGroups).length) {
          case 0:
            // redirect back to login,
            // say that something is wrong with this room?
            break;
          case 1:
            this.selectControlGroup(Object.keys(room.controlGroups)[0]);
            break;
          default:
            break;
        }
      });
    });
  }

  ngOnInit() {}

  goBack = () => {
    this.dialog
      .open(TurnOffRoomDialogComponent)
      .afterClosed()
      .subscribe(result => {
        // if the result is true then send command to turn off room and redirect page, else redirect webpage
        if (result) {
          this.bff.turnOffRoom();
        }

        this.router.navigate(["/login"]);
      });
  };

  selectControlGroup = (cg: string) => {
    console.log("selecting", cg);
    this.router.navigate(["./" + cg], { relativeTo: this.route });
  };
}
