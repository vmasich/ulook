package rest

import (
	"fmt"
	"net/http"

	"bitbucket.org/vmasych/urllookup/pkg/model"
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"github.com/coreos/pkg/capnslog"
	"github.com/gin-gonic/gin"
)

var log = capnslog.NewPackageLogger("bitbucket.org/vmasych/urllookup/pkg/rest", "rest")

type Rest struct {
	Model model.Operations
	Eng   *gin.Engine
}

func New(m model.Operations) *Rest {
	rest := &Rest{
		Model: m,
	}
	rest.Eng = gin.Default()

	v1 := rest.Eng.Group("urlinfo/")

	v1.GET("1/:hostport/:pathquery", rest.CheckURL)
	v1.POST("bulkupdate", rest.UpdateURLs)

	return rest
}

func (r *Rest) Run() {
	r.Eng.Run(":3333")
}

func (r *Rest) CheckURL(c *gin.Context) {

	path := c.Param("pathquery")
	query := c.Request.URL.RawQuery

	url := schema.LookupURL{
		Host:      c.Param("hostport"),
		PathQuery: fmt.Sprintf("%s?%s", path, query),
	}

	log.Infof("url: %v, %v", url, query)
	code := http.StatusNotFound
	msg := ""

	ok, err := r.Model.LookupURL(url)
	if ok {
		code = http.StatusOK
	}
	if err != nil {
		code = http.StatusGatewayTimeout
		msg = err.Error()
	}
	c.String(code, msg)

}

func (r *Rest) UpdateURLs(c *gin.Context) {
	payload := []schema.UpdLookupURL{}
	if err := c.BindJSON(&payload); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Errorf("%v", err)
		return
	}
	r.Model.UpdateURLs(payload)
}
