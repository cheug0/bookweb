-- 管理员表
CREATE TABLE IF NOT EXISTS `admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(32) NOT NULL COMMENT 'MD5加密',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认管理员 (用户名: admin, 密码: admin123)
-- 密码 MD5: 0192023a7bbd73250516f069df18b500
INSERT INTO `admin` (`username`, `password`) VALUES ('admin', '0192023a7bbd73250516f069df18b500');
