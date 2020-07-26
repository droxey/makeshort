package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/chromedp/chromedp"
)

// Profile 1 is the default profile installed with Chrome.
// Most users won't edit this setting.
// Provide a way to override the default, just in case.
var flagProfileName = flag.String("profileName", "Profile 1", "Google Chrome profile name")
var flagSlug = flag.String("slug", "", "URL slug; appears after https://make.sc/")
var flagURL = flag.String("url", "", "URL you want to shorten")
var flagEmail = flag.String("email", os.Getenv("MS_EMAIL"), "Make School email address")
var flagPassword = flag.String("password", os.Getenv("MS_PASSWORD"), "Make School password")

func main() {
	flag.Parse()
	url := "https://www.makeschool.com/login"

	// Get the path to $HOME on any OS.
	usr, _ := user.Current()

	// Path to the default profile on macOS.
	profilePath := usr.HomeDir + `/Library/Application Support/Google/Chrome/Default`

	// Attributes about the browser used to scrape the website.
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("profile-directory", *flagProfileName),
		chromedp.Flag("disable-sync", false),
		chromedp.Flag("save-password", false),
		chromedp.UserDataDir(profilePath),

		// Browser attributes:
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,

		chromedp.Headless,
		chromedp.DisableGPU,
	}

	// Use advanced options declared above when executing chromedp.
	execCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a context with debugging output.
	chromeCtx, cancel := chromedp.NewContext(execCtx)
	defer cancel()

	// List of tasks to run.
	var res string
	err := chromedp.Run(chromeCtx, createShortLink(url, res))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	time.Sleep(5)
}

func createShortLink(urlStr string, res string) chromedp.Tasks {
	emailInputSel := `/html/body/main/div[2]/article/div/div/div[2]/div[1]/form[3]/div[2]/label[1]/input`
	passwordInputSel := `/html/body/main/div[2]/article/div/div/div[2]/div[1]/form[3]/div[2]/label[2]/input`
	loadingSel := `//*[@aria-label="loading animation"]`
	slugInputSel := `/html/body/div/div/form/p[1]/input`
	urlInputSel := `/html/body/div/div/form/p[2]/input`

	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(emailInputSel),
		chromedp.SendKeys(emailInputSel, *flagEmail),
		chromedp.SendKeys(passwordInputSel, *flagPassword),
		chromedp.Submit(emailInputSel),
		chromedp.WaitNotPresent(loadingSel),
		chromedp.Navigate("https://www.makeschool.com/admin/short_links/new"),
		chromedp.SendKeys(slugInputSel, *flagSlug),
		chromedp.SendKeys(urlInputSel, *flagURL),
		chromedp.Submit(`/html/body/div/div/form`),
		chromedp.WaitVisible(`*[@class='ShortLinks__Button']`),
	}
}
