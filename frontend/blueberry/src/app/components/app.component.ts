import { Component } from '@angular/core';
import { BFFService, RoomRef } from '../services/bff.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  public roomRef: RoomRef;

  constructor(public bff: BFFService) {
    this.roomRef = this.bff.getRoom("825943");
  }

  unlock = () => {
    
  }
}
