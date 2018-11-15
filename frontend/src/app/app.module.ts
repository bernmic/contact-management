import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {FlexLayoutModule} from '@angular/flex-layout';
import {
  MatButtonModule,
  MatCardModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule, MatInputModule,
  MatListModule,
  MatToolbarModule
} from '@angular/material';
import { ContactListComponent } from './contact-list/contact-list.component';
import {ContactService} from './contact-list/contact.service';
import {HttpClientModule} from '@angular/common/http';
import {ContactDetailComponent} from './contact-list/contact-detail.component';
import { MessageboxComponent } from './messagebox/messagebox.component';
import {ContactEditComponent} from './contact-list/contact-edit.component';
import {FormsModule} from '@angular/forms';
import {NgxVcardModule} from 'ngx-vcard';

@NgModule({
  declarations: [
    AppComponent,
    ContactDetailComponent,
    ContactEditComponent,
    ContactListComponent,
    MessageboxComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    FlexLayoutModule,
    MatButtonModule,
    MatCardModule,
    MatDialogModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatToolbarModule,
    NgxVcardModule
  ],
  providers: [
    ContactService
  ],
  entryComponents: [
    MessageboxComponent,
    ContactEditComponent
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
