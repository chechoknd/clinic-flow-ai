import { Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardSummary, ScheduleSummary } from '../../core/services/api.models';
import { I18nService } from '../../core/i18n/i18n.service';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard.page.html',
  styleUrl: './dashboard.page.css',
  imports: [RouterLink],
})
export class DashboardPage {
  private readonly api = inject(ApiService);
  readonly i18n = inject(I18nService);

  readonly summary = signal<DashboardSummary | null>(null);
  readonly schedSummary = signal<ScheduleSummary | null>(null);
  readonly activeTab = signal<'schedule' | 'crm'>('schedule');

  constructor() {
    this.api.dashboardSummary().subscribe({
      next: (summary) => this.summary.set(summary),
      error: () => {}
    });
    this.api.scheduleSummary().subscribe({
      next: (summary) => this.schedSummary.set(summary),
      error: () => {}
    });
  }

  changeTab(tab: 'schedule' | 'crm'): void {
    this.activeTab.set(tab);
  }

  metrics() {
    const data = this.summary();
    const sched = this.schedSummary();
    return [
      {
        label: this.i18n.t('dashboard.metric.todayAppointments'),
        value: sched?.todays_appointments ?? 0,
        desc: this.i18n.t('dashboard.metric.todayAppointmentsDesc'),
        color: 'from-blue-500 to-indigo-600',
        route: '/schedule'
      },
      {
        label: this.i18n.t('dashboard.metric.pendingConfirmation'),
        value: sched?.appointments_pending_confirmation ?? 0,
        desc: this.i18n.t('dashboard.metric.pendingConfirmationDesc'),
        color: 'from-amber-500 to-orange-600',
        route: '/schedule'
      },
      {
        label: this.i18n.t('dashboard.metric.hotLeads'),
        value: sched?.hot_leads_without_appointment ?? 0,
        desc: this.i18n.t('dashboard.metric.hotLeadsDesc'),
        color: 'from-rose-500 to-pink-600',
        route: '/leads'
      },
      {
        label: this.i18n.t('dashboard.metric.overdueFollowups'),
        value: sched?.overdue_followups ?? data?.overdue_followups ?? 0,
        desc: this.i18n.t('dashboard.metric.overdueFollowupsDesc'),
        color: 'from-purple-500 to-violet-600',
        route: '/leads'
      },
    ];
  }

  statusRows() {
    const counts = this.summary()?.status_counts ?? {};
    const max = Math.max(...Object.values(counts), 1);
    return Object.entries(counts).map(([status, count]) => ({
      status,
      count,
      width: (count / max) * 100
    }));
  }

  professionalRows() {
    const profs = this.schedSummary()?.appointments_by_professional ?? [];
    const max = Math.max(...profs.map((p) => p.appointment_count), 1);
    return profs.map((p) => ({
      name: p.professional_name,
      count: p.appointment_count,
      width: (p.appointment_count / max) * 100
    }));
  }

  topServicesByDemand() {
    return this.schedSummary()?.top_services_by_schedule_demand ?? [];
  }

  statusLabel(status: string): string {
    return this.i18n.t(`lead.status.${status}`);
  }

  interpolate(key: string, count: number): string {
    return this.i18n.t(key).replace('{{count}}', String(count));
  }

  appointmentCountLabel(count: number): string {
    return count === 1 ? this.i18n.t('dashboard.appointmentsSingular') : this.i18n.t('dashboard.appointmentsPlural');
  }
}
