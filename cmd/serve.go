package cmd

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newServerCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			router := newRouter(cfg)
			return router.Run()
		},
	}
}

func newRouter(cfg *Config) *gin.Engine {
	const realWebHookQuery = "trumpet_to"
	proxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			// already checked
			trumpetTo, _ := url.Parse(request.URL.Query().Get(realWebHookQuery))
			request.Host = trumpetTo.Host
			request.URL = trumpetTo
			request.RequestURI = ""
			request.Header["X-Forwarded-For"] = nil
			request.ContentLength = -1
			delete(request.Header, "Content-Length")
			if logrus.GetLevel() >= logrus.DebugLevel {
				req, err := httputil.DumpRequest(request, true)
				fmt.Printf(
					"\n-------------------- Request --------------------\n%s\nDumpRequestError:%s\n",
					req, err,
				)
			}
		},
		ModifyResponse: func(response *http.Response) error {
			if logrus.GetLevel() >= logrus.DebugLevel {
				resp, err := httputil.DumpResponse(response, true)
				fmt.Printf(
					"\n-------------------- Request --------------------\n%s\nDumpResponseError:%s\n",
					resp, err,
				)
			}
			return nil
		},
	}

	r := gin.Default()
	r.POST("/transformers/:transformer", func(c *gin.Context) {
		if c.ContentType() != binding.MIMEJSON {
			c.String(http.StatusBadRequest, "currently we only accept `%s` content", binding.MIMEJSON)
			return
		}
		transformerName := c.Param("transformer")
		_, err := url.Parse(c.Query(realWebHookQuery))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		transformer, ok := cfg.GetTransformer(transformerName)
		if !ok {
			c.String(http.StatusNotFound, "no such transformer `%s`", transformer)
			return
		}
		req, err := transformer.Exec(c.Request)
		if err != nil {
			c.String(http.StatusInternalServerError, "error when transform data: %s", err)
			return
		}
		proxy.ServeHTTP(c.Writer, req)
	})
	return r
}
