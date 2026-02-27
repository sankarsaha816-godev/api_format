-- Add DEFAULT 'good' to feedback column in chat_logs table
ALTER TABLE chat_logs MODIFY feedback VARCHAR(255) DEFAULT 'good';
