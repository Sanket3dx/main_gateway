CREATE TABLE `user_logins` (
  `user_login_id` integer PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255),
  `password` varchar(255),
  `updated_on` timestamp,
  `updated_by` integer,
  `status` integer,
  `deleted` integer
);

CREATE TABLE `user_types` (
  `user_type_id` integer PRIMARY KEY AUTO_INCREMENT,
  `user_type` varchar(255),
  `updated_on` timestamp,
  `updated_by` integer,
  `deleted` integer
);

CREATE TABLE `Users` (
  `user_id` integer PRIMARY KEY AUTO_INCREMENT,
  `user_type_id` integer,
  `clinic_id` integer,
  `name` varchar(255),
  `gender` varchar(255),
  `dob` date
);
