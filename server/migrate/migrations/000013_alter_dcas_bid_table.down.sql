-- Down Migration for bid table

-- Step 1: Drop the foreign key constraints that have ON DELETE CASCADE.
ALTER TABLE bid
DROP CONSTRAINT IF EXISTS bid_auction_id_fkey;

ALTER TABLE bid
DROP CONSTRAINT IF EXISTS bid_bidder_id_fkey;

-- Step 2: Add the foreign key constraints back WITHOUT ON DELETE CASCADE.
ALTER TABLE bid
ADD CONSTRAINT bid_auction_id_fkey FOREIGN KEY (auction_id)
REFERENCES auctions(id);

ALTER TABLE bid
ADD CONSTRAINT bid_bidder_id_fkey FOREIGN KEY (bidder_id)
REFERENCES users(id);