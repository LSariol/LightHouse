package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LSariol/LightHouse/internal/builder"
	"github.com/LSariol/LightHouse/internal/cli"
	"github.com/LSariol/LightHouse/internal/config"
	"github.com/LSariol/LightHouse/internal/docker"
	"github.com/LSariol/LightHouse/internal/watcher"
	"github.com/LSariol/coveclient"
)

func main() {

	if err := config.Load(); err != nil {
		panic(err)
	}

	// Build Dependencies
	var coveClient *coveclient.Client = watcher.NewCoveClient()

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	dockerHandler, err := docker.NewHandler()
	if err != nil {
		log.Panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var builder *builder.Builder = builder.NewBuilder(dockerHandler)

	var watcher *watcher.Watcher = watcher.NewWatcher(coveClient, client, builder, ctx)

	go watcher.Run()

	//__________________________________________________________

	// containers, _ := dockerHandler.ListContainers(ctx)

	// fmt.Println(containers)

	cmd := cli.NewCLI(watcher)
	go cmd.Run()

	<-ctx.Done()
	log.Println("Shutting Down...")

}
