package logging

import "log"

type Logger struct {
	Level int
}

const (
	LevelDebug   = iota // Shows all messages
	LevelInfo           // Shows Info or lower
	LevelWarning        // Shows Warning or lower
	LevelError          // Shows Error or lower
	LevelOff            // Shows no messages
)

// logging.New(l int) returns a new logger initialized with the level set to l.
// Use the logger.LevelX constants to set the level.
func New(level int) *Logger {
	if level < LevelDebug {
		level = LevelDebug
	}
	return &Logger{Level: level}
}

// Debug messages are shown if level is set to LevelDebug.
func (l *Logger) Debug(s string, args ...any) {
	if l.Level <= LevelDebug {
		if len(args) == 0 {
			log.Println("DEBUG " + s)
		} else {
			log.Printf("DEBUG "+s, args...)
		}
	}
}

// Info messages are shown if level is set to LevelInfo or lower.
func (l *Logger) Info(s string, args ...any) {
	if l.Level <= LevelInfo {
		if len(args) == 0 {
			log.Println("INFO " + s)
		} else {
			log.Printf("INFO "+s, args...)
		}
	}
}

// Warning messages are shown if level is set to LevelWarning or lower.
func (l *Logger) Warning(s string, args ...any) {
	if l.Level <= LevelWarning {
		if len(args) == 0 {
			log.Println("WARNING " + s)
		} else {
			log.Printf("WARNING "+s, args...)
		}
	}
}

// Error messages are shown if level is set to LevelError or lower.
func (l *Logger) Error(s string, args ...any) {
	if l.Level <= LevelError {
		if l.Level <= LevelWarning {
			if len(args) == 0 {
				log.Println("ERROR " + s)
			} else {
				log.Printf("ERROR "+s, args...)
			}
		}
	}
}

// Fatal messages are always shown.
func (l *Logger) Fatal(s string, args ...any) {
	if len(args) == 0 {
		log.Fatalln("FATAL " + s)
	} else {
		log.Fatalf("FATAL "+s, args...)
	}
}
