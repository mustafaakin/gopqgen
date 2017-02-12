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

-- Not that you need something like it, but just imagine
CREATE FUNCTION add(integer, integer) RETURNS integer AS
  'select $1 + $2;' LANGUAGE SQL
IMMUTABLE RETURNS NULL ON NULL INPUT;

CREATE FUNCTION IsUserInCourse(_studentId integer, _courseid integer) RETURNS boolean AS
  'select exists(select 1 from membership where studentId = _studentid AND courseid = _courseid)' LANGUAGE SQL;

CREATE FUNCTION GetStudentsOfCourse(_courseId integer) RETURNS TABLE(name text, email text) AS
  'SELECT s.name, s.email FROM
    course c,
    student s,
    membership m
  WHERE
    c.id = m.courseid AND
    s.id = m.studentid AND
    c.id = _courseId'
LANGUAGE SQL;
