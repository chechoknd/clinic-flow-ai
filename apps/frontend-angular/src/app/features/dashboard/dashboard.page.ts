import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardAction, DashboardSummary } from '../../core/services/api.models';

interface ActionCard {
  leadID: string;
  title: string;
  detail: string;
  meta: string;
  tone: 'urgent' | 'today' | 'new' | 'intent' | 'objection';
  primaryLabel: string;
  primaryPath: string;
  primaryQueryParams?: Record<string, string>;
  secondaryLabel: string;
  secondaryPath: string;
  secondaryQueryParams?: Record<string, string>;
}

interface QueueSummaryItem {
  label: string;
  filter: ActionFilter;
  count: number;
  tone: 'urgent' | 'today' | 'intent' | 'objection' | 'new';
}

type ActionFilter = 'all' | 'overdue_followup' | 'today_followup' | 'high_intent' | 'detected_objection' | 'new_lead';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard.page.html',
  styleUrl: './dashboard.page.css',
  imports: [RouterLink],
})
export class DashboardPage {
  private readonly api = inject(ApiService);

  readonly summary = signal<DashboardSummary | null>(null);
  readonly actions = signal<DashboardAction[]>([]);
  readonly selectedActionFilter = signal<ActionFilter>('all');
  readonly noteDrafts = signal<Record<string, string>>({});
  readonly savedNoteLeads = signal<Record<string, boolean>>({});
  readonly noteSaving = signal<string | null>(null);
  readonly loadingActions = signal(false);
  readonly actionError = signal<string | null>(null);
  readonly actionNotice = signal<string | null>(null);

  readonly visibleActions = computed<DashboardAction[]>(() => {
    const filter = this.selectedActionFilter();
    return filter === 'all' ? this.actions() : this.actions().filter((action) => action.type === filter);
  });
  readonly actionCards = computed<ActionCard[]>(() => {
    return this.visibleActions().map((action) => {
      if (action.type === 'new_lead') {
        return this.newLeadAction(action);
      }
      if (action.type === 'high_intent' || action.type === 'detected_objection') {
        return this.insightAction(action);
      }
      return this.followUpAction(action);
    });
  });
  readonly queueSummary = computed<QueueSummaryItem[]>(() => {
    const counts = this.actions().reduce(
      (acc, action) => {
        acc[action.type] = (acc[action.type] ?? 0) + 1;
        return acc;
      },
      {} as Record<string, number>,
    );

    return [
      { label: 'Vencidos', filter: 'overdue_followup', count: counts['overdue_followup'] ?? 0, tone: 'urgent' },
      { label: 'Hoy', filter: 'today_followup', count: counts['today_followup'] ?? 0, tone: 'today' },
      { label: 'Alta intencion', filter: 'high_intent', count: counts['high_intent'] ?? 0, tone: 'intent' },
      { label: 'Objeciones', filter: 'detected_objection', count: counts['detected_objection'] ?? 0, tone: 'objection' },
      { label: 'Nuevos', filter: 'new_lead', count: counts['new_lead'] ?? 0, tone: 'new' },
    ];
  });

  constructor() {
    this.api.dashboardSummary().subscribe({ next: (summary) => this.summary.set(summary), error: () => {} });
    this.loadActionData();
  }

  metrics() {
    const data = this.summary();
    return [
      { label: 'Leads totales', value: data?.total_leads ?? 0 },
      { label: 'Para hoy', value: data?.pending_followups_today ?? 0 },
      { label: 'Vencidos', value: data?.overdue_followups ?? 0 },
      { label: 'Conversion', value: `${Math.round(data?.conversion_rate ?? 0)}%` },
    ];
  }

  statusRows() {
    const counts = this.summary()?.status_counts ?? {};
    const max = Math.max(...Object.values(counts), 1);
    return Object.entries(counts).map(([status, count]) => ({ status, count, width: (count / max) * 100 }));
  }

  selectActionFilter(filter: ActionFilter): void {
    this.selectedActionFilter.set(this.selectedActionFilter() === filter ? 'all' : filter);
  }

