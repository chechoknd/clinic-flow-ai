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
      phone: '+5710000000',
      whatsapp: '+5710000000',
      address: 'Calle 1',
      communication_tone: 'calido',
    });
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
});
