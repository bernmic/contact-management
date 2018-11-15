import {Component, Inject} from '@angular/core';
import {Contact} from './contact.model';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

@Component({
  selector: 'app-contact-edit',
  templateUrl: './contact-edit.component.html',
  styleUrls: ['./contact-edit.component.scss']
})
export class ContactEditComponent {
  contact: Contact;

  constructor(public thisDialogRef: MatDialogRef<ContactEditComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.contact = data.contact;
    if (this.contact === null) {
      this.contact = new Contact(
        0,
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        null,
        null
      );
    }
  }

  onCloseConfirm() {
    this.thisDialogRef.close(this.contact);
  }
  onCloseCancel() {
    this.thisDialogRef.close(null);
  }
}
