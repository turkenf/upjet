/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

// ConfigureWithNameFn functions configure the base resource fields by using
// given name value.
type ConfigureWithNameFn func(base map[string]interface{}, name string)

// NopConfigureWithName does nothing. It's useful for cases where the external
// name is calculated by provider and doesn't have any effect on spec fields.
func NopConfigureWithName(_ map[string]interface{}, _ string) {}

// ConfigurationOption allows setting optional fields of a Configuration object.
type ConfigurationOption func(*Configuration)

// WithExternalName allows you to set an ExternalNamer for given Configuration.
func WithExternalName(e ExternalNamer) ConfigurationOption {
	return func(c *Configuration) {
		c.ExternalNamer = e
	}
}

// WithTerraformIDFieldName allows you to set TerraformIDFieldName.
func WithTerraformIDFieldName(n string) ConfigurationOption {
	return func(c *Configuration) {
		c.TerraformIDFieldName = n
	}
}

// NewConfiguration returns a new *Configuration.
func NewConfiguration(version, kind, terraformResourceType string, opts ...ConfigurationOption) *Configuration {
	c := &Configuration{
		Version:               version,
		Kind:                  kind,
		TerraformResourceType: terraformResourceType,
		TerraformIDFieldName:  "id",
	}
	for _, f := range opts {
		f(c)
	}
	return c
}

// ExternalNamer contains all information that is necessary for naming operations,
// such as removal of those fields from spec schema and calling Configure function
// to fill attributes with information given in external name.
type ExternalNamer struct {
	// SelfVarPath is the Go path to the variable that an instance of this struct
	// is assigned to. It's necessary since there is no way to know a package
	// path of a given variable in runtime.
	SelfVarPath string

	// Configure name attributes of the given configuration using external name.
	Configure ConfigureWithNameFn

	// OmittedFields are the ones you'd like to be removed from the schema since
	// they are specified via external name. You can omit only the top level fields.
	// No field is omitted by default.
	OmittedFields []string

	// DisableNameInitializer allows you to specify whether the name initializer
	// that sets external name to metadata.name if none specified should be disabled.
	// It needs to be disabled for resources whose external name includes information
	// more than the actual name of the resource, like subscription ID or region
	// etc. which is unlikely to be included in metadata.name
	DisableNameInitializer bool
}

// Configuration is the set of information that you can override at different steps
// of the code generation pipeline.
type Configuration struct {
	// Version is the version CRD will have.
	Version string

	// Kind is the kind of the CRD.
	Kind string

	// TerraformResourceType is the name of the resource type in Terraform,
	// like aws_rds_cluster.
	TerraformResourceType string

	// ExternalNamer allows you to specify a custom ExternalNamer.
	ExternalNamer ExternalNamer

	// TerraformIDFieldName is the name of the ID field in Terraform state of
	// the resource. Its default is "id" and in almost all cases, you don't need
	// to overwrite it.
	TerraformIDFieldName string
}
