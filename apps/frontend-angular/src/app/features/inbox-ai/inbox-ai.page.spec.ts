import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of, throwError } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { InboxAiPage } from './inbox-ai.page';

describe('InboxAiPage', () => {
  let fixture: ComponentFixture<InboxAiPage>;
  let analyzeConversation: ReturnType<typeof vi.fn>;
  let createLead: ReturnType<typeof vi.fn>;
  let updateLead: ReturnType<typeof vi.fn>;
  let writeText: ReturnType<typeof vi.fn>;

  const analysisResponse = {
    analysis_id: 'analysis-1',
    detected_lead: { full_name: 'Lead Demo Inbox', phone: '+573001234567' },
    detected_service: { service_id: 'service-1', service_name: 'Blanqueamiento dental', confidence: 'medium' },
    intent: 'high',
    detected_objections: ['precio'],
    suggested_status: 'Interesado' as const,
    commercial_summary: 'Pregunta por precio y quiere informacion.',
    suggested_reply: 'Claro, podemos orientarte y agendar una valoracion.',
    suggested_next_action: 'Responder y proponer valoracion',
    safety_status: 'passed',
  };

  beforeEach(async () => {
    analyzeConversation = vi.fn(() => of(analysisResponse));
    createLead = vi.fn(() =>
      of({
        id: 'lead-1',
        full_name: 'Lead Demo Inbox',
        phone: '+573001234567',
        service_id: 'service-1',
        status: 'Interesado',
        source: 'whatsapp',
      }),
    );
    updateLead = vi.fn(() => of({ id: 'lead-existing', status: 'Interesado' }));
    writeText = vi.fn(() => Promise.resolve());

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    await TestBed.configureTestingModule({
      imports: [InboxAiPage],
      providers: [
        provideRouter([]),
        {
          provide: ApiService,
          useValue: {
            services: vi.fn(() => of({ data: [{ id: 'service-1', name: 'Blanqueamiento dental', benefits: [], faq: [], common_objections: [] }] })),
            leads: vi.fn(() =>
              of({
                data: [
                  {
                    id: 'lead-existing',
                    full_name: 'Lead Existente',
                    phone: '+573009998888',
                    service_id: 'service-1',
                    status: 'Nuevo',
                    source: 'whatsapp',
                  },
                ],
                pagination: { page: 1, page_size: 20, total: 1, total_pages: 1 },
              }),
            ),
            analyzeConversation,
            createLead,
            updateLead,
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(InboxAiPage);
    fixture.detectChanges();
  });

  it('requires a conversation before analyzing', () => {
    fixture.componentInstance.analyze();

    expect(analyzeConversation).not.toHaveBeenCalled();
    expect(fixture.componentInstance.analyzeForm.controls.conversation_text.touched).toBe(true);
  });

  it('analyzes a pasted conversation and hydrates the reviewed lead form', () => {
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: 'service-1',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });

    fixture.componentInstance.analyze();

    expect(analyzeConversation).toHaveBeenCalledWith({
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
      source: 'whatsapp',
      lead_id: undefined,
      service_id: 'service-1',
    });
    expect(fixture.componentInstance.analysis()).toEqual(analysisResponse);
    expect(fixture.componentInstance.leadForm.getRawValue().full_name).toBe('Lead Demo Inbox');
    expect(fixture.componentInstance.leadForm.getRawValue().status).toBe('Interesado');
  });

  it('creates a lead only after analysis and human-reviewed fields', () => {
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: 'service-1',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });
    fixture.componentInstance.analyze();

    fixture.componentInstance.createLead();

    expect(createLead).toHaveBeenCalledWith({
      full_name: 'Lead Demo Inbox',
      phone: '+573001234567',
      service_id: 'service-1',
      status: 'Interesado',
      source: 'whatsapp',
      notes: 'Pregunta por precio y quiere informacion.',
      next_action_at: undefined,
    });
    expect(fixture.componentInstance.success()).toBe('Lead creado desde el analisis revisado. Ya puedes agendarlo desde el boton de agenda.');
  });

  it('updates an existing lead with the reviewed analysis action', () => {
    fixture.componentInstance.selectExistingLead('lead-existing');
    fixture.componentInstance.analyzeForm.patchValue({
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });
    fixture.componentInstance.analyze();

    fixture.componentInstance.updateExistingLead();

    expect(updateLead).toHaveBeenCalledWith('lead-existing', {
      status: 'Interesado',
      note: 'Pregunta por precio y quiere informacion.',
      next_action_at: undefined,
    });
    expect(fixture.componentInstance.success()).toBe('Lead actualizado con el analisis revisado.');
  });


  it('builds reviewed schedule handoff query params after analysis', () => {
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: 'service-1',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });

    fixture.componentInstance.analyze();

    expect(fixture.componentInstance.canOpenSchedule()).toBe(true);
    expect(fixture.componentInstance.scheduleQueryParams()).toEqual({
      contact_name: 'Lead Demo Inbox',
      contact_phone: '+573001234567',
      service_id: 'service-1',
      lead_id: undefined,
      starts_at_date: undefined,
      admin_notes: 'Pregunta por precio y quiere informacion.',
    });
  });

  it('copies the suggested reply to clipboard', async () => {
    fixture.componentInstance.analysis.set(analysisResponse);

    fixture.componentInstance.copySuggestedReply();
    await fixture.whenStable();

    expect(writeText).toHaveBeenCalledWith('Claro, podemos orientarte y agendar una valoracion.');
    expect(fixture.componentInstance.copied()).toBe(true);
  });

  it('shows an error when analysis fails', () => {
    analyzeConversation.mockReturnValueOnce(throwError(() => new Error('provider unavailable')));
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: '',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });

    fixture.componentInstance.analyze();

    expect(fixture.componentInstance.error()).toBe('No fue posible analizar la conversacion. Revisa el texto o el proveedor AI.');
  });
});
