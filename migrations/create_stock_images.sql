-- Create stock_images table
CREATE TABLE IF NOT EXISTS `stock_images` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `stock_id` BIGINT,
  `stock_no` VARCHAR(100) NOT NULL,
  `filename` VARCHAR(255) NOT NULL,
  `object_key` VARCHAR(500) NOT NULL UNIQUE,
  `public_url` VARCHAR(1000) NOT NULL,
  `file_size` BIGINT DEFAULT 0,
  `content_type` VARCHAR(100) DEFAULT 'image/jpeg',
  `uploaded_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_stock_id` (`stock_id`),
  KEY `idx_stock_no` (`stock_no`),
  CONSTRAINT `fk_stock_images_tenant_info_id` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
