import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material";
import { ActivatedRoute, Router } from "@angular/router";

import { RoomRef } from "src/app/services/bff.service";
import { TurnOffRoomDialogComponent } from "src/app/dialogs/turnOffRoom-dialog/turnOffRoom-dialog.component";
import { Room } from "../../../../../objects/control";

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
    private dialog: MatDialog,
    private router: Router
  ) {
    this.route.data.subscribe(data => {
      this._roomRef = data.roomRef;

      this._roomRef.subject().subscribe(room => {
        if (room.selectedControlGroup) {
          this.router.navigate(["./" + room.selectedControlGroup + "/0"], { relativeTo: this.route });
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
          this._roomRef.turnOff();
        }

        this._roomRef.logout();
      });
  };

  selectControlGroup = (id: string) => {
    this._roomRef.selectControlGroup(id);
  };
}
