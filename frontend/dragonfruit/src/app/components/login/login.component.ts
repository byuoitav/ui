import { Component, OnInit, EventEmitter } from "@angular/core";
import { MatBottomSheet } from "@angular/material";
import { NumpadComponent } from "../../dialogs/numpad/numpad.component";
import { Router } from "@angular/router";
import { BFFService } from "src/app/services/bff.service";

@Component({
  selector: "app-login",
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.scss"]
})
export class LoginComponent implements OnInit {
  roomCode = "";
  keyboardEmitter: EventEmitter<string>;

  constructor(
    private bottomSheet: MatBottomSheet,
    private router: Router,
    private bff: BFFService
  ) {
    this.keyboardEmitter = new EventEmitter<string>();
    this.keyboardEmitter.subscribe(s => {
      this.roomCode = s;
    });
  }

  ngOnInit() {}

  showNumpad = () => {
    // this.bottomSheet
    //   .open(NumpadComponent, {
    //     data: this.keyboardEmitter,
    //     backdropClass: "keyboard-bg"
    //   })
    //   .afterDismissed()
    //   .subscribe(result => {
    //     if (result !== undefined) {
    //       console.log(
    //         "redirecting using the following room code:",
    //         this.roomCode
    //       );
    //       this.goToRoomControl();
    //     }
    //   });
  };

  codeKeyUp(event, index) {
    console.log(event);
    if (event.key === "Backspace") {
      if (index > 0) {
        const elementName = "codeKey" + (index - 1);
        document.getElementById(elementName).focus();
      }
      return;
    }
    if (index >= 0 && index < 5) {
      const elementName = "codeKey" + (index + 1);
      document.getElementById(elementName).focus();
    }
  }

  getCodeChar = (index: number): string => {
    if (this.roomCode.length > index) {
      return this.roomCode.charAt(index);
    }

    return "";
  };

  goToRoomControl = () => {
    // TODO: actually do something with the room code
    this.bff.connectToRoom(this.roomCode);
    // switch (this.roomCode) {
    //   case '1101': {
    this.bff.done.subscribe(e => {
      this.router.navigate(["/key/" + this.roomCode + "/room/" + this.bff.room.id]);
    });
    //     break;
    //   }
    //   case '1102': {
    //     this.router.navigate(['/room/ITB-1101/group/0/tab/0']);
    //     break;
    //   }
    //   default:
    //     break;
    // }
  };
}
