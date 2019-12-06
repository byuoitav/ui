import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { RoomControlComponent } from './components/room-control/room-control.component';
import { SelectionComponent } from './components/selection/selection.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/login',
    pathMatch: 'full'
  },
  {
    path: 'login',
    component: LoginComponent
  },
  {
    path: 'key/:key/room/:id',
    component: SelectionComponent
  },
  {
    path: 'key/:key/room/:id/group/:index/tab/:tabName',
    component: RoomControlComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
