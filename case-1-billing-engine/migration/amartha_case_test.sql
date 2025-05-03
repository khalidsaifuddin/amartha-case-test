--
-- PostgreSQL database dump
--

-- Dumped from database version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.0

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: history_status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.history_status_enum AS ENUM (
    'pending',
    'paid',
    'canceled'
);


ALTER TYPE public.history_status_enum OWNER TO postgres;

--
-- Name: installment_status_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.installment_status_enum AS ENUM (
    'pending',
    'paid',
    'unpaid'
);


ALTER TYPE public.installment_status_enum OWNER TO postgres;

--
-- Name: term_unit_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.term_unit_enum AS ENUM (
    'week',
    'month'
);


ALTER TYPE public.term_unit_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: loan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.loan (
    id integer NOT NULL,
    serial uuid DEFAULT public.uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255),
    updated_at timestamp with time zone,
    updated_by character varying(255),
    deleted_at timestamp with time zone,
    deleted_by character varying(255),
    borrower_name character varying(255),
    loan_name character varying(255),
    description text,
    loan_amount numeric,
    interest_rate numeric,
    loan_term integer,
    loan_term_unit public.term_unit_enum,
    loan_total_amount numeric,
    loan_outstanding numeric,
    loan_principal_only numeric
);


ALTER TABLE public.loan OWNER TO postgres;

--
-- Name: loan_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.loan_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.loan_id_seq OWNER TO postgres;

--
-- Name: loan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.loan_id_seq OWNED BY public.loan.id;


--
-- Name: loan_repayment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.loan_repayment (
    id integer NOT NULL,
    serial uuid DEFAULT public.uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255),
    updated_at timestamp with time zone,
    updated_by character varying(255),
    deleted_at timestamp with time zone,
    deleted_by character varying(255),
    loan_serial uuid,
    installment_number integer,
    amount numeric,
    interest_amount numeric,
    total_amount numeric,
    status public.installment_status_enum,
    repayment_amount numeric,
    repayment_time timestamp with time zone,
    repayment_due_date date,
    loan_repayment_history_serial uuid
);


ALTER TABLE public.loan_repayment OWNER TO postgres;

--
-- Name: loan_installment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.loan_installment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.loan_installment_id_seq OWNER TO postgres;

--
-- Name: loan_installment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.loan_installment_id_seq OWNED BY public.loan_repayment.id;


--
-- Name: loan_repayment_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.loan_repayment_history (
    id integer NOT NULL,
    serial uuid DEFAULT public.uuid_generate_v4(),
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(255),
    updated_at timestamp with time zone,
    updated_by character varying(255),
    deleted_at timestamp with time zone,
    deleted_by character varying(255),
    loan_repayment_serial uuid,
    repayment_amount numeric,
    repayment_time timestamp with time zone,
    related_loan_repayment_serials uuid[],
    status public.history_status_enum,
    loan_serial uuid
);


ALTER TABLE public.loan_repayment_history OWNER TO postgres;

--
-- Name: loan_repayment_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.loan_repayment_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.loan_repayment_history_id_seq OWNER TO postgres;

--
-- Name: loan_repayment_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.loan_repayment_history_id_seq OWNED BY public.loan_repayment_history.id;


--
-- Name: loan id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan ALTER COLUMN id SET DEFAULT nextval('public.loan_id_seq'::regclass);


--
-- Name: loan_repayment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan_repayment ALTER COLUMN id SET DEFAULT nextval('public.loan_installment_id_seq'::regclass);


--
-- Name: loan_repayment_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan_repayment_history ALTER COLUMN id SET DEFAULT nextval('public.loan_repayment_history_id_seq'::regclass);


--
-- Data for Name: loan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.loan (id, serial, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, borrower_name, loan_name, description, loan_amount, interest_rate, loan_term, loan_term_unit, loan_total_amount, loan_outstanding, loan_principal_only) FROM stdin;
1	62efb4bf-19d5-4070-9e61-bc182cee7850	2025-05-02 20:10:10.273608+07	\N	2025-05-03 06:50:30.318948+07	\N	\N	\N	Khalid	Khalid's loan	-	5000000	10	50	week	5500000	4730000	4300000
\.


