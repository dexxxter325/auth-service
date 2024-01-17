CREATE TABLE products (
                            id serial PRIMARY KEY,
                           name varchar(255) NOT NULL UNIQUE,
                           description varchar(255) NOT NULL UNIQUE

);

CREATE TABLE users(
    id serial PRIMARY KEY,
    name varchar(255) not null unique,
    username varchar(255) not null unique,
    password varchar(255) not null

);

/*serial-уникальный id к каждому запросу
  varchar-строка,не превышающая 255 символов
  unique-знач.в столбцах уникальные(без повторов)*/
