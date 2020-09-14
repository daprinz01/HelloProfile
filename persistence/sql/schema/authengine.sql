

create table "applications" (
  "id" bigserial PRIMARY KEY,
  name varchar NOT NULL,
  "description" varchar NOT NULL,
  "icon_url" varchar null,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "uc_applications" UNIQUE ("id", name)
);

create table "users" (
  "id" bigserial primary key,
  "firstname" varchar  null,
  "lastname" varchar  null,
  "username" varchar  null,
  "email" varchar not null,
  "phone" varchar null,
  "is_email_confirmed" BOOLEAN not null DEFAULT FALSE,
  "password" varchar  null,
  "is_password_system_generated" BOOLEAN not null DEFAULT FALSE,
  "address" varchar  null,
  "city" VARCHAR null,
  "state" VARCHAR null,
  "country" VARCHAR null,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_locked_out" BOOLEAN not null default FALSE,
  "image_url" varchar  null,
  "is_active" BOOLEAN not null DEFAULT TRUE,
  CONSTRAINT "uc_users" UNIQUE ("id", "username", "email")
);

create table "languages" (
    "id" bigserial primary key,
    name varchar not null,
    CONSTRAINT "uc_languages" UNIQUE ("id", name)
);

create table "user_languages" (
    "id" bigserial primary key,
    "user_id" bigserial ,
    "language_id" bigserial,
    "proficiency" varchar null,
    CONSTRAINT "uc_user_languages" UNIQUE ("id")
);

create table "timezones" (
    "id" bigserial primary key,
    name varchar not null,
    "zone" varchar not null,
    CONSTRAINT "uc_timezones" UNIQUE ("id", name)
);

create table "user_timezones" (
    "id" bigserial primary key,
    "user_id" bigserial,
    "timezone_id" bigserial,
    CONSTRAINT "uc_user_timezones" UNIQUE ("id", "user_id")
);

create table "roles" (
    "id" bigserial primary key,
    name varchar not null,
    "description" varchar not null,
    CONSTRAINT "uc_roles" UNIQUE ("id", name)
);

create table "user_roles" (
    "id" bigserial primary key,
    "user_id" bigserial,
    "role_id" bigserial,
    CONSTRAINT "uc_user_roles" UNIQUE ("id")
);

create table "applications_roles" (
    "id" bigserial primary key,
    "applications_id" bigserial,
    "roles_id" bigserial,
    CONSTRAINT "uc_applications_roles" UNIQUE ("id")
);

create table "identity_providers" (
    "id" bigserial primary key,
    name varchar not null,
    "client_id" varchar not null,
    "client_secret" varchar not null,
    "image_url" varchar not null,
    CONSTRAINT "uc_identity_providers" UNIQUE ("id", name)
);

create table "user_providers" (
    "id" bigserial primary key,
    "user_id" bigserial,
    "identity_provider_id" bigserial,
    CONSTRAINT "uc_user_providers" UNIQUE ("id")
);
create TABLE "countries" (
    "id" bigserial PRIMARY KEY,
    name VARCHAR not NULL,
    "flag_image_url" VARCHAR null,
"country_code" varchar null,
    CONSTRAINT "uc_countries" UNIQUE ("id", name, "flag_image_url")
);
create TABLE "states" (
    "id" bigserial PRIMARY KEY,
    name VARCHAR not NULL,
    "country_id" bigserial,
    CONSTRAINT "uc_states" UNIQUE ("id", name)
);


-- CREATE  VIEW user_details as
-- SELECT b.firstname, b.lastname, b.email, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active, d.name as language_name, f.name as role_name, k.name as timezone_name, k.zone, m.name as provider_name, m.client_id, m.client_secret, m.image_url as provider_logo
-- from users b
-- full join languages d on d.id = (select language_id from user_languages e where e.user_id = b.id)
-- full join roles f on f.id = (select role_id from user_roles j where j.user_id = b.id)
-- full join timezones k on k.id = (select timezone_id from user_timezones l where l.user_id = b.id)
-- full join identity_providers m on m.id = (select identity_provider_id from user_providers n where n.user_id = b.id);

create table  user_details(
    "id" bigserial not null,
  "firstname" varchar  null,
  "lastname" varchar  null,
  "username" varchar  null,
  "email" varchar not null,
  "phone" varchar null,
  "is_email_confirmed" BOOLEAN not null DEFAULT FALSE,
  "password" varchar  null,
  "is_password_system_generated" BOOLEAN not null DEFAULT FALSE,
  "address" varchar  null,
  "city" VARCHAR null,
  "state" VARCHAR null,
  "country" VARCHAR null,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_locked_out" BOOLEAN not null default FALSE,
  "profile_picture" varchar  null,
  "is_active" BOOLEAN not null DEFAULT TRUE,
  "timezone_name" varchar null,
  "zone" varchar null
);

create table refresh_token(
    "id" bigserial primary key,
    "user_id" bigserial not null,
    "token" varchar not null,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

create table user_login(
    id bigserial primary KEY,
    user_id bigserial,
    application_id bigserial,
    login_time TIMESTAMPtz not NULL DEFAULT (now()),
    login_status BOOLEAN not null DEFAULT FALSE,
    response_code VARCHAR null,
    response_description VARCHAR null,
    device VARCHAR NULL,
    ip_address VARCHAR null,
    longitude DECIMAL NULL,
    latitude DECIMAL NULL,
    resolved BOOLEAN not null default TRUE
    
);

create table otp(
    id bigserial primary key,
    user_id bigserial,
    otp_token varchar null,
    created_at TIMESTAMPtz not NULL DEFAULT (now()),
    is_sms_preferred boolean not null default FALSE,
    is_email_preferred boolean not null default TRUE,
    purpose varchar null
);


Create TABLE language_proficiency(
    id bigserial PRIMARY KEY,
    proficiency VARCHAR null UNIQUE
);
