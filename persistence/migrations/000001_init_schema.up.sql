CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "firstname" varchar NULL,
    "lastname" varchar NULL,
    "username" varchar NULL,
    "email" varchar NOT NULL,
    "phone" varchar NULL,
    "country" VARCHAR NULL,
    "city" VARCHAR NULL,
    "is_email_confirmed" boolean NOT NULL DEFAULT FALSE,
    "password" varchar NULL,
    "is_password_system_generated" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "is_locked_out" boolean NOT NULL DEFAULT FALSE,
    "image_url" varchar NULL,
    "is_active" boolean NOT NULL DEFAULT TRUE,
    CONSTRAINT "uc_users" UNIQUE ("id", "username", "email")
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


CREATE TABLE profiles (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    user_id uuid NOT NULL,
    basic_block_id uuid null,
    contact_block_id uuid null,
    "status" boolean NOT NULL DEFAULT TRUE,
    profile_name varchar NOT NULL,
    page_color VARCHAR not Null,
    font VARCHAR not null,
    is_default boolean NOT NULL DEFAULT FALSE,
    CONSTRAINT "uc_profiles" UNIQUE (user_id, profile_name)
);

CREATE TABLE basic_block(
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    profile_photo_url varchar NULL,
    cover_photo_url VARCHAR null,
    cover_colour VARCHAR NULL,
    fullname varchar NOT NULL,
    title varchar NOT NULL,
    bio varchar NOT NULL
);

Create table contact_block(
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    phone varchar not NULL,
    email varchar not NULL,
    "address" varchar not null,
    website varchar not NULL
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

CREATE OR REPLACE VIEW user_details AS
SELECT
    b.id,
    b.firstname,
    b.lastname,
    b.email,
    b.phone,
    b.username,
    b."password",
    b.country,
    b.city,
    b.image_url AS profile_picture,
    b.is_email_confirmed,
    b.is_locked_out,
    b.is_password_system_generated,
    b.created_at,
    b.is_active,
    f."name" AS role_name
FROM
    users b
    LEFT JOIN roles f ON f.id = (
        SELECT
            role_id
        FROM
            user_roles j
        WHERE
            j.user_id = b.id);

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
    CONSTRAINT "uc_profile_contents" UNIQUE ("title", "profile_id")
);




CREATE INDEX ON email_verification (otp);

ALTER TABLE email_verification
    ADD CONSTRAINT uc_email_verification UNIQUE (otp);

CREATE INDEX ON otp (user_id, otp_token);

CREATE INDEX ON user_login (user_id);

CREATE INDEX ON profiles (user_id);


ALTER TABLE "refresh_token"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


ALTER TABLE "profiles"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");



ALTER TABLE "refresh_token"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "contacts"
    ADD FOREIGN KEY ("contact_category_id") REFERENCES "contact_categories" ("id");


CREATE INDEX ON "users" ("id", "username", "email");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");
CREATE INDEX ON "roles" ("id", "name");

CREATE INDEX ON "user_roles" ("user_id", "role_id");
ALTER TABLE "profile_socials"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "profile_socials"
    ADD FOREIGN KEY ("socials_id") REFERENCES "socials" ("id");
    
ALTER TABLE "profile_contents"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "profile_contents"
    ADD FOREIGN KEY ("call_to_action_id") REFERENCES "call_to_action" ("id");



-- Add necessary constraints to users
ALTER TABLE Users
    ADD CONSTRAINT uc_user_id UNIQUE (Id);

ALTER TABLE Users
    ADD CONSTRAINT uc_user_username UNIQUE (username);

ALTER TABLE Users
    ADD CONSTRAINT uc_user_email UNIQUE (email);

-- Add necessary constraints to roles
ALTER TABLE roles
    ADD CONSTRAINT uc_roles_name UNIQUE (name);








--- seed database with data
insert into roles(
name, description
)values
('admin', 'The super guys in the team, these guys have power to make or break the system. Think again before your push that button'),
('guest', 'These are vistiors to the application that need some form of authorisation to perform an action. It is also the default role of none is specified on account creation');

insert into content("type", "image_url")
VALUES
('links', 'https://helloprofile.io'),
('articles', 'https://helloprofile.io'),
('embedded videos', 'https://helloprofile.io'),
('embedded audios', 'https://helloprofile.io'),
('forms', 'https://helloprofile.io'),
('meetings', 'https://helloprofile.io'),
('events', 'https://helloprofile.io');

insert into "call_to_action"("type", "display_name")
VALUES
('learn_more', 'Learn More'),
('book_now', 'Book Now'),
('sign_up', 'Sign Up');
