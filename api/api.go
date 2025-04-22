package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	// "strconv"
	// "strings"
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

	r.POST("/customer", h.handleCreateCustomer)
	r.GET("/customer/:id", h.handleGetCustomerById)
	r.DELETE("/customer/:id/delete", h.handleDeleteCustomer)
	r.POST("/order", h.handleCreateOrder)
	r.PUT("/order/updatestatus", h.handleUpdateOrderStatus)

	return r
}

// Create Customer
func (h *MessageHandler) handleCreateCustomer(c *gin.Context) {
	var req repo.CreateCustomerParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.querier.CreateCustomer(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, customer)
}

// Delete Customer
func (h *MessageHandler) handleDeleteCustomer(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	// Check if the customer exist
	_, err := h.querier.GetCustomerById(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify customer's existence"})
		}
		return
	}

	// Proceed to delete
	err = h.querier.DeleteCustomer(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})


}

// Get Customer By ID
func (h *MessageHandler) handleGetCustomerById(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	customer, err := h.querier.GetCustomerById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Create Order
func (h *MessageHandler) handleCreateOrder(c *gin.Context) {
	var req repo.CreateOrderParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.querier.CreateOrder(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, order)
}

// Update Status
func (h *MessageHandler) handleUpdateOrderStatus(c*gin.Context) {

	var req repo.UpdateOrderStatusParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	 if err = h.querier.UpdateOrderStatus(c,req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		} else { 
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	    return
	}

	c.JSON(http.StatusOK,"updated successfully")
}