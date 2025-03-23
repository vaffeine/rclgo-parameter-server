package parameter_server

import (
	"errors"
	"fmt"
	"sync"
	"time"

	builtin_interfaces_msg "github.com/vaffeine/rclgo-parameter-server/internal/msgs/builtin_interfaces/msg"
	rcl_interfaces_msg "github.com/vaffeine/rclgo-parameter-server/internal/msgs/rcl_interfaces/msg"
	rcl_interfaces_srv "github.com/vaffeine/rclgo-parameter-server/internal/msgs/rcl_interfaces/srv"

	"github.com/tiiuae/rclgo/pkg/rclgo"
)

type (
	ParameterDescriptor = rcl_interfaces_msg.ParameterDescriptor
	FloatingPointRange  = rcl_interfaces_msg.FloatingPointRange
	IntegerRange        = rcl_interfaces_msg.IntegerRange
)

const (
	ParameterType_PARAMETER_BOOL          uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_BOOL
	ParameterType_PARAMETER_INTEGER       uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER
	ParameterType_PARAMETER_DOUBLE        uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE
	ParameterType_PARAMETER_STRING        uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_STRING
	ParameterType_PARAMETER_BYTE_ARRAY    uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_BYTE_ARRAY
	ParameterType_PARAMETER_BOOL_ARRAY    uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_BOOL_ARRAY
	ParameterType_PARAMETER_INTEGER_ARRAY uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER_ARRAY
	ParameterType_PARAMETER_DOUBLE_ARRAY  uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE_ARRAY
	ParameterType_PARAMETER_STRING_ARRAY  uint8 = rcl_interfaces_msg.ParameterType_PARAMETER_STRING_ARRAY
)

type ParameterServer struct {
	logger                  *rclgo.Logger
	nodeName                string
	parameterEventPublisher *rcl_interfaces_msg.ParameterEventPublisher
	services                services

	mu         sync.Mutex
	parameters map[string]*Parameter
}

func New(node *rclgo.Node, parameters ...Parameter) (*ParameterServer, error) {
	newParameters := make([]rcl_interfaces_msg.Parameter, 0, len(parameters))
	parametersMap := make(map[string]*Parameter, len(parameters))
	for _, param := range parameters {
		newParameters = append(newParameters, rcl_interfaces_msg.Parameter{
			Name:  param.Descriptor.Name,
			Value: param.asRclValue(),
		})
		parametersMap[param.Descriptor.Name] = &param
	}

	parameterEventPublisher, err := rcl_interfaces_msg.NewParameterEventPublisher(node, "/parameter_events", nil)
	if err != nil {
		return nil, fmt.Errorf("create parameter events publisher: %w", err)
	}
	defer func() {
		if err != nil {
			parameterEventPublisher.Close()
		}
	}()
	err = parameterEventPublisher.Publish(&rcl_interfaces_msg.ParameterEvent{
		Stamp:         toRosTimestamp(time.Now()),
		Node:          node.Name(),
		NewParameters: newParameters,
	})
	if err != nil {
		node.Logger().Warnf("failed to publish parameter event: %v", err)
	}
	server := ParameterServer{
		logger:                  node.Logger(),
		nodeName:                node.Name(),
		parameterEventPublisher: parameterEventPublisher,
		parameters:              parametersMap,
	}

	describeParameters, err := rcl_interfaces_srv.NewDescribeParametersService(
		node,
		fmt.Sprintf("%s/describe_parameters", node.Name()),
		nil,
		server.describeParametersImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create describe parameters: %w", err)
	}
	defer func() {
		if err != nil {
			describeParameters.Close()
		}
	}()

	getParameterTypes, err := rcl_interfaces_srv.NewGetParameterTypesService(
		node,
		fmt.Sprintf("%s/get_parameter_types", node.Name()),
		nil,
		server.getParameterTypesImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create get parameter types: %w", err)
	}
	defer func() {
		if err != nil {
			getParameterTypes.Close()
		}
	}()

	getParameters, err := rcl_interfaces_srv.NewGetParametersService(
		node,
		fmt.Sprintf("%s/get_parameters", node.Name()),
		nil,
		server.getParametersImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create get parameters: %w", err)
	}
	defer func() {
		if err != nil {
			getParameters.Close()
		}
	}()

	listParameters, err := rcl_interfaces_srv.NewListParametersService(
		node,
		fmt.Sprintf("%s/list_parameters", node.Name()),
		nil,
		server.listParametersImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create list parameters: %w", err)
	}
	defer func() {
		if err != nil {
			listParameters.Close()
		}
	}()

	setParameters, err := rcl_interfaces_srv.NewSetParametersService(
		node,
		fmt.Sprintf("%s/set_parameters", node.Name()),
		nil,
		server.setParametersImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create set parameters: %w", err)
	}
	defer func() {
		if err != nil {
			setParameters.Close()
		}
	}()

	setParametersAtomically, err := rcl_interfaces_srv.NewSetParametersAtomicallyService(
		node,
		fmt.Sprintf("%s/set_parameters_atomically", node.Name()),
		nil,
		server.setParametersAtomicallyImpl,
	)
	if err != nil {
		return nil, fmt.Errorf("create set parameters atomically: %w", err)
	}
	defer func() {
		if err != nil {
			setParametersAtomically.Close()
		}
	}()

	server.services = services{
		describeParameters:      describeParameters,
		getParameterTypes:       getParameterTypes,
		getParameters:           getParameters,
		listParameters:          listParameters,
		setParameters:           setParameters,
		setParametersAtomically: setParametersAtomically,
	}
	return &server, nil
}

