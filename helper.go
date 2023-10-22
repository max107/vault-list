package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/vault-client-go"
)

func NewHelper() (*Helper, error) {
	// prepare a client with the given base address
	client, err := vault.New(vault.WithEnvironment())
	if err != nil {
		return nil, err
	}

	return &Helper{client: client}, nil
}

type Helper struct {
	client *vault.Client
}

func (h *Helper) ListSecretEngines(ctx context.Context) []string {
	resp, err := h.client.System.MountsListSecretsEngines(ctx)
	if err != nil {
		if vault.IsErrorStatus(err, http.StatusForbidden) {
			return nil
		}

		if vault.IsErrorStatus(err, http.StatusNotFound) {
			return nil
		}

		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(resp.Data))

	result := make([]string, len(resp.Data))
	i := 0
	for engine := range resp.Data {
		result[i] = engine
		i++
	}

	return result
}

func (h *Helper) ListKeys(ctx context.Context, engine, key string, prefix bool) {
	engine = strings.Trim(engine, "/")
	clean := strings.Trim(key, "/")

	metaKey := fmt.Sprintf("%s/metadata/%s/", engine, clean)

	keys, err := h.client.List(ctx, metaKey)
	if err != nil {
		if vault.IsErrorStatus(err, http.StatusForbidden) {
			return
		}

		if vault.IsErrorStatus(err, http.StatusNotFound) {
			return
		}

		log.Fatal(err)
	}

	values, ok := keys.Data["keys"].([]interface{})
	if !ok {
		return
	}

	var wg sync.WaitGroup

	for _, k := range values {
		ks := k.(string)

		var path []string
		if len(clean) > 0 {
			path = append(path, clean)
		}
		path = append(path, strings.Trim(ks, "/"))

		if strings.HasSuffix(ks, "/") {
			wg.Add(1)
			go func() {
				defer wg.Done()

				h.ListKeys(ctx, engine, strings.Join(path, "/"), prefix)
			}()
		} else {
			if prefix {
				fmt.Printf("%s/%s\n", engine, strings.Join(path, "/"))
			} else {
				fmt.Printf("%s\n", strings.Join(path, "/"))
			}
		}
	}

	wg.Wait()
}

func (h *Helper) List(ctx context.Context) {
	engines := h.ListSecretEngines(ctx)

	var wg sync.WaitGroup

	for _, v := range engines {
		wg.Add(1)

		go func(engine string) {
			defer wg.Done()

			h.ListKeys(ctx, engine, "", true)
		}(v)
	}

	wg.Wait()
}

func (h *Helper) ListEngine(ctx context.Context, engine string) {
	h.ListKeys(ctx, engine, "", false)
}
