package logger

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/ivankorobkov/di"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Module(m *di.Module) {
	m.AddConstructor(New)
	m.MarkPackageDep(Config{})
}

type Fields map[string]interface{}
type Payload map[string]interface{}

type Logger interface {
	Prefix(prefix string) Logger
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Panic(err interface{})
	PanicHttp(err interface{}, r *http.Request)
}

type Config struct {
	DefaultPrefix string `yaml:"DefaultPrefix"`
	LogFile       string `yaml:"LogFile"`
	DSN           string `yaml:"DNS"`
	Json          bool   `yaml:"Json"`
}

var sharedRavenClient *raven.Client
var jsonOut bool

func New(config Config) Logger {
	initLogger(config)
	return &logger{prefix: config.DefaultPrefix}
}

func initLogger(config Config) {
	// Check for already inited
	if sharedRavenClient != nil {
		return
	}

	// Init sentry client
	client, err := raven.NewClient(config.DSN, nil)
	if err != nil {
		panic(err)
	}
	sharedRavenClient = client

	if config.Json {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		jsonOut = true
	}

	// Init logrus std logger
	if len(config.LogFile) > 0 {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		jsonOut = true

		output := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    256, // megabytes
			MaxBackups: 7,
			MaxAge:     180, //days
		}
		_, err = output.Write(nil)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(output)
	}
}

type logger struct {
	prefix string
}

func (l *logger) Prefix(prefix string) Logger {
	return &logger{prefix}
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.newEntry(args).Debug(msg)
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.newEntry(args).Info(msg)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.newEntry(args).Error(msg)
}

func (l *logger) Panic(err interface{}) {
	l.sendToSentry(err, nil)
}

func (l *logger) PanicHttp(err interface{}, r *http.Request) {
	l.sendToSentry(err, r)
}

func (l *logger) sendToSentry(err interface{}, r *http.Request) {
	e, ok := err.(error)
	if !ok {
		e = fmt.Errorf("%v", err)
	}

	l.Error("Trace", Fields{"trace": string(debug.Stack()), "err": e})

	trace := raven.NewStacktrace(4, 3, nil)
	var packet *raven.Packet
	if r != nil {
		packet = raven.NewPacket(e.Error(), raven.NewException(e, trace), raven.NewHttp(r))
	} else {
		packet = raven.NewPacket(e.Error(), raven.NewException(e, trace))
	}

	eventId, ch := sharedRavenClient.Capture(packet, nil)
	if err = <-ch; err != nil {
		l.Error("Failed to send an error to sentry", Fields{"err": e})
		return
	}
	l.Error("Sent an error to sentry", Fields{"eventId": eventId, "err": e})
}

func (l *logger) newEntry(args []interface{}) *logrus.Entry {
	if len(args) == 0 {
		return logrus.WithField("prefix", l.prefix)
	}

	entry := logrus.WithField("prefix", l.prefix)
	for _, arg := range args {
		switch v := arg.(type) {
		case Fields:
			entry = entry.WithFields(logrus.Fields(arg.(Fields)))
		case Payload:
			if jsonOut {
				entry = entry.WithField("payload", arg.(Payload))
			} else {
				entry = entry.WithFields(logrus.Fields(arg.(Payload)))
			}
		default:
			fmt.Println("logger: unsupported log type ", v)
		}
	}
	return entry
}
