package main

import (
	"bytes"
	"context"
	"encoding/json"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var router = gin.Default()

const index = "test-web"

var (
	es        *elasticsearch7.Client
	indexes   = []string{index}
	ctx       context.Context
	transport estransport.Interface
)

func init() {
	var err error
	es, err = elasticsearch7.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	transport = es.Transport
	ctx = context.Background()
}

func main() {
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ES Test Web")
	})

	router.GET("/es", func(c *gin.Context) {
		q := `{ "query": { "match_all": {} }}`
		resp, err := esapi.SearchRequest{Index: indexes, Body: strings.NewReader(q)}.Do(ctx, transport)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, string(body))
	})

	router.GET("/es/add", func(c *gin.Context) {

		o, _ := json.Marshal(map[string]string{"hi": "there", "time": time.Now().String()})
		msg := bytes.NewReader(o)
		createRequest := esapi.IndexRequest{Index: index, Body: msg}
		resp, err := createRequest.Do(ctx, transport)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, string(body))
	})

	router.GET("/es/delete/:id", func(c *gin.Context) {
		docID := c.Param("id")
		resp, err := esapi.DeleteRequest{Index: index, DocumentID: docID}.Do(ctx, transport)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, string(body))
	})

	if err := router.Run(); err != nil {
		panic(err)
	}
}