func (s *ParameterServer) Close() error {
	return errors.Join(
		errContext("parameter event publisher", s.parameterEventPublisher.Close()),
		errContext("describe parameters", s.services.describeParameters.Close()),
		errContext("get paramater types", s.services.getParameterTypes.Close()),
		errContext("get paramaters", s.services.getParameters.Close()),
		errContext("list paramaters", s.services.listParameters.Close()),
		errContext("set paramaters", s.services.setParameters.Close()),
		errContext("set paramaters atomically", s.services.setParametersAtomically.Close()),
	)
}

func (s *ParameterServer) Lock() {
	s.mu.Lock()
}

func (s *ParameterServer) Unlock() {
	s.mu.Unlock()
}

func (s *ParameterServer) Get(name string) any {
	return s.parameters[name].Value
}

func errContext(context string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	} else {
		return nil
	}
}

type services struct {
	describeParameters      *rcl_interfaces_srv.DescribeParametersService
	getParameterTypes       *rcl_interfaces_srv.GetParameterTypesService
	getParameters           *rcl_interfaces_srv.GetParametersService
	listParameters          *rcl_interfaces_srv.ListParametersService
	setParameters           *rcl_interfaces_srv.SetParametersService
	setParametersAtomically *rcl_interfaces_srv.SetParametersAtomicallyService
}

func (s *ParameterServer) describeParametersImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.DescribeParameters_Request, sender rcl_interfaces_srv.DescribeParametersServiceResponseSender) {
	descriptors := make([]rcl_interfaces_msg.ParameterDescriptor, 0, len(req.Names))
	for _, name := range req.Names {
		param, ok := s.parameters[name]
		if !ok {
			sender.SendResponse(&rcl_interfaces_srv.DescribeParameters_Response{})
			return
		}
		descriptors = append(descriptors, param.Descriptor)
	}
	sender.SendResponse(&rcl_interfaces_srv.DescribeParameters_Response{
		Descriptors: descriptors,
	})
}

func (s *ParameterServer) getParameterTypesImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.GetParameterTypes_Request, sender rcl_interfaces_srv.GetParameterTypesServiceResponseSender) {
	types := make([]uint8, 0, len(req.Names))
	for _, name := range req.Names {
		param, ok := s.parameters[name]
		if !ok {
			sender.SendResponse(&rcl_interfaces_srv.GetParameterTypes_Response{})
			return
		}
		types = append(types, param.Descriptor.Type)
	}
	sender.SendResponse(&rcl_interfaces_srv.GetParameterTypes_Response{
		Types: types,
	})
}

