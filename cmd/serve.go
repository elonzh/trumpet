package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		r.POST("/transformers/:transformer", func(c *gin.Context) {
			transformerName := c.Param("transformer")
			trumpetTo, err := url.Parse(c.Query("trumpet_to"))
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			transformer, ok := cfg.Transformers[transformerName]
			if !ok {
				c.String(http.StatusNotFound, "no such transformer `%s`", transformer)
				return
			}
			raw, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			rawBody, err := transformer.Exec(string(raw))
			if err != nil {
				c.String(http.StatusInternalServerError, "error when transform data: %s", err)
				return
			}
			proxy := httputil.ReverseProxy{
				Director: func(request *http.Request) {
					request.Host = trumpetTo.Host
					request.URL = trumpetTo
					request.RequestURI = ""
					request.Header["X-Forwarded-For"] = nil
					request.ContentLength = -1
					delete(request.Header, "Content-Length")

					request.Body = ioutil.NopCloser(strings.NewReader(rawBody))
					if cfg.LogLevel >= logrus.DebugLevel {
						req, err := httputil.DumpRequest(request, true)
						fmt.Printf(
							"\n-------------------- Request --------------------\n%s\nDumpRequestError:%s\n",
							req, err,
						)
					}
				},
				ModifyResponse: func(response *http.Response) error {
					if cfg.LogLevel >= logrus.DebugLevel {
						resp, err := httputil.DumpResponse(response, true)
						fmt.Printf(
							"\n-------------------- Request --------------------\n%s\nDumpResponseError:%s\n",
							resp, err,
						)
					}
					return nil
				},
			}
			proxy.ServeHTTP(c.Writer, c.Request)
		})
		return r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
