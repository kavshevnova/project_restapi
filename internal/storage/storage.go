package storage

//здесь будет хранить общая информация для всех реализаций сторадж
import "errors"

//определим общие ошибки для стораджа
var(
	ErrURLNotFound = errors.New("URL not found")
	ErrURLExists = errors.New("URL exists")
)