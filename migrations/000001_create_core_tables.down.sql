-- Filename: migrations/000001_create_core_tables.down.sql

-- Drop foreign key constraints first (reverse order of creation)

ALTER TABLE "coach_tips" DROP CONSTRAINT fk_coach_tips_question;
ALTER TABLE "coach_tips" DROP CONSTRAINT fk_coach_tips_session;

ALTER TABLE "responses" DROP CONSTRAINT fk_responses_question;
ALTER TABLE "responses" DROP CONSTRAINT fk_responses_session;

ALTER TABLE "sessions" DROP CONSTRAINT fk_sessions_teacher;

ALTER TABLE "users" DROP CONSTRAINT fk_users_coach;
ALTER TABLE "users" DROP CONSTRAINT fk_users_school;

-- Then drop the tables (reverse order)
DROP TABLE IF EXISTS "coach_tips";
DROP TABLE IF EXISTS "responses";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "questions";
DROP TABLE IF EXISTS "schools";
DROP TABLE IF EXISTS "users";