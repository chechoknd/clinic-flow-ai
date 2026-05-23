import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { environment } from '../../../environments/environment';
import { AuthService } from './auth.service';
import { LoginResponse } from './auth.models';

const tokenKey = 'clinicflow_access_token';
const userKey = 'clinicflow_user';

const loginResponse: LoginResponse = {
  access_token: 'test-access-token',
  token_type: 'Bearer',
  expires_in: 3600,
  user: {
    id: 'user-1',
    full_name: 'Admin Demo',
    email: 'admin@sonrisaviva.demo',
    role: 'clinic_admin',
    clinic_id: 'clinic-1',
  },
};

describe('AuthService', () => {
  let service: AuthService;
  let http: HttpTestingController;

  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(AuthService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    http.verify();
    localStorage.clear();
  });

  it('starts without an authenticated session', () => {
    expect(service.token()).toBeNull();
    expect(service.user()).toBeNull();
    expect(service.isAuthenticated()).toBe(false);
  });

  it('posts credentials and persists the authenticated session', () => {
    service.login('admin@sonrisaviva.demo', 'clinicflow123').subscribe((response) => {
      expect(response).toEqual(loginResponse);
    });

    const request = http.expectOne(`${environment.apiBaseUrl}/api/auth/login`);
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({
      email: 'admin@sonrisaviva.demo',
      password: 'clinicflow123',
    });

    request.flush(loginResponse);

    expect(service.token()).toBe(loginResponse.access_token);
    expect(service.user()).toEqual(loginResponse.user);
    expect(service.isAuthenticated()).toBe(true);
    expect(localStorage.getItem(tokenKey)).toBe(loginResponse.access_token);
    expect(localStorage.getItem(userKey)).toBe(JSON.stringify(loginResponse.user));
  });

  it('clears the persisted session on logout', () => {
    localStorage.setItem(tokenKey, loginResponse.access_token);
    localStorage.setItem(userKey, JSON.stringify(loginResponse.user));

    service.logout();

    expect(service.token()).toBeNull();
    expect(service.user()).toBeNull();
    expect(service.isAuthenticated()).toBe(false);
    expect(localStorage.getItem(tokenKey)).toBeNull();
    expect(localStorage.getItem(userKey)).toBeNull();
  });
});
