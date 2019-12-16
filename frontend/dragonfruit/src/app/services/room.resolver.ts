import { Injectable } from "@angular/core";
import {
  Router,
  Resolve,
  RouterStateSnapshot,
  ActivatedRouteSnapshot
} from "@angular/router";
import { Observable, of, EMPTY, Subject, BehaviorSubject } from "rxjs";
import { takeUntil } from "rxjs/operators";

import { BFFService, RoomRef } from "./bff.service";
import { Room } from "../objects/control";

@Injectable({
  providedIn: "root"
})
export class RoomResolver implements Resolve<RoomRef> {
  constructor(private bff: BFFService, private router: Router) {}

  resolve(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<RoomRef> | Observable<never> {
    const key = route.paramMap.get("key");
    const unsubscribe = new Subject();

    const roomRef = this.bff.getRoom(key);

    return new Observable(observer => {
      roomRef
        .subject()
        .pipe(takeUntil(unsubscribe))
        .subscribe(
          val => {
            if (val instanceof Room) {
              observer.next(roomRef);
              observer.complete();
              unsubscribe.complete();
            }
          },
          err => {
            this.router.navigate(["/login"], {
              queryParams: {
                error: "Invalid room control key."
              },
              queryParamsHandling: "merge"
            });

            observer.error(err);
            unsubscribe.complete();
          }
        );

      return { unsubscribe() {} };
    });
  }
}
