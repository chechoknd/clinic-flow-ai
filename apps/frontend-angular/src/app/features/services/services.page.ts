import { Component, inject, signal } from '@angular/core';

import { ClinicServiceItem } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-services-page',
  templateUrl: './services.page.html',
  styleUrl: './services.page.css',
})
export class ServicesPage {
  private readonly api = inject(ApiService);
  readonly services = signal<ClinicServiceItem[]>([]);

  constructor() {
    this.api.services().subscribe({ next: (response) => this.services.set(response.data), error: () => {} });
  }
}
