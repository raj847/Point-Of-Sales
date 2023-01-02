package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"vandesar/handler/api"
	"vandesar/middleware"
	"vandesar/repository"
	"vandesar/service"
	"vandesar/utils"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler api.UserAPI
}

// func FlyURL() string {
// 	return "https://final-web-app.fly.dev" // TODO: replace this
// }

func main() {

	//TODO: hapus jika sudah di deploy di fly.io
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:9090/postgres")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		mux := http.NewServeMux()

		err := utils.ConnectDB()
		if err != nil {
			panic(err)
		}

		db := utils.GetDBConnection()

		mux = RunServer(db, mux)

		fmt.Println("Server is running on port 8080")
		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userAPIHandler := api.NewUserAPI(userService)

	apiHandler := APIHandler{
		UserAPIHandler: userAPIHandler,
	}

	MuxRoute(mux, "POST", "/api/v1/users/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Login)))
	MuxRoute(mux, "POST", "/api/v1/users/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Register)))
	MuxRoute(mux, "POST", "/api/v1/users/logout", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Logout)))
	MuxRoute(mux, "DELETE", "/api/v1/users/delete", middleware.Delete(http.HandlerFunc(apiHandler.UserAPIHandler.Delete)), "?user_id=")

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
