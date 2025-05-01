-- Filename: migrations/000001_create_core_tables.up.sql
-- This migration creates the core tables for the application, including users, schools, questions, sessions, responses, and coach tips.
-- It also establishes the necessary foreign key relationships between these tables.

-- Create the "users" table
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

-- Create the "schools" table
CREATE TABLE "schools" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(150) NOT NULL,
  "address" TEXT,
  "district" VARCHAR(100),
  "managment" VARCHAR(100)
);

-- Create the "questions" table
CREATE TABLE "questions" (
  "id" SERIAL PRIMARY KEY,
  "text" TEXT NOT NULL,
  "audio_url" TEXT,
  "image_url" TEXT,
  "required" BOOLEAN DEFAULT true,
  "type" VARCHAR(20) NOT NULL,
  "options" JSONB,
  "is_active" BOOLEAN DEFAULT true
);

-- Create the "sessions" table
CREATE TABLE "sessions" (
  "id" SERIAL PRIMARY KEY,
  "teacher_id" INT NOT NULL,
  "started_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "ended_at" TIMESTAMP
);

-- Create the "responses" table
CREATE TABLE "responses" (
  "id" SERIAL PRIMARY KEY,
  "session_id" INT NOT NULL,
  "question_id" INT NOT NULL,
  "response_text" TEXT,
  "audio_url" TEXT,
  "confidence" INT,
  "submitted_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create the "coach_tips" table
CREATE TABLE "coach_tips" (
  "id" SERIAL PRIMARY KEY,
  "session_id" INT NOT NULL,
  "question_id" INT NOT NULL,
  "tip_text" TEXT NOT NULL,
  "generated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);


-- Add foreign key constraints
ALTER TABLE "users" ADD CONSTRAINT fk_users_school FOREIGN KEY ("school_id") REFERENCES "schools" ("id") ON DELETE SET NULL;
ALTER TABLE "users" ADD CONSTRAINT fk_users_coach FOREIGN KEY ("coach_id") REFERENCES "users" ("id") ON DELETE SET NULL;

ALTER TABLE "sessions" ADD CONSTRAINT fk_sessions_teacher FOREIGN KEY ("teacher_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "responses" ADD CONSTRAINT fk_responses_session FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON DELETE CASCADE;
ALTER TABLE "responses" ADD CONSTRAINT fk_responses_question FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE;

ALTER TABLE "coach_tips" ADD CONSTRAINT fk_coach_tips_session FOREIGN KEY ("session_id") REFERENCES "sessions" ("id") ON DELETE CASCADE;
ALTER TABLE "coach_tips" ADD CONSTRAINT fk_coach_tips_question FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE;