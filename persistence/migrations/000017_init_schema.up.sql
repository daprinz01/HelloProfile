create table user_login(
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    user_id uuid,
    application_id uuid,
    login_time TIMESTAMPtz not NULL DEFAULT (now()),
    login_status BOOLEAN not null DEFAULT FALSE,
    response_code VARCHAR null,
    response_description VARCHAR null,
    device VARCHAR NULL,
    ip_address VARCHAR null,
    longitude DECIMAL NULL,
    latitude DECIMAL NULL,
	 FOREIGN KEY (user_id) REFERENCES users(id),
	 FOREIGN KEY (application_id) REFERENCES applications(id)
    
);

Create INDEX on user_login (user_id, application_id);