package config

/*
Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
*/

var Config = struct {
	a string
	b string
}{
	a: "",
	b: "",
}

func initConfig() {

}
