import { Component, OnInit, Input as AngularInput } from '@angular/core';
import { BFFService } from 'src/app/services/bff.service';
import { ControlGroup } from '../../../../../objects/control';

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
