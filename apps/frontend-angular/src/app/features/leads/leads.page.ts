import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';

import { ClinicServiceItem, Lead, LeadAIInsight, LeadDetail, LeadNote, LeadStatus } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

interface QuickLeadAction {
  id: string;
  label: string;
  status: LeadStatus;
  note: string;
  followUpDays?: number;
  clearNextActionAt?: boolean;
}

@Component({
  selector: 'app-leads-page',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './leads.page.html',
  styleUrl: './leads.page.css',
})
export class LeadsPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute, { optional: true });

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
  readonly selectedLeadDetail = signal<LeadDetail | null>(null);
  readonly loading = signal(false);
  readonly loadingDetail = signal(false);
  readonly saving = signal(false);
  readonly quickActionSaving = signal<string | null>(null);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly phoneCopied = signal(false);
  readonly quickActions: QuickLeadAction[] = [
    {
      id: 'retry_tomorrow',
      label: 'Volver a intentar mañana',
      status: 'Contactado',
      note: 'Contacto realizado. Volver a intentar mañana.',
      followUpDays: 1,
    },
    {
      id: 'interested_two_days',
      label: 'Interesado, seguimiento en 2 días',
      status: 'Interesado',
      note: 'Lead interesado. Programar seguimiento comercial en 2 días.',
      followUpDays: 2,
    },
    {
      id: 'scheduled',
      label: 'Marcar agendado',
      status: 'Agendado',
      note: 'Lead marcado como agendado. Confirmar asistencia antes de la cita.',
      followUpDays: 1,
    },
    {
      id: 'lost',
      label: 'Marcar perdido',
      status: 'Perdido',
      note: 'Lead marcado como perdido. No requiere seguimiento por ahora.',
      clearNextActionAt: true,
    },
    {
      id: 'converted',
      label: 'Marcar convertido',
      status: 'Convertido',
      note: 'Lead convertido. Cerrar seguimiento comercial.',
      clearNextActionAt: true,
    },
  ];

  readonly filteredLeads = computed(() =>
    this.leads().filter((lead) => lead.status === this.selectedStatus()),
  );
  readonly latestAIInsight = computed(() => this.selectedLeadDetail()?.ai_insights?.[0] ?? null);

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
      this.selectedLeadDetail.set(null);
    }
  }

  selectLead(lead: Lead): void {
    this.selectedLead.set(lead);
    this.selectedLeadDetail.set(null);
    this.phoneCopied.set(false);
    this.success.set(null);
    this.error.set(null);
    this.updateForm.reset({
      status: lead.status,
      next_action_at: this.toLocalDateTimeValue(lead.next_action_at),
      note: '',
    });
    this.loadLeadDetail(lead.id);
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
          const nextActionAt = this.toApiDateTime(value.next_action_at);
          const updated: Lead = {
            ...lead,
            status: value.status,
            next_action_at: nextActionAt,
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

  applyQuickAction(action: QuickLeadAction): void {
    const lead = this.selectedLead();
    if (!lead) {
      return;
    }

    const nextActionAt = action.followUpDays ? this.addDays(action.followUpDays).toISOString() : undefined;
    this.quickActionSaving.set(action.id);
    this.error.set(null);
    this.success.set(null);

    this.api
      .updateLead(lead.id, {
        status: action.status,
        note: action.note,
        next_action_at: nextActionAt,
        clear_next_action_at: action.clearNextActionAt,
      })
      .subscribe({
        next: () => {
          const updated: Lead = {
            ...lead,
            status: action.status,
            next_action_at: action.clearNextActionAt ? undefined : nextActionAt,
          };
          this.leads.update((items) => items.map((item) => (item.id === lead.id ? updated : item)));
          this.selectedStatus.set(action.status);
          this.selectLead(updated);
          this.success.set('Accion rapida aplicada correctamente.');
          this.quickActionSaving.set(null);
        },
        error: () => {
          this.error.set('No fue posible aplicar la accion rapida.');
          this.quickActionSaving.set(null);
        },
      });
  }

  whatsappLink(phone: string | undefined): string {
    const digits = (phone || '').replace(/\D/g, '');
    return digits ? `https://wa.me/${digits}` : '';
  }

  copySelectedLeadPhone(): void {
    const phone = this.selectedLeadDetail()?.phone || this.selectedLead()?.phone;
    if (!phone) {
      return;
    }

    void navigator.clipboard.writeText(phone).then(() => this.phoneCopied.set(true));
  }

  insightIntentLabel(insight: LeadAIInsight): string {
    const labels: Record<string, string> = {
      high: 'Alta intencion',
      medium: 'Intencion media',
      low: 'Baja intencion',
    };
    return labels[insight.intent] ?? insight.intent;
  }

  selectedLeadServiceName(): string {
    const detail = this.selectedLeadDetail();
    if (detail?.service?.name) {
      return detail.service.name;
    }
    return this.selectedLead()?.service_name || 'Sin servicio definido';
  }

  recentLeadNotes(detail: LeadDetail): LeadNote[] {
    return [...detail.notes]
      .sort((first, second) => this.noteTimestamp(second) - this.noteTimestamp(first))
      .slice(0, 3);
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
        this.selectRouteLead(response.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('No fue posible cargar los leads. Verifica el backend y la sesion.');
        this.loading.set(false);
      },
    });
  }

  private selectRouteLead(leads: Lead[]): void {
    const leadID = this.route?.snapshot.queryParamMap.get('lead_id');
    if (!leadID) {
      return;
    }
    const lead = leads.find((item) => item.id === leadID);
    if (!lead) {
      return;
    }
    this.selectedStatus.set(lead.status);
    this.selectLead(lead);
  }

  private loadLeadDetail(leadID: string): void {
    this.loadingDetail.set(true);
    this.api.lead(leadID).subscribe({
      next: (detail) => {
        this.selectedLeadDetail.set(detail);
        this.loadingDetail.set(false);
      },
      error: () => {
        this.error.set('No fue posible cargar el detalle del lead.');
        this.loadingDetail.set(false);
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

  private addDays(days: number): Date {
    const date = new Date();
    date.setDate(date.getDate() + days);
    date.setSeconds(0, 0);
    return date;
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

  private noteTimestamp(note: LeadNote): number {
    const date = new Date(note.created_at);
    return Number.isNaN(date.getTime()) ? 0 : date.getTime();
  }
}
