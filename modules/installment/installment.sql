--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.2
-- Dumped by pg_dump version 9.6.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;


DROP TABLE IF EXISTS installment;
DROP TABLE IF EXISTS r_loan_installment;
DROP TABLE IF EXISTS loan;
DROP TABLE IF EXISTS r_loan_history;
DROP TABLE IF EXISTS loan_history;
DROP SEQUENCE IF EXISTS installment_id_seq;
DROP SEQUENCE IF EXISTS r_loan_installment_id_seq;
DROP SEQUENCE IF EXISTS loan_id_seq;
DROP SEQUENCE IF EXISTS r_loan_history_id_seq;
DROP SEQUENCE IF EXISTS loan_history_id_seq;

--
-- Name: installment; Type: TABLE; Schema: public; Owner: mis_amartha
--

CREATE TABLE installment (
    id integer NOT NULL,
    type text,
    presence text,
    "paidInstallment" numeric,
    penalty numeric,
    reserve numeric,
    frequency integer,
    stage text,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    "transactionDate" timestamp with time zone
);



--
-- Name: installment_id_seq; Type: SEQUENCE; Schema: public; Owner: mis_amartha
--

CREATE SEQUENCE installment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- Name: installment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mis_amartha
--

ALTER SEQUENCE installment_id_seq OWNED BY installment.id;


--
-- Name: loan; Type: TABLE; Schema: public; Owner: mis_amartha
--

CREATE TABLE loan (
    id integer NOT NULL,
    "loanPeriod" integer,
    subgroup character varying(30),
    purpose text,
    "urlPic1" text,
    "urlPic2" text,
    "submittedLoanDate" timestamp with time zone,
    "submittedPlafond" numeric,
    "submittedTenor" integer,
    "submittedInstallment" numeric,
    "creditScoreGrade" text,
    "creditScoreValue" numeric,
    tenor integer,
    rate numeric,
    installment numeric,
    plafond numeric,
    stage text,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone,
    "agreementType" character varying(5),
    "groupReserve" numeric,
    "isLWK" boolean,
    "isUPK" boolean,
    "isOld" boolean
);



--
-- Name: loan_id_seq; Type: SEQUENCE; Schema: public; Owner: mis_amartha
--

CREATE SEQUENCE loan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- Name: loan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mis_amartha
--

ALTER SEQUENCE loan_id_seq OWNED BY loan.id;


--
-- Name: r_loan_installment; Type: TABLE; Schema: public; Owner: mis_amartha
--

CREATE TABLE r_loan_installment (
    id integer NOT NULL,
    "loanId" integer,
    "installmentId" integer,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone
);



--
-- Name: r_loan_installment_id_seq; Type: SEQUENCE; Schema: public; Owner: mis_amartha
--

CREATE SEQUENCE r_loan_installment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



--
-- Name: r_loan_installment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mis_amartha
--

ALTER SEQUENCE r_loan_installment_id_seq OWNED BY r_loan_installment.id;


--
-- Name: installment id; Type: DEFAULT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY installment ALTER COLUMN id SET DEFAULT nextval('installment_id_seq'::regclass);


--
-- Name: loan id; Type: DEFAULT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY loan ALTER COLUMN id SET DEFAULT nextval('loan_id_seq'::regclass);


--
-- Name: r_loan_installment id; Type: DEFAULT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY r_loan_installment ALTER COLUMN id SET DEFAULT nextval('r_loan_installment_id_seq'::regclass);


--
-- Data for Name: installment; Type: TABLE DATA; Schema: public; Owner: mis_amartha
--

