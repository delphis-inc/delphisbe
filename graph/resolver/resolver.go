// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import "github.com/nedrocks/delphisbe/internal/backend"

type Resolver struct {
	DAOManager backend.DAOManager
}
