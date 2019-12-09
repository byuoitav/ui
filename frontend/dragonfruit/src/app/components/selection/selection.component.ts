import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material';
import { ActivatedRoute, Router } from '@angular/router';
import { BFFService } from 'src/app/services/bff.service';
import { ControlGroup } from 'src/app/objects/control';
import { TurnOffRoomDialogComponent } from 'src/app/dialogs/turnOffRoom-dialog/turnOffRoom-dialog.component';

@Component({
  selector: 'app-selection',
  templateUrl: './selection.component.html',
  styleUrls: ['./selection.component.scss']
})
export class SelectionComponent implements OnInit {
  roomID = '';
  controlKey = '';

  constructor(
    private route: ActivatedRoute,
    public bff: BFFService,
    private dialog: MatDialog,
    private router: Router) {
    this.route.params.subscribe(params => {
      this.roomID = params["id"]
      this.controlKey = params["key"]

      if (this.bff.room == undefined) {
        this.bff.connectToRoom(this.controlKey);
      }
    });
  }

  ngOnInit() {
  }

  goBack = () => {
    this.dialog.open(TurnOffRoomDialogComponent).afterClosed().subscribe(result => {
      // if the result is true then send command to turn off room and redirect page, else redirect webpage
      console.log(`Dialog result: ${result}`);
      this.router.navigate(['/login']);

    });
  }

  selectControlGroup = (cg: ControlGroup): Promise<boolean> => {
    return new Promise<boolean>(() => {
      const index = cg.id;
      this.bff.room.selectedControlGroup = cg.id;
      this.router.navigate(['/key/' + this.controlKey + '/room/' + this.roomID + '/group/' + index + '/tab/0']);
    });
  }
}
