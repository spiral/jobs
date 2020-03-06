package sqs

import (
	"github.com/spiral/jobs/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	pipe = &jobs.Pipeline{
		"broker": "sqs",
		"name":   "default",
		"queue":  "test",
		"declare": map[string]interface{}{
			"MessageRetentionPeriod": 86400,
		},
	}

	cfg = &Config{
		Key:      "api-key",
		Secret:   "api-secret",
		Region:   "us-west-1",
		Endpoint: "http://localhost:9324",
	}
)

func TestBroker_Init(t *testing.T) {
	b := &Broker{}
	ok, err := b.Init(cfg)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestBroker_StopNotStarted(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)

	b.Stop()
}

func TestBroker_Register(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	assert.NoError(t, b.Register(pipe))
}

func TestBroker_RegisterInvalid(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	assert.Error(t, b.Register(&jobs.Pipeline{
		"broker": "sqs",
		"name":   "default",
	}))
}

func TestBroker_Register_Twice(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	assert.NoError(t, b.Register(pipe))
	assert.Error(t, b.Register(pipe))
}

func TestBroker_Consume_Nil_BeforeServe(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)
	assert.NoError(t, b.Consume(pipe, nil, nil))
}

func TestBroker_Consume_Undefined(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)

	assert.Error(t, b.Consume(pipe, nil, nil))
}

func TestBroker_Consume_BeforeServe(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	exec := make(chan jobs.Handler)
	err := func(id string, j *jobs.Job, err error) {}

	assert.NoError(t, b.Consume(pipe, exec, err))
}

func TestBroker_Consume_Serve_Nil_Stop(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	b.Consume(pipe, nil, nil)

	wait := make(chan interface{})
	go func() {
		assert.NoError(t, b.Serve())
		close(wait)
	}()
	time.Sleep(time.Millisecond * 100)
	b.Stop()

	<-wait
}

func TestBroker_Consume_Serve_Stop(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	exec := make(chan jobs.Handler)
	err := func(id string, j *jobs.Job, err error) {}

	b.Consume(pipe, exec, err)

	wait := make(chan interface{})
	go func() {
		assert.NoError(t, b.Serve())
		close(wait)
	}()
	time.Sleep(time.Millisecond * 100)
	b.Stop()

	<-wait
}

func TestBroker_Consume_Serve_InvalidQueue(t *testing.T) {
	pipe := &jobs.Pipeline{
		"broker": "sqs",
		"name":   "default",
		"queue":  "invalid",
		"declare": map[string]interface{}{
			"VisibilityTimeout": "invalid",
		},
	}

	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	exec := make(chan jobs.Handler)
	err := func(id string, j *jobs.Job, err error) {}

	b.Consume(pipe, exec, err)

	assert.Error(t, b.Serve())
}

func TestBroker_PushToNotRunning(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	_, err := b.Push(pipe, &jobs.Job{})
	assert.Error(t, err)
}

func TestBroker_StatNotRunning(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)
	b.Register(pipe)

	_, err := b.Stat(pipe)
	assert.Error(t, err)
}

func TestBroker_PushToNotRegistered(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)

	ready := make(chan interface{})
	b.Listen(func(event int, ctx interface{}) {
		if event == jobs.EventBrokerReady {
			close(ready)
		}
	})

	go func() { assert.NoError(t, b.Serve()) }()
	defer b.Stop()

	<-ready

	_, err := b.Push(pipe, &jobs.Job{})
	assert.Error(t, err)
}

func TestBroker_StatNotRegistered(t *testing.T) {
	b := &Broker{}
	b.Init(cfg)

	ready := make(chan interface{})
	b.Listen(func(event int, ctx interface{}) {
		if event == jobs.EventBrokerReady {
			close(ready)
		}
	})

	go func() { assert.NoError(t, b.Serve()) }()
	defer b.Stop()

	<-ready

	_, err := b.Stat(pipe)
	assert.Error(t, err)
}
