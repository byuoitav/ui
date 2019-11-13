import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { ControlGroup } from 'src/app/objects/control';

@Component({
  selector: 'app-help-info',
  templateUrl: './help-info.component.html',
  styleUrls: ['./help-info.component.scss']
})
export class HelpInfoComponent implements OnInit {
  info: string;

  constructor(
    public dialogRef: MatDialogRef<HelpInfoComponent>,
    @Inject(MAT_DIALOG_DATA) public cg: ControlGroup
  ) { }

  ngOnInit() {
  }

  close() {
    this.dialogRef.close(this.info);
  }

}
