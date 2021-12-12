CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "firstname" varchar NULL,
    "lastname" varchar NULL,
    "username" varchar NULL,
    "email" varchar NOT NULL,
    "phone" varchar NULL,
    "is_email_confirmed" boolean NOT NULL DEFAULT FALSE,
    "password" varchar NULL,
    "is_password_system_generated" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "is_locked_out" boolean NOT NULL DEFAULT FALSE,
    "image_url" varchar NULL,
    "is_active" boolean NOT NULL DEFAULT TRUE,
    CONSTRAINT "uc_users" UNIQUE ("id", "username", "email")
);

CREATE TABLE "languages" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    CONSTRAINT "uc_languages" UNIQUE ("id", "name")
);

CREATE TABLE "user_languages" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid,
    "language_id" uuid,
    proficiency varchar NULL,
    CONSTRAINT "uc_user_languages" UNIQUE ("id")
);

CREATE TABLE "timezones" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "zone" varchar NOT NULL,
    CONSTRAINT "uc_timezones" UNIQUE ("id", "name")
);

CREATE TABLE "user_timezones" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid,
    "timezone_id" uuid,
    CONSTRAINT "uc_user_timezones" UNIQUE ("id", "user_id")
);

CREATE TABLE "roles" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "description" varchar NOT NULL,
    CONSTRAINT "uc_roles" UNIQUE ("id", "name")
);

CREATE TABLE "user_roles" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid,
    "role_id" uuid,
    CONSTRAINT "uc_user_roles" UNIQUE ("id")
);

CREATE TABLE "identity_providers" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "client_id" varchar NOT NULL,
    "client_secret" varchar NOT NULL,
    "image_url" varchar NOT NULL,
    CONSTRAINT "uc_identity_providers" UNIQUE ("id", "name")
);

CREATE TABLE "user_providers" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid,
    "identity_provider_id" uuid,
    "identifier" varchar NOT NULL,
    CONSTRAINT "uc_user_providers" UNIQUE ("user_id", "identity_provider_id")
);

CREATE TABLE "countries" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "flag_image_url" varchar NULL,
    country_code varchar NULL,
    CONSTRAINT "uc_countries" UNIQUE ("id", "name", "flag_image_url")
);

CREATE TABLE "states" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "country_id" uuid,
    CONSTRAINT "uc_states" UNIQUE ("id", "name")
);

CREATE TABLE recents (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    profile_id uuid NOT NULL,
    title varchar NOT NULL,
    highlights varchar NOT NULL,
    year int NOT NULL,
    link varchar NOT NULL),
CONSTRAINT "uc_recents" UNIQUE (
    "profile_id", title
);

CREATE TABLE profiles (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    status boolean NOT NULL DEFAULT TRUE,
    profile_name varchar NOT NULL,
    fullname varchar NOT NULL,
    title varchar NOT NULL,
    bio varchar NOT NULL,
    company varchar NOT NULL,
    company_address varchar NOT NULL,
    image_url varchar NULL,
    phone varchar NOT NULL,
    email varchar NOT NULL,
    address_id uuid,
    website varchar NULL,
    is_default boolean NOT NULL DEFAULT FALSE,
    color int,
    CONSTRAINT "uc_profiles" UNIQUE (user_id, profile_name)
);

CREATE TABLE address_types (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    CONSTRAINT "uc_address_types" UNIQUE (name))
CREATE TABLE addresses (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid,
    street varchar NOT NULL,
    city varchar NOT NULL,
    state varchar NULL,
    country_id uuid,
    address_type uuid,
    CONSTRAINT "uc_addresses" UNIQUE (user_id, street)
);

CREATE TABLE contact_categories (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    CONSTRAINT "uc_address_types" UNIQUE (name)
);

CREATE TABLE contacts (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    profile_id uuid NOT NULL,
    contact_category_id uuid NOT NULL CONSTRAINT "uc_address_types" UNIQUE (name)
);

CREATE OR REPLACE VIEW user_details AS
SELECT
    b.id,
    b.firstname,
    b.lastname,
    b.email,
    b.phone,
    b.username,
    b."password",
    b.image_url AS profile_picture,
    b.is_email_confirmed,
    b.is_locked_out,
    b.is_password_system_generated,
    b.created_at,
    b.is_active,
    d."name" AS language_name,
    f."name" AS role_name,
    k."name" AS timezone_name,
    k.zone,
    m."name" AS provider_name,
    m.client_id,
    m.client_secret,
    m.image_url AS provider_logo
