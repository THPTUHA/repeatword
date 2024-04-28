USE repeatword;

DELIMITER | 

CREATE PROCEDURE SetWord (
    IN collection_id_p INT,
    IN vob_p JSON
) 
DETERMINISTIC
BEGIN 
    DECLARE pronounces,means, examples,pro,mean,part, parts JSON;
    DECLARE meaning, example TEXT;
    DECLARE i,j,k, vob_id, mean_id, part_id,id_count INT DEFAULT 0;
    DECLARE word VARCHAR(255);

    START TRANSACTION;
    SET word = JSON_UNQUOTE(JSON_EXTRACT(vob_p, '$.word'));
    SET parts = JSON_EXTRACT(vob_p, '$.parts');
    
    -- insert vob
    INSERT INTO
        vobs(word,created_at)
    VALUES
        (word,CURRENT_TIMESTAMP);

    SET vob_id = LAST_INSERT_ID();

    -- check colleciton
    SELECT COUNT(*) INTO id_count FROM collections WHERE id = collection_id_p;
    IF id_count = 0 THEN
        INSERT INTO collections(id, name)
        VALUES (collection_id_p, CONCAT("new collection ", collection_id_p));
    END IF;

    -- insert collection word
    INSERT INTO
        collection_words(vob_id, collection_id)
    VALUES
        (vob_id, collection_id_p);

    -- loop means insert mean, example
    WHILE i < JSON_LENGTH(parts) DO
        SELECT JSON_EXTRACT(parts, CONCAT('$[', i, ']')) INTO part;
        SET pronounces = JSON_EXTRACT(part, '$.pronounces');
        SET means = JSON_EXTRACT(part, '$.means');

        -- insert part
        INSERT INTO
            vob_parts(vob_id, type, title)
        VALUES (
            vob_id,
            JSON_UNQUOTE(JSON_EXTRACT(part, '$.type')),
            JSON_UNQUOTE(JSON_EXTRACT(part, '$.title'))
        );

        SET part_id = LAST_INSERT_ID();
        -- insert pro 
        SET j = 0;
        WHILE j < JSON_LENGTH(pronounces) DO
            SELECT JSON_EXTRACT(pronounces, CONCAT('$[', j, ']')) INTO pro;
            INSERT INTO
                pronounces(audio_src, local_file, region, pro, part_id)
            VALUES
                (
                    JSON_UNQUOTE(JSON_EXTRACT(pro, '$.audio_src')),
                    JSON_UNQUOTE(JSON_EXTRACT(pro, '$.local_file')),
                    JSON_UNQUOTE(JSON_EXTRACT(pro, '$.region')),
                    JSON_UNQUOTE(JSON_EXTRACT(pro, '$.pro')),
                    part_id
                ); 
            SELECT j + 1 INTO j;
        END WHILE;

        -- insert means
        SET j = 0;
        WHILE j < JSON_LENGTH(means) DO
            SELECT JSON_EXTRACT(means, CONCAT('$[', j, ']')) INTO mean;
            INSERT INTO
                means(part_id, meaning, level)
            VALUES (
                part_id,
                JSON_UNQUOTE(JSON_EXTRACT(mean, '$.meaning')),
                JSON_UNQUOTE(JSON_EXTRACT(mean, '$.level'))
            ); 
            SET mean_id = LAST_INSERT_ID();
            -- insert examples
            SET k = 0;

            SET examples = JSON_EXTRACT(mean, '$.examples');

            WHILE k < JSON_LENGTH(examples) DO
                SELECT JSON_UNQUOTE(JSON_EXTRACT(examples, CONCAT('$[', k, ']'))) INTO example;
                
                INSERT INTO
                    examples(mean_id, example)
                VALUES
                    (mean_id, example); 
                
                SELECT k + 1 INTO k;
            END WHILE;
            SELECT j + 1 INTO j;
        END WHILE;
    SELECT
        i + 1 INTO i;
    END WHILE;

    COMMIT;
END |

DELIMITER ;

