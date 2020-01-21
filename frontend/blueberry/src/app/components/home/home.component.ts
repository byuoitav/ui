import { Component, OnInit, Input } from '@angular/core';
import { BFFService, RoomRef } from '../../services/bff.service';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss', '../../colorscheme.scss']
})
export class HomeComponent implements OnInit {
  @Input() room: RoomRef;

  constructor(public bff: BFFService) { }

  ngOnInit() {
  }

}
