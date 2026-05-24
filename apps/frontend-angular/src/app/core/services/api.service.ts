import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { map } from 'rxjs';

import { environment } from '../../../environments/environment';
import {
  AiReplyResponse,
  AnalyzeConversationPayload,
  AnalyzeConversationResponse,
  ClinicProfile,
  ClinicServiceItem,
  CompleteFollowUpPayload,
  CreateLeadPayload,
  CreateServicePayload,
  DashboardSummary,
  FollowUp,
  FollowUpMessagePayload,
  FollowUpMessageResponse,
  Lead,
  PaginatedResponse,
  DashboardSummaryApiResponse,
  ServicesApiResponse,
  RescheduleFollowUpPayload,
  UpdateClinicPayload,
  UpdateLeadPayload,
  UpdateServicePayload,
} from './api.models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiBaseUrl;

  dashboardSummary() {
    return this.http
      .get<DashboardSummaryApiResponse>(`${this.baseUrl}/api/dashboard/summary`)
      .pipe(map((response) => this.mapDashboardSummary(response)));
  }

  clinicCurrent() {
    return this.http.get<ClinicProfile>(`${this.baseUrl}/api/clinics/current`);
  }

  updateClinicCurrent(payload: UpdateClinicPayload) {
    return this.http.put<ClinicProfile>(`${this.baseUrl}/api/clinics/current`, payload);
  }

  services() {
    return this.http
      .get<ServicesApiResponse>(`${this.baseUrl}/api/services`)
      .pipe(map((response) => ({ data: Array.isArray(response) ? response : response.data })));
  }

  createService(payload: CreateServicePayload) {
    return this.http.post<ClinicServiceItem>(`${this.baseUrl}/api/services`, payload);
  }

  updateService(id: string, payload: UpdateServicePayload) {
    return this.http.put<ClinicServiceItem>(`${this.baseUrl}/api/services/${id}`, payload);
  }

  deleteService(id: string) {
    return this.http.delete<void>(`${this.baseUrl}/api/services/${id}`);
  }

  leads() {
    return this.http.get<PaginatedResponse<Lead>>(`${this.baseUrl}/api/leads`);
  }

  createLead(payload: CreateLeadPayload) {
    return this.http.post<Lead>(`${this.baseUrl}/api/leads`, payload);
  }

  updateLead(id: string, payload: UpdateLeadPayload) {
    return this.http.put<{ id: string; status: string }>(`${this.baseUrl}/api/leads/${id}`, payload);
  }

  followups() {
    return this.http.get<PaginatedResponse<FollowUp>>(`${this.baseUrl}/api/followups`);
  }

  completeFollowUp(id: string, payload: CompleteFollowUpPayload) {
    return this.http.post<{ id: string; status: string }>(`${this.baseUrl}/api/followups/${id}/complete`, payload);
  }

  rescheduleFollowUp(id: string, payload: RescheduleFollowUpPayload) {
    return this.http.post<{ id: string; status: string }>(`${this.baseUrl}/api/followups/${id}/reschedule`, payload);
  }

  followUpMessage(payload: FollowUpMessagePayload) {
    return this.http.post<FollowUpMessageResponse>(`${this.baseUrl}/api/ai/follow-up-message`, payload);
  }

  replySuggestion(payload: { patient_message: string; service_id?: string; lead_id?: string }) {
    return this.http.post<AiReplyResponse>(`${this.baseUrl}/api/ai/reply-suggestion`, payload);
  }

  objectionHandler(payload: { objection: string; service_id?: string; lead_id?: string }) {
    return this.http.post<AiReplyResponse>(`${this.baseUrl}/api/ai/objection-handler`, payload);
  }

  analyzeConversation(payload: AnalyzeConversationPayload) {
    return this.http.post<AnalyzeConversationResponse>(`${this.baseUrl}/api/ai/analyze-conversation`, payload);
  }

  private mapDashboardSummary(response: DashboardSummaryApiResponse): DashboardSummary {
    return {
      total_leads: response.leads_total,
      pending_followups_today: response.pending_followups_today,
      overdue_followups: response.overdue_followups,
      upcoming_followups: response.upcoming_followups,
      conversion_rate: response.conversion_rate,
      status_counts: response.leads_by_status,
      top_services: response.top_services.map((service) => ({
        service_name: service.service_name,
        total: service.lead_count,
      })),
    };
  }
}
