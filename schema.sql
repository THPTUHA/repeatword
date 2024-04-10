CREATE TABLE vobs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    word VARCHAR(255)
);

CREATE TABLE means (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vob_id INT,
    meaning VARCHAR(255),
    FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE
);

CREATE TABLE examples (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mean_id INT,
    example TEXT,
    FOREIGN KEY (mean_id) REFERENCES means(id) ON DELETE CASCADE
);

CREATE TABLE collections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE collection_words (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vob_id INT,
    collection_id INT,
    FOREIGN KEY (vob_id) REFERENCES vobs(id) ON DELETE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);