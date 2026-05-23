import { Component, signal, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-ai-assistant-page',
  imports: [ReactiveFormsModule],
  templateUrl: './ai-assistant.page.html',
  styleUrl: './ai-assistant.page.css',
})
export class AiAssistantPage {
  private readonly fb = inject(FormBuilder);
  private readonly api = inject(ApiService);

  readonly loading = signal(false);
  readonly answer = signal('');
  readonly copied = signal(false);
  readonly form = this.fb.nonNullable.group({
    mode: ['reply'],
    message: ['', [Validators.required, Validators.minLength(8)]],
  });

  generate(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.loading.set(true);
    this.copied.set(false);
    const value = this.form.getRawValue();
    const request =
      value.mode === 'objection'
        ? this.api.objectionHandler({ objection: value.message })
        : this.api.replySuggestion({ patient_message: value.message });

    request.subscribe({
      next: (response) => {
        this.answer.set(
          response.suggested_reply ||
            response.message ||
            this.firstVariant(response.variants) ||
            'No se recibio respuesta.',
        );
        this.loading.set(false);
      },
      error: () => {
        this.answer.set('No fue posible generar la respuesta. Verifica el backend o el proveedor AI.');
        this.loading.set(false);
      },
    });
  }

  copy(): void {
    if (!this.answer()) {
      return;
    }

    void navigator.clipboard.writeText(this.answer()).then(() => this.copied.set(true));
  }

  private firstVariant(variants?: Record<string, string>): string | null {
    return variants ? Object.values(variants)[0] ?? null : null;
  }
}
