package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robertacosta/go-integration-example/server"
	"github.com/spf13/cobra"
)

var testServerCmd = &cobra.Command{
	Use:   "testserver",
	Short: "Runs http test server",
	Run:   startTestServer,
}

func startTestServer(cmd *cobra.Command, args []string) {
	testServiceAddr := os.Getenv("TEST_SERVICE_ADDRESS")

	testConfig := server.TestConfig{
		TestServiceAddr: testServiceAddr,
	}

	stop := make(chan os.Signal)
	defer close(stop)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	testSrv, err := server.NewTestServer(testConfig)
	if err != nil {
		log.Printf("unable to initialize test service: %s", err)
	}

	testSrv.Start()

	//// wait for sigint or sigterm
	<-stop

	log.Println("shutting down servers")
	ctx, cncl := context.WithTimeout(context.Background(), time.Duration(5000)*time.Millisecond)
	defer cncl()
	testSrv.Stop(ctx)
}
