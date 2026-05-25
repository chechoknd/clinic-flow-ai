import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardSummary, FollowUp, Lead } from '../../core/services/api.models';

interface ActionCard {
  title: string;
  detail: string;
  meta: string;
  tone: 'urgent' | 'today' | 'new';
  primaryLabel: string;
  primaryPath: string;
  primaryQueryParams?: Record<string, string>;
  secondaryLabel: string;
  secondaryPath: string;
  secondaryQueryParams?: Record<string, string>;
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
  readonly leads = signal<Lead[]>([]);
  readonly followups = signal<FollowUp[]>([]);
  readonly loadingActions = signal(false);
  readonly actionError = signal<string | null>(null);

  readonly actionCards = computed<ActionCard[]>(() => {
    const overdue = this.followups()
      .filter((followup) => this.isOverdue(followup.next_action_at))
      .slice(0, 2)
      .map((followup) => this.followUpAction(followup, 'urgent'));

    const today = this.followups()
      .filter((followup) => this.isToday(followup.next_action_at))
      .filter((followup) => !this.isOverdue(followup.next_action_at))
      .slice(0, 2)
      .map((followup) => this.followUpAction(followup, 'today'));

    const newLeads = this.leads()
      .filter((lead) => lead.status === 'Nuevo')
      .slice(0, 2)
      .map((lead) => this.newLeadAction(lead));

    return [...overdue, ...today, ...newLeads].slice(0, 4);
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

    this.api.followups().subscribe({
      next: (response) => {
        this.followups.set(response.data);
        this.loadingActions.set(false);
      },
      error: () => {
        this.actionError.set('No fue posible cargar las acciones prioritarias.');
        this.loadingActions.set(false);
      },
    });

    this.api.leads().subscribe({
      next: (response) => this.leads.set(response.data),
      error: () => this.actionError.set('No fue posible cargar los leads recientes.'),
    });
  }

  private followUpAction(followup: FollowUp, tone: 'urgent' | 'today'): ActionCard {
    return {
      title: followup.full_name,
      detail: followup.service_name || 'Servicio por confirmar',
      meta: tone === 'urgent' ? `Vencido: ${this.formatDate(followup.next_action_at)}` : `Hoy: ${this.formatDate(followup.next_action_at)}`,
      tone,
      primaryLabel: 'Gestionar',
      primaryPath: '/followups',
      secondaryLabel: 'Responder con AI',
      secondaryPath: '/ai-assistant',
      secondaryQueryParams: this.contextQueryParams(followup),
    };
  }

  private newLeadAction(lead: Lead): ActionCard {
    return {
      title: lead.full_name,
      detail: lead.service_name || 'Sin servicio definido',
      meta: `Nuevo lead - ${lead.phone}`,
      tone: 'new',
      primaryLabel: 'Abrir lead',
      primaryPath: '/leads',
      secondaryLabel: 'Analizar conversacion',
      secondaryPath: '/inbox-ai',
      secondaryQueryParams: this.contextQueryParams(lead),
    };
  }

  private contextQueryParams(lead: Lead): Record<string, string> | undefined {
    const params: Record<string, string> = { lead_id: lead.id };
    if (lead.service_id) {
      params['service_id'] = lead.service_id;
    }
    return params;
  }

  private isOverdue(value: string | undefined): boolean {
    const date = this.parseDate(value);
    return Boolean(date && date.getTime() < Date.now());
  }

  private isToday(value: string | undefined): boolean {
    const date = this.parseDate(value);
    if (!date) {
      return false;
    }
    const now = new Date();
    return date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth() && date.getDate() === now.getDate();
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
