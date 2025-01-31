package stores

// Store defines an interface for managing key-value pairs with operations for
// retrieval, modification, deletion, and listing.
type Store interface {
	Get(name string) (string, error)
	Set(name, value string) error
	Update(name, value string) error
	Delete(name string) error
	List() ([]string, error)
}
