ALTER TABLE notification
DROP CONSTRAINT IF EXISTS fk_auction,
DROP COLUMN IF EXISTS auction_id;
