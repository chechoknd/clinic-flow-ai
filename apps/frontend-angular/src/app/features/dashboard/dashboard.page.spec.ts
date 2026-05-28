import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { DashboardPage } from './dashboard.page';

describe('DashboardPage', () => {
  let fixture: ComponentFixture<DashboardPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DashboardPage],
      providers: [
        provideRouter([]),
        {
          provide: ApiService,
          useValue: {
            dashboardSummary: vi.fn(() =>
              of({
                total_leads: 1,
                pending_followups_today: 0,
                overdue_followups: 0,
                upcoming_followups: 1,
                conversion_rate: 0,
                status_counts: { Interesado: 1 },
                top_services: [{ service_name: 'Blanqueamiento dental', total: 1 }],
              }),
            ),
            dashboardActions: vi.fn(() =>
              of({
                data: [
                  {
                    type: 'overdue_followup',
                    tone: 'urgent',
                    priority: 100,
                    lead_id: 'lead-2',
                    full_name: 'Lead Vencido',
                    phone: '+573001112223',
                    service_id: 'service-1',
                    service_name: 'Blanqueamiento dental',
                    status: 'Interesado',
                    source: 'whatsapp',
                    reason: 'Seguimiento vencido',
                    next_action_at: '2026-05-26T14:30:00Z',
                    created_at: '2026-05-25T12:00:00Z',
                  },
                  {
                    type: 'high_intent',
                    tone: 'intent',
                    priority: 70,
                    lead_id: 'lead-1',
                    full_name: 'Lead Alta Intencion',
                    phone: '+573001112222',
                    service_id: 'service-1',
                    service_name: 'Blanqueamiento dental',
                    status: 'Interesado',
                    source: 'whatsapp',
                    reason: 'Alta intencion detectada por Inbox AI',
                    next_action_at: '2026-05-30T14:30:00Z',
                    created_at: '2026-05-27T12:00:00Z',
                  },
                ],
              }),
            ),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardPage);
    fixture.detectChanges();
  });

  it('links insight actions to the lead commercial detail', () => {
    const links = Array.from(fixture.nativeElement.querySelectorAll('a[href*="leads"]')) as HTMLAnchorElement[];
    const detailLink = links.find((link) => link.href.includes('lead_id=lead-1'));

    expect(detailLink?.href).toContain('/leads');
    expect(detailLink?.href).toContain('lead_id=lead-1');
    expect(detailLink?.href).toContain('service_id=service-1');
  });

  it('shows queue summary counts by action type', () => {
    const text = fixture.nativeElement.textContent as string;

    expect(text).toContain('Vencidos');
    expect(text).toContain('Alta intencion');
    expect(text).toContain('Objeciones');
  });
});
