-- Test data for VAT module (QEA AutoLens)
-- Use these to validate VAT calculations

-- Insert test tenant (if not exists)
INSERT IGNORE INTO `tenant_infos` (`tenant_info_id`, `tenant_name`, `company_registration`, `vat_number`, `created_at`)
VALUES (1, 'QEA AutoLens Demo', 'CR12345678', 'GB123456789', NOW());

-- Test 1: VAT-Qualifying Car Sale (20% VAT)
INSERT INTO `vehicle_sales_vat` 
(`tenant_info_id`, `vehicle_id`, `purchase_price`, `selling_price`, `margin`, `margin_vat`, `vat_scheme`, `sale_date`)
VALUES 
(1, 100, 15000.00, 18000.00, 3000.00, 0.00, 'vat_qualifying', '2024-01-15');

-- Test 2: Margin Scheme Car Sale (VAT on margin only)
INSERT INTO `vehicle_sales_vat` 
(`tenant_info_id`, `vehicle_id`, `purchase_price`, `selling_price`, `margin`, `margin_vat`, `vat_scheme`, `sale_date`)
VALUES 
(1, 101, 12000.00, 13500.00, 1500.00, 250.00, 'margin', '2024-01-20');

-- Test 3: Workshop Transaction (Labour + Parts)
INSERT INTO `workshop_transactions` 
(`tenant_info_id`, `date`, `labour_net`, `labour_vat`, `parts_net`, `parts_vat`, `zero_rated_net`)
VALUES 
(1, '2024-01-10', 500.00, 100.00, 800.00, 160.00, 0.00);

-- Test 4: Workshop Zero-Rated Work (e.g., repairs)
INSERT INTO `workshop_transactions` 
(`tenant_info_id`, `date`, `labour_net`, `labour_vat`, `parts_net`, `parts_vat`, `zero_rated_net`)
VALUES 
(1, '2024-01-12', 0.00, 0.00, 0.00, 0.00, 200.00);

-- Test 5: Purchase with VAT (Reclaimable)
INSERT INTO `vat_purchases` 
(`tenant_info_id`, `date`, `supplier_name`, `net_amount`, `vat_amount`, `gross_amount`, `reclaimable`, `category`, `invoice_number`)
VALUES 
(1, '2024-01-05', 'Car Parts Supplier Ltd', 1000.00, 200.00, 1200.00, TRUE, 'parts', 'INV-001');

-- Test 6: Purchase without VAT (Not Reclaimable)
INSERT INTO `vat_purchases` 
(`tenant_info_id`, `date`, `supplier_name`, `net_amount`, `vat_amount`, `gross_amount`, `reclaimable`, `category`, `invoice_number`)
VALUES 
(1, '2024-01-08', 'Travel Expense', 50.00, 0.00, 50.00, FALSE, 'other', 'TRAV-001');

-- Test 7: EU Acquisition (Goods from EU)
INSERT INTO `ec_transactions` 
(`tenant_info_id`, `date`, `type`, `net_amount`, `vat_due`, `customer_or_supplier_name`, `country_code`, `description`)
VALUES 
(1, '2024-01-25', 'acquisition', 5000.00, 1000.00, 'EU Supplier GmbH', 'DE', 'Car parts from Germany');

-- Test 8: EU Supply (Goods to EU)
INSERT INTO `ec_transactions` 
(`tenant_info_id`, `date`, `type`, `net_amount`, `vat_due`, `customer_or_supplier_name`, `country_code`, `description`)
VALUES 
(1, '2024-01-28', 'supply', 3000.00, 0.00, 'EU Customer SARL', 'FR', 'Vehicle sold to France');

-- Expected VAT Calculation Results for Q1 2024 (Jan 1 - Mar 31):
-- Box 1 (VAT due on sales):
--   - VAT-qualifying car: 18000 / 1.20 = 15000 net, VAT = 3000
--   - Workshop labour: 100
--   - Workshop parts: 160
--   - EU acquisition VAT: 1000
--   Total Box 1: 3000 + 100 + 160 + 1000 = 4260

-- Box 2 (VAT due on EC acquisitions): 1000

-- Box 3 (Total VAT due): 4260 + 1000 = 5260

-- Box 4 (VAT reclaimed): 200 + 1000 = 1200

-- Box 5 (Net VAT due): 5260 - 1200 = 4060

-- Box 6 (Total sales ex VAT):
--   - VAT-qualifying car: 15000
--   - Margin scheme car: 13500 (full price)
--   - Workshop labour net: 500
--   - Workshop parts net: 800
--   - Zero-rated: 200
--   Total Box 6: 15000 + 13500 + 500 + 800 + 200 = 30000

-- Box 7 (Total purchases ex VAT): 1000 + 50 = 1050

-- Box 8 (EC supplies ex VAT): 3000

-- Box 9 (EC acquisitions ex VAT): 5000
