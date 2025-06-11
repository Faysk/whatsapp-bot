package utils

import (
	"io"
	"log"
	"os"
)

// SetupLogger configura o logger padrão com saída para terminal e flags de debug
func SetupLogger() {
	log.SetOutput(getLogOutput())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("📝 Logger inicializado.")
}

// getLogOutput retorna o destino padrão do log (stdout ou arquivo futuramente)
func getLogOutput() io.Writer {
	// 🔄 Futuro: salvar em arquivo de log, como logs/bot.log
	// f, err := os.OpenFile("logs/bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	//     return os.Stdout
	// }
	// return io.MultiWriter(os.Stdout, f)

	return os.Stdout
}
