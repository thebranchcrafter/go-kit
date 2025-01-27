package application_command

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/application"
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/logger"
	"reflect"
	"sync"
)

type Bus interface {
	RegisterCommand(c application.Command, handler CommandHandler) error
	Dispatch(ctx context.Context, c application.Command) error
	DispatchAsync(ctx context.Context, c application.Command) error
	ProcessFailed(ctx context.Context)
}

type CommandBus struct {
	handlers       map[string]CommandHandler
	lock           sync.Mutex
	l              logger.Logger
	failedCommands chan *FailedCommand
}

func InitCommandBus(l logger.Logger) *CommandBus {
	return &CommandBus{
		handlers:       make(map[string]CommandHandler, 0),
		lock:           sync.Mutex{},
		l:              l,
		failedCommands: make(chan *FailedCommand),
	}
}

type FailedCommand struct {
	command        application.Dto
	handler        CommandHandler
	timesProcessed int
}

type CommandAlreadyRegistered struct {
	message     string
	commandName string
}

func (i CommandAlreadyRegistered) Error() string {
	return i.message
}

func NewCommandAlreadyRegistered(message string, commandName string) CommandAlreadyRegistered {
	return CommandAlreadyRegistered{message: message, commandName: commandName}
}

type CommandNotRegistered struct {
	message     string
	commandName string
}

func (i CommandNotRegistered) Error() string {
	return i.message
}

func NewCommandNotRegistered(message string, commandName string) CommandNotRegistered {
	return CommandNotRegistered{message: message, commandName: commandName}
}

func (bus *CommandBus) RegisterCommand(c application.Command, handler CommandHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	commandName, err := bus.commandName(c)
	if err != nil {
		return err
	}

	if _, ok := bus.handlers[*commandName]; ok {
		return NewCommandAlreadyRegistered("Command already registered", *commandName)
	}

	bus.handlers[*commandName] = handler

	return nil
}

func (bus *CommandBus) Dispatch(ctx context.Context, c application.Command) error {
	commandName, err := bus.commandName(c)
	if err != nil {
		return err
	}

	if handler, ok := bus.handlers[*commandName]; ok {
		err := bus.doHandle(ctx, handler, c)
		if err != nil {
			return err
		}

		return nil
	}

	return NewCommandNotRegistered("Command not registered", *commandName)
}

func (bus *CommandBus) DispatchAsync(ctx context.Context, c application.Command) error {
	commandName, err := bus.commandName(c)
	if err != nil {
		return err
	}

	if handler, ok := bus.handlers[*commandName]; ok {
		go bus.doHandleAsync(ctx, handler, c)

		return nil
	}

	return NewCommandNotRegistered("Command not registered", *commandName)
}

func (bus *CommandBus) doHandle(ctx context.Context, handler CommandHandler, c application.Command) error {
	return handler.Handle(ctx, c)
}

func (bus *CommandBus) doHandleAsync(ctx context.Context, handler CommandHandler, command application.Dto) {
	err := bus.doHandle(ctx, handler, command)

	if err != nil {
		bus.failedCommands <- &FailedCommand{
			command:        command,
			handler:        handler,
			timesProcessed: 1,
		}
		bus.l.Error(ctx, "error_message", map[string]interface{}{"error": err.Error()})
	}
}

func (bus *CommandBus) commandName(cmd interface{}) (*string, error) {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return nil, CommandNotValid{"only pointer to commands are allowed"}
	}

	name := value.String()

	return &name, nil
}

func (bus *CommandBus) ProcessFailed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(bus.failedCommands)
			bus.l.Warn(ctx, "exiting safely failed commands consumer", map[string]interface{}{"error": ctx.Err().Error()})
			return
		case failedCommand := <-bus.failedCommands:
			if failedCommand.timesProcessed >= 3 {
				continue
			}

			failedCommand.timesProcessed++
			if err := bus.doHandle(ctx, failedCommand.handler, failedCommand.command); err != nil {
				bus.l.Warn(ctx, "failing processing command", map[string]interface{}{"error": ctx.Err().Error()})
			}
		}
	}
}

type CommandNotValid struct {
	message string
}

func (i CommandNotValid) Error() string {
	return i.message
}
