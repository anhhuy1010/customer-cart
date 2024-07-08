package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/anhhuy1010/customer-cart/constant"
	"github.com/anhhuy1010/customer-cart/grpc"
	"github.com/anhhuy1010/customer-cart/helpers/respond"
	"github.com/anhhuy1010/customer-cart/helpers/util"
	"github.com/anhhuy1010/customer-cart/models"
	request "github.com/anhhuy1010/customer-cart/request/cart"
	pbProduct "github.com/anhhuy1010/customer-menu/grpc/proto/product"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
)

type CartController struct {
}

func (cartCtl CartController) GetProductDetail(productUuid string, availableDate string) (*pbProduct.DetailResponse, error) {
	grpcConn := grpc.GetInstance()
	client := pbProduct.NewProductClient(grpcConn.MenuConnect)
	req := pbProduct.DetailRequest{
		Uuid: productUuid,
		Date: availableDate,
	}

	header := metadata.New(map[string]string{})
	ctx := metadata.NewOutgoingContext(context.TODO(), header)

	productDetail, err := client.Detail(ctx, &req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return productDetail, nil
}

func (cartCtr CartController) Create(c *gin.Context) {
	cartModel := new(models.Carts)
	var req request.CartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cartData := &models.Carts{}
	if req.CartUuid == "" {
		cartData.CartUuid = util.GenerateUUID()
		cartData.Total = 0

		if _, err := cartData.Insert(); err != nil {
			c.JSON(http.StatusInternalServerError, respond.UpdatedFail())
			return
		}
	} else {
		cond := bson.M{"cart_uuid": req.CartUuid}
		cartData, err = cartModel.FindOne(cond)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, respond.ErrorCommon("cart no found!"))
			return
		}
	}
	productDetail, err := cartCtr.GetProductDetail(req.ProductUuid, time.Now().Format("2006-01-02"))
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Get product detail error!"))
		return
	}
	cartItem := models.CartItem{
		CartUuid:     cartData.CartUuid,
		ProductUuid:  req.ProductUuid,
		ProductName:  productDetail.Name,
		Quantity:     1,
		ProductPrice: productDetail.Price,
		ProductTotal: 1 * productDetail.Price,
		CartItemUuid: util.GenerateUUID(),
	}
	_, err = cartItem.Insert()
	if err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Cart Items error!"))
		return
	}

	c.JSON(http.StatusOK, respond.Success(cartItem, "Update successfully"))
}
func (cartCtl CartController) Detail(c *gin.Context) {
	cartItemModel := new(models.CartItem)
	var reqUri request.GetCartRequestUri
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	itemCondition := bson.M{"cart_uuid": reqUri.CartUuid}
	cartItems, err := cartItemModel.Find(itemCondition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Cart items not found!"))
		return
	}

	total := 0.0
	var itemm []request.GetCartItemResponse
	for _, item := range cartItems {
		productTotal := item.ProductPrice * float64(item.Quantity)
		total += productTotal

		itemm = append(itemm, request.GetCartItemResponse{
			ProductUuid:  item.ProductUuid,
			CartItemUuid: item.CartItemUuid,
			CartUuid:     item.CartUuid,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     item.Quantity,
			ProductTotal: productTotal,
		})
	}
	response := request.GetCartResponse{
		CartUuid: reqUri.CartUuid,
		Items:    itemm,
		Total:    total,
	}
	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

func (cartCtl CartController) Delete(c *gin.Context) {
	cartModel := new(models.CartItem)
	var reqUri request.DeleteItemUri
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{
		"cart_uuid":      reqUri.CartUuid,
		"cart_item_uuid": reqUri.CartItemUuid}
	cartitemm, err := cartModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.ErrorCommon("cart no found!"))
		return
	}
	cartitemm.IsDelete = constant.DELETE

	_, err = cartitemm.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(cartitemm.CartItemUuid, "Delete successfully"))
}
