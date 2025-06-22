ALTER TABLE notification
ADD COLUMN auction_id UUID,
ADD CONSTRAINT fk_auction
FOREIGN KEY (auction_id) REFERENCES auctions(id);
