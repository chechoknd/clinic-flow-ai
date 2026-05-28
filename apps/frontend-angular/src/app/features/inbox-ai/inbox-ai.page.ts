import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';

import {
  AnalyzeConversationResponse,
  ClinicServiceItem,
  Lead,
  LeadStatus,
  ReviewedAIAnalysisPayload,
} from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-inbox-ai-page',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './inbox-ai.page.html',
  styleUrl: './inbox-ai.page.css',
})
export class InboxAiPage {
  private readonly api = inject(ApiService);
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute, { optional: true });

  readonly services = signal<ClinicServiceItem[]>([]);
  readonly leads = signal<Lead[]>([]);
  readonly analysis = signal<AnalyzeConversationResponse | null>(null);
  readonly loading = signal(false);
  readonly savingLead = signal(false);
  readonly copied = signal(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);
  readonly reviewedLeadID = signal<string | null>(null);

  readonly analyzeForm = this.fb.nonNullable.group({
    source: ['whatsapp', Validators.required],
    lead_id: [''],
    service_id: [''],
    conversation_text: ['', [Validators.required, Validators.minLength(20), Validators.maxLength(8000)]],
  });

  readonly leadForm = this.fb.nonNullable.group({
    full_name: ['', [Validators.required, Validators.minLength(3)]],
    phone: ['', [Validators.required, Validators.pattern(/^\+[1-9]\d{7,14}$/)]],
    service_id: [''],
    status: ['Nuevo' as LeadStatus, Validators.required],
    source: ['whatsapp', Validators.required],
    notes: [''],
    next_action_at: [''],
  });

  constructor() {
    this.applyRouteContext();
    this.loadInitialData();
  }

  onExistingLeadChange(leadID: string): void {
    if (!leadID) {
      this.clearExistingLead();
      return;
    }

    this.selectExistingLead(leadID);
  }

  selectExistingLead(leadID: string): void {
    const lead = this.leads().find((item) => item.id === leadID);
    if (!lead) {
      return;
    }

    this.analyzeForm.patchValue({
      lead_id: lead.id,
      service_id: lead.service_id || '',
    });
    this.leadForm.patchValue({
      full_name: lead.full_name,
      phone: lead.phone,
      service_id: lead.service_id || '',
      status: lead.status,
      source: lead.source || 'whatsapp',
      next_action_at: this.toLocalDateTimeValue(lead.next_action_at),
    });
  }

  clearExistingLead(): void {
    this.analyzeForm.patchValue({ lead_id: '' });
  }

  selectedLeadName(): string {
    const leadID = this.analyzeForm.controls.lead_id.value;
    const lead = this.leads().find((item) => item.id === leadID);
    return lead ? `${lead.full_name} - ${lead.phone}` : '';
  }

  private loadInitialData(): void {
    this.api.services().subscribe({
      next: (response) => this.services.set(response.data),
      error: () => this.error.set('No fue posible cargar los servicios.'),
    });

    this.api.leads().subscribe({
      next: (response) => {
        this.leads.set(response.data);
        const leadID = this.analyzeForm.controls.lead_id.value;
        if (leadID) {
          this.selectExistingLead(leadID);
        }
      },
      error: () => this.error.set('No fue posible cargar los leads existentes.'),
    });
  }

  analyze(): void {
    if (this.analyzeForm.invalid) {
      this.analyzeForm.markAllAsTouched();
      return;
    }

    const value = this.analyzeForm.getRawValue();
    this.loading.set(true);
    this.copied.set(false);
    this.reviewedLeadID.set(null);
    this.error.set(null);
    this.success.set(null);

    this.api
      .analyzeConversation({
        conversation_text: value.conversation_text.trim(),
        source: value.source,
        lead_id: value.lead_id || undefined,
        service_id: value.service_id || undefined,
      })
      .subscribe({
        next: (analysis) => {
          this.analysis.set(analysis);
          this.hydrateLeadForm(analysis, value.source, value.service_id);
          this.loading.set(false);
        },
        error: () => {
          this.error.set('No fue posible analizar la conversacion. Revisa el texto o el proveedor AI.');
          this.loading.set(false);
        },
      });
  }

  hasExistingLead(): boolean {
    return Boolean(this.analyzeForm.controls.lead_id.value);
  }

  primaryLeadActionLabel(): string {
    if (this.savingLead()) {
      return this.hasExistingLead() ? 'Actualizando...' : 'Guardando...';
    }
    return this.hasExistingLead() ? 'Actualizar lead revisado' : 'Crear lead revisado';
  }

  reviewedActionModeLabel(): string {
    return this.hasExistingLead() ? 'Actualizar lead existente' : 'Crear lead nuevo';
  }

  reviewedServiceName(): string {
    const serviceID = this.leadForm.controls.service_id.value;
    const service = this.services().find((item) => item.id === serviceID);
    return service?.name || this.analysis()?.detected_service.service_name || 'Sin servicio definido';
  }

  reviewedNextStep(): string {
    const nextActionAt = this.leadForm.controls.next_action_at.value;
    if (nextActionAt) {
      return `Seguimiento: ${this.formatDate(nextActionAt)}`;
    }
    return this.analysis()?.suggested_next_action || 'Sin accion sugerida';
  }

  saveReviewedLead(): void {
    if (this.hasExistingLead()) {
      this.updateExistingLead();
      return;
    }
    this.createLead();
  }

  isLeadReady(): boolean {
    const value = this.leadForm.getRawValue();
    return Boolean(this.analysis() && value.full_name.trim() && value.phone.trim());
  }

  isReviewComplete(): boolean {
    return Boolean(this.analysis() && this.leadForm.valid && this.isLeadReady());
  }

  reviewCompletionLabel(): string {
    return this.isReviewComplete() ? 'Revision completa' : 'Revisa nombre, WhatsApp y estado';
  }

  reviewMissingFields(): string[] {
    const fields: string[] = [];
    const value = this.leadForm.getRawValue();
    if (!value.full_name.trim()) {
      fields.push('Nombre');
    }
    if (!value.phone.trim() || this.leadForm.controls.phone.invalid) {
      fields.push('WhatsApp');
    }
    if (!value.status) {
      fields.push('Estado');
    }
    return fields;
  }

  reviewedLeadQueryParams(): Record<string, string> {
    const params: Record<string, string> = {};
    const leadID = this.reviewedLeadID();
    if (leadID) {
      params['lead_id'] = leadID;
    }
    const serviceID = this.leadForm.controls.service_id.value;
    if (serviceID) {
      params['service_id'] = serviceID;
    }
    params['from'] = 'inbox_ai';
    return params;
  }

  createLead(): void {
    if (this.leadForm.invalid || !this.analysis()) {
      this.leadForm.markAllAsTouched();
      return;
    }

    const value = this.leadForm.getRawValue();
    this.savingLead.set(true);
    this.reviewedLeadID.set(null);
    this.error.set(null);
    this.success.set(null);

    this.api
      .createLead({
        full_name: value.full_name.trim(),
        phone: value.phone.trim(),
        service_id: value.service_id || undefined,
        status: value.status,
        source: value.source,
        notes: value.notes.trim() || undefined,
        next_action_at: this.toApiDateTime(value.next_action_at),
        reviewed_ai_analysis: this.reviewedAIAnalysisPayload(),
      })
      .subscribe({
        next: (lead) => {
          this.leads.update((items) => [lead, ...items]);
          this.analyzeForm.patchValue({ lead_id: lead.id });
          this.reviewedLeadID.set(lead.id);
          this.success.set('Lead creado desde el analisis revisado.');
          this.savingLead.set(false);
        },
        error: () => {
          this.error.set('No fue posible crear el lead. Revisa los datos detectados.');
          this.savingLead.set(false);
        },
      });
  }

  updateExistingLead(): void {
    const leadID = this.analyzeForm.controls.lead_id.value;
    if (!leadID || this.leadForm.invalid || !this.analysis()) {
      this.leadForm.markAllAsTouched();
      return;
    }

    const value = this.leadForm.getRawValue();
    this.savingLead.set(true);
    this.reviewedLeadID.set(null);
    this.error.set(null);
    this.success.set(null);

    this.api
      .updateLead(leadID, {
        status: value.status,
        note: value.notes.trim() || undefined,
        next_action_at: this.toApiDateTime(value.next_action_at),
        reviewed_ai_analysis: this.reviewedAIAnalysisPayload(),
      })
      .subscribe({
        next: () => {
          this.leads.update((items) =>
            items.map((lead) =>
              lead.id === leadID
                ? {
                    ...lead,
                    status: value.status,
                    next_action_at: this.toApiDateTime(value.next_action_at),
                  }
                : lead,
            ),
          );
          this.reviewedLeadID.set(leadID);
          this.success.set('Lead actualizado con el analisis revisado.');
          this.savingLead.set(false);
        },
        error: () => {
          this.error.set('No fue posible actualizar el lead. Revisa los datos detectados.');
          this.savingLead.set(false);
        },
      });
  }

  copySuggestedReply(): void {
    const text = this.analysis()?.suggested_reply;
    if (!text) {
      return;
    }

    void navigator.clipboard.writeText(text).then(() => this.copied.set(true));
  }

  private applyRouteContext(): void {
    const leadID = this.route?.snapshot.queryParamMap.get('lead_id') || '';
    const serviceID = this.route?.snapshot.queryParamMap.get('service_id') || '';
    this.analyzeForm.patchValue({
      lead_id: leadID,
      service_id: serviceID,
    });
    this.leadForm.patchValue({ service_id: serviceID });
  }

  private hydrateLeadForm(analysis: AnalyzeConversationResponse, source: string, selectedServiceID: string): void {
    const detectedServiceID = analysis.detected_service.service_id || selectedServiceID;
    const existingLeadID = this.analyzeForm.controls.lead_id.value;
    const existingLead = this.leads().find((lead) => lead.id === existingLeadID);
    this.leadForm.reset({
      full_name: analysis.detected_lead.full_name || existingLead?.full_name || '',
      phone: analysis.detected_lead.phone || existingLead?.phone || '',
      service_id: detectedServiceID || existingLead?.service_id || '',
      status: analysis.suggested_status || 'Nuevo',
      source: existingLead?.source || source,
      notes: analysis.commercial_summary || '',
      next_action_at: this.toLocalDateTimeValue(analysis.suggested_follow_up_at),
    });
  }

  private reviewedAIAnalysisPayload(): ReviewedAIAnalysisPayload | undefined {
    const analysis = this.analysis();
    if (!analysis) {
      return undefined;
    }

    return {
      analysis_id: analysis.analysis_id,
      intent: analysis.intent || 'medium',
      detected_objections: analysis.detected_objections,
      commercial_summary: analysis.commercial_summary,
      suggested_next_action: analysis.suggested_next_action,
      source: this.analyzeForm.controls.source.value,
    };
  }


  getDatePart(value: string | null | undefined): string {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(0, 10);
  }

  getTimePart(value: string | null | undefined): string {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(11, 16);
  }

  updateNextActionAt(control: { setValue(value: string): void }, dateVal: string, timeVal: string): void {
    if (!dateVal || !timeVal) {
      control.setValue('');
    } else {
      control.setValue(`${dateVal}T${timeVal}`);
    }
  }

  private formatDate(value: string): string {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return 'Sin fecha';
    }

    return date.toLocaleString('es-CO', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }

  private toApiDateTime(value: string | undefined): string | undefined {
    if (!value) {
      return undefined;
    }

    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
  }

  private toLocalDateTimeValue(value: string | undefined): string {
    if (!value) {
      return '';
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return '';
    }

    const offsetMs = date.getTimezoneOffset() * 60_000;
    return new Date(date.getTime() - offsetMs).toISOString().slice(0, 16);
  }
}
