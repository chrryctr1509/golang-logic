package handler

import (
	"errors"
	"net/http"

	"github.com/user/tahap2-rest-api/internal/repository"
	"github.com/user/tahap2-rest-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	txService *service.TransactionService
}

func NewTransactionHandler(txService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type PaymentRequest struct {
	Amount  int64  `json:"amount" binding:"required,gt=0"`
	Remarks string `json:"remarks"`
}

type TransferRequest struct {
	TargetUser string `json:"target_user" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	Remarks    string `json:"remarks"`
}

func (h *TransactionHandler) TopUp(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(string)

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
		return
	}

	tx, err := h.txService.TopUp(c.Request.Context(), uid, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"top_up_id":       tx.TransactionID,
			"amount":          tx.Amount,
			"balance_before":  tx.BalanceBefore,
			"balance_after":   tx.BalanceAfter,
		},
	})
}

func (h *TransactionHandler) Payment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(string)

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
		return
	}

	tx, err := h.txService.Payment(c.Request.Context(), uid, req.Amount, req.Remarks)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
			return
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Balance is not enough"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"payment_id":      tx.TransactionID,
			"amount":          tx.Amount,
			"remarks":         tx.Remarks,
			"balance_before":  tx.BalanceBefore,
			"balance_after":   tx.BalanceAfter,
		},
	})
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(string)

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
		return
	}

	tx, err := h.txService.Transfer(c.Request.Context(), uid, req.TargetUser, req.Amount, req.Remarks)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"transfer_id":     tx.TransactionID,
			"amount":          tx.Amount,
			"remarks":         tx.Remarks,
			"balance_before":  tx.BalanceBefore,
			"balance_after":   tx.BalanceAfter,
			"status":          tx.Status,
		},
	})
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(string)

	txns, err := h.txService.GetTransactions(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	results := make([]gin.H, 0, len(txns))
	for _, tx := range txns {
		item := gin.H{
			"amount":        tx.Amount,
			"remarks":       tx.Remarks,
			"balance_before": tx.BalanceBefore,
			"balance_after":  tx.BalanceAfter,
			"status":        tx.Status,
			"created_date":  tx.CreatedDate,
		}
		switch tx.TransactionKind {
		case "TOPUP":
			item["top_up_id"] = tx.TransactionID
		case "PAYMENT":
			item["payment_id"] = tx.TransactionID
		case "TRANSFER":
			item["transfer_id"] = tx.TransactionID
			item["related_user_id"] = tx.RelatedUserID
		}
		results = append(results, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": results,
	})
}
