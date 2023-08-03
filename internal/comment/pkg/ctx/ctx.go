package ctx

// ServiceContext contains the components required by the service.
type ServiceContext struct {
}

// NewServiceContext initialize the components and returns a new ServiceContext instance.
func NewServiceContext() *ServiceContext {
	return &ServiceContext{}
}
