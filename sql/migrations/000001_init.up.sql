CREATE TABLE video (
    id varchar(255) NOT NULL PRIMARY KEY,
    resource_id varchar(255) NOT NULL,
    file_path varchar(255) NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE job (
    id varchar(255) NOT NULL PRIMARY KEY,
    output_bucket_path varchar(255) NOT NULL,
    status varchar(255) NOT NULL,
    video_id varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP(),
    FOREIGN KEY(video_id) REFERENCES video(id)
);