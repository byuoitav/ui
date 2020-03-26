import { Component } from "@angular/core";
import { Router, NavigationStart, NavigationEnd } from '@angular/router';

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"]
})
export class AppComponent {
  title = "dragonfruit";
  loading = false;

  constructor(private router: Router) {
    let vh = window.innerHeight * 0.01;

    document.documentElement.style.setProperty("--vh", `${vh}px`);

    window.addEventListener("resize", () => {
      let vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty("--vh", `${vh}px`);
    });

    this.router.events.subscribe(event => {
      if (event instanceof NavigationStart) {
        if (event.url.includes("/key/")) {
          this.loading = true;
        }
      }
      if (event instanceof NavigationEnd) {
        if (event.url.includes("/key/")) {
          this.loading = false;
        }
      }
    });

  }
}
