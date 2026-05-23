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
  name: string;
  description: string;
  duration_minutes: number;
  price_from: number;
  benefits: string[];
  faq: unknown[];
  common_objections: string[];
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

export interface FollowUp {
  id: string;
  lead_id: string;
  lead_name: string;
  phone: string;
  service_name?: string;
  status: LeadStatus;
  next_action_at: string;
  notes?: string;
}

export interface ClinicProfile {
  id: string;
  name: string;
  clinic_type: string;
  city: string;
  phone: string;
  whatsapp: string;
  address: string;
  communication_tone: string;
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

export interface AiReplyResponse {
  suggested_reply?: string;
  message?: string;
  variants?: Record<string, string>;
}