--
-- Data for Name: loan_repayment; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.loan_repayment (id, serial, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, loan_serial, installment_number, amount, interest_amount, total_amount, status, repayment_amount, repayment_time, repayment_due_date, loan_repayment_history_serial) FROM stdin;
2	916961f0-25f1-4702-9fe8-6e559e4a9729	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:49:16.273375+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	2	100000	10000	110000	paid	110000	2025-05-20 07:00:00+07	2025-05-12	9bdd8c74-4b54-458f-9bfb-29bd7471c48b
1	87706454-8b9b-4ab9-927a-0bf6babf4202	2025-05-02 20:21:30.061592+07	\N	2025-05-03 06:49:16.281+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	1	100000	10000	110000	paid	110000	2025-05-20 07:00:00+07	2025-05-05	9bdd8c74-4b54-458f-9bfb-29bd7471c48b
7	9c70bca7-853c-485d-a9e2-2fff0f15242d	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:50:30.217675+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	7	100000	10000	110000	paid	110000	2025-06-20 07:00:00+07	2025-06-16	687554b4-d110-47c1-8783-8f4cbbaae1eb
6	84835379-e5fb-4df0-b060-4f090f35323a	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:50:30.255036+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	6	100000	10000	110000	paid	110000	2025-06-20 07:00:00+07	2025-06-09	687554b4-d110-47c1-8783-8f4cbbaae1eb
5	f0c2482b-fae6-4dfc-9486-197ae3444e61	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:50:30.275723+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	5	100000	10000	110000	paid	110000	2025-06-20 07:00:00+07	2025-06-02	687554b4-d110-47c1-8783-8f4cbbaae1eb
4	154d433f-6a15-400a-8231-895f8a5ae8a0	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:50:30.31318+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	4	100000	10000	110000	paid	110000	2025-06-20 07:00:00+07	2025-05-26	687554b4-d110-47c1-8783-8f4cbbaae1eb
3	1e4e2133-6576-405b-9813-d2e5d7ee3bb5	2025-05-02 20:23:13.999551+07	\N	2025-05-03 06:49:16.259645+07	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	3	100000	10000	110000	paid	110000	2025-05-20 07:00:00+07	2025-05-19	9bdd8c74-4b54-458f-9bfb-29bd7471c48b
9	2bae948f-29b4-40ed-b561-a7dae0f4e4f6	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	9	100000	10000	110000	unpaid	\N	\N	2025-06-30	\N
8	9254ddc4-4c12-48a4-ad13-2290316618a5	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	8	100000	10000	110000	unpaid	\N	\N	2025-06-23	\N
12	b5e808b2-adda-457b-8c43-3f9a03e54b66	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	12	100000	10000	110000	unpaid	\N	\N	2025-07-21	\N
11	ddf0251e-d30d-4335-840f-053ed9405b7d	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	11	100000	10000	110000	unpaid	\N	\N	2025-07-14	\N
10	c7ca84cd-d27c-42bc-afb0-6a6aa145e4c9	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	10	100000	10000	110000	unpaid	\N	\N	2025-07-07	\N
13	9d3bd758-e874-43ef-836a-56cc1cd907c9	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	13	100000	10000	110000	unpaid	\N	\N	2025-07-28	\N
14	a50e5285-8a1f-4f0d-ac19-e184715f771f	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	14	100000	10000	110000	unpaid	\N	\N	2025-08-04	\N
15	646310fc-a0f8-469e-a7ac-6958148d1a7e	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	15	100000	10000	110000	unpaid	\N	\N	2025-08-11	\N
16	9b70d014-10c3-491d-b8fd-d5f64fc3d743	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	16	100000	10000	110000	unpaid	\N	\N	2025-08-18	\N
17	2619ec19-d7c0-4baf-96c6-5a1985be2b49	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	17	100000	10000	110000	unpaid	\N	\N	2025-08-25	\N
18	4e053ea7-f027-4db8-b483-a49e8d718500	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	18	100000	10000	110000	unpaid	\N	\N	2025-09-01	\N
19	38caa353-2b98-458a-843b-6ddcccadff17	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	19	100000	10000	110000	unpaid	\N	\N	2025-09-08	\N
20	a3f3f252-afee-4312-a6ef-e52ce7906c6f	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	20	100000	10000	110000	unpaid	\N	\N	2025-09-15	\N
21	7afb8d13-b280-43d8-ac5d-433270311413	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	21	100000	10000	110000	unpaid	\N	\N	2025-09-22	\N
22	b4510beb-de31-465b-abbf-527996e58664	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	22	100000	10000	110000	unpaid	\N	\N	2025-09-29	\N
23	f881c65c-e1a9-4d6b-bcd3-9feeeb8d9b79	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	23	100000	10000	110000	unpaid	\N	\N	2025-10-06	\N
24	d0200268-4b29-4c05-b450-8df53df33136	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	24	100000	10000	110000	unpaid	\N	\N	2025-10-13	\N
25	7a6e9e41-7678-46a5-9b2c-ebe32d020f5d	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	25	100000	10000	110000	unpaid	\N	\N	2025-10-20	\N
26	2e15aae8-9ffb-4d8b-87f2-f1071bb5e8f4	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	26	100000	10000	110000	unpaid	\N	\N	2025-10-27	\N
27	f6e0b764-35cf-49bc-b7c1-7782623d5168	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	27	100000	10000	110000	unpaid	\N	\N	2025-11-03	\N
28	dc851364-d562-491b-b8ac-5baf2f841a82	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	28	100000	10000	110000	unpaid	\N	\N	2025-11-10	\N
29	30ca61a0-a90d-4101-9cf7-43a8683d2b53	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	29	100000	10000	110000	unpaid	\N	\N	2025-11-17	\N
30	ad27f770-e1af-40da-bfc1-21aa9db6138c	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	30	100000	10000	110000	unpaid	\N	\N	2025-11-24	\N
31	fc4be224-3b47-4fbb-9de4-397cf070f954	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	31	100000	10000	110000	unpaid	\N	\N	2025-12-01	\N
32	8a3d80e0-7a28-4ba8-9e07-68a51b8c8fcd	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	32	100000	10000	110000	unpaid	\N	\N	2025-12-08	\N
33	6d211423-70e0-47e7-abaa-75b6e87c7517	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	33	100000	10000	110000	unpaid	\N	\N	2025-12-15	\N
34	547570ac-a5f4-4916-82d7-c72487fcf884	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	34	100000	10000	110000	unpaid	\N	\N	2025-12-22	\N
35	f42498ca-2b52-450f-9e41-17932cc5a5c3	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	35	100000	10000	110000	unpaid	\N	\N	2025-12-29	\N
36	319681f0-8f43-46af-a817-0a3673b77936	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	36	100000	10000	110000	unpaid	\N	\N	2026-01-05	\N
37	46bb48b8-0837-4e07-81a8-2a05bc4eca1f	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	37	100000	10000	110000	unpaid	\N	\N	2026-01-12	\N
38	931b030d-c3a3-45a3-8460-b74876cb526a	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	38	100000	10000	110000	unpaid	\N	\N	2026-01-19	\N
39	8bbff701-83e6-49de-aaa5-7142fdc32020	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	39	100000	10000	110000	unpaid	\N	\N	2026-01-26	\N
40	660e3187-b443-4ed7-811d-39caa8b31d12	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	40	100000	10000	110000	unpaid	\N	\N	2026-02-02	\N
41	af481308-6642-4c38-b3eb-c1b052ed4f66	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	41	100000	10000	110000	unpaid	\N	\N	2026-02-09	\N
42	1069863a-cd86-4e55-b004-0089f8b6d742	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	42	100000	10000	110000	unpaid	\N	\N	2026-02-16	\N
43	3a3e57a4-a5d6-40a7-937f-c6d5e1c1925a	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	43	100000	10000	110000	unpaid	\N	\N	2026-02-23	\N
44	f9ed7902-45c9-47e9-beae-850d5c715ead	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	44	100000	10000	110000	unpaid	\N	\N	2026-03-02	\N
45	438df63f-daa2-4526-b71d-f54186f987bf	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	45	100000	10000	110000	unpaid	\N	\N	2026-03-09	\N
46	13ace58e-9957-407a-b921-8682f430e226	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	46	100000	10000	110000	unpaid	\N	\N	2026-03-16	\N
47	aa8ad659-6a52-47a6-b411-563de68e8987	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	47	100000	10000	110000	unpaid	\N	\N	2026-03-23	\N
48	458d0fd3-b3aa-4f08-a1c3-ad0835c4e18e	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	48	100000	10000	110000	unpaid	\N	\N	2026-03-30	\N
49	39883f9c-9730-4b75-973e-d7d676af388e	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	49	100000	10000	110000	unpaid	\N	\N	2026-04-06	\N
50	b80a07ca-cbfd-4a1e-872d-e4b82928c58a	2025-05-02 20:23:13.999551+07	\N	\N	\N	\N	\N	62efb4bf-19d5-4070-9e61-bc182cee7850	50	100000	10000	110000	unpaid	\N	\N	2026-04-13	\N
\.


