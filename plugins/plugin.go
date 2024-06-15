package plugins

import (
	"context"
	"os/exec"
	"vimana/log"
	"vimana/plugins/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type SpacecorePlugin interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status(ctx context.Context) error
	Logs(ctx context.Context) error
}

// Here, we're binding our MySpacecore struct to the plugin interface

// Plugin management
type SpacecorePluginIm struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// concrete implementation, written in Go. This is only used for plugins
	plugin.GRPCPlugin
	Impl proto.ScPluginServer
}

// this is a client-side implementation of the SpacecorePlugin interface
type GRPCClient struct{ client proto.ScPluginClient }

func (p *SpacecorePluginIm) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewScPluginClient(c)}, nil
}

func (p *GRPCClient) Start(ctx context.Context, req *proto.StartRequest) (*proto.StartResponse, error) {
	logger := log.GetLogger(ctx)
	msg, err := p.client.Start(ctx, req)
	logger.Infof("plugin start msg: %v\n", msg)
	return msg, err
}

func (p *GRPCClient) Stop(ctx context.Context, req *proto.StopRequest) (*proto.StopResponse, error) {
	msg, err := p.client.Stop(ctx, req)
	return msg, err
}

func (p *GRPCClient) Status(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	return p.client.Status(ctx, req)
}

func (p *GRPCClient) Logs(ctx context.Context, req *proto.LogsRequest) (*proto.LogsResponse, error) {
	return p.client.Logs(ctx, req)
}

// We're a host! Start by launching the plugin process.
func GetPluginClient(pluginPath string) *plugin.Client {
	pluginMap := map[string]plugin.Plugin{
		"spacecore": &SpacecorePluginIm{},
	}

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "SPACECORE_PLUGIN",
			MagicCookieValue: "v1",
		},
		Plugins:          pluginMap,
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}

	client := plugin.NewClient(clientConfig)
	// defer client.Kill()

	return client
}

// func StartPlugin(client *plugin.Client) (*GRPCClient, error) {
// need a better name representing this function

func SpacecoreGRPCClient(client *plugin.Client) (*GRPCClient, error) {
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return &GRPCClient{}, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("spacecore")
	if err != nil {
		return &GRPCClient{}, err
	}

	// We should have a SpacecorePlugin now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	spacecore := raw.(*GRPCClient) // SpacecorePlugin)

	return spacecore, nil
}
