import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardAction, DashboardSummary } from '../../core/services/api.models';

interface ActionCard {
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
  count: number;
  tone: 'urgent' | 'today' | 'intent' | 'objection' | 'new';
}

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
  readonly loadingActions = signal(false);
  readonly actionError = signal<string | null>(null);

  readonly actionCards = computed<ActionCard[]>(() => {
    return this.actions().map((action) => {
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
      { label: 'Vencidos', count: counts['overdue_followup'] ?? 0, tone: 'urgent' },
      { label: 'Hoy', count: counts['today_followup'] ?? 0, tone: 'today' },
      { label: 'Alta intencion', count: counts['high_intent'] ?? 0, tone: 'intent' },
      { label: 'Objeciones', count: counts['detected_objection'] ?? 0, tone: 'objection' },
      { label: 'Nuevos', count: counts['new_lead'] ?? 0, tone: 'new' },
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
