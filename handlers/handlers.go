package handlers

import (
	"github.com/credo-science/credo-metadata-server/event"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/credo-science/credo-metadata-server/db"
)

func GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func GetEventAll(c *gin.Context) {
	eventType, ok := event.StringToId(c.Params.ByName("event_type"))

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event type."})
		return
	}

	eventId := c.Params.ByName("event_id")
	metadata, err := db.GetEventMetadata(eventType, eventId)
	if err != nil {
		if err == db.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event metadata from database."})
			log.Fatal(err)
			return
		}
	}
	c.JSON(http.StatusOK, metadata)
}

func PutEventAll(c *gin.Context) {
	eventType, ok := event.StringToId(c.Params.ByName("event_type"))

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event type."})
		return
	}

	eventId := c.Params.ByName("event_id")

	metadata := event.Metadata{}
	err := c.BindJSON(&metadata)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid metadata object."})
		return
	}

	err = db.SetEventMetadata(eventType, eventId, metadata)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event metadata."})
		log.Panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Metadata updated."})
}
