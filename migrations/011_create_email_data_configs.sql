-- Migration: Create email_data_configs table
CREATE TABLE IF NOT EXISTS email_data_configs (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    tenant_info_id INT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    access_token VARCHAR(255) NULL,
    client_id VARCHAR(255) NULL,
    client_secret VARCHAR(255) NULL,
    tenant_id VARCHAR(255) NULL,
    mailbox VARCHAR(255) NULL,
    max_results INT DEFAULT 25,
    query VARCHAR(255) NULL,
    since VARCHAR(100) NULL,
    process_attachments BOOLEAN DEFAULT TRUE,
    create_dynamic_table BOOLEAN DEFAULT FALSE,
    table_prefix VARCHAR(100) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_info_id (tenant_info_id)
);
