package interfaces

type ToysService interface {
	ToysRepository
}

type TagsService interface {
	TagsRepository
}

type MastersService interface {
	MastersRepository
}

type CategoriesService interface {
	CategoriesRepository
}
