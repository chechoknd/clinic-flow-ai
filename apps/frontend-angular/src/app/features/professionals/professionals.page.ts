import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { ClinicServiceItem, Professional, WorkingHoursConfig } from '../../core/services/api.models';
import { I18nService } from '../../core/i18n/i18n.service';
import { ApiService } from '../../core/services/api.service';

interface WeekdayState {
  label: string;
  active: boolean;
  start: string;
  end: string;
}

@Component({
  selector: 'app-professionals-page',
  imports: [ReactiveFormsModule],
  templateUrl: './professionals.page.html',
  styleUrl: './professionals.page.css',
})
export class ProfessionalsPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  readonly i18n = inject(I18nService);

  readonly professionals = signal<Professional[]>([]);
  readonly services = signal<ClinicServiceItem[]>([]);
  readonly selectedProfessional = signal<Professional | null>(null);
  
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  
  readonly formTitle = computed(() =>
    this.selectedProfessional() ? this.i18n.t('professionals.editTitle') : this.i18n.t('professionals.newTitle'),
  );

  readonly colorPresets = [
    '#2563EB', // Royal Blue
    '#0D9488', // Teal
    '#7C3AED', // Violet
    '#DB2777', // Pink
    '#EA580C', // Orange
    '#E11D48'  // Rose
  ];

  readonly daysList: { key: keyof WorkingHoursConfig; label: string }[] = [
    { key: 'monday', label: 'weekday.monday' },
    { key: 'tuesday', label: 'weekday.tuesday' },
    { key: 'wednesday', label: 'weekday.wednesday' },
    { key: 'thursday', label: 'weekday.thursday' },
    { key: 'friday', label: 'weekday.friday' },
    { key: 'saturday', label: 'weekday.saturday' },
    { key: 'sunday', label: 'weekday.sunday' }
  ];

  // Flat state to bind weekly working hours inputs easily
  readonly dayStates = signal<Record<string, WeekdayState>>({
    monday: { label: 'Lunes', active: true, start: '09:00', end: '18:00' },
    tuesday: { label: 'Martes', active: true, start: '09:00', end: '18:00' },
    wednesday: { label: 'Miércoles', active: true, start: '09:00', end: '18:00' },
    thursday: { label: 'Jueves', active: true, start: '09:00', end: '18:00' },
    friday: { label: 'Viernes', active: true, start: '09:00', end: '18:00' },
    saturday: { label: 'Sábado', active: false, start: '09:00', end: '13:00' },
    sunday: { label: 'Domingo', active: false, start: '09:00', end: '13:00' },
  });

  // Flat state to bind service assignment checklist
  readonly selectedServiceIds = signal<Record<string, boolean>>({});

  readonly professionalForm = this.fb.nonNullable.group({
    full_name: ['', [Validators.required, Validators.minLength(3)]],
    role_or_specialty: [''],
    calendar_color: ['#2563EB', [Validators.required]],
    is_active: [true],
  });

  constructor() {
    this.loadData();
  }

  loadData(): void {
    this.loading.set(true);
    this.api.professionals().subscribe({
      next: (profs) => {
        this.professionals.set(profs);
        this.loading.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('professionals.loadError'));
        this.loading.set(false);
      }
    });

    this.api.services().subscribe({
      next: (res) => {
        this.services.set(res.data);
      }
    });
  }

  startCreate(): void {
    this.selectedProfessional.set(null);
    this.clearMessages();
    this.professionalForm.reset({
      full_name: '',
      role_or_specialty: '',
      calendar_color: '#2563EB',
      is_active: true
    });
    
    // Reset day states to default weekdays
    this.dayStates.set({
      monday: { label: 'Lunes', active: true, start: '09:00', end: '18:00' },
      tuesday: { label: 'Martes', active: true, start: '09:00', end: '18:00' },
      wednesday: { label: 'Miércoles', active: true, start: '09:00', end: '18:00' },
      thursday: { label: 'Jueves', active: true, start: '09:00', end: '18:00' },
      friday: { label: 'Viernes', active: true, start: '09:00', end: '18:00' },
      saturday: { label: 'Sábado', active: false, start: '09:00', end: '13:00' },
      sunday: { label: 'Domingo', active: false, start: '09:00', end: '13:00' },
    });

    // Reset services checklists
    this.selectedServiceIds.set({});
  }

  selectProfessional(p: Professional): void {
    this.selectedProfessional.set(p);
    this.clearMessages();
    this.professionalForm.reset({
      full_name: p.full_name,
      role_or_specialty: p.role_or_specialty ?? '',
      calendar_color: p.calendar_color ?? '#2563EB',
      is_active: p.is_active
    });

    // Map JSON working hours to flat states
    const states = {
      monday: { label: 'Lunes', active: false, start: '09:00', end: '18:00' },
      tuesday: { label: 'Martes', active: false, start: '09:00', end: '18:00' },
      wednesday: { label: 'Miércoles', active: false, start: '09:00', end: '18:00' },
      thursday: { label: 'Jueves', active: false, start: '09:00', end: '18:00' },
      friday: { label: 'Viernes', active: false, start: '09:00', end: '18:00' },
      saturday: { label: 'Sábado', active: false, start: '09:00', end: '13:00' },
      sunday: { label: 'Domingo', active: false, start: '09:00', end: '13:00' },
    };

    const wh = p.working_hours;
    this.daysList.forEach((d) => {
      const interval = wh[d.key];
      if (interval && interval.length > 0) {
        states[d.key].active = true;
        states[d.key].start = interval[0].start;
        states[d.key].end = interval[0].end;
      }
    });
    this.dayStates.set(states);

    // Map service checklist
    const sIds: Record<string, boolean> = {};
    if (p.service_ids) {
      p.service_ids.forEach((id: string) => {
        sIds[id] = true;
      });
    }
    this.selectedServiceIds.set(sIds);
  }

  toggleDay(day: string): void {
    const states = { ...this.dayStates() };
    states[day].active = !states[day].active;
    this.dayStates.set(states);
  }

  onTimeChange(day: string, type: 'start' | 'end', event: Event): void {
    const input = event.target as HTMLInputElement;
    const states = { ...this.dayStates() };
    if (type === 'start') {
      states[day].start = input.value;
    } else {
      states[day].end = input.value;
    }
    this.dayStates.set(states);
  }

  toggleService(serviceId: string): void {
    const sIds = { ...this.selectedServiceIds() };
    sIds[serviceId] = !sIds[serviceId];
    this.selectedServiceIds.set(sIds);
  }

  selectColor(color: string): void {
    this.professionalForm.patchValue({ calendar_color: color });
  }

  saveProfessional(): void {
    if (this.professionalForm.invalid) {
      this.professionalForm.markAllAsTouched();
      return;
    }

    const selected = this.selectedProfessional();
    const value = this.professionalForm.getRawValue();
    
    // Construct weekly working hours config
    const workingHours: WorkingHoursConfig = {};
    const states = this.dayStates();
    this.daysList.forEach((d) => {
      const state = states[d.key];
      if (state.active) {
        workingHours[d.key] = [{ start: state.start, end: state.end }];
      } else {
        workingHours[d.key] = [];
      }
    });

    // Extract assigned service IDs
    const serviceIds = Object.entries(this.selectedServiceIds())
      .filter(([, checked]) => checked)
      .map(([id]) => id);

    const payload = {
      full_name: value.full_name.trim(),
      role_or_specialty: value.role_or_specialty.trim() || undefined,
      calendar_color: value.calendar_color,
      working_hours: workingHours,
      service_ids: serviceIds
    };

    this.saving.set(true);
    this.clearMessages();

    if (selected) {
      this.api.updateProfessional(selected.id, { ...payload, is_active: value.is_active }).subscribe({
        next: (prof) => {
          this.professionals.update((items) => items.map((item) => (item.id === prof.id ? prof : item)));
          this.selectProfessional(prof);
          this.success.set(this.i18n.t('professionals.updated'));
          this.saving.set(false);
        },
        error: (err) => {
          const errMsg = err.error?.error?.message || this.i18n.t('professionals.updateError');
          this.error.set(errMsg);
          this.saving.set(false);
        }
      });
    } else {
      this.api.createProfessional(payload).subscribe({
        next: (prof) => {
          this.professionals.update((items) => [prof, ...items]);
          this.selectProfessional(prof);
          this.success.set(this.i18n.t('professionals.created'));
          this.saving.set(false);
        },
        error: (err) => {
          const errMsg = err.error?.error?.message || this.i18n.t('professionals.createError');
          this.error.set(errMsg);
          this.saving.set(false);
        }
      });
    }
  }

  getServiceName(serviceId: string): string {
    const service = this.services().find(s => s.id === serviceId);
    return service ? service.name : serviceId;
  }

  getServiceIDs(p: Professional): string[] {
    return p.service_ids || [];
  }

  dayLabel(day: { label: string }): string {
    return this.i18n.t(day.label);
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
  }
}
