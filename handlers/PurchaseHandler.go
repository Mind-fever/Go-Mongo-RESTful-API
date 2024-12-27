package handlers

import (
	"log"
	"net/http"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/services"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	purchaseService services.PurchaseServiceInterface
}

func NewPurchaseHandler(purchaseService services.PurchaseServiceInterface) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}
func (handler *PurchaseHandler) GetPurchases(c *gin.Context) {
	log.Println("Handler: GetPurchases")

	userInfo := utils.GetUserInfoFromContext(c)

	purchases, err := handler.purchaseService.GetPurchases(userInfo.UserId)
	if err != nil {
		log.Printf("error retrieving purchases: %v", err)
		c.Error(err)
		return
	}

	log.Printf("[handler:PurchaseHandler][method:GetPurchases][cantidad:%d][user:%s]", len(purchases), userInfo)
	c.JSON(http.StatusOK, purchases)
}

func (handler *PurchaseHandler) PostPurchase(c *gin.Context) {
	userInfo := utils.GetUserInfoFromContext(c)

	createdPurchaseDTO, err := handler.purchaseService.PostPurchase(userInfo.UserId)
	if err != nil {
		log.Println("Error creating purchase: ", err)
		c.Error(err)
		return
	}
	log.Println("[handler:PurchaseHandler][method:PostPurchase][user:", userInfo, "]")
	c.JSON(http.StatusCreated, createdPurchaseDTO)
}
