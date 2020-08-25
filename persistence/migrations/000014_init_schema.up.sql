-- Add necessary constraints to users
ALTER TABLE Users add CONSTRAINT uc_user_id UNIQUE (Id);
ALTER TABLE Users add CONSTRAINT uc_user_username UNIQUE (username);
ALTER TABLE Users add CONSTRAINT uc_user_email UNIQUE (email);
-- ALTER TABLE Users  Drop COLUMN created_at;
-- ALTER TABLE Users  ADD COLUMN created_at timestamptz not null DEFAULT (now());

-- Add necessary constraints to applications
ALTER TABLE applications add CONSTRAINT uc_application_name UNIQUE (name);

-- Add necessary constraints to languages
ALTER TABLE languages add CONSTRAINT uc_languages_name UNIQUE (name);

-- Add necessary constraints to timezones
ALTER TABLE timezones add CONSTRAINT uc_timezones_name UNIQUE (name);

-- Add necessary constraints to user_timezones
ALTER TABLE user_timezones add CONSTRAINT uc_user_timezones_user_id UNIQUE (user_id);

-- Add necessary constraints to roles
ALTER TABLE roles add CONSTRAINT uc_roles_name UNIQUE (name);

-- Add necessary constraints to identity_providers
ALTER TABLE identity_providers add CONSTRAINT uc_identity_providers_name UNIQUE (name);

-- Add necessary constraints to countries
ALTER TABLE countries add CONSTRAINT uc_countries_flag_image_url UNIQUE (flag_image_url);

-- Add necessary constraints to states
ALTER TABLE states add CONSTRAINT uc_states_name UNIQUE (name);
