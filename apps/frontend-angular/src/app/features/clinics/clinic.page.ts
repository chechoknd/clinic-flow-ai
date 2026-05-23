import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { ClinicProfile, CommunicationTone } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-clinic-page',
  imports: [ReactiveFormsModule],
  templateUrl: './clinic.page.html',
  styleUrl: './clinic.page.css',
})
export class ClinicPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);

  readonly tones: Array<{ value: CommunicationTone; label: string }> = [
    { value: 'amable', label: 'Amable' },
    { value: 'profesional', label: 'Profesional' },
    { value: 'cercano', label: 'Cercano' },
    { value: 'juvenil', label: 'Juvenil' },
    { value: 'elegante', label: 'Elegante' },
  ];
  readonly clinic = signal<ClinicProfile | null>(null);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly clinicForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(3)]],
    city: ['', Validators.required],
    phone: ['', Validators.pattern(/^$|^\+[1-9]\d{7,14}$/)],
    whatsapp: ['', [Validators.required, Validators.pattern(/^\+[1-9]\d{7,14}$/)]],
    address: [''],
    communication_tone: ['amable' as CommunicationTone, Validators.required],
  });

  constructor() {
    this.loadClinic();
  }

  saveClinic(): void {
    if (this.clinicForm.invalid) {
      this.clinicForm.markAllAsTouched();
      return;
    }

    const current = this.clinic();
    const value = this.clinicForm.getRawValue();
    this.saving.set(true);
    this.clearMessages();

    this.api
      .updateClinicCurrent({
        name: value.name.trim(),
        city: value.city.trim(),
        phone: value.phone.trim() || undefined,
        whatsapp: value.whatsapp.trim(),
        address: value.address.trim() || undefined,
        opening_hours: current?.opening_hours ?? {},
        general_faq: current?.general_faq ?? [],
        communication_tone: value.communication_tone,
      })
      .subscribe({
        next: (profile) => {
          this.clinic.set(profile);
          this.hydrateForm(profile);
          this.success.set('Perfil de clinica actualizado correctamente.');
          this.saving.set(false);
        },
        error: () => {
          this.error.set('No fue posible actualizar el perfil. Revisa los datos e intenta de nuevo.');
          this.saving.set(false);
        },
      });
  }

  private loadClinic(): void {
    this.loading.set(true);
    this.api.clinicCurrent().subscribe({
      next: (profile) => {
        this.clinic.set(profile);
        this.hydrateForm(profile);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('No se pudo cargar la clinica. Verifica el backend y la sesion.');
        this.loading.set(false);
      },
    });
  }

  private hydrateForm(profile: ClinicProfile): void {
    this.clinicForm.reset({
      name: profile.name,
      city: profile.city,
      phone: profile.phone ?? '',
      whatsapp: profile.whatsapp,
      address: profile.address ?? '',
      communication_tone: profile.communication_tone,
    });
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
  }
}
