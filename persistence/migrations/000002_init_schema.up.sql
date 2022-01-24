-- seed the database
insert into roles(
name, description
)values
('admin', 'The super guys in the team, these guys have power to make or break the system. Think again before your push that button'),
('guest', 'These are vistiors to the application that need some form of authorisation to perform an action. It is also the default role of none is specified on account creation');