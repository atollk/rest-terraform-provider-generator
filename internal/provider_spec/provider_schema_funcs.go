package provider_spec

import (
	"fmt"
	"log"
)

type RESTOperation struct {
	name string
}

var (
	Create = RESTOperation{"Create"}
	Read   = RESTOperation{"Read"}
	Update = RESTOperation{"Update"}
	Delete = RESTOperation{"Delete"}
)

func (r *ResourceSchema) GetOperationPath(operation RESTOperation, defaults *GlobalDefaults) string {
	var path string
	switch operation {
	case Create:
		if r.Create != nil {
			path = r.Create.Path
		}
	case Read:
		if r.Read != nil {
			path = r.Read.Path
		}
	case Update:
		if r.Update != nil {
			path = r.Update.Path
		}
	case Delete:
		if r.Destroy != nil {
			path = r.Destroy.Path
		}
	default:
		log.Panicf(fmt.Sprintf("%s is not a valid REST operation", operation.name))
	}
	if path == "" {
		if operation == Create {
			path = r.Path
		} else {
			idAttribute := r.IdAttribute
			if idAttribute == "" {
				idAttribute = defaults.IdAttribute
			}
			path = fmt.Sprintf("%s/{%s}", r.Path, idAttribute)
		}
	}
	return path
}

func (r *ResourceSchema) GetOperationMethod(operation RESTOperation, defaults *GlobalDefaults) string {
	var method string
	switch operation {
	case Create:
		if r.Create != nil {
			method = r.Create.Method
		}
	case Read:
		if r.Read != nil {
			method = r.Read.Method
		}
	case Update:
		if r.Update != nil {
			method = r.Update.Method
		}
	case Delete:
		if r.Destroy != nil {
			method = r.Destroy.Method
		}
	default:
		log.Panicf(fmt.Sprintf("%s is not a valid REST operation", operation.name))
	}
	if method == "" {
		switch operation {
		case Create:
			method = defaults.CreateMethod
		case Read:
			method = defaults.ReadMethod
		case Update:
			method = defaults.UpdateMethod
		case Delete:
			method = defaults.DestroyMethod
		default:
			log.Panicf(fmt.Sprintf("%s is not a valid REST operation", operation.name))
		}
	}
	return method
}
