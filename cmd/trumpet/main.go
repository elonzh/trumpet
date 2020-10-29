package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/elonzh/trumpet/transformers"
)

func main() {
	r := gin.Default()
	r.POST("/transformers/:transformer", func(c *gin.Context) {
		transformerName := c.Param("transformer")
		trumpetTo, err := url.Parse(c.Query("trumpet_to"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		transformer, ok := transformers.Get(transformerName)
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

				req, err := httputil.DumpRequest(request, true)
				fmt.Println("-------------------- Request --------------------")
				fmt.Println(string(req))
				if err != nil {
					log.Println(err)
				}

			},
			ModifyResponse: func(response *http.Response) error {
				resp, err := httputil.DumpResponse(response, true)
				fmt.Println("-------------------- Response --------------------")
				fmt.Println(string(resp))
				if err != nil {
					log.Println(err)
				}
				return nil
			},
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Run()
}
