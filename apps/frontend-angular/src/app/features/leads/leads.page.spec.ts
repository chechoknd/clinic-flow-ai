import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { ClinicServiceItem, Lead, PaginatedResponse } from '../../core/services/api.models';
import { LeadsPage } from './leads.page';

const leads: Lead[] = [
  {
    id: 'lead-1',
    full_name: 'Paciente Nuevo',
    phone: '+573001112222',
    status: 'Nuevo',
    service_id: 'service-1',
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

const services: ClinicServiceItem[] = [
  {
    id: 'service-1',
    name: 'Blanqueamiento dental',
    description: 'Tratamiento estetico',
    duration_minutes: 60,
    price_from: '250000.00',
    benefits: [],
    faq: [],
    common_objections: [],
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
  let api: {
    leads: ReturnType<typeof vi.fn>;
    lead: ReturnType<typeof vi.fn>;
    services: ReturnType<typeof vi.fn>;
    createLead: ReturnType<typeof vi.fn>;
    updateLead: ReturnType<typeof vi.fn>;
  };
  let writeText: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    writeText = vi.fn(() => Promise.resolve());
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    api = {
      leads: vi.fn(() => of(response)),
      lead: vi.fn((id: string) =>
        of({
          id,
          full_name: 'Paciente Nuevo',
          phone: '+573001112222',
          service: { id: 'service-1', name: 'Blanqueamiento dental' },
          status: 'Nuevo',
          source: 'whatsapp',
          notes: [
            {
              id: 'note-1',
              body: 'Pidio informacion por WhatsApp.',
              created_at: '2026-05-23T12:00:00Z',
            },
          ],
          ai_insights: [
            {
              id: 'insight-1',
              intent: 'high',
              detected_objections: ['precio'],
              commercial_summary: 'Pregunta por precio y quiere avanzar.',
              suggested_next_action: 'Enviar opciones y proponer valoracion.',
              source: 'whatsapp',
              created_at: '2026-05-23T12:05:00Z',
            },
          ],
          next_action_at: '2026-05-24T14:30:00Z',
          created_at: '2026-05-23T12:00:00Z',
        }),
      ),
      services: vi.fn(() => of({ data: services })),
      createLead: vi.fn((lead: Omit<Lead, 'id'>) => of({ ...lead, id: 'lead-3', created_at: '2026-05-23T12:00:00Z' })),
      updateLead: vi.fn(() => of({ id: 'lead-1', status: 'Contactado' })),
    };

    await TestBed.configureTestingModule({
      imports: [LeadsPage],
      providers: [
        provideRouter([]),
        {
          provide: ApiService,
          useValue: api,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LeadsPage);
    fixture.detectChanges();
  });

  it('loads leads and services from the API service', () => {
    expect(fixture.componentInstance.leads()).toEqual(leads);
    expect(fixture.componentInstance.services()).toEqual(services);
  });

  it('shows only leads for the selected status', () => {
    const text = fixture.nativeElement.textContent as string;

    expect(text).toContain('Paciente Nuevo');
    expect(text).not.toContain('Paciente Interesado');
  });

  it('updates the visible lead list when the selected status changes', () => {
    fixture.componentInstance.selectStatus('Interesado');
    fixture.detectChanges();

    const text = fixture.nativeElement.textContent as string;

    expect(text).not.toContain('Paciente Nuevo');
    expect(text).toContain('Paciente Interesado');
  });

  it('creates a lead with form values and selects it', () => {
    fixture.componentInstance.createForm.setValue({
      full_name: 'Nuevo Paciente',
      phone: '+573005551111',
      service_id: 'service-1',
      status: 'Nuevo',
      source: 'whatsapp',
      next_action_at: '',
      notes: 'Quiere informacion de precio.',
    });

    fixture.componentInstance.createLead();

    expect(api.createLead).toHaveBeenCalledWith({
      full_name: 'Nuevo Paciente',
      phone: '+573005551111',
      service_id: 'service-1',
      status: 'Nuevo',
      source: 'whatsapp',
      notes: 'Quiere informacion de precio.',
      next_action_at: undefined,
    });
    expect(fixture.componentInstance.leads()[0].id).toBe('lead-3');
    expect(fixture.componentInstance.selectedLead()?.id).toBe('lead-3');
    expect(fixture.componentInstance.success()).toBe('Lead creado correctamente.');
  });



  it('builds an AI assistant link with selected lead context', () => {
    fixture.componentInstance.selectLead(leads[0]);
    fixture.detectChanges();

    const links = Array.from(
      fixture.nativeElement.querySelectorAll('a[href*="ai-assistant"]'),
    ) as HTMLAnchorElement[];
    const contextualLink = links.find((link) => link.href.includes('lead_id=lead-1'));

    expect(contextualLink?.href).toContain('lead_id=lead-1');
    expect(contextualLink?.href).toContain('service_id=service-1');
  });

  it('loads selected lead detail with notes and reviewed AI insight', () => {
    fixture.componentInstance.selectLead(leads[0]);
    fixture.detectChanges();

    const text = fixture.nativeElement.textContent as string;

    expect(api.lead).toHaveBeenCalledWith('lead-1');
    expect(text).toContain('Detalle comercial');
    expect(text).toContain('Alta intencion');
    expect(text).toContain('Pregunta por precio y quiere avanzar.');
    expect(text).toContain('Pidio informacion por WhatsApp.');
  });

  it('opens manual WhatsApp and copies selected lead phone', async () => {
    fixture.componentInstance.selectLead(leads[0]);
    fixture.detectChanges();

    const whatsappLink = fixture.nativeElement.querySelector('[data-testid="leads-detail-whatsapp-link"]') as HTMLAnchorElement;

    expect(whatsappLink.href).toBe('https://wa.me/573001112222');

    fixture.componentInstance.copySelectedLeadPhone();
    await fixture.whenStable();

    expect(writeText).toHaveBeenCalledWith('+573001112222');
    expect(fixture.componentInstance.phoneCopied()).toBe(true);
  });

  it('applies quick post-contact actions with status, note, and follow-up date', () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-05-27T15:00:00.000Z'));
    fixture.componentInstance.selectLead(leads[0]);

    fixture.componentInstance.applyQuickAction(fixture.componentInstance.quickActions[1]);

    expect(api.updateLead).toHaveBeenCalledWith('lead-1', {
      status: 'Interesado',
      note: 'Lead interesado. Programar seguimiento comercial en 2 días.',
      next_action_at: '2026-05-29T15:00:00.000Z',
      clear_next_action_at: undefined,
    });
    expect(fixture.componentInstance.selectedStatus()).toBe('Interesado');
    expect(fixture.componentInstance.selectedLead()?.next_action_at).toBe('2026-05-29T15:00:00.000Z');
    expect(fixture.componentInstance.success()).toBe('Accion rapida aplicada correctamente.');
    vi.useRealTimers();
  });

  it('clears next action when marking a lead as lost', () => {
    fixture.componentInstance.selectLead({ ...leads[0], next_action_at: '2026-05-29T15:00:00.000Z' });

    fixture.componentInstance.applyQuickAction(fixture.componentInstance.quickActions[3]);

    expect(api.updateLead).toHaveBeenCalledWith('lead-1', {
      status: 'Perdido',
      note: 'Lead marcado como perdido. No requiere seguimiento por ahora.',
      next_action_at: undefined,
      clear_next_action_at: true,
    });
    expect(fixture.componentInstance.selectedStatus()).toBe('Perdido');
    expect(fixture.componentInstance.selectedLead()?.next_action_at).toBeUndefined();
  });

  it('updates the selected lead status, note, and next action', () => {
    fixture.componentInstance.selectLead(leads[0]);
    fixture.componentInstance.updateForm.setValue({
      status: 'Contactado',
      next_action_at: '2026-05-24T09:30',
      note: 'Se envio informacion por WhatsApp.',
    });

    fixture.componentInstance.updateLead();

    expect(api.updateLead).toHaveBeenCalledWith('lead-1', {
      status: 'Contactado',
      note: 'Se envio informacion por WhatsApp.',
      next_action_at: '2026-05-24T14:30:00.000Z',
    });
    expect(fixture.componentInstance.selectedStatus()).toBe('Contactado');
    expect(fixture.componentInstance.selectedLead()?.status).toBe('Contactado');
    expect(fixture.componentInstance.success()).toBe('Lead actualizado correctamente.');
  });
});
