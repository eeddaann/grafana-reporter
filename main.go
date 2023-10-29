package main

import (
	"flag"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	linkFlag := flag.String("url", "http://127.0.0.1:3000", "website url")
	flag.Parse()
	link := *linkFlag
	if path, exists := launcher.LookPath(); exists {
		log.Printf("using browser found at: %v", path)

		cmd := exec.Command(path)

		parser := launcher.NewURLParser()
		cmd.Stderr = parser
		utils.E(cmd.Start())
		u := launcher.New().Bin(path).Leakless(false).Headless(true).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()

		defer browser.Close()
		page := browser.MustPage("")
		var e proto.NetworkResponseReceived
		page.MustNavigate(link)
		wait := page.WaitEvent(&e)
		page.MustSetExtraHeaders("Authorization", "Basic YWRtaW46YWRtaW4=") // default grafana password
		wait()
		page.MustWaitStable()
		height := page.MustEval(`() => document.body.clientHeight`).String()
		heightToInt, err := strconv.ParseInt(height, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		page.MustSetViewport(1080, int(heightToInt), 8, false)
		page.MustEval(`() => document.body.style.zoom = 0.25`)
		time.Sleep(5 * time.Second)

		utils.Dump(e.Response.RequestHeaders)
		log.Printf("request headers: %v", e.Response.RequestHeaders)
		utils.Dump(e.Response.Status, e.Response.URL, e.Response.Headers)
		log.Printf("response status: %v url: %v", e.Response.Status, e.Response.URL)
		page.MustScreenshotFullPage("sample.png")

	}
}
