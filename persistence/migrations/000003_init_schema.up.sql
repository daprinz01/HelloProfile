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
    CONSTRAINT "uc_supported_socials" UNIQUE ("socials_id", "profile_id")
);

CREATE TABLE "content" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "type" varchar NOT NULL,
    "image_url" varchar NOT NULL,
    CONSTRAINT "uc_content" UNIQUE ("type")
);

insert into content("type", "image_url")
VALUES
('links', 'https://helloprofile.io'),
('articles', 'https://helloprofile.io'),
('embedded videos', 'https://helloprofile.io'),
('embedded audios', 'https://helloprofile.io'),
('forms', 'https://helloprofile.io'),
('meetings', 'https://helloprofile.io');


CREATE TABLE "call_to_action" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "type" varchar NOT NULL,
    "display_name" VARCHAR not null,
    CONSTRAINT "uc_call_to_action" UNIQUE ("type")
);

insert into "call_to_action"("type", "display_name")
VALUES
('learn_more', 'Learn More'),
('book_now', 'Book Now'),
('sign_up', 'Sign Up');

CREATE TABLE "link_content" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "title" varchar NOT NULL,
    "display_title" VARCHAR not null,
    "description" varchar NOT NULL,
    "url" varchar NOT NULL,
    "profile_id" uuid NOT NULL,
    "call_to_action_id" uuid not null,
    "order" int not null default 0,
    CONSTRAINT "uc_link_content" UNIQUE ("title", "profile_id")
);

CREATE TABLE "meeting_content" (
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
    "title" varchar NOT NULL,
    "display_title" VARCHAR not null,
    "description" varchar NOT NULL,
    "url" varchar NOT NULL,
    "profile_id" uuid NOT NULL,
    "call_to_action_id" uuid not null,
    "order" int not null default 0,
    CONSTRAINT "uc_meeting_content" UNIQUE ("title", "profile_id")
);

CREATE TABLE "embedded_content"(
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    "title" VARCHAR NOT NULL,
    "content_url" VARCHAR not NULL,
    "link_url" VARCHAR NOT NULL,
    "is_video" BOOLEAN NOT NULL DEFAULT TRUE,
    "profile_id" uuid NOT NULL,
    "call_to_action_id" uuid not null,
    "order" int not null default 0,
    CONSTRAINT "uc_embedded_video_content" UNIQUE("title")
);



ALTER TABLE "profile_socials"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "profile_socials"
    ADD FOREIGN KEY ("socials_id") REFERENCES "socials" ("id");
    
ALTER TABLE "link_content"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "link_content"
    ADD FOREIGN KEY ("call_to_action_id") REFERENCES "call_to_action" ("id");

ALTER TABLE "meeting_content"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "meeting_content"
    ADD FOREIGN KEY ("call_to_action_id") REFERENCES "call_to_action" ("id");

ALTER TABLE "embedded_content"
    ADD FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id");

ALTER TABLE "embedded_content"
    ADD FOREIGN KEY ("call_to_action_id") REFERENCES "call_to_action" ("id");