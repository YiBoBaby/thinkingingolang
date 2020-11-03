# 表结构
CREATE TABLE
IF
	NOT EXISTS `t_area_info` (
	`id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '自增主键',
	`area_id` BIGINT NOT NULL COMMENT '区域ID',
	`parent_id` BIGINT NOT NULL COMMENT '父区域ID',
	`area_name` VARCHAR ( 100 ) DEFAULT NULL COMMENT '区域名称',
	`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
	PRIMARY KEY ( `id` ) USING BTREE
	) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;
CREATE TABLE
IF
	NOT EXISTS `t_area_closure_info` (
	`id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '自增主键',
	`ancestor` BIGINT NOT NULL COMMENT '祖先id',
	`descendant` BIGINT NOT NULL COMMENT '后代id',
	`depth` TINYINT ( 4 ) NOT NULL DEFAULT '0' COMMENT '层级深度',
	PRIMARY KEY ( `id` ) USING BTREE
	) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;


# 插入
INSERT INTO t_area_closure_info ( ancestor, descendant, depth ) SELECT
t.ancestor,
2001,
t.depth + 1
FROM
	t_area_closure_info AS t
WHERE
	t.descendant = 1001 UNION ALL
SELECT
	2001,
	2001,
	1


# 查询后代节点
EXPLAIN SELECT
	descendant
FROM
	t_area_closure_info
WHERE
	ancestor = 1001
	AND depth <= 1;


# 查询祖先节点
EXPLAIN SELECT DISTINCT
	( ancestor )
FROM
	t_area_closure_info
WHERE
	descendant = 3001;


# 删除节点及其子节点数据
EXPLAIN DELETE t1
FROM
	`t_area_closure_info` t1,
	`t_area_closure_info` t2
WHERE
	t1.id = t2.id
	AND t2.ancestor = 11;