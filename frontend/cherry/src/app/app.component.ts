import { Component, ViewEncapsulation } from "@angular/core";
import { MatDialog } from '@angular/material';
import { trigger, transition, animate } from "@angular/animations";
import { Http } from "@angular/http";
import { Output } from '../objects/status.objects';
import { BFFService } from '../services/bff.service';
import { HelpDialog } from "./dialogs/help.dialog";
import { MobileControlComponent } from "./dialogs/mobilecontrol/mobilecontrol.component";


const HIDDEN = "hidden";
const QUERY = "query";
const LOADING = "indeterminate";
const BUFFER = "buffer";

@Component({
  selector: "cherry",
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  encapsulation: ViewEncapsulation.None,
  animations: [
    trigger("delay", [
      transition(":enter", [animate(500)]),
      transition(":leave", [animate(500)])
    ])
  ]
})
export class AppComponent {
  public loaded: boolean;
  public unlocking = false;
  public progressMode: string = QUERY;
  public power: boolean;
  public selectedTabIndex: number;

  constructor(
    public bff: BFFService,
    public dialog: MatDialog
    // private http: Http
  ) {
    this.loaded = false;
    this.power = true
  }

  public isPoweredOff(): boolean {
    if (!this.loaded) {
      return true;
    }
    // return Output.isPoweredOn(this.data.panel.preset.displays);
    return false;
  }

  public unlock() {
    this.unlocking = true;
    this.progressMode = QUERY;
    const room = this.bff.getRoom;
    
    

    // this.command.powerOnDefault(this.data.panel.preset).subscribe(success => {
    //   if (!success) {
    //     this.reset();
    //     console.warn("failed to turn on");
    //   } else {
    //     // switch direction of loading bar
    //     this.progressMode = LOADING;

    //     this.reset();
    //   }
    // });
  }

  public powerOff() {
    this.progressMode = QUERY;

    // this.command.powerOff(this.data.panel.preset).subscribe(success => {
    //   if (!success) {
    //     console.warn("failed to turn off");
    //   } else {
    //     this.reset();
    //   }
    // });
  }

  private reset() {
    // select displays tab
    this.selectedTabIndex = 0;

    // // reset mix levels to 100
    // this.data.panel.preset.audioDevices.forEach(a => (a.mixlevel = 100));

    // // reset masterVolume level
    // this.data.panel.preset.masterVolume = 30;

    // // reset masterMute
    // this.data.panel.preset.masterMute = false;

    // stop showing progress bar
    this.unlocking = false;
  }

  public openHelpDialog() {
    const dialogRef = this.dialog.open(HelpDialog, {
      width: "70vw",
      disableClose: true
    });
  }

  public openMobileControlDialog() {
    const dialogRef = this.dialog.open(MobileControlComponent, {
      width: "70vw",
      height: "52.5vw"
    });
  }

  powerStatus() {
    return this.power;
  }

  setPower() {
    this.power = !this.power;
  }
}