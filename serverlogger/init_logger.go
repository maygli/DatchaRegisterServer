package serverlogger

import (
	"datcha/servercommon"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"time"
)

const (
	LevelRequest  = slog.Level(slog.LevelError - 1)
	REQUEST_LEVEL = "REQUEST"
)

const (
	LOGGER_CFG_PATH    = "$.logger"
	LOG_FILE_PREFIX    = "log"
	LOG_JSON_FORMAT    = "json"
	LEVEL_DEBUG        = "debug"
	LEVEL_INFO         = "info"
	LEVEL_WARNING      = "warning"
	LEVEL_REQUEST      = "request"
	LEVEL_ERROR        = "error"
	LEVEL_DISABLED     = "disabled"
	LOG_LEVEL_DISABLED = slog.LevelError + 53
)

type LoggerConfiguration struct {
	LogFolder       string `json:"log_folder" env:"${SERVER_NAME}_LOG_FOLDER" default:"data/logs"`
	LogLevel        string `json:"log_level" env:"${SERVER_NAME}_LOG_LEVEL" default:"error"`
	IsSaveToFile    bool   `json:"is_save_to_file" env:"${SERVER_NAME}_LOG_SAVE_TO_FILE" default:"true"`
	IsPrintToStdout bool   `json:"is_print_to_std_out" env:"${SERVER_NAME}_LOG_PRINT_TO_STD_OUT" default:"true"`
	LogFormat       string `json:"log_format" env:"${SERVER_NAME}_LOG_FORMAT" default:"json"`
}

func NewLoggerConfiguration(cfgReader *servercommon.ConfigurationReader) (LoggerConfiguration, error) {
	cfg := LoggerConfiguration{}
	err := cfgReader.ReadConfiguration(&cfg, LOGGER_CFG_PATH)
	return cfg, err
}

func getLogLevel(levelStr string) (slog.Level, error) {
	switch levelStr {
	case LEVEL_DEBUG:
		return slog.LevelDebug, nil
	case LEVEL_INFO:
		return slog.LevelInfo, nil
	case LEVEL_WARNING:
		return slog.LevelWarn, nil
	case LEVEL_REQUEST:
		return LevelRequest, nil
	case LEVEL_ERROR:
		return slog.LevelError, nil
	case LEVEL_DISABLED:
		return LOG_LEVEL_DISABLED, nil
	default:
		msg := fmt.Sprintf("unrecognize log level '%s'", levelStr)
		return LOG_LEVEL_DISABLED, errors.New(msg)
	}
}

func getLogFileName(folder string, fileFormat string) string {
	t := time.Now()
	fileName := fmt.Sprintf("%s%s.%s", LOG_FILE_PREFIX, t.Format("20060102150405"), fileFormat)
	return path.Join(folder, fileName)
}

func InitLogger(cfgReader *servercommon.ConfigurationReader) error {
	cfg, err := NewLoggerConfiguration(cfgReader)
	if err != nil {
		return err
	}
	logLevel, err := getLogLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	if logLevel > slog.LevelError {
		// In case of log disabled create just default log with level more than error
		// Dont create any files of folders
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))
		return nil
	}
	if cfg.IsSaveToFile {
		isExists, err := servercommon.IsFileExists(cfg.LogFolder)
		if err != nil {
			return err
		}
		if !isExists {
			err = os.MkdirAll(cfg.LogFolder, os.ModePerm)
		}
	}
	var logWriter io.Writer = os.Stdout
	var stdOutWriter io.Writer = nil
	var fileWriter io.Writer = nil
	if cfg.IsPrintToStdout {
		stdOutWriter = os.Stdout
		logWriter = stdOutWriter
	}
	if cfg.IsSaveToFile {
		fileName := getLogFileName(cfg.LogFolder, cfg.LogFormat)
		fileWriter, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		logWriter = fileWriter
	}
	if fileWriter != nil && stdOutWriter != nil {
		logWriter = io.MultiWriter(fileWriter, stdOutWriter)
	}
	if logWriter == nil {
		logWriter = os.Stdout
		logLevel = LOG_LEVEL_DISABLED
	}
	var handler slog.Handler = nil
	opt := slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel := level.String()
				if level == LevelRequest {
					levelLabel = REQUEST_LEVEL
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
	if cfg.LogFormat == LOG_JSON_FORMAT {
		handler = slog.NewJSONHandler(logWriter, &opt)
	} else {
		handler = slog.NewTextHandler(logWriter, &opt)
	}
	slog.SetDefault(slog.New(handler))
	return nil
}
