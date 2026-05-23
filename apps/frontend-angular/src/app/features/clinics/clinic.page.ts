import { Component, inject, signal } from '@angular/core';

import { ClinicProfile } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-clinic-page',
  templateUrl: './clinic.page.html',
  styleUrl: './clinic.page.css',
})
export class ClinicPage {
  private readonly api = inject(ApiService);
  readonly clinic = signal<ClinicProfile | null>(null);

  constructor() {
    this.api.clinicCurrent().subscribe({ next: (profile) => this.clinic.set(profile), error: () => {} });
  }
}
