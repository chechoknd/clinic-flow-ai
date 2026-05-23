import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of, throwError } from 'rxjs';

import { ApiService } from '../../core/services/api.service';
import { AiAssistantPage } from './ai-assistant.page';

describe('AiAssistantPage', () => {
  let fixture: ComponentFixture<AiAssistantPage>;
  let replySuggestion: ReturnType<typeof vi.fn>;
  let objectionHandler: ReturnType<typeof vi.fn>;
  let writeText: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    replySuggestion = vi.fn(() => of({ suggested_reply: 'Respuesta comercial sugerida.' }));
    objectionHandler = vi.fn(() => of({ variants: { short: 'Manejo corto de objecion.' } }));
    writeText = vi.fn(() => Promise.resolve());

    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    await TestBed.configureTestingModule({
      imports: [AiAssistantPage],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            queryParamMap: of(convertToParamMap({ lead_id: 'lead-1', service_id: 'service-1' })),
          },
        },
        {
          provide: ApiService,
          useValue: {
            replySuggestion,
            objectionHandler,
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(AiAssistantPage);
    fixture.detectChanges();
  });



  it('shows when lead or service context is active', () => {
    expect(fixture.componentInstance.leadId()).toBe('lead-1');
    expect(fixture.componentInstance.serviceId()).toBe('service-1');
    expect(fixture.componentInstance.hasContext()).toBe(true);
  });

  it('marks the message as required before generating a response', () => {
    fixture.componentInstance.generate();

    expect(replySuggestion).not.toHaveBeenCalled();
    expect(objectionHandler).not.toHaveBeenCalled();
    expect(fixture.componentInstance.form.controls.message.touched).toBe(true);
  });

  it('generates a reply suggestion from the patient message', () => {
    fixture.componentInstance.form.controls.message.setValue('Quiero saber cuanto cuesta el tratamiento.');

    fixture.componentInstance.generate();

    expect(replySuggestion).toHaveBeenCalledWith({
      patient_message: 'Quiero saber cuanto cuesta el tratamiento.',
      lead_id: 'lead-1',
      service_id: 'service-1',
    });
    expect(fixture.componentInstance.answer()).toBe('Respuesta comercial sugerida.');
    expect(fixture.componentInstance.loading()).toBe(false);
  });

  it('uses the objection handler mode and falls back to the first variant', () => {
    fixture.componentInstance.form.setValue({
      mode: 'objection',
      message: 'Me parece muy costoso para hacerlo ahora.',
    });

    fixture.componentInstance.generate();

    expect(objectionHandler).toHaveBeenCalledWith({
      objection: 'Me parece muy costoso para hacerlo ahora.',
      lead_id: 'lead-1',
      service_id: 'service-1',
    });
    expect(fixture.componentInstance.answer()).toBe('Manejo corto de objecion.');
  });

  it('shows a readable error when the AI request fails', () => {
    replySuggestion.mockReturnValueOnce(throwError(() => new Error('provider unavailable')));
    fixture.componentInstance.form.controls.message.setValue('Necesito informacion del servicio.');

    fixture.componentInstance.generate();

    expect(fixture.componentInstance.answer()).toBe(
      'No fue posible generar la respuesta. Verifica el backend o el proveedor AI.',
    );
    expect(fixture.componentInstance.loading()).toBe(false);
  });

  it('copies the generated answer to the clipboard', async () => {
    fixture.componentInstance.answer.set('Texto listo para WhatsApp.');

    fixture.componentInstance.copy();
    await fixture.whenStable();

    expect(writeText).toHaveBeenCalledWith('Texto listo para WhatsApp.');
    expect(fixture.componentInstance.copied()).toBe(true);
  });
});
