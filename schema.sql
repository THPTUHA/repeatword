CREATE SCHEMA IF NOT EXISTS repeatword;
USE repeatword;

CREATE TABLE IF NOT EXISTS vobs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    word VARCHAR(255)
);


ALTER TABLE vobs ADD CONSTRAINT unique_word UNIQUE (word);

CREATE TABLE IF NOT EXISTS vob_parts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vob_id INT,
    type VARCHAR(255),
    title VARCHAR(255),
    FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pronounces (
    id INT AUTO_INCREMENT PRIMARY KEY,
    audio_src VARCHAR(255),
    local_file VARCHAR(255),
    region VARCHAR(255),
    pro VARCHAR(255),
    part_id INT,
    FOREIGN KEY (part_id) REFERENCES vob_parts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS means (
    id INT AUTO_INCREMENT PRIMARY KEY,
    part_id INT,
    meaning VARCHAR(255),
    level VARCHAR(255),
    FOREIGN KEY (part_id) REFERENCES vob_parts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS examples (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mean_id INT,
    example TEXT,
    FOREIGN KEY (mean_id) REFERENCES means(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS collections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS collection_words (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vob_id INT,
    collection_id INT,
    FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);