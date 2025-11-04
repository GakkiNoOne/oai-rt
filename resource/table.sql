-- RT Token 管理表
CREATE TABLE `rt_rts` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `biz_id` varchar(255) NOT NULL COMMENT '业务ID（唯一标识）',
  `user_name` varchar(255) DEFAULT NULL COMMENT '用户名',
  `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
  `type` varchar(50) DEFAULT NULL COMMENT '账号类型（如：free, team）',
  `rt` text NOT NULL COMMENT 'Refresh Token',
  `at` text COMMENT 'Access Token',
  `proxy` varchar(255) DEFAULT NULL COMMENT '代理地址',
  `client_id` varchar(255) DEFAULT NULL COMMENT 'OpenAI Client ID',
  `tag` varchar(255) DEFAULT NULL COMMENT '标签',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用（1:启用, 0:禁用）',
  `last_rt` text COMMENT '上一次的 Refresh Token',
  `refresh_result` text COMMENT '刷新结果',
  `user_info` text COMMENT '用户信息（JSON）',
  `account_info` text COMMENT '账号信息（JSON）',
  `last_refresh_time` datetime DEFAULT NULL COMMENT '最后刷新时间',
  `memo` text COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_rt_rts_biz_id` (`biz_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='RT Token 管理表';

-- 系统配置表
CREATE TABLE `rt_system_configs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `config_key` varchar(255) NOT NULL COMMENT '配置键',
  `config_value` text COMMENT '配置值',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_rt_system_configs_config_key` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统配置表';

-- 如果表已存在但缺少唯一索引，执行以下语句：
-- ALTER TABLE rt_rts ADD UNIQUE INDEX `uni_rt_rts_biz_id` (`biz_id`);

