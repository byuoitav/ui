import { Component, Inject } from "@angular/core";
import { HELP_TAB } from "../../../../objects/control";
import { Button } from "protractor";
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from "@angular/material";
import { RoomRef } from "../../services/bff.service";
import { ConfirmHelpDialog } from './confirmhelp.dialog'

@Component({
    selector: "help",
    template: `
        <h1 mat-dialog-title class="text">Help</h1>
        
        
        <div mat-dialog-content class="text">
            <p *ngIf="!isAfterHours()">
                Please call AV Support at 801-422-7671 for help, or request help by pressing <i>Request Help</i> to send support staff to you.
            </p>
            <p *ngIf="isAfterHours()">
                No technicians are currently available. For emergencies please call 801-422-7671.
            </p>
        </div>

        <div mat-dialog-actions class="items secondary-theme">
            <button mat-raised-button
                color="warn"
                (click)="cancel()">
                Cancel
            </button>
            <button mat-raised-button
                *ngIf="!isAfterHours()"
                color="primary"
                (click)="requestHelp()"
                (press)="requestHelp()">
                Request Help
            </button>
        </div>
    `,
    styles: [
        `
        .text {
            text-align: center;
            font-family: Roboto, "Helvetica Neue", sans-serif;
            }
    
            .items {
            display: flex;
            flex-direction: row;
            justify-content: center;
            align-items: center;
            }
        `
    ]
})
export class HelpDialog {
    constructor(
        public dialogRef: MatDialogRef<HelpDialog>,
        @Inject(MAT_DIALOG_DATA) public roomRef: RoomRef,
        public dialog: MatDialog,
    ) { }

    public cancel() {
        this.dialogRef.close();
    }

    public requestHelp() {
      this.roomRef.requestHelp(location.hostname + "needs help");
      const dialogRef = this.dialog.open(ConfirmHelpDialog, {
        width: "70vw",
        disableClose: true
      });
    }

    public isAfterHours(): boolean {
        let date = new Date();
        let DayOfTheWeek = date.getDay();
        let CurrentHour = date.getHours();
    
        switch(DayOfTheWeek) {
          // Sunday
          case 0: { return true; }
          // Monday
          case 1: {
            if(CurrentHour < 7 || CurrentHour >= 19) { return true; }
            else { return false; }
          }
          // Tuesday
          case 2: {
            if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
            else { return false; }
          }
          // Wednesday
          case 3: {
            if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
            else { return false; }
          }
          // Thursday
          case 4: {
            if(CurrentHour < 7 || CurrentHour >= 21) { return true; }
            else { return false; }
          }
          // Friday
          case 5: {
            if(CurrentHour < 7 || CurrentHour >= 20) { return true; }
            else { return false; }
          }
          // Saturday
          case 6: {
            if(CurrentHour < 8 || CurrentHour >= 12) { return true; }
            else { return false; }
          }
          default: { return false; }
        }
      }
}