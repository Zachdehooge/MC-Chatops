package functions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

var startTime = time.Now()

func Uptime() string {
	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ServerStatus() string {

	godotenv.Load()
	log.Print("Getting bot token from .env file")
	var serverip = os.Getenv("SERVERIP")

	log.Print("Fetching server information")
	url := fmt.Sprintf("http://%s", serverip)
	cmd := exec.Command("curl", "-o", "/dev/null", "-s", "-w", "%{http_code}", url)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}

	return string(out)
}

func StartServer() string {
	return "Starting Server..."
}

func StopServer() string {
	return "Stopping Server..."
}
