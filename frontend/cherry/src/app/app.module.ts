import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { NgModule } from '@angular/core';
import { HttpModule } from "@angular/http";
import { NgxQRCodeModule } from 'ngx-qrcode2';


import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { MatIconModule, MatDialogModule, MatTabsModule, MatButtonModule, MatSliderModule, MatProgressBarModule, MatProgressSpinnerModule } from "@angular/material";
import { BFFService } from '../services/bff.service';
import { DisplayComponent } from './components/display/display.component';
import { AudioControlComponent } from './components/audiocontrol/audiocontrol.component';
import { VolumeComponent } from './components/volume/volume.component';
import { ScreenControlComponent } from './components/screencontrol/screencontrol.component';
import { MobileControlComponent } from './dialogs/mobilecontrol/mobilecontrol.component';
import "hammerjs";
import { HelpDialog } from './dialogs/help.dialog';
import { LockScreenAudioComponent } from './components/lockscreenaudio/lockscreenaudio.component';
import { LockScreenScreenControlComponent } from './components/lockscreenscreencontrol/lockscreenscreencontrol.component';


@NgModule({
  declarations: [
    AppComponent,
    DisplayComponent,
    AudioControlComponent,
    VolumeComponent,
    ScreenControlComponent,
    MobileControlComponent,
    HelpDialog,
    LockScreenAudioComponent,
    LockScreenScreenControlComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpModule,
    MatSliderModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    NgxQRCodeModule,
    MatIconModule,
    AppRoutingModule,
    MatDialogModule,
    MatTabsModule,
    MatButtonModule
  ],
  providers: [
    BFFService
  ],
  entryComponents: [MobileControlComponent, HelpDialog],
  bootstrap: [AppComponent]
})
export class AppModule { }
