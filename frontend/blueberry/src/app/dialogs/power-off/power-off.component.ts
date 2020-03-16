import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { RoomRef, BFFService } from 'src/app/services/bff.service';

@Component({
  selector: 'app-power-off',
  templateUrl: './power-off.component.html',
  styleUrls: ['./power-off.component.scss']
})
export class PowerOffComponent implements OnInit {

  constructor(public ref: MatDialogRef<PowerOffComponent>,
    @Inject(MAT_DIALOG_DATA) public roomRef: RoomRef,
    public bff: BFFService) { }

  ngOnInit() {
  }

  cancel() {
    this.ref.close("cancel");
  }

  shutDown(option: string) {
    this.ref.close(option);
  }
}
