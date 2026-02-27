-- Add subscription_id to tenant_subscriptions
ALTER TABLE `tenant_subscriptions`
  ADD COLUMN `subscription_id` VARCHAR(100) NULL AFTER `id`;

CREATE INDEX `idx_tenant_subscriptions_subscription_id` ON `tenant_subscriptions` (`subscription_id`);
