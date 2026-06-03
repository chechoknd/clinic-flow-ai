import { Component, computed, inject, signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import {
  Appointment,
  AppointmentStatus,
  ClinicServiceItem,
  Professional,
  TimeSlot,
} from '../../core/services/api.models';
import { I18nService } from '../../core/i18n/i18n.service';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-schedule-page',
  imports: [ReactiveFormsModule, DatePipe],
  templateUrl: './schedule.page.html',
  styleUrl: './schedule.page.css',
})
export class SchedulePage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  readonly i18n = inject(I18nService);

  // Core State
  readonly professionals = signal<Professional[]>([]);
  readonly services = signal<ClinicServiceItem[]>([]);
  readonly appointments = signal<Appointment[]>([]);
  readonly selectedDate = signal<string>(new Date().toISOString().split('T')[0]);
  readonly selectedView = signal<'day' | 'week'>('day');

  // Selected Professional for Weekly View (defaults to first active)
  readonly activeWeeklyProfessional = signal<Professional | null>(null);

  // Loading / Messages
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  // Dialog & Detail Modals
  readonly isCreateModalOpen = signal(false);
  readonly isDetailModalOpen = signal(false);
  readonly selectedAppointment = signal<Appointment | null>(null);

  // Availability check
  readonly availableSlots = signal<TimeSlot[]>([]);
  readonly checkingAvailability = signal(false);

  // AI Assistant suggested copy
  readonly aiSuggestedMessage = signal<string | null>(null);
  readonly loadingAiMessage = signal(false);

  // Forms
  readonly appointmentForm = this.fb.nonNullable.group({
    contact_name: ['', [Validators.required, Validators.minLength(3)]],
    contact_phone: ['', [Validators.required]],
    professional_id: ['', [Validators.required]],
    service_id: ['', [Validators.required]],
    starts_at_date: ['', [Validators.required]],
    starts_at_time: ['09:00', [Validators.required]],
    duration_mins: [60, [Validators.required, Validators.min(1)]],
    admin_notes: [''],
    lead_id: [''],
  });

  readonly statusForm = this.fb.nonNullable.group({
    status: ['scheduled' as AppointmentStatus, [Validators.required]],
    admin_notes: [''],
  });

  readonly rescheduleForm = this.fb.nonNullable.group({
    starts_at_date: ['', [Validators.required]],
    starts_at_time: ['09:00', [Validators.required]],
    admin_note: [''],
  });

  readonly statusLabels: Record<AppointmentStatus, string> = {
    scheduled: 'appointment.status.scheduled',
    confirmed: 'appointment.status.confirmed',
    pending_confirmation: 'appointment.status.pending_confirmation',
    rescheduled: 'appointment.status.rescheduled',
    no_show: 'appointment.status.no_show',
    cancelled: 'appointment.status.cancelled',
    completed: 'appointment.status.completed',
    converted_from_lead: 'appointment.status.converted_from_lead',
  };

  readonly selectedDateLabel = computed(() => {
    const date = new Date(this.selectedDate() + 'T00:00:00');
    return date.toLocaleDateString(this.locale(), { weekday: 'long', day: 'numeric', month: 'long' });
  });

  readonly selectedDateAppointments = computed(() =>
    this.appointments().filter((appt) => appt.starts_at.split('T')[0] === this.selectedDate()),
  );

  readonly dailyStats = computed(() => {
    const appointments = this.selectedDateAppointments();
    return {
      total: appointments.length,
      pending: appointments.filter((appt) => appt.status === 'pending_confirmation' || appt.confirmation_status === 'pending').length,
      confirmed: appointments.filter((appt) => appt.status === 'confirmed' || appt.confirmation_status === 'confirmed').length,
      availableSlots: this.availableSlots().length,
    };
  });

  // Calendar Timeline Rows (08:00 to 20:00)
  readonly hoursList = Array.from({ length: 13 }, (_, i) => {
    const hr = 8 + i;
    return `${hr.toString().padStart(2, '0')}:00`;
  });

  // Color Mapping helper for CSS border/background cards
  readonly statusColors: Record<AppointmentStatus, { bg: string; border: string; text: string }> = {
    scheduled: { bg: 'bg-blue-50/90', border: 'border-blue-200', text: 'text-blue-700' },
    confirmed: { bg: 'bg-emerald-50/90', border: 'border-emerald-200', text: 'text-emerald-700' },
    pending_confirmation: { bg: 'bg-amber-50/90', border: 'border-amber-200', text: 'text-amber-700' },
    rescheduled: { bg: 'bg-purple-50/90', border: 'border-purple-200', text: 'text-purple-700' },
    no_show: { bg: 'bg-rose-50/90', border: 'border-rose-200', text: 'text-rose-700' },
    cancelled: { bg: 'bg-slate-50/90', border: 'border-slate-200', text: 'text-slate-700' },
    completed: { bg: 'bg-teal-50/90', border: 'border-teal-200', text: 'text-teal-700' },
    converted_from_lead: { bg: 'bg-indigo-50/90', border: 'border-indigo-200', text: 'text-indigo-700' },
  };

  // Weekdates List helper for Weekly View
  readonly weekDates = computed(() => {
    const baseDate = new Date(this.selectedDate() + 'T00:00:00');
    const day = baseDate.getDay();
    const sundayOffset = day === 0 ? -6 : 1 - day; // Align Lunes as start of week

    const dates: { dateStr: string; label: string; dayLabel: string }[] = [];
    for (let i = 0; i < 7; i++) {
      const d = new Date(baseDate);
      d.setDate(baseDate.getDate() + sundayOffset + i);
      const dateStr = d.toISOString().split('T')[0];
      dates.push({
        dateStr,
        label: `${d.getDate()} ${d.toLocaleString(this.locale(), { month: 'short' })}`,
        dayLabel: this.i18n.t(`weekday.short.${d.getDay()}`),
      });
    }
    return dates;
  });

  constructor() {
    this.loadCatalog();
    this.loadAppointments();

    this.route.queryParams.subscribe((params) => {
      if (params['contact_name'] || params['contact_phone'] || params['service_id'] || params['lead_id']) {
        const startsAtTime = params['starts_at_time'] || '09:00';
        if (params['starts_at_date']) {
          this.selectedDate.set(params['starts_at_date']);
        }
        this.openCreateDialog(undefined, startsAtTime);
        this.appointmentForm.patchValue({
          contact_name: params['contact_name'] || '',
          contact_phone: params['contact_phone'] || '',
          service_id: params['service_id'] || '',
          lead_id: params['lead_id'] || '',
          starts_at_date: params['starts_at_date'] || this.selectedDate(),
          admin_notes: params['admin_notes'] || '',
        });
        this.applyServiceDuration(params['service_id'] || '');
      }
    });
  }

  loadCatalog(): void {
    this.api.professionals().subscribe({
      next: (profs) => {
        const activeOnly = profs.filter((p) => p.is_active);
        this.professionals.set(activeOnly);
        if (activeOnly.length > 0 && !this.activeWeeklyProfessional()) {
          this.activeWeeklyProfessional.set(activeOnly[0]);
        }
      },
    });

    this.api.services().subscribe({
      next: (res) => {
        this.services.set(res.data.filter((s) => s.is_active));
      },
    });
  }

  loadAppointments(): void {
    this.loading.set(true);
    this.error.set(null);

    // If weekly view, load dates from Monday to Sunday of current week
    let dateFrom = this.selectedDate();
    let dateTo = this.selectedDate();

    if (this.selectedView() === 'week') {
      const dates = this.weekDates();
      dateFrom = dates[0].dateStr;
      dateTo = dates[6].dateStr;
    }

    this.api.appointments({ date_from: dateFrom, date_to: dateTo }).subscribe({
      next: (res) => {
        this.appointments.set(res.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('schedule.loadError'));
        this.loading.set(false);
      },
    });
  }

  changeDate(days: number): void {
    const d = new Date(this.selectedDate() + 'T00:00:00');
    const step = this.selectedView() === 'week' ? days * 7 : days;
    d.setDate(d.getDate() + step);
    this.selectedDate.set(d.toISOString().split('T')[0]);
    this.loadAppointments();
  }

  goToToday(): void {
    this.selectedDate.set(new Date().toISOString().split('T')[0]);
    this.loadAppointments();
  }

  onDateChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.value) {
      this.selectedDate.set(input.value);
      this.loadAppointments();
    }
  }

  toggleView(view: 'day' | 'week'): void {
    this.selectedView.set(view);
    this.loadAppointments();
  }

  selectWeeklyProfessional(p: Professional): void {
    this.activeWeeklyProfessional.set(p);
  }

  gridTemplateColumns(): string {
    const columns = this.selectedView() === 'day' ? Math.max(this.professionals().length, 1) : 7;
    return `repeat(${columns}, minmax(0, 1fr))`;
  }

  statusLabel(status: AppointmentStatus): string {
    return this.i18n.t(this.statusLabels[status] ?? status);
  }

  sourceLabel(source: string): string {
    const labels: Record<string, string> = {
      manual: this.i18n.t('source.manual'),
      whatsapp: this.i18n.t('source.whatsapp'),
      instagram: this.i18n.t('source.instagram'),
      facebook: this.i18n.t('source.facebook'),
      web: this.i18n.t('source.web'),
      llamada: this.i18n.t('source.llamada'),
      otro: this.i18n.t('source.otro'),
      lead_conversion: this.i18n.t('source.lead_conversion'),
    };
    return labels[source] ?? source;
  }

  onServiceChange(): void {
    this.applyServiceDuration(this.appointmentForm.controls.service_id.value);
    this.checkAvailability();
  }

  private applyServiceDuration(serviceID: string): void {
    const service = this.services().find((item) => item.id === serviceID);
    if (service?.duration_minutes) {
      this.appointmentForm.patchValue({ duration_mins: service.duration_minutes });
    }
  }

  // Visual absolute positioning math
  getSlotStyle(startsAt: string, endsAt: string) {
    const start = new Date(startsAt);
    const end = new Date(endsAt);

    const dayStartHour = 8;
    const startHour = start.getHours() + start.getMinutes() / 60;
    const endHour = end.getHours() + end.getMinutes() / 60;

    const duration = Math.max(0.5, endHour - startHour); // Minimum 30 mins slot visual height
    const topOffset = Math.max(0, (startHour - dayStartHour) * 4.5); // 4.5rem per hour row
    const height = duration * 4.5;

    return {
      top: `${topOffset}rem`,
      height: `${height}rem`,
    };
  }

  // Dialog triggers
  openCreateDialog(professionalId?: string, timeStr?: string): void {
    this.clearMessages();
    this.aiSuggestedMessage.set(null);
    this.availableSlots.set([]);

    const defaultProfessionalID = professionalId || this.activeWeeklyProfessional()?.id || this.professionals()[0]?.id || '';

    this.appointmentForm.reset({
      contact_name: '',
      contact_phone: '',
      professional_id: defaultProfessionalID,
      service_id: '',
      starts_at_date: this.selectedDate(),
      starts_at_time: timeStr || '09:00',
      duration_mins: 60,
      admin_notes: '',
      lead_id: '',
    });
    this.isCreateModalOpen.set(true);
  }

  closeCreateDialog(): void {
    this.isCreateModalOpen.set(false);
  }

  openDetailDialog(appt: Appointment): void {
    this.clearMessages();
    this.selectedAppointment.set(appt);
    this.aiSuggestedMessage.set(null);

    this.statusForm.reset({
      status: appt.status,
      admin_notes: appt.admin_notes ?? '',
    });

    const dateStr = appt.starts_at.split('T')[0];
    const timeStr = appt.starts_at.split('T')[1].substring(0, 5);

    this.rescheduleForm.reset({
      starts_at_date: dateStr,
      starts_at_time: timeStr,
      admin_note: '',
    });

    this.isDetailModalOpen.set(true);
  }

  closeDetailDialog(): void {
    this.isDetailModalOpen.set(false);
    this.selectedAppointment.set(null);
  }

  // Real-time availability checker
  checkAvailability(): void {
    const profId = this.appointmentForm.controls.professional_id.value;
    const serviceId = this.appointmentForm.controls.service_id.value;
    const date = this.appointmentForm.controls.starts_at_date.value;
    const duration = this.appointmentForm.controls.duration_mins.value;

    if (!profId || !date) {
      return;
    }

    this.checkingAvailability.set(true);
    this.api.availability(profId, date, date, serviceId || undefined, duration).subscribe({
      next: (res) => {
        this.availableSlots.set(res.slots);
        this.checkingAvailability.set(false);
      },
      error: () => {
        this.checkingAvailability.set(false);
      },
    });
  }

  selectAvailableSlot(slot: TimeSlot): void {
    const timeStr = slot.starts_at.split('T')[1].substring(0, 5);
    this.appointmentForm.patchValue({ starts_at_time: timeStr });
  }

  // Create Appointment submit handler
  saveAppointment(): void {
    if (this.appointmentForm.invalid) {
      this.appointmentForm.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    this.clearMessages();

    const value = this.appointmentForm.getRawValue();
    const startsAt = new Date(`${value.starts_at_date}T${value.starts_at_time}:00`);

    const payload = {
      professional_id: value.professional_id,
      service_id: value.service_id,
      contact_name: value.contact_name.trim(),
      contact_phone: value.contact_phone.trim() || undefined,
      starts_at: startsAt.toISOString(),
      duration_minutes: value.duration_mins,
      lead_id: value.lead_id || undefined,
      admin_notes: value.admin_notes.trim() || undefined,
    };

    this.api.createAppointment(payload).subscribe({
      next: () => {
        this.success.set(this.i18n.t('schedule.created'));
        this.loadAppointments();
        setTimeout(() => this.closeCreateDialog(), 1000);
        this.saving.set(false);
      },
      error: (err) => {
        const errMsg = err.error?.error?.message || this.i18n.t('schedule.createError');
        this.error.set(errMsg);
        this.saving.set(false);
      },
    });
  }

  // Detail panel status transition handler
  updateStatus(): void {
    const appt = this.selectedAppointment();
    if (!appt) return;

    this.saving.set(true);
    this.clearMessages();

    const value = this.statusForm.getRawValue();
    this.api.updateAppointmentStatus(appt.id, value.status, value.admin_notes || undefined).subscribe({
      next: (updated) => {
        this.success.set(this.i18n.t('schedule.statusUpdated'));
        this.appointments.update((items) => items.map((item) => item.id === updated.id ? updated : item));
        this.selectedAppointment.set(updated);
        this.saving.set(false);
      },
      error: (err) => {
        this.error.set(err.error?.error?.message || this.i18n.t('schedule.statusError'));
        this.saving.set(false);
      },
    });
  }

  // Detail panel reschedule handler
  reschedule(): void {
    const appt = this.selectedAppointment();
    if (!appt) return;

    this.saving.set(true);
    this.clearMessages();

    const value = this.rescheduleForm.getRawValue();
    const startsAt = new Date(`${value.starts_at_date}T${value.starts_at_time}:00`);

    const payload = {
      starts_at: startsAt.toISOString(),
      duration_minutes: this.appointmentDurationMinutes(appt),
      admin_note: value.admin_note.trim() || undefined,
    };

    this.api.rescheduleAppointment(appt.id, payload).subscribe({
      next: (updated) => {
        this.success.set(this.i18n.t('schedule.rescheduled'));
        this.appointments.update((items) => items.map((item) => item.id === updated.id ? updated : item));
        this.selectedAppointment.set(updated);

        // Update form values
        const dateStr = updated.starts_at.split('T')[0];
        const timeStr = updated.starts_at.split('T')[1].substring(0, 5);
        this.rescheduleForm.patchValue({
          starts_at_date: dateStr,
          starts_at_time: timeStr,
          admin_note: '',
        });

        this.saving.set(false);
      },
      error: (err) => {
        this.error.set(err.error?.error?.message || this.i18n.t('schedule.rescheduleError'));
        this.saving.set(false);
      },
    });
  }

  // Quick confirm status action
  quickConfirm(appt: Appointment): void {
    this.api.updateAppointmentStatus(appt.id, 'confirmed', this.i18n.t('schedule.quickConfirmNote')).subscribe({
      next: (updated) => {
        this.appointments.update((items) => items.map((item) => item.id === updated.id ? updated : item));
      },
    });
  }

  // AI suggestions message copy helpers
  generateAiMessage(type: 'confirmation' | 'no_show'): void {
    const appt = this.selectedAppointment();
    if (!appt) return;

    this.loadingAiMessage.set(true);
    this.aiSuggestedMessage.set(null);

    const payload = {
      patient_message: type === 'confirmation'
        ? this.i18n.t('schedule.aiConfirmationSeed').replace('{{service}}', appt.service.name)
        : this.i18n.t('schedule.aiNoShowSeed').replace('{{service}}', appt.service.name),
      service_id: appt.service.id,
      lead_id: appt.lead?.id || undefined,
    };

    this.api.replySuggestion(payload).subscribe({
      next: (res) => {
        const dateObj = new Date(appt.starts_at);
        const dayStr = dateObj.toLocaleDateString(this.locale(), { weekday: 'long', day: 'numeric', month: 'long' });
        const hrStr = dateObj.toLocaleTimeString(this.locale(), { hour: '2-digit', minute: '2-digit' });

        let msg = res.suggested_reply || res.suggested_message || res.message || '';

        // Clean and replace generic placeholders with real scheduled dates in draft
        msg = msg
          .replace(/\[Fecha\]/gi, dayStr)
          .replace(/\[Hora\]/gi, hrStr)
          .replace(/\[Nombre del Paciente\]/gi, appt.contact_name)
          .replace(/\[Nombre del Profesional\]/gi, appt.professional.full_name)
          .replace(/\[Servicio\]/gi, appt.service.name);

        this.aiSuggestedMessage.set(msg);
        this.loadingAiMessage.set(false);
      },
      error: () => {
        this.aiSuggestedMessage.set(this.i18n.t('schedule.aiError'));
        this.loadingAiMessage.set(false);
      },
    });
  }

  private appointmentDurationMinutes(appt: Appointment): number {
    const starts = new Date(appt.starts_at).getTime();
    const ends = new Date(appt.ends_at).getTime();
    const minutes = Math.round((ends - starts) / 60000);
    return Number.isFinite(minutes) && minutes > 0 ? minutes : 60;
  }

  copyToClipboard(text: string): void {
    void navigator.clipboard.writeText(text);
    this.success.set(this.i18n.t('schedule.copied'));
    setTimeout(() => this.success.set(null), 3000);
  }

  private locale(): string {
    return this.i18n.language() === 'en' ? 'en-US' : 'es-CO';
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
  }
}
