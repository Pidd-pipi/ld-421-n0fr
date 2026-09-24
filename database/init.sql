CREATE DATABASE IF NOT EXISTS `lab_equipment` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'lab_user'@'%' IDENTIFIED BY 'lab_password';
GRANT ALL PRIVILEGES ON `lab_equipment`.* TO 'lab_user'@'%';
FLUSH PRIVILEGES;
