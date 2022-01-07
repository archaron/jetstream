package jetstream

import (
	"github.com/im-kulikov/helium/module"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"strconv"
)

type (
	// Config alias
	Config = nats.Options

	StreamerConfig = []nats.JSOpt

	// Client alias
	Client = nats.Conn

	// Error is constant error
	Error string
)

const (
	configKey = "nats"

	// ErrEmptyConfig when given empty options
	ErrEmptyConfig = Error("nats empty config")
)

// Error returns error message string
func (e Error) Error() string { return string(e) }

var (
	// Module is default Nats client
	Module = module.Module{
		{Constructor: NewDefaultConfig},
		{Constructor: NewConnection},
		{Constructor: NewDefaultStreamerConfig},
		{Constructor: NewStreamer},
	}
)

func fetchAddresses(key string, v *viper.Viper) []string {
	var (
		addresses []string
	)

	for i := 0; ; i++ {
		addr := v.GetString(key + "_" + strconv.Itoa(i))
		if addr == "" {
			break
		}

		addresses = append(addresses, addr)
	}

	if len(addresses) == 0 {
		addresses = v.GetStringSlice(key)
	}

	return addresses
}

func NewDefaultConfig(v *viper.Viper) (*Config, error) {
	if !v.IsSet(configKey) {
		return nil, ErrEmptyConfig
	}

	var servers []string
	if addresses := fetchAddresses("nats.servers", v); len(addresses) > 0 {
		servers = addresses
	}

	v.SetDefault(configKey+".allow_reconnect", true)
	v.SetDefault(configKey+".max_reconnect", nats.DefaultMaxReconnect)
	v.SetDefault(configKey+".reconnect_wait", nats.DefaultReconnectWait)
	v.SetDefault(configKey+".reconnect_jitter", nats.DefaultReconnectJitter)
	v.SetDefault(configKey+".reconnect_jitter_tls", nats.DefaultReconnectJitterTLS)
	v.SetDefault(configKey+".timeout", nats.DefaultTimeout)
	v.SetDefault(configKey+".ping_interval", nats.DefaultPingInterval)
	v.SetDefault(configKey+".max_pings_out", nats.DefaultMaxPingOut)
	v.SetDefault(configKey+".sub_chan_len", nats.DefaultMaxChanLen)
	v.SetDefault(configKey+".reconnect_buf_size", nats.DefaultReconnectBufSize)
	v.SetDefault(configKey+".drain_timeout", nats.DefaultDrainTimeout)

	config := &Config{
		Url:                         v.GetString(configKey + ".url"),
		Servers:                     servers,
		NoRandomize:                 v.GetBool(configKey + ".no_randomize"),
		NoEcho:                      false,
		Name:                        v.GetString(configKey + ".name"),
		Verbose:                     v.GetBool(configKey + ".verbose"),
		Pedantic:                    v.GetBool(configKey + ".pedantic"),
		Secure:                      v.GetBool(configKey + ".secure"),
		AllowReconnect:              v.GetBool(configKey + ".allow_reconnect"),
		MaxReconnect:                v.GetInt(configKey + ".max_reconnect"),
		ReconnectWait:               v.GetDuration(configKey + ".reconnect_wait"),
		ReconnectJitter:             v.GetDuration(configKey + ".reconnect_jitter"),
		ReconnectJitterTLS:          v.GetDuration(configKey + ".reconnect_jitter_tls"),
		Timeout:                     v.GetDuration(configKey + ".timeout"),
		DrainTimeout:                v.GetDuration(configKey + ".drain_timeout"),
		FlusherTimeout:              v.GetDuration(configKey + ".flusher_timeout"),
		PingInterval:                v.GetDuration(configKey + ".ping_interval"),
		MaxPingsOut:                 v.GetInt(configKey + ".max_pings_out"),
		ReconnectBufSize:            v.GetInt(configKey + ".reconnect_buf_size"),
		SubChanLen:                  v.GetInt(configKey + ".sub_chan_len"),
		UserJWT:                     nil,
		Nkey:                        v.GetString(configKey + ".nkey"),
		User:                        v.GetString(configKey + ".user"),
		Password:                    v.GetString(configKey + ".password"),
		Token:                       v.GetString(configKey + ".token"),
		UseOldRequestStyle:          v.GetBool(configKey + ".use_old_request_style"),
		NoCallbacksAfterClientClose: v.GetBool(configKey + ".no_callbacks_after_client_close"),
		RetryOnFailedConnect:        v.GetBool(configKey + ".retry_on_failed_connect"),
		Compression:                 v.GetBool(configKey + ".compression"),
		InboxPrefix:                 v.GetString(configKey + ".inbox_prefix"),
	}

	//var tls *tls.Config

	if v.IsSet(configKey + ".tls") {
		cert := nats.ClientCert(v.GetString(configKey+".tls.cert"), v.GetString(configKey+".tls.key"))
		err := cert(config)
		if err != nil {
			return nil, err
		}

		if v.IsSet(configKey + ".tls.cacert") {
			ca := nats.RootCAs(v.GetString(configKey + ".tls.cacert"))
			err := ca(config)
			if err != nil {
				return nil, err
			}
		}

	}

	return config, nil
}

func NewDefaultStreamerConfig(v *viper.Viper) StreamerConfig {
	jsKey := configKey + ".jetstream"
	if !v.IsSet(jsKey) {
		return nil
	}

	opts := make(StreamerConfig, 0)

	if v.IsSet(jsKey + ".prefix") {
		opts = append(opts, nats.APIPrefix(v.GetString(jsKey+".prefix")))
	}

	if v.IsSet(jsKey + ".publish_async_max_pending") {
		opts = append(opts, nats.PublishAsyncMaxPending(v.GetInt(jsKey+".publish_async_max_pending")))
	}

	return opts
}

// NewConnection of nats client
func NewConnection(opts *Config) (bus *Client, err error) {
	if opts == nil {
		return nil, ErrEmptyConfig
	}

	if bus, err = opts.Connect(); err != nil {
		return nil, err
	}

	return bus, nil
}

// NewStreamer is jetstream client
func NewStreamer(bus *Client, opts StreamerConfig) (nats.JetStreamContext, error) {
	return bus.JetStream(opts...)
}
