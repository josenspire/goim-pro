package internal

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/edwingeng/slog"
)

const (
	// CriticalValue indicates when the low 36 bits are about to run out
	CriticalValue int64 = (1 << 36) * 80 / 100
	// RenewInterval indicates how often renew retries are performed
	RenewInterval int64 = 0x1FFFFFFF
	// PanicValue indicates when Next starts to panic
	PanicValue int64 = (1 << 36) * 98 / 100
)

// WUID is for internal use only.
type WUID struct {
	slog.Logger
	Tag         string
	Section     int8
	H28Verifier func(h28 int64) error

	N int64

	sync.Mutex
	Renew func() error
}

// NewWUID is for internal use only.
func NewWUID(tag string, logger slog.Logger, opts ...Option) (w *WUID) {
	w = &WUID{Tag: tag}
	if logger != nil {
		w.Logger = logger
	} else {
		w.Logger = slog.NewConsoleLogger()
	}
	for _, opt := range opts {
		opt(w)
	}
	return
}

// Next is for internal use only.
func (this *WUID) Next() int64 {
	x := atomic.AddInt64(&this.N, 1)
	v := x & 0x0FFFFFFFFF
	if v >= PanicValue {
		atomic.CompareAndSwapInt64(&this.N, x, x&(0x07FFFFFF<<36)|PanicValue)
		panic(fmt.Errorf("<wuid> the low 36 bits are about to run out. tag: %s", this.Tag))
	}
	if v >= CriticalValue && v&RenewInterval == 0 {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					this.Warnf("<wuid> panic, renew failed. tag: %s, reason: %+v", this.Tag, r)
				}
			}()

			err := this.RenewNow()
			if err != nil {
				this.Warnf("<wuid> renew failed. tag: %s, reason: %+v", this.Tag, err)
			} else {
				this.Infof("<wuid> renew succeeded. tag: %s", this.Tag)
			}
		}()
	}
	return x
}

// RenewNow reacquires the high 28 bits from your data store immediately
func (this *WUID) RenewNow() error {
	this.Lock()
	f := this.Renew
	this.Unlock()
	return f()
}

// Reset is for internal use only.
func (this *WUID) Reset(n int64) {
	if n < 0 {
		panic(fmt.Errorf("n should never be negative. tag: %s", this.Tag))
	}
	if this.Section == 0 {
		atomic.StoreInt64(&this.N, n)
	} else {
		atomic.StoreInt64(&this.N, n&0x0FFFFFFFFFFFFFFF|int64(this.Section)<<60)
	}
}

// VerifyH28 is for internal use only.
func (this *WUID) VerifyH28(h28 int64) error {
	if h28 <= 0 {
		return errors.New("h28 must be positive. tag: " + this.Tag)
	}

	if this.Section == 0 {
		if h28 > 0x07FFFFFF {
			return errors.New("h28 should not exceed 0x07FFFFFF. tag: " + this.Tag)
		}
	} else {
		if h28 > 0x00FFFFFF {
			return errors.New("h28 should not exceed 0x00FFFFFF. tag: " + this.Tag)
		}
	}

	if this.Section == 0 {
		if h28 == atomic.LoadInt64(&this.N)>>36 {
			return fmt.Errorf("h28 should be a different value other than %d. tag: %s", h28, this.Tag)
		}
	} else {
		if h28 == atomic.LoadInt64(&this.N)>>36&0x00FFFFFF {
			return fmt.Errorf("h28 should be a different value other than %d. tag: %s", h28, this.Tag)
		}
	}

	if this.H28Verifier != nil {
		if err := this.H28Verifier(h28); err != nil {
			return err
		}
	}

	return nil
}

// Option is for internal use only.
type Option func(*WUID)

// WithSection is for internal use only.
func WithSection(section int8) Option {
	if section < 1 || section > 7 {
		panic("section must be in between [1, 7]")
	}
	return func(w *WUID) {
		w.Section = section
	}
}

// WithH28Verifier is for internal use only.
func WithH28Verifier(cb func(h28 int64) error) Option {
	return func(w *WUID) {
		w.H28Verifier = cb
	}
}
