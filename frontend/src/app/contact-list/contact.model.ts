export class Phone {
  constructor(
    public id: number,
    public name: string,
    public number: string,
    public contact_id: number
  ) {}
}

export class Contact {
  constructor(
    public id: number,
    public firstname: string,
    public lastname: string,
    public company: string,
    public address1: string,
    public address2: string,
    public zipcode: string,
    public city: string,
    public country: string,
    public tag: string,
    public email: string,
    public web: string,
    public birthday: Date,
    public phones: Phone[]
) {}
}
