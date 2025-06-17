package services

import (
	"context"
	"log"
	"sync"

	"github.com/br4tech/concurrent-generators/internal/core/domain"
	"github.com/br4tech/concurrent-generators/internal/core/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("processing-service")

// ProcessingService implementa a lógica de processamento de pedidos.
type ProcessingService struct {
	logger ports.OrderProcessor // Adaptador de logging
}

// NewProcessingService cria uma nova instância do serviço de processamento.
func NewProcessingService(logger ports.OrderProcessor) *ProcessingService {
	return &ProcessingService{
		logger: logger,
	}
}

// StartWorkers inicia os workers que irão processar os pedidos do canal.
// Esta função demonstra o princípio de "Single Responsibility" (S do SOLID).
func (s *ProcessingService) StartWorkers(ctx context.Context, numWorkers int, orders <-chan domain.Order, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.worker(ctx, i+1, wg, orders)
	}
}

// worker é a função que executa em uma goroutine para processar pedidos.
func (s *ProcessingService) worker(ctx context.Context, id int, wg *sync.WaitGroup, orders <-chan domain.Order) {
	defer wg.Done()
	log.Printf("Worker %d iniciado", id)

	for {
		select {
		case order, ok := <-orders:
			if !ok {
				log.Printf("Worker %d finalizando.", id)
				return // Canal fechado
			}

			// Inicia um novo span de trace para o processamento do pedido
			workerCtx, span := tracer.Start(ctx, "ProcessOrder", trace.WithAttributes(
				attribute.Int("worker.id", id),
				attribute.Int("order.id", order.ID),
			))

			log.Printf("Worker %d recebeu o pedido: %d", id, order.ID)
			if err := s.logger.Process(workerCtx, order); err != nil {
				log.Printf("Worker %d: Erro ao processar o pedido %d: %v", id, order.ID, err)
				span.RecordError(err)
			}
			span.End()

		case <-ctx.Done():
			log.Printf("Worker %d recebendo sinal de cancelamento.", id)
			return // Contexto cancelado
		}
	}
}
