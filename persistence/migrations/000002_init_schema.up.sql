-- -- create table "applications"(
-- --   "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --   "name" varchar NOT NULL,
-- --   "description" varchar NOT NULL,
-- --   "created_at" timestamptz NOT NULL DEFAULT (now()),

-- --   CONSTRAINT "uc_applications" UNIQUE ("id", "name")
-- -- );

-- -- create table "users"(
-- --   "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --   "firstname" varchar  null,
-- --   "lastname" varchar  null,
-- --   "username" varchar  null,
-- --   "email" varchar not null,
-- --   "is_email_confirmed" BOOLEAN not null DEFAULT FALSE,
-- --   "password" varchar  null,
-- --   "is_password_system_generated" BOOLEAN not null DEFAULT FALSE,
-- --   "address" varchar  null,
-- --   "city" VARCHAR null,
-- --   "state" VARCHAR null,
-- --   "country" VARCHAR null,
-- --   "created_at" timestamptz NOT NULL DEFAULT (now()),
-- --   "is_locked_out" BOOLEAN not null default FALSE,
-- --   "image_url" varchar  null,
-- --   CONSTRAINT "uc_users" UNIQUE ("id", "username", "email")
-- -- );

-- -- create table "languages"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" varchar not null,
-- --     CONSTRAINT "uc_languages" UNIQUE ("id", "name")
-- -- );

-- -- create table "user_languages"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "user_id" uuid ,
-- --     "language_id" bigserial,
-- --     CONSTRAINT "uc_user_languages" UNIQUE ("id")
-- -- );

-- -- create table "timezones"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" varchar not null,
-- --     "zone" varchar not null,
-- --     CONSTRAINT "uc_timezones" UNIQUE ("id", "name")
-- -- );

-- -- create table "user_timezones"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "user_id" uuid,
-- --     "timezone_id" bigserial,
-- --     CONSTRAINT "uc_user_timezones" UNIQUE ("id", "user_id")
-- -- );

-- -- create table "roles"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" varchar not null,
-- --     "description" varchar not null,
-- --     CONSTRAINT "uc_roles" UNIQUE ("id", "name")
-- -- );

-- -- create table "user_roles"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "user_id" uuid,
-- --     "role_id" bigserial,
-- --     CONSTRAINT "uc_user_roles" UNIQUE ("id")
-- -- );

-- -- create table "applications_roles"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "applications_id" bigserial,
-- --     "roles_id" bigserial,
-- --     CONSTRAINT "uc_applications_roles" UNIQUE ("id")
-- -- );

-- -- create table "identity_providers"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" varchar not null,
-- --     "client_id" varchar not null,
-- --     "client_secret" varchar not null,
-- --     "image_url" varchar not null,
-- --     CONSTRAINT "uc_identity_providers" UNIQUE ("id", "name")
-- -- );

-- -- create table "user_providers"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "user_id" uuid,
-- --     "identity_provider_id" bigserial,
-- --     CONSTRAINT "uc_user_providers" UNIQUE ("id")
-- -- );
-- -- create TABLE "countries"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" VARCHAR not NULL,
-- --     "flag_image_url" VARCHAR null,

-- --     CONSTRAINT "uc_countries" UNIQUE ("id", "name", "flag_image_url")
-- -- );
-- -- create TABLE "states"(
-- --     "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
-- --     "name" VARCHAR not NULL,
-- --     "country_id" bigserial,
-- --     CONSTRAINT "uc_states" UNIQUE ("id", "name")
-- -- );

-- -- CREATE VIEW user_details as
-- -- SELECT b.firstname, b.lastname, b.email, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, d."name" as language_name, f."name" as role_name, k."name" as timezone_name, k.zone, m."name" as provider_name, m.client_id, m.client_secret, m.image_url
-- -- from users b
-- -- full join languages d on d.id = (select language_id from user_languages e where e.user_id = b.id)
-- -- full join roles f on f.id = (select role_id from user_roles j where j.user_id = b.id)
-- -- full join timezones k on k.id = (select timezone_id from user_timezones l where l.user_id = b.id)
-- -- full join identity_providers m on m.id = (select identity_provider_id from user_providers n where n.user_id = b.id);



-- alter table "applications_roles" add foreign key ("applications_id") references "applications" ("id");
-- alter table "applications_roles" add foreign key ("roles_id") references roles ("id");



-- alter table "user_roles" add foreign key ("user_id") references "users" ("id");
-- alter table "user_roles" add foreign key ("role_id") references "roles" ("id");

-- alter table "user_languages" add foreign key ("user_id") references "users" ("id");
-- alter table "user_languages" add foreign key ("language_id") references "languages" ("id");

-- alter table "user_providers" add foreign key ("user_id") references "users" ("id");
-- alter table "user_providers" add foreign key ("identity_provider_id") references "identity_providers" ("id");



-- ALTER TABLE "states" ADD FOREIGN key ("country_id") REFERENCES "countries" ("id");

-- CREATE INDEX on "countries" ("id");
-- CREATE INDEX on "applications" ("id","name");
-- CREATE INDEX on "users" ("id","username", "email");
-- CREATE INDEX on "languages" ("id");
-- CREATE INDEX on "user_languages" ("user_id", "language_id");
-- CREATE INDEX on "timezones" ("id");
-- CREATE INDEX on "user_timezones" ("user_id", "timezone_id");
-- CREATE INDEX on "roles" ("id", "name");
-- CREATE INDEX on "user_roles" ("user_id", "role_id");
-- CREATE INDEX on "applications_roles" ("applications_id", "roles_id");
-- CREATE INDEX on "identity_providers" ("id", "name");
-- CREATE INDEX on "user_providers" ("user_id", "identity_provider_id");

