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
    analysis_id: '11111111-1111-4111-8111-111111111111',
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

  it('shows a human review preview before saving the analyzed lead action', () => {
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: 'service-1',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });

    fixture.componentInstance.analyze();
    fixture.detectChanges();

    const preview = fixture.nativeElement.querySelector('[data-testid="inbox-ai-review-preview"]') as HTMLElement;
    const lead = fixture.nativeElement.querySelector('[data-testid="inbox-ai-preview-lead"]') as HTMLElement;
    const action = fixture.nativeElement.querySelector('[data-testid="inbox-ai-preview-action"]') as HTMLElement;
    const nextStep = fixture.nativeElement.querySelector('[data-testid="inbox-ai-preview-next-step"]') as HTMLElement;
    const completion = fixture.nativeElement.querySelector('[data-testid="inbox-ai-review-completion"]') as HTMLElement;
    const missingFields = fixture.nativeElement.querySelector('[data-testid="inbox-ai-review-missing-fields"]');

    expect(preview.textContent).toContain('Vista previa antes de guardar');
    expect(completion.textContent).toContain('Revision completa');
    expect(missingFields).toBeNull();
    expect(lead.textContent).toContain('Lead Demo Inbox');
    expect(action.textContent).toContain('Crear lead nuevo');
    expect(preview.textContent).toContain('Interesado - Blanqueamiento dental');
    expect(nextStep.textContent).toContain('Responder y proponer valoracion');
  });

  it('shows pending review when required reviewed lead fields are invalid', () => {
    fixture.componentInstance.analysis.set(analysisResponse);
    fixture.componentInstance.leadForm.patchValue({
      full_name: '',
      phone: '3001234567',
      status: 'Interesado',
    });
    fixture.detectChanges();

    const completion = fixture.nativeElement.querySelector('[data-testid="inbox-ai-review-completion"]') as HTMLElement;
    const missingFields = fixture.nativeElement.querySelector('[data-testid="inbox-ai-review-missing-fields"]') as HTMLElement;

    expect(fixture.componentInstance.isReviewComplete()).toBe(false);
    expect(completion.textContent).toContain('Revisa nombre, WhatsApp y estado');
    expect(missingFields.textContent).toContain('Nombre');
    expect(missingFields.textContent).toContain('WhatsApp');
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
      reviewed_ai_analysis: {
        analysis_id: '11111111-1111-4111-8111-111111111111',
        intent: 'high',
        detected_objections: ['precio'],
        commercial_summary: 'Pregunta por precio y quiere informacion.',
        suggested_next_action: 'Responder y proponer valoracion',
        source: 'whatsapp',
      },
    });
    expect(fixture.componentInstance.success()).toBe('Lead creado desde el analisis revisado.');
    expect(fixture.componentInstance.reviewedLeadID()).toBe('lead-1');
  });

  it('links to the reviewed lead detail after creating a lead', () => {
    fixture.componentInstance.analyzeForm.setValue({
      source: 'whatsapp',
      lead_id: '',
      service_id: 'service-1',
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento dental.',
    });
    fixture.componentInstance.analyze();

    fixture.componentInstance.createLead();
    fixture.detectChanges();

    const detailLink = fixture.nativeElement.querySelector('[data-testid="inbox-ai-reviewed-lead-detail-link"]') as HTMLAnchorElement;

    expect(detailLink.href).toContain('/leads');
    expect(detailLink.href).toContain('lead_id=lead-1');
    expect(detailLink.href).toContain('service_id=service-1');
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
      reviewed_ai_analysis: {
        analysis_id: '11111111-1111-4111-8111-111111111111',
        intent: 'high',
        detected_objections: ['precio'],
        commercial_summary: 'Pregunta por precio y quiere informacion.',
        suggested_next_action: 'Responder y proponer valoracion',
        source: 'whatsapp',
      },
    });
    expect(fixture.componentInstance.success()).toBe('Lead actualizado con el analisis revisado.');
    expect(fixture.componentInstance.reviewedLeadID()).toBe('lead-existing');
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
