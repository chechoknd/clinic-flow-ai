import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { Lead, PaginatedResponse } from '../../core/services/api.models';
import { LeadsPage } from './leads.page';

const leads: Lead[] = [
  {
    id: 'lead-1',
    full_name: 'Paciente Nuevo',
    phone: '+573001112222',
    status: 'Nuevo',
    service_name: 'Blanqueamiento dental',
  },
  {
    id: 'lead-2',
    full_name: 'Paciente Interesado',
    phone: '+573003334444',
    status: 'Interesado',
    service_name: 'Ortodoncia',
  },
];

const response: PaginatedResponse<Lead> = {
  data: leads,
  pagination: {
    page: 1,
    page_size: 20,
    total: 2,
    total_pages: 1,
  },
};

describe('LeadsPage', () => {
  let fixture: ComponentFixture<LeadsPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LeadsPage],
      providers: [
        provideRouter([]),
        {
          provide: ApiService,
          useValue: {
            leads: () => of(response),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LeadsPage);
    fixture.detectChanges();
  });

  it('loads leads from the API service', () => {
    expect(fixture.componentInstance.leads()).toEqual(leads);
  });

  it('shows only leads for the selected status', () => {
    const text = fixture.nativeElement.textContent as string;

    expect(text).toContain('Paciente Nuevo');
    expect(text).not.toContain('Paciente Interesado');
  });

  it('updates the visible lead list when the selected status changes', () => {
    fixture.componentInstance.selectedStatus.set('Interesado');
    fixture.detectChanges();

    const text = fixture.nativeElement.textContent as string;

    expect(text).not.toContain('Paciente Nuevo');
    expect(text).toContain('Paciente Interesado');
  });
});
