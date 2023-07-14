SHOW DATABASES;

-- drop database children_math_redis;

CREATE DATABASE children_math_redis;

USE children_math_redis;

CREATE TABLE users(
	id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	user_name varchar(20),
	email VARCHAR(320),
	password CHAR(60),
	points INT,
	CONSTRAINT UQ_users_email
	UNIQUE (email)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- DROP TABLE problems;

CREATE TABLE problems(
	id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	user_id INT UNSIGNED,
	num1 INT,
	num2 INT,
	wrong_ans INT,
	operator CHAR(1),
	CONSTRAINT fk_problems_user_id
	FOREIGN KEY(user_id)
	REFERENCES users(id)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

ALTER TABLE problems
ADD INDEX idx_problems_user_id(user_id);

SHOW INDEX from users;