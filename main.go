package main

import (
	"fmt"
	"github.com/dearkk/component/market"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"plugin"
)

var store *gorm.DB
var klog *log.Entry
var gWebService *restful.WebService

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	klog = log.WithFields(log.Fields{
		"fields": "market",
	})
	//log.SetReportCaller(true)
}

func initStore(my *market.Mysql) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		my.User, my.Password, my.IP, my.Port, my.Database)
	var err error
	store, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		klog.Errorln("error open mysql: ", err)
		os.Exit(-1)
	}
}

func main() {
	initLog()
	var cfg market.Config
	f, err := os.OpenFile("/Users/kun/workspace/bean/market/conf.yaml", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		klog.Errorln(err.Error())
		os.Exit(-1)
	}
	contentByte, _ := ioutil.ReadAll(f)
	yaml.Unmarshal(contentByte, &cfg)
	klog.Printf("cfg: %+v\n", cfg)

	initStore(&cfg.Mysql)

	container := restful.DefaultContainer
	container.EnableContentEncoding(true)
	gWebService = new(restful.WebService)
	gWebService.Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	addr := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)

	for _, module := range cfg.Modules {
		if module.Enable {
			file := fmt.Sprintf("%s/%s.so", cfg.ModuleDir, module.Name)
			klog.Printf("model: %s, enable: %v, file: %s\n", module.Name, module.Enable, file)
			loadPlugin(module.Name, file, &module.Params)
		} else {
			klog.Printf("model: %s, enable: %v\n", module.Name, module.Enable)
		}
	}

	klog.Infof("bind local address: %s", addr)
	container.Add(gWebService)
	enableSwagger(container, addr)
	err = http.ListenAndServe(addr, nil)
	klog.Infof("market exit! %s", err)
}

func loadPlugin(module string, file string, params *[]market.Param) {
	p, err := plugin.Open(file)
	if err != nil {
		klog.Errorln("error open bin: ", err)
		os.Exit(-1)
	}
	lookup, err := p.Lookup("Load")
	if err != nil {
		klog.Errorln("error lookup Hello: ", err)
		os.Exit(-1)
	}
	klog.Println("lookup: ", file)

	l := log.WithFields(log.Fields{
		"fields": module,
	})

	if load, ok := lookup.(func() market.Load); ok {
		f := load()
		f.Start(store, l, params, addRoute)
		klog.Printf("success to load bin: %s, %v\n", file, ok)
	} else {
		klog.Printf("failed to load bin: %s, %v\n", file, ok)
	}
}

func addRoute(module string, routes []market.Route) {
	for _, route := range routes {
		route.Path = fmt.Sprintf("/%s%s", module, route.Path)
		klog.Infof("addRoute: %s\n", route.Path)
		r := gWebService.
			POST(route.Path).
			To(route.Handle).
			Doc(route.Name).
			Metadata(restfulspec.KeyOpenAPITags, []string{route.Tag})
		if route.Reads != nil {
			r.Reads(route.Reads)
		}
		if route.Writes != nil {
			r.Writes(route.Writes)
		}
		gWebService.Route(r)
	}
}

func enableSwagger(container *restful.Container, addr string) {

	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json"}
	container.Add(restfulspec.NewOpenAPIService(config))

	// Open http://localhost:80/apidocs/?url=http://localhost:80/apidocs.json
	klog.Printf("http://%s/apidocs/?url=http://%s/apidocs.json", addr, addr)
	swaggerPath := "/Users/kun/xh_work/xh-access/build/swagger-ui/dist"
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir(swaggerPath))))

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      container}
	container.Filter(cors.Filter)
}
