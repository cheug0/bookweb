-- 长尾词插件数据表
-- 运行此SQL创建长尾词表

CREATE TABLE IF NOT EXISTS `article_langtail` (
  `langid` int unsigned NOT NULL AUTO_INCREMENT,
  `sourceid` int NOT NULL,
  `langname` varchar(50) NOT NULL DEFAULT '',
  `sourcename` varchar(50) NOT NULL DEFAULT '',
  `uptime` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`langid`),
  KEY `sourceid` (`sourceid`,`langid`),
  UNIQUE KEY `langname`(`langname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
