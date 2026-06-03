import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';

import { I18nService } from '../../core/i18n/i18n.service';
import { ClinicServiceItem, Lead, LeadStatus } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-leads-page',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './leads.page.html',
  styleUrl: './leads.page.css',
})
export class LeadsPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  readonly i18n = inject(I18nService);

  readonly statuses: LeadStatus[] = [
    'Nuevo',
    'Contactado',
    'Interesado',
    'Agendado',
    'No Respondio',
    'Perdido',
    'Convertido',
  ];
  readonly selectedStatus = signal<LeadStatus>('Nuevo');
  readonly leads = signal<Lead[]>([]);
  readonly services = signal<ClinicServiceItem[]>([]);
  readonly selectedLead = signal<Lead | null>(null);
  readonly createModalOpen = signal(false);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly statusCounts = computed(() => {
    const counts = this.statuses.reduce(
      (acc, status) => ({ ...acc, [status]: 0 }),
      {} as Record<LeadStatus, number>,
    );
    this.leads().forEach((lead) => {
      counts[lead.status] = (counts[lead.status] ?? 0) + 1;
    });
    return counts;
  });
  readonly selectedStatusCount = computed(() => this.statusCounts()[this.selectedStatus()] ?? 0);
  readonly filteredLeads = computed(() =>
    this.leads().filter((lead) => lead.status === this.selectedStatus()),
  );

  readonly createForm = this.fb.nonNullable.group({
    full_name: ['', [Validators.required, Validators.minLength(3)]],
    phone: ['', [Validators.required, Validators.pattern(/^\+[1-9]\d{7,14}$/)]],
    service_id: [''],
    status: ['Nuevo' as LeadStatus, Validators.required],
    source: ['whatsapp', Validators.required],
    next_action_at: [''],
    notes: [''],
  });

  readonly updateForm = this.fb.nonNullable.group({
    status: ['Nuevo' as LeadStatus, Validators.required],
    next_action_at: [''],
    note: [''],
  });

  constructor() {
    this.loadInitialData();
  }

  openCreateModal(): void {
    this.createModalOpen.set(true);
    this.success.set(null);
    this.error.set(null);
  }

  closeCreateModal(): void {
    if (this.saving()) {
      return;
    }
    this.createModalOpen.set(false);
  }

  selectStatus(status: LeadStatus): void {
    this.selectedStatus.set(status);
    const current = this.selectedLead();
    if (current && current.status !== status) {
      this.selectedLead.set(null);
    }
  }

  selectLead(lead: Lead): void {
    this.selectedLead.set(lead);
    this.success.set(null);
    this.error.set(null);
    this.updateForm.reset({
      status: lead.status,
      next_action_at: this.toLocalDateTimeValue(lead.next_action_at),
      note: '',
    });
  }

  createLead(): void {
    if (this.createForm.invalid) {
      this.createForm.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    this.error.set(null);
    this.success.set(null);
    const value = this.createForm.getRawValue();

    this.api
      .createLead({
        full_name: value.full_name.trim(),
        phone: value.phone.trim(),
        service_id: value.service_id || undefined,
        status: value.status,
        source: value.source,
        notes: value.notes.trim() || undefined,
        next_action_at: this.toApiDateTime(value.next_action_at),
      })
      .subscribe({
        next: (lead) => {
          this.leads.update((items) => [lead, ...items]);
          this.selectedStatus.set(lead.status);
          this.selectLead(lead);
          this.createForm.reset({
            full_name: '',
            phone: '',
            service_id: '',
            status: 'Nuevo',
            source: 'whatsapp',
            next_action_at: '',
            notes: '',
          });
          this.createModalOpen.set(false);
          this.success.set(this.i18n.t('leads.created'));
          this.saving.set(false);
        },
        error: () => {
          this.error.set(this.i18n.t('leads.createError'));
          this.saving.set(false);
        },
      });
  }

  updateLead(): void {
    const lead = this.selectedLead();
    if (!lead) {
      return;
    }
    if (this.updateForm.invalid) {
      this.updateForm.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    this.error.set(null);
    this.success.set(null);
    const value = this.updateForm.getRawValue();

    this.api
      .updateLead(lead.id, {
        status: value.status,
        note: value.note.trim() || undefined,
        next_action_at: this.toApiDateTime(value.next_action_at),
      })
      .subscribe({
        next: () => {
          const updated: Lead = {
            ...lead,
            status: value.status,
            next_action_at: this.toApiDateTime(value.next_action_at),
          };
          this.leads.update((items) => items.map((item) => (item.id === lead.id ? updated : item)));
          this.selectedStatus.set(value.status);
          this.selectLead(updated);
          this.success.set(this.i18n.t('leads.updated'));
          this.saving.set(false);
        },
        error: () => {
          this.error.set(this.i18n.t('leads.updateError'));
          this.saving.set(false);
        },
      });
  }

  statusLabel(status: LeadStatus): string {
    return this.i18n.t(`lead.status.${status}`);
  }

  formatDate(value: string | undefined): string {
    if (!value) {
      return this.i18n.t('leads.noAction');
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return this.i18n.t('leads.noAction');
    }
    return date.toLocaleString(this.i18n.language() === 'en' ? 'en-US' : 'es-CO', {
      day: '2-digit',
      month: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }

  private loadInitialData(): void {
    this.loading.set(true);
    this.api.services().subscribe({
      next: (response) => this.services.set(response.data),
      error: () => this.error.set(this.i18n.t('leads.servicesError')),
    });

    this.api.leads().subscribe({
      next: (response) => {
        this.leads.set(response.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('leads.loadError'));
        this.loading.set(false);
      },
    });
  }

  private toApiDateTime(value: string | undefined): string | undefined {
    if (!value) {
      return undefined;
    }

    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
  }

  private toLocalDateTimeValue(value: string | undefined): string {
    if (!value) {
      return '';
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return '';
    }

    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(0, 16);
  }
}
