import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { ClinicProfile } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';
import { ClinicPage } from './clinic.page';

const clinic: ClinicProfile = {
  id: 'clinic-1',
  name: 'Sonrisa Viva',
  clinic_type: 'dental',
  city: 'Bogota',
  phone: '+5710000000',
  whatsapp: '+573001234567',
  address: 'Calle 1',
  opening_hours: { monday: '08:00-17:00' },
  general_faq: [],
  communication_tone: 'profesional',
};

describe('ClinicPage', () => {
  let fixture: ComponentFixture<ClinicPage>;
  let api: {
    clinicCurrent: ReturnType<typeof vi.fn>;
    updateClinicCurrent: ReturnType<typeof vi.fn>;
  };

  beforeEach(async () => {
    api = {
      clinicCurrent: vi.fn(() => of(clinic)),
      updateClinicCurrent: vi.fn((payload) => of({ ...clinic, ...payload })),
    };

    await TestBed.configureTestingModule({
      imports: [ClinicPage],
      providers: [
        {
          provide: ApiService,
          useValue: api,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ClinicPage);
    fixture.detectChanges();
  });

  it('loads and displays the clinic profile', () => {
    expect(fixture.componentInstance.clinic()).toEqual(clinic);
    expect(fixture.componentInstance.clinicForm.controls.name.value).toBe('Sonrisa Viva');

    const text = fixture.nativeElement.textContent as string;
    expect(text).toContain('Sonrisa Viva');
    expect(text).toContain('+573001234567');
  });

  it('validates required commercial fields', () => {
    fixture.componentInstance.clinicForm.controls.name.setValue('');

    fixture.componentInstance.saveClinic();

    expect(api.updateClinicCurrent).not.toHaveBeenCalled();
    expect(fixture.componentInstance.clinicForm.controls.name.touched).toBe(true);
  });

  it('updates the clinic profile with normalized optional fields', () => {
    fixture.componentInstance.clinicForm.setValue({
      name: 'Sonrisa Viva Norte',
      city: 'Medellin',
      phone: '',
      whatsapp: '+573009998877',
      address: 'Carrera 10',
      communication_tone: 'cercano',
    });

    fixture.componentInstance.saveClinic();

    expect(api.updateClinicCurrent).toHaveBeenCalledWith({
      name: 'Sonrisa Viva Norte',
      city: 'Medellin',
      phone: undefined,
      whatsapp: '+573009998877',
      address: 'Carrera 10',
      opening_hours: { monday: '08:00-17:00' },
      general_faq: [],
      communication_tone: 'cercano',
    });
    expect(fixture.componentInstance.clinic()?.name).toBe('Sonrisa Viva Norte');
    expect(fixture.componentInstance.success()).toBe('Perfil de clinica actualizado correctamente.');
  });
});
