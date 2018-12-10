package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	contactsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "counter_contacts",
			Help: "total of requests to endpoint /contact",
		},
		[]string{"status", "method"},
	)
	contactCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "counter_contact",
			Help: "total of requests to endpoint /contact/:id",
		},
		[]string{"status", "method"},
	)
)

/*----------------------------------------------------------------------------------------*/

func InitAndStartRouter(db *DB) {
	gin.SetMode("debug")

	router := gin.New()
	router.Use(CorsMiddleware())
	router.Use(ginrus.Ginrus(log.New(), time.RFC3339, false))
	router.Use(gin.Recovery())
	router.Use(static.Serve("/", static.LocalFile("static", true)))
	router.GET("/api/contact", db.Contacts)
	router.GET("/api/contact/:id", db.Contact)
	router.POST("/api/contact", db.CreateContact)
	router.PUT("/api/contact/:id", db.ModifyContact)
	router.DELETE("/api/contact/:id", db.RemoveContact)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run(":8080")

}

/*----------------------------------------------------------------------------------------*/

func (db *DB) Contacts(c *gin.Context) {
	contacts, err := db.FindAllContacts()
	if err == nil {
		contactsCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusOK), "method": http.MethodGet}).Inc()
		c.JSON(http.StatusOK, contacts)
		return
	}
	contactsCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusInternalServerError), "method": http.MethodGet}).Inc()
	respondWithError(http.StatusInternalServerError, "Cound not read contacts", c)
}

func (db *DB) Contact(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorf("Error parsing id: %v", err)
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodGet}).Inc()
		respondWithError(http.StatusBadRequest, "invalid id", c)
	}
	contact, err := db.FindContactById(id)
	if err != nil {
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusNotFound), "method": http.MethodGet}).Inc()
		respondWithError(http.StatusNotFound, "contact not found", c)
		return
	}
	contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusOK), "method": http.MethodGet}).Inc()
	c.JSON(http.StatusOK, contact)
}

func (db *DB) CreateContact(c *gin.Context) {
	contact := &Contact{}
	err := c.BindJSON(contact)
	if err != nil {
		log.Warn("cannot decode request", err)
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodPost}).Inc()
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	contact, err = db.InsertContact(contact)
	if err != nil {
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodPost}).Inc()
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusCreated), "method": http.MethodPost}).Inc()
	c.JSON(http.StatusCreated, contact)
}

func (db *DB) ModifyContact(c *gin.Context) {
	contact := &Contact{}
	err := c.BindJSON(contact)
	if err != nil {
		log.Warn("cannot decode request", err)
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodPut}).Inc()
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	contact, err = db.UpdateContact(contact)
	if err != nil {
		respondWithError(http.StatusBadRequest, "bad request", c)
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodPut}).Inc()
		return
	}
	contactsCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusOK), "method": http.MethodPut}).Inc()
	c.JSON(http.StatusOK, contact)
}

func (db *DB) RemoveContact(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorf("Error parsing id: %v", err)
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodDelete}).Inc()
		respondWithError(http.StatusBadRequest, "invalid id", c)
	}
	if db.DeleteContact(id) != nil {
		contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusBadRequest), "method": http.MethodDelete}).Inc()
		respondWithError(http.StatusBadRequest, "cannot delete contact", c)
		return
	}
	contactCounter.With(prometheus.Labels{"status": strconv.Itoa(http.StatusOK), "method": http.MethodDelete}).Inc()
	c.JSON(http.StatusOK, "")
}

/*----------------------------------------------------------------------------------------*/

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD")
		if c.Request.Method == "OPTIONS" {
			c.Data(http.StatusOK, "text/plain", nil)
			c.Abort()
		}
		c.Next()
	}
}
