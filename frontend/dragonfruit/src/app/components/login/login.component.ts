import { Component, OnInit, EventEmitter } from "@angular/core";
import { MatBottomSheet } from "@angular/material";
import { NumpadComponent } from "../../dialogs/numpad/numpad.component";
import {
  Router,
  Event,
  NavigationStart,
  NavigationEnd,
  NavigationError,
  NavigationCancel
} from "@angular/router";
import { BFFService } from "src/app/services/bff.service";

@Component({
  selector: "app-login",
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.scss"]
})
export class LoginComponent implements OnInit {
  roomCode = "";
  loggingIn = false;
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

    // subscribe to routing events so that we can
    // show the loading indicator when logging in
    this.router.events.subscribe(event => {
      switch (true) {
        case event instanceof NavigationStart:
          this.loggingIn = true;
          break;
        case event instanceof NavigationEnd:
        case event instanceof NavigationCancel:
        case event instanceof NavigationError:
          this.loggingIn = false;
          break;
        default:
      }
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
    console.log("hello");
    // TODO: actually do something with the room code
    // this.bff.connectToRoom(this.roomCode);
    this.bff.getRoom(this.roomCode);
    // switch (this.roomCode) {
    //   case '1101': {
    /*
    this.bff.done.subscribe(e => {
      console.log("world");
      if (this.bff.room.selectedControlGroup) {
        this.router.navigate([
          "/key/" +
            this.roomCode +
            "/room/" +
            this.bff.room.id +
            "/group/" +
            this.bff.room.selectedControlGroup +
            "/tab/0"
        ]);
      } else {
        this.router.navigate([
          "/key/" + this.roomCode + "/room/" + this.bff.room.id
        ]);
      }
    });
    */
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
