package main

import (
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

/*----------------------------------------------------------------------------------------*/

func InitAndStartRouter(db *DB) {
	gin.SetMode("debug")

	router := gin.New()

	router.Use(ginrus.Ginrus(log.New(), time.RFC3339, false))
	router.Use(gin.Recovery())
	router.GET("/api/contact", db.Contacts)
	router.GET("/api/contact/:id", db.Contact)
	router.POST("/api/contact", db.CreateContact)
	router.PUT("/api/contact/:id", db.ModifyContact)
	router.DELETE("/api/contact/:id", db.RemoveContact)
	router.Run(":8080")
}

/*----------------------------------------------------------------------------------------*/

func (db *DB) Contacts(c *gin.Context) {
	contacts, err := db.FindAllContacts()
	if err == nil {
		c.JSON(http.StatusOK, contacts)
		return
	}
	respondWithError(http.StatusInternalServerError, "Cound not read contacts", c)
}

func (db *DB) Contact(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorf("Error parsing id: %v", err)
		respondWithError(http.StatusBadRequest, "invalid id", c)
	}
	contact, err := db.FindContactById(id)
	if err != nil {
		respondWithError(http.StatusNotFound, "contact not found", c)
		return
	}
	c.JSON(http.StatusOK, contact)
}

func (db *DB) CreateContact(c *gin.Context) {
	contact := &Contact{}
	err := c.BindJSON(contact)
	if err != nil {
		log.Warn("cannot decode request", err)
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	contact, err = db.InsertContact(contact)
	if err != nil {
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	c.JSON(http.StatusCreated, contact)
}

func (db *DB) ModifyContact(c *gin.Context) {
	contact := &Contact{}
	err := c.BindJSON(contact)
	if err != nil {
		log.Warn("cannot decode request", err)
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	contact, err = db.UpdateContact(contact)
	if err != nil {
		respondWithError(http.StatusBadRequest, "bad request", c)
		return
	}
	c.JSON(http.StatusOK, contact)
}

func (db *DB) RemoveContact(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Errorf("Error parsing id: %v", err)
		respondWithError(http.StatusBadRequest, "invalid id", c)
	}
	if db.DeleteContact(id) != nil {
		respondWithError(http.StatusBadRequest, "cannot delete contact", c)
		return
	}
	c.JSON(http.StatusOK, "")
}

/*----------------------------------------------------------------------------------------*/

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}
