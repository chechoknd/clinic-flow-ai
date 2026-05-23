import { Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { FollowUp } from '../../core/services/api.models';
import { ApiService } from '../../core/services/api.service';

@Component({
  selector: 'app-followups-page',
  imports: [RouterLink],
  templateUrl: './followups.page.html',
  styleUrl: './followups.page.css',
})
export class FollowupsPage {
  private readonly api = inject(ApiService);
  readonly followups = signal<FollowUp[]>([]);

  constructor() {
    this.api.followups().subscribe({ next: (response) => this.followups.set(response.data), error: () => {} });
  }
}
