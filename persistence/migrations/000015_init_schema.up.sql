DROP VIEW user_details;

CREATE or REPLACE VIEW user_details as
SELECT b.id, b.firstname, b.lastname, b.email, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active, d."name" as language_name, f."name" as role_name, k."name" as timezone_name, k.zone, m."name" as provider_name, m.client_id, m.client_secret, m.image_url as provider_logo
from users b
full join languages d on d.id = (select language_id from user_languages e where e.user_id = b.id)
full join roles f on f.id = (select role_id from user_roles j where j.user_id = b.id)
full join timezones k on k.id = (select timezone_id from user_timezones l where l.user_id = b.id)
full join identity_providers m on m.id = (select identity_provider_id from user_providers n where n.user_id = b.id);