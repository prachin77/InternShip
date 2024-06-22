create DATABASE IF NOT EXISTS task6;

use task6;

CREATE TABLE UserInfo(
    username varchar(20) ,
    email VARCHAR(20),
    password varchar(20),
    userid varchar(50) NOT NULL
);