import {AfterViewInit, Component, ElementRef, HostListener, OnInit, ViewChild} from '@angular/core';
import {Contact} from './contact.model';
import {ContactService} from './contact.service';
import {isNullOrUndefined} from 'util';

@Component({
  selector: 'app-contact-list',
  templateUrl: './contact-list.component.html',
  styleUrls: ['./contact-list.component.scss']
})
export class ContactListComponent implements OnInit, AfterViewInit {

  prefixes: string[] = ["A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"];
  contacts: Contact[];
  selectedContact: Contact;
  listHeight = 500;

  constructor(private contactService: ContactService) { }

  ngOnInit() {
    this.contactService.getAllContacts().subscribe(contacts => {
      this.contacts = contacts;
    });
  }

  filterContactsStartwith(prefix: string): Contact[] {
    if (isNullOrUndefined(this.contacts)) {
      return this.contacts;
    }
    let filteredContacts: Contact[] = new Array(0);
    for (let contact of this.contacts) {
      if (contact.lastname.startsWith(prefix)) {
        filteredContacts.push(contact);
      }
    }
    return filteredContacts;
  }

  contactClicked(contact: Contact) {
    this.contactService.getContact(contact.id).subscribe(c => {
      this.selectedContact = c;
    });
  }

  @HostListener('window:resize', ['$event'])
  onResize(event) {
    this.listHeight = event.target.innerHeight;
  }

  ngAfterViewInit(): void {
    window.dispatchEvent(new Event('resize'));
  }
}
