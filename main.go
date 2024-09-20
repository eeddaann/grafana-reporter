package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/signintech/gopdf"
)

type Slide struct {
	Link string `json:"link"`
	Zoom string `json:"zoom"`
	Name string `json:"name"`
}

func GetFolderDate() string {
	// Use layout string for time format.
	const layout = "01-02-2006"
	// Place now in the string.
	t := time.Now()
	return t.Format(layout)
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

func GeneratePdf(path string, images []string) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, image := range images {
		pdf.AddPage()
		pdf.Image(image, 0, 0, nil)
	}
	pdf.WritePdf(filepath.Join(path, "report.pdf"))
}

func takeScreenshot(slide *Slide, browser *rod.Browser) {
	log.Printf("navigating to: %v", slide.Link)
	page := browser.MustPage("")
	var e proto.NetworkResponseReceived
	page.MustSetExtraHeaders("Authorization", "Basic YWRtaW46YWRtaW4=") // default grafana password
	page.MustNavigate(slide.Link)
	wait := page.WaitEvent(&e)
	wait()
	page.MustWaitStable()
	height := page.MustEval(`() => document.body.clientHeight`).String()
	heightToInt, err := strconv.ParseInt(height, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	page.MustSetViewport(1080, int(heightToInt), 8, false)
	page.MustEval(`() => document.body.style.zoom =` + slide.Zoom)
	time.Sleep(5 * time.Second)

	utils.Dump(e.Response.Status, e.Response.URL, e.Response.Headers)
	log.Printf("response status: %v url: %v", e.Response.Status, e.Response.URL)
	page.MustScreenshotFullPage(filepath.Join(GetFolderDate(), slide.Name+".png"))
}

func main() {
	configPath := flag.String("conf", "conf.json", "config path")
	flag.Parse()
	if path, exists := launcher.LookPath(); exists {
		log.Printf("using browser found at: %v", path)

		cmd := exec.Command(path)

		parser := launcher.NewURLParser()
		cmd.Stderr = parser
		utils.E(cmd.Start())
		u := launcher.New().Bin(path).Leakless(false).Headless(true).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.Close()
		content, err := os.ReadFile(*configPath)
		if err != nil {
			log.Fatal(err)
		}
		var slides []Slide
		if err := json.Unmarshal(content, &slides); err != nil {
			log.Fatal(err)
		}
		err = makeDirectoryIfNotExists(GetFolderDate())
		if err != nil {
			log.Fatal(err)
		}
		var images []string
		for _, slide := range slides {
			takeScreenshot(&slide, browser)
			images = append(images, filepath.Join(GetFolderDate(), slide.Name+".png"))
		}
		GeneratePdf(GetFolderDate(), images)
	}
}
