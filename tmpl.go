package main

import (
	"bytes"
	"log"
	"text/template"
)

const tmplText = `
package main

import (
	"log"

	"golang.org/x/net/context"
	_ "github.com/lib/pq"

	pb "../." // The relative import so we do not need to specify an absolute path
)

type server struct {
	db *sql.DB
}

{{range .Fns}}
func(d *server)	{{.Name}}(c context.Context, arg *pb.{{.InputName}}) (*pb.{{.OutputName}}, error){
	var sql = "{{.Function.Query}}"
	{{ if .Function.IsOutArray }}
	rows, err := d.db.Query(sql, {{ range .Function.Inputs }} arg.{{.Name}}, {{ end }})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	{{ else }}
	row, err := d.db.QueryRow(sql, {{ range .Function.Inputs }} arg.{{.Name}}, {{ end }})
	if err != nil {

	}
	out := &pb.{{.OutputName}}
	err := row.Scan({{ range .Function.Outputs }} &out.{{.Name}}, {{ end }}))
	defer rows.Close()
	{{ end}}
}
{{end}}
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

`

var tmpl *template.Template

func init() {
	_tmpl := template.New("gopqgen-template")
	_tmpl, err := _tmpl.Parse(tmplText)

	tmpl = _tmpl // re assign to global variable
	if err != nil {
		log.Fatal("Could not parse the template", err)
	}
}

func generateTemplate(summary protoSummary) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, &summary)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println(buf.String())

	return buf.String(), nil
}
