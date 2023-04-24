package repository

import "errors"

type repo struct {
	shortToFull map[string]string
	fullToShort map[string]string
}

func NewInMemoryRepository() *repo {
	return &repo{
		shortToFull: make(map[string]string),
		fullToShort: make(map[string]string),
	}
}

func (r *repo) ExistsShortURL(fullURL string) *string {
	if shortURL, ok := r.fullToShort[fullURL]; ok {
		return &shortURL
	}
	return nil
}

func (r *repo) SaveURL(shortURL, fullURL string) error {
	if _, ok := r.shortToFull[shortURL]; ok {
		return errors.New("already exists")
	}
	r.shortToFull[shortURL] = fullURL
	r.fullToShort[fullURL] = shortURL
	return nil
}

func (r *repo) GetURL(shortURL string) (string, error) {
	var url string
	var ok bool
	if url, ok = r.shortToFull[shortURL]; !ok {
		return "", errors.New("no full url")
	}
	return url, nil
}
