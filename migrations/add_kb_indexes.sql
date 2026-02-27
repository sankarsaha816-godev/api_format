-- Add indexes to speed up knowledge_bases queries
-- Run this migration after deploying the chatbot changes

-- Index for tenant_info_id lookups (most common filter)
ALTER TABLE knowledge_bases ADD INDEX idx_kb_tenant (tenant_info_id);

-- Index for embedding status checks (used in ensureKBEmbeddings)
ALTER TABLE knowledge_bases ADD INDEX idx_kb_tenant_embedding (tenant_info_id, embedding(10));

-- Index for FindBestMatch query
ALTER TABLE knowledge_bases ADD INDEX idx_kb_tenant_embedding_answer (tenant_info_id, embedding(10), answer(50));

-- Index for chat_logs to speed up inserts
ALTER TABLE chat_logs ADD INDEX idx_chat_log_tenant_email (tenant_info_id, email);
