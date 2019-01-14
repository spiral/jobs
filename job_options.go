package jobs

import "time"

// Options carry information about how to handle given job.
type Options struct {
	// Pipeline manually specified pipeline.
	Pipeline string `json:"pipeline,omitempty"`

	// Delay defines time duration to delay execution for. Defaults to none.
	Delay int `json:"delay,omitempty"`

	// Maximum job retries. Defaults to none.
	MaxAttempts int `json:"maxAttempts,omitempty"`

	// RetryDelay defines for how long job should be waiting until next retry. Defaults to none.
	RetryDelay int `json:"retryDelay,omitempty"`

	// Reserve defines for how broker should wait until treating job are failed. Defaults to 30 min.
	Timeout int `json:"timeout,omitempty"`
}

// Merge merges job options.
func (o *Options) Merge(from *Options) {
	if o.Pipeline == "" {
		o.Pipeline = from.Pipeline
	}

	if o.MaxAttempts == 0 {
		o.MaxAttempts = from.MaxAttempts
	}

	if o.Timeout == 0 {
		o.Timeout = from.Timeout
	}

	if o.RetryDelay == 0 {
		o.RetryDelay = from.RetryDelay
	}

	if o.Delay == 0 {
		o.Delay = from.Delay
	}
}

// CanRetry must return true if broker is allowed to re-run the job.
func (o *Options) CanRetry(attempts int) bool {
	// MaxAttempts 1 and 0 has identical effect
	return o.MaxAttempts > (attempts + 1)
}

// RetryDuration returns retry delay duration in a form of time.Duration.
func (o *Options) RetryDuration() time.Duration {
	return time.Second * time.Duration(o.RetryDelay)
}

// DelayDuration returns delay duration in a form of time.Duration.
func (o *Options) DelayDuration() time.Duration {
	return time.Second * time.Duration(o.Delay)
}

// TimeoutDuration returns timeout duration in a form of time.Duration.
func (o *Options) TimeoutDuration() time.Duration {
	if o.Timeout == 0 {
		return 30 * time.Minute
	}

	return time.Second * time.Duration(o.Timeout)
}
