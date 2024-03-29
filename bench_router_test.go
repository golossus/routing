package routing

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"
)

var testRoutes = []string{
	"/",
	"/cmd.html",
	"/code.html",
	"/contrib.html",
	"/contribute.html",
	"/debugging_with_gdb.html",
	"/docs.html",
	"/effective_go.html",
	"/files.log",
	"/gccgo_contribute.html",
	"/gccgo_install.html",
	"/go-logo-black.png",
	"/go-logo-blue.png",
	"/go-logo-white.png",
	"/go1.1.html",
	"/go1.2.html",
	"/go1.html",
	"/go1compat.html",
	"/go_faq.html",
	"/go_mem.html",
	"/go_spec.html",
	"/help.html",
	"/ie.css",
	"/install-source.html",
	"/install.html",
	"/logo-153x55.png",
	"/Makefile",
	"/root.html",
	"/share.png",
	"/sieve.gif",
	"/tos.html",
	"/articles",
	"/articles/go_command.html",
	"/articles/index.html",
	"/articles/wiki",
	"/articles/wiki/edit.html",
	"/articles/wiki/final-noclosure.go",
	"/articles/wiki/final-noerror.go",
	"/articles/wiki/final-parsetemplate.go",
	"/articles/wiki/final-template.go",
	"/articles/wiki/final.go",
	"/articles/wiki/get.go",
	"/articles/wiki/http-sample.go",
	"/articles/wiki/index.html",
	"/articles/wiki/Makefile",
	"/articles/wiki/notemplate.go",
	"/articles/wiki/part1-noerror.go",
	"/articles/wiki/part1.go",
	"/articles/wiki/part2.go",
	"/articles/wiki/part3-errorhandling.go",
	"/articles/wiki/part3.go",
	"/articles/wiki/test.bash",
	"/articles/wiki/test_edit.good",
	"/articles/wiki/test_Test.txt.good",
	"/articles/wiki/test_view.good",
	"/articles/wiki/view.html",
	"/codewalk",
	"/codewalk/codewalk.css",
	"/codewalk/codewalk.js",
	"/codewalk/codewalk.xml",
	"/codewalk/functions.xml",
	"/codewalk/markov.go",
	"/codewalk/markov.xml",
	"/codewalk/pig.go",
	"/codewalk/popout.png",
	"/codewalk/run",
	"/codewalk/sharemem.xml",
	"/codewalk/urlpoll.go",
	"/devel",
	"/devel/release.html",
	"/devel/weekly.html",
	"/gopher",
	"/gopher/appenginegopher.jpg",
	"/gopher/appenginegophercolor.jpg",
	"/gopher/appenginelogo.gif",
	"/gopher/bumper.png",
	"/gopher/bumper192x108.png",
	"/gopher/bumper320x180.png",
	"/gopher/bumper480x270.png",
	"/gopher/bumper640x360.png",
	"/gopher/doc.png",
	"/gopher/frontpage.png",
	"/gopher/gopherbw.png",
	"/gopher/gophercolor.png",
	"/gopher/gophercolor16x16.png",
	"/gopher/help.png",
	"/gopher/pkg.png",
	"/gopher/project.png",
	"/gopher/ref.png",
	"/gopher/run.png",
	"/gopher/talks.png",
	"/gopher/pencil",
	"/gopher/pencil/gopherhat.jpg",
	"/gopher/pencil/gopherhelmet.jpg",
	"/gopher/pencil/gophermega.jpg",
	"/gopher/pencil/gopherrunning.jpg",
	"/gopher/pencil/gopherswim.jpg",
	"/gopher/pencil/gopherswrench.jpg",
	"/play",
	"/play/fib.go",
	"/play/hello.go",
	"/play/life.go",
	"/play/peano.go",
	"/play/pi.go",
	"/play/sieve.go",
	"/play/solitaire.go",
	"/play/tree.go",
	"/progs",
	"/progs/cgo1.go",
	"/progs/cgo2.go",
	"/progs/cgo3.go",
	"/progs/cgo4.go",
	"/progs/defer.go",
	"/progs/defer.out",
	"/progs/defer2.go",
	"/progs/defer2.out",
	"/progs/eff_bytesize.go",
	"/progs/eff_bytesize.out",
	"/progs/eff_qr.go",
	"/progs/eff_sequence.go",
	"/progs/eff_sequence.out",
	"/progs/eff_unused1.go",
	"/progs/eff_unused2.go",
	"/progs/error.go",
	"/progs/error2.go",
	"/progs/error3.go",
	"/progs/error4.go",
	"/progs/go1.go",
	"/progs/gobs1.go",
	"/progs/gobs2.go",
	"/progs/image_draw.go",
	"/progs/image_package1.go",
	"/progs/image_package1.out",
	"/progs/image_package2.go",
	"/progs/image_package2.out",
	"/progs/image_package3.go",
	"/progs/image_package3.out",
	"/progs/image_package4.go",
	"/progs/image_package4.out",
	"/progs/image_package5.go",
	"/progs/image_package5.out",
	"/progs/image_package6.go",
	"/progs/image_package6.out",
	"/progs/interface.go",
	"/progs/interface2.go",
	"/progs/interface2.out",
	"/progs/json1.go",
	"/progs/json2.go",
	"/progs/json2.out",
	"/progs/json3.go",
	"/progs/json4.go",
	"/progs/json5.go",
	"/progs/run",
	"/progs/slices.go",
	"/progs/timeout1.go",
	"/progs/timeout2.go",
	"/progs/update.bash",
}

func BenchmarkRouter(b *testing.B) {
	benchRouter(b, false)
}
func BenchmarkRouter_SortedByWeight(b *testing.B) {
	benchRouter(b, true)
}

func BenchmarkRouter_GenerateURL(b *testing.B) {
	mainRouter := Router{}

	_ = mainRouter.Register(http.MethodGet, "/path1", testHandlerFunc)
	_ = mainRouter.As("path1.id.name").Register(http.MethodGet, "/path1/{id}/{name:[a-z]{1,5}}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/path1/{file:.*}", testHandlerFunc)
	_ = mainRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)

	postsRouter := Router{}
	_ = postsRouter.Register(http.MethodGet, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", testHandlerFunc)

	_ = mainRouter.Prefix("/posts/{id}", &postsRouter)

	paramsKey := newURLParameterBag(2)
	paramsKey.add("id", "2")
	paramsKey.add("name", "john")
	for i := 0; i < b.N; i++ {
		_, _ = mainRouter.GenerateURL("path1.id.name", paramsKey)
	}
}

func benchRouter(b *testing.B, prioritizeByWeight bool) {
	router := Router{}
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	before := m.HeapAlloc
	handler := func(response http.ResponseWriter, request *http.Request) {}
	for _, routes := range testRoutes {
		_ = router.Register(http.MethodGet, routes, handler)
	}
	if prioritizeByWeight {
		router.PrioritizeByWeight()
	}
	runtime.ReadMemStats(m)
	after := m.HeapAlloc

	st := reflect.ValueOf(b)
	value := st.MethodByName("ReportMetric")

	if !value.IsNil() {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(float64(after - before))
		inputs[1] = reflect.ValueOf("memory")
		value.Call(inputs)
	}

	request1, _ := http.NewRequest(http.MethodGet, "/play", nil)
	request2, _ := http.NewRequest(http.MethodGet, "/articles/wiki", nil)
	request3, _ := http.NewRequest(http.MethodGet, "/gopher/pencil/gopherswrench.jpg", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, request1)
		router.ServeHTTP(nil, request2)
		router.ServeHTTP(nil, request3)
	}
}
