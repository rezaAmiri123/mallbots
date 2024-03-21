package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/bojand/ghz/cmd/ghz"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rezaAmiri123/mallbots/customers/internal/agent"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("cannot load config:", err.Error())
	}

	appPrefix := os.Getenv("APP_PREFIX")
	agentCfg := agent.Config{}
	err = envconfig.Process(appPrefix, &agentCfg)
	if err != nil {
		log.Fatal("error reading environment variables: ", err)
	}

	ag, err := agent.NewAgent(agentCfg)
	if err != nil {
		log.Fatal("cannot load agent config:", err)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	ag.Shutdown()
}
