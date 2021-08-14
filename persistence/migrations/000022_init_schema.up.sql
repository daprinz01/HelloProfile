-- DROP VIEW user_details;

-- CREATE or REPLACE VIEW user_details as
-- SELECT b.id, b.firstname, b.lastname, b.email, b.phone, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active,  k."name" as timezone_name, k.zone
-- from users b
-- full join timezones k on k.id = (select l.timezone_id from user_timezones l where l.user_id = b.id);

-- ALTER TABLE user_languages add colume proficiency VARCHAR null;

-- Create TABLE language_proficiency(
--     id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
--     proficiency VARCHAR null UNIQUE
-- );

