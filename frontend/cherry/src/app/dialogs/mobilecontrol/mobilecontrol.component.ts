import { Component, OnInit, Inject } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from "@angular/material";


@Component({
  selector: 'mobilecontrol',
  templateUrl: './mobilecontrol.component.html',
  styleUrls: ['./mobilecontrol.component.scss']
})
export class MobileControlComponent implements OnInit {
  url: string;
  key: string;
  public qrCode: string;
  public elementType: 'url';


  constructor(
    public ref: MatDialogRef<MobileControlComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {
      url: string;
      key: string;
    }
  ) {
    this.qrCode = this.data.url + "/key/" + this.data.key;
   }

  ngOnInit() {
  }

  public cancel() {
    this.ref.close()
  }
}
