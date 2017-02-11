CREATE TABLE student (
  id SERIAL PRIMARY KEY,
  name TEXT,
  email TEXT,
  password BYTEA,
  age INT
);

CREATE TABLE teacher (
  id  SERIAL PRIMARY KEY,
  name TEXT
);

CREATE TYPE severity AS ENUM ('Easy', 'Medium', 'Hard');

CREATE TABLE course(
  id        SERIAL PRIMARY KEY,
  title     TEXT,
  teacherId INT REFERENCES teacher(id),
  severity  severity
);

CREATE TABLE membership (
  courseId INT REFERENCES course(id),
  studentId INT REFERENCES student(id)
);