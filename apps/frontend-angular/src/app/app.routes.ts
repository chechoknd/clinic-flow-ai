import { Routes } from '@angular/router';

import { authGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';
import { AppLayout } from './core/layout/app-layout';

export const routes: Routes = [
  {
    path: 'login',
    canActivate: [guestGuard],
    loadComponent: () => import('./features/auth/login.page').then((m) => m.LoginPage),
  },
  {
    path: '',
    component: AppLayout,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard.page').then((m) => m.DashboardPage),
      },
      {
        path: 'leads',
        loadComponent: () => import('./features/leads/leads.page').then((m) => m.LeadsPage),
      },
      {
        path: 'followups',
        loadComponent: () =>
          import('./features/followups/followups.page').then((m) => m.FollowupsPage),
      },
      {
        path: 'services',
        loadComponent: () =>
          import('./features/services/services.page').then((m) => m.ServicesPage),
      },
      {
        path: 'clinic',
        loadComponent: () => import('./features/clinics/clinic.page').then((m) => m.ClinicPage),
      },
      {
        path: 'ai-assistant',
        loadComponent: () =>
          import('./features/ai-assistant/ai-assistant.page').then((m) => m.AiAssistantPage),
      },
      { path: '', pathMatch: 'full', redirectTo: 'dashboard' },
    ],
  },
  { path: '**', redirectTo: 'dashboard' },
];
