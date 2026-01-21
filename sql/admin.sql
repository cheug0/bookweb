-- 管理员表
CREATE TABLE IF NOT EXISTS `admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL COMMENT 'bcrypt加密',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认管理员 (用户名: admin, 密码: admin123)
-- 密码使用 bcrypt 加密，可使用 go run ./cmd/genpwd/ 生成新密码
INSERT INTO `admin` (`username`, `password`) VALUES ('admin', '$2a$10$RS.yPuMJr8Q9AxULUM1AmedQwJT17qxBcnxS12kzB/oUh9gyh2EKe');
