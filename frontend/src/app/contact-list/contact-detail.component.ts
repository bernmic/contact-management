import {Component, Input} from '@angular/core';
import {Contact} from './contact.model';
import {MatDialog} from '@angular/material';
import {MessageboxComponent} from '../messagebox/messagebox.component';
import {ContactEditComponent} from './contact-edit.component';
import {VCard} from 'ngx-vcard';
import {isNullOrUndefined} from 'util';

@Component({
  selector: 'app-contact-detail',
  templateUrl: './contact-detail.component.html',
  styleUrls: ['./contact-detail.component.scss']
})
export class ContactDetailComponent {
  @Input() contact: Contact;

  constructor(private dialog: MatDialog) {}

  confirmDelete() {
    let dialogRef = this.dialog.open(MessageboxComponent, {
      width: '400px',
      data: {
        title: "Delete Contact",
        content: "Delete contact " + this.contact.firstname + " " + this.contact.lastname + ". Are you sure?"
      }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log("Delete contact " + this.contact.firstname + " " + this.contact.lastname);
      }
    });
  }

  editContact() {
    let dialogRef = this.dialog.open(ContactEditComponent, {
      width: '400px',
      data: {
        title: "Edit Contact",
        contact: this.contact
      }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log("Edit contact " + result.firstname + " " + result.lastname);
      }
    });
  }

  addContact() {
    let dialogRef = this.dialog.open(ContactEditComponent, {
      width: '400px',
      data: {
        title: "Add Contact",
        contact: null
      }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        console.log("Add contact " + result.firstname + " " + result.lastname);
      }
    });
  }

  exportContact(): VCard {
    const phonenumbers: string[] = [];
    if (!isNullOrUndefined(this.contact.phones)) {
      for (let phone of this.contact.phones) {
        phonenumbers.push(phone.number);
      }
    }
    const vcard: VCard = {
      name: {
        firstNames: this.contact.firstname,
        lastNames: this.contact.lastname
      },
      organization: this.contact.company,
      address: [
        {
          street: this.contact.address1,
          locality: this.contact.city,
          postalCode: this.contact.zipcode,
          countryName: this.contact.country
        }
      ],
      email: [this.contact.email],
      url: this.contact.web,
      birthday: this.contact.birthday,
      telephone: phonenumbers
    }
    return vcard;
  }
}
