-- Alter vat_transactions table to match new structure
-- Step 1: Check if type column exists and rename it to transaction
SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'vat_transactions' 
  AND COLUMN_NAME = 'type');

SET @sql = IF(@column_exists > 0, 
  'ALTER TABLE `vat_transactions` CHANGE COLUMN `type` `transaction` VARCHAR(50) NOT NULL COMMENT ''sales, service, acquisition, overhead, cost, purchase''',
  'SELECT ''Column type does not exist, skipping rename''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 2: Add transaction column if it doesn't exist (in case we're running on a fresh table)
SET @column_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'vat_transactions' 
  AND COLUMN_NAME = 'transaction');

SET @sql = IF(@column_exists = 0, 
  'ALTER TABLE `vat_transactions` ADD COLUMN `transaction` VARCHAR(50) NOT NULL COMMENT ''sales, service, acquisition, overhead, cost, purchase'' AFTER `date`',
  'SELECT ''Column transaction already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 3: Drop old columns if they exist
ALTER TABLE `vat_transactions` 
  DROP COLUMN IF EXISTS `vehicle_id`,
  DROP COLUMN IF EXISTS `customer_id`,
  DROP COLUMN IF EXISTS `supplier_id`,
  DROP COLUMN IF EXISTS `reclaimable`;

-- Step 4: Add new columns if they don't exist
ALTER TABLE `vat_transactions`
  ADD COLUMN IF NOT EXISTS `purchase_price` DECIMAL(15,2) DEFAULT 0 AFTER `net_amount`,
  ADD COLUMN IF NOT EXISTS `name` VARCHAR(255) COMMENT 'regno for sales/service, overhead category for OH' AFTER `vat_amount`,
  ADD COLUMN IF NOT EXISTS `detail` VARCHAR(255) COMMENT 'stock num for sales, OH detail for OH' AFTER `name`;

-- Step 5: Update indexes
DROP INDEX IF EXISTS `idx_type` ON `vat_transactions`;

-- Add new index on transaction if it doesn't exist
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS 
  WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'vat_transactions' 
  AND INDEX_NAME = 'idx_transaction');

SET @sql = IF(@index_exists = 0,
  'CREATE INDEX `idx_transaction` ON `vat_transactions` (`transaction`)',
  'SELECT ''Index idx_transaction already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

DROP INDEX IF EXISTS `idx_vehicle_id` ON `vat_transactions`;

