import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { environment } from '../../../environments/environment';
import { ApiService } from './api.service';

const baseUrl = environment.apiBaseUrl;

describe('ApiService', () => {
  let service: ApiService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(ApiService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    http.verify();
  });

  it('requests and maps the dashboard summary endpoint', () => {
    service.dashboardSummary().subscribe((summary) => {
      expect(summary).toEqual({
        total_leads: 3,
        pending_followups_today: 1,
        overdue_followups: 0,
        upcoming_followups: 2,
        conversion_rate: 25,
        status_counts: { Nuevo: 2, Contactado: 1 },
        top_services: [{ service_name: 'Blanqueamiento dental', total: 3 }],
      });
    });

    const request = http.expectOne(`${baseUrl}/api/dashboard/summary`);

    expect(request.request.method).toBe('GET');
    request.flush({
      leads_total: 3,
      pending_followups_today: 1,
      overdue_followups: 0,
      upcoming_followups: 2,
      conversion_rate: 25,
      leads_by_status: { Nuevo: 2, Contactado: 1 },
      top_services: [
        {
          service_id: 'service-1',
          service_name: 'Blanqueamiento dental',
          lead_count: 3,
        },
      ],
    });
  });

  it('requests the current clinic endpoint', () => {
    service.clinicCurrent().subscribe();

    const request = http.expectOne(`${baseUrl}/api/clinics/current`);

    expect(request.request.method).toBe('GET');
    request.flush({
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
      phone: '+5710000000',
      whatsapp: '+5710000000',
      address: 'Calle 1',
      communication_tone: 'profesional',
    });
  });

  it('updates the current clinic endpoint', () => {
    const payload = {
      name: 'Sonrisa Viva',
      city: 'Bogota',
      country_code: 'CO' as const,
      currency_code: 'COP' as const,
      phone: '+5710000000',
      whatsapp: '+573001234567',
      address: 'Calle 1',
      opening_hours: {},
      general_faq: [],
      communication_tone: 'profesional' as const,
    };

    service.updateClinicCurrent(payload).subscribe();

    const request = http.expectOne(`${baseUrl}/api/clinics/current`);

    expect(request.request.method).toBe('PUT');
    expect(request.request.body).toEqual(payload);
    request.flush({ id: 'clinic-1', clinic_type: 'dental', ...payload });
  });

  it('requests list endpoints for services, leads, and followups', () => {
    service.services().subscribe((response) => {
      expect(response).toEqual({ data: [] });
    });
    service.leads().subscribe();
    service.followups().subscribe();

    const servicesRequest = http.expectOne(`${baseUrl}/api/services`);
    const leadsRequest = http.expectOne(`${baseUrl}/api/leads`);
    const followupsRequest = http.expectOne(`${baseUrl}/api/followups`);

    expect(servicesRequest.request.method).toBe('GET');
    expect(leadsRequest.request.method).toBe('GET');
    expect(followupsRequest.request.method).toBe('GET');

    servicesRequest.flush([]);
    leadsRequest.flush({ data: [], pagination: { page: 1, page_size: 20, total: 0, total_pages: 0 } });
    followupsRequest.flush({ data: [], pagination: { page: 1, page_size: 20, total: 0, total_pages: 0 } });
  });

  it('posts service catalog create, update, and delete payloads', () => {
    const createPayload = {
      name: 'Ortodoncia',
      description: 'Alineacion dental',
      duration_minutes: 45,
      price_from: 120000,
      benefits: ['Mejora estetica'],
      faq: [],
      common_objections: ['Precio'],
    };
    const updatePayload = { ...createPayload, is_active: true };

    service.createService(createPayload).subscribe();
    service.updateService('service-1', updatePayload).subscribe();
    service.deleteService('service-1').subscribe();

    const createRequest = http.expectOne(`${baseUrl}/api/services`);
    const updateRequest = http.expectOne(
      (request) => request.url === `${baseUrl}/api/services/service-1` && request.method === 'PUT',
    );
    const deleteRequest = http.expectOne(
      (request) => request.url === `${baseUrl}/api/services/service-1` && request.method === 'DELETE',
    );

    expect(createRequest.request.method).toBe('POST');
    expect(createRequest.request.body).toEqual(createPayload);
    expect(updateRequest.request.method).toBe('PUT');
    expect(updateRequest.request.body).toEqual(updatePayload);
    expect(deleteRequest.request.method).toBe('DELETE');

    createRequest.flush({ id: 'service-1', ...createPayload, is_active: true });
    updateRequest.flush({ id: 'service-1', ...updatePayload });
    deleteRequest.flush(null);
  });

  it('posts lead create and update payloads to lead endpoints', () => {
    const createPayload = {
      full_name: 'Lead Demo',
      phone: '+573001112222',
      service_id: 'service-1',
      status: 'Nuevo' as const,
      source: 'whatsapp',
      notes: 'Consulta inicial.',
    };
    const updatePayload = {
      status: 'Contactado' as const,
      note: 'Se contacto por WhatsApp.',
      next_action_at: '2026-05-24T14:30:00.000Z',
    };

    service.createLead(createPayload).subscribe();
    service.updateLead('lead-1', updatePayload).subscribe();

    const createRequest = http.expectOne(`${baseUrl}/api/leads`);
    const updateRequest = http.expectOne(`${baseUrl}/api/leads/lead-1`);

    expect(createRequest.request.method).toBe('POST');
    expect(createRequest.request.body).toEqual(createPayload);
    expect(updateRequest.request.method).toBe('PUT');
    expect(updateRequest.request.body).toEqual(updatePayload);

    createRequest.flush({ id: 'lead-1', ...createPayload, created_at: '2026-05-23T12:00:00Z' });
    updateRequest.flush({ id: 'lead-1', status: 'Contactado' });
  });

  it('posts follow-up actions and AI follow-up message payloads', () => {
    const completePayload = { status: 'Contactado' as const, note: 'Seguimiento completado.' };
    const reschedulePayload = {
      next_action_at: '2026-05-24T14:30:00.000Z',
      note: 'Volver a contactar manana.',
    };
    const messagePayload = {
      lead_id: 'lead-1',
      service_id: 'service-1',
      last_contact_note: 'Pidio informacion de precio.',
    };

    service.completeFollowUp('lead-1', completePayload).subscribe();
    service.rescheduleFollowUp('lead-1', reschedulePayload).subscribe();
    service.followUpMessage(messagePayload).subscribe();

    const completeRequest = http.expectOne(`${baseUrl}/api/followups/lead-1/complete`);
    const rescheduleRequest = http.expectOne(`${baseUrl}/api/followups/lead-1/reschedule`);
    const messageRequest = http.expectOne(`${baseUrl}/api/ai/follow-up-message`);

    expect(completeRequest.request.method).toBe('POST');
    expect(completeRequest.request.body).toEqual(completePayload);
    expect(rescheduleRequest.request.method).toBe('POST');
    expect(rescheduleRequest.request.body).toEqual(reschedulePayload);
    expect(messageRequest.request.method).toBe('POST');
    expect(messageRequest.request.body).toEqual(messagePayload);

    completeRequest.flush({ id: 'lead-1', status: 'completed' });
    rescheduleRequest.flush({ id: 'lead-1', status: 'rescheduled' });
    messageRequest.flush({
      generation_id: 'gen-1',
      suggested_message: 'Hola, seguimos atentos a tu consulta.',
      recommended_timing: 'Hoy',
      next_step: 'Enviar WhatsApp',
      safety_status: 'safe',
    });
  });

  it('posts reply and objection payloads to AI endpoints', () => {
    const replyPayload = { patient_message: 'Quiero informacion', service_id: 'service-1' };
    const objectionPayload = { objection: 'Esta muy caro', lead_id: 'lead-1' };

    service.replySuggestion(replyPayload).subscribe();
    service.objectionHandler(objectionPayload).subscribe();

    const replyRequest = http.expectOne(`${baseUrl}/api/ai/reply-suggestion`);
    const objectionRequest = http.expectOne(`${baseUrl}/api/ai/objection-handler`);

    expect(replyRequest.request.method).toBe('POST');
    expect(replyRequest.request.body).toEqual(replyPayload);
    expect(objectionRequest.request.method).toBe('POST');
    expect(objectionRequest.request.body).toEqual(objectionPayload);

    replyRequest.flush({ suggested_reply: 'Claro, con gusto.' });
    objectionRequest.flush({ suggested_reply: 'Entiendo tu inquietud.' });
  });
  it('posts conversation analysis payload to the AI endpoint', () => {
    const payload = {
      conversation_text: 'Paciente: Hola, quiero saber cuanto cuesta el blanqueamiento.',
      source: 'whatsapp',
      service_id: 'service-1',
    };

    service.analyzeConversation(payload).subscribe();

    const request = http.expectOne(`${baseUrl}/api/ai/analyze-conversation`);

    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(payload);
    request.flush({
      analysis_id: 'analysis-1',
      detected_lead: { full_name: 'Lead Demo', phone: '+573001234567' },
      detected_service: { service_id: 'service-1', service_name: 'Blanqueamiento dental', confidence: 'medium' },
      intent: 'high',
      detected_objections: ['precio'],
      suggested_status: 'Interesado',
      commercial_summary: 'Pregunta por precio.',
      suggested_reply: 'Claro, podemos ayudarte.',
      suggested_next_action: 'Responder y proponer valoracion',
      safety_status: 'passed',
    });
  });

});
