package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	querier repo.Querier
}

func NewMessageHandler(querier repo.Querier) *MessageHandler {
	return &MessageHandler{
		querier: querier,
	}
}

// Register the endpoints
func (h *MessageHandler) WireHttpHandler() http.Handler {
	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.POST("/message", h.handleCreateMessage)
	r.GET("/message/:id", h.handleGetMessage)
	r.GET("/thread/:id/messages", h.handleGetThreadMessages)// Example of paginated : GET /thread/6e32eb10-1cbe-4f09-a82e-f706dfe9cf93/messages?limit=10&offset=20
	r.DELETE("/message/:id", h.handleDeleteMessageById)
	r.PATCH("/message/:id", h.handleEditMessage)
	r.POST("/thread", h.handleCreateThread)
	r.GET("/thread/:id", h.handlerGetThreadID)
	r.GET("/thread/:id/messages/count", h.handlerCountMessagesInThread)
	r.GET("/thread/:id/messages/latest", h.handleGetLatestMessageInThread)
	r.GET("/thread/:id/messages/search", h.handleSearchMessagesByKeyword) // GET /thread/abc123/messages/search?keyword=hello
	r.DELETE("/thread/:id", h.handleDeleteThread)
	r.GET("/thread/:id/messages/page", h.handleGetMessagesByThreadPaginatedHandler)


	return r
}

// Create a message
func (h *MessageHandler) handleCreateMessage(c *gin.Context) {
	var req repo.CreateMessageParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the thread exist
	if _,err := h.querier.GetThreadID(c,req.Thread); err != nil{
		if errors.Is(err,sql.ErrNoRows){
			//dont continue if no row is found
			c.JSON(http.StatusBadRequest,gin.H{"error":"Thread id not found"})
			return
		}
		//return if there was an error encountered running this query
		c.JSON(http.StatusInternalServerError,err.Error())
	}

	//if the thread exist, then proceed to create the message

	message, err := h.querier.CreateMessage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// Get a message
func (h *MessageHandler) handleGetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// Get message by thread
func (h *MessageHandler) handleGetThreadMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	messages, err := h.querier.GetMessagesByThread(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread":   id,
		"topic":    "example",
		"messages": messages,
	})
}

// Delete message
func (h *MessageHandler) handleDeleteMessageById(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	// Check if the message exists
	_, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify message existence"})
		}
		return
	}

	// Proceed to delete
	err = h.querier.DeleteMessage(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// Funtion to Edit Content of Message
func (h *MessageHandler) handleEditMessage(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))

	// If the id in the endpoint is present
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req repo.EditMessageParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Doing the editing of the content in a message
	message, err := h.querier.EditMessage(c, id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		} else { 
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	    return
	}

	c.JSON(http.StatusOK, message)

}

// Create a Thread
func (h *MessageHandler) handleCreateThread(c *gin.Context) {
	var req repo.CreateThreadParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.querier.CreateThread(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// Get Thread by ID
func (h *MessageHandler) handlerGetThreadID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetThreadID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// Count Messages in a thread
func (h *MessageHandler) handlerCountMessagesInThread(c *gin.Context) {
	threadID := strings.TrimSpace(c.Param("id")) // The thread here is the thread's id

	if threadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread id is required"})
		return
	}

	count, err := h.querier.CountMessagesInThread(c, threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread_id": threadID,
		"message_count": count,
	})
}

// Get Latest Message in a Thread
func (h *MessageHandler) handleGetLatestMessageInThread(c *gin.Context) {
	threadID := strings.TrimSpace(c.Param("id"))
	if threadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread ID is required"})
		return
	}

	// Check if the thread exists
	if _, err := h.querier.GetThreadID(c, threadID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the latest message
	message, err := h.querier.GetLatestMessageInThread(c, threadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No messages found in this thread"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, message)
}


// Search messages by keyword
func (h *MessageHandler) handleSearchMessagesByKeyword(c *gin.Context) {
	threadID := strings.TrimSpace(c.Param("id"))
	keyword := strings.TrimSpace(c.Query("keyword"))

	if threadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread ID is required"})
		return
	}

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword is required"})
		return
	}

	// Check if thread exists
	if _, err := h.querier.GetThreadID(c, threadID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Perform the search
	messages, err := h.querier.SearchMessagesByKeyword(c, repo.SearchMessagesByKeywordParams{
		Thread:  threadID,
		Column2: keyword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread_id": threadID,
		"keyword":   keyword,
		"matches":   messages,
	})
}

// Delete a thread
func (h *MessageHandler) handleDeleteThread(c *gin.Context) {
	threadID := strings.TrimSpace(c.Param("id"))

	if threadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread ID is required"})
		return
	}

	if err := h.querier.DeleteThread(c, threadID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Thread and messages deleted"})
}

// Get Messade by thread paginated
func (h *MessageHandler) handleGetMessagesByThreadPaginatedHandler(c *gin.Context) {
		thread := c.Param("id")
		if thread == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread is required"})
			return
		}

		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		req := repo.GetMessagesByThreadPaginatedParams{
			Thread: thread,
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		messages, err := h.querier.GetMessagesByThreadPaginated(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
