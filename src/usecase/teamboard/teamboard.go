package teamboard

import (
	"log/slog"
)

// type TeamBoard struct {
// 	db     repository.DBInterface
// 	logger *slog.Logger
// }

type TeamBoardAggregation struct {
	log                  *slog.Logger
	/*

	Думаю не стоит держать бд, будем использовать интерфейсы, например:
	subscriptionSaver    SubscriptionSaver
	subscriptionModifier SubscriptionModifier
	subscriptionRemover  SubscriptionRemover
	subscriptionProvider SubscriptionProvider
	
	*/
}

/*

Паттерн интерфейсов по месту использования:

type SubscriptionSaver interface {
	SaveSubscription(
		ctx context.Context,
		serviceName string,
		price int,
		userID string,
		startDate time.Time,
		stopDate *time.Time,
	) (subID int64, err error)
}

func (a *SubscriptionAggregation) Create(
	ctx context.Context, 
	serviceName string, 
	price int, 
	userID string, 
	startDate time.Time, 
	stopDate *time.Time,
) (int64, error) {
	const op = "subsaggregation.Create"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
		slog.String("serviceName", serviceName),
	)

	id, err := a.subscriptionSaver.SaveSubscription(ctx, serviceName, price, userID, startDate, stopDate)
	if err != nil {
		log.Warn("failed to save subscription", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

*/


func New(log *slog.Logger) *TeamBoardAggregation {
	return &TeamBoardAggregation{
		log: log,
	}
}
