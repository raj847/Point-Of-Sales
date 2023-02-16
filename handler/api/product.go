package api

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"vandesar/entity"
	"vandesar/service"
)

type ProductAPI struct {
	productService    *service.ProductService
	cashierRepository service.CashierRepository
}

func NewProductAPI(
	productService *service.ProductService,
	cashierRepository service.CashierRepository,
) *ProductAPI {
	return &ProductAPI{
		productService:    productService,
		cashierRepository: cashierRepository,
	}
}

func (p *ProductAPI) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)

	product := r.URL.Query()

	productID, foundID := product["product_id"]
	productSearch, foundObject := product["search"]

	if foundID == true {
		pID, _ := strconv.Atoi(productID[0])
		productByID, err := p.productService.GetProductByID(r.Context(), pID)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		if productByID.UserID != adminId {
			w.WriteHeader(401)
			_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error unauthorized user id"))
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(productByID)
		return
	} else if foundObject == true {
		ProductBySearch, err := p.productService.GetProductBySearch(r.Context(), productSearch[0])
		if err != nil {
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		var productsFiltered []entity.Product
		for _, v := range ProductBySearch {
			if v.UserID == adminId {
				productsFiltered = append(productsFiltered, v)
			}
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(productsFiltered)
		return
	}

	list, err := p.productService.GetProducts(r.Context(), adminId)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(list)
}

func (p *ProductAPI) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid product request"))
		return
	}

	if product.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid name request"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.productService.AddProduct(r.Context(), &entity.Product{
		UserID: adminIdUint,
		Code:   product.Code,
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
		Modal:  product.Modal,
	})
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    adminIdUint,
		"product_id": prod.ID,
		"message":    "success create new product",
	})
}

func (p *ProductAPI) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
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
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    adminIdUint,
		"product_id": prodID,
		"message":    "success delete product",
	})
}

func (p *ProductAPI) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	adminIdUint := r.Context().Value("id").(uint)
	if adminIdUint == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	id := r.URL.Query().Get("product_id")
	idInt, err := strconv.Atoi(id)

	products, err := p.productService.UpdateProduct(r.Context(), &entity.Product{
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
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":    adminIdUint,
		"product_id": products.ID,
		"message":    "success update product",
	})
}
