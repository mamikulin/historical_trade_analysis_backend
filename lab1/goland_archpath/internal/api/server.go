package api

import (
	"archpath/internal/app/handler"
	"archpath/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {

	log.Println("Server start up")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("../../templates/*")
	r.Static("./static", "../../resources")

	r.GET("/", handler.GetCommodities)
	r.GET("/commodity/:id", handler.GetCommodity)
	r.GET("/analysis", handler.GetAnalysisPage)

	r.Run()
	log.Println("Server down")
}
