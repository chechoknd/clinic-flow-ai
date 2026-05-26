import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { ClinicProfile, ClinicServiceItem } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';
import { ServicesPage } from './services.page';

const clinic: ClinicProfile = {
  id: 'clinic-1',
  name: 'Sonrisa Viva',
  clinic_type: 'dental',
  city: 'Bogota',
  country_code: 'CO',
  currency_code: 'COP',
  currency: {
    code: 'COP',
    symbol: '$',
    locale: 'es-CO',
    decimal_digits: 0,
    thousand_separator: '.',
    decimal_separator: ',',
    symbol_position: 'before',
  },
  whatsapp: '+573001234567',
  opening_hours: {},
  general_faq: [],
  communication_tone: 'profesional',
};

const services: ClinicServiceItem[] = [
  {
    id: 'service-1',
    name: 'Blanqueamiento dental',
    description: 'Tratamiento estetico',
    duration_minutes: 60,
    price_from: 250000,
    currency_code: 'COP',
    benefits: ['Sonrisa mas clara'],
    faq: [],
    common_objections: ['Precio'],
    is_active: true,
  },
];

describe('ServicesPage', () => {
  let fixture: ComponentFixture<ServicesPage>;
  let api: {
    clinicCurrent: ReturnType<typeof vi.fn>;
    services: ReturnType<typeof vi.fn>;
    createService: ReturnType<typeof vi.fn>;
    updateService: ReturnType<typeof vi.fn>;
    deleteService: ReturnType<typeof vi.fn>;
  };

  beforeEach(async () => {
    api = {
      clinicCurrent: vi.fn(() => of(clinic)),
      services: vi.fn(() => of({ data: services })),
      createService: vi.fn((payload) => of({ id: 'service-2', ...payload, is_active: true })),
      updateService: vi.fn((id, payload) => of({ id, ...payload })),
      deleteService: vi.fn(() => of(undefined)),
    };

    await TestBed.configureTestingModule({
      imports: [ServicesPage],
      providers: [
        {
          provide: ApiService,
          useValue: api,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ServicesPage);
    fixture.detectChanges();
  });

  it('loads services from the API service', () => {
    expect(fixture.componentInstance.services()).toEqual(services);
  });

  it('selects a service and hydrates the form', () => {
    fixture.componentInstance.selectService(services[0]);

    expect(fixture.componentInstance.selectedService()).toEqual(services[0]);
    expect(fixture.componentInstance.serviceForm.controls.name.value).toBe('Blanqueamiento dental');
    expect(fixture.componentInstance.serviceForm.controls.benefits.value).toBe('Sonrisa mas clara');
  });

  it('creates a service with normalized array fields', () => {
    fixture.componentInstance.startCreate();
    fixture.componentInstance.serviceForm.setValue({
      name: 'Ortodoncia',
      description: 'Alineacion dental',
      duration_minutes: 45,
      price_from: 120000,
      benefits: 'Mejora mordida\nMejora estetica',
      common_objections: 'Precio',
      is_active: true,
    });

    fixture.componentInstance.saveService();

    expect(api.createService).toHaveBeenCalledWith({
      name: 'Ortodoncia',
      description: 'Alineacion dental',
      duration_minutes: 45,
      price_from: 120000,
      benefits: ['Mejora mordida', 'Mejora estetica'],
      faq: [],
      common_objections: ['Precio'],
    });
    expect(fixture.componentInstance.services()[0].id).toBe('service-2');
    expect(fixture.componentInstance.success()).toBe('Servicio creado correctamente.');
  });

  it('rejects decimal service prices for zero-decimal currencies', () => {
    fixture.componentInstance.startCreate();
    fixture.componentInstance.serviceForm.setValue({
      name: 'Limpieza',
      description: 'Control preventivo',
      duration_minutes: 30,
      price_from: 120000.5,
      benefits: '',
      common_objections: '',
      is_active: true,
    });

    fixture.componentInstance.saveService();

    expect(fixture.componentInstance.serviceForm.controls.price_from.hasError('currencyDecimals')).toBe(true);
    expect(api.createService).not.toHaveBeenCalled();
  });

  it('allows two decimal service prices for two-decimal currencies', () => {
    fixture.componentInstance.clinic.set({
      ...clinic,
      country_code: 'PE',
      currency_code: 'PEN',
      currency: {
        code: 'PEN',
        symbol: 'S/',
        locale: 'es-PE',
        decimal_digits: 2,
        thousand_separator: ',',
        decimal_separator: '.',
        symbol_position: 'before',
      },
    });
    fixture.componentInstance.serviceForm.controls.price_from.updateValueAndValidity();
    fixture.componentInstance.startCreate();
    fixture.componentInstance.serviceForm.setValue({
      name: 'Consulta estetica',
      description: 'Valoracion comercial',
      duration_minutes: 30,
      price_from: 120.5,
      benefits: '',
      common_objections: '',
      is_active: true,
    });

    fixture.componentInstance.saveService();

    expect(api.createService).toHaveBeenCalledWith(expect.objectContaining({ price_from: 120.5 }));
  });

  it('shows the active currency and formatted price preview', () => {
    fixture.componentInstance.serviceForm.controls.price_from.setValue(250000);
    fixture.detectChanges();

    const text = fixture.nativeElement.textContent as string;
    expect(text).toContain('COP');
    expect(text).toContain('Vista previa');
    expect(text).toContain('COP');
  });

  it('shows a specific backend price precision error', () => {
    api.createService.mockReturnValueOnce(
      throwError(() => ({
        error: {
          error: {
            code: 'INVALID_REQUEST',
            message: 'price_from has too many decimal places for clinic currency',
          },
        },
      })),
    );
    fixture.componentInstance.startCreate();
    fixture.componentInstance.serviceForm.setValue({
      name: 'Limpieza',
      description: 'Control preventivo',
      duration_minutes: 30,
      price_from: 120000,
      benefits: '',
      common_objections: '',
      is_active: true,
    });

    fixture.componentInstance.saveService();

    expect(fixture.componentInstance.error()).toBe('La moneda COP no permite decimales en el precio.');
  });

  it('updates the selected service', () => {
    fixture.componentInstance.selectService(services[0]);
    fixture.componentInstance.serviceForm.setValue({
      name: 'Blanqueamiento premium',
      description: 'Tratamiento estetico guiado',
      duration_minutes: 75,
      price_from: 280000,
      benefits: 'Resultado visible',
      common_objections: 'Sensibilidad',
      is_active: false,
    });

    fixture.componentInstance.saveService();

    expect(api.updateService).toHaveBeenCalledWith('service-1', {
      name: 'Blanqueamiento premium',
      description: 'Tratamiento estetico guiado',
      duration_minutes: 75,
      price_from: 280000,
      benefits: ['Resultado visible'],
      faq: [],
      common_objections: ['Sensibilidad'],
      is_active: false,
    });
    expect(fixture.componentInstance.selectedService()?.name).toBe('Blanqueamiento premium');
    expect(fixture.componentInstance.success()).toBe('Servicio actualizado correctamente.');
  });

  it('deletes the selected service', () => {
    fixture.componentInstance.selectService(services[0]);

    fixture.componentInstance.deleteSelected();

    expect(api.deleteService).toHaveBeenCalledWith('service-1');
    expect(fixture.componentInstance.services()).toEqual([]);
    expect(fixture.componentInstance.selectedService()).toBeNull();
    expect(fixture.componentInstance.success()).toBe('Servicio eliminado correctamente.');
  });
});
