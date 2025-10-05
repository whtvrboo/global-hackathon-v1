-- GOD-TIER DEMO SEED DATA FOR FOLIO
-- This creates the perfect demo data to showcase the curation engine

-- Clear existing data (in reverse dependency order)
DELETE FROM list_comments;
DELETE FROM list_likes;
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

-- Insert the demo curator user
INSERT INTO users (id, username, name, email, picture, created_at) VALUES
('demo-user-1', 'alex_curator', 'Alex the Curator', 'alex@folio-demo.com', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '1 year');

-- Insert additional demo users for social features
INSERT INTO users (id, username, name, email, picture, created_at) VALUES
('demo-user-2', 'bookworm_sarah', 'Sarah Chen', 'sarah@folio-demo.com', 'https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '8 months'),
('demo-user-3', 'mike_reader', 'Mike Rodriguez', 'mike@folio-demo.com', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '6 months'),
('demo-user-4', 'emma_books', 'Emma Thompson', 'emma@folio-demo.com', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face', NOW() - INTERVAL '4 months');

-- Insert the demo books
INSERT INTO books (id, title, authors, isbn, published_date, description, cover_url, page_count, language, created_at) VALUES
-- Sci-Fi Masterpieces
('demo-book-1', 'Project Hail Mary', ARRAY['Andy Weir'], '9780593135204', '2021-05-04', 'A lone astronaut must save the earth from disaster in this incredible new science-based thriller from the #1 New York Times bestselling author of The Martian.', 'https://covers.openlibrary.org/b/isbn/9780593135204-L.jpg', 496, 'English', NOW()),
('demo-book-2', 'The Martian', ARRAY['Andy Weir'], '9780553418026', '2014-02-11', 'Six days ago, astronaut Mark Watney became one of the first people to walk on Mars. Now, he''s sure he''ll be the first person to die there.', 'https://covers.openlibrary.org/b/isbn/9780553418026-L.jpg', 369, 'English', NOW()),
('demo-book-3', 'Dune', ARRAY['Frank Herbert'], '9780441172719', '1965-08-01', 'Set on the desert planet Arrakis, Dune is the story of the boy Paul Atreides, heir to a noble family tasked with ruling an inhospitable world where the only thing of value is the "spice" melange.', 'https://covers.openlibrary.org/b/isbn/9780441172719-L.jpg', 688, 'English', NOW()),
('demo-book-4', 'The Expanse: Leviathan Wakes', ARRAY['James S.A. Corey'], '9780316129084', '2011-06-15', 'Humanity has colonized the solar system—Mars, the Moon, the Asteroid Belt and beyond—but the stars are still out of our reach.', 'https://covers.openlibrary.org/b/isbn/9780316129084-L.jpg', 592, 'English', NOW()),
('demo-book-5', 'Hyperion', ARRAY['Dan Simmons'], '9780553288200', '1989-05-26', 'On the world called Hyperion, beyond the law of the Hegemony of Man, there waits the creature called the Shrike.', 'https://covers.openlibrary.org/b/isbn/9780553288200-L.jpg', 482, 'English', NOW()),

-- Modern Stoic's Library
('demo-book-6', 'Meditations', ARRAY['Marcus Aurelius'], '9780486298238', '180-01-01', 'A series of personal writings by Marcus Aurelius, Roman Emperor from 161 to 180 AD, recording his private notes to himself and ideas on Stoic philosophy.', 'https://covers.openlibrary.org/b/isbn/9780486298238-L.jpg', 256, 'English', NOW()),
('demo-book-7', 'Letters from a Stoic', ARRAY['Seneca'], '9780140442106', '65-01-01', 'A selection of Seneca''s most significant letters that illuminate the intellectual and spiritual life of one of the most influential philosophers of the ancient world.', 'https://covers.openlibrary.org/b/isbn/9780140442106-L.jpg', 254, 'English', NOW()),
('demo-book-8', 'The Daily Stoic', ARRAY['Ryan Holiday', 'Stephen Hanselman'], '9780735211735', '2016-10-18', '366 meditations on wisdom, perseverance, and the art of living from the bestselling author of The Obstacle is the Way.', 'https://covers.openlibrary.org/b/isbn/9780735211735-L.jpg', 416, 'English', NOW()),
('demo-book-9', 'A Guide to the Good Life', ARRAY['William B. Irvine'], '9780195374612', '2008-11-01', 'One of the great fears many of us face is that despite all our effort and striving, we will discover at the end that we have wasted our life.', 'https://covers.openlibrary.org/b/isbn/9780195374612-L.jpg', 336, 'English', NOW()),
('demo-book-10', 'The Obstacle Is the Way', ARRAY['Ryan Holiday'], '9781591846352', '2014-05-01', 'The Obstacle is the Way has become a cult classic, beloved by men and women around the world who apply its wisdom to become more successful at whatever they do.', 'https://covers.openlibrary.org/b/isbn/9781591846352-L.jpg', 224, 'English', NOW()),

-- Japanese Fiction
('demo-book-11', 'Norwegian Wood', ARRAY['Haruki Murakami'], '9780375704024', '1987-09-04', 'A poignant coming-of-age, unusual love story, limned with the same air of desolation that marked all of Murakami''s early work.', 'https://covers.openlibrary.org/b/isbn/9780375704024-L.jpg', 296, 'English', NOW()),
('demo-book-12', 'The Wind-Up Bird Chronicle', ARRAY['Haruki Murakami'], '9780679775430', '1994-09-12', 'A "hypnotic" (The New York Times Book Review) novel from one of Japan''s most celebrated authors.', 'https://covers.openlibrary.org/b/isbn/9780679775430-L.jpg', 607, 'English', NOW()),
('demo-book-13', 'Snow Country', ARRAY['Yasunari Kawabata'], '9780679755333', '1935-01-01', 'A tale of wasted love set amid the desolate beauty of western Japan, where snow falls on the only two geishas in a hot-spring town.', 'https://covers.openlibrary.org/b/isbn/9780679755333-L.jpg', 175, 'English', NOW()),
('demo-book-14', 'The Memory Police', ARRAY['Yoko Ogawa'], '9781101870600', '1994-01-01', 'A haunting, Orwellian novel about the terrors of state surveillance, from the acclaimed author of The Housekeeper and the Professor.', 'https://covers.openlibrary.org/b/isbn/9781101870600-L.jpg', 274, 'English', NOW()),
('demo-book-15', 'Convenience Store Woman', ARRAY['Sayaka Murata'], '9780802128256', '2016-07-26', 'A brilliant, deeply satisfying story of a woman who has been unable to feel at home in the world, and the convenience store where she works.', 'https://covers.openlibrary.org/b/isbn/9780802128256-L.jpg', 163, 'English', NOW());

-- Insert followers relationships
INSERT INTO followers (follower_id, following_id, created_at) VALUES
('demo-user-2', 'demo-user-1', NOW() - INTERVAL '6 months'),
('demo-user-3', 'demo-user-1', NOW() - INTERVAL '4 months'),
('demo-user-4', 'demo-user-1', NOW() - INTERVAL '2 months'),
('demo-user-1', 'demo-user-2', NOW() - INTERVAL '5 months'),
('demo-user-1', 'demo-user-3', NOW() - INTERVAL '3 months');

-- Insert the three GOD-TIER demo lists
INSERT INTO lists (id, user_id, name, description, is_public, header_image_url, theme_color, items_count, created_at, updated_at) VALUES
-- List 1: Sci-Fi Masterpieces of the 21st Century
('demo-list-1', 'demo-user-1', 'Sci-Fi Masterpieces of the 21st Century', 'A carefully curated collection of the most groundbreaking science fiction novels that have defined the 21st century. These books don''t just tell stories—they expand our understanding of what''s possible.', true, 'https://images.unsplash.com/photo-1446776877081-d282a0f896e2?w=800&h=400&fit=crop', '#1e40af', 5, NOW() - INTERVAL '3 months', NOW() - INTERVAL '1 week'),

-- List 2: The Modern Stoic''s Library
('demo-list-2', 'demo-user-1', 'The Modern Stoic''s Library', 'Ancient wisdom for modern life. These books offer timeless guidance on living with purpose, resilience, and inner peace in our chaotic world.', true, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=800&h=400&fit=crop', '#059669', 5, NOW() - INTERVAL '2 months', NOW() - INTERVAL '2 weeks'),

-- List 3: A Beginner''s Guide to Japanese Fiction
('demo-list-3', 'demo-user-1', 'A Beginner''s Guide to Japanese Fiction', 'Discover the unique beauty of Japanese literature. From Murakami''s surreal worlds to Kawabata''s delicate prose, these books offer a perfect introduction to one of the world''s most distinctive literary traditions.', true, 'https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?w=800&h=400&fit=crop', '#dc2626', 5, NOW() - INTERVAL '1 month', NOW() - INTERVAL '3 days');

-- Insert list items with thoughtful notes and perfect ordering
INSERT INTO list_items (list_id, book_id, notes, item_order, created_at) VALUES
-- Sci-Fi Masterpieces (ordered by impact and accessibility)
('demo-list-1', 'demo-book-1', 'Weir''s masterpiece that makes hard science accessible and emotionally resonant. Rocky might be the best alien character ever written.', 0, NOW() - INTERVAL '3 months'),
('demo-list-1', 'demo-book-2', 'The book that launched a thousand memes and proved that science fiction can be both accurate and entertaining. Watney''s humor in the face of certain death is unforgettable.', 1, NOW() - INTERVAL '3 months'),
('demo-list-1', 'demo-book-3', 'The foundation of modern space opera. Herbert created an entire universe with its own ecology, politics, and philosophy. Essential reading.', 2, NOW() - INTERVAL '3 months'),
('demo-list-1', 'demo-book-4', 'Gritty, realistic space opera that feels like it could happen tomorrow. The political intrigue is as compelling as the science.', 3, NOW() - INTERVAL '3 months'),
('demo-list-1', 'demo-book-5', 'A literary masterpiece disguised as science fiction. Simmons weaves together multiple narratives to create something truly unique.', 4, NOW() - INTERVAL '3 months'),

-- Modern Stoic''s Library (ordered from ancient to modern)
('demo-list-2', 'demo-book-6', 'The emperor-philosopher''s private thoughts on virtue, duty, and the nature of reality. Timeless wisdom that feels surprisingly modern.', 0, NOW() - INTERVAL '2 months'),
('demo-list-2', 'demo-book-7', 'Seneca''s letters offer practical advice on everything from dealing with anger to finding inner peace. His voice feels like a wise friend.', 1, NOW() - INTERVAL '2 months'),
('demo-list-2', 'demo-book-8', 'The perfect daily companion. Each page offers a bite-sized piece of wisdom to start your day with intention and purpose.', 2, NOW() - INTERVAL '2 months'),
('demo-list-2', 'demo-book-9', 'A modern introduction to Stoicism that makes ancient philosophy accessible. Irvine explains why Stoicism is more relevant than ever.', 3, NOW() - INTERVAL '2 months'),
('demo-list-2', 'demo-book-10', 'Holiday shows how to turn obstacles into opportunities. Essential reading for anyone facing challenges in life or business.', 4, NOW() - INTERVAL '2 months'),

-- Japanese Fiction (ordered from accessible to more challenging)
('demo-list-3', 'demo-book-11', 'Murakami''s most accessible novel. A beautiful, melancholic coming-of-age story that perfectly captures the feeling of young love and loss.', 0, NOW() - INTERVAL '1 month'),
('demo-list-3', 'demo-book-12', 'Murakami at his most surreal and ambitious. A dreamlike journey through memory, identity, and the subconscious mind.', 1, NOW() - INTERVAL '1 month'),
('demo-list-3', 'demo-book-13', 'Kawabata''s Nobel Prize-winning masterpiece. Every sentence is crafted with the precision of a haiku. Pure poetry in prose form.', 2, NOW() - INTERVAL '1 month'),
('demo-list-3', 'demo-book-14', 'A haunting meditation on memory, loss, and the power of the state. Ogawa creates an atmosphere of quiet dread that''s impossible to forget.', 3, NOW() - INTERVAL '1 month'),
('demo-list-3', 'demo-book-15', 'A brilliant exploration of conformity and individuality. Murata''s deadpan humor makes this a surprisingly funny read about a serious subject.', 4, NOW() - INTERVAL '1 month');

-- Insert some likes on the lists
INSERT INTO list_likes (list_id, user_id, created_at) VALUES
('demo-list-1', 'demo-user-2', NOW() - INTERVAL '2 months'),
('demo-list-1', 'demo-user-3', NOW() - INTERVAL '1 month'),
('demo-list-1', 'demo-user-4', NOW() - INTERVAL '3 weeks'),
('demo-list-2', 'demo-user-2', NOW() - INTERVAL '1 month'),
('demo-list-2', 'demo-user-3', NOW() - INTERVAL '2 weeks'),
('demo-list-3', 'demo-user-2', NOW() - INTERVAL '3 weeks'),
('demo-list-3', 'demo-user-4', NOW() - INTERVAL '1 week');

-- Insert some comments on the lists
INSERT INTO list_comments (list_id, user_id, content, created_at, updated_at) VALUES
('demo-list-1', 'demo-user-2', 'This list is incredible! I''ve read three of these and they''re all amazing. Adding the rest to my TBR immediately.', NOW() - INTERVAL '2 months'),
('demo-list-1', 'demo-user-3', 'Perfect curation! The notes really help explain why each book is special. Rocky from Project Hail Mary is my favorite character ever.', NOW() - INTERVAL '1 month'),
('demo-list-2', 'demo-user-2', 'Exactly what I needed. The Daily Stoic has been life-changing for me. Thank you for this thoughtful collection.', NOW() - INTERVAL '1 month'),
('demo-list-2', 'demo-user-3', 'Marcus Aurelius and Seneca are timeless. This list is a perfect introduction to Stoic philosophy.', NOW() - INTERVAL '2 weeks'),
('demo-list-3', 'demo-user-2', 'I''ve been wanting to explore Japanese fiction but didn''t know where to start. This is perfect!', NOW() - INTERVAL '3 weeks'),
('demo-list-3', 'demo-user-4', 'Murakami is a genius. Norwegian Wood broke my heart in the best way possible.', NOW() - INTERVAL '1 week');

-- Insert some reading logs for Alex the Curator
INSERT INTO logs (id, user_id, book_id, status, rating, review, notes, start_date, finish_date, is_public, created_at, updated_at) VALUES
('demo-log-1', 'demo-user-1', 'demo-book-1', 'read', 5, 'Absolutely brilliant! Weir has outdone himself with this one. The science is fascinating, the humor is perfect, and Rocky is one of the most endearing characters I''ve ever encountered. This book made me laugh, cry, and think deeply about what it means to be human.', 'Read in one sitting - couldn''t put it down!', '2023-01-15', '2023-01-16', true, '2023-01-16', '2023-01-16'),
('demo-log-2', 'demo-user-1', 'demo-book-6', 'read', 5, 'Reading the private thoughts of a Roman emperor is a profound experience. Marcus Aurelius'' wisdom is timeless and his reflections on virtue, duty, and the nature of reality are as relevant today as they were 1800 years ago. This is a book I''ll return to throughout my life.', 'Daily meditation companion', '2023-02-01', '2023-02-15', true, '2023-02-15', '2023-02-15'),
('demo-log-3', 'demo-user-1', 'demo-book-11', 'read', 4, 'Murakami''s most accessible novel, but no less powerful for it. The melancholy atmosphere and themes of love, loss, and memory are beautifully rendered. There''s something haunting about the way he captures the feeling of being young and lost in the world.', 'Perfect introduction to Murakami', '2023-03-01', '2023-03-05', true, '2023-03-05', '2023-03-05');

-- Insert likes on Alex's logs
INSERT INTO log_likes (user_id, log_id, created_at) VALUES
('demo-user-2', 'demo-log-1', NOW() - INTERVAL '2 months'),
('demo-user-3', 'demo-log-1', NOW() - INTERVAL '1 month'),
('demo-user-4', 'demo-log-1', NOW() - INTERVAL '3 weeks'),
('demo-user-2', 'demo-log-2', NOW() - INTERVAL '1 month'),
('demo-user-3', 'demo-log-2', NOW() - INTERVAL '2 weeks'),
('demo-user-2', 'demo-log-3', NOW() - INTERVAL '3 weeks'),
('demo-user-4', 'demo-log-3', NOW() - INTERVAL '1 week');

-- Insert comments on Alex's logs
INSERT INTO log_comments (user_id, log_id, content, created_at, updated_at) VALUES
('demo-user-2', 'demo-log-1', 'I completely agree! Rocky is such a wonderful character. The way Weir makes hard science accessible is incredible.', NOW() - INTERVAL '2 months'),
('demo-user-3', 'demo-log-1', 'This book made me want to learn more about space and science. Weir is a genius at making complex concepts entertaining.', NOW() - INTERVAL '1 month'),
('demo-user-2', 'demo-log-2', 'Marcus Aurelius is so wise. I love how you can open this book to any page and find something profound.', NOW() - INTERVAL '1 month'),
('demo-user-4', 'demo-log-3', 'Murakami has such a unique voice. Norwegian Wood is beautiful and heartbreaking.', NOW() - INTERVAL '1 week');

-- Update counts (triggers should handle this, but let's ensure they're correct)
UPDATE lists SET 
    items_count = (SELECT COUNT(*) FROM list_items WHERE list_id = lists.id);

UPDATE logs SET 
    likes_count = (SELECT COUNT(*) FROM log_likes WHERE log_id = logs.id),
    comments_count = (SELECT COUNT(*) FROM log_comments WHERE log_id = logs.id);
