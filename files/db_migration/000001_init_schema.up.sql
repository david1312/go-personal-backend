BEGIN;

create table customers
(
    id              bigint primary key AUTO_INCREMENT,
    uid             varchar(36) not null UNIQUE,
    name            varchar(100) not null,
    password        varchar(200) not null,
    email           varchar(200) not null,
    email_verified_token varchar(64),
    email_verified_at timestamp NULL DEFAULT NULL,
    email_verified_sent   tinyint NOT NULL default 0,
    email_change_code  varchar(6) NULL DEFAULT null,
    email_change_eligible      boolean not null default false,
    gender          enum('LAKI-LAKI', 'PEREMPUAN'),
    is_active       boolean not null default true,
    phone           varchar(20),
    phone_verified_at timestamp NULL DEFAULT NULL,
    avatar          varchar(100),
    birthdate       DATE NULL DEFAULT NULL,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      timestamp NULL DEFAULT NULL
);

create table carts
(
    id              bigint primary key AUTO_INCREMENT,
    customer_uid             varchar(36) not null UNIQUE,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      timestamp NULL DEFAULT NULL

);


create table outlets
(
    id              bigint primary key AUTO_INCREMENT,
    name             varchar(50) not null,
    city      varchar(50) not null,
    districts      varchar(50) not null,
    address      varchar(100) not null,
    latitude     varchar(20) not null,
    longitude varchar(20) not null,
    gmap_url varchar(50) not null
);



CREATE TABLE `payment_category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `payment_method` (
  `id` varchar(50) NOT NULL ,
  `id_payment_category` int NOT NULL ,
  `description` varchar(50) NOT NULL,
   is_default      boolean not null default false,
   icon     varchar(50) not null default 'default.png',
  PRIMARY KEY (`id`),
 foreign key (id_payment_category) references payment_category(id)
);

create table outlet_ratings
(
    id              bigint primary key AUTO_INCREMENT,
    customer_id              bigint not null,
    outlet_id              bigint not null,
    comment varchar(200),
    rating tinyint not null,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    
    foreign key (customer_id) references customers (id),
    foreign key (outlet_id) references outlets(id)
);

create table outlet_ratings_img
(
    id_ratings              bigint,
    image              varchar(50) not null,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    
    foreign key (id_ratings) references outlet_ratings (id)
);

create table merchants
(
    id              int primary key AUTO_INCREMENT,
    outlet_id              bigint not null,
    password        varchar(200) not null,
    email           varchar(200) null default null,
    phone varchar(250) NULL DEFAULT NULL,
    avatar          varchar(100) default 'default.png',
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      timestamp NULL DEFAULT NULL,
    foreign key (outlet_id) references outlets(id)
);


COMMIT;
