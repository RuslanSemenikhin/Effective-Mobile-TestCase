package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	crud "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/db/querys/gen"
	"github.com/google/uuid"
)

func ListSubscriptions(
	ctx context.Context,
	db *sql.DB,
	serviceName, userUuid *string,
) ([]crud.ListSubscriptionsRow, error) {
	log.Println(`start function 'ListSubscriptions' into controller/usecase`)

	q := crud.New(db)
	resp, err := q.ListSubscriptions(ctx, crud.ListSubscriptionsParams{
		Column1: *userUuid,
		Column2: *serviceName,
	})
	if err != nil {
		log.Printf("error occured while executig db query - 'ListSubscriptions', error - '%s'", err.Error())
		return []crud.ListSubscriptionsRow{}, err
	}

	return resp, nil
}

func TotalPrice(
	ctx context.Context,
	db *sql.DB,
	startDate, stopDate string,
	userUuid, serviceName *string,
) (int64, error) {
	log.Println(`start function 'ListTotalPrice' into controller/usecase`)

	dateSlc, err := transformDate(startDate, stopDate)
	if err != nil {
		log.Printf("error occured while transform date from string-type to time.Time-type, %v", err)
		return 0, err
	}

	start, stop := dateSlc[0], dateSlc[1]
	if start.After(stop) {
		log.Printf("start date must be less than the stop date")
		return 0, errors.New("start date more than stop date")
	}

	q := crud.New(db)
	totalPrice, err := q.TotalPriceSubscriptions(
		ctx,
		crud.TotalPriceSubscriptionsParams{
			Column1: start,
			Column2: stop,
			Column3: *userUuid,
			Column4: *serviceName,
		},
	)
	if err != nil {
		log.Printf("error occured while executig db query - 'TotalPriceSubscriptions', error - '%s'", err.Error())
		return 0, err
	}

	return totalPrice, nil
}

func AddSubscription(
	ctx context.Context,
	db *sql.DB,
	serviceName, userUuid string,
	price int64,
	startDate, stopDate string,
) (crud.ServicesSubscription, error) {
	log.Println(`start function 'AddSubscription' into controller/usecase`)
	q := crud.New(db)
	exists, err := q.ExistingSubscriptionService(ctx, crud.ExistingSubscriptionServiceParams{
		UserUuid: userUuid,
		Name:     serviceName,
	})
	if err != nil {
		log.Printf("error occured while executig db query - 'ExistingSubscriptionService', error - '%s'", err.Error())
		return crud.ServicesSubscription{}, err
	}

	if exists {
		log.Printf("user with uuid - '%s' have subscription on service - '%s'", userUuid, serviceName)
		return crud.ServicesSubscription{}, errors.New("user have subscription")
	}

	dateSlc, err := transformDate(startDate, stopDate)
	if err != nil {
		log.Printf("error occured while transform date from string-type to time.Time-type, %v", err)
		return crud.ServicesSubscription{}, err
	}

	start, stop := dateSlc[0], dateSlc[1]
	if start.After(stop) {
		log.Printf("start date must be less than the stop date")
		return crud.ServicesSubscription{}, errors.New("start date more than stop date")
	}

	serviceUuid, err := q.ExistingService(ctx, serviceName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error occured while executig db query - 'ExistingService', error - '%s'", err.Error())
			return crud.ServicesSubscription{}, err
		}
	}

	flag := false

	if len(serviceUuid) == 0 {
		serviceUuid = uuid.New().String()
		if _, err := q.AddService(ctx, crud.AddServiceParams{
			Uuid:  serviceUuid,
			Name:  serviceName,
			Price: int32(price),
		}); err != nil {
			log.Printf("error occured while executig db query - 'AddService', error - '%s'", err.Error())
			return crud.ServicesSubscription{}, err
		}
		flag = true
	}

	sub, err := q.AddSubscription(ctx, crud.AddSubscriptionParams{
		UserUuid:    userUuid,
		ServiceUuid: serviceUuid,
		StartDate:   start,
		StopDate:    stop,
	})
	if err != nil {
		log.Printf("error occured while executig db query - 'AddSubscription', error - '%s'", err.Error())
		if flag {
			if err := q.DeleteService(ctx, serviceUuid); err != nil {
				log.Printf("error occured while executig db query - 'DeleteService' in method 'AddSubscription', error - '%s'", err.Error())
				return crud.ServicesSubscription{}, err
			}
			return crud.ServicesSubscription{}, err
		}
		return crud.ServicesSubscription{}, err
	}

	log.Println(`finished successfuly function 'AddSubscription' into controller/usecase`)
	return sub, nil
}

func UpdatedSubscription(
	ctx context.Context,
	db *sql.DB,
	serviceName, userUuid string,
	price *int64,
	startDate, stopDate *string,
) string {
	return "'UpdatedSubscription' method into controller"
}

func DeleteSubscription(
	ctx context.Context,
	db *sql.DB,
	serviceName, userUuid string,
) string {
	return "'DeleteSubscription' method into controller"
}

func transformDate(dates ...string) ([]time.Time, error) {
	var (
		errSlc  []error
		dateSlc []time.Time
	)
	for _, date := range dates {
		dt, err := time.Parse("01.2006", date)
		if err != nil {
			errSlc = append(errSlc, err)
			continue
		}
		goodDate := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
		dateSlc = append(dateSlc, goodDate)
	}
	if len(errSlc) > 0 {
		err := errorHandler(errSlc)
		return nil, err
	}

	return dateSlc, nil
}

func errorHandler(errorsSlc []error) error {
	var resErr strings.Builder
	for _, er := range errorsSlc {
		resErr.WriteString(er.Error())
		resErr.WriteString(" | ")
	}
	return errors.New(strings.TrimRight(resErr.String(), " | "))
}
