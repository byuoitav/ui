import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";

import { AppComponent } from "./app.component";
import { LoginComponent } from "./components/login/login.component";
import { RoomControlComponent } from "./components/room-control/room-control.component";
import { SelectionComponent } from "./components/selection/selection.component";

const routes: Routes = [
  {
    path: "",
    redirectTo: "/login",
    pathMatch: "full"
  },
  {
    path: "",
    component: AppComponent,
    children: [
      {
        path: "login",
        component: LoginComponent
      },
      {
        path: "key/:key",
        resolve: {
          roomRef: RoomResolver
        },
        children: [
          {
            path: "",
            component: SelectionComponent
          },
          {
            path: ":groupid",
            component: RoomControlComponent
          }
        ]
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
