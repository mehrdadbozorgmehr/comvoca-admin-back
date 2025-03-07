CREATE DATABASE comvoca_stage;
CREATE DATABASE comvoca_dev;

CREATE DATABASE housing_comply_dev;
CREATE DATABASE housing_comply_stage;


SELECT datname FROM pg_database;

CREATE USER comvoca_stage_user WITH PASSWORD 'qT&(URF~Xp4r';
CREATE USER comvoca_dev_user WITH PASSWORD 'v+6H4af7t?Mv';
CREATE USER housing_comply_dev_user WITH PASSWORD 'WLHUzFQ=9$VZ';
CREATE USER housing_comply_stage_user WITH PASSWORD 'H?[)U\XE\]fg2';


GRANT ALL PRIVILEGES ON DATABASE comvoca_stage TO comvoca_stage_user;
GRANT ALL PRIVILEGES ON DATABASE comvoca_dev TO comvoca_dev_user;
GRANT ALL PRIVILEGES ON DATABASE housing_comply_dev TO housing_comply_dev_user;
GRANT ALL PRIVILEGES ON DATABASE housing_comply_stage TO housing_comply_stage_user;

ALTER DATABASE comvoca_stage OWNER TO comvoca_stage_user;
ALTER DATABASE comvoca_dev OWNER TO comvoca_dev_user;
ALTER DATABASE housing_comply_dev OWNER TO housing_comply_dev_user;
ALTER DATABASE housing_comply_stage OWNER TO housing_comply_stage_user;


-- Revoke access for comvoca_stage_user
REVOKE CONNECT ON DATABASE comvoca_dev FROM comvoca_stage_user;
REVOKE CONNECT ON DATABASE housing_comply_dev FROM comvoca_stage_user;
REVOKE CONNECT ON DATABASE housing_comply_stage FROM comvoca_stage_user;

-- Revoke access for comvoca_dev_user
REVOKE CONNECT ON DATABASE comvoca_stage FROM comvoca_dev_user;
REVOKE CONNECT ON DATABASE housing_comply_dev FROM comvoca_dev_user;
REVOKE CONNECT ON DATABASE housing_comply_stage FROM comvoca_dev_user;

-- Revoke access for housing_comply_dev_user
REVOKE CONNECT ON DATABASE comvoca_stage FROM housing_comply_dev_user;
REVOKE CONNECT ON DATABASE comvoca_dev FROM housing_comply_dev_user;
REVOKE CONNECT ON DATABASE housing_comply_stage FROM housing_comply_dev_user;

-- Revoke access for housing_comply_stage_user
REVOKE CONNECT ON DATABASE comvoca_stage FROM housing_comply_stage_user;
REVOKE CONNECT ON DATABASE comvoca_dev FROM housing_comply_stage_user;
REVOKE CONNECT ON DATABASE housing_comply_dev FROM housing_comply_stage_user;