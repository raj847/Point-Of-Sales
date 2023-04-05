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
	RekapAPIHandler       *api.RekapAPI
	ExpensesAPIHandler    *api.ExpensesAPI
}

func main() {
	err := os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	minioClientConn, err := service.NewMinioClient()
	if err != nil {
		log.Fatalf("cannot connect to minio: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	transactRepo := repository.NewTransactionRepository(db)
	rekapRepo := repository.NewRekapRepository(db)
	exRepo := repository.NewExpensesRepository(db)

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	transactService := service.NewTransactionService(transactRepo)
	expensesService := service.NewExpensesService(exRepo)

	rekapService := service.NewRekapService(rekapRepo, transactRepo, userRepo, minioClientConn)

	// run cron job for generate pdf
	service.DoRekapEachMonth(rekapService)
	service.DoRekapEveryDay(rekapService)

	userAPIHandler := api.NewUserAPI(userService, minioClientConn)
	productAPIHandler := api.NewProductAPI(productService, userService)
	transactionAPIHandler := api.NewTransactionAPI(transactService)

	rekapApiHandler := api.NewRekapAPI(rekapRepo)
	expensesApiHandler := api.NewExpensesAPI(expensesService, userService)

	apiHandler := APIHandler{
		UserAPIHandler:        userAPIHandler,
		ProductAPIHandler:     productAPIHandler,
		TransactionAPIHandler: transactionAPIHandler,
		RekapAPIHandler:       rekapApiHandler,
		ExpensesAPIHandler:    expensesApiHandler,
	}

	MuxRoute(mux, "POST", "/api/v1/users/admin/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminRegister)))
	MuxRoute(mux, "POST", "/api/v1/users/admin/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminLogin)))
	MuxRoute(mux, "PUT", "/api/v1/users/kasir/badalakingkong", middleware.Put(middleware.MustCashier(http.HandlerFunc(apiHandler.UserAPIHandler.Badalakingkong))))
	MuxRoute(mux, "GET", "/api/v1/cashiers",
		middleware.Get(
			middleware.MustAdmin(
				http.HandlerFunc(apiHandler.UserAPIHandler.GetAllCashiers),
			),
		),
	)

	// MuxRoute(mux, "PATCH", "/api/v1/cashier/ongleng",
	// 	middleware.Patch(
	// 		middleware.Auth(
	// 			middleware.MustCashier(
	// 				http.HandlerFunc(apiHandler.UserAPIHandler.UpdateOnline)))))

	MuxRoute(mux, "DELETE", "/api/v1/cashier/delete",
		middleware.Delete(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.UserAPIHandler.DeleteCashier)))),
		"?cashier_id=",
	)

	MuxRoute(mux, "PUT", "/api/v1/users/admin/change-password",
		middleware.Put(
			middleware.MustAdmin(
				http.HandlerFunc(apiHandler.UserAPIHandler.ChangeAdminPassword))))

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

	MuxRoute(mux, "POST", "/api/v1/products/createmany",
		middleware.Post(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ProductAPIHandler.CreateNewProductAkeh)))))

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
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.GetAllTransactionsByAdmin),
				),
			),
		),
	)

	MuxRoute(mux, "GET", "/api/v1/transactions/admin/debt",
		middleware.Get(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.GetAllTransactionsByAdminDebt),
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

	MuxRoute(mux, "PUT", "/api/v1/transactions/updatedebt",
		middleware.Put(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.TransactionAPIHandler.UpdateTransactionDebt),
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

	// rekap

	MuxRoute(mux, "GET", "/api/v1/rekap/months",
		middleware.Get(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.RekapAPIHandler.ListRekapPerMonth)))))

	MuxRoute(mux, "GET", "/api/v1/rekap/days",
		middleware.Get(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.RekapAPIHandler.ListRekapPerDays)))))

	// MuxRoute(mux, "POST", "/api/v1/users/admin/check", middleware.Post(middleware.Auth(middleware.MustAdmin(http.HandlerFunc(apiHandler.UserAPIHandler.CheckTokenAdmin)))))

	// MuxRoute(mux, "POST", "/api/v1/users/cashier/check", middleware.Post(middleware.Auth(middleware.MustCashier(http.HandlerFunc(apiHandler.UserAPIHandler.CheckTokenCashier)))))

	MuxRoute(mux, "POST", "/api/v1/users/checker", middleware.Post(middleware.Checker(http.HandlerFunc(apiHandler.UserAPIHandler.CheckToken))))

	MuxRoute(mux, "DELETE", "/api/v1/bebans/delete",
		middleware.Delete(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.DeleteBeban)))),
		"?beban_id=",
	)

	MuxRoute(mux, "PUT", "/api/v1/bebans/update",
		middleware.Put(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.UpdateBeban)))),
		"?beban_id=",
	)

	MuxRoute(mux, "POST", "/api/v1/bebans/create",
		middleware.Post(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.CreateNewBeban)))))

	MuxRoute(mux, "GET", "/api/v1/bebans/get",
		middleware.Get(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.GetAllBeban)))),
		"?beban_id=",
	)

	MuxRoute(mux, "DELETE", "/api/v1/prives/delete",
		middleware.Delete(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.DeletePrive)))),
		"?prive_id=",
	)

	MuxRoute(mux, "PUT", "/api/v1/prives/update",
		middleware.Put(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.UpdatePrive)))),
		"?prive_id=",
	)

	MuxRoute(mux, "POST", "/api/v1/prives/create",
		middleware.Post(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.CreateNewPrive)))))

	MuxRoute(mux, "GET", "/api/v1/prives/get",
		middleware.Get(
			middleware.Auth(
				middleware.MustAdmin(
					http.HandlerFunc(apiHandler.ExpensesAPIHandler.GetAllPrive)))),
		"?prive_id=",
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
