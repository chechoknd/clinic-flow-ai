export type UserRole = 'superadmin' | 'clinic_admin' | 'assistant';

export interface AuthUser {
  id: string;
  full_name: string;
  email: string;
  role: UserRole;
  clinic_id: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: 'Bearer';
  expires_in: number;
  user: AuthUser;
}
