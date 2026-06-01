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
  Professional,
  CreateProfessionalPayload,
  UpdateProfessionalPayload,
  Appointment,
  CreateAppointmentPayload,
  UpdateAppointmentPayload,
  RescheduleAppointmentPayload,
  ConvertLeadPayload,
  AvailabilityResponse,
  ScheduleSummary,
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

  professionals() {
    return this.http.get<Professional[]>(`${this.baseUrl}/api/professionals`);
  }

  professional(id: string) {
    return this.http.get<Professional>(`${this.baseUrl}/api/professionals/${id}`);
  }

  createProfessional(payload: CreateProfessionalPayload) {
    return this.http.post<Professional>(`${this.baseUrl}/api/professionals`, payload);
  }

  updateProfessional(id: string, payload: UpdateProfessionalPayload) {
    return this.http.put<Professional>(`${this.baseUrl}/api/professionals/${id}`, payload);
  }

  appointments(params?: { date?: string; date_from?: string; date_to?: string; professional_id?: string; service_id?: string; status?: string }) {
    const httpParams: Record<string, string> = {};
    if (params) {
      Object.entries(params).forEach(([key, val]) => {
        if (val) {
          httpParams[key] = val;
        }
      });
    }
    return this.http.get<PaginatedResponse<Appointment>>(`${this.baseUrl}/api/appointments`, { params: httpParams });
  }

  appointment(id: string) {
    return this.http.get<Appointment>(`${this.baseUrl}/api/appointments/${id}`);
  }

  createAppointment(payload: CreateAppointmentPayload) {
    return this.http.post<Appointment>(`${this.baseUrl}/api/appointments`, payload);
  }

  updateAppointment(id: string, payload: UpdateAppointmentPayload) {
    return this.http.put<Appointment>(`${this.baseUrl}/api/appointments/${id}`, payload);
  }

  updateAppointmentStatus(id: string, status: string, adminNote?: string) {
    return this.http.post<Appointment>(`${this.baseUrl}/api/appointments/${id}/status`, { status, admin_note: adminNote });
  }

  rescheduleAppointment(id: string, payload: RescheduleAppointmentPayload) {
    return this.http.post<Appointment>(`${this.baseUrl}/api/appointments/${id}/reschedule`, payload);
  }

  convertLeadToAppointment(leadId: string, payload: ConvertLeadPayload) {
    return this.http.post<{ appointment_id: string; lead_status: string }>(`${this.baseUrl}/api/leads/${leadId}/convert-to-appointment`, payload);
  }

  availability(professionalId: string, dateFrom: string, dateTo: string, serviceId?: string, durationMinutes?: number) {
    const params: Record<string, string> = {
      professional_id: professionalId,
      date_from: dateFrom,
      date_to: dateTo,
    };
    if (serviceId) {
      params['service_id'] = serviceId;
    }
    if (durationMinutes) {
      params['duration_minutes'] = durationMinutes.toString();
    }
    return this.http.get<AvailabilityResponse>(`${this.baseUrl}/api/schedule/availability`, { params });
  }

  scheduleSummary(date?: string) {
    const params: Record<string, string> = {};
    if (date) {
      params['date'] = date;
    }
    return this.http.get<ScheduleSummary>(`${this.baseUrl}/api/dashboard/schedule-summary`, { params });
  }
}
