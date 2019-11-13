import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { ControlGroup } from 'src/app/objects/control';
import { BFFService } from 'src/app/services/bff.service';

@Component({
  selector: 'app-present',
  templateUrl: './present.component.html',
  styleUrls: ['./present.component.scss']
})
export class PresentComponent implements OnInit {
  @AngularInput() cg: ControlGroup;

  constructor(
    private bff: BFFService
  ) { }

  ngOnInit() {
  }

}