COPY installment (id, type, presence, "paidInstallment", penalty, reserve, frequency, stage, "createdAt", "updatedAt", "deletedAt", "transactionDate") FROM stdin;
1	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
2	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
3	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
4	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
5	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
6	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
7	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
8	NORMAL	ALFA	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
9	NORMAL	SAKIT	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
10	NORMAL	ALFA	78000.0	0	3000	1	SUCCESS	2017-04-03 16:45:05+07	2017-04-03 16:45:05+07	\N	2017-04-03 16:45:05+07
11	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
12	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
13	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
14	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
15	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
16	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
17	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
18	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
19	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
20	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
21	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
22	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
23	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
24	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
25	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
26	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
27	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
28	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
29	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
30	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
31	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
32	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
33	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
34	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
35	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
36	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
37	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
38	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
39	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
40	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
41	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
42	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
43	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
44	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
45	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
46	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
47	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
48	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
49	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
50	NORMAL	HADIR	78000.0	0	3000	1	APPROVE	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
51	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
52	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
53	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
54	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
55	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
56	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
57	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
58	NORMAL	ALFA	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
59	NORMAL	SAKIT	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
60	NORMAL	ALFA	78000.0	0	3000	1	SUCCESS	2017-04-03 16:45:05+07	2017-04-03 16:45:05+07	\N	2017-04-03 16:45:05+07
61	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
62	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
63	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
64	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
65	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
66	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
67	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
68	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
69	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
70	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
71	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
72	NORMAL	HADIR	78000.0	0	3000	1	SUCCESS	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
73	NORMAL	HADIR	78000.0	0	3000	28	APPROVE	2017-03-30 18:54:36+07	2015-01-01 12:34:45+07	\N	2017-03-30 00:00:00+07
\.


--
-- Name: installment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mis_amartha
--

SELECT pg_catalog.setval('installment_id_seq', 73, true);


--
-- Data for Name: loan; Type: TABLE DATA; Schema: public; Owner: mis_amartha
--

COPY loan (id, "loanPeriod", subgroup, purpose, "urlPic1", "urlPic2", "submittedLoanDate", "submittedPlafond", "submittedTenor", "submittedInstallment", "creditScoreGrade", "creditScoreValue", tenor, rate, installment, plafond, stage, "createdAt", "updatedAt", "deletedAt", "agreementType", "groupReserve", "isLWK", "isUPK", "isOld") FROM stdin;
1	1		modal kreditan sepre			2017-04-24 00:00:00+07	3000000	50	81000	B	93.76305941397311	50	0.3	78000	3000000	INSTALLMENT	2017-05-04 15:21:25.996749+07	2017-05-04 15:21:25.996749+07	\N		0	\N	\N	\N
2	1		modal kreditan sepre			2017-04-24 00:00:00+07	3000000	50	81000	B	93.76305941397311	50	0.3	78000	3000000	INSTALLMENT	2017-05-04 15:21:25.996749+07	2017-05-04 15:21:25.996749+07	\N		0	\N	\N	\N
\.


--
-- Name: loan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mis_amartha
--

SELECT pg_catalog.setval('loan_id_seq', 3, true);


--
-- Data for Name: r_loan_installment; Type: TABLE DATA; Schema: public; Owner: mis_amartha
--

