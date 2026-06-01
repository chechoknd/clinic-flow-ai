import { Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardSummary, ScheduleSummary } from '../../core/services/api.models';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard.page.html',
  styleUrl: './dashboard.page.css',
  imports: [RouterLink],
})
export class DashboardPage {
  private readonly api = inject(ApiService);

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
        label: 'Citas de Hoy',
        value: sched?.todays_appointments ?? 0,
        desc: 'Agendadas para hoy',
        color: 'from-blue-500 to-indigo-600',
        route: '/schedule'
      },
      {
        label: 'Por Confirmar',
        value: sched?.appointments_pending_confirmation ?? 0,
        desc: 'Citas pendientes de confirmación',
        color: 'from-amber-500 to-orange-600',
        route: '/schedule'
      },
      {
        label: 'Leads sin Cita',
        value: sched?.hot_leads_without_appointment ?? 0,
        desc: 'Interesados listos para agendar',
        color: 'from-rose-500 to-pink-600',
        route: '/leads'
      },
      {
        label: 'Seguimientos Vencidos',
        value: sched?.overdue_followups ?? data?.overdue_followups ?? 0,
        desc: 'Vencidos o atrasados hoy',
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
}
