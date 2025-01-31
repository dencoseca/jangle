package stores

type Store interface {
	Get(name string) (string, error)
	Set(name, value string) error
	Update(name, value string) error
	Delete(name string) error
	List() ([]string, error)
}
