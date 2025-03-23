// Code generated by rclgo-gen. DO NOT EDIT.

package rcl_interfaces_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	primitives "github.com/tiiuae/rclgo/pkg/rclgo/primitives"
	builtin_interfaces_msg "github.com/vaffeine/rclgo-parameter-server/internal/msgs/builtin_interfaces/msg"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <rcl_interfaces/msg/parameter_event.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("rcl_interfaces/ParameterEvent", ParameterEventTypeSupport)
	typemap.RegisterMessage("rcl_interfaces/msg/ParameterEvent", ParameterEventTypeSupport)
}

type ParameterEvent struct {
	Stamp builtin_interfaces_msg.Time `yaml:"stamp"`// The time stamp when this parameter event occurred.
	Node string `yaml:"node"`// Fully qualified ROS path to node.
	NewParameters []Parameter `yaml:"new_parameters"`// New parameters that have been set for this node.
	ChangedParameters []Parameter `yaml:"changed_parameters"`// Parameters that have been changed during this event.
	DeletedParameters []Parameter `yaml:"deleted_parameters"`// Parameters that have been deleted during this event.
}

// NewParameterEvent creates a new ParameterEvent with default values.
func NewParameterEvent() *ParameterEvent {
	self := ParameterEvent{}
	self.SetDefaults()
	return &self
}

func (t *ParameterEvent) Clone() *ParameterEvent {
	c := &ParameterEvent{}
	c.Stamp = *t.Stamp.Clone()
	c.Node = t.Node
	if t.NewParameters != nil {
		c.NewParameters = make([]Parameter, len(t.NewParameters))
		CloneParameterSlice(c.NewParameters, t.NewParameters)
	}
	if t.ChangedParameters != nil {
		c.ChangedParameters = make([]Parameter, len(t.ChangedParameters))
		CloneParameterSlice(c.ChangedParameters, t.ChangedParameters)
	}
	if t.DeletedParameters != nil {
		c.DeletedParameters = make([]Parameter, len(t.DeletedParameters))
		CloneParameterSlice(c.DeletedParameters, t.DeletedParameters)
	}
	return c
}

func (t *ParameterEvent) CloneMsg() types.Message {
	return t.Clone()
}

func (t *ParameterEvent) SetDefaults() {
	t.Stamp.SetDefaults()
	t.Node = ""
	t.NewParameters = nil
	t.ChangedParameters = nil
	t.DeletedParameters = nil
}

func (t *ParameterEvent) GetTypeSupport() types.MessageTypeSupport {
	return ParameterEventTypeSupport
}

// ParameterEventPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type ParameterEventPublisher struct {
	*rclgo.Publisher
}

// NewParameterEventPublisher creates and returns a new publisher for the
// ParameterEvent
func NewParameterEventPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*ParameterEventPublisher, error) {
	pub, err := node.NewPublisher(topic_name, ParameterEventTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &ParameterEventPublisher{pub}, nil
}

func (p *ParameterEventPublisher) Publish(msg *ParameterEvent) error {
	return p.Publisher.Publish(msg)
}

// ParameterEventSubscription wraps rclgo.Subscription to provide type safe helper
// functions
type ParameterEventSubscription struct {
	*rclgo.Subscription
}

// ParameterEventSubscriptionCallback type is used to provide a subscription
// handler function for a ParameterEventSubscription.
type ParameterEventSubscriptionCallback func(msg *ParameterEvent, info *rclgo.MessageInfo, err error)

// NewParameterEventSubscription creates and returns a new subscription for the
// ParameterEvent
func NewParameterEventSubscription(node *rclgo.Node, topic_name string, opts *rclgo.SubscriptionOptions, subscriptionCallback ParameterEventSubscriptionCallback) (*ParameterEventSubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg ParameterEvent
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, ParameterEventTypeSupport, opts, callback)
	if err != nil {
		return nil, err
	}
	return &ParameterEventSubscription{sub}, nil
}

func (s *ParameterEventSubscription) TakeMessage(out *ParameterEvent) (*rclgo.MessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneParameterEventSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneParameterEventSlice(dst, src []ParameterEvent) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var ParameterEventTypeSupport types.MessageTypeSupport = _ParameterEventTypeSupport{}

type _ParameterEventTypeSupport struct{}

func (t _ParameterEventTypeSupport) New() types.Message {
	return NewParameterEvent()
}

func (t _ParameterEventTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.rcl_interfaces__msg__ParameterEvent
	return (unsafe.Pointer)(C.rcl_interfaces__msg__ParameterEvent__create())
}

func (t _ParameterEventTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.rcl_interfaces__msg__ParameterEvent__destroy((*C.rcl_interfaces__msg__ParameterEvent)(pointer_to_free))
}

func (t _ParameterEventTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*ParameterEvent)
	mem := (*C.rcl_interfaces__msg__ParameterEvent)(dst)
	builtin_interfaces_msg.TimeTypeSupport.AsCStruct(unsafe.Pointer(&mem.stamp), &m.Stamp)
	primitives.StringAsCStruct(unsafe.Pointer(&mem.node), m.Node)
	Parameter__Sequence_to_C(&mem.new_parameters, m.NewParameters)
	Parameter__Sequence_to_C(&mem.changed_parameters, m.ChangedParameters)
	Parameter__Sequence_to_C(&mem.deleted_parameters, m.DeletedParameters)
}

func (t _ParameterEventTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*ParameterEvent)
	mem := (*C.rcl_interfaces__msg__ParameterEvent)(ros2_message_buffer)
	builtin_interfaces_msg.TimeTypeSupport.AsGoStruct(&m.Stamp, unsafe.Pointer(&mem.stamp))
	primitives.StringAsGoStruct(&m.Node, unsafe.Pointer(&mem.node))
	Parameter__Sequence_to_Go(&m.NewParameters, mem.new_parameters)
	Parameter__Sequence_to_Go(&m.ChangedParameters, mem.changed_parameters)
	Parameter__Sequence_to_Go(&m.DeletedParameters, mem.deleted_parameters)
}

func (t _ParameterEventTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__rcl_interfaces__msg__ParameterEvent())
}

type CParameterEvent = C.rcl_interfaces__msg__ParameterEvent
type CParameterEvent__Sequence = C.rcl_interfaces__msg__ParameterEvent__Sequence

func ParameterEvent__Sequence_to_Go(goSlice *[]ParameterEvent, cSlice CParameterEvent__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]ParameterEvent, cSlice.size)
	src := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range src {
		ParameterEventTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(&src[i]))
	}
}
func ParameterEvent__Sequence_to_C(cSlice *CParameterEvent__Sequence, goSlice []ParameterEvent) {
	if len(goSlice) == 0 {
		cSlice.data = nil
		cSlice.capacity = 0
		cSlice.size = 0
		return
	}
	cSlice.data = (*C.rcl_interfaces__msg__ParameterEvent)(C.malloc(C.sizeof_struct_rcl_interfaces__msg__ParameterEvent * C.size_t(len(goSlice))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity
	dst := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range goSlice {
		ParameterEventTypeSupport.AsCStruct(unsafe.Pointer(&dst[i]), &goSlice[i])
	}
}
func ParameterEvent__Array_to_Go(goSlice []ParameterEvent, cSlice []CParameterEvent) {
	for i := 0; i < len(cSlice); i++ {
		ParameterEventTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func ParameterEvent__Array_to_C(cSlice []CParameterEvent, goSlice []ParameterEvent) {
	for i := 0; i < len(goSlice); i++ {
		ParameterEventTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
