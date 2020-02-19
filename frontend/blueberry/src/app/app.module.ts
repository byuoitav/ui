import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { HttpModule } from "@angular/http";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import {
  MatSliderModule,
  MatIconModule,
  MatButtonModule,
  MatMenuModule,
  MatDialogModule,
  MatGridListModule,
  MatChipsModule,
  MatProgressSpinnerModule,
  MatSlideToggleModule,
  MatToolbarModule,
  MatTabsModule,
  MAT_DIALOG_DATA
} from "@angular/material";
import { UiSwitchModule } from "ngx-ui-switch";
import "hammerjs";

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './components/app.component';
import { BFFService } from './services/bff.service';
import { HomeComponent } from './components/home/home.component';
import { WheelComponent } from './components/wheel/wheel.component';
import { VolumeComponent } from './components/volume/volume.component';
import { ProjectorComponent } from './components/projector/projector.component';
import { AudioComponent } from './components/audio/audio.component';

import { MobileComponent } from './dialogs/mobile/mobile.component';

import { NgxQRCodeModule } from 'ngx-qrcode2';
import { SharingComponent } from './dialogs/sharing/sharing.component';
import { HelpComponent } from './dialogs/help/help.component';
import { ActivityButtonComponent } from './components/activity-button/activity-button.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    WheelComponent,
    VolumeComponent,
    ProjectorComponent,
    AudioComponent,
    MobileComponent,
    SharingComponent,
    HelpComponent,
    ActivityButtonComponent,
  ],
  imports: [
    BrowserModule,
    HttpModule,
    BrowserAnimationsModule,
    MatSliderModule,
    MatIconModule,
    MatButtonModule,
    MatMenuModule,
    MatDialogModule,
    MatGridListModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatSlideToggleModule,
    UiSwitchModule,
    AppRoutingModule,
    MatToolbarModule,
    MatTabsModule,
    NgxQRCodeModule
  ],
  entryComponents: [
    MobileComponent,
    HelpComponent,
    SharingComponent
  ],
  providers: [
    BFFService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
