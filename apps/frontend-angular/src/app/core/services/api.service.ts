import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { environment } from '../../../environments/environment';
import {
  AiReplyResponse,
  ClinicProfile,
  ClinicServiceItem,
  DashboardSummary,
  FollowUp,
  Lead,
  PaginatedResponse,
} from './api.models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiBaseUrl;

  dashboardSummary() {
    return this.http.get<DashboardSummary>(`${this.baseUrl}/api/dashboard/summary`);
  }

  clinicCurrent() {
    return this.http.get<ClinicProfile>(`${this.baseUrl}/api/clinics/current`);
  }

  services() {
    return this.http.get<{ data: ClinicServiceItem[] }>(`${this.baseUrl}/api/services`);
  }

  leads() {
    return this.http.get<PaginatedResponse<Lead>>(`${this.baseUrl}/api/leads`);
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
}
