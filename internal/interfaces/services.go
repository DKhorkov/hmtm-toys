package interfaces

//go:generate mockgen -source=services.go -destination=../../mocks/services/toys_service.go -exclude_interfaces=MastersService,CategoriesService,TagsService,SsoService -package=mockservices
type ToysService interface {
	ToysRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/tags_service.go -exclude_interfaces=MastersService,CategoriesService,ToysService,SsoService -package=mockservices
type TagsService interface {
	TagsRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/masters_service.go -exclude_interfaces=TagsService,CategoriesService,ToysService,SsoService -package=mockservices
type MastersService interface {
	MastersRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/categories_service.go -exclude_interfaces=TagsService,MastersService,ToysService,SsoService -package=mockservices
type CategoriesService interface {
	CategoriesRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/sso_service.go -exclude_interfaces=TagsService,MastersService,ToysService,CategoriesService -package=mockservices
type SsoService interface {
	SsoRepository
}
