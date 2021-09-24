package application

import (
	"log"
)

func (s Service) Query(search string) error {
	data, err := s.supplier.Supply(search)
	if err != nil {
		return err
	}
	for _, v := range data {
		log.Println(v.Host, v.Url)
	}
	return nil
}
