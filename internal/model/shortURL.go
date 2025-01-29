package model

// Сущность базы данных
type ShortURL struct {
	ID        int
	URL       string
	ShortURL  string
	Username  string
	IsDeleted bool
}

// Создает ShortURL из входящего массива Aliases
func NewShortURLs(aliases []string) []ShortURL {
	ret := []ShortURL{}
	for i, v := range aliases {
		ret = append(ret, ShortURL{})
		ret[i].ShortURL = v
	}
	return ret
}
