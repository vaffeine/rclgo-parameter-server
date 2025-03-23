package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/vaffeine/rclgo-parameter-server/pkg/parameter_server"
)

func main() {
	rclArgs, _, err := rclgo.ParseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}

	ctx, err := rclgo.NewContext(rclgo.ClockTypeROSTime, rclArgs)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	node, err := ctx.NewNode("rclgo_parameter_server_example", "")
	if err != nil {
		panic(err)
	}
	defer node.Close()

	parameterServer, err := parameter_server.New(
		node,
		parameter_server.Parameter{
			Descriptor: parameter_server.ParameterDescriptor{
				Name:          "bitrate",
				Type:          parameter_server.ParameterType_PARAMETER_INTEGER,
				Description:   "Bitrate to use for compression of video streams from each camera",
				ReadOnly:      false,
				DynamicTyping: false,
				IntegerRange: []parameter_server.IntegerRange{{
					FromValue: 0,
					ToValue:   4294967295,
					Step:      1,
				}},
			},
			Value: int64(8000000),
			Callbacks: []parameter_server.OnParameterChangedCallback{
				func(newValue any) { fmt.Printf("new value: %v\n", newValue) },
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer parameterServer.Close()

	if err := ctx.Spin(context.Background()); err != nil {
		panic(err)
	}
}
