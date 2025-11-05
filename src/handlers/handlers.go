package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	v1 = "/v1"
)

type serverApi struct {
	log            *slog.Logger
	router         *gin.RouterGroup
	tba            TeamBoardAggregation
	createColumnUC CreateColumnUseCase
	// ua     UserAggregation
}

type TeamBoardAggregation interface {
}

/*

Пример интерфейса по SubsAggregation :

type SubsAggregation interface {
	Create(
		ctx context.Context,
		serviceName string,
		price int,
		userID string,
		startDate time.Time,
		stopDate *time.Time,
	) (subID int64, err error)

	GetByID(
		ctx context.Context,
		subID int64,
	) (sub models.Subscription, err error)

	GetAll(
		ctx context.Context,
		subName string,
		userID string,
		offset,
		limit int,
	) (subs []models.Subscription, err error)

	Modify(
		ctx context.Context,
		subID int64,
		patch dto.SubscriptionPatch,
	) (sub models.Subscription, err error)

	GetTotalCost(
		ctx context.Context,
		start time.Time,
		end time.Time,
		subName string,
		userID string,
	) (totalCost int64, err error)

	Delete(
		ctx context.Context,
		subIDs []int64,
	) (err error)
}

*/

func RegisterHandlers(log *slog.Logger, router *gin.RouterGroup, tba TeamBoardAggregation) {

	s := &serverApi{
		log:    log,
		router: router,
		tba:    tba,
		// ua:     ua,
	}
	router.GET("/healthcheck", s.Healthcheck)

	// Это будет регистрация сваггера
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	/*
		path := router.Group(v1)


		Тут будет регистрация всех хендлеров, например:

		path.POST("/subscriptions", s.Create)
		path.GET("/subscriptions/:id", s.GetByID)
		path.GET("/subscriptions", s.GetAll_ByUIDAndSubName)
		path.PATCH("/subscriptions/:id", s.Modify)
		path.DELETE("/subscriptions/:id", s.Delete)
		path.GET("/subscriptions/aggregate", s.GetTotalCost)
	*/
}

func (s *serverApi) Healthcheck(c *gin.Context) {
	// const op = "handlers.Healthcheck"
	// log := s.log.With("op", op, "method", c.Request.Method)
	// log.Info(c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
