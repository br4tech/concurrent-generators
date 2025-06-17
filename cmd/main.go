package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/br4tech/concurrent-generators/internal/adapters/generator"
	"github.com/br4tech/concurrent-generators/internal/adapters/logger"
	"github.com/br4tech/concurrent-generators/internal/core/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// initTracer inicializa o provedor de trace OpenTelemetry.
func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("concurrent-generator-service"),
		semconv.ServiceVersion("v0.1.0"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func main() {
	// --- Configuração da Observabilidade (Tracing) ---
	tp, err := initTracer()
	if err != nil {
		log.Fatal("Falha ao inicializar o tracer:", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Erro ao desligar o tracer: %v", err)
		}
	}()

	// --- Configuração do Contexto para Cancelamento ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura de sinais do sistema para um desligamento gracioso (graceful shutdown)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopChan
		log.Println("Recebido sinal de interrupção, iniciando graceful shutdown...")
		cancel()
	}()

	// --- Injeção de Dependência ---
	// Criamos nosso adaptador de log...
	stdoutLogger := logger.NewStdoutLogger()
	// ...e o injetamos no nosso serviço do core.
	processingSvc := services.NewProcessingService(stdoutLogger)

	// --- Pipeline de Execução ---
	log.Println("Iniciando a aplicação...")

	// 1. Inicia os geradores simultâneos.
	// O canal 'orders' receberá os dados de todas as goroutines geradoras.
	const numGenerators = 5
	const numOrdersPerGenerator = 10
	ordersChan := generator.GenerateOrders(ctx, numGenerators, numOrdersPerGenerator)

	// 2. Inicia os workers para processar os pedidos do canal.
	// Usamos um WaitGroup para garantir que a main() espere os workers finalizarem.
	var wg sync.WaitGroup
	const numWorkers = 3
	processingSvc.StartWorkers(ctx, numWorkers, ordersChan, &wg)

	// 3. Aguarda a finalização de todos os workers.
	wg.Wait()

	log.Println("Aplicação finalizada com sucesso.")
}
