-- jieqi_article_article definition

CREATE TABLE `jieqi_article_article` (
  `articleid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `siteid` smallint(6) unsigned NOT NULL DEFAULT '0',
  `postdate` int(11) unsigned NOT NULL DEFAULT '0',
  `lastupdate` int(11) unsigned NOT NULL DEFAULT '0',
  `articlename` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `keywords` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `initial` char(1) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `authorid` int(11) unsigned NOT NULL DEFAULT '0',
  `author` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `posterid` int(11) unsigned NOT NULL DEFAULT '0',
  `poster` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `agentid` int(11) unsigned NOT NULL DEFAULT '0',
  `agent` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `sortid` smallint(3) unsigned NOT NULL DEFAULT '0',
  `typeid` smallint(3) unsigned NOT NULL DEFAULT '0',
  `intro` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `notice` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `setting` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `lastvolumeid` int(11) unsigned NOT NULL DEFAULT '0',
  `lastvolume` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `lastchapterid` int(11) unsigned NOT NULL DEFAULT '0',
  `lastchapter` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `chapters` smallint(6) unsigned NOT NULL DEFAULT '0',
  `size` int(11) unsigned NOT NULL DEFAULT '0',
  `lastvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `dayvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `weekvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `monthvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `allvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `lastvote` int(11) unsigned NOT NULL DEFAULT '0',
  `dayvote` int(11) unsigned NOT NULL DEFAULT '0',
  `weekvote` int(11) unsigned NOT NULL DEFAULT '0',
  `monthvote` int(11) unsigned NOT NULL DEFAULT '0',
  `allvote` int(11) unsigned NOT NULL DEFAULT '0',
  `vipvotetime` int(11) NOT NULL DEFAULT '0',
  `vipvotenow` int(11) NOT NULL DEFAULT '0',
  `vipvotepreview` int(11) NOT NULL DEFAULT '0',
  `goodnum` int(11) unsigned NOT NULL DEFAULT '0',
  `badnum` int(11) unsigned NOT NULL DEFAULT '0',
  `toptime` int(11) unsigned NOT NULL DEFAULT '0',
  `saleprice` int(11) unsigned NOT NULL DEFAULT '0',
  `salenum` int(11) unsigned NOT NULL DEFAULT '0',
  `totalcost` int(11) unsigned NOT NULL DEFAULT '0',
  `articletype` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `permission` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `firstflag` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `fullflag` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `imgflag` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `power` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `display` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`articleid`),
  KEY `articlename` (`articlename`),
  KEY `posterid` (`posterid`),
  KEY `authorid` (`authorid`),
  KEY `agentid` (`agentid`),
  KEY `initial` (`initial`),
  KEY `sortid` (`sortid`,`typeid`),
  KEY `display` (`display`),
  KEY `size` (`size`),
  KEY `lastupdate` (`lastupdate`),
  KEY `author` (`author`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



-- jieqi_article_chapter definition

CREATE TABLE `jieqi_article_chapter` (
  `chapterid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `siteid` smallint(6) unsigned NOT NULL DEFAULT '0',
  `articleid` int(11) unsigned NOT NULL DEFAULT '0',
  `articlename` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `volumeid` int(11) unsigned NOT NULL DEFAULT '0',
  `posterid` int(11) unsigned NOT NULL DEFAULT '0',
  `poster` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `postdate` int(11) unsigned NOT NULL DEFAULT '0',
  `lastupdate` int(11) unsigned NOT NULL DEFAULT '0',
  `chaptername` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `chapterorder` smallint(6) unsigned NOT NULL DEFAULT '0',
  `size` int(11) unsigned NOT NULL DEFAULT '0',
  `saleprice` int(11) unsigned NOT NULL DEFAULT '0',
  `salenum` int(11) unsigned NOT NULL DEFAULT '0',
  `totalcost` int(11) unsigned NOT NULL DEFAULT '0',
  `attachment` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `isvip` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `chaptertype` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `power` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `display` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`chapterid`),
  KEY `articleid` (`articleid`),
  KEY `volumeid` (`volumeid`),
  KEY `chapterorder` (`chapterorder`),
  KEY `display` (`display`),
  KEY `articlename` (`articlename`,`chaptername`),
  KEY `lastupdate` (`lastupdate`)
) ENGINE=MyISAM AUTO_INCREMENT=310 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- 小说分类表

CREATE TABLE `sort` (
  `sortid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号',
  `weight` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `caption` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类名称',
  `shortname` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类简称',
  PRIMARY KEY (`sortid`),
  KEY `idx_weight` (`weight`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入初始分类数据
INSERT INTO `sort` (`sortid`, `shortname`, `caption`, `weight`) VALUES
(1, 'xuanhuan', '玄幻魔法', 1),
(2, 'wuxia',    '武侠修真', 2),
(3, 'dushi',    '都市言情', 3),
(4, 'lishi',    '历史军事', 4),
(5, 'kehuan',   '科幻灵异', 5),
(6, 'youxi',    '游戏竞技', 6),
(7, 'nvsheng',  '女生耽美', 7),
(8, 'qita',     '其他类型', 8);


-- 用户表
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `PASSWORD` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_login_time` datetime DEFAULT NULL,
  `current_login_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 书架表
CREATE TABLE `bookcase` (
  `caseid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `articleid` int(11) unsigned NOT NULL DEFAULT '0',
  `articlename` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `userid` int(11) unsigned NOT NULL DEFAULT '0',
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `chapterid` int(11) unsigned NOT NULL DEFAULT '0',
  `chaptername` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `chapterorder` smallint(6) unsigned NOT NULL DEFAULT '0',
  `joindate` int(11) unsigned NOT NULL DEFAULT '0',
  `lastvisit` int(11) unsigned NOT NULL DEFAULT '0',
  `flag` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`caseid`),
  KEY `articleid` (`articleid`),
  KEY `userid` (`userid`),
  KEY `flag` (`flag`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 书签表
CREATE TABLE `bookmark` (
  `bookid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `articleid` int(11) unsigned NOT NULL DEFAULT '0',
  `articlename` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `userid` int(11) unsigned NOT NULL DEFAULT '0',
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `chapterid` int(11) unsigned NOT NULL DEFAULT '0',
  `chaptername` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `chapterorder` smallint(6) unsigned NOT NULL DEFAULT '0',
  `joindate` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`bookid`),
  KEY `articleid` (`articleid`),
  KEY `userid` (`userid`),
  KEY `chapterid` (`chapterid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;