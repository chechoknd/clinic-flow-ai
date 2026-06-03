import { Component, computed, inject, signal } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

import { AuthService } from '../auth/auth.service';
import { AppLanguage, I18nService } from '../i18n/i18n.service';

interface NavItem {
  labelKey:
    | 'nav.dashboard'
    | 'nav.schedule'
    | 'nav.professionals'
    | 'nav.leads'
    | 'nav.inboxAi'
    | 'nav.followups'
    | 'nav.services'
    | 'nav.aiAssistant'
    | 'nav.clinic';
  path: string;
}

@Component({
  selector: 'app-layout',
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './app-layout.html',
  styleUrl: './app-layout.css',
})
export class AppLayout {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  readonly i18n = inject(I18nService);

  readonly menuOpen = signal(false);
  readonly toggleOpen = (open: boolean) => !open;
  readonly userLabel = computed(() => this.auth.user()?.full_name ?? this.i18n.t('common.commercialTeam'));
  readonly navItems: NavItem[] = [
    { labelKey: 'nav.dashboard', path: '/dashboard' },
    { labelKey: 'nav.schedule', path: '/schedule' },
    { labelKey: 'nav.professionals', path: '/professionals' },
    { labelKey: 'nav.leads', path: '/leads' },
    { labelKey: 'nav.inboxAi', path: '/inbox-ai' },
    { labelKey: 'nav.followups', path: '/followups' },
    { labelKey: 'nav.services', path: '/services' },
    { labelKey: 'nav.aiAssistant', path: '/ai-assistant' },
    { labelKey: 'nav.clinic', path: '/clinic' },
  ];

  setLanguage(language: AppLanguage): void {
    this.i18n.setLanguage(language);
  }

  logout(): void {
    this.auth.logout();
    void this.router.navigate(['/login']);
  }
}
