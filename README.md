# gopqgen

This project aims to generate common CRUD operations for a given PostgreSQL database and expose them as a gRPC API.

## Why?

- Are you tired of writing CRUD?
- Do you realize there is actually nothing such as One-To-Many in RMDBS and all you do is deceive yourselves with Hibernate and others? 
- You realize SQL is already simple enough
- You use PostgreSQL 
- You realize you can define functions inside DB instead of the your application code
- How many times you debugged an ORM?
- Do you want another micro-service to your next-gen-kubernetes-based application?

**Note:** It is a work on progress.

## Installation

```bash
go get -u github.com/mustafaakin/gopqgen
```

## TODO Listgit 

- [X] List tables and fields
- [X] Find enums
- [X] List views
- [ ] Composite types
- [X] Generate getters for indexes (pkey and other composite keys)
- [X] User defined functions
- [ ] Generate functions from references foreign keys
- [ ] Update function
- [ ] gRPC Server implementation

## Example

For given this schema, it generates an almost valid gRPC right now.

SQL:

```sql
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
```

gRPC Proto:

```proto
syntax = "proto3";

// Enums
enum severity {
    UNKNOWN = 0;
    Easy = 1;
    Medium = 2;
    Hard = 3;
}

// Messages, Field Types
message teacher {
    int32 id = 1;
    string name = 2;
}

message AddArg {
    int32 var1 = 1;
    int32 var2 = 2;
}

message GetstudentsofcourseOut {
    string name = 1;
    string email = 2;
}

message IsuserincourseArg {
    int32 _studentid = 1;
    int32 _courseid = 2;
}

message course {
    int32 id = 1;
    string title = 2;
    int32 teacherid = 3;
    severity severity = 4;
}

message membership {
    int32 courseid = 1;
    int32 studentid = 2;
}

message student {
    int32 id = 1;
    string name = 2;
    string email = 3;
    bytes password = 4;
    int32 age = 5;
}

// Service Definition
service DatabaseService {
    // SELECT id, title, teacherid, severity FROM course
    rpc ListCourse(VoidRequest) returns (stream course) {}
    // SELECT * FROM getstudentsofcourse($1)
    rpc Getstudentsofcourse(int32) returns (GetstudentsofcourseOut) {}
    // SELECT id, title, teacherid, severity FROM course WHERE `id` = $1
    rpc GetCourseById(int32) returns (course) {}
    // SELECT id, name, email, password, age FROM student WHERE `id` = $1
    rpc GetStudentById(int32) returns (student) {}
    // SELECT id, name FROM teacher WHERE `id` = $1
    rpc GetTeacherById(int32) returns (teacher) {}
    // SELECT * FROM add($1, $2)
    rpc Add(AddArg) returns (int32) {}
    // SELECT * FROM isuserincourse($1, $2)
    rpc Isuserincourse(IsuserincourseArg) returns (bool) {}
    // SELECT courseid, studentid FROM membership
    rpc ListMembership(VoidRequest) returns (stream membership) {}
    // SELECT id, name, email, password, age FROM student
    rpc ListStudent(VoidRequest) returns (stream student) {}
    // SELECT id, name FROM teacher
    rpc ListTeacher(VoidRequest) returns (stream teacher) {}
}
```