FROM
    users b
    LEFT JOIN languages d ON d.id = (
        SELECT
            language_id
        FROM
            user_languages e
        WHERE
            e.user_id = b.id)
    LEFT JOIN roles f ON f.id = (
        SELECT
            role_id
        FROM
            user_roles j
        WHERE
            j.user_id = b.id)
    LEFT JOIN timezones k ON k.id = (
        SELECT
            timezone_id
        FROM
            user_timezones l
        WHERE
            l.user_id = b.id)
    LEFT JOIN identity_providers m ON m.id = (
        SELECT
            identity_provider_id
        FROM
            user_providers n
        WHERE
            n.user_id = b.id);

-- CREATE or REPLACE VIEW user_details as
-- SELECT b.id, b.firstname, b.lastname, b.email, b.phone, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active,  k."name" as timezone_name, k.zone
-- from users b
-- left join timezones k on k.id = (select l.timezone_id from user_timezones l where l.user_id = b.id);
CREATE TABLE language_proficiency (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    proficiency varchar NULL UNIQUE
);

CREATE TABLE refresh_token (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 () UNIQUE,
    "user_id" uuid NOT NULL,
    "token" varchar NOT NULL UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE user_login (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid,
    login_time timestamptz NOT NULL DEFAULT (now()),
    login_status boolean NOT NULL DEFAULT FALSE,
    response_code varchar NULL,
    response_description varchar NULL,
    device varchar NULL,
    ip_address varchar NULL,
    longitude DECIMAL NULL,
    latitude DECIMAL NULL,
    resolved boolean NOT NULL DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE otp (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid,
    otp_token varchar NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    is_sms_preferred boolean NOT NULL DEFAULT FALSE,
    is_email_preferred boolean NOT NULL DEFAULT TRUE,
    purpose varchar NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE email_verification (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "email" varchar,
    "otp" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON email_verification (otp);

ALTER TABLE email_verification
    ADD CONSTRAINT uc_email_verification UNIQUE (otp);

CREATE INDEX ON otp (user_id, otp_token);

CREATE INDEX ON user_login (user_id);

CREATE INDEX ON profiles (user_id);

CREATE INDEX ON recents (profile_id);

ALTER TABLE "refresh_token"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles"
    ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "user_languages"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_languages"
    ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("id");

ALTER TABLE "user_providers"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_providers"
    ADD FOREIGN KEY ("identity_provider_id") REFERENCES "identity_providers" ("id");

ALTER TABLE "states"
    ADD FOREIGN KEY ("country_id") REFERENCES "countries" ("id");

ALTER TABLE "profiles"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "addresses"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "profiles"
    ADD FOREIGN KEY ("address_id") REFERENCES "addresses" ("id");

ALTER TABLE "addresses"
    ADD FOREIGN KEY ("country_id") REFERENCES "states" ("id");

ALTER TABLE "recents"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "refresh_token"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("contact_category_id") REFERENCES "contact_categories" ("id");

CREATE INDEX ON "countries" ("id");

CREATE INDEX ON "users" ("id", "username", "email");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "languages" ("id");

CREATE INDEX ON "user_languages" ("user_id", "language_id");

CREATE INDEX ON "timezones" ("id");

CREATE INDEX ON "user_timezones" ("user_id", "timezone_id");

CREATE INDEX ON "roles" ("id", "name");

CREATE INDEX ON "user_roles" ("user_id", "role_id");

CREATE INDEX ON "identity_providers" ("id", "name");

CREATE INDEX ON "user_providers" ("user_id", "identity_provider_id");

-- Add necessary constraints to users
ALTER TABLE Users
    ADD CONSTRAINT uc_user_id UNIQUE (Id);

ALTER TABLE Users
    ADD CONSTRAINT uc_user_username UNIQUE (username);

ALTER TABLE Users
    ADD CONSTRAINT uc_user_email UNIQUE (email);

-- Add necessary constraints to languages
ALTER TABLE languages
    ADD CONSTRAINT uc_languages_name UNIQUE (name);

-- Add necessary constraints to timezones
ALTER TABLE timezones
    ADD CONSTRAINT uc_timezones_name UNIQUE (name);

-- Add necessary constraints to user_timezones
ALTER TABLE user_timezones
    ADD CONSTRAINT uc_user_timezones_user_id UNIQUE (user_id);

-- Add necessary constraints to roles
ALTER TABLE roles
    ADD CONSTRAINT uc_roles_name UNIQUE (name);

-- Add necessary constraints to identity_providers
ALTER TABLE identity_providers
    ADD CONSTRAINT uc_identity_providers_name UNIQUE (name);

-- Add necessary constraints to countries
ALTER TABLE countries
    ADD CONSTRAINT uc_countries_flag_image_url UNIQUE (flag_image_url);

-- Add necessary constraints to states
ALTER TABLE states
    ADD CONSTRAINT uc_states_name UNIQUE (name);

ALTER TABLE language_proficiency
    ADD CONSTRAINT uc_language_proficiency UNIQUE (proficiency);

