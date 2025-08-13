package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
) (crud.ListSubscriptionsRow, error) {
	changedDate := false
	changedSub := false
	q := crud.New(db)
	subs, err := q.ListSubscriptions(ctx, crud.ListSubscriptionsParams{
		Column1: userUuid,
		Column2: serviceName,
	})
	if err != nil {
		log.Printf("error occured while executig db query - 'ListSubscriptions', error - '%s'", err.Error())
		return crud.ListSubscriptionsRow{}, err
	}
	currentSub := subs[0]

	if startDate != nil || stopDate != nil {
		var (
			newStart time.Time
			newStop  time.Time
		)

		if startDate != nil {
			dtSlc, err := transformDate(*startDate)
			if err != nil {
				return crud.ListSubscriptionsRow{}, err
			}
			newStart = dtSlc[0]
		}

		if stopDate != nil {
			dtSlc, err := transformDate(*stopDate)
			if err != nil {
				return crud.ListSubscriptionsRow{}, err
			}
			newStop = dtSlc[0]
		}

		if !newStart.IsZero() && !newStop.IsZero() {
			if newStart.After(newStop) {
				log.Printf("start date must be less than the stop date")
				return crud.ListSubscriptionsRow{}, errors.New("start date more than stop date")
			}
			if newStart == currentSub.StartDate && newStop == currentSub.StopDate {

			}
		} else if !newStart.IsZero() {
			if newStart.After(currentSub.StopDate) {
				log.Printf("new date start - '%v' more then old date stop - '%v' for service - '%s' user - '%s'", newStart, currentSub.StopDate, serviceName, userUuid)
				return crud.ListSubscriptionsRow{}, fmt.Errorf("new date start - '%v' more then old date stop - '%v'", newStart, currentSub.StopDate)
			}
			newStop = currentSub.StopDate
		} else if !newStop.IsZero() {
			if currentSub.StartDate.After(newStop) {
				log.Printf("new date stop - '%v' less then old date start - '%v' for service - '%s' user - '%s'", newStop, currentSub.StartDate, serviceName, userUuid)
				return crud.ListSubscriptionsRow{}, fmt.Errorf("new date stop - '%v' less then old date start - '%v'", newStop, currentSub.StartDate)
			}
			newStart = currentSub.StartDate
		}

		quantity, err := q.ChangeSubscriptionData(ctx, crud.ChangeSubscriptionDataParams{
			Column1: newStart,
			Column2: newStop,
			Column3: userUuid,
			Column4: serviceName,
		})
		if err != nil {
			log.Printf("error occured while executig db query - 'ChangeSubscriptionData', error - '%s'", err.Error())
			return crud.ListSubscriptionsRow{}, err
		}

		if quantity == 0 {
			log.Printf("service - '%s' for user - '%s' does not exists", serviceName, userUuid)
			return crud.ListSubscriptionsRow{}, err
		}
		changedSub = true
		changedDate = true
	}

	if price != nil {
		quantity, err := q.UpdateService(ctx, crud.UpdateServiceParams{
			Column1: int32(*price),
			Column2: serviceName,
		})
		if err != nil {
			log.Printf("error occured while executig db query - 'UpdateService', error - '%s'", err.Error())
			if changedDate {
				q.ChangeSubscriptionData(ctx, crud.ChangeSubscriptionDataParams{ // transaction-rollback
					Column1: currentSub.StartDate,
					Column2: currentSub.StopDate,
					Column3: userUuid,
					Column4: serviceName,
				})
			}
			return crud.ListSubscriptionsRow{}, err
		}

		if quantity == 0 {
			log.Printf("srvice - '%s' not found", serviceName)
			if changedDate {
				q.ChangeSubscriptionData(ctx, crud.ChangeSubscriptionDataParams{ // transaction-rollback
					Column1: currentSub.StartDate,
					Column2: currentSub.StopDate,
					Column3: userUuid,
					Column4: serviceName,
				})
			}
			return crud.ListSubscriptionsRow{}, err
		}
		changedSub = true
	}

	if changedSub {
		allData, _ := q.ListSubscriptions(ctx, crud.ListSubscriptionsParams{
			Column1: userUuid,
			Column2: serviceName,
		})
		return allData[0], nil
	}
	log.Printf("for service - '%s' user - '%s' nothing to change", serviceName, userUuid)
	return crud.ListSubscriptionsRow{}, nil
}

func DeleteSubscription(
	ctx context.Context,
	db *sql.DB,
	serviceName, userUuid string,
) error {
	log.Println(`start function 'DeleteSubscription' into controller/usecase`)
	q := crud.New(db)
	row, err := q.DeleteSubscription(ctx, crud.DeleteSubscriptionParams{
		UserUuid: userUuid,
		Name:     serviceName,
	})
	if err != nil {
		log.Printf("error occured while executig db query - 'DeleteSubscription' in method 'DeleteSubscription', error - '%s'", err.Error())
		return err
	}

	if row == 0 {
		log.Printf("subscription on service - '%s' for user - '%s' does not exists", serviceName, userUuid)
		return errors.New("subscription does not exists")
	}
	return nil
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
