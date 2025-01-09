package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AppLogger struct {
	mu         sync.Mutex
	wstream    *os.File
	dir        string
	prefix     string
	useConsole bool
	theday     string
}

var (
	applog *AppLogger
)

func NewLogger(dir, prefix string, useConsole bool) *AppLogger {
	logger := &AppLogger{
		dir:        dir,
		prefix:     prefix,
		useConsole: useConsole,
		theday:     time.Now().Format("02"),
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	logger.createWStream(false)
	return logger
}

func (l *AppLogger) createWStream(flNewDay bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	filename := fmt.Sprintf("%s_%s_%02d%02d%02d.log", l.prefix, now.Format("20060102"), now.Hour(), now.Minute(), now.Second())
	fpath := filepath.Join(l.dir, filename)

	file, err := os.Create(fpath)
	if err != nil {
		fmt.Printf("ERROR: Cannot create log stream: %v\n", fpath)
		l.wstream = nil
		return
	}

	l.wstream = file

	sstamp := now.Format("20060102 150405 ")
	l.wstream.WriteString("====================================================================\n")
	if flNewDay {
		l.wstream.WriteString(sstamp + "BEGIN LOG for a NEW DAY\n")
	} else {
		l.wstream.WriteString(sstamp + "BEGIN LOG\n")
	}
	l.wstream.WriteString("====================================================================\n")
}

func (l *AppLogger) WriteString(str string) {
	now := time.Now()
	sstamp := now.Format("20060102 150405 .000 ")

	if l.useConsole {
		fmt.Println(sstamp + str)
	}

	if l.wstream == nil {
		return
	}

	if now.Format("02") != l.theday {
		l.wstream.WriteString("====================================================================\n")
		l.wstream.WriteString(sstamp + "END LOG for the OLD DAY\n")
		l.wstream.WriteString("====================================================================\n")
		l.wstream.Close()
		l.createWStream(true)
		l.theday = now.Format("02")
	}

	l.wstream.WriteString(sstamp + str + "\n")
}

func (l *AppLogger) WriteLine(v ...interface{}) {
	sargs := fmt.Sprint(v...)
	l.WriteString(sargs)
}

func (l *AppLogger) End() {
	if l.wstream != nil {
		l.wstream.Close()
	}
}

func AppLog(v ...any) {
	if applog == nil {
		applog = NewLogger("./log", "log", true)
	}
	applog.WriteLine(v...)
}

// ------------------------------ HANDLE REQUEST RESPONSE REST API -----------------------------------------------------

func BodyToJSONSlices(body io.Reader) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	dec := json.NewDecoder(body)

	for {
		var body map[string]interface{}
		if err := dec.Decode(&body); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to decode request body: %w", err)
		}

		result = append(result, body)
	}

	return result, nil
}

func LogRequest(r *http.Request, bodyBytes []byte) map[string]interface{} {
	bodyJSON, err := BodyToJSONSlices(bytes.NewBuffer(bodyBytes))
	if err != nil {
		rawBody := string(bodyBytes)
		rawBody = strings.ReplaceAll(rawBody, "\n", "")
		rawBody = strings.ReplaceAll(rawBody, "\r", "")
		bodyJSON = []map[string]interface{}{
			{"raw_body": rawBody},
		}
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr
	}

	userAgent := r.UserAgent()

	currentTime := time.Now().Format("2006-01-02 15:04:05 .000")

	headersCopy := make(http.Header)
	for key, values := range r.Header {
		if key == "Authorization" {
			continue
		}
		headersCopy[key] = values
	}

	requestLog := map[string]interface{}{
		"body":      bodyJSON,
		"method":    r.Method,
		"url":       r.URL.String(),
		"headers":   headersCopy,
		"client_ip": clientIP,
		"time":      currentTime,
		"agent":     userAgent,
	}

	return requestLog
}

func LogResponseError(status string, message string) map[string]interface{} {

	responseLog := map[string]interface{}{
		"status":  status,
		"message": message,
	}
	return responseLog
}

func LogResponseSuccess(data interface{}) map[string]interface{} {
	responseLog := map[string]interface{}{
		"status": "success",
		"data":   data,
	}
	return responseLog
}

func LogResponseSuccessMap(responseLog map[string]interface{}) map[string]interface{} {
	sensitiveKeys := map[string]bool{"token": true, "pwd": true, "password": true}

	response := make(map[string]interface{})
	for k, v := range responseLog {
		if _, exists := sensitiveKeys[k]; exists {
			response[k] = "*"
		} else {
			response[k] = v
		}
	}

	return response
}

func LogRequestResponse(requestLog, responseLog map[string]interface{}) string {

	logData := map[string]interface{}{
		"request":  requestLog,
		"response": responseLog,
	}

	logJSON, err := json.Marshal(logData)
	if err != nil {
		log.Printf("Failed to encode log to JSON: %v", err)
		return `{"error": "failed to generate log"}`
	}

	return string(logJSON)
}

func respondWithError(requestLog map[string]interface{}, message string, err error) (interface{}, error) {
	emptyErr := fmt.Sprintf("%s: %v", message, err)
	if err == nil {
		emptyErr = message
	} else if message == err.Error() {
		emptyErr = message
	}

	AppLog("ini apa", emptyErr)
	responseLog := LogResponseError("error", emptyErr)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return nil, fmt.Errorf("%s", message)
}

func respondWithSuccess(requestLog map[string]interface{}, data interface{}) (interface{}, error) {
	responseLog := LogResponseSuccess(data)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func respondWithSuccessStruct(requestLog map[string]interface{}, mapData map[string]interface{}) (interface{}, error) {
	checkMapData := LogResponseSuccessMap(mapData)
	AppLog(LogRequestResponse(requestLog, checkMapData))
	return mapData, nil
}
