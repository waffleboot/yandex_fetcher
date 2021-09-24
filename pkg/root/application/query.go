package application

func (s Service) Query(search string) error {
	data, err := s.supplier.Supply(search)
	if err != nil {
		return err
	}
	for _, v := range data {
		// log.Println(v.Host, v.Url)
		s.benchmark.Benchmark(v.Url)
	}
	return nil
}
