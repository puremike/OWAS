ALTER TABLE bid
DROP CONSTRAINT IF EXISTS bid_auction_id_fkey; 

ALTER TABLE bid
DROP CONSTRAINT IF EXISTS bid_bidder_id_fkey;

ALTER TABLE bid
ADD CONSTRAINT bid_auction_id_fkey FOREIGN KEY (auction_id)
REFERENCES auctions(id) ON DELETE CASCADE;

ALTER TABLE bid
ADD CONSTRAINT bid_bidder_id_fkey FOREIGN KEY (bidder_id)
REFERENCES users(id) ON DELETE CASCADE;