-- Create databases
CREATE DATABASE gamification;
CREATE DATABASE lesson;
CREATE DATABASE identity;

-- Create users
CREATE USER gamification_user WITH PASSWORD 'gamipass';
CREATE USER lesson_user WITH PASSWORD 'lessonpass';
CREATE USER identity_user WITH PASSWORD 'idpass';

-- Grant access to their respective DBs
GRANT ALL PRIVILEGES ON DATABASE gamification TO gamification_user;
GRANT ALL PRIVILEGES ON DATABASE lesson TO lesson_user;
GRANT ALL PRIVILEGES ON DATABASE identity TO identity_user;

GRANT USAGE ON SCHEMA public TO gamification_user;
GRANT CREATE ON SCHEMA public TO gamification_user;
