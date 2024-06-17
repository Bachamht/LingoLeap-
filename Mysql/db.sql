CREATE TABLE IF NOT EXISTS `user_info` (
    `user_id` INT NOT NULL AUTO_INCREMENT COMMENT '用户身份号',
    `user_phone` VARCHAR(50) NOT NULL COMMENT '用户手机号',
    `user_name` VARCHAR(50) NOT NULL COMMENT '用户姓名',
    `user_password` VARCHAR(50) NOT NULL COMMENT '用户密码',
    `user_mailbox` VARCHAR(50) NOT NULL COMMENT '用户邮箱',
    PRIMARY KEY (user_id)
)ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `session_info` (
    `session_id` INT NOT NULL AUTO_INCREMENT COMMENT '会话号',
    `user_id` VARCHAR(50) NOT NULL COMMENT '用户身份号',
    `role` VARCHAR(50) NOT NULL COMMENT '会话角色',
    `content` VARCHAR(200) NOT NULL  COMMENT '会话内容',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`session_id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;