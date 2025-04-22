-- Filename: migrations/000001_create_core_tables.up.sql
-- This migration creates the core tables for the application, including users, schools, questions, sessions, responses, and coach tips.
-- It also establishes the necessary foreign key relationships between these tables.
-- This migration is essential for setting up the database schema and ensuring data integrity across related tables.
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(100) NOT NULL,
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "role" VARCHAR(10) NOT NULL,
  "age" INT,
  "school_id" INT,
  "coach_id" INT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "schools" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(150) NOT NULL,
  "address" TEXT,
  "district" VARCHAR(100),
  "managment" VARCHAR(100)
);

CREATE TABLE "questions" (
  "id" SERIAL PRIMARY KEY,
  "text" TEXT NOT NULL,
  "type" VARCHAR(20) NOT NULL,
  "options" JSONB,
  "is_active" BOOLEAN DEFAULT true,
  "created_by" INT NOT NULL
);

CREATE TABLE "sessions" (
  "id" SERIAL PRIMARY KEY,
  "teacher_id" INT NOT NULL,
  "started_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "ended_at" TIMESTAMP
);

CREATE TABLE "responses" (
  "id" SERIAL PRIMARY KEY,
  "session_id" INT NOT NULL,
  "question_id" INT NOT NULL,
  "response_text" TEXT,
  "audio_url" TEXT,
  "confidence" INT,
  "submitted_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "coach_tips" (
  "id" SERIAL PRIMARY KEY,
  "session_id" INT NOT NULL,
  "question_id" INT NOT NULL,
  "tip_text" TEXT NOT NULL,
  "generated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE "users" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("coach_id") REFERENCES "users" ("id");

ALTER TABLE "questions" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("teacher_id") REFERENCES "users" ("id");

ALTER TABLE "responses" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "responses" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");

ALTER TABLE "coach_tips" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "coach_tips" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");
