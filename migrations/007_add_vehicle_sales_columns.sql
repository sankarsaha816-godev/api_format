-- Migration: Add fields required for margin VAT to vehicle_sales
-- Adds: vehicle_id, purchase_price, selling_price, margin, margin_vat, vat_scheme

ALTER TABLE vehicle_sales
  ADD COLUMN vehicle_id INT NULL,
  ADD COLUMN purchase_price DECIMAL(13,2) NULL,
  ADD COLUMN selling_price DECIMAL(13,2) NULL,
  ADD COLUMN margin DECIMAL(13,2) NULL,
  ADD COLUMN margin_vat DECIMAL(13,2) NULL,
  ADD COLUMN vat_scheme VARCHAR(50) NULL;

-- NOTE: margin and margin_vat are stored as plain columns to maintain compatibility.
-- If you prefer generated/computed columns, replace the above with GENERATED columns
-- or add triggers to keep them in sync.
