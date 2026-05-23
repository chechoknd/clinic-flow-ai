import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { FollowUp, PaginatedResponse } from '../../core/services/api.models';
import { FollowupsPage } from './followups.page';

const followups: FollowUp[] = [
  {
    id: 'lead-1',
    full_name: 'Paciente Seguimiento',
    phone: '+573001112222',
    status: 'Contactado',
    service_id: 'service-1',
    service_name: 'Blanqueamiento dental',
    source: 'whatsapp',
    next_action_at: '2026-05-24T14:30:00.000Z',
  },
];

const response: PaginatedResponse<FollowUp> = {
  data: followups,
  pagination: { page: 1, page_size: 20, total: 1, total_pages: 1 },
};

describe('FollowupsPage', () => {
  let fixture: ComponentFixture<FollowupsPage>;
  let api: {
    followups: ReturnType<typeof vi.fn>;
    completeFollowUp: ReturnType<typeof vi.fn>;
    rescheduleFollowUp: ReturnType<typeof vi.fn>;
    followUpMessage: ReturnType<typeof vi.fn>;
  };
  let writeText: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    writeText = vi.fn(() => Promise.resolve());
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    api = {
      followups: vi.fn(() => of(response)),
      completeFollowUp: vi.fn(() => of({ id: 'lead-1', status: 'completed' })),
      rescheduleFollowUp: vi.fn(() => of({ id: 'lead-1', status: 'rescheduled' })),
      followUpMessage: vi.fn(() =>
        of({
          generation_id: 'gen-1',
          suggested_message: 'Mensaje de seguimiento sugerido.',
          recommended_timing: 'Hoy',
          next_step: 'Enviar WhatsApp',
          safety_status: 'safe',
        }),
      ),
    };

    await TestBed.configureTestingModule({
      imports: [FollowupsPage],
      providers: [
        provideRouter([]),
        {
          provide: ApiService,
          useValue: api,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(FollowupsPage);
    fixture.detectChanges();
  });

  it('loads pending follow-ups', () => {
    expect(fixture.componentInstance.followups()).toEqual(followups);
  });

  it('selects a follow-up and prepares the action form', () => {
    fixture.componentInstance.selectFollowUp(followups[0]);

    expect(fixture.componentInstance.selectedFollowUp()).toEqual(followups[0]);
    expect(fixture.componentInstance.actionForm.controls.status.value).toBe('Contactado');
    expect(fixture.componentInstance.canGenerateMessage()).toBe(true);
  });

  it('completes the selected follow-up and removes it from the list', () => {
    fixture.componentInstance.selectFollowUp(followups[0]);
    fixture.componentInstance.actionForm.controls.note.setValue('Se envio informacion.');

    fixture.componentInstance.completeSelected();

    expect(api.completeFollowUp).toHaveBeenCalledWith('lead-1', {
      status: 'Contactado',
      note: 'Se envio informacion.',
    });
    expect(fixture.componentInstance.followups()).toEqual([]);
    expect(fixture.componentInstance.success()).toBe('Seguimiento completado correctamente.');
  });

  it('reschedules the selected follow-up', () => {
    fixture.componentInstance.selectFollowUp(followups[0]);
    const nextActionAt = new Date('2026-05-25T10:00').toISOString();
    fixture.componentInstance.actionForm.setValue({
      status: 'Contactado',
      next_action_at: '2026-05-25T10:00',
      note: 'Reprogramar contacto.',
    });

    fixture.componentInstance.rescheduleSelected();

    expect(api.rescheduleFollowUp).toHaveBeenCalledWith('lead-1', {
      next_action_at: nextActionAt,
      note: 'Reprogramar contacto.',
    });
    expect(fixture.componentInstance.selectedFollowUp()?.next_action_at).toBe(nextActionAt);
    expect(fixture.componentInstance.success()).toBe('Seguimiento reprogramado correctamente.');
  });

  it('generates and copies an AI follow-up message', async () => {
    fixture.componentInstance.selectFollowUp(followups[0]);
    fixture.componentInstance.actionForm.controls.note.setValue('Pidio precio.');

    fixture.componentInstance.generateMessage();

    expect(api.followUpMessage).toHaveBeenCalledWith({
      lead_id: 'lead-1',
      service_id: 'service-1',
      last_contact_note: 'Pidio precio.',
    });
    expect(fixture.componentInstance.generatedMessage()).toBe('Mensaje de seguimiento sugerido.');

    fixture.componentInstance.copyGeneratedMessage();
    await fixture.whenStable();

    expect(writeText).toHaveBeenCalledWith('Mensaje de seguimiento sugerido.');
    expect(fixture.componentInstance.copied()).toBe(true);
  });
});
