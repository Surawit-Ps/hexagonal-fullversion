package unused

// import ("fmt"
// "hexagonal2/core/ports"
// "hexagonal2/core/entity")

// type subService struct{
// 	repo ports.SubscriptionRepository
// }

// func NewSubService(repo ports.SubscriptionRepository) subService{
// 	return subService{repo: repo}
// }

// func(r subService)CreateSubscription(subscription entity.Subscription) error{
// 	err := r.repo.CreateSubscription(subscription)
// 	if err != nil{
// 		fmt.Print(err)
// 		return err
// 	}
// 	return nil
// }

// func(r subService)GetSubscriptionByUserID(userID string)(*entity.Subscription,error){
// 	sub,err := r.repo.GetSubscriptionByUserID(userID)
// 	if err != nil{
// 		fmt.Print(err)
// 		return nil,err
// 	}
// 	return sub,nil
// }