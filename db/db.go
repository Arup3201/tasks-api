package db

const DBFile = "tasks.json"

type DB interface {
	Write(p []byte)
	Read() []byte
}

func Add(file string) {

}

func Get(file string) {

}
