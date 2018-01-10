package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
	"github.com/wallix/awless/sync/repo"
	tstore "github.com/wallix/triplestore"
)

type server struct {
	port string
	gph  cloud.GraphAPI
}

func New(port string) *server {
	return &server{port: port}
}

func (s *server) Start() error {
	g, err := sync.LoadAllLocalGraphs(config.GetAWSProfile())
	if err != nil {
		return fmt.Errorf("cannot load local graphs: %s", err)
	}

	s.gph = g

	log.Printf("Starting browsing on http://localhost%s\n", s.port)
	return http.ListenAndServe(s.port, s.routes())
}

func (s *server) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/resources/{id}", s.showResourceHandler)
	r.HandleFunc("/resources", s.listResourcesHandler)
	r.HandleFunc("/rdf", s.rdfHandler)
	r.HandleFunc("/graph", s.graphHandler)
	r.HandleFunc("/", s.homeHandler)
	return r
}

func (s *server) homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("home").Parse(homeTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) rdfHandler(w http.ResponseWriter, r *http.Request) {
	tris, err := loadLocalTriples()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var encErr error
	if r.FormValue("namespaced") == "true" {
		encErr = tstore.NewLenientNTEncoderWithContext(w, tstore.RDFContext).Encode(tris...)
	} else {
		encErr = tstore.NewLenientNTEncoder(w).Encode(tris...)
	}

	if encErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) graphHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("graph").Parse(graphVizTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tris, err := loadLocalTriples()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data bytes.Buffer
	if err := tstore.NewDotGraphEncoder(&data, "cloud-rel:parentOf").Encode(tris...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) showResourceHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("show").Parse(showResourceTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resId := mux.Vars(r)["id"]
	res, err := s.gph.FindWithProperties(map[string]interface{}{properties.ID: resId})
	if err != nil && len(res) != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resource := newResource(res[0])
	deps, _ := s.gph.ResourceRelations(res[0], rdf.DependingOnRel, false)
	resource.AddDependsOn(deps...)
	applies, _ := s.gph.ResourceRelations(res[0], rdf.ApplyOn, false)
	resource.AddAppliesOn(applies...)

	parents, _ := s.gph.ResourceRelations(res[0], rdf.ParentOf, true)
	resource.AddDependsOn(parents...)

	children, _ := s.gph.ResourceRelations(res[0], rdf.ChildrenOfRel, true)
	resource.AddDependsOn(children...)

	if err := t.Execute(w, resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) listResourcesHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("resources").Parse(resourcesTpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resourcesByTypes := make(map[string][]*Resource)

	for _, typ := range append(awsservices.ResourceTypes, "region") {
		gRes, err := s.gph.Find(cloud.NewQuery(typ))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, r := range gRes {
			resourcesByTypes[typ] = append(resourcesByTypes[typ], newResource(r))
		}
	}

	if err := t.Execute(w, resourcesByTypes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Resource struct {
	Id, Type   string
	Properties map[string]interface{}
	Parents    []*Resource
	Children   []*Resource
	DependsOn  []*Resource
	AppliesOn  []*Resource
}

func (r *Resource) AddDependsOn(gr ...cloud.Resource) {
	for _, res := range gr {
		r.DependsOn = append(r.DependsOn, newResource(res))
	}
}

func (r *Resource) AddAppliesOn(gr ...cloud.Resource) {
	for _, res := range gr {
		r.AppliesOn = append(r.AppliesOn, newResource(res))
	}
}

func (r *Resource) AddParents(gr ...cloud.Resource) {
	for _, res := range gr {
		r.Parents = append(r.Parents, newResource(res))
	}
}

func (r *Resource) AddChildren(gr ...*graph.Resource) {
	for _, res := range gr {
		r.Children = append(r.Children, newResource(res))
	}
}

func newResource(r cloud.Resource) *Resource {
	return &Resource{Id: r.Id(), Type: r.Type(), Properties: r.Properties()}
}

func loadLocalTriples() ([]tstore.Triple, error) {
	path := filepath.Join(repo.BaseDir(), "*", fmt.Sprintf("*%s", ".triples"))
	files, _ := filepath.Glob(path)

	var readers []io.Reader
	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			return nil, fmt.Errorf("loading '%s': %s", f, err)
		}
		readers = append(readers, reader)
	}

	dec := tstore.NewDatasetDecoder(tstore.NewAutoDecoder, readers...)
	return dec.Decode()
}

const homeTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
	<ul>
	<li><a href="/resources">View resources and their relations</a></li>
	<li><a href="/rdf">View RDF</a></li>
	<li><a href="/rdf?namespaced=true">View namespaced RDF</a></li>
	<li><a href="/graph">View DOT graph (experimental)</a></li>
	</ul>
	</body>
</html>`

const showResourceTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
	<h2>{{.Type}}: {{.Id}}</h2>
	<ul>
	{{range $name, $val := .Properties}}
	  <li><b>{{$name}}:</b> {{$val}}</li>
        {{end}}
        </ul>

	{{if (len .Parents) gt 0}}
	<h4>Parents:</h4>
	<ul>
	{{range .Parents}}
	  <li>{{.Type}} <a href="/resources/{{ urlquery .Id}}">{{.Id}}</a></li>
	{{end}}
	</ul>
	{{end}}

	{{if (len .Children) gt 0}}
	<h4>Children:</h4>
	<ul>
	{{range .Children}}
	  <li>{{.Type}} <a href="/resources/{{ urlquery .Id}}">{{.Id}}</a></li>
	{{end}}
	</ul>
	{{end}}


	{{if (len .DependsOn) gt 0}}
	<h4>Depends on:</h4>
	<ul>
	{{range .DependsOn}}
	  <li>{{.Type}} <a href="/resources/{{ urlquery .Id}}">{{.Id}}</a></li>
	{{end}}
	</ul>
	{{end}}

	{{if (len .AppliesOn) gt 0}}
	<h4>Applies on:</h4>
	<ul>
	{{range .AppliesOn}}
	  <li>{{.Type}} <a href="/resources/{{ urlquery .Id}}">{{.Id}}</a></li>
	{{end}}
	</ul>
	{{end}}
	</body>
</html>`

const resourcesTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		{{range $type, $resource := .}}
		<h2>{{$type}} ({{len .}})</h2>
		<ul>
		  {{range $resource}}
		         {{ $name := index .Properties "Name" }}
			 <li>
			  {{if $name}} {{if (ne (print $name) "")}}<b>Name:</b> {{$name}}, {{end}}{{end}}
			  <b>Id: </b><a href="/resources/{{ urlquery .Id}}">{{.Id}}</a>
			 </li>
		  {{end}}
		</ul>
		{{end}}
	</body>
</html>`

const graphVizTpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
        <script src="https://github.com/mdaines/viz.js/releases/download/v1.7.1/viz-lite.js"></script>
	<script>
	   document.body.innerHTML += Viz("{{ . }}");
	</script>
	</body>
</html>`
