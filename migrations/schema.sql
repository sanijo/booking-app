--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3 (Ubuntu 15.3-1.pgdg20.04+1)
-- Dumped by pg_dump version 15.3 (Ubuntu 15.3-1.pgdg20.04+1)

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
-- Name: models; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.models (
    id integer NOT NULL,
    model_name character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.models OWNER TO postgres;

--
-- Name: models_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.models_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.models_id_seq OWNER TO postgres;

--
-- Name: models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.models_id_seq OWNED BY public.models.id;


--
-- Name: rent; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rent (
    id integer NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    phone character varying(255) DEFAULT ''::character varying NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    model_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rent OWNER TO postgres;

--
-- Name: rent_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rent_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rent_id_seq OWNER TO postgres;

--
-- Name: rent_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rent_id_seq OWNED BY public.rent.id;


--
-- Name: rent_restrictions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rent_restrictions (
    id integer NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    model_id bigint NOT NULL,
    rent_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    restriction_id bigint NOT NULL
);


ALTER TABLE public.rent_restrictions OWNER TO postgres;

--
-- Name: rent_restrictions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rent_restrictions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rent_restrictions_id_seq OWNER TO postgres;

--
-- Name: rent_restrictions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rent_restrictions_id_seq OWNED BY public.rent_restrictions.id;


--
-- Name: restriction_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.restriction_types (
    id integer NOT NULL,
    restriction_name character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.restriction_types OWNER TO postgres;

--
-- Name: restriction_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.restriction_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.restriction_types_id_seq OWNER TO postgres;

--
-- Name: restriction_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.restriction_types_id_seq OWNED BY public.restriction_types.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255) DEFAULT ''::character varying NOT NULL,
    last_name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    access_level integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: models id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models ALTER COLUMN id SET DEFAULT nextval('public.models_id_seq'::regclass);


--
-- Name: rent id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent ALTER COLUMN id SET DEFAULT nextval('public.rent_id_seq'::regclass);


--
-- Name: rent_restrictions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent_restrictions ALTER COLUMN id SET DEFAULT nextval('public.rent_restrictions_id_seq'::regclass);


--
-- Name: restriction_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restriction_types ALTER COLUMN id SET DEFAULT nextval('public.restriction_types_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- Name: rent rent_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent
    ADD CONSTRAINT rent_pkey PRIMARY KEY (id);


--
-- Name: rent_restrictions rent_restrictions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent_restrictions
    ADD CONSTRAINT rent_restrictions_pkey PRIMARY KEY (id);


--
-- Name: restriction_types restriction_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restriction_types
    ADD CONSTRAINT restriction_types_pkey PRIMARY KEY (id);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: rent_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rent_email_idx ON public.rent USING btree (email);


--
-- Name: rent_last_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rent_last_name_idx ON public.rent USING btree (last_name);


--
-- Name: rent_restrictions_model_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rent_restrictions_model_id_idx ON public.rent_restrictions USING btree (model_id);


--
-- Name: rent_restrictions_rent_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rent_restrictions_rent_id_idx ON public.rent_restrictions USING btree (rent_id);


--
-- Name: rent_restrictions_start_date_end_date_model_id_rent_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rent_restrictions_start_date_end_date_model_id_rent_id_idx ON public.rent_restrictions USING btree (start_date, end_date, model_id, rent_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: rent rent_models_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent
    ADD CONSTRAINT rent_models_id_fk FOREIGN KEY (model_id) REFERENCES public.models(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rent_restrictions rent_restrictions_models_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent_restrictions
    ADD CONSTRAINT rent_restrictions_models_id_fk FOREIGN KEY (model_id) REFERENCES public.models(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rent_restrictions rent_restrictions_rent_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent_restrictions
    ADD CONSTRAINT rent_restrictions_rent_id_fk FOREIGN KEY (rent_id) REFERENCES public.rent(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rent_restrictions rent_restrictions_restriction_types_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rent_restrictions
    ADD CONSTRAINT rent_restrictions_restriction_types_id_fk FOREIGN KEY (restriction_id) REFERENCES public.restriction_types(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

