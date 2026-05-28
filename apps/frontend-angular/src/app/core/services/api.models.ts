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
  price_from?: string;
  currency_code?: string;
  benefits: string[];
  faq: unknown[];
  common_objections: string[];
  is_active?: boolean;
}

export interface CreateServicePayload {
  name: string;
  description?: string;
  duration_minutes?: number;
  price_from?: string;
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

export type ContactOutcome =
  | 'attempted_no_answer'
  | 'asked_price'
  | 'interested'
  | 'scheduled'
  | 'lost_price'
  | 'lost_timing'
  | 'lost_trust'
  | 'converted'
  | 'follow_up_requested'
  | 'other';

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

export interface LeadDetail {
  id: string;
  full_name: string;
  phone: string;
  service?: {
    id: string;
    name: string;
  };
  status: LeadStatus;
  source: string;
  notes: LeadNote[];
  ai_insights?: LeadAIInsight[];
  next_action_at?: string;
  created_at: string;
}

export interface LeadNote {
  id: string;
  body: string;
  contact_outcome?: ContactOutcome;
  created_at: string;
}

export interface LeadAIInsight {
  id: string;
  analysis_id?: string;
  intent: 'low' | 'medium' | 'high' | string;
  detected_objections: string[];
  commercial_summary?: string;
  suggested_next_action?: string;
  source?: string;
  created_at: string;
}

export interface CreateLeadPayload {
  full_name: string;
  phone: string;
  service_id?: string;
  status: LeadStatus;
  source: string;
  notes?: string;
  next_action_at?: string;
  reviewed_ai_analysis?: ReviewedAIAnalysisPayload;
}

export interface UpdateLeadPayload {
  status: LeadStatus;
  note?: string;
  contact_outcome?: ContactOutcome;
  next_action_at?: string;
  clear_next_action_at?: boolean;
  reviewed_ai_analysis?: ReviewedAIAnalysisPayload;
}

export interface ReviewedAIAnalysisPayload {
  analysis_id?: string;
  intent: string;
  detected_objections?: string[];
  commercial_summary?: string;
  suggested_next_action?: string;
  source?: string;
}

export type FollowUp = Lead;

export interface CompleteFollowUpPayload {
  status?: LeadStatus;
  note?: string;
  contact_outcome?: ContactOutcome;
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
export type CountryCode = 'CO' | 'PE' | 'AR' | 'CL';
export type CurrencyCode = 'COP' | 'PEN' | 'ARS' | 'CLP';

export interface CurrencyMetadata {
  code: CurrencyCode;
  symbol: string;
  locale: string;
  decimal_digits: number;
  thousand_separator: string;
  decimal_separator: string;
  symbol_position: 'before' | 'after';
}

export interface ClinicProfile {
  id: string;
  name: string;
  clinic_type: string;
  city: string;
  country_code: CountryCode;
  currency_code: CurrencyCode;
  currency: CurrencyMetadata;
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
  country_code: CountryCode;
  currency_code: CurrencyCode;
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

export type DashboardActionType = 'overdue_followup' | 'today_followup' | 'new_lead' | 'high_intent' | 'detected_objection';
export type DashboardActionTone = 'urgent' | 'today' | 'new' | 'intent' | 'objection';

export interface DashboardAction {
  type: DashboardActionType;
  tone: DashboardActionTone;
  priority: number;
  lead_id: string;
  full_name: string;
  phone: string;
  service_id?: string;
  service_name?: string;
  status: LeadStatus;
  source: string;
  reason: string;
  next_action_at?: string;
  created_at: string;
}

export interface DashboardActionsResponse {
  data: DashboardAction[];
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
