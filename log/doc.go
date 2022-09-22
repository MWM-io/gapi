/*
Package log provides a simple and extensible logger.

Usage

First, you need to build your logger instance, using log.NewDefaultLogger,
specifying what's your output (log.EntryWriter).
This default logger will provide timestamp, stacktrace and tracing for all your logs.
Use log.NewLogger if you need other options. (see EntryOption)

You can then use this logger instance in 3 different ways:
 - pass your logger instance across all the functions that need it
 - push it in a context, and retrieve it in the functions that need it (see With Context)
 - set it globally, and use the global instance when needed

	// Setup your log output
	var logOutput log.EntryWriter
	var ctx context.Context

	// use instance
	logger := log.NewDefaultLogger(logOutput)
	logger.Log("my log")

	// use global
	log.SetGlobalLogger(logger)
	log.Log("my log")

	// use with context
	ctx = log.WithContext(ctx, logger)
	log.LogC(ctx, "my log")

With Context

If you can, you should rather use logging function with context.
The global functions with context will first try to get the logger from the context,
but if none is found, it will use the global logger.

Package design

Its flow is based around the Entry struct.
You can log anything, (ie with log.LogAny) the input will be converted into an Entry.
A set of interfaces is available so the data you are logging can implement some of them to improve readability.

The second part is the EntryWriter interface.
It corresponds to the output of the logger.

*/
package log
