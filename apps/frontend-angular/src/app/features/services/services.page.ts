import { Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { ClinicServiceItem } from '../../core/services/api.models';
import { I18nService } from '../../core/i18n/i18n.service';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-services-page',
  imports: [ReactiveFormsModule],
  templateUrl: './services.page.html',
  styleUrl: './services.page.css',
})
export class ServicesPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  readonly i18n = inject(I18nService);

  readonly services = signal<ClinicServiceItem[]>([]);
  readonly selectedService = signal<ClinicServiceItem | null>(null);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly deleting = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly formTitle = computed(() =>
    this.selectedService() ? this.i18n.t('services.editTitle') : this.i18n.t('services.newTitle'),
  );

  readonly serviceForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(3)]],
    description: [''],
    duration_minutes: [60, [Validators.required, Validators.min(1)]],
    price_from: [0, [Validators.required, Validators.min(0)]],
    benefits: [''],
    common_objections: [''],
    is_active: [true],
  });

  constructor() {
    this.loadServices();
  }

  startCreate(): void {
    this.selectedService.set(null);
    this.clearMessages();
    this.serviceForm.reset({
      name: '',
      description: '',
      duration_minutes: 60,
      price_from: 0,
      benefits: '',
      common_objections: '',
      is_active: true,
    });
  }

  selectService(service: ClinicServiceItem): void {
    this.selectedService.set(service);
    this.clearMessages();
    this.serviceForm.reset({
      name: service.name,
      description: service.description ?? '',
      duration_minutes: service.duration_minutes ?? 60,
      price_from: service.price_from ?? 0,
      benefits: this.joinLines(service.benefits),
      common_objections: this.joinLines(service.common_objections),
      is_active: service.is_active ?? true,
    });
  }

  saveService(): void {
    if (this.serviceForm.invalid) {
      this.serviceForm.markAllAsTouched();
      return;
    }

    const selected = this.selectedService();
    const value = this.serviceForm.getRawValue();
    const payload = {
      name: value.name.trim(),
      description: value.description.trim() || undefined,
      duration_minutes: Number(value.duration_minutes),
      price_from: Number(value.price_from),
      benefits: this.splitLines(value.benefits),
      faq: [],
      common_objections: this.splitLines(value.common_objections),
    };

    this.saving.set(true);
    this.clearMessages();

    if (selected) {
      this.api.updateService(selected.id, { ...payload, is_active: value.is_active }).subscribe({
        next: (service) => {
          this.services.update((items) => items.map((item) => (item.id === service.id ? service : item)));
          this.selectService(service);
          this.success.set(this.i18n.t('services.updated'));
          this.saving.set(false);
        },
        error: () => {
          this.error.set(this.i18n.t('services.updateError'));
          this.saving.set(false);
        },
      });
      return;
    }

    this.api.createService(payload).subscribe({
      next: (service) => {
        this.services.update((items) => [service, ...items]);
        this.selectService(service);
        this.success.set(this.i18n.t('services.created'));
        this.saving.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('services.createError'));
        this.saving.set(false);
      },
    });
  }

  deleteSelected(): void {
    const selected = this.selectedService();
    if (!selected) {
      return;
    }

    this.deleting.set(true);
    this.clearMessages();
    this.api.deleteService(selected.id).subscribe({
      next: () => {
        this.services.update((items) => items.filter((item) => item.id !== selected.id));
        this.startCreate();
        this.success.set(this.i18n.t('services.deleted'));
        this.deleting.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('services.deleteError'));
        this.deleting.set(false);
      },
    });
  }

  private loadServices(): void {
    this.loading.set(true);
    this.api.services().subscribe({
      next: (response) => {
        this.services.set(response.data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set(this.i18n.t('services.loadError'));
        this.loading.set(false);
      },
    });
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
  }

  private splitLines(value: string): string[] {
    return value
      .split('\n')
      .map((item) => item.trim())
      .filter(Boolean);
  }

  private joinLines(values: string[] | undefined): string {
    return values?.join('\n') ?? '';
  }
}