func (s *ParameterServer) getParametersImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.GetParameters_Request, sender rcl_interfaces_srv.GetParametersServiceResponseSender) {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := make([]rcl_interfaces_msg.ParameterValue, 0, len(req.Names))
	for _, name := range req.Names {
		param, ok := s.parameters[name]
		if !ok {
			sender.SendResponse(&rcl_interfaces_srv.GetParameters_Response{})
			return
		}
		value := param.asRclValue()
		values = append(values, value)
	}
	sender.SendResponse(&rcl_interfaces_srv.GetParameters_Response{
		Values: values,
	})
}

func (s *ParameterServer) listParametersImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.ListParameters_Request, sender rcl_interfaces_srv.ListParametersServiceResponseSender) {
	names := make([]string, 0, len(s.parameters))
	for _, param := range s.parameters {
		// TODO: handle prefixes
		names = append(names, param.Descriptor.Name)
	}
	sender.SendResponse(&rcl_interfaces_srv.ListParameters_Response{
		Result: rcl_interfaces_msg.ListParametersResult{
			Names:    names,
			Prefixes: []string{},
		},
	})
}

func (s *ParameterServer) setParametersImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.SetParameters_Request, sender rcl_interfaces_srv.SetParametersServiceResponseSender) {
	s.mu.Lock()
	results := make([]rcl_interfaces_msg.SetParametersResult, 0, len(req.Parameters))
	changedParameters := make([]rcl_interfaces_msg.Parameter, 0, len(req.Parameters))
	for _, paramReq := range req.Parameters {
		param, ok := s.parameters[paramReq.Name]
		if !ok {
			results = append(results, rcl_interfaces_msg.SetParametersResult{
				Successful: false,
				Reason:     "invalid parameter name",
			})
			continue
		}
		if err := param.validate(paramReq.Value); err != nil {
			results = append(results, rcl_interfaces_msg.SetParametersResult{
				Successful: false,
				Reason:     err.Error(),
			})
			continue
		}
		param.set(paramReq.Value)
		for _, callback := range param.Callbacks {
			callback(param.Value)
		}
		results = append(results, rcl_interfaces_msg.SetParametersResult{Successful: true})
		changedParameters = append(changedParameters, paramReq)
	}
	s.mu.Unlock()

	err := s.parameterEventPublisher.Publish(&rcl_interfaces_msg.ParameterEvent{
		Stamp:             toRosTimestamp(time.Now()),
		Node:              s.nodeName,
		ChangedParameters: changedParameters,
	})
	if err != nil {
		s.logger.Warnf("failed to publish parameter event: %v", err)
	}

	sender.SendResponse(&rcl_interfaces_srv.SetParameters_Response{
		Results: results,
	})
}

func (s *ParameterServer) setParametersAtomicallyImpl(_ *rclgo.ServiceInfo, req *rcl_interfaces_srv.SetParametersAtomically_Request, sender rcl_interfaces_srv.SetParametersAtomicallyServiceResponseSender) {
	for _, paramReq := range req.Parameters {
		param, ok := s.parameters[paramReq.Name]
		if !ok {
			sender.SendResponse(&rcl_interfaces_srv.SetParametersAtomically_Response{
				Result: rcl_interfaces_msg.SetParametersResult{
					Successful: false,
					Reason:     fmt.Sprintf("set parameter %s: invalid parameter name", paramReq.Name),
				},
			})
			return
		}
		if err := param.validate(paramReq.Value); err != nil {
			sender.SendResponse(&rcl_interfaces_srv.SetParametersAtomically_Response{
				Result: rcl_interfaces_msg.SetParametersResult{
					Successful: false,
					Reason:     fmt.Sprintf("set parameter %s: %s", paramReq.Name, err.Error()),
				},
			})
			return
		}
	}

	s.mu.Lock()
	for _, paramReq := range req.Parameters {
		param := s.parameters[paramReq.Name]
		param.set(paramReq.Value)
		for _, callback := range param.Callbacks {
			callback(param.Value)
		}
	}
	s.mu.Unlock()

	err := s.parameterEventPublisher.Publish(&rcl_interfaces_msg.ParameterEvent{
		Stamp:             toRosTimestamp(time.Now()),
		Node:              s.nodeName,
		ChangedParameters: req.Parameters,
	})
	if err != nil {
		s.logger.Warnf("failed to publish parameter event: %v", err)
	}

	sender.SendResponse(&rcl_interfaces_srv.SetParametersAtomically_Response{
		Result: rcl_interfaces_msg.SetParametersResult{
			Successful: true,
		},
	})
}

