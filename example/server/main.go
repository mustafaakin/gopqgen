package main

import (
	"log"

	"golang.org/x/net/context"
	_ "github.com/lib/pq"

	pb "../." // The relative import so we do not need to specify an absolute path
	"flag"
	"database/sql"
)

type server struct {
	db *sql.DB
}


func(d *server)	Add(c context.Context, arg *pb.AddArg) (*pb.AddOut, error){
	var sql = "SELECT * FROM add($1, $2)"

	row, err := d.db.QueryRow(sql,  arg.Arg0,  arg.Arg1, )
	if err != nil {

	}
	out := &pb.AddOut{}
	err := row.Scan( &out.Out, ))
	defer rows.Close()

}

func(d *server)	GetCourseById(c context.Context, arg *pb.GetCourseByIdArg) (*pb.GetCourseByIdOut, error){
	var sql = "SELECT id, title, teacherid, severity FROM course WHERE `id` = $1"

	row, err := d.db.QueryRow(sql,  arg.id, )
	if err != nil {

	}
	out := &pb.GetCourseByIdOut
	err := row.Scan( &out.output, ))
	defer rows.Close()

}

func(d *server)	GetStudentById(c context.Context, arg *pb.GetStudentByIdArg) (*pb.GetStudentByIdOut, error){
	var sql = "SELECT id, name, email, password, age FROM student WHERE `id` = $1"

	row, err := d.db.QueryRow(sql,  arg.id, )
	if err != nil {

	}
	out := &pb.GetStudentByIdOut
	err := row.Scan( &out.output, ))
	defer rows.Close()

}

func(d *server)	GetTeacherById(c context.Context, arg *pb.GetTeacherByIdArg) (*pb.GetTeacherByIdOut, error){
	var sql = "SELECT id, name FROM teacher WHERE `id` = $1"

	row, err := d.db.QueryRow(sql,  arg.id, )
	if err != nil {

	}
	out := &pb.GetTeacherByIdOut
	err := row.Scan( &out.output, ))
	defer rows.Close()

}

func(d *server)	Getstudentsofcourse(c context.Context, arg *pb.GetstudentsofcourseArg) (*pb.GetstudentsofcourseOut, error){
	var sql = "SELECT * FROM getstudentsofcourse($1)"

	rows, err := d.db.Query(sql,  arg._courseid, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

}

func(d *server)	Isuserincourse(c context.Context, arg *pb.IsuserincourseArg) (*pb.IsuserincourseOut, error){
	var sql = "SELECT * FROM isuserincourse($1, $2)"

	row, err := d.db.QueryRow(sql,  arg._studentid,  arg._courseid, )
	if err != nil {

	}
	out := &pb.IsuserincourseOut
	err := row.Scan( &out.Out, ))
	defer rows.Close()

}

func(d *server)	ListCourse(c context.Context, arg *pb.VoidArg) (*pb.ListCourseOut, error){
	var sql = "SELECT id, title, teacherid, severity FROM course"

	rows, err := d.db.Query(sql, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

}

func(d *server)	ListMembership(c context.Context, arg *pb.VoidArg) (*pb.ListMembershipOut, error){
	var sql = "SELECT courseid, studentid FROM membership"

	rows, err := d.db.Query(sql, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

}

func(d *server)	ListStudent(c context.Context, arg *pb.VoidArg) (*pb.ListStudentOut, error){
	var sql = "SELECT id, name, email, password, age FROM student"

	rows, err := d.db.Query(sql, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

}

func(d *server)	ListTeacher(c context.Context, arg *pb.VoidArg) (*pb.ListTeacherOut, error){
	var sql = "SELECT id, name FROM teacher"

	rows, err := d.db.Query(sql, )
	if err != nil {
		return nil, err
	}
	defer rows.Close()

}

var dsn  = flag.String("dsn","user=postgres dbname=gopqgen sslmode=disable","The data source name, like how to connect to db")
var port = flag.Int("port", 3000, "port to serve requests from")

func main(){
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatal("Could not connect to Postgres server:", err)
	}
	s := &server{db:db}

	grpc.serve...etc
}