DELIMITER | 
CREATE FUNCTION GetVobsRandom (
    collection_id_p INT,
    lm INT,
    recent_day_num_p INT
) 
RETURNS JSON
DETERMINISTIC
BEGIN
    DECLARE vobs JSON;

    SELECT JSON_ARRAYAGG(JSON_OBJECT(
            'parts', parts,
            'word', v.word,
            'id', v.id
        )) INTO vobs
    FROM collections c
    LEFT JOIN collection_words cw ON c.id = cw.collection_id
    JOIN (
        SELECT id, word 
        FROM vobs 
        WHERE 
            (recent_day_num_p != -1 AND created_at >= CONCAT(DATE_SUB(CURDATE(), INTERVAL recent_day_num_p DAY), ' 00:00:00'))
            OR recent_day_num_p = -1
        ORDER BY RAND() 
        LIMIT lm
    )AS v ON v.id = cw.vob_id 
    LEFT JOIN (
        SELECT vob_id, JSON_ARRAYAGG(
            JSON_OBJECT(
                'pronounces', pronounces,
                'means', means,
                'type', vob_parts.type,
                'title', vob_parts.title
            )
        ) AS parts
        FROM vob_parts
        LEFT JOIN (
            SELECT part_id, JSON_ARRAYAGG(
                JSON_OBJECT(
                    'audio_src', audio_src,
                    'local_file', local_file,
                    'region', region,
                    'pro', pro
                )
            ) AS pronounces
            FROM pronounces
            GROUP BY part_id
        ) AS audio_data ON vob_parts.id = audio_data.part_id
        LEFT JOIN (
            SELECT part_id, JSON_ARRAYAGG(
                JSON_OBJECT(
                    'meaning', meaning,
                    'level', level,
                    'examples', examples
                )
            ) AS means
            FROM (
                SELECT means.part_id, meaning, level, JSON_ARRAYAGG(example) AS examples
                FROM means
                LEFT JOIN examples ON means.id = examples.mean_id
                GROUP BY means.part_id, meaning, level
            ) AS mean_data
            GROUP BY part_id
        ) AS mean_data ON vob_parts.id = mean_data.part_id 
        GROUP BY vob_parts.vob_id
    ) AS part_data ON v.id = part_data.vob_id 
    GROUP BY c.id;
    RETURN vobs;
END |
DELIMITER ;

DELIMITER | 
CREATE FUNCTION GetVobDict (
    word VARCHAR(255)
) 
RETURNS JSON
DETERMINISTIC
BEGIN
    DECLARE vobs JSON;

    SELECT JSON_OBJECT(
            'parts', parts,
            'word', v.word,
            'id', v.id
        ) INTO vobs
    FROM vobs v
    LEFT JOIN (
        SELECT vob_id, JSON_ARRAYAGG(
            JSON_OBJECT(
                'pronounces', pronounces,
                'means', means,
                'type', vob_parts.type,
                'title', vob_parts.title
            )
        ) AS parts
        FROM vob_parts
        LEFT JOIN (
            SELECT part_id, JSON_ARRAYAGG(
                JSON_OBJECT(
                    'audio_src', audio_src,
                    'local_file', local_file,
                    'region', region,
                    'pro', pro
                )
            ) AS pronounces
            FROM pronounces
            GROUP BY part_id
        ) AS audio_data ON vob_parts.id = audio_data.part_id
        LEFT JOIN (
            SELECT part_id, JSON_ARRAYAGG(
                JSON_OBJECT(
                    'meaning', meaning,
                    'level', level,
                    'examples', examples
                )
            ) AS means
            FROM (
                SELECT means.part_id, meaning, level, JSON_ARRAYAGG(example) AS examples
                FROM means
                LEFT JOIN examples ON means.id = examples.mean_id
                GROUP BY means.part_id, meaning, level
            ) AS mean_data
            GROUP BY part_id
        ) AS mean_data ON vob_parts.id = mean_data.part_id 
        GROUP BY vob_parts.vob_id
    ) AS part_data ON v.id = part_data.vob_id
    WHERE v.word = word;
    RETURN vobs;
END |
DELIMITER ;
