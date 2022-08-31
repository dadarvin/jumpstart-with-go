-- Database Query to create table
CREATE DATABASE IF NOT EXISTS `entry_task` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
USE `entry_task`;

-- Dumping structure for table entry_task.tbluser
CREATE TABLE IF NOT EXISTS `user` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `username` varchar(50) UNIQUE NOT NULL,
    `nickname` varchar(50) NOT NULL,
    `password` varchar(255) NOT NULL,
    `profile_picture` varchar(255),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
