import { Component, OnInit, Inject } from '@angular/core';
import { ControlGroup } from 'src/app/objects/control';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

@Component({
  selector: 'app-mobile',
  templateUrl: './mobile.component.html',
  styleUrls: ['./mobile.component.scss']
})
export class MobileComponent implements OnInit {
  public qrcode: string;
  public url: string;
  public key: string;
  public elementType: 'url';
  public cg: ControlGroup;

  constructor(public dialogRef: MatDialogRef<MobileComponent>, @Inject(MAT_DIALOG_DATA) public data: ControlGroup) {
    this.cg = data;
    this.url = "rooms-stg.byu.edu";
    this.key = "103236";
    this.qrcode = "http://" + this.url + "/key/" + this.key;
  }

  ngOnInit() {
    this.dialogRef.disableClose = true;
  }

  cancel = () => {
    this.dialogRef.close();
  }
}
