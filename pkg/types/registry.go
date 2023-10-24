package types

import (
	"reflect"
	"time"

	"github.com/pkg/errors"
)

var global = NewTypeRegistry()

func Global() TypeRegistry {
	return global
}

type TypeRegistry interface {
	RegisterType(ResourceTypeDescriptor) error
	NewObject(ResourceType) (Resource, error)
	DescriptorFor(resourceType ResourceType) (ResourceTypeDescriptor, error)
}

func NewTypeRegistry() TypeRegistry {
	return &typeRegistry{
		descriptors: make(map[ResourceType]ResourceTypeDescriptor),
	}
}

type typeRegistry struct {
	descriptors map[ResourceType]ResourceTypeDescriptor
}

func (t *typeRegistry) DescriptorFor(resType ResourceType) (ResourceTypeDescriptor, error) {
	typDesc, ok := t.descriptors[resType]
	if !ok {
		return ResourceTypeDescriptor{}, errors.Errorf("invalid resource type %q", resType)
	}
	return typDesc, nil
}

func (t *typeRegistry) RegisterType(res ResourceTypeDescriptor) error {
	if res.Resource.GetSpec() == nil {
		return errors.New("spec in the object cannot be nil")
	}
	if previous, ok := t.descriptors[res.Name]; ok {
		return errors.Errorf("duplicate registration of ResourceType under name %q: previous=%#v new=%#v", res.Name, previous, reflect.TypeOf(res.Resource).Elem().String())
	}
	t.descriptors[res.Name] = res
	return nil
}

func (t *typeRegistry) NewObject(resType ResourceType) (Resource, error) {
	typDesc, ok := t.descriptors[resType]
	if !ok {
		return nil, errors.Errorf("invalid resource type %q", resType)
	}
	return typDesc.NewObject(), nil
}

type ResourceSpec any

type ResourceMeta interface {
	GetName() string
	GetCreationTime() time.Time
}

type Resource interface {
	GetMeta() ResourceMeta
	SetMeta(ResourceMeta)
	GetSpec() ResourceSpec
	SetSpec(ResourceSpec) error
	Descriptor() ResourceTypeDescriptor
}

type ResourceType string

type ResourceTypeDescriptor struct {
	// Name
	Name ResourceType
	// Resource a created element of this type
	Resource Resource
}

func (d ResourceTypeDescriptor) NewObject() Resource {
	specType := reflect.TypeOf(d.Resource.GetSpec()).Elem()
	newSpec := reflect.New(specType).Interface().(ResourceSpec)

	resType := reflect.TypeOf(d.Resource).Elem()
	resource := reflect.New(resType).Interface().(Resource)

	if err := resource.SetSpec(newSpec); err != nil {
		panic(errors.Wrap(err, "could not set spec on the new resource"))
	}
	return resource
}
