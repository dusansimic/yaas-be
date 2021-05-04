-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.3
-- PostgreSQL version: 13.0
-- Project Site: pgmodeler.io
-- Model Author: Dušan Simić <dusan.simic1810@gmail.com>
-- object: yaas | type: ROLE --
-- DROP ROLE IF EXISTS yaas;
-- CREATE ROLE yaas WITH 
	-- UNENCRYPTED PASSWORD 'yaaspass';
-- ddl-end --


-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: yaas | type: DATABASE --
-- DROP DATABASE IF EXISTS yaas;
-- CREATE DATABASE yaas
	-- OWNER = yaas;
-- ddl-end --


-- object: yaas | type: SCHEMA --
-- DROP SCHEMA IF EXISTS yaas CASCADE;
-- CREATE SCHEMA yaas;
-- ddl-end --
-- ALTER SCHEMA yaas OWNER TO yaas;
-- ddl-end --

-- SET search_path TO pg_catalog,public,yaas;
-- ddl-end --

-- object: public.domain | type: TABLE --
-- DROP TABLE IF EXISTS public.domain CASCADE;
CREATE TABLE public.domain (
	idd serial NOT NULL,
	idu integer NOT NULL,
	code varchar(22) NOT NULL,
	name varchar(256) NOT NULL,
	description text,
	CONSTRAINT "Domain_Code_uq" UNIQUE (code),
	CONSTRAINT "Domain_Name_uq" UNIQUE (code,name),
	CONSTRAINT "Domain_pk" PRIMARY KEY (idd)

);
-- ddl-end --
ALTER TABLE public.domain OWNER TO yaas;
-- ddl-end --

-- object: public.people | type: TABLE --
-- DROP TABLE IF EXISTS public.people CASCADE;
CREATE TABLE public.people (
	idu serial NOT NULL,
	username varchar(32) NOT NULL,
	password_hash bytea NOT NULL,
	password_salt bytea NOT NULL,
	CONSTRAINT "Username_uq" UNIQUE (username),
	CONSTRAINT "User_pk" PRIMARY KEY (idu)

);
-- ddl-end --
ALTER TABLE public.people OWNER TO yaas;
-- ddl-end --

-- object: public.record | type: TABLE --
-- DROP TABLE IF EXISTS public.record CASCADE;
CREATE TABLE public.record (
	idr serial NOT NULL,
	idd integer NOT NULL,
	idrf integer NOT NULL,
	idc integer NOT NULL,
	idb integer NOT NULL,
	idos integer NOT NULL,
	iddv integer NOT NULL,
	"timestamp" timestamptz NOT NULL,
	url text NOT NULL,
	path text NOT NULL,
	CONSTRAINT "Record_pk" PRIMARY KEY (idr)

);
-- ddl-end --
ALTER TABLE public.record OWNER TO yaas;
-- ddl-end --

-- object: public.referrer | type: TABLE --
-- DROP TABLE IF EXISTS public.referrer CASCADE;
CREATE TABLE public.referrer (
	idrf serial NOT NULL,
	name varchar(45) NOT NULL,
	domain varchar(256) NOT NULL,
	CONSTRAINT "Referrer_pk" PRIMARY KEY (idrf),
	CONSTRAINT "Referrer_Name_uq" UNIQUE (name),
	CONSTRAINT "Referrer_Domain_uq" UNIQUE (domain)

);
-- ddl-end --
ALTER TABLE public.referrer OWNER TO yaas;
-- ddl-end --

-- object: public.device | type: TABLE --
-- DROP TABLE IF EXISTS public.device CASCADE;
CREATE TABLE public.device (
	iddv serial NOT NULL,
	name varchar(45) NOT NULL,
	CONSTRAINT "Device_Name_uq" UNIQUE (name),
	CONSTRAINT "Device_pk" PRIMARY KEY (iddv)

);
-- ddl-end --
ALTER TABLE public.device OWNER TO yaas;
-- ddl-end --

-- object: public.country | type: TABLE --
-- DROP TABLE IF EXISTS public.country CASCADE;
CREATE TABLE public.country (
	idc serial NOT NULL,
	name varchar(45) NOT NULL,
	iso varchar(2) NOT NULL,
	iso3 varchar(3) NOT NULL,
	numcode integer NOT NULL,
	phoneprefix integer NOT NULL,
	CONSTRAINT "Country_ISO_uq" UNIQUE (iso),
	CONSTRAINT "Country_pk" PRIMARY KEY (idc),
	CONSTRAINT "Country_ISO3_uq" UNIQUE (iso3),
	CONSTRAINT "Country_Num_Code_uq" UNIQUE (numcode),
	CONSTRAINT "Country_Name_uq" UNIQUE (name)

);
-- ddl-end --
ALTER TABLE public.country OWNER TO yaas;
-- ddl-end --

-- object: public.browser | type: TABLE --
-- DROP TABLE IF EXISTS public.browser CASCADE;
CREATE TABLE public.browser (
	idb serial NOT NULL,
	name varchar(45) NOT NULL,
	CONSTRAINT "Browser_Name_uq" UNIQUE (name),
	CONSTRAINT "Browser_pk" PRIMARY KEY (idb)

);
-- ddl-end --
ALTER TABLE public.browser OWNER TO yaas;
-- ddl-end --

-- object: public.operatingsystem | type: TABLE --
-- DROP TABLE IF EXISTS public.operatingsystem CASCADE;
CREATE TABLE public.operatingsystem (
	idos serial NOT NULL,
	name varchar(45) NOT NULL,
	CONSTRAINT "OperatingSystem_pk" PRIMARY KEY (idos),
	CONSTRAINT "OperatingSystem_Name_uq" UNIQUE (name)

);
-- ddl-end --
ALTER TABLE public.operatingsystem OWNER TO yaas;
-- ddl-end --

-- object: domain_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS domain_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT domain_fk FOREIGN KEY (idd)
REFERENCES public.domain (idd) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: referrer_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS referrer_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT referrer_fk FOREIGN KEY (idrf)
REFERENCES public.referrer (idrf) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: country_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS country_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT country_fk FOREIGN KEY (idc)
REFERENCES public.country (idc) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: browser_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS browser_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT browser_fk FOREIGN KEY (idb)
REFERENCES public.browser (idb) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: operatingsystem_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS operatingsystem_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT operatingsystem_fk FOREIGN KEY (idos)
REFERENCES public.operatingsystem (idos) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: device_fk | type: CONSTRAINT --
-- ALTER TABLE public.record DROP CONSTRAINT IF EXISTS device_fk CASCADE;
ALTER TABLE public.record ADD CONSTRAINT device_fk FOREIGN KEY (iddv)
REFERENCES public.device (iddv) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.domain DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.domain ADD CONSTRAINT people_fk FOREIGN KEY (idu)
REFERENCES public.people (idu) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --


