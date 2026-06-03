export interface PaginatedResponse<T> {
  data: T[];
  pagination?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

export interface ClinicServiceItem {
  id: string;
  clinic_id?: string;
  name: string;
  description?: string;
  duration_minutes?: number;
  price_from?: number;
  benefits: string[];
  faq: unknown[];
  common_objections: string[];
  is_active?: boolean;
}

export interface CreateServicePayload {
  name: string;
  description?: string;
  duration_minutes?: number;
  price_from?: number;
  benefits: string[];
  faq: unknown[];
  common_objections: string[];
}

export interface UpdateServicePayload extends CreateServicePayload {
  is_active: boolean;
}

export type LeadStatus =
  | 'Nuevo'
  | 'Contactado'
  | 'Interesado'
  | 'Agendado'
  | 'No Respondio'
  | 'Perdido'
  | 'Convertido';

export interface Lead {
  id: string;
  full_name: string;
  phone: string;
  service_id?: string;
  service_name?: string;
  status: LeadStatus;
  source?: string;
  notes?: string;
  next_action_at?: string;
  created_at?: string;
}

export interface CreateLeadPayload {
  full_name: string;
  phone: string;
  service_id?: string;
  status: LeadStatus;
  source: string;
  notes?: string;
  next_action_at?: string;
}

export interface UpdateLeadPayload {
  status: LeadStatus;
  note?: string;
  next_action_at?: string;
}

export type FollowUp = Lead;

export interface CompleteFollowUpPayload {
  status?: LeadStatus;
  note?: string;
}

export interface RescheduleFollowUpPayload {
  next_action_at: string;
  note?: string;
}

export interface FollowUpMessagePayload {
  lead_id: string;
  service_id: string;
  last_contact_note?: string;
}

export interface FollowUpMessageResponse {
  generation_id: string;
  suggested_message: string;
  recommended_timing: string;
  next_step: string;
  safety_status: string;
}

export type CommunicationTone = 'amable' | 'profesional' | 'cercano' | 'juvenil' | 'elegante';

export interface ClinicProfile {
  id: string;
  name: string;
  clinic_type: string;
  city: string;
  phone?: string;
  whatsapp: string;
  address?: string;
  opening_hours?: Record<string, unknown>;
  general_faq?: unknown[];
  communication_tone: CommunicationTone;
}

export interface UpdateClinicPayload {
  name: string;
  city: string;
  phone?: string;
  whatsapp: string;
  address?: string;
  opening_hours: Record<string, unknown>;
  general_faq: unknown[];
  communication_tone: CommunicationTone;
}

export type ServicesApiResponse = ClinicServiceItem[] | { data: ClinicServiceItem[] };

export interface DashboardSummaryApiResponse {
  leads_total: number;
  pending_followups_today: number;
  overdue_followups: number;
  upcoming_followups: number;
  conversion_rate: number;
  leads_by_status: Record<string, number>;
  top_services: Array<{ service_id: string; service_name: string; lead_count: number }>;
}

export interface DashboardSummary {
  total_leads: number;
  pending_followups_today: number;
  overdue_followups: number;
  upcoming_followups: number;
  conversion_rate: number;
  status_counts: Record<string, number>;
  top_services: Array<{ service_name: string; total: number }>;
}

export interface AnalyzeConversationPayload {
  conversation_text: string;
  source?: string;
  lead_id?: string;
  service_id?: string;
}

export interface AnalyzeConversationResponse {
  analysis_id: string;
  detected_lead: {
    full_name: string;
    phone: string;
  };
  detected_service: {
    service_id: string;
    service_name: string;
    confidence: string;
  };
  intent: string;
  detected_objections: string[];
  suggested_status: LeadStatus;
  commercial_summary: string;
  suggested_reply: string;
  suggested_next_action: string;
  suggested_follow_up_at?: string;
  safety_status: string;
}

export interface AiReplyResponse {
  suggested_reply?: string;
  suggested_message?: string;
  message?: string;
  variants?: Record<string, string>;
}

export interface DayWorkingInterval {
  start: string;
  end: string;
}

export interface WorkingHoursConfig {
  monday?: DayWorkingInterval[];
  tuesday?: DayWorkingInterval[];
  wednesday?: DayWorkingInterval[];
  thursday?: DayWorkingInterval[];
  friday?: DayWorkingInterval[];
  saturday?: DayWorkingInterval[];
  sunday?: DayWorkingInterval[];
}

export interface Professional {
  id: string;
  clinic_id: string;
  full_name: string;
  role_or_specialty?: string;
  calendar_color?: string;
  working_hours: WorkingHoursConfig;
  is_active: boolean;
  service_ids: string[];
  created_at?: string;
  updated_at?: string;
}

export interface CreateProfessionalPayload {
  full_name: string;
  role_or_specialty?: string;
  calendar_color?: string;
  working_hours: WorkingHoursConfig;
  service_ids: string[];
}

export interface UpdateProfessionalPayload extends CreateProfessionalPayload {
  is_active: boolean;
}

export type AppointmentStatus =
  | 'scheduled'
  | 'confirmed'
  | 'pending_confirmation'
  | 'rescheduled'
  | 'no_show'
  | 'cancelled'
  | 'completed'
  | 'converted_from_lead';

export type ConfirmationStatus = 'pending' | 'confirmed' | 'not_required' | 'failed';

export interface ProfessionalSummary {
  id: string;
  full_name: string;
  calendar_color?: string;
}

export interface ServiceSummary {
  id: string;
  name: string;
}

export interface LeadSummary {
  id: string;
  full_name: string;
}

export interface Appointment {
  id: string;
  clinic_id: string;
  professional: ProfessionalSummary;
  service: ServiceSummary;
  lead?: LeadSummary;
  contact_name: string;
  contact_phone?: string;
  starts_at: string;
  ends_at: string;
  status: AppointmentStatus;
  confirmation_status: ConfirmationStatus;
  source: string;
  admin_notes?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CreateAppointmentPayload {
  professional_id: string;
  service_id: string;
  lead_id?: string;
  contact_name: string;
  contact_phone?: string;
  starts_at: string;
  ends_at?: string;
  duration_minutes?: number;
  status?: AppointmentStatus;
  source?: string;
  admin_notes?: string;
}

export interface UpdateAppointmentPayload {
  professional_id: string;
  service_id: string;
  contact_name: string;
  contact_phone?: string;
  starts_at: string;
  ends_at?: string;
  duration_minutes?: number;
  status: AppointmentStatus;
  confirmation_status: ConfirmationStatus;
  admin_notes?: string;
}

export interface RescheduleAppointmentPayload {
  starts_at: string;
  ends_at?: string;
  duration_minutes?: number;
  admin_note?: string;
}

export interface ConvertLeadPayload {
  professional_id: string;
  service_id: string;
  starts_at: string;
  ends_at?: string;
  duration_minutes?: number;
  status?: AppointmentStatus;
  admin_notes?: string;
  update_lead_status: boolean;
}

export interface TimeSlot {
  starts_at: string;
  ends_at: string;
}

export interface AvailabilityResponse {
  professional_id: string;
  slots: TimeSlot[];
}

export interface ProfessionalCount {
  professional_id: string;
  professional_name: string;
  appointment_count: number;
}

export interface ServiceCount {
  service_id: string;
  service_name: string;
  appointment_count: number;
}

export interface ScheduleSummary {
  date: string;
  todays_appointments: number;
  appointments_pending_confirmation: number;
  available_slots: number;
  hot_leads_without_appointment: number;
  overdue_followups: number;
  appointments_by_professional: ProfessionalCount[];
  leads_converted_to_appointments: number;
  top_services_by_schedule_demand: ServiceCount[];
}
