// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo/v2"
)

type LoggerSuite struct{}

var _ = gc.Suite(&LoggerSuite{})

func (*LoggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetDefaultContext()
}

func (s *LoggerSuite) TestRootLogger(c *gc.C) {
	root := loggo.Logger{}
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.LogLevel(), gc.Equals, loggo.WARNING)
	c.Check(root.IsErrorEnabled(), gc.Equals, true)
	c.Check(root.IsWarningEnabled(), gc.Equals, true)
	c.Check(root.IsInfoEnabled(), gc.Equals, false)
	c.Check(root.IsDebugEnabled(), gc.Equals, false)
	c.Check(root.IsTraceEnabled(), gc.Equals, false)
}

func (s *LoggerSuite) TestWithLabels(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")
	loggerWithLabels := logger.WithLabels(loggo.Labels{"foo": "bar"})
	loggerWithTagsAndLabels := logger.
		ChildWithTags("withTags", "tag1", "tag2").
		WithLabels(loggo.Labels{"hello": "world"})

	logger.Logf(loggo.INFO, "without labels")
	loggerWithLabels.Logf(loggo.INFO, "with labels")
	loggerWithTagsAndLabels.Logf(loggo.INFO, "with tags and labels")

	logs := writer.Log()
	c.Assert(logs, gc.HasLen, 3)
	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)
	c.Check(logs[1].Message, gc.Equals, "with labels")
	c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{"foo": "bar"})
	c.Check(logs[2].Message, gc.Equals, "with tags and labels")
	c.Check(logs[2].Labels, gc.DeepEquals, loggo.Labels{
		"logger-tags": "tag1,tag2",
		"hello":       "world",
	})
}

func (s *LoggerSuite) TestSetLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")

	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.TRACE)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.TRACE)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.TRACE)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, true)
	logger.SetLogLevel(loggo.DEBUG)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.INFO)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.INFO)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.WARNING)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.ERROR)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	// This is added for completeness, but not really expected to be used.
	logger.SetLogLevel(loggo.CRITICAL)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.CRITICAL)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.CRITICAL)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, false)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
}

func (s *LoggerSuite) TestModuleLowered(c *gc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *LoggerSuite) TestLevelsInherited(c *gc.C) {
	root := loggo.GetLogger("")
	first := loggo.GetLogger("first")
	second := loggo.GetLogger("first.second")

	root.SetLogLevel(loggo.ERROR)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.ERROR)

	first.SetLogLevel(loggo.DEBUG)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)

	second.SetLogLevel(loggo.INFO)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.INFO)

	first.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.INFO)
}

func (s *LoggerSuite) TestParent(c *gc.C) {
	logger := loggo.GetLogger("a.b.c")
	b := logger.Parent()
	a := b.Parent()
	root := a.Parent()

	c.Check(b.Name(), gc.Equals, "a.b")
	c.Check(a.Name(), gc.Equals, "a")
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Parent(), gc.DeepEquals, root)
}

func (s *LoggerSuite) TestParentSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	b := logger.Parent()

	c.Check(b, gc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChild(c *gc.C) {
	root := loggo.GetLogger("")

	a := root.Child("a")
	logger := a.Child("b.c")

	c.Check(a.Name(), gc.Equals, "a")
	c.Check(logger.Name(), gc.Equals, "a.b.c")
	c.Check(logger.Parent(), gc.DeepEquals, a.Child("b"))
}

func (s *LoggerSuite) TestChildSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a")
	b := logger.Child("b")

	c.Check(b, gc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChildSameContextWithLabels(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a", "parent")
	b := logger.ChildWithTags("b", "child")

	c.Check(ctx.GetAllLoggerLabels(), gc.DeepEquals, []string{"child", "parent"})
	c.Check(logger.Tags(), gc.DeepEquals, []string{"parent"})
	c.Check(b.Tags(), gc.DeepEquals, []string{"child"})
}

func (s *LoggerSuite) TestRoot(c *gc.C) {
	logger := loggo.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.DeepEquals, logger)
}

func (s *LoggerSuite) TestRootSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.DeepEquals, logger)
	c.Check(root.Child("a.b.c"), gc.Not(gc.DeepEquals), loggo.GetLogger("a.b.c"))
}
