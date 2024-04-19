package services

type Service interface {
	FeedService
	UserService
	AuthService
	NotificationServices
}

func NewService(addr string) Service {
	return &service{
		FeedService:          NewFeedService(addr),
		UserService:          NewUserService(addr),
		AuthService:          NewAuthService(addr),
		NotificationServices: NewNotificationService(addr),
	}
}

type service struct {
	FeedService
	UserService
	AuthService
	NotificationServices
}
