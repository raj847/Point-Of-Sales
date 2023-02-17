package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"vandesar/handler/api"
	"vandesar/middleware"
	"vandesar/repository"
	"vandesar/service"
	"vandesar/utils"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler        *api.UserAPI
	ProductAPIHandler     *api.ProductAPI
	TransactionAPIHandler *api.TransactionAPI
}

func main() {
	err := os.Setenv("DATABASE_URL", "postgres://root:secret@localhost:5432/pos")
	if err != nil {
		log.Fatalf("cannot set env: %v", err)
	}

	mux := http.NewServeMux()

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	db := utils.GetDBConnection()
	mux = RunServer(db, mux)

	handler := cors.Default().Handler(mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	transactRepo := repository.NewTransactionRepository(db)

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	transactService := service.NewTransactionService(transactRepo)

	userAPIHandler := api.NewUserAPI(userService)
	productAPIHandler := api.NewProductAPI(productService,userService)
	transactionAPIHandler := api.NewTransactionAPI(transactService)

	apiHandler := APIHandler{
		UserAPIHandler:        userAPIHandler,
		ProductAPIHandler:     productAPIHandler,
		TransactionAPIHandler: transactionAPIHandler,
	}

	MuxRoute(mux, "POST", "/api/v1/users/admin/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminRegister)))
	MuxRoute(mux, "POST", "/api/v1/users/admin/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminLogin)))

	MuxRoute(mux, "POST", "/api/v1/users/cashier/register",
		middleware.Post(
			middleware.MustAdmin(
				http.HandlerFunc(apiHandler.UserAPIHandler.CashierRegister))))

	MuxRoute(mux, "POST", "/api/v1/users/cashier/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.CashierLogin)))

	MuxRoute(mux, "POST", "/api/v1/products/create",
		middleware.Post(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ProductAPIHandler.CreateNewProduct)))))

	MuxRoute(mux, "GET", "/api/v1/products",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.ProductAPIHandler.GetAllProducts),
			),
		),
		"?product_id=", "&search=",
	)

	MuxRoute(mux, "PUT", "/api/v1/products/update",
		middleware.Put(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ProductAPIHandler.UpdateProduct)))),
		"?product_id=",
	)

	MuxRoute(mux, "DELETE", "/api/v1/products/delete",
		middleware.Delete(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ProductAPIHandler.DeleteProduct)))),
		"?product_id=",
	)

	MuxRoute(mux, "POST", "/api/v1/transactions/create",
		middleware.Post(
			middleware.Auth(
				middleware.MustCashier(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.CreateTransaction)))))

	MuxRoute(mux, "GET", "/api/v1/transactions/admin",
		middleware.Get(
			middleware.Auth(
				middleware.MustCashier(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.GetAllTransactionsByAdmin),
				),
			),
		),
	)

	MuxRoute(mux, "GET", "/api/v1/transactions/cashier",
		middleware.Get(
			middleware.Auth(
				middleware.MustCashier(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.GetAllTransactionsByCashier),
				),
			),
		),
	)

	MuxRoute(mux, "PUT", "/api/v1/transactions/update",
		middleware.Put(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.UpdateTransaction),
				),
			),
		),
		"?transaction_id=",
	)

	MuxRoute(mux, "DELETE", "/api/v1/transactions/delete",
		middleware.Delete(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.DeleteTransaction),
				),
			),
		),
		"?transaction_id=",
	)

	return mux
}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
