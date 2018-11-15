import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {Contact} from './contact.model';
import {HttpClient} from '@angular/common/http';

@Injectable()
export class ContactService {
  constructor(private http: HttpClient) {}

  getAllContacts(): Observable<Contact[]> {
    return this.http.get<Contact[]>("http://localhost:8080/api/contact");
  }

  getContact(id: number): Observable<Contact> {
    return this.http.get<Contact>("http://localhost:8080/api/contact/" + id);
  }

}
