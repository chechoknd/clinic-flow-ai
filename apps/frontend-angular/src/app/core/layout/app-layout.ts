import { Component, computed, inject, signal } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

import { AuthService } from '../auth/auth.service';

interface NavItem {
  label: string;
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

  readonly menuOpen = signal(false);
  readonly toggleOpen = (open: boolean) => !open;
  readonly userLabel = computed(() => this.auth.user()?.full_name ?? 'Equipo comercial');
  readonly navItems: NavItem[] = [
    { label: 'Dashboard', path: '/dashboard' },
    { label: 'Agenda', path: '/schedule' },
    { label: 'Profesionales', path: '/professionals' },
    { label: 'Leads', path: '/leads' },
    { label: 'Inbox AI', path: '/inbox-ai' },
    { label: 'Seguimientos', path: '/followups' },
    { label: 'Servicios', path: '/services' },
    { label: 'Asistente AI', path: '/ai-assistant' },
    { label: 'Clinica', path: '/clinic' },
  ];

  logout(): void {
    this.auth.logout();
    void this.router.navigate(['/login']);
  }
}
