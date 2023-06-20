SHOW DATABASES;

-- drop database children_math_redis;

CREATE DATABASE children_math_redis;

USE children_math_redis;

CREATE TABLE users(
	id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	user_name VARCHAR(20),
	`password` CHAR(60),
	points INT,
	CONSTRAINT UQ_users_name
	UNIQUE (user_name)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- DROP TABLE problems;

CREATE TABLE problems(
	id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	user_name VARCHAR(20),
	num1 INT,
	num2 INT,
	wrong_ans INT,
	operator CHAR(1),
	CONSTRAINT fk_problems_username
	FOREIGN KEY(user_name)
	REFERENCES users(user_name)
)ENGINE=INNODB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

ALTER TABLE problems
ADD INDEX idx_problems_username(user_name);

SHOW INDEX from users;