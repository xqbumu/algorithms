drop TABLE IF EXISTS tbl_demo;

create TABLE IF NOT EXISTS tbl_demo(  
    ID INTEGER NOT NULL PRIMARY KEY,
    NAME TEXT
);

insert INTO tbl_demo (ID, NAME) VALUES (1, 'Alpha');

insert INTO tbl_demo (ID, NAME) VALUES (2, 'Beta');

insert INTO tbl_demo (ID, NAME) VALUES (3, 'Gamma');


-- @label: data
SELECT * FROM tbl_demo WHERE ID > ?low? AND ID < ?high?;
