USE repeatword;

DELIMITER |

CREATE PROCEDURE CreateWord (
    IN word_p VARCHAR(255),
    IN collection_id_p INT,
    IN means_p JSON
) 
DETERMINISTIC
BEGIN 
    DECLARE mean, examples JSON;
    DECLARE example VARCHAR(255);
    DECLARE i,j, vob_id, mean_id INT DEFAULT 0;

    -- insert vob
    INSERT INTO
        vobs(word)
    VALUES
        (word_p);

    SET vob_id = LAST_INSERT_ID();

    -- insert collection word
    INSERT INTO
        collection_words(vob_id, collection_id)
    VALUES
        (vob_id, collection_id_p);

    -- loop means insert mean, example
    WHILE i < JSON_LENGTH(means_p) DO
    SELECT
        JSON_EXTRACT(means_p, CONCAT('$[', i, ']')) INTO mean;

        -- insert means
        INSERT INTO
            means(vob_id, meaning, type)
        VALUES (
            vob_id,
            JSON_UNQUOTE(JSON_EXTRACT(mean, '$.meaning')),
            JSON_UNQUOTE(JSON_EXTRACT(mean, '$.type'))
        );

        SET mean_id = LAST_INSERT_ID();

        -- insert examples
        SET j = 0;

        SET examples = JSON_EXTRACT(mean, '$.examples');

        WHILE j < JSON_LENGTH(examples) DO
            SELECT JSON_EXTRACT(examples, CONCAT('$[', j, ']')) INTO example;
            
            INSERT INTO
                examples(mean_id, example)
            VALUES
                (mean_id, example); 
            
            SELECT j + 1 INTO j;

        END WHILE;

    SELECT
        i + 1 INTO i;

    END WHILE;
END |

DELIMITER ;