type OnParameterChangedCallback func(newValue any)

type Parameter struct {
	Descriptor           rcl_interfaces_msg.ParameterDescriptor
	Value                any
	AdditionalValidation func(value rcl_interfaces_msg.ParameterValue) error
	Callbacks            []OnParameterChangedCallback
}

func (param *Parameter) asRclValue() rcl_interfaces_msg.ParameterValue {
	value := rcl_interfaces_msg.ParameterValue{Type: param.Descriptor.Type}
	switch param.Descriptor.Type {
	case rcl_interfaces_msg.ParameterType_PARAMETER_BOOL:
		value.BoolValue = param.Value.(bool)
	case rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER:
		value.IntegerValue = param.Value.(int64)
	case rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE:
		value.DoubleValue = param.Value.(float64)
	case rcl_interfaces_msg.ParameterType_PARAMETER_STRING:
		value.StringValue = param.Value.(string)
	case rcl_interfaces_msg.ParameterType_PARAMETER_BYTE_ARRAY:
		value.ByteArrayValue = param.Value.([]byte)
	case rcl_interfaces_msg.ParameterType_PARAMETER_BOOL_ARRAY:
		value.BoolArrayValue = param.Value.([]bool)
	case rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER_ARRAY:
		value.IntegerArrayValue = param.Value.([]int64)
	case rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE_ARRAY:
		value.DoubleArrayValue = param.Value.([]float64)
	case rcl_interfaces_msg.ParameterType_PARAMETER_STRING_ARRAY:
		value.StringArrayValue = param.Value.([]string)
	}
	return value
}

func (param *Parameter) set(value rcl_interfaces_msg.ParameterValue) {
	switch param.Descriptor.Type {
	case rcl_interfaces_msg.ParameterType_PARAMETER_BOOL:
		param.Value = value.BoolValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER:
		param.Value = value.IntegerValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE:
		param.Value = value.DoubleValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_STRING:
		param.Value = value.StringValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_BYTE_ARRAY:
		param.Value = value.ByteArrayValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_BOOL_ARRAY:
		param.Value = value.BoolArrayValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER_ARRAY:
		param.Value = value.IntegerArrayValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE_ARRAY:
		param.Value = value.DoubleArrayValue
	case rcl_interfaces_msg.ParameterType_PARAMETER_STRING_ARRAY:
		param.Value = value.StringArrayValue
	}
}

func (param *Parameter) validate(value rcl_interfaces_msg.ParameterValue) error {
	if param.Descriptor.ReadOnly {
		return errors.New("read-only parameter")
	}
	if param.Descriptor.Type != value.Type {
		return errors.New("invalid parameter type")
	}
	switch param.Descriptor.Type {
	case rcl_interfaces_msg.ParameterType_PARAMETER_INTEGER:
		intValue := value.IntegerValue
		if len(param.Descriptor.IntegerRange) == 0 {
			break
		}
		paramRange := param.Descriptor.IntegerRange[0]
		if intValue > paramRange.ToValue || intValue < paramRange.FromValue {
			return errors.New("invalid parameter value range")
		}
		if uint64(intValue-paramRange.FromValue)%paramRange.Step != 0 {
			return errors.New("invalid parameter value step")
		}
	case rcl_interfaces_msg.ParameterType_PARAMETER_DOUBLE:
		doubleValue := value.DoubleValue
		if len(param.Descriptor.FloatingPointRange) == 0 {
			break
		}
		paramRange := param.Descriptor.FloatingPointRange[0]
		if doubleValue > paramRange.ToValue || doubleValue < paramRange.FromValue {
			return errors.New("invalid parameter value range")
		}
		if uint64(doubleValue-paramRange.FromValue)%uint64(paramRange.Step) != 0 {
			return errors.New("invalid parameter value step")
		}
	}
	if param.AdditionalValidation != nil {
		if err := param.AdditionalValidation(value); err != nil {
			return err
		}
	}
	return nil
}

func toRosTimestamp(t time.Time) builtin_interfaces_msg.Time {
	return builtin_interfaces_msg.Time{
		Sec:     int32(t.Unix()),
		Nanosec: uint32(t.Nanosecond()),
	}
}
