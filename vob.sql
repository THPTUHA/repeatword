USE repeatword;

-- DROP SCHEMA repeatword;
-- CREATE SCHEMA repeatword;
-- CREATE TABLE vobs (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     word VARCHAR(255)
-- );

-- CREATE TABLE means (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     vob_id INT,
--     meaning VARCHAR(255),
--     type VARCHAR(255),
--     FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE
-- );

-- CREATE TABLE examples (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     mean_id INT,
--     example TEXT,
--     FOREIGN KEY (mean_id) REFERENCES means(id) ON DELETE CASCADE
-- );

-- CREATE TABLE collections (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255)
-- );

-- CREATE TABLE collection_words (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     vob_id INT,
--     collection_id INT,
--     FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE,
--     FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
-- );

-- ### Add word
-- INSERT INTO collections (name) VALUES ("s1");
-- INSERT INTO vobs (word) VALUES ("crew");
-- SET @vob_id = LAST_INSERT_ID();
-- INSERT INTO means (vob_id, meaning) VALUES (@vob_id, "a group of people who work together, especially all those who work on and operate a ship, aircraft, etc.");
-- SET @mean_id = LAST_INSERT_ID();
-- INSERT INTO examples (mean_id, example) VALUES 
--     (@mean_id, "an ambulance crew"),
--     (@mean_id, "a film crew");
-- INSERT INTO collection_words (vob_id, collection_id) VALUES 
--     (@vob_id, 1);
-- END
-- SELECT v.word, m.meaning, e.example FROM vobs v, means m, examples e
-- WHERE v.id = m.vob_id AND m.id = e.mean_id;
SELECT * FROM means;
SELECT * FROM examples;
-- DELETE FROM vobs;