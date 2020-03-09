import { Component, OnInit, Inject } from '@angular/core';
import { ControlGroup } from '../../../../../objects/control';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef, BFFService } from 'src/app/services/bff.service';

@Component({
  selector: 'app-mobile',
  templateUrl: './mobile.component.html',
  styleUrls: ['./mobile.component.scss']
})
export class MobileComponent implements OnInit {
  public elementType: 'url';

  constructor(public ref: MatDialogRef<MobileComponent>, @Inject(MAT_DIALOG_DATA) public data: ControlGroup) {
  }

  ngOnInit() {
    this.ref.disableClose = true;
  }

  cancel = () => {
    this.ref.close();
  }

  getQRCode() {
    return "http://" + this.data.controlInfo.url + "/key/" + this.data.controlInfo.key;
  }
}
