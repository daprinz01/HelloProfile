-- -- Add necessary constraints to users
-- ALTER TABLE Users add UNIQUE (Id);
-- ALTER TABLE Users add UNIQUE (name);
-- ALTER TABLE Users add UNIQUE (username);
-- ALTER TABLE Users add UNIQUE (email);
-- ALTER TABLE Users  Drop COLUMN created_at;
-- ALTER TABLE Users  ADD COLUMN created_at timestamptz not null DEFAULT (now());

-- -- Add necessary constraints to applications
-- ALTER TABLE applications add UNIQUE (name);

-- -- Add necessary constraints to languages
-- ALTER TABLE languages add UNIQUE (name);

-- -- Add necessary constraints to timezones
-- ALTER TABLE timezones add UNIQUE (name);

-- -- Add necessary constraints to user_timezones
-- ALTER TABLE user_timezones add UNIQUE (user_id);

-- -- Add necessary constraints to roles
-- ALTER TABLE roles add UNIQUE (name);

-- -- Add necessary constraints to identity_providers
-- ALTER TABLE identity_providers add UNIQUE (name)

-- -- Add necessary constraints to countries
-- ALTER TABLE countries add UNIQUE (flag_image_url)

-- -- Add necessary constraints to states
-- ALTER TABLE states add UNIQUE (name)