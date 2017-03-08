package salt

import (
	"github.com/unchartedsoftware/veldt"
)

// Meta contains information about how to request metadata from a salt server
type Meta struct {
	rmqConfig *Configuration // The configuration defining how we connect to the RabbitMQ server
}

// NewSaltMeta instantiates and returns a pointer to a new generator.
func NewSaltMeta(rmqConfig *Configuration) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		return &Meta{rmqConfig}, nil
	}
}

// Create creates a metadata request
func (meta *Meta) Create (uri string) ([]byte, error) {
	connection, err := NewConnection(meta.rmqConfig)
	if err != nil {
		return nil, err
	}

	return connection.QueryMetadata([]byte(uri))
}

// Parse gets the arguments a metadata constructor will need to create
// metadata requests.  Currently, there is no such information needed,
// so this is a no-op.
func (meta *Meta) Parse (params map[string]interface{}) error {
	return nil
}
