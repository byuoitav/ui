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
import { Room, isRoom } from "../../../../objects/control";

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
            if (isRoom(val)) {
              if (val.controlGroups[val.selectedControlGroup].poweredOn) {
                observer.next(roomRef);
                observer.complete();
                unsubscribe.complete();
              } else {
                roomRef.setPower(true);
              }
            }
          },
          err => {
            if (err.code === 4000) {
              this.router.navigate(["/login"], {
                queryParams: {
                  error: err.reason
                },
                queryParamsHandling: "merge"
              });
            } else {
              this.router.navigate(["/login"]);
            }

            observer.error(err);
            unsubscribe.complete();
          }
        );

      return { unsubscribe() {} };
    });
  }
}
