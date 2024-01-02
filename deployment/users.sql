CREATE TABLE users (
   id uuid DEFAULT gen_random_uuid() ,
   email character varying(255),
   first_name character varying(255),
   last_name character varying(255),
   password character varying(60),
   active BOOLEAN,
   created_at timestamp without time zone,
   updated_at timestamp without time zone
);
