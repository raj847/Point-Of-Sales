package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"vandesar/entity"
	"vandesar/service"
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

func NewProductAPI(productService service.ProductService) *productAPI {
	return &productAPI{productService}
}

func (p *productAPI) GetProduct(w http.ResponseWriter, r *http.Request) {
	// usidS := fmt.Sprintf("%s", r.Context().Value("id"))
	// usid, err := strconv.Atoi(usidS)

	userIdWIthAdminId := r.Context().Value("id").(string)

	userIdStr := strings.Split(userIdWIthAdminId, "|")[0]  // user id
	adminIdStr := strings.Split(userIdWIthAdminId, "|")[1] // admin id

	// convert to int
	userId, _ := strconv.Atoi(userIdStr)
	adminId, _ := strconv.Atoi(adminIdStr)

	if userId == 0 || adminId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	// productID := r.URL.Query().Get("product_id")
	// productName := r.URL.Query().Get("search")
	product := r.URL.Query()

	productID, foundID := product["product_id"]
	productSearch, foundObject := product["search"]

	if foundID == true {
		pID, _ := strconv.Atoi(productID[0])
		productbyID, err := p.productService.GetProductByID(r.Context(), pID)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(productbyID)
		return
	} else if foundObject == true {
		ProductBySearch, err := p.productService.GetProductBySearch(r.Context(), productSearch[0])
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ProductBySearch)
		// fmt.Println(productSearch)
		return
	}

	list, err := p.productService.GetProducts(r.Context(), userId, adminId)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(list)

	// listProduct, err := p.productService.GetProductBySearch(r.Context(), productName)
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	// 	return
	// }
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	// 	return
	// }
	// if len(productID) == 0 {
	// 	list, err := p.productService.GetProducts(r.Context(), userId)
	// 	if err != nil {
	// 		w.WriteHeader(500)
	// 		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	// 		return
	// 	}
	// 	w.WriteHeader(200)
	// 	json.NewEncoder(w).Encode(list)
	// 	return
	// }

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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	prod, err := p.productService.AddProduct(r.Context(), &entity.Product{
		UserID: userId,
		Code:   product.Code,
		Name:   product.Name,
		Price:  product.Price,
		Stock:  product.Stock,
		Modal:  product.Modal,
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
	id := r.URL.Query().Get("product_id")
	idInt, err := strconv.Atoi(id)

	products, err := p.productService.UpdateProduct(r.Context(), &entity.Product{
		ID:    idInt,
		Name:  product.Name,
		Code:  product.Code,
		Price: product.Price,
		Stock: product.Stock,
		Modal: product.Modal,
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
