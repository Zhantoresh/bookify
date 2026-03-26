-- Migration 000004: Already handled in migration 000003
-- This migration can be skipped or left as no-op for backward compatibility

CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