  actionNote(leadID: string): string {
    return this.noteDrafts()[leadID] ?? '';
  }

  setActionNote(leadID: string, value: string): void {
    this.noteDrafts.update((drafts) => ({ ...drafts, [leadID]: value }));
  }

  hasSavedActionNote(leadID: string): boolean {
    return Boolean(this.savedNoteLeads()[leadID]);
  }

  saveActionNote(action: DashboardAction): void {
    const note = this.actionNote(action.lead_id).trim();
    if (!note) {
      return;
    }

    this.noteSaving.set(action.lead_id);
    this.actionError.set(null);
    this.actionNotice.set(null);

    this.api
      .updateLead(action.lead_id, {
        status: action.status,
        note,
      })
      .subscribe({
        next: () => {
          this.noteDrafts.update((drafts) => ({ ...drafts, [action.lead_id]: '' }));
          this.savedNoteLeads.update((leads) => ({ ...leads, [action.lead_id]: true }));
          this.actionNotice.set('Nota guardada correctamente.');
          this.noteSaving.set(null);
        },
        error: () => {
          this.actionError.set('No fue posible guardar la nota rapida.');
          this.noteSaving.set(null);
        },
      });
  }

  private loadActionData(): void {
    this.loadingActions.set(true);
    this.actionError.set(null);

    this.api.dashboardActions(4).subscribe({
      next: (response) => {
        this.actions.set(response.data);
        this.loadingActions.set(false);
      },
      error: () => {
        this.actionError.set('No fue posible cargar las acciones prioritarias.');
        this.loadingActions.set(false);
      },
    });
  }

  private followUpAction(followup: DashboardAction): ActionCard {
    return {
      leadID: followup.lead_id,
      title: followup.full_name,
      detail: followup.service_name || 'Servicio por confirmar',
      meta:
        followup.type === 'overdue_followup'
          ? `Vencido: ${this.formatDate(followup.next_action_at)}`
          : `Hoy: ${this.formatDate(followup.next_action_at)}`,
      tone: followup.tone,
      primaryLabel: 'Gestionar',
      primaryPath: '/followups',
      secondaryLabel: 'Responder con AI',
      secondaryPath: '/ai-assistant',
      secondaryQueryParams: this.contextQueryParams(followup),
    };
  }

  private newLeadAction(lead: DashboardAction): ActionCard {
    return {
      leadID: lead.lead_id,
      title: lead.full_name,
      detail: lead.service_name || 'Sin servicio definido',
      meta: `Nuevo lead - ${lead.phone}`,
      tone: 'new',
      primaryLabel: 'Abrir lead',
      primaryPath: '/leads',
      primaryQueryParams: this.contextQueryParams(lead),
      secondaryLabel: 'Analizar conversacion',
      secondaryPath: '/inbox-ai',
      secondaryQueryParams: this.contextQueryParams(lead),
    };
  }

  private insightAction(action: DashboardAction): ActionCard {
    return {
      leadID: action.lead_id,
      title: action.full_name,
      detail: action.service_name || 'Servicio por confirmar',
      meta: `${action.reason} - ${action.phone}`,
      tone: action.tone,
      primaryLabel: 'Abrir lead',
      primaryPath: '/leads',
      primaryQueryParams: this.contextQueryParams(action),
      secondaryLabel: 'Responder con AI',
      secondaryPath: '/ai-assistant',
      secondaryQueryParams: this.contextQueryParams(action),
    };
  }

  private contextQueryParams(lead: DashboardAction): Record<string, string> | undefined {
    const params: Record<string, string> = { lead_id: lead.lead_id };
    if (lead.service_id) {
      params['service_id'] = lead.service_id;
    }
    return params;
  }

  private parseDate(value: string | undefined): Date | null {
    if (!value) {
      return null;
    }
    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? null : date;
  }

  private formatDate(value: string | undefined): string {
    const date = this.parseDate(value);
    if (!date) {
      return 'Sin fecha';
    }
    return date.toLocaleString('es-CO', {
      day: '2-digit',
      month: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }
}
