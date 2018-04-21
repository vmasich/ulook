package rest

import (
	"net/http"

	"bitbucket.org/vmasych/urllookup/pkg/model"
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"github.com/coreos/pkg/capnslog"
	"github.com/gin-gonic/gin"
)

var log = capnslog.NewPackageLogger("bitbucket.org/vmasych/urllookup/pkg/rest", "rest")

type Rest struct {
	Model model.MyModel
	Eng   *gin.Engine
}

func New(m model.MyModel) *Rest {
	rest := &Rest{
		Model: m,
	}
	rest.Eng = gin.Default()

	v1 := rest.Eng.Group("urlinfo/1")

	v1.GET("/:hostport/:pathquery", rest.CheckURL)
	v1.POST("/", rest.UpdateURLs)

	return rest
}

func (r *Rest) Run() {
	r.Eng.Run(":3333")
}

func (r *Rest) CheckURL(c *gin.Context) {
	url := schema.MyUrl{
		Host:      c.Param("hostport"),
		PathQuery: c.Param("pathquery"),
	}
	code := http.StatusNotFound
	if r.Model.CheckURL(url) {
		code = http.StatusOK
	}

	c.String(code, "")
}

func (r *Rest) UpdateURLs(c *gin.Context) {
	payload := []schema.UpdateMyUrl{}
	if err := c.BindJSON(&payload); err != nil {
		c.String(http.StatusBadRequest, "")
		log.Errorf("%v", err)
		return
	}
	r.Model.UpdateURLs(payload)
}
