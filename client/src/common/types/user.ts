export type User = {
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  password: string;
  role: string;
  email_verified: boolean;
  registration_provider: string;
  created_at: Date;
  updated_at: Date;
};

export type CreateUserData = {
  firstname: string;
  lastname: string;
  username?: string;
  email: string;
  password: string;
};

export type LoginUserData = {
  email: string;
  password: string;
};

export type UpdateUserData = {
  firstname: string;
  lastname: string;
  username?: string;
  email: string;
};
