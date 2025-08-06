package inmemoryQueue

import (
	"queue-broker/api/core"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestGetNoSuchQueue(t *testing.T) {
	r := NewInMemoryQueueRepository(1, 1)

	type args struct {
		queueName  string
		timeoutSec int
	}
	tests := []struct {
		name    string
		args    args
		want    core.Message
		wantErr bool
	}{
		{
			name: "Unexsisting",
			args: args{
				queueName:  "",
				timeoutSec: 0,
			},
			want:    core.Message{},
			wantErr: true,
		},
		{
			name: "Unexsisting",
			args: args{
				queueName:  "121312",
				timeoutSec: 0,
			},
			want:    core.Message{},
			wantErr: true,
		},
		{
			name: "Unexsisting",
			args: args{
				queueName:  "fdsfdsfsd",
				timeoutSec: 0,
			},
			want:    core.Message{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Get(tt.args.queueName, tt.args.timeoutSec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutNoSuchQueueAndOverLimit(t *testing.T) {
	r := NewInMemoryQueueRepository(1, 2)

	type args struct {
		queueName string
		message   core.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Unexsisting",
			args: args{
				queueName: "adssadasd",
				message:   core.Message{},
			},
			wantErr: false,
		},
		{
			name: "Unexsisting",
			args: args{
				queueName: "121313",
				message:   core.Message{},
			},
			wantErr: false,
		},
		{
			name: "Unexsisting",
			args: args{
				queueName: "mfkdomfosm",
				message:   core.Message{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Put(tt.args.queueName, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPutQueueOverflow(t *testing.T) {
	r := NewInMemoryQueueRepository(2, 1)

	type args struct {
		queueName string
		message   core.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Has space",
			args: args{
				queueName: "1",
				message:   core.Message{byte(65)},
			},
			wantErr: false,
		},
		{
			name: "Has space",
			args: args{
				queueName: "1",
				message:   core.Message{byte(65)},
			},
			wantErr: false,
		},
		{
			name: "No free space",
			args: args{
				queueName: "1",
				message:   core.Message{byte(65)},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Put(tt.args.queueName, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCorrectWithoutTimeout(t *testing.T) {
	r := NewInMemoryQueueRepository(5, 5)
	qName := "Name1"
	msg := core.Message("mmasdas")

	// Создадим очередь и сразу опустошим
	if err := r.Put(qName, msg); err != nil {
		t.Errorf("Put() error = %v", err)
	}
	if _, err := r.Get(qName, 0); err != nil {
		t.Errorf("Get() error = %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		got, err := r.Get(qName, 0)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if !reflect.DeepEqual(got, msg) {
			t.Errorf("Get() got = %v, want %v", got, msg)
		}
	}()

	time.Sleep(6 * time.Second)

	if err := r.Put(qName, msg); err != nil {
		t.Errorf("Put2() error = %v", err)
	}

	wg.Wait()
}

func TestCorrectWithTimeout(t *testing.T) {
	r := NewInMemoryQueueRepository(5, 5)
	qName := "Name2"
	msg := core.Message("mmasdas")

	// Создадим очередь и сразу опустошим
	if err := r.Put(qName, msg); err != nil {
		t.Errorf("Put() error = %v", err)
	}
	if _, err := r.Get(qName, 0); err != nil {
		t.Errorf("Get() error = %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		got, err := r.Get(qName, 4)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if !reflect.DeepEqual(got, msg) {
			t.Errorf("Get() got = %v, want %v", got, msg)
		}
	}()

	time.Sleep(2 * time.Second)

	if err := r.Put(qName, msg); err != nil {
		t.Errorf("Put2() error = %v", err)
	}

	wg.Wait()
}
