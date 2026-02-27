-- Credit management tables for valuation APIs

CREATE TABLE IF NOT EXISTS `plan_tiers` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `monthly_credits` INT NOT NULL DEFAULT 0,
  `burst_limit_per_minute` INT NOT NULL DEFAULT 0,
  `allow_overage` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `tenant_subscriptions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `plan_tier_id` VARCHAR(50) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active',
  `period_start` DATETIME NOT NULL,
  `period_end` DATETIME NOT NULL,
  `credits_total` INT NOT NULL DEFAULT 0,
  `credits_used` INT NOT NULL DEFAULT 0,
  `credits_remaining` INT NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_subscriptions_tenant` (`tenant_info_id`),
  KEY `idx_tenant_subscriptions_status` (`status`),
  KEY `idx_tenant_subscriptions_plan` (`plan_tier_id`),
  CONSTRAINT `fk_tenant_subscriptions_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tenant_subscriptions_plan` FOREIGN KEY (`plan_tier_id`) REFERENCES `plan_tiers` (`id`) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS `credit_usage_ledger` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `user_id` BIGINT NULL,
  `endpoint` VARCHAR(255) NOT NULL,
  `credits_charged` INT NOT NULL DEFAULT 0,
  `request_id` VARCHAR(100) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'charged',
  `reason` VARCHAR(255) NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uniq_credit_usage_request` (`request_id`),
  KEY `idx_credit_usage_tenant` (`tenant_info_id`),
  KEY `idx_credit_usage_created` (`created_at`),
  CONSTRAINT `fk_credit_usage_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);

-- Seed tiers
INSERT INTO `plan_tiers` (`id`, `monthly_credits`, `burst_limit_per_minute`, `allow_overage`)
VALUES
  ('individual', 100, 5, 0),
  ('business', 2000, 50, 0),
  ('enterprise', 50000, 300, 1)
ON DUPLICATE KEY UPDATE
  `monthly_credits` = VALUES(`monthly_credits`),
  `burst_limit_per_minute` = VALUES(`burst_limit_per_minute`),
  `allow_overage` = VALUES(`allow_overage`);
