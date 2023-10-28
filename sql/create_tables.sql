CREATE TABLE school(
    idschool serial PRIMARY KEY,
    name TEXT,
    coords POINT
);

CREATE TABLE restaurant(
    idrestaurant serial PRIMARY KEY,
    url TEXT,
    name TEXT,
    gpscoord POINT
);

CREATE TABLE cat_suggestions(
    idcat serial PRIMARY KEY,
    namecat TEXT
);

CREATE TABLE suggestions_restaurant(
    idsuggestion serial PRIMARY KEY,
    keyword TEXT,
    idrestaurant INT,
    CONSTRAINT fk_idrestaurant_sr FOREIGN KEY (idrestaurant) REFERENCES restaurant(idrestaurant) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE uselessfoodname (
    idfoodname serial PRIMARY KEY,
    name TEXT
);

CREATE TABLE meal (
    idmeal serial PRIMARY KEY,
    typemeal TEXT,
    foodies JSONB,
    day DATE,
    idrestaurant INT,
    CONSTRAINT fk_idrestaurant_meal FOREIGN KEY  (idrestaurant) REFERENCES restaurant(idrestaurant) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE users (
    iduser serial PRIMARY KEY,
    mail TEXT,
    password TEXT,
    name TEXT,
    idschool INT,
    ical TEXT,
    nonce TEXT,
    name_modified TIMESTAMP,
    token TEXT,
    salt TEXT,
    CONSTRAINT fk_idschool_user FOREIGN KEY (idschool) REFERENCES school(idschool) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE federal_credentials(
    credential_id serial PRIMARY KEY,
    user_id INT,
    provider TEXT,
    created_at TIMESTAMP,
    custom_name TEXT,
    CONSTRAINT fk_iduser_fc FOREIGN KEY (user_id) REFERENCES users(iduser) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE favoriterestaurant (
    idrestaurant INT,
    iduser INT,
    CONSTRAINT ffk_iduser_fr FOREIGN KEY (iduser) REFERENCES users(iduser) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_idrestaurant_fr FOREIGN KEY (idrestaurant) REFERENCES restaurant(idrestaurant) ON DELETE CASCADE ON UPDATE CASCADE
);
