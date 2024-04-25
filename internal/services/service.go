package services

type Service interface {
	FeedService
	UserService
	AuthService
	NotificationService
}

func NewService(addr string) Service {
	return &service{
		FeedService:         NewFeedService(addr),
		UserService:         NewUserService(addr),
		AuthService:         NewAuthService(addr),
		NotificationService: NewNotificationService(addr),
	}
}

type service struct {
	FeedService
	UserService
	AuthService
	NotificationService
}
