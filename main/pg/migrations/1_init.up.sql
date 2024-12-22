CREATE TABLE clicks (
    banner_id INT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    click_count INT NOT NULL,
    PRIMARY KEY (banner_id, timestamp)
);

