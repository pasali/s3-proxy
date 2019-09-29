package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"net/http"
)

var (
	configDirs = []string{"/etc/s3-proxy/", "$HOME/.s3-proxy/", "."}
)

type Config struct {
	Sites []Site
}

type Site struct {
	Host    string
	Bucket  string
	Users   []User
	Options Options
}

type User struct {
	Name     string
	Password string
}

type Options struct {
	CORS     bool
	Gzip     bool
	Website  bool
	Prefix   string
	ForceSSL bool
	Proxied  bool
}

func init() {
	for _, value := range configDirs {
		viper.AddConfigPath(value)
	}

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func ConfiguredProxyHandler() (http.Handler, error) {

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Sites) == 0 {
		return nil, errors.New("must specify one or more configurations")
	}

	handler := NewHostDispatchingHandler()

	for i, site := range cfg.Sites {
		err = site.validateWithHost()

		if err != nil {
			msg := fmt.Sprintf("%v in configuration at position %d", err, i)
			return nil, errors.New(msg)
		}

		handler.HandleHost(site.Host, createSiteHandler(site))
	}

	return handler, nil
}

func createSiteHandler(s Site) http.Handler {
	var handler http.Handler

	proxy := NewS3Proxy(s.Bucket)
	handler = NewProxyHandler(proxy, s.Options.Prefix)

	if s.Options.Website {
		cfg, err := proxy.GetWebsiteConfig()
		if err != nil {
			fmt.Printf("warning: site for bucket %s configured with "+
				"website option but received error when retrieving "+
				"website config\n\t%v", s.Bucket, err)
		} else {
			handler = NewWebsiteHandler(handler, cfg)
		}
	}

	if s.Options.CORS {
		handler = corsHandler(handler)
	}

	if s.Options.Gzip {
		handler = handlers.CompressHandler(handler)
	}

	if len(s.Users) > 0 {
		handler = NewBasicAuthHandler(s.Users, handler)
	} else {
		fmt.Printf("warning: site for bucket %s has no configured users\n", s.Bucket)
	}

	if s.Options.ForceSSL {
		handler = NewSSLRedirectHandler(handler)
	}

	if s.Options.Proxied {
		handler = handlers.ProxyHeaders(handler)
	}

	return handler
}

func corsHandler(next http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"HEAD", "GET", "OPTIONS"}),
	)(next)
}

func (s Site) validateWithHost() error {
	if s.Host == "" {
		return errors.New("Host not specified")
	}

	return nil
}
