package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *userDeviceResolver) Platform(ctx context.Context, obj *model.UserDevice) (model.Platform, error) {
	switch strings.ToLower(obj.Platform) {
	case "ios":
		return model.PlatformIos, nil
	case "android":
		return model.PlatformAndroid, nil
	case "web":
		return model.PlatformWeb, nil
	default:
		return model.PlatformUnknown, nil
	}
}

// UserDevice returns generated.UserDeviceResolver implementation.
func (r *Resolver) UserDevice() generated.UserDeviceResolver { return &userDeviceResolver{r} }

type userDeviceResolver struct{ *Resolver }
