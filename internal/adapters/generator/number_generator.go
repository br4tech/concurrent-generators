package generator

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/br4tech/concurrent-generators/internal/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("number-generator")

// GenerateOrders inicia múltiplos geradores que enviam pedidos para um canal.
// Esta é a nossa fábrica de geradores simultâneos.
func GenerateOrders(ctx context.Context, numGenerators int, numOrdersPerGenerator int) <-chan domain.Order {
	orders := make(chan domain.Order)
	var wg sync.WaitGroup

	// Inicia um span para a fase de geração de pedidos
	genCtx, span := tracer.Start(ctx, "GenerateAllOrders")
	span.SetAttributes(
		attribute.Int("num.generators", numGenerators),
		attribute.Int("orders.per.generator", numOrdersPerGenerator),
	)
	defer span.End()

	for i := 0; i < numGenerators; i++ {
		wg.Add(1)
		go func(generatorID int) {
			defer wg.Done()
			log.Printf("Gerador %d iniciado", generatorID)

			for j := 0; j < numOrdersPerGenerator; j++ {
				order := domain.Order{
					ID:    (generatorID * numOrdersPerGenerator) + j,
					Value: j,
				}

				// Criamos um span para cada pedido gerado
				_, orderSpan := tracer.Start(genCtx, "CreateOrder")
				orderSpan.SetAttributes(
					attribute.Int("generator.id", generatorID),
					attribute.Int("order.id", order.ID),
				)

				select {
				case orders <- order:
					// Simula um trabalho de geração
					time.Sleep(10 * time.Millisecond)
					orderSpan.End()
				case <-ctx.Done():
					log.Printf("Gerador %d recebendo sinal de cancelamento.", generatorID)
					orderSpan.RecordError(ctx.Err())
					orderSpan.End()
					return
				}
			}
		}(i + 1)
	}

	// Inicia uma goroutine para fechar o canal 'orders' quando todos os geradores terminarem.
	// Isso é crucial para que os workers saibam quando parar.
	go func() {
		wg.Wait()
		close(orders)
		log.Println("Todos os geradores finalizaram. Canal de pedidos fechado.")
	}()

	return orders
}
