--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2 (Debian 14.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: books; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.books (
    site character varying(15),
    id integer,
    hash_code integer,
    title text,
    writer_id integer,
    type character varying(20),
    update_date character varying(30),
    update_chapter text,
    status character varying(10),
    is_downloaded boolean DEFAULT false,
    checksum text,
    writer_checksum text
);


ALTER TABLE public.books OWNER TO test;

--
-- Name: errors; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.errors (
    site character varying(15),
    id integer,
    data text
);


ALTER TABLE public.errors OWNER TO test;

--
-- Name: writers; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.writers (
    id integer NOT NULL,
    name text,
    checksum text
);


ALTER TABLE public.writers OWNER TO test;

--
-- Name: writers_id_seq; Type: SEQUENCE; Schema: public; Owner: test
--

CREATE SEQUENCE public.writers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.writers_id_seq OWNER TO test;

--
-- Name: writers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: test
--

ALTER SEQUENCE public.writers_id_seq OWNED BY public.writers.id;


--
-- Name: writers id; Type: DEFAULT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.writers ALTER COLUMN id SET DEFAULT nextval('public.writers_id_seq'::regclass);


--
-- Name: writers writers_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.writers
    ADD CONSTRAINT writers_pkey PRIMARY KEY (id);


--
-- Name: books_index; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX books_index ON public.books USING btree (site, id, hash_code);


--
-- Name: books_status; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX books_status ON public.books USING btree (status);


--
-- Name: books_title; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX books_title ON public.books USING btree (title);


--
-- Name: books_writer; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX books_writer ON public.books USING btree (writer_id);


--
-- Name: books_writer_checksum; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX books_writer_checksum ON public.books USING btree (writer_checksum);


--
-- Name: checksum_index; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX checksum_index ON public.books USING btree (checksum);


--
-- Name: errors_index; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX errors_index ON public.errors USING btree (site, id);


--
-- Name: writers_checksum_index; Type: INDEX; Schema: public; Owner: test
--

CREATE INDEX writers_checksum_index ON public.writers USING btree (checksum);


--
-- Name: writers_id; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX writers_id ON public.writers USING btree (id);


--
-- Name: writers_name; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX writers_name ON public.writers USING btree (name);


--
-- PostgreSQL database dump complete
--

