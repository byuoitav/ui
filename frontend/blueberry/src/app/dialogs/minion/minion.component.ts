import { Component, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';


@Component({
  selector: 'app-minion',
  templateUrl: './minion.component.html',
  styleUrls: ['./minion.component.scss']
})
export class MinionComponent implements OnInit {

    constructor(public ref: MatDialogRef<MinionComponent>) { }

  ngOnInit() {
  }

  cancel = () => {
    this.ref.close()
  }

}
