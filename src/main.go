package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var availableResolutions = []string{"3840x2160", "1280x720"}

func getBestResolution(clientSpeedMbps float64) string {
	const highResThreshold = 50.0
	const lowResThreshold = 10.0

	if clientSpeedMbps >= highResThreshold {
		return availableResolutions[0]
	} else if clientSpeedMbps >= lowResThreshold {
		return availableResolutions[0]
	} else {
		return availableResolutions[1]
	}
}

func speedHandler(w http.ResponseWriter, r *http.Request) {

	videoPath := "videos/video-001-" + availableResolutions[0] + ".mp4"

	// Open the video file
	videoFile, err := os.Open(videoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer videoFile.Close()

	// Get the video file's information
	fileStat, err := videoFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// Set the content length header
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	// Copy the video data to the response writer
	io.Copy(w, videoFile)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	clientSpeedStr := r.Header.Get("X-Connection-Speed")

	clientSpeedMbps, err := strconv.ParseFloat(clientSpeedStr, 64)
	if err != nil {
		log.Println("Error parsing connection speed:", err)
		clientSpeedMbps = 0.0
	}

	// get best resolution by client speed
	bestResolution := getBestResolution(clientSpeedMbps)
	videoPath := "videos/video-001-" + bestResolution + ".mp4"
	fmt.Printf("Client Speed: %f / Best: %s\n", clientSpeedMbps, bestResolution)

	// Open the video file
	videoFile, err := os.Open(videoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer videoFile.Close()

	// Get the video file's information
	fileStat, err := videoFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// Set the content length header
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	// Copy the video data to the response writer
	io.Copy(w, videoFile)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/speed", speedHandler)
	http.HandleFunc("/video", videoHandler)
	http.Handle("/", fs)

	log.Println("Server listening on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
