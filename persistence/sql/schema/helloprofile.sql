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
    name varchar NOT NULL,
    CONSTRAINT "uc_languages" UNIQUE ("id", name)
);

CREATE TABLE "user_languages" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid NOT NULL,
    "language_id" uuid NOT NULL,
    "proficiency" varchar NULL,
    CONSTRAINT "uc_user_languages" UNIQUE ("id")
);

CREATE TABLE "timezones" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    "zone" varchar NOT NULL,
    CONSTRAINT "uc_timezones" UNIQUE ("id", name)
);

CREATE TABLE "user_timezones" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid NOT NULL,
    "timezone_id" uuid NOT NULL,
    CONSTRAINT "uc_user_timezones" UNIQUE ("id", "user_id")
);

CREATE TABLE "roles" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    "description" varchar NOT NULL,
    CONSTRAINT "uc_roles" UNIQUE ("id", name)
);

CREATE TABLE "user_roles" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid NOT NULL,
    "role_id" uuid NOT NULL,
    CONSTRAINT "uc_user_roles" UNIQUE ("id")
);

CREATE TABLE "identity_providers" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    "client_id" varchar NOT NULL,
    "client_secret" varchar NOT NULL,
    "image_url" varchar NOT NULL,
    CONSTRAINT "uc_identity_providers" UNIQUE ("id", name)
);

CREATE TABLE "user_providers" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid NOT NULL,
    "identity_provider_id" uuid NOT NULL,
    CONSTRAINT "uc_user_providers" UNIQUE ("id")
);

CREATE TABLE "countries" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    "flag_image_url" varchar NULL,
    "country_code" varchar NULL,
    CONSTRAINT "uc_countries" UNIQUE ("id", name, "flag_image_url")
);

CREATE TABLE "states" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    "country_id" uuid NOT NULL,
    CONSTRAINT "uc_states" UNIQUE ("id", name)
);

-- CREATE  VIEW user_details as
-- SELECT b.firstname, b.lastname, b.email, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active, d.name as language_name, f.name as role_name, k.name as timezone_name, k.zone, m.name as provider_name, m.client_id, m.client_secret, m.image_url as provider_logo
-- from users b
-- full join languages d on d.id = (select language_id from user_languages e where e.user_id = b.id)
-- full join roles f on f.id = (select role_id from user_roles j where j.user_id = b.id)
-- full join timezones k on k.id = (select timezone_id from user_timezones l where l.user_id = b.id)
-- full join identity_providers m on m.id = (select identity_provider_id from user_providers n where n.user_id = b.id);
CREATE TABLE user_details (
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
    "profile_picture" varchar NULL,
    "is_active" boolean NOT NULL DEFAULT TRUE,
    "timezone_name" varchar NULL,
    "zone" varchar NULL
);

CREATE TABLE refresh_token (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "user_id" uuid NOT NULL,
    "token" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE user_login (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    login_time timestamptz NOT NULL DEFAULT (now()),
    login_status boolean NOT NULL DEFAULT FALSE,
    response_code varchar NULL,
    response_description varchar NULL,
    device varchar NULL,
    ip_address varchar NULL,
    longitude DECIMAL NULL,
    latitude DECIMAL NULL,
    resolved boolean NOT NULL DEFAULT TRUE
);

CREATE TABLE otp (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    otp_token varchar NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    is_sms_preferred boolean NOT NULL DEFAULT FALSE,
    is_email_preferred boolean NOT NULL DEFAULT TRUE,
    purpose varchar NULL
);

CREATE TABLE language_proficiency (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    proficiency varchar NULL UNIQUE
);

CREATE TABLE email_verification (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "email" varchar,
    "otp" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE recents (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    profile_id uuid NOT NULL,
    title varchar NOT NULL,
    highlights varchar NOT NULL,
    year int NOT NULL,
    link varchar NOT NULL,
CONSTRAINT "uc_recents" UNIQUE (
    "profile_id", title
));

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

    
CREATE TABLE addresses (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid,
    street varchar NOT NULL,
    city varchar NOT NULL,
    state varchar NULL,
    country varchar null,
    isPrimaryAddress BOOLEAN DEFAULT FALSE,
    CONSTRAINT "uc_addresses" UNIQUE (user_id, street)
);

CREATE TABLE contact_categories (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    name varchar NOT NULL,
    CONSTRAINT "uc_contact_categories" UNIQUE (name)
);

CREATE TABLE contacts (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    profile_id uuid NOT NULL,
    contact_category_id uuid NOT NULL, CONSTRAINT "uc_contacts" UNIQUE (user_id, profile_id)
);

-- Add other components
CREATE TABLE "socials" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "name" varchar NOT NULL,
    "placeholder" varchar not null,
    "image_url" varchar NOT NULL,
    CONSTRAINT "uc_socials" UNIQUE ("name")
);

CREATE TABLE "profile_socials" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "username" varchar NOT NULL,
    "socials_id" uuid NOT NULL,
    "profile_id" uuid not null,
    "order" int not null default 0,
    CONSTRAINT "uc_supported_socials" UNIQUE ("socials_id", "profile_id")
);

CREATE TABLE "content" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "type" varchar NOT NULL,
    "image_url" varchar NOT NULL,
    CONSTRAINT "uc_content" UNIQUE ("type")
);



CREATE TABLE "call_to_action" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "type" varchar NOT NULL,
    "display_name" VARCHAR not null,
    CONSTRAINT "uc_call_to_action" UNIQUE ("type")
);



CREATE TABLE "profile_contents" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "title" varchar NOT NULL,
    "display_title" VARCHAR not null,
    "description" varchar NOT NULL,
    "url" varchar NOT NULL,
    "profile_id" uuid NOT NULL,
    "call_to_action_id" uuid not null,
    "content_id" uuid not null,
    "order" int not null default 0,
    CONSTRAINT "uc_link_content" UNIQUE ("title", "profile_id")
);
