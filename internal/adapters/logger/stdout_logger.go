package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/br4tech/concurrent-generators/internal/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("stdout-logger")

// StdoutLogger é um adaptador que processa um pedido imprimindo-o no stdout.
type StdoutLogger struct{}

// NewStdoutLogger cria um novo logger.
func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{}
}

// Process implementa a interface ports.OrderProcessor.
func (l *StdoutLogger) Process(ctx context.Context, order domain.Order) error {
	_, span := tracer.Start(ctx, "LogToStdout")
	defer span.End()

	// Simula um trabalho de processamento/IO
	time.Sleep(50 * time.Millisecond)

	logMessage := fmt.Sprintf("Pedido processado com sucesso: ID %d, Valor %d", order.ID, order.Value)
	log.Println(logMessage)

	span.SetAttributes(attribute.String("log.message", logMessage))

	return nil
}
