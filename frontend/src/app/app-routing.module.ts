import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {ContactListComponent} from './contact-list/contact-list.component';

const routes: Routes = [
  { path: '', redirectTo: 'list', pathMatch: 'full' },
  { path: 'list', component: ContactListComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
