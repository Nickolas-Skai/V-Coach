-- Filename: migrations/000001_create_core_tables.down.sql

-- Drop foreign key constraints first (reverse order of creation)

ALTER TABLE "coach_tips" DROP CONSTRAINT coach_tips_question_id_fkey;
ALTER TABLE "coach_tips" DROP CONSTRAINT coach_tips_session_id_fkey;

ALTER TABLE "responses" DROP CONSTRAINT responses_question_id_fkey;
ALTER TABLE "responses" DROP CONSTRAINT responses_session_id_fkey;

ALTER TABLE "sessions" DROP CONSTRAINT sessions_teacher_id_fkey;

ALTER TABLE "questions" DROP CONSTRAINT questions_created_by_fkey;

ALTER TABLE "users" DROP CONSTRAINT users_coach_id_fkey;
ALTER TABLE "users" DROP CONSTRAINT users_school_id_fkey;

-- Then drop the tables (reverse order)
DROP TABLE IF EXISTS "coach_tips";
DROP TABLE IF EXISTS "responses";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "questions";
DROP TABLE IF EXISTS "schools";
DROP TABLE IF EXISTS "users";
