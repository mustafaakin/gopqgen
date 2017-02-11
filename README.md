# gopqgen

This project aims to generate common CRUD operations for a given PostgreSQL database and expose them as a gRPC API.

- Are you tired of writing CRUD?
- Do you realize there is actually nothing such as One-To-Many in RMDBS and all you do is deceive yourselves with Hibernate and others?
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
- [ ] Generate getters for indexes (pkey and other composite keys)
- [ ] User defined functions
- [ ] Generate functions from references foreign keys
- [ ] Update function
- [ ] gRPC Server implementation

## Example

For given this schema, it generates an almost valid gRPC right now.

SQL:

```
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

message teacher {
    int32 id = 1;
    string name = 2;
}

message course {
    int32 id = 1;
    string title = 2;
    int32 teacherid = 3;
    severity severity = 4;
}

// Service Definition
service DatabaseService {
    rpc GetTeacherById(int32) returns (teacher) {}
    rpc ListCourse(VoidRequest) returns (stream course) {}
    rpc ListMembership(VoidRequest) returns (stream membership) {}
    rpc ListStudent(VoidRequest) returns (stream student) {}
    rpc ListTeacher(VoidRequest) returns (stream teacher) {}
    rpc GetCourseById(int32) returns (course) {}
    rpc GetStudentById(int32) returns (student) {}
}
```