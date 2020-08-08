package albums

type Service interface {
	GetAlbums() ([]Album, error)
}

type Repository interface {
	GetAlbums() ([]Album, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) (Service, error) {
	return &service{r}, nil
}

func (s *service) GetAlbums() ([]Album, error) {
	allAlbums, err := s.r.GetAlbums()
	if err != nil {
		return nil, err
	}
	return allAlbums, nil
}
