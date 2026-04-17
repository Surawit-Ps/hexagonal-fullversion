package sql

import ("hexagonal2/core/entity"
"gorm.io/gorm"
"github.com/google/uuid"
e "hexagonal2/pkg/errors"
"time")

type subscriptionRepositoryDB struct{
	db *gorm.DB
}

type SubscriptionModel struct{
	ID string `gorm:"primaryKey"`
	UserID string
	SubscriptionID string
	Status string
	ExpiryDate time.Time
}

func EnToGormSub(s entity.Subscription)SubscriptionModel{
	return SubscriptionModel{
		ID: s.ID,
		UserID: s.UserID,
		SubscriptionID: s.SubscriptionID,
		Status: s.Status,
		ExpiryDate: s.ExpiryDate,
	}
}

func GormToEnSub(s SubscriptionModel)entity.Subscription{
	return entity.Subscription{
		ID: s.ID,
		UserID: s.UserID,
		SubscriptionID: s.SubscriptionID,
		Status: s.Status,
		ExpiryDate: s.ExpiryDate,
	}
}


func NewSubscriptionRepositoryDB(db *gorm.DB) *subscriptionRepositoryDB{
	return &subscriptionRepositoryDB{db:db}
}

func (r subscriptionRepositoryDB)CreateSubscription(subscription entity.Subscription) error{
	subscription.ID = uuid.New().String()
	subscription.ExpiryDate = time.Now().AddDate(0, 1, 0) // Set expiry date to 1 month from now
	subDB := EnToGormSub(subscription)
	result := r.db.Create(&subDB)
	if result.Error != nil {
		return e.ErrInternalServer
	}
	return nil
}

func (r subscriptionRepositoryDB)GetSubscriptionByUserID(userID string) (*entity.Subscription, error){
	var subscription SubscriptionModel
	result := r.db.Where("user_id = ?", userID).First(&subscription)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrInternalServer
	}
	subscriptionEntity := GormToEnSub(subscription)
	return &subscriptionEntity, nil
}