--
-- PostgreSQL database dump
--

-- Dumped from database version 12.22 (Ubuntu 12.22-0ubuntu0.20.04.2)
-- Dumped by pg_dump version 12.22 (Ubuntu 12.22-0ubuntu0.20.04.2)

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
-- Name: coach_tips; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.coach_tips (
    id integer NOT NULL,
    session_id integer NOT NULL,
    question_id integer NOT NULL,
    tip_text text NOT NULL,
    generated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.coach_tips OWNER TO vcoach;

--
-- Name: coach_tips_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.coach_tips_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.coach_tips_id_seq OWNER TO vcoach;

--
-- Name: coach_tips_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.coach_tips_id_seq OWNED BY public.coach_tips.id;


--
-- Name: questions; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.questions (
    id integer NOT NULL,
    text text NOT NULL,
    audio_url text,
    image_url text,
    required boolean DEFAULT true,
    type character varying(20) NOT NULL,
    options jsonb,
    is_active boolean DEFAULT true
);


ALTER TABLE public.questions OWNER TO vcoach;

--
-- Name: questions_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.questions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.questions_id_seq OWNER TO vcoach;

--
-- Name: questions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.questions_id_seq OWNED BY public.questions.id;


--
-- Name: responses; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.responses (
    id integer NOT NULL,
    session_id bigint NOT NULL,
    question_id integer NOT NULL,
    response_text text,
    audio_url text,
    confidence integer,
    submitted_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.responses OWNER TO vcoach;

--
-- Name: responses_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.responses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.responses_id_seq OWNER TO vcoach;

--
-- Name: responses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.responses_id_seq OWNED BY public.responses.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO vcoach;

--
-- Name: schools; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.schools (
    id integer NOT NULL,
    name character varying(150) NOT NULL,
    address text,
    district character varying(100),
    managment character varying(100)
);


ALTER TABLE public.schools OWNER TO vcoach;

--
-- Name: schools_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.schools_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.schools_id_seq OWNER TO vcoach;

--
-- Name: schools_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.schools_id_seq OWNED BY public.schools.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.sessions (
    id bigint NOT NULL,
    teacher_id integer NOT NULL,
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    ended_at timestamp without time zone
);


ALTER TABLE public.sessions OWNER TO vcoach;

--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sessions_id_seq OWNER TO vcoach;

--
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: vcoach
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    role character varying(10) NOT NULL,
    age integer,
    school_id integer,
    coach_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO vcoach;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: vcoach
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO vcoach;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: vcoach
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: coach_tips id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.coach_tips ALTER COLUMN id SET DEFAULT nextval('public.coach_tips_id_seq'::regclass);


--
-- Name: questions id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.questions ALTER COLUMN id SET DEFAULT nextval('public.questions_id_seq'::regclass);


--
-- Name: responses id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.responses ALTER COLUMN id SET DEFAULT nextval('public.responses_id_seq'::regclass);


--
-- Name: schools id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.schools ALTER COLUMN id SET DEFAULT nextval('public.schools_id_seq'::regclass);


--
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: coach_tips; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.coach_tips (id, session_id, question_id, tip_text, generated_at) FROM stdin;
\.


--
-- Data for Name: questions; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.questions (id, text, audio_url, image_url, required, type, options, is_active) FROM stdin;
1	Did you have any Limited Devices Session? If yes, mention how many	\N	\N	t	radio	["0", "1", "2", "3", "4", "5 or more"]	t
2	Which lessons from the platform did you use for your regular lesson(s) or limited devices lesson(s)? If none were taught, write 'NA'	\N	\N	t	text	[]	t
3	List the name of the lessons that were covered	\N	\N	t	text	[]	t
4	Were there any unplugged sessions? If yes, mention how many. *Note that you need to upload picture evidence later	\N	\N	t	radio	["0", "1", "2", "3", "4", "5 or more"]	t
5	If none were taught, why? *Select all that apply*	\N	\N	t	checkbox	["Doesn't apply, lessons were taught", "There was not enough time this month", "Not enough devices", "No internet access", "Teacher needs help organizing", "No interest from the students", "Teacher is on long leave", "Teacher is subject teaching but not science and technology", "Teacher is no longer in the primary level"]	t
6	Upload the evidence from the unplugged session	\N	\N	f	file	[]	t
7	Add a brief description for the unplugged/limited devices session(s), and the number of students who participated. *If none were taught, write 'NA'	\N	\N	t	text	[]	t
8	How many hours per week do you spend teaching computer science curriculum?	\N	\N	t	radio	["No lesson", "Less than 30 min", "Between 30 min and 1 hour", "Between 1 hour and 2 hours", "More than 2 hours"]	t
9	How easy/difficult did your students find the curriculum?	\N	\N	t	radio	["No lesson taught", "1 (Very Difficult)", "2", "3", "4", "5 (Very Easy)"]	t
10	How easy/difficult was it for you to implement the curriculum?	\N	\N	t	radio	["No lesson taught", "1 (Very Difficult)", "2", "3", "4", "5 (Very Easy)"]	t
11	How would you rate the curriculum overall?	\N	\N	t	radio	["No lesson taught", "1 (Positive)", "2", "3", "4", "5 (Negative)"]	t
12	List here any questions and/or comments the teachers would like to report	\N	\N	t	text	[]	t
13	Has the teacher experienced any issue with equipment provided by the Ministry or Code Caribbean? (You can select more than one)	\N	\N	t	checkbox	["None", "Laptop", "Micro-bit"]	t
14	If yes, please explain the issue.	\N	\N	f	text	[]	t
\.


--
-- Data for Name: responses; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.responses (id, session_id, question_id, response_text, audio_url, confidence, submitted_at) FROM stdin;
8	5	1	5 or more	\N	\N	2025-05-11 02:34:24
9	5	1	5 or more	\N	\N	2025-05-11 02:42:55
10	6	1	5 or more	\N	\N	2025-05-11 02:47:46
11	7	1	5 or more	\N	\N	2025-05-11 14:00:54
13	10	1	4	\N	\N	2025-05-11 20:49:46
14	10	2	sdfsdf	\N	\N	2025-05-11 20:49:49
15	10	3	asdfasdf	\N	\N	2025-05-11 20:49:51
16	10	4	4	\N	\N	2025-05-11 20:49:54
17	10	5	Teacher needs help organizing	\N	\N	2025-05-11 20:49:58
18	10	1	3	\N	\N	2025-05-11 21:01:17
19	10	2	asdf	\N	\N	2025-05-11 21:01:20
20	10	3	asdf	\N	\N	2025-05-11 21:01:22
21	10	4	3	\N	\N	2025-05-11 21:01:25
22	10	5	Teacher is on long leave	\N	\N	2025-05-11 21:01:27
23	10	1	3	\N	\N	2025-05-11 21:32:29
24	10	2	d	\N	\N	2025-05-11 21:32:34
25	10	3	a	\N	\N	2025-05-11 21:32:35
26	10	4	3	\N	\N	2025-05-11 21:32:37
27	10	5	Teacher is on long leave	\N	\N	2025-05-11 21:32:40
28	10	1	3	\N	\N	2025-05-11 21:43:57
29	10	2	s	\N	\N	2025-05-11 21:44:07
30	10	3	s	\N	\N	2025-05-11 21:44:09
31	10	4	2	\N	\N	2025-05-11 21:44:11
32	10	5	Teacher needs help organizing	\N	\N	2025-05-11 21:44:13
33	10	6	./uploads/ChatGPT Image May 9, 2025, 08_41_09 PM.png	\N	\N	2025-05-11 21:44:21
34	10	7	s	\N	\N	2025-05-11 21:44:26
35	10	8	Between 1 hour and 2 hours	\N	\N	2025-05-11 21:44:29
36	10	9	4	\N	\N	2025-05-11 21:44:31
37	10	10	4	\N	\N	2025-05-11 21:44:32
38	10	11	4	\N	\N	2025-05-11 21:44:34
39	10	12	ssd	\N	\N	2025-05-11 21:44:36
40	10	13	Laptop	\N	\N	2025-05-11 21:44:39
41	10	14	sdf	\N	\N	2025-05-11 21:44:41
42	11	1	1	\N	\N	2025-05-11 21:56:21
43	11	2	Na	\N	\N	2025-05-11 21:56:33
44	11	3	All	\N	\N	2025-05-11 21:56:38
45	11	4	2	\N	\N	2025-05-11 21:56:43
46	11	5	There was not enough time this month	\N	\N	2025-05-11 21:56:52
47	11	6	./uploads/Screenshot_20250509_110457_Pinterest.jpg	\N	\N	2025-05-11 21:57:54
48	11	7	it was two devices and two students	\N	\N	2025-05-11 21:59:48
49	11	8	Between 30 min and 1 hour	\N	\N	2025-05-11 22:01:34
50	11	9	3	\N	\N	2025-05-11 22:02:33
51	11	10	3	\N	\N	2025-05-11 22:03:15
52	11	11	4	\N	\N	2025-05-11 22:03:32
53	11	12	lessons need to be a little bit longer	\N	\N	2025-05-11 22:04:14
54	11	13	None	\N	\N	2025-05-11 22:04:35
55	11	14	no issues	\N	\N	2025-05-11 22:04:50
56	11	1	4	\N	\N	2025-05-11 22:06:25
57	19	1	4	\N	\N	2025-05-12 15:45:25
58	19	2	yes, there is.	\N	\N	2025-05-12 15:46:22
59	19	3	math, english, science.	\N	\N	2025-05-12 15:46:33
60	19	4	4	\N	\N	2025-05-12 15:46:43
61	19	1	2	\N	\N	2025-05-12 16:50:15
62	19	1	3	\N	\N	2025-05-12 17:36:55
63	19	2	asdf	\N	\N	2025-05-12 17:36:56
64	19	3	asdf	\N	\N	2025-05-12 17:36:58
65	19	4	2	\N	\N	2025-05-12 17:37:00
66	19	5	No interest from the students	\N	\N	2025-05-12 17:37:02
67	19	6	./uploads/ub_background.png	\N	\N	2025-05-12 17:37:13
68	19	7	asdf	\N	\N	2025-05-12 17:37:15
69	19	8	More than 2 hours	\N	\N	2025-05-12 17:37:18
70	19	9	4	\N	\N	2025-05-12 17:37:20
71	19	10	5 (Very Easy)	\N	\N	2025-05-12 17:37:22
72	19	11	5 (Negative)	\N	\N	2025-05-12 17:37:23
73	19	12	asdf	\N	\N	2025-05-12 17:37:25
74	19	13	Micro-bit	\N	\N	2025-05-12 17:37:26
75	19	14	asdf	\N	\N	2025-05-12 17:37:28
76	19	1	4	\N	\N	2025-05-12 17:46:58
77	19	1	4	\N	\N	2025-05-12 17:51:31
78	19	1	5 or more	\N	\N	2025-05-12 17:51:51
79	19	1	5 or more	\N	\N	2025-05-12 17:52:36
80	19	1	5 or more	\N	\N	2025-05-12 18:03:07
81	20	1	4	\N	\N	2025-05-12 21:23:36
82	20	2	SDF	\N	\N	2025-05-12 21:23:39
83	20	3	asdf	\N	\N	2025-05-12 21:23:45
84	20	4	3	\N	\N	2025-05-12 21:23:47
85	20	5	Teacher is on long leave	\N	\N	2025-05-12 21:23:50
86	20	6	./uploads/code.png	\N	\N	2025-05-12 21:24:06
87	20	7	kjb	\N	\N	2025-05-12 21:24:20
88	21	1	2	\N	\N	2025-05-12 21:37:44
89	21	2	bbb	\N	\N	2025-05-12 21:37:48
90	22	1	4	\N	\N	2025-05-13 01:07:40
91	22	2	g	\N	\N	2025-05-13 01:07:44
92	22	3	g	\N	\N	2025-05-13 01:07:45
93	22	4	3	\N	\N	2025-05-13 01:07:47
94	22	5	Teacher is on long leave	\N	\N	2025-05-13 01:07:48
95	22	6	./uploads/Colorful Pastel Simple Comparative SWOT Analysis Weakness Threats Strengths Opportunities Chart Graph (2).png	\N	\N	2025-05-13 01:08:00
96	22	7	NA	\N	\N	2025-05-13 01:08:06
97	22	8	Between 1 hour and 2 hours	\N	\N	2025-05-13 01:08:07
98	22	9	4	\N	\N	2025-05-13 01:08:09
99	22	10	4	\N	\N	2025-05-13 01:08:10
100	22	11	4	\N	\N	2025-05-13 01:08:12
101	22	12	asdf	\N	\N	2025-05-13 01:08:15
102	22	13	Micro-bit	\N	\N	2025-05-13 01:08:16
103	22	14	asdf	\N	\N	2025-05-13 01:08:18
104	27	1	5 or more	\N	\N	2025-05-13 04:35:22
105	27	2	asdf	\N	\N	2025-05-13 04:35:23
106	27	3	asdf	\N	\N	2025-05-13 04:35:24
107	27	4	1	\N	\N	2025-05-13 04:35:26
108	27	5	No internet access	\N	\N	2025-05-13 04:35:29
109	27	6	./uploads/RD_signature.png	\N	\N	2025-05-13 04:35:37
110	27	7	asdf	\N	\N	2025-05-13 04:35:38
111	27	8	Between 30 min and 1 hour	\N	\N	2025-05-13 04:35:40
112	27	9	4	\N	\N	2025-05-13 04:35:42
113	27	10	4	\N	\N	2025-05-13 04:35:43
114	27	11	5 (Negative)	\N	\N	2025-05-13 04:35:44
115	27	12	asdf	\N	\N	2025-05-13 04:35:46
116	27	13	Laptop	\N	\N	2025-05-13 04:35:47
117	27	14	asdf	\N	\N	2025-05-13 04:35:49
118	27	1	3	\N	\N	2025-05-13 04:58:36
119	27	2	asdf	\N	\N	2025-05-13 04:58:37
120	27	3	asdf	\N	\N	2025-05-13 04:58:38
121	27	4	0	\N	\N	2025-05-13 04:58:39
122	27	5	Teacher needs help organizing	\N	\N	2025-05-13 04:58:41
123	27	6	./uploads/VACANCY.png .jpg	\N	\N	2025-05-13 04:58:50
124	27	7	asdf	\N	\N	2025-05-13 04:58:52
125	27	8	Between 30 min and 1 hour	\N	\N	2025-05-13 04:58:54
126	27	9	4	\N	\N	2025-05-13 04:58:55
127	27	10	4	\N	\N	2025-05-13 04:58:57
128	27	11	5 (Negative)	\N	\N	2025-05-13 04:58:58
129	27	12	asdf	\N	\N	2025-05-13 04:58:59
130	27	13	Micro-bit	\N	\N	2025-05-13 04:59:00
131	27	14	asdf	\N	\N	2025-05-13 04:59:01
132	27	14	asdf	\N	\N	2025-05-13 05:00:44
133	27	1	3	\N	\N	2025-05-13 05:03:08
134	27	2	asdf	\N	\N	2025-05-13 05:04:35
135	27	3	asdf	\N	\N	2025-05-13 05:04:37
136	27	4	0	\N	\N	2025-05-13 05:04:38
137	27	5	No interest from the students	\N	\N	2025-05-13 05:04:40
138	27	6	./uploads/draw_text.png	\N	\N	2025-05-13 05:04:47
139	27	6	./uploads/draw_text.png	\N	\N	2025-05-13 05:04:48
140	27	7	asdf	\N	\N	2025-05-13 05:04:50
141	27	8	Between 1 hour and 2 hours	\N	\N	2025-05-13 05:04:52
142	27	9	4	\N	\N	2025-05-13 05:04:54
143	27	10	4	\N	\N	2025-05-13 05:04:56
144	27	11	4	\N	\N	2025-05-13 05:04:59
145	27	12	asdf	\N	\N	2025-05-13 05:05:00
146	27	13	Micro-bit	\N	\N	2025-05-13 05:05:01
147	27	14	sadf	\N	\N	2025-05-13 05:05:03
148	27	1	4	\N	\N	2025-05-13 05:37:06
149	27	2	asdf	\N	\N	2025-05-13 05:37:08
150	27	3	asdf	\N	\N	2025-05-13 05:37:09
151	27	4	0	\N	\N	2025-05-13 05:37:10
152	27	5	No internet access	\N	\N	2025-05-13 05:37:12
153	27	6	./uploads/Screenshot 2025-01-26 131207.png	\N	\N	2025-05-13 05:37:32
154	27	7	asdf	\N	\N	2025-05-13 05:37:34
155	27	8	Between 1 hour and 2 hours	\N	\N	2025-05-13 05:37:35
156	27	9	4	\N	\N	2025-05-13 05:37:38
157	27	10	5 (Very Easy)	\N	\N	2025-05-13 05:37:40
158	27	11	4	\N	\N	2025-05-13 05:37:41
159	27	12	asdf	\N	\N	2025-05-13 05:37:42
160	27	13	Micro-bit	\N	\N	2025-05-13 05:37:45
161	27	14	asdf	\N	\N	2025-05-13 05:37:46
162	27	1	2	\N	\N	2025-05-13 05:39:32
163	27	2	asdf	\N	\N	2025-05-13 05:39:33
164	27	3	asdf	\N	\N	2025-05-13 05:39:35
165	27	4	1	\N	\N	2025-05-13 05:39:37
166	27	5	Teacher needs help organizing	\N	\N	2025-05-13 05:39:39
167	27	6	./uploads/Screenshot 2025-01-25 120324.png	\N	\N	2025-05-13 05:39:42
168	27	7	asdf	\N	\N	2025-05-13 05:39:45
169	27	8	Between 30 min and 1 hour	\N	\N	2025-05-13 05:39:47
170	27	9	4	\N	\N	2025-05-13 05:39:49
171	27	10	4	\N	\N	2025-05-13 05:39:50
172	27	11	4	\N	\N	2025-05-13 05:39:51
173	27	12	asdf	\N	\N	2025-05-13 05:39:54
174	27	13	Micro-bit	\N	\N	2025-05-13 05:39:56
175	27	14	asdf	\N	\N	2025-05-13 05:39:57
176	27	14	asdf	\N	\N	2025-05-13 05:39:58
177	27	1	5 or more	\N	\N	2025-05-13 05:42:07
178	27	2	asdf	\N	\N	2025-05-13 05:42:12
179	27	3	asdf	\N	\N	2025-05-13 05:42:13
180	27	4	1	\N	\N	2025-05-13 05:42:15
181	27	5	Teacher needs help organizing	\N	\N	2025-05-13 05:42:16
182	27	6	./uploads/Screenshot 2025-01-27 003830.png	\N	\N	2025-05-13 05:42:20
183	27	7	asdf	\N	\N	2025-05-13 05:42:21
184	27	8	More than 2 hours	\N	\N	2025-05-13 05:42:23
185	27	9	5 (Very Easy)	\N	\N	2025-05-13 05:42:24
186	27	10	5 (Very Easy)	\N	\N	2025-05-13 05:42:26
187	27	11	5 (Negative)	\N	\N	2025-05-13 05:42:27
188	27	12	asdf	\N	\N	2025-05-13 05:42:30
189	27	13	Micro-bit	\N	\N	2025-05-13 05:42:33
190	27	14	maybe there isn't.	\N	\N	2025-05-13 05:42:40
191	28	1	5 or more	\N	\N	2025-05-13 05:43:50
192	28	2	zds	\N	\N	2025-05-13 05:43:53
193	28	3	asdf	\N	\N	2025-05-13 05:43:54
194	28	4	1	\N	\N	2025-05-13 05:43:55
195	28	5	Teacher is subject teaching but not science and technology	\N	\N	2025-05-13 05:43:56
196	28	6	./uploads/Screenshot 2025-01-26 230242.png	\N	\N	2025-05-13 05:44:02
197	28	7	asdf	\N	\N	2025-05-13 05:44:04
198	28	8	Between 30 min and 1 hour	\N	\N	2025-05-13 05:44:05
199	28	9	4	\N	\N	2025-05-13 05:44:07
200	28	10	5 (Very Easy)	\N	\N	2025-05-13 05:44:09
201	28	11	5 (Negative)	\N	\N	2025-05-13 05:44:10
202	28	12	asdf	\N	\N	2025-05-13 05:44:11
203	28	13	Micro-bit	\N	\N	2025-05-13 05:44:13
204	28	14	asdf	\N	\N	2025-05-13 05:44:14
205	28	1	4	\N	\N	2025-05-13 05:54:25
206	28	2	dvcsdv	\N	\N	2025-05-13 05:54:31
207	28	3	asdfasdfadfasdf	\N	\N	2025-05-13 05:54:33
208	28	4	2	\N	\N	2025-05-13 05:54:35
209	28	5	Teacher needs help organizing	\N	\N	2025-05-13 05:54:37
210	28	6	./uploads/Screenshot 2025-02-02 002827.png	\N	\N	2025-05-13 05:54:43
211	28	7	asdf	\N	\N	2025-05-13 05:54:46
212	28	8	Between 1 hour and 2 hours	\N	\N	2025-05-13 05:54:47
213	28	9	4	\N	\N	2025-05-13 05:54:49
214	28	10	3	\N	\N	2025-05-13 05:54:52
215	28	11	5 (Negative)	\N	\N	2025-05-13 05:54:53
216	28	12	asdf	\N	\N	2025-05-13 05:54:55
217	28	13	Micro-bit	\N	\N	2025-05-13 05:54:57
218	28	14	asdf	\N	\N	2025-05-13 05:54:58
219	42	1	4	\N	\N	2025-05-14 06:52:04
220	42	2	asdf	\N	\N	2025-05-14 06:52:20
221	42	3	asdf	\N	\N	2025-05-14 06:52:22
222	42	4	3	\N	\N	2025-05-14 06:52:33
223	42	5	Teacher is on long leave	\N	\N	2025-05-14 06:52:35
224	42	6	./uploads/Screenshot 2025-02-02 002823.png	\N	\N	2025-05-14 06:52:44
225	42	7	and i briefed the description for the ipod limited device sessions and the number of students who are participated. if no were taught, write any.	\N	\N	2025-05-14 06:54:29
226	42	8	More than 2 hours	\N	\N	2025-05-14 06:54:43
227	42	9	5 (Very Easy)	\N	\N	2025-05-14 06:54:45
228	42	10	5 (Very Easy)	\N	\N	2025-05-14 06:54:46
229	42	11	5 (Negative)	\N	\N	2025-05-14 06:54:47
230	42	12	nn	\N	\N	2025-05-14 06:54:50
231	42	13	Micro-bit	\N	\N	2025-05-14 06:54:58
232	42	14	nn	\N	\N	2025-05-14 06:55:02
233	48	1	5 or more	\N	\N	2025-05-14 19:09:07
234	48	2	cvffffffff	\N	\N	2025-05-14 19:09:13
235	48	1	5 or more	\N	\N	2025-05-14 20:05:58
236	48	1	5 or more	\N	\N	2025-05-14 20:11:50
237	48	1	5 or more	\N	\N	2025-05-14 20:13:56
238	48	1	5 or more	\N	\N	2025-05-14 20:16:25
239	48	1	5 or more	\N	\N	2025-05-14 20:16:42
240	48	1	5 or more	\N	\N	2025-05-14 20:19:35
241	48	1	5 or more	\N	\N	2025-05-14 20:19:42
242	50	1	3	\N	\N	2025-05-14 20:50:10
243	50	1	3	\N	\N	2025-05-14 20:51:26
244	50	1	3	\N	\N	2025-05-14 20:51:41
245	50	2	dd	\N	\N	2025-05-14 20:51:51
246	50	3	dd	\N	\N	2025-05-14 20:51:54
247	50	4	3	\N	\N	2025-05-14 20:52:03
248	50	5	Teacher is subject teaching but not science and technology	\N	\N	2025-05-14 20:52:05
249	50	6	./uploads/Screenshot 2025-01-27 010941.png	\N	\N	2025-05-14 20:52:13
250	50	7	lnknb	\N	\N	2025-05-14 20:52:17
251	50	8	More than 2 hours	\N	\N	2025-05-14 20:52:19
252	50	9	5 (Very Easy)	\N	\N	2025-05-14 20:52:20
253	50	10	5 (Very Easy)	\N	\N	2025-05-14 20:52:21
254	50	11	5 (Negative)	\N	\N	2025-05-14 20:52:22
255	50	12	xxx	\N	\N	2025-05-14 20:52:25
256	50	13	Micro-bit	\N	\N	2025-05-14 20:52:26
257	50	14	xxxx	\N	\N	2025-05-14 20:52:28
258	51	1	2	\N	\N	2025-05-16 06:58:39
259	51	2	N/A	\N	\N	2025-05-16 06:58:51
260	51	3	8	\N	\N	2025-05-16 06:59:04
261	51	4	1	\N	\N	2025-05-16 07:01:42
262	51	5	Not enough devices	\N	\N	2025-05-16 07:02:15
263	51	5	Not enough devices	\N	\N	2025-05-16 07:02:32
264	51	6	./uploads/IMG_1117.jpeg	\N	\N	2025-05-16 07:03:00
265	51	7	N/A	\N	\N	2025-05-16 07:03:10
266	51	8	Between 1 hour and 2 hours	\N	\N	2025-05-16 07:03:32
267	51	9	1 (Very Difficult)	\N	\N	2025-05-16 07:03:38
268	51	10	4	\N	\N	2025-05-16 07:03:50
269	51	11	1 (Positive)	\N	\N	2025-05-16 07:03:57
270	51	12	N/A	\N	\N	2025-05-16 07:04:03
271	51	13	None	\N	\N	2025-05-16 07:04:07
272	51	14	N/A	\N	\N	2025-05-16 07:04:12
273	51	14	N/A	\N	\N	2025-05-16 07:04:29
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.schema_migrations (version, dirty) FROM stdin;
2	f
\.


--
-- Data for Name: schools; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.schools (id, name, address, district, managment) FROM stdin;
1	Anglican Primary School	123 Church Street, Belmopan	Cayo	Religious
2	Saint Joseph R.C. School	456 River Road, Orange Walk	Orange Walk	Catholic
3	Belize High School	789 Constitution Ave, Belize City	Belize	Private
4	Harmony Government School	12 Pine Street, Dangriga	Stann Creek	Government
5	Sunrise Academy	34 Sunrise Blvd, Corozal	Corozal	Private
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.sessions (id, teacher_id, started_at, ended_at) FROM stdin;
4	1	2023-10-01 00:00:00	2023-10-02 00:00:00
5	1	2023-10-01 00:00:00	2023-10-02 00:00:00
6	1	2023-10-01 00:00:00	2023-10-02 00:00:00
7	1	2023-10-01 00:00:00	2023-10-02 00:00:00
8	1	2023-10-01 00:00:00	2023-10-02 00:00:00
9	1	2023-10-01 00:00:00	2023-10-02 00:00:00
10	1	2023-10-01 00:00:00	2023-10-02 00:00:00
11	1	2023-10-01 00:00:00	2023-10-02 00:00:00
12	1	2023-10-01 00:00:00	2023-10-02 00:00:00
19	1	2023-10-01 00:00:00	2023-10-02 00:00:00
20	1	2023-10-01 00:00:00	2023-10-02 00:00:00
21	1	2023-10-01 00:00:00	2023-10-02 00:00:00
22	1	2023-10-01 00:00:00	2023-10-02 00:00:00
23	1	2023-10-01 00:00:00	2023-10-02 00:00:00
24	1	2023-10-01 00:00:00	2023-10-02 00:00:00
25	1	2023-10-01 00:00:00	2023-10-02 00:00:00
26	1	2023-10-01 00:00:00	2023-10-02 00:00:00
27	1	2023-10-01 00:00:00	2023-10-02 00:00:00
28	21	2025-05-13 05:43:21	\N
29	18	2025-05-13 08:31:50	\N
30	18	2025-05-13 20:05:49	\N
32	1	2023-10-01 00:00:00	2023-10-02 00:00:00
36	1	2023-10-01 00:00:00	2023-10-02 00:00:00
41	20	2025-05-14 06:03:14	\N
42	20	2025-05-14 06:05:14	\N
43	20	2025-05-14 06:58:30	\N
44	21	2025-05-14 07:56:15	\N
45	18	2025-05-14 08:57:18	\N
46	18	2025-05-14 13:57:51	\N
47	18	2025-05-14 15:41:41	\N
48	18	2025-05-14 19:08:50	\N
49	18	2025-05-14 20:41:12	\N
50	27	2025-05-14 20:50:04	\N
51	28	2025-05-16 06:57:43	\N
52	27	2025-05-16 07:00:45	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: vcoach
--

COPY public.users (id, name, email, password_hash, role, age, school_id, coach_id, created_at) FROM stdin;
1	D'Alesseo Requena	2021154699@ub.edu.bz	iloveyou143		20	\N	\N	2025-05-01 03:50:53.48851
2	Amilcar Vasquez 	asdf@gmail.com	$2a$10$hakTKNZBNBSCeajnlukBw.DNPbLzjlwUE33yge5d79n8a.qillU4i	coach	\N	\N	\N	2025-05-01 20:02:53.462872
11	Amilcar Vasquez 	2022156572@ub.edu.bz	iloveyou143	coach	23	1	\N	2025-05-01 22:41:38.278811
13	test	test@test.com	fatguy_123	coach	23	1	\N	2025-05-01 22:51:20.720161
14	Test2_withhasing 	test2@test.com	$2a$10$I.de7ZghBQxIGQtFJ0YCM.A2Zyc0TMmoOcSN1TecqUlGdVyhmrN3K	teacher	42	1	\N	2025-05-01 22:58:22.272203
15	Render Test	RenderTest@gmail.com	$2a$10$9zJ1j7aOKj8wEsegWV8WNOHC/Bav/XfeZK25GJ2wrvTXAvO/d/Y96	coach	100	1	\N	2025-05-01 23:05:31.364369
16	Render Test twp	RenderTest2@gmail.com	$2a$10$Akcc4.u5ohNaDTnSD9npg./CcrMnyddGKW.i90jia8e5OqgZB7kXm	teacher	56	1	\N	2025-05-01 23:12:24.85084
17	testlogin	testlogin@gmail.com	$2a$10$tRbrhmQCIgB.H8Py0TGhuO7wTRLh8YRKo4yiJdZIV5CtzAm3fYFEG	teacher	111	1	\N	2025-05-02 00:08:37.397524
18	Coachtesting	coach@gmail.com	$2a$10$7xT/Cq1c3JBtk9virpkLTOMM19yBiF.dMTK3NjGxceBmi9uL/J1.C	coach	11	4	\N	2025-05-09 22:38:02.766007
19	Kamera 	kamera_requena@yahoo.com	$2a$10$gM6IqaZClN8lytOlHsyNaeggkbEXArsUeVBRFjgJcXxnt7pP4z/ja	teacher	36	1	\N	2025-05-11 21:54:27.930044
20	Coachtesting3	coach234@gmail.com	$2a$10$vrblO7B62xjaHsTlZ7ZnkejkU7RCWDOPb0nF/r.RSw.lJ/xTpGA4S	coach	111	1	\N	2025-05-12 01:44:29.476954
21	testinginterview	sss@gmail.com	$2a$10$vW2wIu8dGfIr2D5a51.tje5I.eMCivFtxGM6Ox5QTYhXDIuCN8E8a	teacher	47	1	\N	2025-05-12 05:06:30.591235
22	lolkjhg	cjhvjhvuvoach@gmail.com	$2a$10$EYN1a468wlfK4emlyQnFn.U1nPxkY4ldKaQd71qFumDuLENtuhUVy	student	18	1	\N	2025-05-14 18:48:14.88208
23	lkajbdsfipbjsdfbsadf	ijasbfildjsfibfslaifjds@gmail.com	$2a$10$6bZTvJIY8mE7bGOGIbb97e/slqa0kQAYSG0cJe.CKyxR01n2O96n.	student	18	1	\N	2025-05-14 18:50:10.309751
24	asdfasdfasdfasdf	asdfsadfasdfasdfasdfasdf@gmail.com	$2a$10$ySEon2yv80L/ed2CgF6aT.JR6dPGsOv38pr7KHluG/a11EGsVXtrO	student	41	1	\N	2025-05-14 18:55:56.658338
25	loll	asdfasdfads@gmail.com	$2a$10$uVyw7rtwcaNpN5RwZBUIZeK7cJpDNKPRdTThuEC78OOSqvgCDdWH6	student	44	1	\N	2025-05-14 20:31:30.105134
26	LEGOHOME	legooohome@gmail.com	$2a$10$kMwxNlO2DotfjtkwWotDgOCqPH0HszF5qlsUPbTFaCaT3krv0Fpp6	student	19	1	\N	2025-05-14 20:48:30.244494
27	demcom	demo@gmail.com	$2a$10$xjex5EmmIWh55EnjNgnRxuQ4jgQdn2gU3OHfCNK7utIqHeOxl3ulS	student	20	1	\N	2025-05-14 20:49:33.941942
28	Daniel	monterodaniel155@gmail.com	$2a$10$nA66jdRCy2aDg/m1G1PvzOKCkxtBzm9vnyZ5Jgrbb5IseYQCAGUKO	coach	20	1	\N	2025-05-16 06:56:53.013133
\.


--
-- Name: coach_tips_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.coach_tips_id_seq', 1, false);


--
-- Name: questions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.questions_id_seq', 14, true);


--
-- Name: responses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.responses_id_seq', 273, true);


--
-- Name: schools_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.schools_id_seq', 5, true);


--
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.sessions_id_seq', 52, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: vcoach
--

SELECT pg_catalog.setval('public.users_id_seq', 28, true);


--
-- Name: coach_tips coach_tips_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.coach_tips
    ADD CONSTRAINT coach_tips_pkey PRIMARY KEY (id);


--
-- Name: questions questions_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pkey PRIMARY KEY (id);


--
-- Name: responses responses_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.responses
    ADD CONSTRAINT responses_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: schools schools_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.schools
    ADD CONSTRAINT schools_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: coach_tips fk_coach_tips_question; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.coach_tips
    ADD CONSTRAINT fk_coach_tips_question FOREIGN KEY (question_id) REFERENCES public.questions(id) ON DELETE CASCADE;


--
-- Name: coach_tips fk_coach_tips_session; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.coach_tips
    ADD CONSTRAINT fk_coach_tips_session FOREIGN KEY (session_id) REFERENCES public.sessions(id) ON DELETE CASCADE;


--
-- Name: responses fk_responses_question; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.responses
    ADD CONSTRAINT fk_responses_question FOREIGN KEY (question_id) REFERENCES public.questions(id) ON DELETE CASCADE;


--
-- Name: responses fk_responses_session; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.responses
    ADD CONSTRAINT fk_responses_session FOREIGN KEY (session_id) REFERENCES public.sessions(id) ON DELETE CASCADE;


--
-- Name: sessions fk_sessions_teacher; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_sessions_teacher FOREIGN KEY (teacher_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users fk_users_coach; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_coach FOREIGN KEY (coach_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: users fk_users_school; Type: FK CONSTRAINT; Schema: public; Owner: vcoach
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_school FOREIGN KEY (school_id) REFERENCES public.schools(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