--
-- Data for Name: loan_repayment_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.loan_repayment_history (id, serial, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, loan_repayment_serial, repayment_amount, repayment_time, related_loan_repayment_serials, status, loan_serial) FROM stdin;
9	9bdd8c74-4b54-458f-9bfb-29bd7471c48b	2025-05-03 06:46:07.362203+07		2025-05-03 06:49:16.204314+07		\N	\N	1e4e2133-6576-405b-9813-d2e5d7ee3bb5	330000	2025-05-20 07:00:00+07	{1e4e2133-6576-405b-9813-d2e5d7ee3bb5,916961f0-25f1-4702-9fe8-6e559e4a9729,87706454-8b9b-4ab9-927a-0bf6babf4202}	paid	62efb4bf-19d5-4070-9e61-bc182cee7850
10	687554b4-d110-47c1-8783-8f4cbbaae1eb	2025-05-03 06:50:30.202267+07		2025-05-03 06:50:30.202267+07		\N	\N	9c70bca7-853c-485d-a9e2-2fff0f15242d	440000	2025-06-20 07:00:00+07	{9c70bca7-853c-485d-a9e2-2fff0f15242d,84835379-e5fb-4df0-b060-4f090f35323a,f0c2482b-fae6-4dfc-9486-197ae3444e61,154d433f-6a15-400a-8231-895f8a5ae8a0}	paid	62efb4bf-19d5-4070-9e61-bc182cee7850
\.


--
-- Name: loan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.loan_id_seq', 1, true);


--
-- Name: loan_installment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.loan_installment_id_seq', 50, true);


--
-- Name: loan_repayment_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.loan_repayment_history_id_seq', 10, true);


--
-- Name: loan_repayment loan_installment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan_repayment
    ADD CONSTRAINT loan_installment_pkey PRIMARY KEY (id);


--
-- Name: loan loan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan
    ADD CONSTRAINT loan_pkey PRIMARY KEY (id);


--
-- Name: loan_repayment_history loan_repayment_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan_repayment_history
    ADD CONSTRAINT loan_repayment_history_pkey PRIMARY KEY (id);


--
-- Name: loan_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX loan_unique ON public.loan USING btree (serial);


--
-- Name: loan_repayment loan_installment_loan_serial_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.loan_repayment
    ADD CONSTRAINT loan_installment_loan_serial_fkey FOREIGN KEY (loan_serial) REFERENCES public.loan(serial);


--
-- PostgreSQL database dump complete
--

