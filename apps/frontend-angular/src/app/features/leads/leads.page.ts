import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { Lead } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-leads-page',
  imports: [RouterLink],
  templateUrl: './leads.page.html',
  styleUrl: './leads.page.css',
})
export class LeadsPage {
  private readonly api = inject(ApiService);

  readonly statuses = ['Nuevo', 'Contactado', 'Interesado', 'Agendado', 'No Respondio', 'Perdido', 'Convertido'];
  readonly selectedStatus = signal('Nuevo');
  readonly leads = signal<Lead[]>([]);
  readonly filteredLeads = computed(() =>
    this.leads().filter((lead) => lead.status === this.selectedStatus()),
  );

  constructor() {
    this.api.leads().subscribe({ next: (response) => this.leads.set(response.data), error: () => {} });
  }
}
