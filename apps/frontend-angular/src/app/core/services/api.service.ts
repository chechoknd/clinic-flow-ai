import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { map } from 'rxjs';

import { environment } from '../../../environments/environment';
import {
  AiReplyResponse,
  ClinicProfile,
  CreateLeadPayload,
  DashboardSummary,
  FollowUp,
  Lead,
  PaginatedResponse,
  DashboardSummaryApiResponse,
  ServicesApiResponse,
  UpdateLeadPayload,
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

  services() {
    return this.http
      .get<ServicesApiResponse>(`${this.baseUrl}/api/services`)
      .pipe(map((response) => ({ data: Array.isArray(response) ? response : response.data })));
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

  replySuggestion(payload: { patient_message: string; service_id?: string; lead_id?: string }) {
    return this.http.post<AiReplyResponse>(`${this.baseUrl}/api/ai/reply-suggestion`, payload);
  }

  objectionHandler(payload: { objection: string; service_id?: string; lead_id?: string }) {
    return this.http.post<AiReplyResponse>(`${this.baseUrl}/api/ai/objection-handler`, payload);
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
