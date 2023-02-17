package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vandesar/entity"
	"vandesar/service"

	"gorm.io/gorm"
)

type ProductAPI struct {
	productService    *service.ProductService
	userService *service.UserService
}

func NewProductAPI(
	productService *service.ProductService,
	userService *service.UserService,
) *ProductAPI {
	return &ProductAPI{
		productService:    productService,
		userService:       userService,
	}
}

func (p *ProductAPI) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)

	product := r.URL.Query()
	productID, foundProductId := product["product_id"]
	productSearch, foundProductSearch := product["search"]

	if foundProductId {
		pID, _ := strconv.Atoi(productID[0])
		productByID, err := p.productService.GetProductByID(r.Context(), pID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if productByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error product not found"))
			return
		}

		if productByID.UserID != adminId {
			WriteJSON(w, http.StatusUnauthorized, entity.NewErrorResponse("error unauthorized user id"))
		}

		WriteJSON(w, http.StatusOK, productByID)
		return
	}

	if foundProductSearch {
		ProductBySearch, err := p.productService.SearchProducts(r.Context(), productSearch[0])
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		var productsFiltered []entity.Product
		for _, v := range ProductBySearch {
			if v.UserID == adminId {
				productsFiltered = append(productsFiltered, v)
			}
		}

		WriteJSON(w, http.StatusOK, productsFiltered)
		return
	}

	list, err := p.productService.GetProducts(r.Context(), adminId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (p *ProductAPI) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid product request"))
		return
	}

	if product.Name == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid name request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.productService.AddProduct(r.Context(), entity.Product{
		UserID: adminIdUint,
		Code:   product.Code,
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
		Modal:  product.Modal,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":    adminIdUint,
		"product_id": prod.ID,
		"message":    "success create new product",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (p *ProductAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	productID := r.URL.Query().Get("product_id")
	prodID, _ := strconv.Atoi(productID)
	err := p.productService.DeleteProduct(r.Context(), prodID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":    adminIdUint,
		"product_id": prodID,
		"message":    "success delete product",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (p *ProductAPI) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid product request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("product_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid product id"))
		return
	}

	products, err := p.productService.UpdateProduct(r.Context(), entity.Product{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		Code:   product.Code,
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
		UserID: adminIdUint,
		Modal:  product.Modal,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id":    adminIdUint,
		"product_id": products.ID,
		"message":    "success update product",
	}

	WriteJSON(w, http.StatusOK, response)
}
