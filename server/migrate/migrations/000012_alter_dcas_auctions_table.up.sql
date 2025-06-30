-- (Run this in your PostgreSQL client, then use the name it returns in the DROP CONSTRAINT statement below)
-- SELECT conname
-- FROM pg_constraint
-- WHERE conrelid = 'auctions'::regclass
--   AND confrelid = 'users'::regclass
--   AND contype = 'f';

ALTER TABLE auctions
DROP CONSTRAINT IF EXISTS auctions_seller_id_fkey;

ALTER TABLE auctions
ADD CONSTRAINT auctions_seller_id_fkey FOREIGN KEY (seller_id)
REFERENCES users(id) ON DELETE CASCADE;