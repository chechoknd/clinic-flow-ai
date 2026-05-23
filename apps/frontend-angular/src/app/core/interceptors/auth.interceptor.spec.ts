import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { HttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { AuthService } from '../auth/auth.service';
import { authInterceptor } from './auth.interceptor';

describe('authInterceptor', () => {
  let http: HttpClient;
  let controller: HttpTestingController;
  let token: string | null;

  beforeEach(() => {
    token = null;

    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting(),
        {
          provide: AuthService,
          useValue: {
            token: () => token,
          },
        },
      ],
    });

    http = TestBed.inject(HttpClient);
    controller = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    controller.verify();
  });

  it('does not add an authorization header when there is no token', () => {
    http.get('/api/leads').subscribe();

    const request = controller.expectOne('/api/leads');

    expect(request.request.headers.has('Authorization')).toBe(false);
    request.flush({ data: [] });
  });

  it('adds the bearer token when a session token exists', () => {
    token = 'session-token';

    http.get('/api/leads').subscribe();

    const request = controller.expectOne('/api/leads');

    expect(request.request.headers.get('Authorization')).toBe('Bearer session-token');
    request.flush({ data: [] });
  });
});
