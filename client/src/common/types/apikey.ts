export type APIKey = {
  id: string;
  user_id: string;
  name: string;
  revoked: boolean;
  key: string;
  key_type: string;
  role: string;
  permissions: Array<string>;
  expires_at: Date;
  created_at: Date;
  updated_at: Date;
};

export type CreateAPIKeyData = {
  name: string;
  key_type?: string;
  role?: string;
  permissions?: Array<string>;
  expires_at: Date | string;
};
