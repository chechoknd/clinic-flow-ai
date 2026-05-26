import { Component, computed, inject, signal } from '@angular/core';
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';

import { ClinicProfile, ClinicServiceItem, CurrencyMetadata } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';
import { defaultCurrency, formatClinicCurrency } from '../../shared/currency-format';

@Component({
  selector: 'app-services-page',
  imports: [ReactiveFormsModule],
  templateUrl: './services.page.html',
  styleUrl: './services.page.css',
})
export class ServicesPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);

  readonly services = signal<ClinicServiceItem[]>([]);
  readonly clinic = signal<ClinicProfile | null>(null);
  readonly selectedService = signal<ClinicServiceItem | null>(null);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly deleting = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly formTitle = computed(() => (this.selectedService() ? 'Editar servicio' : 'Nuevo servicio'));
  readonly priceInputValue = signal(0);
  readonly currentCurrency = computed(() => this.clinic()?.currency ?? defaultCurrency);
  readonly priceStep = computed(() => (this.currentCurrency().decimal_digits === 0 ? '1' : '0.01'));
  readonly priceDecimalHelp = computed(() => {
    const currency = this.currentCurrency();
    return currency.decimal_digits === 0
      ? `Valor en ${currency.code}, sin decimales.`
      : `Valor en ${currency.code}, hasta ${currency.decimal_digits} decimales.`;
  });
  readonly formattedPricePreview = computed(() => formatClinicCurrency(this.priceInputValue(), this.currentCurrency()));

  readonly serviceForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(3)]],
    description: [''],
    duration_minutes: [60, [Validators.required, Validators.min(1)]],
    price_from: [0, [Validators.required, Validators.min(0), this.pricePrecisionValidator()]],
    benefits: [''],
    common_objections: [''],
    is_active: [true],
  });

  constructor() {
    this.priceInputValue.set(this.serviceForm.controls.price_from.value);
    this.serviceForm.controls.price_from.valueChanges.subscribe((value) => {
      this.priceInputValue.set(Number(value ?? 0));
    });
    this.loadClinic();
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
          this.success.set('Servicio actualizado correctamente.');
          this.saving.set(false);
        },
        error: (err: unknown) => {
          this.error.set(this.serviceErrorMessage(err, 'No fue posible actualizar el servicio.'));
          this.saving.set(false);
        },
      });
      return;
    }

    this.api.createService(payload).subscribe({
      next: (service) => {
        this.services.update((items) => [service, ...items]);
        this.selectService(service);
        this.success.set('Servicio creado correctamente.');
        this.saving.set(false);
      },
      error: (err: unknown) => {
        this.error.set(this.serviceErrorMessage(err, 'No fue posible crear el servicio. Revisa los datos e intenta de nuevo.'));
        this.saving.set(false);
      },
    });
  }

  formatPrice(value: number | null | undefined, currencyCode?: string): string {
    const clinicCurrency = this.currentCurrency();
    const currency: CurrencyMetadata = currencyCode && currencyCode !== clinicCurrency.code
      ? { ...clinicCurrency, code: currencyCode as CurrencyMetadata['code'] }
      : clinicCurrency;

    return formatClinicCurrency(value, currency);
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
        this.success.set('Servicio eliminado correctamente.');
        this.deleting.set(false);
      },
      error: () => {
        this.error.set('No fue posible eliminar el servicio.');
        this.deleting.set(false);
      },
    });
  }

  private loadClinic(): void {
    this.api.clinicCurrent().subscribe({
      next: (profile) => {
        this.clinic.set(profile);
        this.serviceForm.controls.price_from.updateValueAndValidity();
      },
      error: () => this.error.set('No fue posible cargar la moneda de la clinica.'),
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
        this.error.set('No fue posible cargar el catalogo de servicios.');
        this.loading.set(false);
      },
    });
  }

  private clearMessages(): void {
    this.error.set(null);
    this.success.set(null);
  }

  private serviceErrorMessage(err: unknown, fallback: string): string {
    const apiMessage = this.extractApiErrorMessage(err);
    if (apiMessage?.includes('too many decimal places')) {
      const currency = this.currentCurrency();
      return currency.decimal_digits === 0
        ? `La moneda ${currency.code} no permite decimales en el precio.`
        : `La moneda ${currency.code} permite maximo ${currency.decimal_digits} decimales.`;
    }
    if (apiMessage?.includes('greater than or equal to zero')) {
      return 'El precio no puede ser negativo.';
    }
    return fallback;
  }

  private extractApiErrorMessage(err: unknown): string | null {
    if (!err || typeof err !== 'object') {
      return null;
    }
    const maybeHttpError = err as { error?: { error?: { message?: unknown } } };
    const message = maybeHttpError.error?.error?.message;
    return typeof message === 'string' ? message : null;
  }

  private pricePrecisionValidator(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const value = control.value;
      if (value === null || value === undefined || value === '') {
        return null;
      }

      const amount = Number(value);
      if (!Number.isFinite(amount)) {
        return { invalidPrice: true };
      }

      const decimalDigits = this.currentCurrency().decimal_digits;
      const factor = 10 ** decimalDigits;
      return Math.round(amount * factor) === amount * factor
        ? null
        : { currencyDecimals: { allowed: decimalDigits } };
    };
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
