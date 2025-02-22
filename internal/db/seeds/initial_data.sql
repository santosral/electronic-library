WITH word_list AS (
    SELECT unnest(ARRAY[
        'The', 'Great', 'Mysterious', 'Ancient', 'Forbidden', 'Lost', 'Secret', 'Hidden',
        'Legend', 'Endless', 'Silent', 'Fateful', 'Invisible', 'Eternal', 'Dark', 'Bright',
        'Whispers', 'Shadows', 'Dream', 'Destiny', 'Cursed', 'Wings', 'Journey', 'Promise',
        'Empire', 'Crown', 'Song', 'Flame', 'Tale', 'Book', 'World', 'Force', 'Storm',
        'Revenge', 'Throne', 'Sea', 'Mountain', 'Magic', 'Kingdom', 'Fate', 'Soul', 'Truth',
        'Wanderer', 'Crossroads', 'Path', 'Legendary', 'Phoenix', 'Lands', 'War', 'Heart',
        'Time', 'Rising', 'Fallen', 'King', 'Queen'
    ]) AS word
),
combinations AS (
    SELECT
        w1.word || ' ' || w2.word || ' ' || w3.word AS title
    FROM word_list w1
    CROSS JOIN word_list w2
    CROSS JOIN word_list w3
    WHERE w1.word != w2.word
      AND w1.word != w3.word
      AND w2.word != w3.word
)

INSERT INTO book_details (title, available_copies)
SELECT 
    title,
    CASE 
        WHEN RANDOM() < 0.5 THEN 0
        ELSE 100
    END AS available_copies
FROM combinations
ORDER BY title;
