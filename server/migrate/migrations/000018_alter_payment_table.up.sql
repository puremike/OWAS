ALTER TABLE payment
DROP CONSTRAINT IF EXISTS payment_auction_id_fkey;

ALTER TABLE payment
ADD CONSTRAINT payment_auction_id_fkey FOREIGN KEY (auction_id)
REFERENCES auctions(id) ON DELETE CASCADE;