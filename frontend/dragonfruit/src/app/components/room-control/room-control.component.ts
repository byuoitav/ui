import { Component, OnInit, HostListener } from '@angular/core';
import { BFFService } from 'src/app/services/bff.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Preset } from 'src/app/objects/database';
import { ControlGroup, CONTROL_TAB, AUDIO_TAB, PRESENT_TAB, HELP_TAB } from 'src/app/objects/control';
import { MatTabChangeEvent, MatTab } from '@angular/material';

@Component({
  selector: 'app-room-control',
  templateUrl: './room-control.component.html',
  styleUrls: ['./room-control.component.scss']
})
export class RoomControlComponent implements OnInit {
  controlGroup: ControlGroup;
  groupIndex: string;
  roomID: string;

  tabPosition = 'below';
  selectedTab: number;

  @HostListener('window:resize', ['$event'])
  onResize(event) {
    if (window.innerWidth >= 768) {
      this.tabPosition = 'above';
    } else {
      this.tabPosition = 'below';
    }
  }

  constructor(
    public bff: BFFService,
    public route: ActivatedRoute,
    private router: Router
  ) {
    this.route.params.subscribe(params => {
      this.roomID = params['id'];
      this.groupIndex = params['index'];
      this.selectedTab = +params['tabName'];
      if (this.bff.room === undefined) {
        this.bff.connectToRoom(this.roomID);

        this.bff.done.subscribe(e => {
          this.controlGroup = this.bff.room.controlGroups[this.groupIndex];
          if (this.controlGroup.id === 'Third') {
            this.setExtraDisplays();
          }
        });
      } else {
        this.controlGroup = this.bff.room.controlGroups[this.groupIndex];
        if (this.controlGroup.id === 'Third') {
          this.setExtraDisplays();
        }
      }

      // this.bff.done.subscribe(() => {
      //   this.controlGroup = this.bff.room.controlGroups[this.groupIndex];
      //   if (this.bff.room.selectedGroup === undefined) {
      //     this.bff.room.selectedGroup = this.controlGroup.name;
      //   }
      // });
    });
  }

  // for testing only
  setExtraDisplays() {
    this.controlGroup.displays = [
      {
        id: '111 - A',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D1',
            icon: 'tv'
          }
        ]
      },
      {
        id: '111 - B',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D2',
            icon: 'tv'
          }
        ]
      },
      {
        id: '111 - C',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D3',
            icon: 'tv'
          }
        ]
      },
      {
        id: '21 - A',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D4',
            icon: 'tv'
          },
          {
            name: 'D5',
            icon: 'tv'
          }
        ]
      },
      {
        id: '21 - B',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D6',
            icon: 'tv'
          }
        ]
      },
      {
        id: '31 - A',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D11',
            icon: 'videocam'
          },
          {
            name: 'D12',
            icon: 'videocam'
          },
          {
            name: 'D13',
            icon: 'videocam'
          }
        ]
      },
      {
        id: '31 - B',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D6',
            icon: 'tv'
          }
        ]
      },
      {
        id: '22 - A',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D7',
            icon: 'tv'
          },
          {
            name: 'D8',
            icon: 'tv'
          }
        ]
      },
      {
        id: '22 - B',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D9',
            icon: 'tv'
          },
          {
            name: 'D10',
            icon: 'tv'
          }
        ]
      },
      {
        id: '3 wide',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D11',
            icon: 'tv'
          },
          {
            name: 'D12',
            icon: 'tv'
          },
          {
            name: 'D13',
            icon: 'tv'
          }
        ]
      },
      {
        id: '4+',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D14',
            icon: 'tv'
          },
          {
            name: 'D15',
            icon: 'tv'
          },
          {
            name: 'D16',
            icon: 'tv'
          },
          {
            name: 'D17',
            icon: 'tv'
          },
          {
            name: 'D18',
            icon: 'tv'
          },
          {
            name: 'D19',
            icon: 'tv'
          }
        ]
      },
      {
        id: 'End 2 - A',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D20',
            icon: 'tv'
          }
        ]
      },
      {
        id: 'End 2 - B',
        input: this.roomID + '-VIA1',
        blanked: false,
        outputs: [
          {
            name: 'D21',
            icon: 'tv'
          }
        ]
      }
    ];
  }

  ngOnInit() {
    if (window.innerWidth >= 768) {
      this.tabPosition = 'above';
    } else {
      this.tabPosition = 'below';
    }
  }

  goBack = () => {
    this.router.navigate(['/room/' + this.roomID]);
  }

  tabChange(index: number) {
    this.selectedTab = index;
    const currentURL = decodeURI(window.location.pathname);
    const newURL = currentURL.substr(0, currentURL.lastIndexOf('/') + 1) + (this.selectedTab);
    this.router.navigate([newURL]);
  }
}
