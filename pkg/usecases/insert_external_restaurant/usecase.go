package insertexternalrestaurant

import "time"

type Repository interface {
	InsertExternalCusine(externalUUID, x, y, placeUrl string, updatedAt time.Time, placeName string) error
}

type UseCase struct {
	Repository Repository
}

func NewUseCase(Repository Repository) *UseCase {
	return &UseCase{
		Repository: Repository,
	}
}

func (u *UseCase) InsertExternalRestaurant() error {

	// array로 받고 array로 넘기기

	return u.Repository.InsertExternalCusine()
}
