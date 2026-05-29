package app

import (
	"fmt"
)

// Resources type
type Resources struct {
	// CPU is the number of CPU cores
	CPU string `json:"cpu,omitempty"`
	// Memory is the amount of memory
	Memory string `json:"memory,omitempty"`
}

// custom resource type
type CustomResource struct {
	// Name of the custom resource
	Name string `json:"name,omitempty"`
	// Value of the custom resource
	Value string `json:"value,omitempty"`
}

// Limit represents a resource limit
type Limit struct {
	// Max defines the resource limits
	Max Resources `json:"max,omitempty"`
	// Min defines the resource request
	Min Resources `json:"min,omitempty"`
	// Default stes resource limits to pods without defined resource limits
	Default Resources `json:"default,omitempty"`
	// DefaultRequest sets the resource requests for pods without defined resource requests
	DefaultRequest Resources `json:"defaultRequest,omitempty"`
	// MaxLimitRequestRatio set the max limit request ratio
	MaxLimitRequestRatio Resources `json:"maxLimitRequestRatio,omitempty"`
	Type                 string    `json:"type"`
}

// Limits type
type Limits []Limit

// quota type
type Quotas struct {
	// Pods is the pods quota
	Pods string `json:"pods,omitempty"`
	// CPULimits is the CPU quota
	CPULimits string `json:"limits.cpu,omitempty"`
	// CPURequests is the CPU requests quota
	CPURequests string `json:"requests.cpu,omitempty"`
	// MemoryLimits is the memory quota
	MemoryLimits string `json:"limits.memory,omitempty"`
	// MemoryRequests is the memory requests quota
	MemoryRequests string `json:"requests.memory,omitempty"`
	// CustomResource is a list of custom resource quotas
	CustomQuotas []CustomResource `json:"customQuotas,omitempty"`
}

// Namespace type represents the fields of a Namespace
type Namespace struct {
	// Protected if set to true no changes can be applied to the namespace
	Protected bool `json:"protected"`
	// Limits to set on the namespace
	Limits Limits `json:"limits,omitempty"`
	// Labels to set to the namespace
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations to set on the namespace
	Annotations map[string]string `json:"annotations,omitempty"`
	// Quotas to set on the namespace
	Quotas   *Quotas `json:"quotas,omitempty"`
	disabled bool
}

func (n *Namespace) Disable() {
	n.disabled = true
}

// print prints the namespace
func (n *Namespace) print() {
	fmt.Println("\tprotected: ", n.Protected)
	fmt.Println("\tdisabled: ", n.disabled)
	fmt.Println("\tlabels:")
	printMap(n.Labels, 2)
	fmt.Println("\tannotations:")
	printMap(n.Annotations, 2)
	fmt.Println("-------------------")
}
