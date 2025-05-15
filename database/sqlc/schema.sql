--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: books; Type: TABLE; Schema: public; Owner: book_spider
--

CREATE TABLE public.books (
    site character varying(15) NOT NULL,
    id integer NOT NULL,
    hash_code integer NOT NULL,
    title text,
    writer_id integer,
    type character varying(20),
    update_date character varying(30),
    update_chapter text,
    status character varying(10) NOT NULL,
    is_downloaded boolean DEFAULT false NOT NULL,
    checksum text,
    writer_checksum text
);


ALTER TABLE public.books OWNER TO book_spider;

--
-- Name: errors; Type: TABLE; Schema: public; Owner: book_spider
--

CREATE TABLE public.errors (
    site character varying(15),
    id integer,
    data text
);


ALTER TABLE public.errors OWNER TO book_spider;

--
-- Name: writers; Type: TABLE; Schema: public; Owner: book_spider
--

CREATE TABLE public.writers (
    id integer NOT NULL,
    name text,
    checksum text
);


ALTER TABLE public.writers OWNER TO book_spider;

--
-- Name: writers_id_seq; Type: SEQUENCE; Schema: public; Owner: book_spider
--

CREATE SEQUENCE public.writers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.writers_id_seq OWNER TO book_spider;

--
-- Name: writers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: book_spider
--

ALTER SEQUENCE public.writers_id_seq OWNED BY public.writers.id;


--
-- Name: writers id; Type: DEFAULT; Schema: public; Owner: book_spider
--

ALTER TABLE ONLY public.writers ALTER COLUMN id SET DEFAULT nextval('public.writers_id_seq'::regclass);


--
-- Name: writers writers_pkey; Type: CONSTRAINT; Schema: public; Owner: book_spider
--

ALTER TABLE ONLY public.writers
    ADD CONSTRAINT writers_pkey PRIMARY KEY (id);


--
-- Name: books__checksum; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books__checksum ON public.books USING btree (checksum, writer_checksum);


--
-- Name: books__is_downloaded; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books__is_downloaded ON public.books USING btree (site, is_downloaded);


--
-- Name: books__status; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books__status ON public.books USING btree (status, is_downloaded);


--
-- Name: books__vendor_reference; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books__vendor_reference ON public.books USING btree (site, id, hash_code DESC);


--
-- Name: books_index; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE UNIQUE INDEX books_index ON public.books USING btree (site, id, hash_code);


--
-- Name: books_status; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books_status ON public.books USING btree (status);


--
-- Name: books_title; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books_title ON public.books USING btree (title);


--
-- Name: books_writer; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books_writer ON public.books USING btree (writer_id);


--
-- Name: books_writer_checksum; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX books_writer_checksum ON public.books USING btree (writer_checksum);


--
-- Name: checksum_index; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX checksum_index ON public.books USING btree (checksum);


--
-- Name: errors_index; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE UNIQUE INDEX errors_index ON public.errors USING btree (site, id);


--
-- Name: writers__name; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX writers__name ON public.writers USING btree (name);


--
-- Name: writers_checksum_index; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE INDEX writers_checksum_index ON public.writers USING btree (checksum);


--
-- Name: writers_id; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE UNIQUE INDEX writers_id ON public.writers USING btree (id);


--
-- Name: writers_name; Type: INDEX; Schema: public; Owner: book_spider
--

CREATE UNIQUE INDEX writers_name ON public.writers USING btree (name);


--
-- PostgreSQL database dump complete
--

