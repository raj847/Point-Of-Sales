package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vandesar/entity"
	"vandesar/service"

	"github.com/google/uuid"
)

type ProductAPI interface {
	GetProduct(w http.ResponseWriter, r *http.Request)
	CreateNewProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
}

type productAPI struct {
	productService service.ProductService
}

func NewProductAPI(categoryService service.ProductService) *productAPI {
	return &productAPI{categoryService}
}

func (p *productAPI) GetProduct(w http.ResponseWriter, r *http.Request) {
	// usidS := fmt.Sprintf("%s", r.Context().Value("id"))
	// usid, err := strconv.Atoi(usidS)
	userId := r.Context().Value("id").(int)
	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	productID := r.URL.Query().Get("product_id")
	if len(productID) == 0 {
		list, err := p.productService.GetProducts(r.Context(), userId)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(list)
		return
	}
	pID, _ := strconv.Atoi(productID)
	product, err := p.productService.GetProductByID(r.Context(), pID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(product)
	// TODO: answer here
}

func (p *productAPI) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid product request"))
		return
	}
	if product.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid name request"))
		return
	}
	// usidS := fmt.Sprintf("%s", r.Context().Value("id"))
	// usid, err := strconv.Atoi(usidS)
	// fmt.Println(usid)
	userId := r.Context().Value("id").(int)
	fmt.Println(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.productService.AddProduct(r.Context(), &entity.Product{
		UserID: userId,
		Code:   uuid.NewString(),
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userId,
		"product_id": prod.ID,
		"message":    "success create new product"})

	// TODO: answer here
}

func (p *productAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// usidS := fmt.Sprintf("%s", r.Context().Value("id"))
	// usid, err := strconv.Atoi(usidS)
	userId := r.Context().Value("id").(int)
	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	productID := r.URL.Query().Get("product_id")

	prodID, _ := strconv.Atoi(productID)
	err := p.productService.DeleteProduct(r.Context(), prodID)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userId,
		"product_id": prodID,
		"message":    "success delete product"})
	// TODO: answer here
}

func (p *productAPI) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	userId := r.Context().Value("id").(int)
	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	products, err := p.productService.UpdateProduct(r.Context(), &entity.Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	})

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userId,
		"product_id": products.ID,
		"message":    "success update product"})

	// TODO: answer here
}
