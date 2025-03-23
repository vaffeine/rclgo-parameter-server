// Code generated by rclgo-gen. DO NOT EDIT.

package rcl_interfaces_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	primitives "github.com/tiiuae/rclgo/pkg/rclgo/primitives"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <rcl_interfaces/msg/list_parameters_result.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("rcl_interfaces/ListParametersResult", ListParametersResultTypeSupport)
	typemap.RegisterMessage("rcl_interfaces/msg/ListParametersResult", ListParametersResultTypeSupport)
}

type ListParametersResult struct {
	Names []string `yaml:"names"`// The resulting parameters under the given prefixes.
	Prefixes []string `yaml:"prefixes"`// The resulting prefixes under the given prefixes.TODO(wjwwood): link to prefix definition and rules.
}

// NewListParametersResult creates a new ListParametersResult with default values.
func NewListParametersResult() *ListParametersResult {
	self := ListParametersResult{}
	self.SetDefaults()
	return &self
}

func (t *ListParametersResult) Clone() *ListParametersResult {
	c := &ListParametersResult{}
	if t.Names != nil {
		c.Names = make([]string, len(t.Names))
		copy(c.Names, t.Names)
	}
	if t.Prefixes != nil {
		c.Prefixes = make([]string, len(t.Prefixes))
		copy(c.Prefixes, t.Prefixes)
	}
	return c
}

func (t *ListParametersResult) CloneMsg() types.Message {
	return t.Clone()
}

func (t *ListParametersResult) SetDefaults() {
	t.Names = nil
	t.Prefixes = nil
}

func (t *ListParametersResult) GetTypeSupport() types.MessageTypeSupport {
	return ListParametersResultTypeSupport
}

// ListParametersResultPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type ListParametersResultPublisher struct {
	*rclgo.Publisher
}

// NewListParametersResultPublisher creates and returns a new publisher for the
// ListParametersResult
func NewListParametersResultPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*ListParametersResultPublisher, error) {
	pub, err := node.NewPublisher(topic_name, ListParametersResultTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &ListParametersResultPublisher{pub}, nil
}

func (p *ListParametersResultPublisher) Publish(msg *ListParametersResult) error {
	return p.Publisher.Publish(msg)
}

// ListParametersResultSubscription wraps rclgo.Subscription to provide type safe helper
// functions
type ListParametersResultSubscription struct {
	*rclgo.Subscription
}

// ListParametersResultSubscriptionCallback type is used to provide a subscription
// handler function for a ListParametersResultSubscription.
type ListParametersResultSubscriptionCallback func(msg *ListParametersResult, info *rclgo.MessageInfo, err error)

// NewListParametersResultSubscription creates and returns a new subscription for the
// ListParametersResult
func NewListParametersResultSubscription(node *rclgo.Node, topic_name string, opts *rclgo.SubscriptionOptions, subscriptionCallback ListParametersResultSubscriptionCallback) (*ListParametersResultSubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg ListParametersResult
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, ListParametersResultTypeSupport, opts, callback)
	if err != nil {
		return nil, err
	}
	return &ListParametersResultSubscription{sub}, nil
}

func (s *ListParametersResultSubscription) TakeMessage(out *ListParametersResult) (*rclgo.MessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneListParametersResultSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneListParametersResultSlice(dst, src []ListParametersResult) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var ListParametersResultTypeSupport types.MessageTypeSupport = _ListParametersResultTypeSupport{}

type _ListParametersResultTypeSupport struct{}

func (t _ListParametersResultTypeSupport) New() types.Message {
	return NewListParametersResult()
}

func (t _ListParametersResultTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.rcl_interfaces__msg__ListParametersResult
	return (unsafe.Pointer)(C.rcl_interfaces__msg__ListParametersResult__create())
}

func (t _ListParametersResultTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.rcl_interfaces__msg__ListParametersResult__destroy((*C.rcl_interfaces__msg__ListParametersResult)(pointer_to_free))
}

func (t _ListParametersResultTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*ListParametersResult)
	mem := (*C.rcl_interfaces__msg__ListParametersResult)(dst)
	primitives.String__Sequence_to_C((*primitives.CString__Sequence)(unsafe.Pointer(&mem.names)), m.Names)
	primitives.String__Sequence_to_C((*primitives.CString__Sequence)(unsafe.Pointer(&mem.prefixes)), m.Prefixes)
}

func (t _ListParametersResultTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*ListParametersResult)
	mem := (*C.rcl_interfaces__msg__ListParametersResult)(ros2_message_buffer)
	primitives.String__Sequence_to_Go(&m.Names, *(*primitives.CString__Sequence)(unsafe.Pointer(&mem.names)))
	primitives.String__Sequence_to_Go(&m.Prefixes, *(*primitives.CString__Sequence)(unsafe.Pointer(&mem.prefixes)))
}

func (t _ListParametersResultTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__rcl_interfaces__msg__ListParametersResult())
}

type CListParametersResult = C.rcl_interfaces__msg__ListParametersResult
type CListParametersResult__Sequence = C.rcl_interfaces__msg__ListParametersResult__Sequence

func ListParametersResult__Sequence_to_Go(goSlice *[]ListParametersResult, cSlice CListParametersResult__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]ListParametersResult, cSlice.size)
	src := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range src {
		ListParametersResultTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(&src[i]))
	}
}
func ListParametersResult__Sequence_to_C(cSlice *CListParametersResult__Sequence, goSlice []ListParametersResult) {
	if len(goSlice) == 0 {
		cSlice.data = nil
		cSlice.capacity = 0
		cSlice.size = 0
		return
	}
	cSlice.data = (*C.rcl_interfaces__msg__ListParametersResult)(C.malloc(C.sizeof_struct_rcl_interfaces__msg__ListParametersResult * C.size_t(len(goSlice))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity
	dst := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range goSlice {
		ListParametersResultTypeSupport.AsCStruct(unsafe.Pointer(&dst[i]), &goSlice[i])
	}
}
func ListParametersResult__Array_to_Go(goSlice []ListParametersResult, cSlice []CListParametersResult) {
	for i := 0; i < len(cSlice); i++ {
		ListParametersResultTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func ListParametersResult__Array_to_C(cSlice []CListParametersResult, goSlice []ListParametersResult) {
	for i := 0; i < len(goSlice); i++ {
		ListParametersResultTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
