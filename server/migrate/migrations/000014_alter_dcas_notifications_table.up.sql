ALTER TABLE notification
DROP CONSTRAINT IF EXISTS notification_user_id_fkey;

ALTER TABLE notification
DROP CONSTRAINT IF EXISTS notification_auction_id_fkey;

ALTER TABLE notification
ADD CONSTRAINT notification_user_id_fkey FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE notification
ADD CONSTRAINT notification_auction_id_fkey FOREIGN KEY (auction_id)
REFERENCES auctions(id) ON DELETE CASCADE;