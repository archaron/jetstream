package jetstream

import (
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func RunBasicJetStreamServer() *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	opts.JetStream = true
	return natsserver.RunServer(&opts)
}

func TestConnection(t *testing.T) {

	t.Run("must fail on empty", func(t *testing.T) {
		v := viper.New()
		c, err := NewDefaultConfig(v)
		require.Nil(t, c)
		require.Error(t, err)
		require.Equal(t, err, ErrEmptyConfig)
	})

	t.Run("servers should be nil", func(t *testing.T) {
		v := viper.New()
		v.SetDefault(configKey+".url", "something")

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Nil(t, c.Servers)
	})

	t.Run("servers should be slice of string", func(t *testing.T) {
		v := viper.New()
		v.SetDefault(configKey+".servers", "something")

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Len(t, c.Servers, 1)
		require.Equal(t, c.Servers[0], "something")
	})

	t.Run("should be ok", func(t *testing.T) {
		v := viper.New()
		url := "something"
		v.SetDefault(configKey+".url", url)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Equal(t, c.Url, url)
	})

	t.Run("should fetch "+configKey+".servers_0", func(t *testing.T) {
		v := viper.New()
		v.SetDefault(configKey+".servers_0", nats.DefaultURL)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Len(t, c.Servers, 1)
		require.Equal(t, c.Servers[0], nats.DefaultURL)
	})

	t.Run("should fetch "+configKey+".servers", func(t *testing.T) {
		v := viper.New()
		v.SetDefault(configKey+".servers", []string{nats.DefaultURL})

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Len(t, c.Servers, 1)
		require.Equal(t, c.Servers[0], nats.DefaultURL)
	})

	t.Run("should fail for empty config", func(t *testing.T) {
		c, err := NewConnection(nil)
		require.Nil(t, c)
		require.EqualError(t, err, ErrEmptyConfig.Error())
	})

	t.Run("should fail client", func(t *testing.T) {
		v := viper.New()

		v.SetDefault(configKey+".url", nats.DefaultURL)

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.Equal(t, c.Url, nats.DefaultURL)

		cli, err := NewConnection(c)
		require.Nil(t, cli)
		require.Error(t, err)
	})

	t.Run("should not fail with test server", func(t *testing.T) {
		v := viper.New()
		serve := RunBasicJetStreamServer()
		defer serve.Shutdown()

		v.SetDefault(configKey+".url", serve.ClientURL())

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)

		cli, err := NewConnection(c)
		require.NoError(t, err)
		require.NotNil(t, cli)
	})

	t.Run("should be able to get jetstream context", func(t *testing.T) {
		v := viper.New()
		serve := RunBasicJetStreamServer()
		defer serve.Shutdown()

		v.SetDefault(configKey+".url", serve.ClientURL())

		c, err := NewDefaultConfig(v)
		require.NoError(t, err)

		cli, err := NewConnection(c)
		require.NoError(t, err)
		require.NotNil(t, cli)

		v.SetDefault(configKey+".jetstream.enable", true)

		sconf := NewDefaultStreamerConfig(v)
		require.NoError(t, err)
		require.NotNil(t, sconf)

		js, err := NewStreamer(cli, sconf)
		require.NoError(t, err)
		require.NotNil(t, js)

	})

}
