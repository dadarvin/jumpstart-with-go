-- Database Query to create table
CREATE DATABASE IF NOT EXISTS `entry_task` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
USE `entry_task`;

-- Dumping structure for table entry_task.tbluser
CREATE TABLE IF NOT EXISTS `tbluser` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '0',
  `nickname` varchar(50) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
