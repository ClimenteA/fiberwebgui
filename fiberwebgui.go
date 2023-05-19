package fiberwebgui

import (
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

type FiberApp interface {
	Listen(string) error
}

func getExistingPath(paths []string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func findBrowserOnLinux() string {
	paths := []string{
		"/usr/bin/google-chrome",
		"/usr/bin/microsoft-edge-stable",
		"/usr/bin/microsoft-edge",
		"/usr/bin/brave-browser",
	}
	return getExistingPath(paths)
}

func findBrowserOnMac() string {
	paths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	}
	return getExistingPath(paths)
}

func findBrowserOnWindows() string {
	paths := []string{
		"C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
		"C:\\Program Files\\Microsoft\\Edge\\Application\\msedge.exe",
		"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
		"C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application\\brave.exe",
	}
	return getExistingPath(paths)
}

func getBrowserPath() string {

	var browserPath string

	if runtime.GOOS == "windows" {
		browserPath = findBrowserOnWindows()
	}
	if runtime.GOOS == "linux" {
		browserPath = findBrowserOnLinux()
	}
	if runtime.GOOS == "darwin" {
		browserPath = findBrowserOnMac()
	}

	return browserPath

}

func getFreePort() int {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
	l, _ := net.ListenTCP("tcp", addr)
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	return port
}

func startBrowser(guiWg *sync.WaitGroup, browserClosed chan bool, browserPath, port string, windowSizeFlag string) {
	tempDir, _ := os.MkdirTemp("", "fiberwebgui")

	url := "http://127.0.0.1:" + port
	browserExecPath := browserPath
	userDataDir := "--user-data-dir=" + tempDir
	newWindow := "--new-window"
	noFirstRun := "--no-first-run"
	appUrl := "--app=" + url

	log.Println("Browser started with: ", browserExecPath, userDataDir, newWindow, noFirstRun, windowSizeFlag, appUrl)

	cmd := exec.Command(browserExecPath, userDataDir, newWindow, noFirstRun, windowSizeFlag, appUrl)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	os.RemoveAll(tempDir)
	log.Println("Browser stopped!")
	browserClosed <- true
	guiWg.Done()

}

func startFiberServer(guiWg *sync.WaitGroup, browserClosed chan bool, app FiberApp, port string) {
	log.Println("Server started...")

	go func() {
		closed := <-browserClosed
		if closed {
			log.Println("Server stopped!")
			guiWg.Done()
		}
	}()

	log.Fatal(app.Listen(":" + port))

}

func run(app FiberApp, windowSizeFlag string, port int) {
	browserPath := getBrowserPath()
	portstr := strconv.Itoa(port)
	browserClosed := make(chan bool)
	var guiWg sync.WaitGroup
	guiWg.Add(2)
	go startBrowser(&guiWg, browserClosed, browserPath, portstr, windowSizeFlag)
	go startFiberServer(&guiWg, browserClosed, app, portstr)
	guiWg.Wait()
}

func Run(app FiberApp) {
	port := getFreePort()
	run(app, "--start-maximized", port)
}

func RunOnPort(app FiberApp, port int) {
	run(app, "--start-maximized", port)
}

func RunWithSize(app FiberApp, width, height int) {
	port := getFreePort()
	windowSize := "--window-size=" + strconv.Itoa(width) + "," + strconv.Itoa(height)
	run(app, windowSize, port)
}

func RunWithSizeOnPort(app FiberApp, width, height, port int) {
	windowSize := "--window-size=" + strconv.Itoa(width) + "," + strconv.Itoa(height)
	run(app, windowSize, port)
}
