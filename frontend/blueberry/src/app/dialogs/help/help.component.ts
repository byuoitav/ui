import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { ControlGroup } from 'src/app/objects/control';

@Component({
  selector: 'app-help',
  templateUrl: './help.component.html',
  styleUrls: ['./help.component.scss']
})
export class HelpComponent implements OnInit {

  constructor(
    public ref: MatDialogRef<HelpComponent>,
    @Inject(MAT_DIALOG_DATA) public cg: ControlGroup
  ) { }

  ngOnInit() {
  }

  cancel() {
    this.ref.close(false);
  }

  sendForHelp() {
    this.ref.close(true);
  }
}
