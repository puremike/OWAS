ALTER TABLE notification
DROP CONSTRAINT IF EXISTS notification_user_id_fkey;

ALTER TABLE notification
ADD CONSTRAINT notification_user_id_fkey FOREIGN KEY (user_id)
REFERENCES users(id);