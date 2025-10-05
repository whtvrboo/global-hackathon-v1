-- Comprehensive seed data for Folio
-- This creates realistic demo data to showcase all features

-- Clear existing data (in reverse dependency order)
DELETE FROM log_comments;
DELETE FROM log_likes;
DELETE FROM list_items;
DELETE FROM lists;
DELETE FROM logs;
DELETE FROM followers;
DELETE FROM books;
DELETE FROM users;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE books_id_seq RESTART WITH 1;

-- Insert realistic users
INSERT INTO users (id, username, name, email, picture, created_at) VALUES
('user-1', 'alex_reader', 'Alex Johnson', 'alex@example.com', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '2 years'),
('user-2', 'bookworm_sarah', 'Sarah Chen', 'sarah@example.com', 'https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '18 months'),
('user-3', 'mike_literature', 'Mike Rodriguez', 'mike@example.com', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '1 year'),
('user-4', 'emma_books', 'Emma Thompson', 'emma@example.com', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '8 months'),
('user-5', 'david_reader', 'David Kim', 'david@example.com', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '6 months');

-- Load comprehensive book database
\i books_seed.sql

-- Insert followers relationships
INSERT INTO followers (follower_id, following_id, created_at) VALUES
('user-1', 'user-2', NOW() - INTERVAL '1 year'),
('user-1', 'user-3', NOW() - INTERVAL '10 months'),
('user-2', 'user-1', NOW() - INTERVAL '8 months'),
('user-2', 'user-4', NOW() - INTERVAL '6 months'),
('user-3', 'user-1', NOW() - INTERVAL '9 months'),
('user-3', 'user-5', NOW() - INTERVAL '4 months'),
('user-4', 'user-2', NOW() - INTERVAL '5 months'),
('user-4', 'user-3', NOW() - INTERVAL '3 months'),
('user-5', 'user-1', NOW() - INTERVAL '2 months'),
('user-5', 'user-4', NOW() - INTERVAL '1 month');

-- Insert diverse reading logs
INSERT INTO logs (id, user_id, book_id, status, rating, review, notes, start_date, finish_date, is_public, created_at, updated_at) VALUES
-- Alex's logs (user-1)
('log-1', 'user-1', 'book-1', 'read', 5, 'Absolutely stunning! The way Taylor Jenkins Reid weaves together Hollywood glamour with deep emotional truths is masterful. Evelyn Hugo is one of the most complex characters I''ve ever encountered.', 'Read for book club. Couldn''t put it down!', '2023-01-15', '2023-01-20', true, '2023-01-20', '2023-01-20'),
('log-2', 'user-1', 'book-2', 'read', 4, 'Andy Weir does it again! The science is fascinating and the story is gripping. Ryland Grace is such a lovable protagonist.', 'Great audiobook experience', '2023-02-01', '2023-02-10', true, '2023-02-10', '2023-02-10'),
('log-3', 'user-1', 'book-6', 'read', 5, 'A masterpiece of science fiction. The world-building is incredible and the political intrigue is fascinating.', 'Re-read for the movie. Even better the second time!', '2023-03-01', '2023-03-15', true, '2023-03-15', '2023-03-15'),
('log-4', 'user-1', 'book-3', 'reading', 4, 'Really enjoying this so far. The concept is intriguing and Matt Haig''s writing is beautiful.', 'About halfway through', '2023-04-01', NULL, true, '2023-04-01', '2023-04-01'),
('log-5', 'user-1', 'book-7', 'want_to_read', NULL, NULL, 'Heard amazing things about this retelling', NULL, NULL, true, '2023-04-05', '2023-04-05');

-- Sarah's logs (user-2)
('log-6', 'user-2', 'book-1', 'read', 4, 'Beautiful writing and compelling characters. The structure of the story is really clever.', 'Perfect beach read', '2023-01-10', '2023-01-18', true, '2023-01-18', '2023-01-18'),
('log-7', 'user-2', 'book-4', 'read', 5, 'Ishiguro at his finest. Klara''s perspective is so unique and touching. Made me think about AI and humanity in new ways.', 'Book club pick - great discussion!', '2023-02-15', '2023-02-25', true, '2023-02-25', '2023-02-25'),
('log-8', 'user-2', 'book-8', 'read', 5, 'Circe is such a powerful character. Miller''s writing is absolutely gorgeous.', 'Couldn''t stop reading', '2023-03-05', '2023-03-12', true, '2023-03-12', '2023-03-12'),
('log-9', 'user-2', 'book-9', 'reading', 3, 'Sweet and charming so far. The characters are endearing.', 'Light reading for bedtime', '2023-04-02', NULL, true, '2023-04-02', '2023-04-02'),
('log-10', 'user-2', 'book-10', 'want_to_read', NULL, NULL, 'Love Susanna Clarke''s writing', NULL, NULL, true, '2023-04-08', '2023-04-08');

-- Mike's logs (user-3)
('log-11', 'user-3', 'book-2', 'read', 5, 'Incredible! The science is spot-on and the story is both funny and touching. Weir is a genius.', 'Best sci-fi I''ve read in years', '2023-01-20', '2023-01-28', true, '2023-01-28', '2023-01-28'),
('log-12', 'user-3', 'book-6', 'read', 4, 'Classic for a reason. The world-building is phenomenal.', 'Finally got around to reading this classic', '2023-02-10', '2023-02-25', true, '2023-02-25', '2023-02-25'),
('log-13', 'user-3', 'book-11', 'read', 3, 'Fun mystery with great characters. The humor is delightful.', 'Light and entertaining', '2023-03-15', '2023-03-20', true, '2023-03-20', '2023-03-20'),
('log-14', 'user-3', 'book-12', 'read', 4, 'Great psychological thriller. The twist was unexpected!', 'Page-turner', '2023-04-01', '2023-04-05', true, '2023-04-05', '2023-04-05'),
('log-15', 'user-3', 'book-5', 'want_to_read', NULL, NULL, 'Heard this is amazing', NULL, NULL, true, '2023-04-10', '2023-04-10');

-- Emma's logs (user-4)
('log-16', 'user-4', 'book-3', 'read', 4, 'Beautiful and thought-provoking. Made me reflect on my own life choices.', 'Perfect for a rainy day', '2023-01-25', '2023-02-02', true, '2023-02-02', '2023-02-02'),
('log-17', 'user-4', 'book-7', 'read', 5, 'Absolutely heartbreaking and beautiful. Miller''s writing is pure poetry.', 'Cried multiple times', '2023-02-20', '2023-02-28', true, '2023-02-28', '2023-02-28'),
('log-18', 'user-4', 'book-8', 'read', 5, 'Circe''s journey is so powerful. Another masterpiece from Miller.', 'Even better than Song of Achilles', '2023-03-10', '2023-03-18', true, '2023-03-18', '2023-03-18'),
('log-19', 'user-4', 'book-9', 'read', 4, 'Sweet and heartwarming. Perfect comfort read.', 'Made me smile throughout', '2023-04-01', '2023-04-08', true, '2023-04-08', '2023-04-08'),
('log-20', 'user-4', 'book-1', 'want_to_read', NULL, NULL, 'Everyone is talking about this one', NULL, NULL, true, '2023-04-12', '2023-04-12');

-- David's logs (user-5)
('log-21', 'user-5', 'book-4', 'read', 4, 'Thought-provoking and beautifully written. Klara''s perspective is fascinating.', 'Great for book club discussion', '2023-02-01', '2023-02-10', true, '2023-02-10', '2023-02-10'),
('log-22', 'user-5', 'book-10', 'read', 5, 'Absolutely mesmerizing. Clarke''s writing is incredible.', 'Couldn''t put it down', '2023-03-01', '2023-03-08', true, '2023-03-08', '2023-03-08'),
('log-23', 'user-5', 'book-11', 'read', 3, 'Fun and light mystery. The characters are charming.', 'Good beach read', '2023-03-20', '2023-03-25', true, '2023-03-25', '2023-03-25'),
('log-24', 'user-5', 'book-12', 'reading', 4, 'Really enjoying this thriller so far. The pacing is perfect.', 'About 2/3 through', '2023-04-05', NULL, true, '2023-04-05', '2023-04-05'),
('log-25', 'user-5', 'book-2', 'want_to_read', NULL, NULL, 'Love Andy Weir''s books', NULL, NULL, true, '2023-04-15', '2023-04-15');

-- Insert some rereads (multiple logs for same book by same user)
INSERT INTO logs (id, user_id, book_id, status, rating, review, notes, start_date, finish_date, is_public, created_at, updated_at) VALUES
('log-26', 'user-1', 'book-1', 'read', 5, 'Even better the second time! Picked up so many details I missed.', 'Reread for book club discussion', '2023-03-01', '2023-03-05', true, '2023-03-05', '2023-03-05'),
('log-27', 'user-2', 'book-8', 'read', 5, 'Still amazing on reread. Circe''s character development is even more powerful.', 'Reread before the TV adaptation', '2023-04-01', '2023-04-08', true, '2023-04-08', '2023-04-08');

-- Insert likes on logs
INSERT INTO log_likes (user_id, log_id, created_at) VALUES
-- Likes on Alex's logs
('user-2', 'log-1', NOW() - INTERVAL '2 months'),
('user-3', 'log-1', NOW() - INTERVAL '1 month'),
('user-4', 'log-1', NOW() - INTERVAL '3 weeks'),
('user-2', 'log-2', NOW() - INTERVAL '1 month'),
('user-5', 'log-2', NOW() - INTERVAL '2 weeks'),
('user-3', 'log-3', NOW() - INTERVAL '3 weeks'),
('user-4', 'log-3', NOW() - INTERVAL '2 weeks'),

-- Likes on Sarah's logs
('user-1', 'log-6', NOW() - INTERVAL '2 months'),
('user-3', 'log-6', NOW() - INTERVAL '1 month'),
('user-1', 'log-7', NOW() - INTERVAL '1 month'),
('user-4', 'log-7', NOW() - INTERVAL '3 weeks'),
('user-5', 'log-7', NOW() - INTERVAL '2 weeks'),
('user-1', 'log-8', NOW() - INTERVAL '3 weeks'),
('user-3', 'log-8', NOW() - INTERVAL '2 weeks'),

-- Likes on Mike's logs
('user-1', 'log-11', NOW() - INTERVAL '2 months'),
('user-2', 'log-11', NOW() - INTERVAL '1 month'),
('user-4', 'log-11', NOW() - INTERVAL '3 weeks'),
('user-2', 'log-12', NOW() - INTERVAL '1 month'),
('user-5', 'log-12', NOW() - INTERVAL '2 weeks'),
('user-1', 'log-13', NOW() - INTERVAL '3 weeks'),
('user-4', 'log-13', NOW() - INTERVAL '2 weeks'),

-- Likes on Emma's logs
('user-2', 'log-16', NOW() - INTERVAL '2 months'),
('user-3', 'log-16', NOW() - INTERVAL '1 month'),
('user-1', 'log-17', NOW() - INTERVAL '1 month'),
('user-2', 'log-17', NOW() - INTERVAL '3 weeks'),
('user-5', 'log-17', NOW() - INTERVAL '2 weeks'),
('user-1', 'log-18', NOW() - INTERVAL '3 weeks'),
('user-3', 'log-18', NOW() - INTERVAL '2 weeks'),

-- Likes on David's logs
('user-1', 'log-21', NOW() - INTERVAL '2 months'),
('user-4', 'log-21', NOW() - INTERVAL '1 month'),
('user-2', 'log-22', NOW() - INTERVAL '1 month'),
('user-3', 'log-22', NOW() - INTERVAL '3 weeks'),
('user-1', 'log-23', NOW() - INTERVAL '2 weeks');

-- Insert comments on logs
INSERT INTO log_comments (user_id, log_id, content, created_at, updated_at) VALUES
-- Comments on Alex's logs
('user-2', 'log-1', 'I completely agree! Evelyn Hugo is such a complex character. The ending had me in tears.', NOW() - INTERVAL '2 months'),
('user-3', 'log-1', 'Taylor Jenkins Reid is a master storyteller. Have you read Daisy Jones & The Six?', NOW() - INTERVAL '1 month'),
('user-4', 'log-1', 'This is on my TBR! Your review makes me want to read it even more.', NOW() - INTERVAL '3 weeks'),
('user-2', 'log-2', 'The science in this book is incredible! Weir really knows how to make complex concepts accessible.', NOW() - INTERVAL '1 month'),
('user-5', 'log-2', 'Rocky is the best character! 😄', NOW() - INTERVAL '2 weeks'),

-- Comments on Sarah's logs
('user-1', 'log-7', 'Ishiguro never disappoints. Klara''s perspective is so unique and touching.', NOW() - INTERVAL '1 month'),
('user-4', 'log-7', 'This book made me think about AI and consciousness in completely new ways.', NOW() - INTERVAL '3 weeks'),
('user-1', 'log-8', 'Circe is one of my all-time favorite characters. Miller''s writing is pure magic.', NOW() - INTERVAL '3 weeks'),
('user-3', 'log-8', 'The way she transforms throughout the story is incredible.', NOW() - INTERVAL '2 weeks'),

-- Comments on Mike's logs
('user-1', 'log-11', 'Andy Weir is a genius! The humor mixed with hard science is perfect.', NOW() - INTERVAL '2 months'),
('user-2', 'log-11', 'I loved the audiobook version. The narrator was fantastic!', NOW() - INTERVAL '1 month'),
('user-4', 'log-11', 'The ending had me on the edge of my seat!', NOW() - INTERVAL '3 weeks'),

-- Comments on Emma's logs
('user-2', 'log-17', 'I cried so much reading this! Patroclus and Achilles'' love story is beautiful.', NOW() - INTERVAL '1 month'),
('user-1', 'log-17', 'Miller''s retelling is absolutely stunning. The prose is like poetry.', NOW() - INTERVAL '3 weeks'),
('user-3', 'log-17', 'This book destroyed me emotionally in the best way possible.', NOW() - INTERVAL '2 weeks'),

-- Comments on David's logs
('user-1', 'log-22', 'Piranesi is such a unique book. Clarke''s imagination is incredible.', NOW() - INTERVAL '1 month'),
('user-2', 'log-22', 'The world-building is so detailed and immersive. I felt like I was there.', NOW() - INTERVAL '3 weeks'),
('user-4', 'log-22', 'This book stayed with me for weeks after reading it.', NOW() - INTERVAL '2 weeks');

-- Insert custom lists
INSERT INTO lists (id, user_id, name, description, is_public, items_count, created_at, updated_at) VALUES
-- Alex's lists
('list-1', 'user-1', '2023 Favorites', 'My favorite books read this year', true, 3, NOW() - INTERVAL '3 months', NOW() - INTERVAL '1 month'),
('list-2', 'user-1', 'Sci-Fi Masterpieces', 'The best science fiction books I''ve read', true, 2, NOW() - INTERVAL '2 months', NOW() - INTERVAL '2 weeks'),
('list-3', 'user-1', 'Book Club Picks', 'Books I want to suggest for our book club', false, 2, NOW() - INTERVAL '1 month', NOW() - INTERVAL '1 week'),

-- Sarah's lists
('list-4', 'user-2', 'Greek Mythology Retellings', 'Modern retellings of Greek myths', true, 2, NOW() - INTERVAL '2 months', NOW() - INTERVAL '3 weeks'),
('list-5', 'user-2', 'Beach Reads', 'Perfect books for vacation reading', true, 2, NOW() - INTERVAL '1 month', NOW() - INTERVAL '2 weeks'),
('list-6', 'user-2', 'TBR Priority', 'Books I really want to read soon', false, 3, NOW() - INTERVAL '3 weeks', NOW() - INTERVAL '1 week'),

-- Mike's lists
('list-7', 'user-3', 'Thriller Recommendations', 'Gripping thrillers that kept me up all night', true, 2, NOW() - INTERVAL '1 month', NOW() - INTERVAL '2 weeks'),
('list-8', 'user-3', 'Classic Sci-Fi', 'Must-read science fiction classics', true, 1, NOW() - INTERVAL '2 weeks', NOW() - INTERVAL '1 week'),

-- Emma's lists
('list-9', 'user-4', 'Literary Fiction Gems', 'Beautifully written literary fiction', true, 3, NOW() - INTERVAL '2 months', NOW() - INTERVAL '1 month'),
('list-10', 'user-4', 'Comfort Reads', 'Books that make me feel warm and happy', true, 1, NOW() - INTERVAL '1 month', NOW() - INTERVAL '2 weeks'),

-- David's lists
('list-11', 'user-5', 'Mind-Bending Fiction', 'Books that completely changed my perspective', true, 2, NOW() - INTERVAL '1 month', NOW() - INTERVAL '2 weeks'),
('list-12', 'user-5', 'Quick Reads', 'Short books perfect for busy schedules', true, 2, NOW() - INTERVAL '2 weeks', NOW() - INTERVAL '1 week');

-- Insert list items
INSERT INTO list_items (list_id, book_id, notes, item_order, created_at) VALUES
-- Alex's list items
('list-1', 'book-1', 'Absolutely stunning!', 0, NOW() - INTERVAL '3 months'),
('list-1', 'book-2', 'Andy Weir at his best', 1, NOW() - INTERVAL '2 months'),
('list-1', 'book-6', 'A masterpiece', 2, NOW() - INTERVAL '1 month'),
('list-2', 'book-2', 'Hard science fiction done right', 0, NOW() - INTERVAL '2 months'),
('list-2', 'book-6', 'The gold standard of sci-fi', 1, NOW() - INTERVAL '1 month'),
('list-3', 'book-3', 'Great for discussion', 0, NOW() - INTERVAL '1 month'),
('list-3', 'book-7', 'Emotional and beautiful', 1, NOW() - INTERVAL '1 week'),

-- Sarah's list items
('list-4', 'book-7', 'Heartbreaking and beautiful', 0, NOW() - INTERVAL '2 months'),
('list-4', 'book-8', 'Circe is incredible', 1, NOW() - INTERVAL '1 month'),
('list-5', 'book-1', 'Perfect vacation read', 0, NOW() - INTERVAL '1 month'),
('list-5', 'book-9', 'Sweet and charming', 1, NOW() - INTERVAL '2 weeks'),
('list-6', 'book-3', 'Really want to read this', 0, NOW() - INTERVAL '3 weeks'),
('list-6', 'book-5', 'Heard amazing things', 1, NOW() - INTERVAL '2 weeks'),
('list-6', 'book-10', 'Love Susanna Clarke', 2, NOW() - INTERVAL '1 week'),

-- Mike's list items
('list-7', 'book-11', 'Fun and engaging mystery', 0, NOW() - INTERVAL '1 month'),
('list-7', 'book-12', 'Great psychological thriller', 1, NOW() - INTERVAL '2 weeks'),
('list-8', 'book-6', 'The classic that started it all', 0, NOW() - INTERVAL '2 weeks'),

-- Emma's list items
('list-9', 'book-3', 'Thought-provoking and beautiful', 0, NOW() - INTERVAL '2 months'),
('list-9', 'book-4', 'Ishiguro at his finest', 1, NOW() - INTERVAL '1 month'),
('list-9', 'book-10', 'Mesmerizing and unique', 2, NOW() - INTERVAL '2 weeks'),
('list-10', 'book-9', 'Perfect comfort read', 0, NOW() - INTERVAL '1 month'),

-- David's list items
('list-11', 'book-4', 'Changed how I think about AI', 0, NOW() - INTERVAL '1 month'),
('list-11', 'book-10', 'Completely mind-bending', 1, NOW() - INTERVAL '2 weeks'),
('list-12', 'book-11', 'Quick and entertaining', 0, NOW() - INTERVAL '2 weeks'),
('list-12', 'book-12', 'Fast-paced thriller', 1, NOW() - INTERVAL '1 week');

-- Update the counts in logs table (triggers should handle this, but let's ensure they're correct)
UPDATE logs SET 
    likes_count = (SELECT COUNT(*) FROM log_likes WHERE log_id = logs.id),
    comments_count = (SELECT COUNT(*) FROM log_comments WHERE log_id = logs.id);

-- Update the counts in lists table (triggers should handle this, but let's ensure they're correct)
UPDATE lists SET 
    items_count = (SELECT COUNT(*) FROM list_items WHERE list_id = lists.id);


