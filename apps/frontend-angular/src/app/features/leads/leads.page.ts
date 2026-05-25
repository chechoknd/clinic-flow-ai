import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';

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
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly filteredLeads = computed(() =>
    this.leads().filter((lead) => lead.status === this.selectedStatus()),
  );

  statusCount(status: LeadStatus): number {
    return this.leads().filter((lead) => lead.status === status).length;
  }

  urgencyLabel(lead: Lead): string {
    if (this.isOverdue(lead.next_action_at)) {
      return 'Vencido';
    }
    if (this.isToday(lead.next_action_at)) {
      return 'Hoy';
    }
    return '';
  }


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
          this.success.set('Lead creado correctamente.');
          this.saving.set(false);
        },
        error: () => {
          this.error.set('No fue posible crear el lead. Revisa los datos e intenta de nuevo.');
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
          this.success.set('Lead actualizado correctamente.');
          this.saving.set(false);
        },
        error: () => {
          this.error.set('No fue posible actualizar el lead. Intenta de nuevo.');
          this.saving.set(false);
        },
      });
  }

  private loadInitialData(): void {
    this.loading.set(true);
    this.api.services().subscribe({
      next: (response) => this.services.set(response.data),
      error: () => this.error.set('No fue posible cargar el catalogo de servicios.'),
    });

    this.api.leads().subscribe({
      next: (response) => {
        this.leads.set(response.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('No fue posible cargar los leads. Verifica el backend y la sesion.');
        this.loading.set(false);
      },
    });
  }


  formatDate(value: string | null | undefined): string {
    if (!value) {
      return 'Sin fecha';
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return 'Sin fecha';
    }
    try {
      return date.toLocaleString('es-CO', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true,
      });
    } catch {
      return 'Sin fecha';
    }
  }

  getDatePart(value: string | null | undefined): string {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(0, 10);
  }

  getTimePart(value: string | null | undefined): string {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(11, 16);
  }

  updateNextActionAt(control: { setValue(value: string): void }, dateVal: string, timeVal: string): void {
    if (!dateVal || !timeVal) {
      control.setValue('');
    } else {
      control.setValue(`${dateVal}T${timeVal}`);
    }
  }

  private isOverdue(value: string | undefined): boolean {
    const date = this.parseDate(value);
    return Boolean(date && date.getTime() < Date.now());
  }

  private isToday(value: string | undefined): boolean {
    const date = this.parseDate(value);
    if (!date) {
      return false;
    }
    const now = new Date();
    return date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth() && date.getDate() === now.getDate();
  }

  private parseDate(value: string | undefined): Date | null {
    if (!value) {
      return null;
    }
    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? null : date;
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
