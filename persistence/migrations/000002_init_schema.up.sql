create table saved_profiles(
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    profile_id uuid not null,
    first_name varchar not null,
    last_name varchar not null,
    email varchar not null,
    is_added boolean not null default false,
    is_terms_agreed boolean not null default true,
    CONSTRAINT "uc_saved_profiles" UNIQUE ("email", "profile_id")
);

ALTER TABLE "saved_profiles"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");