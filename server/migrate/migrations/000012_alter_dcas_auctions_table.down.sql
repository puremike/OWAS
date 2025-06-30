ALTER TABLE auctions
DROP CONSTRAINT IF EXISTS auctions_seller_id_fkey;

ALTER TABLE auctions
ADD CONSTRAINT auctions_seller_id_fkey FOREIGN KEY (seller_id)
REFERENCES users(id);