COPY r_loan_installment (id, "loanId", "installmentId", "createdAt", "updatedAt", "deletedAt") FROM stdin;
1	1	1	2017-04-01 13:49:30+07	\N	\N
2	1	2	2017-04-01 13:49:30+07	\N	\N
3	1	3	2017-04-01 13:49:30+07	\N	\N
4	1	4	2017-04-01 13:49:30+07	\N	\N
5	1	5	2017-04-01 13:49:30+07	\N	\N
6	1	6	2017-04-01 13:49:30+07	\N	\N
7	1	7	2017-04-01 13:49:30+07	\N	\N
8	1	8	2017-04-01 13:49:30+07	\N	\N
9	1	9	2017-04-01 13:49:30+07	\N	\N
10	1	10	2017-04-01 13:49:30+07	\N	\N
11	1	11	2017-04-01 13:49:30+07	\N	\N
12	1	12	2017-04-01 13:49:30+07	\N	\N
13	1	13	2017-04-01 13:49:30+07	\N	\N
14	1	14	2017-04-01 13:49:30+07	\N	\N
15	1	15	2017-04-01 13:49:30+07	\N	\N
16	1	16	2017-04-01 13:49:30+07	\N	\N
17	1	17	2017-04-01 13:49:30+07	\N	\N
18	1	18	2017-04-01 13:49:30+07	\N	\N
19	1	19	2017-04-01 13:49:30+07	\N	\N
20	1	20	2017-04-01 13:49:30+07	\N	\N
21	1	21	2017-04-01 13:49:30+07	\N	\N
22	1	22	2017-04-01 13:49:30+07	\N	\N
23	1	23	2017-04-01 13:49:30+07	\N	\N
24	1	24	2017-04-01 13:49:30+07	\N	\N
25	1	25	2017-04-01 13:49:30+07	\N	\N
26	1	26	2017-04-01 13:49:30+07	\N	\N
27	1	27	2017-04-01 13:49:30+07	\N	\N
28	1	28	2017-04-01 13:49:30+07	\N	\N
29	1	29	2017-04-01 13:49:30+07	\N	\N
30	1	30	2017-04-01 13:49:30+07	\N	\N
31	1	31	2017-04-01 13:49:30+07	\N	\N
32	1	32	2017-04-01 13:49:30+07	\N	\N
33	1	33	2017-04-01 13:49:30+07	\N	\N
34	1	34	2017-04-01 13:49:30+07	\N	\N
35	1	35	2017-04-01 13:49:30+07	\N	\N
36	1	36	2017-04-01 13:49:30+07	\N	\N
37	1	37	2017-04-01 13:49:30+07	\N	\N
38	1	38	2017-04-01 13:49:30+07	\N	\N
39	1	39	2017-04-01 13:49:30+07	\N	\N
40	1	40	2017-04-01 13:49:30+07	\N	\N
41	1	41	2017-04-01 13:49:30+07	\N	\N
42	1	42	2017-04-01 13:49:30+07	\N	\N
43	1	43	2017-04-01 13:49:30+07	\N	\N
44	1	44	2017-04-01 13:49:30+07	\N	\N
45	1	45	2017-04-01 13:49:30+07	\N	\N
46	1	46	2017-04-01 13:49:30+07	\N	\N
47	1	47	2017-04-01 13:49:30+07	\N	\N
48	1	48	2017-04-01 13:49:30+07	\N	\N
49	1	49	2017-04-01 13:49:30+07	\N	\N
50	1	50	2017-04-01 13:49:30+07	\N	\N
51	2	51	2017-04-01 13:49:30+07	\N	\N
52	2	52	2017-04-01 13:49:30+07	\N	\N
53	2	53	2017-04-01 13:49:30+07	\N	\N
54	2	54	2017-04-01 13:49:30+07	\N	\N
55	2	55	2017-04-01 13:49:30+07	\N	\N
56	2	56	2017-04-01 13:49:30+07	\N	\N
57	2	57	2017-04-01 13:49:30+07	\N	\N
58	2	58	2017-04-01 13:49:30+07	\N	\N
59	2	59	2017-04-01 13:49:30+07	\N	\N
60	2	60	2017-04-01 13:49:30+07	\N	\N
61	2	61	2017-04-01 13:49:30+07	\N	\N
62	2	62	2017-04-01 13:49:30+07	\N	\N
63	2	63	2017-04-01 13:49:30+07	\N	\N
64	2	64	2017-04-01 13:49:30+07	\N	\N
65	2	65	2017-04-01 13:49:30+07	\N	\N
66	2	66	2017-04-01 13:49:30+07	\N	\N
67	2	67	2017-04-01 13:49:30+07	\N	\N
68	2	68	2017-04-01 13:49:30+07	\N	\N
69	2	69	2017-04-01 13:49:30+07	\N	\N
70	2	70	2017-04-01 13:49:30+07	\N	\N
71	2	71	2017-04-01 13:49:30+07	\N	\N
72	2	72	2017-04-01 13:49:30+07	\N	\N
\.


