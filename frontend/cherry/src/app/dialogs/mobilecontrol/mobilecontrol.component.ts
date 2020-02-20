import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef } from "@angular/material";


@Component({
  selector: 'mobilecontrol',
  templateUrl: './mobilecontrol.component.html',
  styleUrls: ['./mobilecontrol.component.scss']
})
export class MobileControlComponent implements OnInit {
  public value: string;
  public elementType: 'url';


  constructor(
    public ref: MatDialogRef<MobileControlComponent>
  ) {
    this.value = "rooms.stg.byu.edu/key/111111"
   }

  ngOnInit() {
  }

  public cancel() {
    this.ref.close()
  }
}
