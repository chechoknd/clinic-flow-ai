import { Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ApiService } from '../../core/services/api.service';
import { DashboardSummary } from '../../core/services/api.models';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard.page.html',
  styleUrl: './dashboard.page.css',
  imports: [RouterLink],
})
export class DashboardPage {
  private readonly api = inject(ApiService);

  readonly summary = signal<DashboardSummary | null>(null);

  constructor() {
    this.api.dashboardSummary().subscribe({ next: (summary) => this.summary.set(summary), error: () => {} });
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
}