--
-- Name: r_loan_installment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mis_amartha
--

SELECT pg_catalog.setval('r_loan_installment_id_seq', 73, true);


--
-- Name: installment installment_pkey; Type: CONSTRAINT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY installment
    ADD CONSTRAINT installment_pkey PRIMARY KEY (id);


--
-- Name: loan loan_pkey; Type: CONSTRAINT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY loan
    ADD CONSTRAINT loan_pkey PRIMARY KEY (id);


--
-- Name: r_loan_installment r_loan_installment_pkey; Type: CONSTRAINT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY r_loan_installment
    ADD CONSTRAINT r_loan_installment_pkey PRIMARY KEY (id);


--
-- Name: installment_id__del_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX installment_id__del_idx ON installment USING btree (id, "deletedAt");


--
-- Name: installment_id_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX installment_id_idx ON installment USING btree (id);


--
-- Name: loan_id__del_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX loan_id__del_idx ON loan USING btree (id, "deletedAt");


--
-- Name: loan_id_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX loan_id_idx ON loan USING btree (id);


--
-- Name: r_loan_installment_id__del_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX r_loan_installment_id__del_idx ON r_loan_installment USING btree (id, "deletedAt");


--
-- Name: r_loan_installment_id_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX r_loan_installment_id_idx ON r_loan_installment USING btree (id);


--
-- Name: loan_history; Type: TABLE; Schema: public; Owner: mis_amartha
--

CREATE TABLE loan_history (
    id integer NOT NULL,
    "stageFrom" text,
    "stageTo" text,
    remark text,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone
);


--
-- Name: loan_history_id_seq; Type: SEQUENCE; Schema: public; Owner: mis_amartha
--

CREATE SEQUENCE loan_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loan_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mis_amartha
--

ALTER SEQUENCE loan_history_id_seq OWNED BY loan_history.id;


--
-- Name: r_loan_history; Type: TABLE; Schema: public; Owner: mis_amartha
--

CREATE TABLE r_loan_history (
    id integer NOT NULL,
    "loanId" integer,
    "loanHistoryId" integer,
    "createdAt" timestamp with time zone,
    "updatedAt" timestamp with time zone,
    "deletedAt" timestamp with time zone
);


--
-- Name: r_loan_history_id_seq; Type: SEQUENCE; Schema: public; Owner: mis_amartha
--

CREATE SEQUENCE r_loan_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: r_loan_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mis_amartha
--

ALTER SEQUENCE r_loan_history_id_seq OWNED BY r_loan_history.id;


--
-- Name: loan_history id; Type: DEFAULT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY loan_history ALTER COLUMN id SET DEFAULT nextval('loan_history_id_seq'::regclass);


--
-- Name: r_loan_history id; Type: DEFAULT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY r_loan_history ALTER COLUMN id SET DEFAULT nextval('r_loan_history_id_seq'::regclass);



--
-- Name: loan_history loan_history_pkey; Type: CONSTRAINT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY loan_history
    ADD CONSTRAINT loan_history_pkey PRIMARY KEY (id);


--
-- Name: r_loan_history r_loan_history_pkey; Type: CONSTRAINT; Schema: public; Owner: mis_amartha
--

ALTER TABLE ONLY r_loan_history
    ADD CONSTRAINT r_loan_history_pkey PRIMARY KEY (id);


--
-- Name: loan_history_id_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX loan_history_id_idx ON loan_history USING btree (id);


--
-- Name: r_loan_history_id__del_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX r_loan_history_id__del_idx ON r_loan_history USING btree (id, "deletedAt");


--
-- Name: r_loan_history_id_idx; Type: INDEX; Schema: public; Owner: mis_amartha
--

CREATE INDEX r_loan_history_id_idx ON r_loan_history USING btree (id);
--
-- PostgreSQL database dump complete
--

