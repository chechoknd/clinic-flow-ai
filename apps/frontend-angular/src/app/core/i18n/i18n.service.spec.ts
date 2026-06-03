import { TestBed } from '@angular/core/testing';

import { I18nService } from './i18n.service';

describe('I18nService', () => {
  beforeEach(() => {
    localStorage.clear();
    TestBed.configureTestingModule({});
  });

  afterEach(() => {
    localStorage.clear();
  });

  it('defaults to Spanish', () => {
    const service = TestBed.inject(I18nService);

    expect(service.language()).toBe('es');
    expect(service.t('login.submit')).toBe('Entrar');
  });

  it('switches and persists English', () => {
    const service = TestBed.inject(I18nService);

    service.setLanguage('en');

    expect(service.language()).toBe('en');
    expect(service.t('login.submit')).toBe('Sign in');
    expect(localStorage.getItem('clinicflow_language')).toBe('en');
  });
});
