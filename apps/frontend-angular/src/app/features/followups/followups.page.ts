import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';

import { FollowUp, LeadStatus } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

interface FollowUpGroup {
  label: string;
  description: string;
  tone: 'urgent' | 'today' | 'upcoming';
  items: FollowUp[];
}

@Component({
  selector: 'app-followups-page',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './followups.page.html',
  styleUrl: './followups.page.css',
})
export class FollowupsPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);

  readonly statuses: LeadStatus[] = ['Contactado', 'Interesado', 'Agendado', 'No Respondio', 'Perdido', 'Convertido'];
  readonly followups = signal<FollowUp[]>([]);
  readonly selectedFollowUp = signal<FollowUp | null>(null);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly generating = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly generatedMessage = signal('');
  readonly copied = signal(false);
  readonly canGenerateMessage = computed(() => Boolean(this.selectedFollowUp()?.service_id));
  readonly followupGroups = computed<FollowUpGroup[]>(() => [
    {
      label: 'Vencidos',
      description: 'Requieren accion primero.',
      tone: 'urgent',
      items: this.followups().filter((followup) => this.isOverdue(followup.next_action_at)),
    },
    {
      label: 'Hoy',
      description: 'Contactos programados para el dia.',
      tone: 'today',
      items: this.followups().filter((followup) => this.isToday(followup.next_action_at) && !this.isOverdue(followup.next_action_at)),
    },
    {
      label: 'Proximos',
      description: 'Seguimientos futuros.',
      tone: 'upcoming',
      items: this.followups().filter((followup) => !this.isToday(followup.next_action_at) && !this.isOverdue(followup.next_action_at)),
    },
  ]);

  readonly actionForm = this.fb.nonNullable.group({
    status: ['Contactado' as LeadStatus, Validators.required],
    next_action_at: ['', Validators.required],
    note: [''],
  });

  constructor() {
    this.loadFollowUps();
  }

  selectFollowUp(followup: FollowUp): void {
    this.selectedFollowUp.set(followup);
    this.generatedMessage.set('');
    this.copied.set(false);
    this.success.set(null);
    this.error.set(null);
    this.actionForm.reset({
      status: followup.status,
      next_action_at: this.toLocalDateTimeValue(followup.next_action_at),
      note: '',
    });
  }

  completeSelected(): void {
    const followup = this.selectedFollowUp();
    if (!followup) {
      return;
    }

    this.saving.set(true);
    this.clearMessages();
    const value = this.actionForm.getRawValue();

    this.api
      .completeFollowUp(followup.id, {
        status: value.status,
        note: value.note.trim() || undefined,
      })
      .subscribe({
        next: () => {
          this.followups.update((items) => items.filter((item) => item.id !== followup.id));
          this.selectedFollowUp.set(null);
          this.success.set('Seguimiento completado correctamente.');
          this.saving.set(false);
        },
        error: () => {
          this.error.set('No fue posible completar el seguimiento.');
          this.saving.set(false);
        },
      });
  }

  rescheduleSelected(): void {
    const followup = this.selectedFollowUp();
    if (!followup) {
      return;
    }
    if (this.actionForm.controls.next_action_at.invalid) {
      this.actionForm.controls.next_action_at.markAsTouched();
      return;
    }

    this.saving.set(true);
    this.clearMessages();
    const value = this.actionForm.getRawValue();
    const nextActionAt = this.toApiDateTime(value.next_action_at);
    if (!nextActionAt) {
      this.error.set('La fecha de proxima accion no es valida.');
      this.saving.set(false);
      return;
    }

    this.api
      .rescheduleFollowUp(followup.id, {
        next_action_at: nextActionAt,
        note: value.note.trim() || undefined,
      })
      .subscribe({
        next: () => {
          const updated = { ...followup, next_action_at: nextActionAt };
          this.followups.update((items) => items.map((item) => (item.id === followup.id ? updated : item)));
          this.selectFollowUp(updated);
          this.success.set('Seguimiento reprogramado correctamente.');
          this.saving.set(false);
        },
        error: () => {
          this.error.set('No fue posible reprogramar el seguimiento.');
          this.saving.set(false);
        },
      });
  }

  generateMessage(): void {
    const followup = this.selectedFollowUp();
    if (!followup || !followup.service_id) {
      this.error.set('Selecciona un seguimiento con servicio para generar mensaje.');
      return;
    }

    this.generating.set(true);
    this.clearMessages();
    const note = this.actionForm.controls.note.value.trim();

    this.api
      .followUpMessage({
        lead_id: followup.id,
        service_id: followup.service_id,
        last_contact_note: note || undefined,
      })
      .subscribe({
        next: (response) => {
          this.generatedMessage.set(response.suggested_message);
          this.generating.set(false);
        },
        error: () => {
          this.error.set('No fue posible generar el mensaje de seguimiento.');
          this.generating.set(false);
        },
      });
  }

  copyGeneratedMessage(): void {
    if (!this.generatedMessage()) {
      return;
    }

    void navigator.clipboard.writeText(this.generatedMessage()).then(() => this.copied.set(true));
  }

  formatDate(value: string | null | undefined): string {
    if (!value) {
      return 'Sin fecha';
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return 'Sin fecha';
    }
    return date.toLocaleString('es-CO', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }

  getDatePart(value: string | null | undefined): string {
    if (!value) {
      return '';
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return '';
    }
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(0, 10);
  }

  getTimePart(value: string | null | undefined): string {
    if (!value) {
      return '';
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return '';
    }
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(11, 16);
  }

  updateNextActionAt(control: { setValue(value: string): void }, dateValue: string, timeValue: string): void {
    control.setValue(dateValue && timeValue ? `${dateValue}T${timeValue}` : '');
  }

  private loadFollowUps(): void {
    this.loading.set(true);
    this.api.followups().subscribe({
      next: (response) => {
        this.followups.set(response.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('No fue posible cargar los seguimientos. Verifica el backend y la sesion.');
        this.loading.set(false);
      },
    });
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
    this.copied.set(false);
  }

  private toApiDateTime(value: string | undefined): string | undefined {
    if (!value) {
      return undefined;
    }

    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
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
