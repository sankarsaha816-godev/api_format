-- Migration: add supplier_id, deposit_customer_name, deposit_customer_id to stocks
ALTER TABLE `stocks`
  ADD COLUMN `supplier_id` INT DEFAULT NULL,
  ADD COLUMN `deposit_customer_name` VARCHAR(150) DEFAULT NULL,
  ADD COLUMN `deposit_customer_id` INT DEFAULT NULL;
